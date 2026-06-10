package ai_model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"seo-backend/internal/domain/ai_model"
	model_family "seo-backend/internal/domain/modelfamily"
	"seo-backend/internal/domain/paginate"
	"seo-backend/internal/infrastructure/db/repositories"
	"seo-backend/internal/models"

	"github.com/go-redis/redis/v8"
)

type Service struct {
	repo ai_model.Repository
}

func NewService(db *sql.DB, redisClient *redis.Client) ai_model.Service {
	return &Service{
		repo: repositories.ModelRepository(db, redisClient),
	}
}

func (s *Service) Create(ctx context.Context, req *ai_model.CreateRequest, userCtx models.UserContext) (*ai_model.Response, error) {
	// Business logic validation
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if req.DisplayName == "" {
		return nil, errors.New("display_name is required")
	}

	// Check existing
	exists, err := s.repo.Exists(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ai_model.ErrDuplicate
	}

	// Create entity
	entity := &ai_model.AIModel{
		Name:        req.Name,
		DisplayName: &req.DisplayName,
		Description: req.Description,
		ProviderID:  req.ProviderID,
		FamilyID:    req.FamilyID,
		// SchemaID:      req.SchemaID, // <-- DIHAPUS
		TeamID:        req.TeamID,
		SystemPrompt:  req.SystemPrompt,
		MaxTokens:     req.MaxTokens,
		Temperature:   req.Temperature,
		ContextWindow: req.ContextWindow,
		IsActive:      req.IsActive,
		IsDefault:     req.IsDefault,
		CreatedBy:     req.CreatedBy,
	}

	err = s.repo.Create(ctx, entity)
	if err != nil {
		return nil, err
	}

	// If set as default, update provider defaults
	if req.IsDefault != nil && *req.IsDefault {
		providerID := ""
		if req.ProviderID != nil {
			providerID = *req.ProviderID
		}
		if providerID != "" {
			err = s.repo.SetDefaultForProvider(ctx, providerID, entity.ID)
			if err != nil {
				return nil, err
			}
		}
	}

	return ai_model.ToResponse(entity), nil
}

func (s *Service) Update(ctx context.Context, id string, req *ai_model.UpdateRequest, userCtx models.UserContext) (*ai_model.Response, error) {
	log.Print("[AIModel Update] START")
	log.Print("ID:", id)
	log.Print("User Role:", userCtx.GetRole())
	log.Print("Team ID:", userCtx.GetTeamID())

	// Get existing
	log.Print("Fetching existing model...")
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Print("ERROR failed to get existing model:", err)
		return nil, err
	}
	log.Print("Existing model found:", existing.ID)
	log.Print("Existing Family ID:", getStringValue(existing.FamilyID))

	// Update fields
	if req.FamilyID != nil {
		log.Print("Updating FamilyID:", *req.FamilyID)
		existing.FamilyID = req.FamilyID
	}
	if req.Name != nil {
		log.Print("Updating Name:", *req.Name)
		existing.Name = *req.Name
	}
	if req.DisplayName != nil {
		log.Print("Updating DisplayName:", *req.DisplayName)
		existing.DisplayName = req.DisplayName
	}
	if req.Description != nil {
		log.Print("Updating Description:", *req.Description)
		existing.Description = req.Description
	}
	if req.SystemPrompt != nil {
		log.Print("Updating SystemPrompt:", *req.SystemPrompt)
		existing.SystemPrompt = req.SystemPrompt
	}
	if req.MaxTokens != nil {
		log.Print("Updating MaxTokens:", *req.MaxTokens)
		existing.MaxTokens = req.MaxTokens
	}
	if req.Temperature != nil {
		log.Print("Updating Temperature:", *req.Temperature)
		existing.Temperature = req.Temperature
	}
	if req.ContextWindow != nil {
		log.Print("Updating ContextWindow:", *req.ContextWindow)
		existing.ContextWindow = req.ContextWindow
	}
	if req.IsActive != nil {
		log.Print("Updating IsActive:", *req.IsActive)
		existing.IsActive = req.IsActive
	}
	if req.IsDefault != nil {
		log.Print("Updating IsDefault:", *req.IsDefault)
		existing.IsDefault = req.IsDefault
	}
	if req.TeamID != nil {
		log.Print("Updating TeamID:", *req.TeamID)
		existing.TeamID = req.TeamID
	}

	// Save to database
	log.Print("Saving to database...")
	err = s.repo.Update(ctx, existing)
	if err != nil {
		log.Print("ERROR failed to update:", err)
		return nil, err
	}
	log.Print("Database update successful")

	// If set as default, update provider defaults
	if req.IsDefault != nil && *req.IsDefault {
		log.Print("Setting as default model...")
		providerID := ""
		if existing.ProviderID != nil {
			providerID = *existing.ProviderID
		}

		if providerID != "" {
			err = s.repo.SetDefaultForProvider(ctx, providerID, id)
			if err != nil {
				log.Print("ERROR failed to set default:", err)
				return nil, err
			}
			log.Print("Default model set successfully")
		}
	}

	log.Print("[AIModel Update] SUCCESS")
	return ai_model.ToResponse(existing), nil
}

// Helper function
func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func (s *Service) Delete(ctx context.Context, id string, userCtx models.UserContext) error {
	return s.repo.Delete(ctx, id)
}

// Untuk mendapatkan response AIModel saja
func (s *Service) GetByID(ctx context.Context, id string, userCtx models.UserContext) (*ai_model.Response, error) {
	entity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return ai_model.ToResponse(entity), nil
}

func (s *Service) GetSchemaByID(ctx context.Context, id string, userCtx models.UserContext) (*model_family.ModelFamilyWithSchema, error) {
	log.Printf("Getting model family with schema by ID: %s", id)

	// Panggil repository
	entity, err := s.repo.GetByIDWithSchema(ctx, id)
	if err != nil {
		log.Printf("Failed to get model family from repository: %v", err)
		return nil, err
	}

	// Check if entity is nil
	if entity == nil {
		log.Printf("Model family not found with ID: %s", id)
		return nil, model_family.ErrNotFound
	}

	// Check if schema exists
	if entity.Schema.ID == "" {
		log.Printf("Model family has no schema associated - ID: %s", id)
		return nil, fmt.Errorf("schema not found for model family id: %s", id)
	}

	log.Printf("Successfully retrieved model family with schema by ID: %s, Name: %s", id, entity.Name)

	return entity, nil
}

func (s *Service) GetAll(ctx context.Context, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[ai_model.Response], error) {
	result, err := s.repo.GetAll(ctx, userCtx, params)
	if err != nil {
		return nil, err
	}

	responses := make([]ai_model.Response, len(result))
	for i, entity := range result {
		responses[i] = *ai_model.ToResponse(&entity)
	}

	return &paginate.PaginatedResult[ai_model.Response]{
		Data: responses,
	}, nil
}

func (s *Service) GetAllWithStatus(ctx context.Context, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[ai_model.ModelWithStatus], error) {
	result, err := s.repo.GetAllWithStatus(ctx, userCtx, params)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) GetByFamily(ctx context.Context, familyID string, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[ai_model.Response], error) {
	result, err := s.repo.GetByFamily(ctx, familyID, userCtx, params)
	if err != nil {
		return nil, err
	}

	responses := make([]ai_model.Response, len(result.Data))
	for i, entity := range result.Data {
		responses[i] = *ai_model.ToResponse(&entity)
	}

	return &paginate.PaginatedResult[ai_model.Response]{
		Data:        responses,
		TotalItems:  result.TotalItems,
		TotalPages:  result.TotalPages,
		CurrentPage: result.CurrentPage,
		Limit:       result.Limit,
		Offset:      result.Offset,
	}, nil
}

func (s *Service) GetByProvider(ctx context.Context, providerID string, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[ai_model.Response], error) {
	result, err := s.repo.GetByProvider(ctx, providerID, userCtx, params)
	if err != nil {
		return nil, err
	}

	responses := make([]ai_model.Response, len(result.Data))
	for i, entity := range result.Data {
		responses[i] = *ai_model.ToResponse(&entity)
	}

	return &paginate.PaginatedResult[ai_model.Response]{
		Data:        responses,
		TotalItems:  result.TotalItems,
		TotalPages:  result.TotalPages,
		CurrentPage: result.CurrentPage,
		Limit:       result.Limit,
		Offset:      result.Offset,
	}, nil
}

func (s *Service) GetByTeam(ctx context.Context, teamID string, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[ai_model.Response], error) {
	result, err := s.repo.GetByTeam(ctx, teamID, userCtx, params)
	if err != nil {
		return nil, err
	}

	responses := make([]ai_model.Response, len(result.Data))
	for i, entity := range result.Data {
		responses[i] = *ai_model.ToResponse(&entity)
	}

	return &paginate.PaginatedResult[ai_model.Response]{
		Data:        responses,
		TotalItems:  result.TotalItems,
		TotalPages:  result.TotalPages,
		CurrentPage: result.CurrentPage,
		Limit:       result.Limit,
		Offset:      result.Offset,
	}, nil
}

func (s *Service) GetDefault(ctx context.Context, userCtx models.UserContext) ([]ai_model.Response, error) {
	entities, err := s.repo.GetDefault(ctx, userCtx)
	if err != nil {
		return nil, err
	}

	responses := make([]ai_model.Response, len(entities))
	for i, entity := range entities {
		responses[i] = *ai_model.ToResponse(&entity)
	}

	return responses, nil
}

func (s *Service) SetAsDefault(ctx context.Context, id string, userCtx models.UserContext) error {
	// Get the model first
	model, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Determine provider ID
	providerID := ""
	if model.ProviderID != nil {
		providerID = *model.ProviderID
	} else if model.FamilyID != nil {
		// Would need to get provider from family repository
		return errors.New("cannot determine provider from family")
	}

	if providerID == "" {
		return errors.New("provider ID not found")
	}

	return s.repo.SetDefaultForProvider(ctx, providerID, id)
}

func (s *Service) Validate(ctx context.Context, model *ai_model.AIModel) error {
	if model.Name == "" {
		return errors.New("name is required")
	}
	if model.DisplayName == nil {
		return errors.New("display_name is required")
	}
	return nil
}
