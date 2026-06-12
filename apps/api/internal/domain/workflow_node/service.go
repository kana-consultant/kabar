// internal/domain/workflow_node/service.go
package workflow_node

import "context"

type WorkflowNodeService interface {
	GetByID(ctx context.Context, id string) (*WorkflowNode, error)
	GetByWorkflowID(ctx context.Context, workflowID string) ([]*WorkflowNode, error)
	Create(ctx context.Context, node *WorkflowNode) error
	CreateBatch(ctx context.Context, nodes []WorkflowNode) ([]WorkflowNode, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	ReorderNodes(ctx context.Context, workflowID string, nodeIDs []string) error
	SaveBatch(ctx context.Context, workflowID string, req BatchSaveRequest) (*BatchSaveResult, error)
}

// BatchSaveRequest represents batch operation request from frontend
type BatchSaveRequest struct {
	ToCreate []BatchCreateNode `json:"toCreate"`
	ToUpdate []BatchUpdateNode `json:"toUpdate"`
	ToDelete []string          `json:"toDelete"`
}

type BatchCreateNode struct {
	TempID          string                 `json:"tempId"`
	AdapterConfigID string                 `json:"adapterConfigId"`
	StepOrder       int                    `json:"stepOrder"`
	InputMapping    map[string]interface{} `json:"inputMapping"`
	NextNodeID      *string                `json:"nextNodeId"`
}

type BatchUpdateNode struct {
	ID      string                 `json:"id"`
	Updates map[string]interface{} `json:"updates"`
}

type BatchSaveResult struct {
	Created []BatchCreateResult `json:"created"`
	Updated []string            `json:"updated"`
	Deleted []string            `json:"deleted"`
	Errors  []BatchError        `json:"errors,omitempty"`
}

type BatchCreateResult struct {
	TempID string       `json:"tempId"`
	Node   WorkflowNode `json:"node"`
}

type BatchError struct {
	Operation string `json:"operation"`
	ID        string `json:"id,omitempty"`
	Message   string `json:"message"`
}
