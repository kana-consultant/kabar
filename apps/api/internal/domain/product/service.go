// internal/domain/product/service.go
package product

import (
	"context"

	"seo-backend/internal/domain/draft"
	"seo-backend/internal/domain/workflow_node"
	"seo-backend/internal/models"
)

// ProductService defines the product business logic interface
type ProductService interface {
	// CreateProduct creates a new product
	CreateProduct(ctx context.Context, req ProductRequest, userCtx models.UserContext) (string, error)

	// UpdateProduct updates an existing product
	UpdateProduct(ctx context.Context, id string, req ProductRequest, userCtx models.UserContext) error

	// DeleteProduct deletes a product
	DeleteProduct(ctx context.Context, id string) error

	// GetByID retrieves a product by ID with its adapter config
	GetByID(ctx context.Context, productID string, userCtx models.UserContext) (*ProductRequest, error)

	// GetAllProducts retrieves all products with filters
	GetAllProducts(ctx context.Context, userCtx models.UserContext) ([]Product, error)

	// UpdateConnectionStatus updates product connection status
	UpdateConnectionStatus(ctx context.Context, productID string, isConnected bool) error

	GetProductConfig(ctx context.Context, productID string, draft draft.DraftDataPost) (*ProductConfig, error)
	ReorderNodesWithBatch(nodes []*workflow_node.WorkflowNode) ([]*workflow_node.WorkflowNode, map[int][]*workflow_node.WorkflowNode)
	ReorderNodesByLevel(nodes []*workflow_node.WorkflowNode) []*workflow_node.WorkflowNode
	LoadAdapterConfig(ctx context.Context, cfg *ProductConfig) error
	ParseCustomHeaders(cfg *ProductConfig) error
	SendWithRetry(cfg ProductConfig, requestBody interface{}) (interface{}, error)
}
