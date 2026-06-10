package provider

import (
	"context"
	"database/sql"
	model_family "seo-backend/internal/domain/modelfamily"
	"seo-backend/internal/domain/paginate"
	"seo-backend/internal/models"
)

// Service defines the business logic interface for APIProvider
type Service interface {
	// Create operations
	Create(ctx context.Context, req *CreateRequest, userCtx models.UserContext) (*Response, error)

	// Update operations
	Update(ctx context.Context, id string, req *UpdateRequest, userCtx models.UserContext) (*APIProvider, error)

	// Delete operations
	Delete(ctx context.Context, id string, userCtx models.UserContext) error
	HardDelete(ctx context.Context, id string, userCtx models.UserContext) error

	// Read operations
	GetByID(ctx context.Context, id string, userCtx models.UserContext) (*Response, error)
	GetByName(ctx context.Context, name string, userCtx models.UserContext) (*Response, error)
	GetAll(ctx context.Context, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[Response], error)
	GetActive(ctx context.Context, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[Response], error)

	// Utility operations
	UpdateHeaders(ctx context.Context, id string, headers map[string]string, userCtx models.UserContext) error
	ToggleActive(ctx context.Context, id string, userCtx models.UserContext) error
	Validate(ctx context.Context, provider *APIProvider) error

	// Model Family operations
	CreateModelFamily(ctx context.Context, providerID string, families []model_family.ModelFamilyWithSchema, userCtx models.UserContext) ([]model_family.Response, error)
	UpdateModelFamily(ctx context.Context, familyID string, family *model_family.ModelFamily, userCtx models.UserContext) (*model_family.Response, error)
	DeleteModelFamily(ctx context.Context, familyID string, userCtx models.UserContext) error
	CreateOrUpdateModelFamilies(ctx context.Context, tx *sql.Tx, providerID string, families []model_family.ModelFamilyWithSchema) error
	GetModelFamiliesByProvider(ctx context.Context, providerID string, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[model_family.Response], error)
}
