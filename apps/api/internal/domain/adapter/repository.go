package adapter

import (
	"context"
	"database/sql"
	"seo-backend/internal/domain/product"
)

// AdapterConfigRepository interface
type AdapterConfigRepository interface {
	// Insert with transaction
	InsertWithTx(ctx context.Context, tx *sql.Tx, productID string, config *product.AdapterConfig) error

	// Update with transaction
	UpdateWithTx(ctx context.Context, tx *sql.Tx, productID string, cfg product.AdapterConfig) error

	// Delete with transaction
	DeleteWithTx(ctx context.Context, tx *sql.Tx, productID string) error

	// Get by product ID
	GetByProductID(ctx context.Context, productID string) (*product.AdapterConfig, error)

	// Load for product
	LoadForProduct(ctx context.Context, product *product.Product) error

	// Get with defaults
	GetOrDefault(ctx context.Context, productID string) (*product.AdapterConfig, error)
}

// AdapterConfigDTO for repository operations
type UpdateAdapterConfigRequest struct {
	EndpointPath    *string                `json:"endpointPath,omitempty"`
	HTTPMethod      *string                `json:"httpMethod,omitempty"`
	CustomHeaders   map[string]string      `json:"customHeaders,omitempty"`
	FieldMapping    map[string]interface{} `json:"fieldMapping,omitempty"`
	ResponseMapping interface{}            `json:"responseMapping,omitempty"`
	MetaConfig      map[string]interface{} `json:"metaConfig,omitempty"`
	SitemapConfig   map[string]interface{} `json:"sitemapConfig,omitempty"`
	TimeoutSeconds  *int                   `json:"timeoutSeconds,omitempty"`
	RetryCount      *int                   `json:"retryCount,omitempty"`
}
