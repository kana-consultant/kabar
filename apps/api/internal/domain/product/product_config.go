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

func (c *ProductConfig) HasWorkflowNodes() bool {
	return len(c.WorkflowNodes) > 0
}

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

func (c *ProductConfig) GetAllExecutionResults() map[string]interface{} {
	if c.ExecutionResults == nil {
		return make(map[string]interface{})
	}
	return c.ExecutionResults
}

func (c *ProductConfig) GetAllVariables() map[string]interface{} {
	if c.Variables == nil {
		return make(map[string]interface{})
	}
	return c.Variables
}

func (c *ProductConfig) ParseTemplate(template string) (interface{}, error) {
	if !strings.HasPrefix(template, "{{") || !strings.HasSuffix(template, "}}") {
		return nil, fmt.Errorf("invalid template format: %s", template)
	}

	trimmed := strings.TrimPrefix(strings.TrimSuffix(template, "}}"), "{{")
	parts := strings.SplitN(trimmed, ".", 2)

	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid template: %s (format: {{node-id.field.path}})", template)
	}

	nodeID := parts[0]
	fieldPath := parts[1]

	// Coba ambil dari execution results
	if c.HasExecutionResult(nodeID) {
		return c.GetExecutionResultField(nodeID, fieldPath)
	}

	// Coba ambil dari variables
	if c.HasVariable(nodeID) {
		val := c.GetVariable(nodeID)
		if fieldPath == "" {
			return val, nil
		}
		fieldParts := strings.Split(fieldPath, ".")
		current := val
		for _, part := range fieldParts {
			if currentMap, ok := current.(map[string]interface{}); ok {
				if v, exists := currentMap[part]; exists {
					current = v
				} else {
					return nil, fmt.Errorf("field '%s' not found", part)
				}
			} else {
				return nil, fmt.Errorf("cannot navigate to '%s'", part)
			}
		}
		return current, nil
	}

	return nil, fmt.Errorf("node or variable '%s' not found", nodeID)
}

// ============================================================
// 11. ✅ DIPAKAI - HasExecutionResult
// ============================================================
func (c *ProductConfig) HasExecutionResult(nodeID string) bool {
	if c.ExecutionResults == nil {
		return false
	}
	_, exists := c.ExecutionResults[nodeID]
	return exists
}

// ============================================================
// 12. ✅ DIPAKAI - HasVariable
// ============================================================
func (c *ProductConfig) HasVariable(key string) bool {
	if c.Variables == nil {
		return false
	}
	_, exists := c.Variables[key]
	return exists
}
