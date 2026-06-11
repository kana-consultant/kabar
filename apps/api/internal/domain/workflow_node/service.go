package workflow_node

import "context"

type WorkflowNodeService interface {
	GetByID(ctx context.Context, id string) (*WorkflowNode, error)
	GetByWorkflowID(ctx context.Context, workflowID string) ([]*WorkflowNode, error)
	Create(ctx context.Context, node *WorkflowNode) error
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	ReorderNodes(ctx context.Context, workflowID string, nodeIDs []string) error
}
