package product

import (
	"seo-backend/internal/domain/workflow"
	"time"
)

type UserContext interface {
	GetUserID() string
	GetTeamID() string
	GetRole() string
	IsAdmin() bool
}

// ProductFilters for filtering products
type ProductFilters struct {
	Status     string
	SyncStatus string
	Platform   string
	TeamID     string
	UserID     string
	Search     string
	OrderBy    string
	Limit      int
	Offset     int
}

type ConnectionTestResult struct {
	Success     bool   `json:"success"`
	StatusCode  int    `json:"status_code"`
	StatusText  string `json:"status_text"`
	Message     string `json:"message"`
	ProductName string `json:"product_name"`
	Endpoint    string `json:"endpoint"`
	Method      string `json:"method"`
	Response    string `json:"response"`
	TestedAt    string `json:"tested_at"`
}

type ProductBasicInfo struct {
	ID            string
	Name          string
	CustomHeaders string
	Platform      string
	APIEndpoint   string
	APIKey        string
}

type ProductPlatform string

const (
	PlatformWordpress ProductPlatform = "wordpress"
	PlatformShopify   ProductPlatform = "shopify"
	PlatformCustom    ProductPlatform = "custom"
)

type ProductStatus string

const (
	ProductStatusConnected    ProductStatus = "connected"
	ProductStatusPending      ProductStatus = "pending"
	ProductStatusError        ProductStatus = "error"
	ProductStatusDisconnected ProductStatus = "disconnected"
)

type SyncStatus string

const (
	SyncStatusIdle    SyncStatus = "idle"
	SyncStatusSyncing SyncStatus = "syncing"
	SyncStatusSuccess SyncStatus = "success"
	SyncStatusFailed  SyncStatus = "failed"
)

// FieldMapping untuk struktur mapping field
type FieldMapping struct {
	ID           string `json:"id"`
	SourceField  string `json:"sourceField"`
	TargetField  string `json:"targetField"`
	IsRequired   bool   `json:"isRequired"`
	DefaultValue string `json:"defaultValue,omitempty"`
}

// NestedMapping untuk mapping bersarang
type NestedMapping struct {
	ID           string         `json:"id"`
	SourceField  string         `json:"sourceField"`
	TargetField  string         `json:"targetField"`
	IsRequired   bool           `json:"isRequired"`
	DefaultValue string         `json:"defaultValue,omitempty"`
	Children     []FieldMapping `json:"children,omitempty"`
	IsExpanded   bool           `json:"isExpanded,omitempty"`
}

// AdapterConfig untuk konfigurasi adapter
type AdapterConfig struct {
	ID              string      `json:"id"`
	ProductID       string      `json:"productId"`
	EndpointPath    string      `json:"endpointPath"`
	HTTPMethod      string      `json:"httpMethod"`
	CustomHeaders   string      `json:"customHeaders"`
	FieldMapping    string      `json:"fieldMapping"`
	ResponseMapping interface{} `json:"responseMapping"` // bisa string atau object
	MetaConfig      string      `json:"metaConfig,omitempty"`
	SitemapConfig   string      `json:"sitemapConfig,omitempty"`
	TimeoutSeconds  int         `json:"timeoutSeconds,omitempty"`
	RetryCount      int         `json:"retryCount,omitempty"`
	CreatedAt       time.Time   `json:"createdAt,omitempty"`
	UpdatedAt       time.Time   `json:"updatedAt,omitempty"`
}

// WorkflowDefinition untuk definisi workflow
type WorkflowDefinition struct {
	ID        string         `json:"id"`
	ProductID string         `json:"productId"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	Nodes     []WorkflowNode `json:"nodes,omitempty"`
}

// WorkflowNode untuk node dalam workflow
type WorkflowNode struct {
	ID              string                 `json:"id"`
	WorkflowID      string                 `json:"workflowId"`
	AdapterConfigID string                 `json:"adapterConfigId"`
	StepOrder       int                    `json:"stepOrder"`
	InputMapping    map[string]interface{} `json:"inputMapping"`
	NextNodeID      *string                `json:"nextNodeId,omitempty"`
	CreatedAt       time.Time              `json:"createdAt"`
	AdapterConfig   *AdapterConfig         `json:"adapterConfig,omitempty"`
}

// Product struct utama
type Product struct {
	ID              string                        `json:"id"`
	Name            string                        `json:"name"`
	Platform        string                        `json:"platform"`
	APIEndpoint     string                        `json:"apiEndpoint"`
	APIKeyEncrypted string                        `json:"apiKey,omitempty"`
	Status          ProductStatus                 `json:"status"`
	LastSync        *time.Time                    `json:"lastSync,omitempty"`
	SyncStatus      SyncStatus                    `json:"syncStatus"`
	CreatedBy       *string                       `json:"createdBy,omitempty"`
	TeamID          *string                       `json:"teamId,omitempty"`
	UserID          *string                       `json:"userId,omitempty"`
	CreatedAt       time.Time                     `json:"createdAt"`
	UpdatedAt       time.Time                     `json:"updatedAt"`
	WorkflowID      string                        `json:"workflow_id"`
	AdapterConfig   *AdapterConfig                `json:"adapterConfig,omitempty"`
	AdapterConfigs  []AdapterConfig               `json:"adapterConfigs,omitempty"`
	Workflows       []workflow.WorkflowDefinition `json:"workflows,omitempty"`
}

// ProductWithDetails untuk response dengan data lengkap
type ProductWithDetails struct {
	Product
	AdapterConfigs   []AdapterConfig               `json:"adapterConfigs"`
	Workflows        []workflow.WorkflowDefinition `json:"workflows"`
	ActiveWorkflowID string                        `json:"activeWorkflowId,omitempty"`
}

// Request/Response types
type CreateProductRequest struct {
	Name          string                      `json:"name"`
	Platform      string                      `json:"platform"`
	APIEndpoint   string                      `json:"apiEndpoint"`
	APIKey        string                      `json:"apiKey"`
	TeamID        string                      `json:"teamId,omitempty"`
	AdapterConfig *CreateAdapterConfigRequest `json:"adapterConfig,omitempty"`

	// TAMBAHKAN FIELD INI UNTUK DEBUG
	AdapterConfigs []AdapterConfig               `json:"adapterConfigs,omitempty"`
	Workflows      []workflow.WorkflowDefinition `json:"workflows,omitempty"`
}

type CreateAdapterConfigRequest struct {
	EndpointPath    string      `json:"endpointPath"`
	HTTPMethod      string      `json:"httpMethod"`
	CustomHeaders   string      `json:"customHeaders"`
	FieldMapping    string      `json:"fieldMapping"`
	ResponseMapping interface{} `json:"responseMapping,omitempty"`
	MetaConfig      string      `json:"metaConfig,omitempty"`
	SitemapConfig   string      `json:"sitemapConfig,omitempty"`
	TimeoutSeconds  int         `json:"timeoutSeconds,omitempty"`
	RetryCount      int         `json:"retryCount,omitempty"`
}

// Batch save request untuk workflow nodes
type BatchCreateNode struct {
	TempID          string                 `json:"tempId"`
	AdapterConfigID string                 `json:"adapterConfigId"`
	StepOrder       int                    `json:"stepOrder"`
	InputMapping    map[string]interface{} `json:"inputMapping"`
	NextNodeID      *string                `json:"nextNodeId,omitempty"`
}

type BatchUpdateNode struct {
	ID      string                 `json:"id"`
	Updates map[string]interface{} `json:"updates"`
}

type BatchSaveRequest struct {
	ToCreate []BatchCreateNode `json:"toCreate"`
	ToUpdate []BatchUpdateNode `json:"toUpdate"`
	ToDelete []string          `json:"toDelete"`
}

type BatchCreateResult struct {
	TempID string      `json:"tempId"`
	Node   interface{} `json:"node"`
}

type BatchError struct {
	Operation string `json:"operation"`
	ID        string `json:"id,omitempty"`
	Message   string `json:"message"`
}

type BatchSaveResult struct {
	Created []BatchCreateResult `json:"created"`
	Updated []string            `json:"updated"`
	Deleted []string            `json:"deleted"`
	Errors  []BatchError        `json:"errors,omitempty"`
}
