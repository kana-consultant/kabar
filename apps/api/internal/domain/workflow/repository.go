// internal/domain/workflow/workflow.go
package workflow

import (
	"context"
	"database/sql"
	"seo-backend/internal/domain/workflow_node"
	"time"
)

type WorkflowDefinition struct {
	ID        string                       `json:"id"`
	ProductID string                       `json:"product_id"` // fix: snake_case
	Name      string                       `json:"name"`
	CreatedAt time.Time                    `json:"created_at"` // fix: snake_case
	UpdatedAt time.Time                    `json:"updated_at"` // fix: snake_case
	Nodes     []workflow_node.WorkflowNode `json:"nodes"`
}

type WorkflowDefinitionRepository interface {
	GetByID(ctx context.Context, id string) (*WorkflowDefinition, error)
	GetByProductID(ctx context.Context, productID string) ([]*WorkflowDefinition, error)
	Insert(ctx context.Context, wf *WorkflowDefinition) error
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	InsertWithTx(ctx context.Context, tx *sql.Tx, node *WorkflowDefinition) error
	UpsertWithTx(ctx context.Context, tx *sql.Tx, wf *WorkflowDefinition) error
}
