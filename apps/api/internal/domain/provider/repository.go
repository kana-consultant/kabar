package provider

import (
	"context"
	"database/sql"
)

// Repository interface for API Provider data access
type Repository interface {
	// Create operations
	Create(ctx context.Context, tx *sql.Tx, provider *APIProvider) (string, error)

	// Read operations
	GetByID(ctx context.Context, id string) (*APIProvider, error)
	GetAll(ctx context.Context, userRole string) ([]APIProvider, error)
	GetActiveProviders(ctx context.Context) ([]APIProvider, error)

	// Update operations
	Update(ctx context.Context, tx *sql.Tx, id string, updates map[string]interface{}) error

	// Delete operations
	Delete(ctx context.Context, tx *sql.Tx, id string) error

	// Utility operations
	CheckProviderUsage(ctx context.Context, providerID string) (int, error)
	ExistsByName(ctx context.Context, name string) (bool, error)

	// Transaction management
	BeginTx(ctx context.Context) (*sql.Tx, error)
}
