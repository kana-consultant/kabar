package workflow_node

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

type WorkflowNode struct {
	ID              string          `json:"id"`
	WorkflowID      string          `json:"workflowId"`
	AdapterConfigID string          `json:"adapterConfigId"`
	PreviousNodeIDs []string        `json:"previousNodeIds,omitempty"`
	StepOrder       int             `json:"stepOrder"`
	InputMapping    json.RawMessage `json:"inputMapping"`
	NextNodeID      *string         `json:"nextNodeId"`
	CreatedAt       time.Time       `json:"createdAt"`
	AdapterConfig   *AdapterConfig  `json:"adapterConfig,omitempty"` // optional, for joined data
}

type AdapterConfig struct {
	ID              string          `json:"id"`
	ProductID       string          `json:"productId"`
	EndpointPath    string          `json:"endpointPath"`
	HTTPMethod      string          `json:"httpMethod"`
	CustomHeaders   string          `json:"customHeaders"`
	FieldMapping    string          `json:"fieldMapping"`
	ResponseMapping json.RawMessage `json:"responseMapping"`
	MetaConfig      string          `json:"metaConfig"`
	SitemapConfig   string          `json:"sitemapConfig"`
	TimeoutSeconds  int             `json:"timeoutSeconds"`
	RetryCount      int             `json:"retryCount"`
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

	// Methods with transaction support
	InsertBatchWithTx(ctx context.Context, tx *sql.Tx, nodes []WorkflowNode) ([]WorkflowNode, error)
	InsertWithTx(ctx context.Context, tx *sql.Tx, node *WorkflowNode) error
	UpdateWithTx(ctx context.Context, tx *sql.Tx, id string, updates map[string]interface{}) error
}
