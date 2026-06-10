package provider

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	model_family "seo-backend/internal/domain/modelfamily"
	"seo-backend/internal/domain/paginate"
	"seo-backend/internal/domain/provider"
	"seo-backend/internal/domain/request_schema"
	"seo-backend/internal/models"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type ServiceImpl struct {
	repo       provider.Repository
	familyRepo model_family.Repository
	schemaRepo request_schema.Repository
	db         *sql.DB
	redis      *redis.Client
}

func NewService(db *sql.DB, repo provider.Repository, familyRepo model_family.Repository, schemaRepo request_schema.Repository, redis *redis.Client) provider.Service {
	return &ServiceImpl{
		db:         db,
		repo:       repo,
		familyRepo: familyRepo,
		schemaRepo: schemaRepo,
		redis:      redis,
	}
}

// Create creates a new API provider with optional model families
func (s *ServiceImpl) Create(ctx context.Context, req *provider.CreateRequest, userCtx models.UserContext) (*provider.Response, error) {
	if !s.isAdmin(userCtx.GetRole()) {
		return nil, errors.New("access denied: admin role required")
	}

	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	if _, err := url.ParseRequestURI(req.BaseURL); err != nil {
		return nil, provider.ErrInvalidBaseURL
	}

	if req.Families != nil {
		if err := s.validateFamilyForCreate(req.Families); err != nil {
			return nil, fmt.Errorf("invalid family data: %w", err)
		}
	}

	exists, err := s.repo.Exists(ctx, req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check provider existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("provider with name '%s' already exists", req.Name)
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	authType := "bearer"
	if req.AuthType != nil {
		authType = *req.AuthType
	}

	authHeader := "Authorization"
	if req.AuthHeader != nil {
		authHeader = *req.AuthHeader
	}

	authPrefix := "Bearer"
	if req.AuthPrefix != nil {
		authPrefix = *req.AuthPrefix
	}

	defaultHeaders := json.RawMessage("{}")
	if req.DefaultHeaders != nil {
		defaultHeaders = req.DefaultHeaders
	}

	entity := &provider.APIProvider{
		Name:           req.Name,
		DisplayName:    req.DisplayName,
		Description:    req.Description,
		BaseURL:        req.BaseURL,
		AuthType:       &authType,
		AuthHeader:     &authHeader,
		AuthPrefix:     &authPrefix,
		DefaultHeaders: defaultHeaders,
		IsActive:       &isActive,
	}

	if err := s.repo.CreateWithTx(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}

	if req.Families != nil {
		if err := s.createModelFamiliesWithTx(ctx, tx, entity.ID.String(), req.Families); err != nil {
			return nil, fmt.Errorf("failed to create model families: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	go s.clearProviderCache(userCtx.GetTeamID())

	return provider.ToResponse(entity), nil
}

func (s *ServiceImpl) createModelFamiliesWithTx(ctx context.Context, tx *sql.Tx, providerID string, families []model_family.ModelFamilyWithSchema) error {
	if len(families) == 0 {
		return nil
	}

	// STEP 1: Set provider ID + create schema dulu, semua dalam TX yang sama
	for i := range families {
		families[i].ProviderID = providerID

		if families[i].Schema.Name == "" {
			families[i].Schema.Name = fmt.Sprintf("%s_schema_%d", families[i].Name, time.Now().Unix())
		}

		log.Printf("Creating schema for family: %s", families[i].Name)
		schemaID, err := s.createSchemaWithTx(ctx, tx, providerID, &families[i].Schema)
		if err != nil {
			return fmt.Errorf("failed to create schema for family '%s': %w", families[i].Name, err)
		}

		families[i].SchemaID = schemaID
		log.Printf("Schema created with ID: %s for family: %s", schemaID, families[i].Name)

		if families[i].SchemaID == "" {
			return fmt.Errorf("schema ID is empty for family '%s'", families[i].Name)
		}
	}

	// STEP 2: Baru create family, semua dalam TX yang sama
	for i := range families {
		log.Printf("Creating family '%s' with schema_id: %s", families[i].Name, families[i].SchemaID)

		createFamily := &model_family.ModelFamily{
			ProviderID:   providerID,
			SchemaID:     families[i].SchemaID,
			Name:         families[i].Name,
			DisplayName:  families[i].DisplayName,
			Description:  families[i].Description,
			MaxTokens:    families[i].MaxTokens,
			Temperature:  families[i].Temperature,
			SystemPrompt: families[i].SystemPrompt,
		}

		if err := s.familyRepo.CreateWithTx(ctx, tx, createFamily); err != nil {
			return fmt.Errorf("failed to create family '%s': %w", families[i].Name, err)
		}

		families[i].ID = createFamily.ID
		log.Printf("Family created: %s (id: %s)", families[i].Name, families[i].ID)
	}

	return nil
}

// createSchemaWithTx creates a new request schema WITH transaction
func (s *ServiceImpl) createSchemaWithTx(ctx context.Context, tx *sql.Tx, providerID string, schema *request_schema.RequestSchema) (string, error) {
	now := time.Now()

	log.Printf("[createSchemaWithTx] Start - ProviderID: %s", providerID)

	schema.ID = uuid.New().String()
	schema.ProviderID = providerID
	schema.CreatedAt = now
	schema.UpdatedAt = now

	// Pastikan field-field wajib tidak kosong
	if schema.Name == "" {
		schema.Name = fmt.Sprintf("schema_%d", now.Unix())
	}

	if schema.EndpointPath == "" {
		schema.EndpointPath = "/v1/chat/completions"
	}

	if schema.RequestTemplate == nil {
		defaultTemplate := "{\"model\":\"{model}\",\"messages\":\"{prompt}\"}"
		schema.RequestTemplate = &defaultTemplate
	}

	if schema.ResponseTextPath == nil {
		defaultPath := "choices[0].message.content"
		schema.ResponseTextPath = &defaultPath
	}

	log.Printf(
		"[createSchemaWithTx] Schema Data => ID=%s ProviderID=%s Name=%s Endpoint=%s RequestTemplate=%v ResponseTextPath=%v",
		schema.ID,
		schema.ProviderID,
		schema.Name,
		schema.EndpointPath,
		schema.RequestTemplate,
		schema.ResponseTextPath,
	)

	if err := s.schemaRepo.CreateWithTx(ctx, tx, schema); err != nil {
		log.Printf("[createSchemaWithTx] ERROR CreateWithTx: %v", err)
		return "", err
	}

	log.Printf("[createSchemaWithTx] SUCCESS - SchemaID: %s", schema.ID)

	return schema.ID, nil
}

// Update updates an existing provider
func (s *ServiceImpl) Update(ctx context.Context, id string, req *provider.UpdateRequest, userCtx models.UserContext) (*provider.APIProvider, error) {
	if !s.isAdmin(userCtx.GetRole()) {
		return nil, errors.New("access denied: admin role required")
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}
	if existing == nil {
		return nil, provider.ErrNotFound
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update provider fields
	updates := map[string]interface{}{
		"name":            nil,
		"display_name":    nil,
		"description":     nil,
		"base_url":        nil,
		"auth_type":       nil,
		"auth_header":     nil,
		"auth_prefix":     nil,
		"default_headers": nil,
		"is_active":       nil,
	}

	if req.Name != nil {
		updates["name"] = *req.Name
	}

	if req.DisplayName != nil {
		updates["display_name"] = *req.DisplayName
	}

	if req.Description != nil {
		updates["description"] = *req.Description
	}

	if req.BaseURL != nil {
		if _, err := url.ParseRequestURI(*req.BaseURL); err != nil {
			return nil, provider.ErrInvalidBaseURL
		}
		updates["base_url"] = *req.BaseURL
	}

	if req.AuthType != nil {
		updates["auth_type"] = *req.AuthType
	}

	if req.AuthHeader != nil {
		updates["auth_header"] = *req.AuthHeader
	}

	if req.AuthPrefix != nil {
		updates["auth_prefix"] = *req.AuthPrefix
	}

	if req.DefaultHeaders != nil {
		updates["default_headers"] = req.DefaultHeaders
	}

	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) > 0 {
		if err := s.repo.UpdateWithTx(ctx, tx, id, updates); err != nil {
			return nil, fmt.Errorf("failed to update provider: %w", err)
		}
	}

	// Create or update model families (including delete removed ones)
	if req.Families != nil {
		if err := s.CreateOrUpdateModelFamilies(ctx, tx, id, req.Families); err != nil {
			return nil, fmt.Errorf("failed to create/update model families: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	go s.clearProviderCache(userCtx.GetTeamID())

	return s.repo.GetByID(ctx, id)
}

// CreateOrUpdateModelFamilies creates, updates, or deletes model families for a provider.
func (s *ServiceImpl) CreateOrUpdateModelFamilies(ctx context.Context, tx *sql.Tx, providerID string, families []model_family.ModelFamilyWithSchema) error {
	// STEP 1: Set provider ID
	for i := range families {
		families[i].ProviderID = providerID
	}

	// STEP 2: Fetch semua existing families dulu SEBELUM upsert
	existingFamilies, err := s.familyRepo.GetByProviderID(ctx, providerID)
	if err != nil {
		return fmt.Errorf("failed to get existing families for provider '%s': %w", providerID, err)
	}

	// Build map existing families by name untuk lookup cepat
	existingByName := make(map[string]*model_family.ModelFamilyWithSchema, len(existingFamilies))
	for i, f := range existingFamilies {
		existingByName[f.Name] = &existingFamilies[i]
	}

	// STEP 3: Handle schema
	for i := range families {
		existing, alreadyExists := existingByName[families[i].Name]

		if alreadyExists && families[i].Schema.Name == "" {
			families[i].SchemaID = existing.SchemaID
			log.Printf("Reusing existing schema for family: %s, schemaID: %s", families[i].Name, families[i].SchemaID)
			continue
		}

		if families[i].Schema.Name == "" {
			families[i].Schema.Name = fmt.Sprintf("%s_schema_%d", families[i].Name, time.Now().Unix())
			log.Printf("Generated schema name: %s", families[i].Schema.Name)
		}

		log.Printf("Creating new schema for family: %s", families[i].Name)
		schemaID, err := s.createSchemaWithTx(ctx, tx, providerID, &families[i].Schema)
		if err != nil {
			return fmt.Errorf("failed to create schema for family '%s': %w", families[i].Name, err)
		}
		families[i].SchemaID = schemaID
		log.Printf("Created schema ID: %s for family: %s", schemaID, families[i].Name)
	}

	// STEP 4: Upsert families
	requestedNames := make(map[string]bool, len(families))
	for i := range families {
		requestedNames[families[i].Name] = true

		if families[i].SchemaID == "" {
			return fmt.Errorf("family '%s' has no schema_id after schema step", families[i].Name)
		}

		existing, alreadyExists := existingByName[families[i].Name]

		if alreadyExists {
			// Update existing family
			existing.DisplayName = families[i].DisplayName
			existing.Description = families[i].Description
			existing.SchemaID = families[i].SchemaID
			existing.MaxTokens = families[i].MaxTokens
			existing.Temperature = families[i].Temperature
			existing.SystemPrompt = families[i].SystemPrompt

			if err := s.familyRepo.UpdateWithTx(ctx, tx, existing); err != nil {
				return fmt.Errorf("failed to update family '%s': %w", families[i].Name, err)
			}
			families[i].ID = existing.ID
			log.Printf("Updated family: %s (id: %s)", families[i].Name, families[i].ID)
		} else {
			// Create new family
			createFamily := &model_family.ModelFamily{
				ProviderID:   providerID,
				SchemaID:     families[i].SchemaID,
				Name:         families[i].Name,
				DisplayName:  families[i].DisplayName,
				Description:  families[i].Description,
				MaxTokens:    families[i].MaxTokens,
				Temperature:  families[i].Temperature,
				SystemPrompt: families[i].SystemPrompt,
			}

			if err := s.familyRepo.CreateWithTx(ctx, tx, createFamily); err != nil {
				return fmt.Errorf("failed to create family '%s': %w", families[i].Name, err)
			}
			families[i].ID = createFamily.ID
			log.Printf("Created family: %s (id: %s)", families[i].Name, families[i].ID)
		}
	}

	// STEP 5: Delete families yang tidak ada di request
	for _, dbFamily := range existingFamilies {
		if !requestedNames[dbFamily.Name] {
			log.Printf("Deleting removed family: %s (id: %s)", dbFamily.Name, dbFamily.ID)
			if err := s.familyRepo.DeleteWithTx(ctx, tx, dbFamily.ID); err != nil {
				return fmt.Errorf("failed to delete family '%s': %w", dbFamily.Name, err)
			}
		}
	}

	return nil
}

// Delete soft deletes a provider
func (s *ServiceImpl) Delete(ctx context.Context, id string, userCtx models.UserContext) error {
	if !s.isAdmin(userCtx.GetRole()) {
		return errors.New("access denied: admin role required")
	}
	if id == "" {
		return provider.ErrInvalidID
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, provider.ErrNotFound) {
			return errors.New("provider not found")
		}
		return fmt.Errorf("failed to get provider: %w", err)
	}

	usageCount, err := s.repo.CheckProviderUsage(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check provider usage: %w", err)
	}
	if usageCount > 0 {
		return fmt.Errorf("cannot delete provider: still referenced by %d api_key(s)", usageCount)
	}

	// STEP 1: Ambil semua families untuk dapat schema IDs
	existingFamilies, err := s.familyRepo.GetByProviderID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get model families: %w", err)
	}

	// STEP 2: Delete families dulu
	for _, family := range existingFamilies {
		if err := s.familyRepo.DeleteWithTx(ctx, tx, family.ID); err != nil {
			return fmt.Errorf("failed to delete family '%s': %w", family.Name, err)
		}
		log.Printf("Deleted family: %s (id: %s)", family.Name, family.ID)
	}

	// STEP 3: Delete schemas
	for _, family := range existingFamilies {
		if family.SchemaID == "" {
			continue
		}
		if err := s.schemaRepo.DeleteWithTx(ctx, tx, family.SchemaID); err != nil {
			return fmt.Errorf("failed to delete schema for family '%s': %w", family.Name, err)
		}
		log.Printf("Deleted schema: %s for family: %s", family.SchemaID, family.Name)
	}

	// STEP 4: Hard delete provider
	if err := s.repo.HardDeleteWithTx(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete provider: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	go s.clearProviderCache(userCtx.GetTeamID())

	return nil
}

// HardDelete permanently deletes a provider
func (s *ServiceImpl) HardDelete(ctx context.Context, id string, userCtx models.UserContext) error {
	if !s.isAdmin(userCtx.GetRole()) {
		return errors.New("access denied: admin role required")
	}
	if id == "" {
		return provider.ErrInvalidID
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, provider.ErrNotFound) {
			return errors.New("provider not found")
		}
		return fmt.Errorf("failed to get provider: %w", err)
	}

	usageCount, err := s.repo.CheckProviderUsage(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check provider usage: %w", err)
	}
	if usageCount > 0 {
		return fmt.Errorf("cannot delete provider: still referenced by %d model_families or model(s)", usageCount)
	}

	if err := s.repo.HardDeleteWithTx(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to hard delete provider: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	go s.clearProviderCache(userCtx.GetTeamID())

	return nil
}

// GetByID retrieves a provider by ID
func (s *ServiceImpl) GetByID(ctx context.Context, id string, userCtx models.UserContext) (*provider.Response, error) {
	if id == "" {
		return nil, provider.ErrInvalidID
	}

	entity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, provider.ErrNotFound) {
			return nil, errors.New("provider not found")
		}
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	if !s.isAdmin(userCtx.GetRole()) && entity.IsActive != nil && !*entity.IsActive {
		return nil, errors.New("provider not found")
	}

	families, err := s.familyRepo.GetByProviderID(ctx, id)
	if err != nil {
		log.Printf("Warning: failed to get families for provider %s: %v", id, err)
	}
	entity.Families = families

	return provider.ToResponse(entity), nil
}

// GetByName retrieves a provider by name
func (s *ServiceImpl) GetByName(ctx context.Context, name string, userCtx models.UserContext) (*provider.Response, error) {
	if name == "" {
		return nil, provider.ErrInvalidName
	}

	entity, err := s.repo.GetByName(ctx, name)
	if err != nil {
		if errors.Is(err, provider.ErrNotFound) {
			return nil, errors.New("provider not found")
		}
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	if !s.isAdmin(userCtx.GetRole()) && entity.IsActive != nil && !*entity.IsActive {
		return nil, errors.New("provider not found")
	}

	return provider.ToResponse(entity), nil
}

// GetAll retrieves all providers with pagination
func (s *ServiceImpl) GetAll(ctx context.Context, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[provider.Response], error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	cacheKey := s.getCacheKey(userCtx.GetTeamID(), "all", params)
	cached, err := s.getFromCache(ctx, cacheKey)
	if err == nil && cached != nil {
		var result paginate.PaginatedResult[provider.Response]
		if err := json.Unmarshal(cached, &result); err == nil {
			log.Printf("cache hit for key: %s", cacheKey)
			return &result, nil
		}
	}

	result, err := s.repo.GetAll(ctx, userCtx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get providers: %w", err)
	}

	responses := provider.ToResponseList(result.Data)

	paginatedResult := &paginate.PaginatedResult[provider.Response]{
		Data:        responses,
		TotalItems:  result.TotalItems,
		TotalPages:  result.TotalPages,
		CurrentPage: result.CurrentPage,
		Limit:       result.Limit,
		Offset:      result.Offset,
	}

	go s.setToCache(ctx, cacheKey, paginatedResult, 5*time.Minute)

	return paginatedResult, nil
}

// GetActive retrieves all active providers with pagination
func (s *ServiceImpl) GetActive(ctx context.Context, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[provider.Response], error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}
	if params.Offset < 0 {
		params.Offset = 0
	}

	cacheKey := s.getCacheKey(userCtx.GetTeamID(), "active", params)
	cached, err := s.getFromCache(ctx, cacheKey)
	if err == nil && cached != nil {
		var result paginate.PaginatedResult[provider.Response]
		if err := json.Unmarshal(cached, &result); err == nil {
			log.Printf("cache hit for key: %s", cacheKey)
			return &result, nil
		}
	}

	entities, err := s.repo.GetActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active providers: %w", err)
	}

	start := params.Offset
	end := params.Offset + params.Limit
	if start > len(entities) {
		start = len(entities)
	}
	if end > len(entities) {
		end = len(entities)
	}

	paginatedEntities := entities[start:end]
	total := int(len(entities))

	totalPages := int(total) / params.Limit
	if int(total)%params.Limit > 0 {
		totalPages++
	}
	currentPage := (params.Offset / params.Limit) + 1

	responses := provider.ToResponseList(paginatedEntities)

	paginatedResult := &paginate.PaginatedResult[provider.Response]{
		Data:        responses,
		TotalItems:  total,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		Limit:       params.Limit,
		Offset:      params.Offset,
	}

	go s.setToCache(ctx, cacheKey, paginatedResult, 5*time.Minute)

	return paginatedResult, nil
}

// UpdateHeaders updates provider's default headers
func (s *ServiceImpl) UpdateHeaders(ctx context.Context, id string, headers map[string]string, userCtx models.UserContext) error {
	if !s.isAdmin(userCtx.GetRole()) {
		return errors.New("access denied: admin role required")
	}
	if id == "" {
		return provider.ErrInvalidID
	}

	headersJSON, err := json.Marshal(headers)
	if err != nil {
		return fmt.Errorf("failed to marshal headers: %w", err)
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := s.repo.UpdateDefaultHeadersWithTx(ctx, tx, id, headersJSON); err != nil {
		if errors.Is(err, provider.ErrNotFound) {
			return errors.New("provider not found")
		}
		return fmt.Errorf("failed to update headers: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	go s.clearProviderCache(userCtx.GetTeamID())

	return nil
}

// ToggleActive toggles provider active status
func (s *ServiceImpl) ToggleActive(ctx context.Context, id string, userCtx models.UserContext) error {
	if !s.isAdmin(userCtx.GetRole()) {
		return errors.New("access denied: admin role required")
	}
	if id == "" {
		return provider.ErrInvalidID
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	provider_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, provider.ErrNotFound) {
			return errors.New("provider not found")
		}
		return fmt.Errorf("failed to get provider: %w", err)
	}

	currentStatus := true
	if provider_.IsActive != nil {
		currentStatus = *provider_.IsActive
	}

	if err := s.repo.ToggleActiveWithTx(ctx, tx, id, !currentStatus); err != nil {
		return fmt.Errorf("failed to toggle provider status: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	go s.clearProviderCache(userCtx.GetTeamID())

	return nil
}

// Validate validates a provider entity
func (s *ServiceImpl) Validate(ctx context.Context, p *provider.APIProvider) error {
	if p.Name == "" {
		return provider.ErrInvalidName
	}
	if p.BaseURL == "" {
		return provider.ErrInvalidBaseURL
	}
	if _, err := url.ParseRequestURI(p.BaseURL); err != nil {
		return provider.ErrInvalidBaseURL
	}
	return nil
}

// CreateModelFamily creates a single model family for a provider
func (s *ServiceImpl) CreateModelFamily(
	ctx context.Context,
	providerID string,
	families []model_family.ModelFamilyWithSchema,
	userCtx models.UserContext,
) ([]model_family.Response, error) {

	if !s.isAdmin(userCtx.GetRole()) {
		return nil, errors.New("access denied: admin role required")
	}

	if providerID == "" {
		return nil, errors.New("provider ID is required")
	}

	if err := s.validateFamilyForCreate(families); err != nil {
		return nil, err
	}

	_, err := s.repo.GetByID(ctx, providerID)
	if err != nil {
		if errors.Is(err, provider.ErrNotFound) {
			return nil, errors.New("provider not found")
		}
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create schema terlebih dahulu
	for i := range families {
		schemaID, err := s.createSchemaWithTx(
			ctx,
			tx,
			providerID,
			&families[i].Schema,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to create schema for family '%s': %w",
				families[i].Name,
				err,
			)
		}

		families[i].SchemaID = schemaID
		families[i].ProviderID = providerID
	}

	// Baru create family
	if err := s.familyRepo.CreateBatchWithTx(ctx, tx, families); err != nil {
		return nil, fmt.Errorf("failed to create model family: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	go s.clearProviderCache(userCtx.GetTeamID())

	return model_family.ToResponseListSchem(families), nil
}

// UpdateModelFamily updates an existing model family
func (s *ServiceImpl) UpdateModelFamily(ctx context.Context, familyID string, family *model_family.ModelFamily, userCtx models.UserContext) (*model_family.Response, error) {
	if !s.isAdmin(userCtx.GetRole()) {
		return nil, errors.New("access denied: admin role required")
	}
	if familyID == "" {
		return nil, errors.New("family ID is required")
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	existing, err := s.familyRepo.GetByID(ctx, familyID)
	if err != nil {
		if errors.Is(err, model_family.ErrNotFound) {
			return nil, errors.New("model family not found")
		}
		return nil, fmt.Errorf("failed to get model family: %w", err)
	}

	if family.Name != "" {
		existing.Name = family.Name
	}
	if family.DisplayName != "" {
		existing.DisplayName = family.DisplayName
	}
	if family.Description != nil {
		existing.Description = family.Description
	}
	if family.SchemaID != "" {
		existing.SchemaID = family.SchemaID
	}
	if family.MaxTokens > 0 {
		existing.MaxTokens = family.MaxTokens
	}
	if family.Temperature > 0 {
		existing.Temperature = family.Temperature
	}
	if family.SystemPrompt != "" {
		existing.SystemPrompt = family.SystemPrompt
	}

	if err := s.familyRepo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("failed to update model family: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	go s.clearProviderCache(userCtx.GetTeamID())

	return model_family.ToResponse(existing), nil
}

// DeleteModelFamily deletes a model family
func (s *ServiceImpl) DeleteModelFamily(ctx context.Context, familyID string, userCtx models.UserContext) error {
	if !s.isAdmin(userCtx.GetRole()) {
		return errors.New("access denied: admin role required")
	}
	if familyID == "" {
		return errors.New("family ID is required")
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = s.familyRepo.GetByID(ctx, familyID)
	if err != nil {
		if errors.Is(err, model_family.ErrNotFound) {
			return errors.New("model family not found")
		}
		return fmt.Errorf("failed to get model family: %w", err)
	}

	if err := s.familyRepo.Delete(ctx, familyID); err != nil {
		return fmt.Errorf("failed to delete model family: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	go s.clearProviderCache(userCtx.GetTeamID())

	return nil
}

// GetModelFamiliesByProvider retrieves all model families for a provider with pagination
func (s *ServiceImpl) GetModelFamiliesByProvider(ctx context.Context, providerID string, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[model_family.Response], error) {
	if providerID == "" {
		return nil, errors.New("provider ID is required")
	}

	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}
	if params.Offset < 0 {
		params.Offset = 0
	}

	cacheKey := s.getCacheKey(userCtx.GetTeamID(), "families", params, providerID)
	cached, err := s.getFromCache(ctx, cacheKey)
	if err == nil && cached != nil {
		var result paginate.PaginatedResult[model_family.Response]
		if err := json.Unmarshal(cached, &result); err == nil {
			log.Printf("cache hit for key: %s", cacheKey)
			return &result, nil
		}
	}

	_, err = s.repo.GetByID(ctx, providerID)
	if err != nil {
		if errors.Is(err, provider.ErrNotFound) {
			return nil, errors.New("provider not found")
		}
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	allFamilies, err := s.familyRepo.GetByProviderID(ctx, providerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get model families: %w", err)
	}

	start := params.Offset
	end := params.Offset + params.Limit
	if start > len(allFamilies) {
		start = len(allFamilies)
	}
	if end > len(allFamilies) {
		end = len(allFamilies)
	}

	paginatedFamilies := allFamilies[start:end]
	total := int(len(allFamilies))

	totalPages := int(total) / params.Limit
	if int(total)%params.Limit > 0 {
		totalPages++
	}
	currentPage := (params.Offset / params.Limit) + 1

	responses := model_family.ToResponseListSchem(paginatedFamilies)

	paginatedResult := &paginate.PaginatedResult[model_family.Response]{
		Data:        responses,
		TotalItems:  total,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		Limit:       params.Limit,
		Offset:      params.Offset,
	}

	go s.setToCache(ctx, cacheKey, paginatedResult, 5*time.Minute)

	return paginatedResult, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Private Helper Methods
// ─────────────────────────────────────────────────────────────────────────────

func (s *ServiceImpl) validateCreateRequest(req *provider.CreateRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.DisplayName == "" {
		return errors.New("display_name is required")
	}
	if req.BaseURL == "" {
		return errors.New("base_url is required")
	}
	return nil
}

func (s *ServiceImpl) isAdmin(role string) bool {
	return role == "admin" || role == "superadmin"
}

func (s *ServiceImpl) validateFamilyForCreate(families []model_family.ModelFamilyWithSchema) error {
	for i, family := range families {
		if family.Name == "" {
			return fmt.Errorf("family[%d]: name is required", i)
		}
		if family.DisplayName == "" {
			return fmt.Errorf("family[%d]: display_name is required", i)
		}
	}
	return nil
}

func (s *ServiceImpl) buildUpdates(req *provider.UpdateRequest) map[string]interface{} {
	updates := map[string]interface{}{}

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.DisplayName != nil {
		updates["display_name"] = *req.DisplayName
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.BaseURL != nil {
		updates["base_url"] = strings.TrimSuffix(*req.BaseURL, "/")
	}
	if req.AuthType != nil {
		updates["auth_type"] = *req.AuthType
	}
	if req.AuthHeader != nil {
		updates["auth_header"] = *req.AuthHeader
	}
	if req.AuthPrefix != nil {
		updates["auth_prefix"] = *req.AuthPrefix
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.DefaultHeaders != nil {
		updates["default_headers"] = req.DefaultHeaders
	}

	return updates
}

// Redis cache helper methods
func (s *ServiceImpl) getCacheKey(teamID string, cacheType string, params paginate.PaginationParams, extra ...string) string {
	key := fmt.Sprintf("api_providers:%s:%s:limit_%d:offset_%d", teamID, cacheType, params.Limit, params.Offset)

	if len(extra) > 0 && extra[0] != "" {
		key = fmt.Sprintf("%s:provider_%s", key, extra[0])
	}

	return key
}

func (s *ServiceImpl) getFromCache(ctx context.Context, key string) ([]byte, error) {
	if s.redis == nil {
		return nil, errors.New("redis client not available")
	}

	return s.redis.Get(ctx, key).Bytes()
}

func (s *ServiceImpl) setToCache(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if s.redis == nil {
		return errors.New("redis client not available")
	}

	data, err := json.Marshal(value)
	if err != nil {
		log.Printf("failed to marshal cache data: %v", err)
		return err
	}

	return s.redis.Set(ctx, key, data, expiration).Err()
}

func (s *ServiceImpl) clearProviderCache(teamID string) {
	if s.redis == nil {
		log.Printf("redis client not available, skipping cache clear")
		return
	}

	pattern := fmt.Sprintf("api_providers:%s:*", teamID)
	ctx := context.Background()
	iter := s.redis.Scan(ctx, 0, pattern, 0).Iterator()

	var keysToDelete []string
	for iter.Next(ctx) {
		keysToDelete = append(keysToDelete, iter.Val())
	}

	if err := iter.Err(); err != nil {
		log.Printf("failed to scan redis keys: %v", err)
		return
	}

	if len(keysToDelete) > 0 {
		if err := s.redis.Del(ctx, keysToDelete...).Err(); err != nil {
			log.Printf("failed to delete cache keys: %v", err)
		} else {
			log.Printf("cleared %d cache keys for team %s", len(keysToDelete), teamID)
		}
	}
}
