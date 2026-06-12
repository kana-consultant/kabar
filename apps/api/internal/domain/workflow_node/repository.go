package workflow_node

import (
	"context"
	"encoding/json"
	"time"
)

type WorkflowNode struct {
	ID              string          `json:"id"`
	WorkflowID      string          `json:"workflowId"`
	AdapterConfigID string          `json:"adapterConfigId"`
	StepOrder       int             `json:"stepOrder"`
	InputMapping    json.RawMessage `json:"inputMapping"`
	NextNodeID      *string         `json:"nextNodeId"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
}
type WorkflowNodeRepository interface {
	GetByID(ctx context.Context, id string) (*WorkflowNode, error)
	GetByWorkflowID(ctx context.Context, workflowID string) ([]*WorkflowNode, error)
	Insert(ctx context.Context, node *WorkflowNode) error
	InsertBatch(ctx context.Context, nodes []WorkflowNode) ([]WorkflowNode, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	DeleteByWorkflowID(ctx context.Context, workflowID string) error
	ReorderNodes(ctx context.Context, workflowID string, nodeIDs []string) error
}
