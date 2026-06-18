package product

import (
	"fmt"
	"seo-backend/internal/domain/workflow_node"
	"strings"
	"time"
)

// ProductConfig - VERSI SEDERHANA (HANYA YANG DIPAKAI)
type ProductConfig struct {
	// Basic product info
	ProductID   string
	APIEndpoint string
	APIKey      string

	// Adapter config
	AdapterEndpoint string
	HTTPMethod      string
	FullURL         string
	BaseURL         string

	// Config strings (JSON)
	FieldMappingStr  string
	MetaConfigStr    string
	SitemapConfigStr string
	CustomHeadersStr string

	// Parsed configs
	CustomHeaders map[string]string
	FieldMapping  []FieldMapping         `json:"fieldMapping,omitempty"`
	MetaConfig    map[string]interface{} `json:"metaConfig,omitempty"`
	SitemapConfig map[string]interface{} `json:"sitemapConfig,omitempty"`

	// Timeout and retry
	Timeout    int
	RetryCount int

	// Workflow nodes (reordered)
	WorkflowNodes    []workflow_node.WorkflowNode          `json:"workflowNodes,omitempty"`
	WorkflowLevelMap map[int][]*workflow_node.WorkflowNode `json:"workflowLevelMap,omitempty"`

	// Execution context - untuk data passing antar nodes
	ExecutionResults map[string]interface{} `json:"executionResults,omitempty"`
	Variables        map[string]interface{} `json:"variables,omitempty"`
	CurrentNodeID    string                 `json:"currentNodeId,omitempty"`

	// Metadata
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

// ============================================================
// 1. ✅ DIPAKAI - BuildFullURL
// ============================================================
func (c *ProductConfig) BuildFullURL() string {
	if c.APIEndpoint == "" {
		return ""
	}
	base := strings.TrimRight(c.APIEndpoint, "/")
	if c.AdapterEndpoint != "" {
		return base + "/" + strings.TrimLeft(c.AdapterEndpoint, "/")
	}
	return base
}

// ============================================================
// 2. ✅ DIPAKAI - HasWorkflowNodes
// ============================================================
func (c *ProductConfig) HasWorkflowNodes() bool {
	return c.WorkflowNodes != nil && len(c.WorkflowNodes) > 0
}

// ============================================================
// 3. 🔑 PENTING - Untuk data passing antar nodes (AKAN DIPAKAI)
// ============================================================
func (c *ProductConfig) SetExecutionResult(nodeID string, result interface{}) {
	if c.ExecutionResults == nil {
		c.ExecutionResults = make(map[string]interface{})
	}
	c.ExecutionResults[nodeID] = result
}

func (c *ProductConfig) GetExecutionResult(nodeID string) interface{} {
	if c.ExecutionResults == nil {
		return nil
	}
	return c.ExecutionResults[nodeID]
}

func (c *ProductConfig) GetExecutionResultField(nodeID, fieldPath string) (interface{}, error) {
	result := c.GetExecutionResult(nodeID)
	if result == nil {
		return nil, fmt.Errorf("node %s result not found", nodeID)
	}

	if fieldPath == "" {
		return result, nil
	}

	parts := strings.Split(fieldPath, ".")
	current := result

	for _, part := range parts {
		if currentMap, ok := current.(map[string]interface{}); ok {
			val, exists := currentMap[part]
			if !exists {
				return nil, fmt.Errorf("field '%s' not found", part)
			}
			current = val
		} else {
			return nil, fmt.Errorf("cannot navigate to '%s'", part)
		}
	}
	return current, nil
}

func (c *ProductConfig) SetVariable(key string, value interface{}) {
	if c.Variables == nil {
		c.Variables = make(map[string]interface{})
	}
	c.Variables[key] = value
}

func (c *ProductConfig) GetVariable(key string) interface{} {
	if c.Variables == nil {
		return nil
	}
	return c.Variables[key]
}

// ============================================================
// 4. 🔧 UTILITY - Untuk debugging
// ============================================================
func (c *ProductConfig) Validate() error {
	if c.ProductID == "" {
		return fmt.Errorf("product ID is required")
	}
	if c.APIEndpoint == "" {
		return fmt.Errorf("API endpoint is required")
	}
	return nil
}

func (c *ProductConfig) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"product_id":        c.ProductID,
		"api_endpoint":      c.APIEndpoint,
		"full_url":          c.FullURL,
		"timeout":           c.Timeout,
		"retry_count":       c.RetryCount,
		"has_workflow":      c.HasWorkflowNodes(),
		"total_nodes":       len(c.WorkflowNodes),
		"execution_results": c.ExecutionResults,
		"variables":         c.Variables,
	}
}
