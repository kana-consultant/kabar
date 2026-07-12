package user

import (
	"context"
	"database/sql"
	"seo-backend/internal/models"
)

// Repository interface for user data access
type Repository interface {
	// Read operations
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error) // ← Tambahan
	GetAll(ctx context.Context, query string, args []interface{}) ([]models.User, error)

	// Write operations
	Create(ctx context.Context, req CreateUserRequest, passwordHash []byte) (*models.User, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error

	// Utility operations
	EmailExists(ctx context.Context, email string) (bool, error)

	// Transaction management (optional, for future use)
	BeginTx(ctx context.Context) (*sql.Tx, error)
}
