package model_family

import (
	"context"
	"seo-backend/internal/domain/paginate"
	"seo-backend/internal/models"
)

// Service defines the business logic interface for ModelFamily
type Service interface {
	Create(ctx context.Context, req *CreateRequest) (*Response, error)
	Update(ctx context.Context, id string, req *UpdateRequest) (*Response, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Response, error)
	GetByProviderAndName(ctx context.Context, providerID string, name string) (*Response, error)
	GetAll(ctx context.Context, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[Response], error)
	GetByProvider(ctx context.Context, providerID string) ([]Response, error)
	GetBySchema(ctx context.Context, schemaID string) ([]Response, error)
	Validate(ctx context.Context, modelFamily *ModelFamily) error
}
