package ai_model

import (
	"context"
	"errors"
	model_family "seo-backend/internal/domain/modelfamily"
	"seo-backend/internal/domain/paginate"
	"seo-backend/internal/models"
)

// Errors
var (
	ErrNotFound          = errors.New("AI model not found")
	ErrDuplicate         = errors.New("AI model with this name already exists")
	ErrDuplicateInFamily = errors.New("AI model with this name already exists in the family")
	ErrDatabase          = errors.New("database error")
)

// ModelWithStatus for API responses with limited fields
type ModelWithStatus struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ProviderID  string `json:"provider_id"`
	DisplayName string `json:"display_name"`
	IsActive    *bool  `json:"is_active"`
	IsDefault   *bool  `json:"is_default"`
}

// Repository interface
type Repository interface {
	Create(ctx context.Context, model *AIModel) error
	GetByID(ctx context.Context, id string) (*AIModel, error)
	GetByIDWithSchema(ctx context.Context, id string) (*model_family.ModelFamilyWithSchema, error)
	GetAll(ctx context.Context, userCtx models.UserContext, params paginate.PaginationParams) ([]AIModel, error)
	GetByFamily(ctx context.Context, familyID string, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[AIModel], error)
	GetByProvider(ctx context.Context, providerID string, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[AIModel], error)
	GetByTeam(ctx context.Context, teamID string, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[AIModel], error)
	GetDefault(ctx context.Context, userCtx models.UserContext) ([]AIModel, error)
	GetAllWithStatus(ctx context.Context, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[ModelWithStatus], error)
	Update(ctx context.Context, model *AIModel) error
	Delete(ctx context.Context, id string) error
	HardDelete(ctx context.Context, id string) error
	Exists(ctx context.Context, name string) (bool, error)
	ExistsByFamilyAndName(ctx context.Context, familyID, name string) (bool, error)
	Count(ctx context.Context) (int64, error)
	CountByFamily(ctx context.Context, familyID string) (int64, error)
	SetDefaultForProvider(ctx context.Context, providerID, modelID string) error
}
