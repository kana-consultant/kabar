package ai_model

import (
	"context"
	model_family "seo-backend/internal/domain/modelfamily"
	"seo-backend/internal/domain/paginate"
	"seo-backend/internal/models"
)

// Service defines the business logic interface for AIModel
type Service interface {
	Create(ctx context.Context, req *CreateRequest, userCtx models.UserContext) (*Response, error)
	Update(ctx context.Context, id string, req *UpdateRequest, userCtx models.UserContext) (*Response, error)
	Delete(ctx context.Context, id string, userCtx models.UserContext) error
	GetByID(ctx context.Context, id string, userCtx models.UserContext) (*Response, error)
	GetSchemaByID(ctx context.Context, id string, userCtx models.UserContext) (*model_family.ModelFamilyWithSchema, error)
	GetAll(ctx context.Context, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[Response], error)
	GetAllWithStatus(ctx context.Context, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[ModelWithStatus], error)
	GetByFamily(ctx context.Context, familyID string, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[Response], error)
	GetByProvider(ctx context.Context, providerID string, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[Response], error)
	GetByTeam(ctx context.Context, teamID string, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[Response], error)
	GetDefault(ctx context.Context, userCtx models.UserContext) ([]Response, error)
	SetAsDefault(ctx context.Context, id string, userCtx models.UserContext) error
	Validate(ctx context.Context, model *AIModel) error
}
