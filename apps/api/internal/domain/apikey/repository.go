package apikey

import (
	"context"
	"database/sql"
	"seo-backend/internal/models"
)

// Repository interface for API Key data access
type Repository interface {
	// Basic CRUD with transaction
	Create(ctx context.Context, tx *sql.Tx, key *APIKey) (string, error)
	GetByID(ctx context.Context, id string) (*APIKey, error)
	GetAll(ctx context.Context, filter models.UserContext) ([]APIKeyDetail, error)
	Update(ctx context.Context, tx *sql.Tx, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, tx *sql.Tx, id string) error

	// Transaction management
	BeginTx(ctx context.Context) (*sql.Tx, error)

	// Encryption key management (repository doesn't encrypt, just stores)
	UpdateKey(ctx context.Context, tx *sql.Tx, id string, encryptedKey string) error
}
