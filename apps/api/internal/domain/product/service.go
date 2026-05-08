// internal/domain/product/service.go
package product

import (
	"context"

	"seo-backend/internal/models"
)

// ProductService defines the product business logic interface
type ProductService interface {
	// CreateProduct creates a new product
	CreateProduct(ctx context.Context, req models.CreateProductRequest, userCtx models.UserContext) (string, error)

	// UpdateProduct updates an existing product
	UpdateProduct(ctx context.Context, id string, updates map[string]interface{}, userCtx models.UserContext) error

	// DeleteProduct deletes a product
	DeleteProduct(ctx context.Context, id string) error

	// GetByID retrieves a product by ID with its adapter config
	GetByID(ctx context.Context, id string) (*models.Product, error)

	// GetAllProducts retrieves all products with filters
	GetAllProducts(ctx context.Context) ([]models.Product, int, error)

	// UpdateConnectionStatus updates product connection status
	UpdateConnectionStatus(ctx context.Context, productID string, isConnected bool) error
}