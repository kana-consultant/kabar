package product

import (
	"context"
	"database/sql"
)

type ProductRepository interface {
	// Transaction management
	BeginTx(ctx context.Context) (*sql.Tx, error)

	// Read operations
	GetByID(ctx context.Context, id string) (*Product, error)
	GetByIDWithTx(ctx context.Context, tx *sql.Tx, id string) (*Product, error)
	GetAll(ctx context.Context, query string, args []interface{}) ([]Product, error)
	GetAllWithFilters(ctx context.Context) ([]Product, int, error)
	GetProductsByTeamID(ctx context.Context, teamID string) ([]Product, error)
	GetProductsByUserID(ctx context.Context, userID string) ([]Product, error)
	GetProductBasicInfo(ctx context.Context, id string) (*ProductBasicInfo, error)
	GetProductBasicInfoWithTx(ctx context.Context, tx *sql.Tx, id string) (*ProductBasicInfo, error)

	// Write operations
	InsertProductWithTx(ctx context.Context, tx *sql.Tx, req CreateProductRequest) (string, error)
	UpdateProductWithTx(ctx context.Context, tx *sql.Tx, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	DeleteWithTx(ctx context.Context, tx *sql.Tx, id string) error
	DeleteProductWithTx(ctx context.Context, tx *sql.Tx, id string) error

	// Business operations
	UpdateConnectionStatus(ctx context.Context, productID string, isConnected bool) error
	UpdateConnectionStatusWithTx(ctx context.Context, tx *sql.Tx, productID string, isConnected bool) error
}
