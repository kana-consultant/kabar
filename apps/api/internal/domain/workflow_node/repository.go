package workflow_node

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

type NodeAdapterConfig struct {
	EndpointPath string `json:"endpoint_path"`
	HTTPMethod   string `json:"http_method"`
	FieldMapping string `json:"field_mapping"`
}

type WorkflowNode struct {
	ID              string   `json:"id"`
	WorkflowID      string   `json:"workflow_id"`
	AdapterConfigID string   `json:"adapter_config_id"`
	PreviousNodeIDs []string `json:"previous_node_ids,omitempty"`
	StepOrder       int      `json:"step_order"`

	NextNodeIDs   []string           `json:"next_node_ids"`
	CreatedAt     time.Time          `json:"created_at"`
	AdapterConfig *NodeAdapterConfig `json:"adapter_config,omitempty"`
}

type WorkflowNodeCreate struct {
	ID              string
	WorkflowID      string
	AdapterConfigID string
	PreviousNodeIDs []string
	StepOrder       int
	InputMapping    json.RawMessage
	NextNodeIDs     []string
	CreatedAt       time.Time
	EndpointPath    *string
	HTTPMethod      string
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
	GetByWorkflowIDs(ctx context.Context, workflowIDs []string) ([]*WorkflowNodeCreate, error)
	Insert(ctx context.Context, node *WorkflowNode) error
	InsertBatch(ctx context.Context, nodes []WorkflowNode) ([]WorkflowNode, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	DeleteByWorkflowID(ctx context.Context, workflowID string) error
	ReorderNodes(ctx context.Context, workflowID string, nodeIDs []string) error
	UpsertWithTx(ctx context.Context, tx *sql.Tx, nodes []WorkflowNode) error

	// Methods with transaction support
	InsertBatchWithTx(ctx context.Context, tx *sql.Tx, nodes []WorkflowNodeCreate) ([]WorkflowNodeCreate, error)
	InsertWithTx(ctx context.Context, tx *sql.Tx, node *WorkflowNode) error
	UpdateWithTx(ctx context.Context, tx *sql.Tx, id string, updates map[string]interface{}) error
}
