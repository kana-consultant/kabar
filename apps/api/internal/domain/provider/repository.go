package provider

import (
	"context"
	"database/sql"
	"encoding/json"

	"seo-backend/internal/domain/paginate"
	"seo-backend/internal/models"
)

// Repository defines the data access interface for APIProvider
type Repository interface {
	// Transaction management
	BeginTx(ctx context.Context) (*sql.Tx, error)

	// Create operations
	Create(ctx context.Context, provider *APIProvider) error
	CreateWithTx(ctx context.Context, tx *sql.Tx, provider *APIProvider) error

	// Read operations
	GetByID(ctx context.Context, id string) (*APIProvider, error)
	GetByName(ctx context.Context, name string) (*APIProvider, error)
	GetAll(ctx context.Context, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[APIProvider], error)
	GetActive(ctx context.Context) ([]APIProvider, error)

	// Update operations
	Update(ctx context.Context, provider *APIProvider) error
	UpdateWithTx(ctx context.Context, tx *sql.Tx, id string, updates map[string]interface{}) error
	UpdateDefaultHeaders(ctx context.Context, id string, headers json.RawMessage) error
	UpdateDefaultHeadersWithTx(ctx context.Context, tx *sql.Tx, id string, headers json.RawMessage) error
	ToggleActive(ctx context.Context, id string, isActive bool) error
	ToggleActiveWithTx(ctx context.Context, tx *sql.Tx, id string, isActive bool) error

	// Delete operations
	Delete(ctx context.Context, id string) error
	DeleteWithTx(ctx context.Context, tx *sql.Tx, id string) error
	HardDelete(ctx context.Context, id string) error
	HardDeleteWithTx(ctx context.Context, tx *sql.Tx, id string) error

	// Utility operations
	Exists(ctx context.Context, name string) (bool, error)
	Count(ctx context.Context) (int64, error)
	CheckProviderUsage(ctx context.Context, providerID string) (int, error)
}
