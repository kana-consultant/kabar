// internal/domain/adapter/service.go
package adapter

import (
	"context"

	"seo-backend/internal/models"
)

// AdapterConfigService defines the adapter configuration business logic interface
type AdapterConfigService interface {
	// GetAdapterConfig gets config with defaults
	GetAdapterConfig(ctx context.Context, productID string) (*models.AdapterConfig, error)

	// UpdateAdapterConfig updates adapter configuration
	UpdateAdapterConfig(ctx context.Context, productID string, updates map[string]interface{}) error

	// CreateOrUpdateAdapterConfig creates or updates full config
	CreateOrUpdateAdapterConfig(ctx context.Context, productID string, config *models.AdapterConfig) error

	// LoadConfigForProduct loads and attaches config to product
	LoadConfigForProduct(ctx context.Context, product *models.Product) error
}
