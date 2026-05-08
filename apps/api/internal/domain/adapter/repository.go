package adapter

import (
	"context"
	"database/sql"
	"seo-backend/internal/models"
)

// AdapterConfigRepository interface
// internal/domain/adapter/repository.go

type AdapterConfigRepository interface {
	// Insert with transaction
	InsertWithTx(ctx context.Context, tx *sql.Tx, productID string, config *models.AdapterConfig) error

	// Update with transaction
	UpdateWithTx(ctx context.Context, tx *sql.Tx, productID string, updates map[string]interface{}) error

	// Delete with transaction
	DeleteWithTx(ctx context.Context, tx *sql.Tx, productID string) error // TAMBAHKAN INI

	// Get by product ID
	GetByProductID(ctx context.Context, productID string) (*models.AdapterConfig, error)

	// Load for product
	LoadForProduct(ctx context.Context, product *models.Product) error

	// Get with defaults
	GetOrDefault(ctx context.Context, productID string) (*models.AdapterConfig, error)
}

// AdapterConfigDTO for repository operations
type UpdateAdapterConfigRequest struct {
	EndpointPath   *string
	HTTPMethod     *string
	CustomHeaders  map[string]string
	FieldMapping   map[string]interface{}
	TimeoutSeconds *int
	RetryCount     *int
}
