package request_schema

import (
	"context"
)

// Service defines the business logic interface for RequestSchema
type Service interface {
	Create(ctx context.Context, req *CreateRequest) (*Response, error)
	Update(ctx context.Context, id string, req *UpdateRequest) (*Response, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Response, error)
	GetByProviderAndName(ctx context.Context, providerID string, name string) (*Response, error)
	GetAll(ctx context.Context, page, limit int) (*ListResponse, error)
	GetByProvider(ctx context.Context, providerID string) ([]Response, error)
	Validate(ctx context.Context, requestSchema *RequestSchema) error
}
