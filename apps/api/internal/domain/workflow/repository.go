// internal/domain/workflow/workflow.go
package workflow

import (
	"context"
	"time"
)

type WorkflowDefinition struct {
	ID        string    `json:"id"`
	ProductID string    `json:"productId"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type WorkflowDefinitionRepository interface {
	GetByID(ctx context.Context, id string) (*WorkflowDefinition, error)
	GetByProductID(ctx context.Context, productID string) ([]*WorkflowDefinition, error)
	Insert(ctx context.Context, wf *WorkflowDefinition) error
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
}
