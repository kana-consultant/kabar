package model_family

import (
	"context"
	"log"
	"seo-backend/internal/domain/ai_model"
	model_family "seo-backend/internal/domain/modelfamily"
	"seo-backend/internal/domain/paginate"
	"seo-backend/internal/models"
)

type serviceImpl struct {
	repo     model_family.Repository
	ai_model ai_model.Repository
}

func NewService(repo model_family.Repository) model_family.Service {
	return &serviceImpl{
		repo: repo,
	}
}

func (s *serviceImpl) Create(ctx context.Context, req *model_family.CreateRequest) (*model_family.Response, error) {
	// Validate request
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Check if family already exists
	exists, err := s.repo.Exists(ctx, req.ProviderID, req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, model_family.ErrDuplicate
	}

	// Create entity with new fields
	entity := model_family.NewModelFamily(
		req.ProviderID,
		req.SchemaID,
		req.Name,
		req.DisplayName,
		req.Description,
		req.MaxTokens,
		req.Temperature,
		req.SystemPrompt,
	)

	if err := s.repo.Create(ctx, entity); err != nil {
		return nil, err
	}

	return model_family.ToResponse(entity), nil
}

func (s *serviceImpl) GetByID(ctx context.Context, id string) (*model_family.Response, error) {
	if id == "" {
		return nil, model_family.ErrInvalidID
	}

	entity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return model_family.ToResponse(entity), nil
}

func (s *serviceImpl) GetByProviderAndName(ctx context.Context, providerID string, name string) (*model_family.Response, error) {
	if providerID == "" {
		return nil, model_family.ErrInvalidProviderID
	}
	if name == "" {
		return nil, model_family.ErrInvalidName
	}

	entity, err := s.repo.GetByProviderAndName(ctx, providerID, name)
	if err != nil {
		return nil, err
	}

	return model_family.ToResponse(entity), nil
}

func (s *serviceImpl) GetAll(ctx context.Context, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[model_family.Response], error) {
	// Validate pagination parameters
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	result, err := s.repo.GetAll(ctx, userCtx, params)
	if err != nil {
		return nil, err
	}

	// Convert ModelWithStatus ke Response model_family
	responses := make([]model_family.Response, len(result))
	for i, item := range result {
		responses[i] = model_family.Response{
			ID:           item.ID,
			Name:         item.Name,
			DisplayName:  item.DisplayName,
			ProviderID:   item.ProviderID,
			SchemaID:     item.SchemaID,
			Description:  item.Description,
			MaxTokens:    item.MaxTokens,
			Temperature:  item.Temperature,
			SystemPrompt: item.SystemPrompt,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		}
	}

	return &paginate.PaginatedResult[model_family.Response]{
		Data: responses,
	}, nil
}

func (s *serviceImpl) GetByProvider(ctx context.Context, providerID string) ([]model_family.Response, error) {
	if providerID == "" {
		return nil, model_family.ErrInvalidProviderID
	}

	entities, err := s.repo.GetByProviderID(ctx, providerID)
	if err != nil {
		return nil, err
	}

	return model_family.ToResponseListFromSchema(entities), nil
}

func (s *serviceImpl) GetBySchema(ctx context.Context, schemaID string) ([]model_family.Response, error) {
	if schemaID == "" {
		return nil, model_family.ErrInvalidSchemaID
	}

	entities, err := s.repo.GetBySchema(ctx, schemaID)
	if err != nil {
		return nil, err
	}

	return model_family.ToResponseList(entities), nil
}

func (s *serviceImpl) Update(ctx context.Context, id string, req *model_family.UpdateRequest) (*model_family.Response, error) {
	log.Printf("[ModelFamily Update] START | id=%s", id)
	log.Printf("[ModelFamily Update] Request payload: SchemaID=%v, Name=%v, DisplayName=%v, Description=%v, MaxTokens=%v, Temperature=%v, SystemPrompt=%v",
		req.SchemaID, req.Name, req.DisplayName, req.Description, req.MaxTokens, req.Temperature, req.SystemPrompt)

	if id == "" {
		log.Printf("[ModelFamily Update] ERROR: empty ID provided")
		return nil, model_family.ErrInvalidID
	}

	// Get existing entity
	log.Printf("[ModelFamily Update] Fetching existing family with id=%s", id)
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("[ModelFamily Update] ERROR: failed to get existing family | id=%s | err=%v", id, err)
		return nil, err
	}
	log.Printf("[ModelFamily Update] Existing family found | id=%s | name=%s | display_name=%s | schema_id=%v | max_tokens=%d | temperature=%.2f",
		existing.ID, existing.Name, existing.DisplayName, existing.SchemaID, existing.MaxTokens, existing.Temperature)

	// Update fields if provided
	if req.SchemaID != nil {
		log.Printf("[ModelFamily Update] Updating SchemaID: %v -> %v", existing.SchemaID, *req.SchemaID)
		existing.SchemaID = *req.SchemaID
	}

	if req.Name != nil {
		log.Printf("[ModelFamily Update] Checking name conflict | provider_id=%s | new_name=%s | old_name=%s",
			existing.ProviderID, *req.Name, existing.Name)

		// Check if new name conflicts
		exists, err := s.repo.Exists(ctx, existing.ProviderID, *req.Name)
		if err != nil {
			log.Printf("[ModelFamily Update] ERROR: failed to check name existence | provider_id=%s | name=%s | err=%v",
				existing.ProviderID, *req.Name, err)
			return nil, err
		}

		if exists && existing.Name != *req.Name {
			log.Printf("[ModelFamily Update] ERROR: duplicate name | provider_id=%s | name=%s",
				existing.ProviderID, *req.Name)
			return nil, model_family.ErrDuplicate
		}

		log.Printf("[ModelFamily Update] Updating Name: %s -> %s", existing.Name, *req.Name)
		existing.Name = *req.Name
	}

	if req.DisplayName != nil {
		log.Printf("[ModelFamily Update] Updating DisplayName: %s -> %s", existing.DisplayName, *req.DisplayName)
		existing.DisplayName = *req.DisplayName
	}

	if req.Description != nil {
		log.Printf("[ModelFamily Update] Updating Description: %v -> %v", existing.Description, *req.Description)
		existing.Description = req.Description
	}

	if req.MaxTokens != nil {
		log.Printf("[ModelFamily Update] Updating MaxTokens: %d -> %d", existing.MaxTokens, *req.MaxTokens)
		existing.MaxTokens = *req.MaxTokens
	}

	if req.Temperature != nil {
		log.Printf("[ModelFamily Update] Updating Temperature: %.2f -> %.2f", existing.Temperature, *req.Temperature)
		existing.Temperature = *req.Temperature
	}

	if req.SystemPrompt != nil {
		log.Printf("[ModelFamily Update] Updating SystemPrompt: %s -> %s", existing.SystemPrompt, *req.SystemPrompt)
		existing.SystemPrompt = *req.SystemPrompt
	}

	// Save to database
	log.Printf("[ModelFamily Update] Saving updated family to database | id=%s", id)
	if err := s.repo.Update(ctx, existing); err != nil {
		log.Printf("[ModelFamily Update] ERROR: failed to update family | id=%s | err=%v", id, err)
		return nil, err
	}

	log.Printf("[ModelFamily Update] SUCCESS | id=%s | name=%s | display_name=%s | max_tokens=%d | temperature=%.2f",
		existing.ID, existing.Name, existing.DisplayName, existing.MaxTokens, existing.Temperature)

	return model_family.ToResponse(existing), nil
}

func (s *serviceImpl) Delete(ctx context.Context, id string) error {
	if id == "" {
		return model_family.ErrInvalidID
	}
	return s.repo.Delete(ctx, id)
}

func (s *serviceImpl) Validate(ctx context.Context, mf *model_family.ModelFamily) error {
	if mf.ProviderID == "" {
		return model_family.ErrInvalidProviderID
	}
	if mf.SchemaID == "" {
		return model_family.ErrInvalidSchemaID
	}
	if mf.Name == "" {
		return model_family.ErrInvalidName
	}
	if mf.DisplayName == "" {
		return model_family.ErrInvalidDisplayName
	}

	return nil
}

// Private helper methods
func (s *serviceImpl) validateCreateRequest(req *model_family.CreateRequest) error {
	if req.ProviderID == "" {
		return model_family.ErrInvalidProviderID
	}
	if req.SchemaID == "" {
		return model_family.ErrInvalidSchemaID
	}
	if req.Name == "" {
		return model_family.ErrInvalidName
	}
	if req.DisplayName == "" {
		return model_family.ErrInvalidDisplayName
	}

	return nil
}
