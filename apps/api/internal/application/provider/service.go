package provider

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"seo-backend/internal/domain/provider"
)

// Service handles provider business logic
type Service struct {
	repo provider.Repository
	db   *sql.DB
}

// NewService creates a new provider service
func NewService(db *sql.DB, repo provider.Repository) provider.ProviderService {
	return &Service{
		db:   db,
		repo: repo,
	}
}

// CreateProvider creates a new provider (admin only)
func (s *Service) CreateProvider(ctx context.Context, req provider.CreateProviderRequest, userCtx provider.UserContext) (string, error) {
	// Authorization check
	if !s.isAdmin(userCtx.GetUserRole()) {
		return "", errors.New("access denied: admin role required")
	}

	// Validation
	if err := s.validateCreateRequest(req); err != nil {
		return "", err
	}

	// Check if provider name already exists
	exists, err := s.repo.ExistsByName(ctx, req.Name)
	if err != nil {
		return "", fmt.Errorf("failed to check provider existence: %w", err)
	}
	if exists {
		return "", fmt.Errorf("provider with name '%s' already exists", req.Name)
	}

	// Begin transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Prepare entity
	imageEndpoint := s.emptyToNil(req.ImageEndpoint)
	responseImagePath := s.emptyToNil(req.ResponseImagePath)

	provider := &provider.APIProvider{
		Name:              req.Name,
		DisplayName:       req.DisplayName,
		Description:       req.Description,
		BaseURL:           req.BaseURL,
		AuthType:          req.AuthType,
		AuthHeader:        req.AuthHeader,
		AuthPrefix:        sql.NullString{String: req.AuthPrefix, Valid: req.AuthPrefix != ""},
		TextEndpoint:      req.TextEndpoint,
		ImageEndpoint:     imageEndpoint,
		DefaultHeaders:    req.DefaultHeaders,
		RequestTemplate:   req.RequestTemplate,
		ResponseTextPath:  req.ResponseTextPath,
		ResponseImagePath: responseImagePath,
		IsActive:          true,
	}

	// Create via repository
	id, err := s.repo.Create(ctx, tx, provider)
	if err != nil {
		return "", fmt.Errorf("failed to create provider: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return id, nil
}

// GetProviderByID retrieves a provider by ID
func (s *Service) GetProviderByID(ctx context.Context, id string) (*provider.APIProvider, error) {
	if id == "" {
		return nil, errors.New("provider id is required")
	}

	provider, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	if provider == nil {
		return nil, errors.New("provider not found")
	}

	return provider, nil
}

// GetAllProviders retrieves all providers based on user role
func (s *Service) GetAllProviders(ctx context.Context, userCtx provider.UserContext) ([]provider.APIProvider, error) {
	providers, err := s.repo.GetAll(ctx, userCtx.GetUserRole())
	if err != nil {
		return nil, fmt.Errorf("failed to get providers: %w", err)
	}

	return providers, nil
}

// UpdateProvider updates a provider (admin only)
func (s *Service) UpdateProvider(ctx context.Context, id string, req provider.UpdateProviderRequest, userCtx provider.UserContext) error {
	// Authorization check
	if !s.isAdmin(userCtx.GetUserRole()) {
		return errors.New("access denied: admin role required")
	}

	if id == "" {
		return errors.New("provider id is required")
	}

	// Begin transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if provider exists
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get provider: %w", err)
	}
	if existing == nil {
		return errors.New("provider not found")
	}

	// Build updates map
	updates := s.buildUpdates(req)

	// Update via repository
	if len(updates) > 0 {
		if err := s.repo.Update(ctx, tx, id, updates); err != nil {
			return fmt.Errorf("failed to update provider: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// DeleteProvider deletes a provider (admin only)
func (s *Service) DeleteProvider(ctx context.Context, id string, userCtx provider.UserContext) error {
	// Authorization check
	if !s.isAdmin(userCtx.GetUserRole()) {
		return errors.New("access denied: admin role required")
	}

	if id == "" {
		return errors.New("provider id is required")
	}

	// Begin transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if provider exists
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get provider: %w", err)
	}
	if existing == nil {
		return errors.New("provider not found")
	}

	// Check if provider is used by any models
	usageCount, err := s.repo.CheckProviderUsage(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check provider usage: %w", err)
	}
	if usageCount > 0 {
		return fmt.Errorf("cannot delete provider because it is used by %d model(s)", usageCount)
	}

	// Delete via repository
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete provider: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Helper: Validate create request
func (s *Service) validateCreateRequest(req provider.CreateProviderRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.DisplayName == "" {
		return errors.New("displayName is required")
	}
	if req.BaseURL == "" {
		return errors.New("baseUrl is required")
	}
	if req.AuthType == "" {
		return errors.New("authType is required")
	}
	if req.AuthHeader == "" {
		return errors.New("authHeader is required")
	}
	if req.TextEndpoint == "" {
		return errors.New("textEndpoint is required")
	}
	if req.ResponseTextPath == "" {
		return errors.New("responseTextPath is required")
	}
	return nil
}

// Helper: Build updates map from request
func (s *Service) buildUpdates(req provider.UpdateProviderRequest) map[string]interface{} {
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.DisplayName != nil {
		updates["displayName"] = *req.DisplayName
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.BaseURL != nil {
		updates["baseUrl"] = *req.BaseURL
	}
	if req.AuthType != nil {
		updates["authType"] = *req.AuthType
	}
	if req.AuthHeader != nil {
		updates["authHeader"] = *req.AuthHeader
	}
	if req.AuthPrefix != nil {
		updates["authPrefix"] = *req.AuthPrefix
	}
	if req.TextEndpoint != nil {
		updates["textEndpoint"] = *req.TextEndpoint
	}
	if req.ImageEndpoint != nil {
		updates["imageEndpoint"] = s.emptyToNil(*req.ImageEndpoint)
	}
	if len(req.DefaultHeaders) > 0 {
		updates["defaultHeaders"] = req.DefaultHeaders
	}
	if len(req.RequestTemplate) > 0 {
		updates["requestTemplate"] = req.RequestTemplate
	}
	if req.ResponseTextPath != nil {
		updates["responseTextPath"] = *req.ResponseTextPath
	}
	if req.ResponseImagePath != nil {
		updates["responseImagePath"] = s.emptyToNil(*req.ResponseImagePath)
	}
	if req.IsActive != nil {
		updates["isActive"] = *req.IsActive
	}

	return updates
}

// Helper: Convert empty string to nil pointer
func (s *Service) emptyToNil(str string) *string {
	if str == "" {
		return nil
	}
	return &str
}

// Helper: Check if user is admin
func (s *Service) isAdmin(role string) bool {
	return role == "admin" || role == "super_admin"
}
