// internal/domain/apikey/service.go
package apikey

import (
	"context"
	"seo-backend/internal/models"
)

// Service defines the API key business logic interface
type Service interface {
	// CreateAPIKey creates a new API key
	CreateAPIKey(ctx context.Context, req CreateAPIKeyRequest, userCtx models.UserContext) (string, error)

	// GetAPIKeyByID retrieves an API key by ID
	GetAPIKeyByID(ctx context.Context, id string) (*APIKey, error)

	// GetAllAPIKeys retrieves all API keys with filters
	GetAllAPIKeys(ctx context.Context, userCtx models.UserContext) ([]APIKeyDetail, error)

	// UpdateAPIKey updates an API key
	UpdateAPIKey(ctx context.Context, id string, req UpdateAPIKeyRequest, userCtx models.UserContext) error

	// DeleteAPIKey deletes an API key
	DeleteAPIKey(ctx context.Context, id string, userCtx models.UserContext) error
}
