package workflow_node

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

type WorkflowNode struct {
	ID              string          `json:"id"`
	WorkflowID      string          `json:"workflow_id"`
	AdapterConfigID string          `json:"adapter_config_id"`
	PreviousNodeIDs []string        `json:"previous_node_ids,omitempty"`
	StepOrder       int             `json:"step_order"`
	InputMapping    json.RawMessage `json:"input_mapping"`
	NextNodeID      *string         `json:"next_node_id"`
	CreatedAt       time.Time       `json:"created_at"`
	AdapterConfig   *AdapterConfig  `json:"adapter_config,omitempty"` // optional, for joined data
}

type AdapterConfig struct {
	ID              string          `json:"id"`
	ProductID       string          `json:"product_id"`
	EndpointPath    string          `json:"endpoint_path"`
	HTTPMethod      string          `json:"http_method"`
	CustomHeaders   string          `json:"custom_headers"`
	FieldMapping    string          `json:"field_mapping"`
	ResponseMapping json.RawMessage `json:"response_mapping"`
	MetaConfig      string          `json:"meta_config"`
	SitemapConfig   string          `json:"sitemap_config"`
	TimeoutSeconds  int             `json:"timeout_seconds"`
	RetryCount      int             `json:"retry_count"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
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
