package product

import (
	"encoding/json"
	"fmt"
	"log"
	"seo-backend/internal/domain/workflow_node"
	"strings"
	"time"
)

// ProductConfig - VERSI SEDERHANA (NO VARIABLES, ONLY EXECUTION RESULTS)
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

	// Config strings
	FieldMappingStr  string
	MetaConfigStr    string
	SitemapConfigStr string
	CustomHeadersStr string

	// Parsed configs
	CustomHeaders map[string]string
	FieldMapping  []FieldMapping
	MetaConfig    map[string]interface{}
	SitemapConfig map[string]interface{}

	// Timeout and retry
	Timeout    int
	RetryCount int

	// Workflow nodes
	WorkflowNodes    []workflow_node.WorkflowNode
	WorkflowLevelMap map[int][]*workflow_node.WorkflowNode

	// Execution context (ONLY SOURCE OF TRUTH)
	ExecutionResults map[string]interface{}
	CurrentNodeID    string

	// Metadata
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *ProductConfig) HasWorkflowNodes() bool {
	return len(c.WorkflowNodes) > 0
}

// ============================================================
// EXECUTION RESULTS (ONLY STORAGE)
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

func (c *ProductConfig) HasExecutionResult(nodeID string) bool {
	if c.ExecutionResults == nil {
		return false
	}
	_, exists := c.ExecutionResults[nodeID]
	return exists
}

func (c *ProductConfig) GetAllExecutionResults() map[string]interface{} {
	if c.ExecutionResults == nil {
		return make(map[string]interface{})
	}
	return c.ExecutionResults
}

// ============================================================
// FIELD ACCESS (DOT NOTATION)
// ============================================================

func (c *ProductConfig) GetExecutionResultField(nodeID, fieldPath string) (interface{}, error) {
	result := c.GetExecutionResult(nodeID)
	if result == nil {
		return nil, fmt.Errorf("node %s result not found", nodeID)
	}

	// ===== PARSE []byte/string KE map[string]interface{} =====
	var parsed map[string]interface{}

	switch data := result.(type) {
	case []byte:
		if err := json.Unmarshal(data, &parsed); err != nil {
			return nil, fmt.Errorf("failed to parse node result ([]byte): %v", err)
		}
		result = parsed

	case string:
		if err := json.Unmarshal([]byte(data), &parsed); err != nil {
			// Bukan JSON, return as is
			parsed = nil
		} else {
			result = parsed
		}

	case map[string]interface{}:
		// Sudah dalam format yang benar
		break

	default:
		return nil, fmt.Errorf("unknown result type: %T", result)
	}
	// ===== END PARSE =====

	if fieldPath == "" {
		return result, nil
	}

	parts := strings.Split(fieldPath, ".")
	current := result

	for _, part := range parts {
		currentMap, ok := current.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("cannot navigate through non-object at '%s'", part)
		}

		val, exists := currentMap[part]
		if !exists {
			return nil, fmt.Errorf("field '%s' not found", part)
		}

		current = val
	}

	return current, nil
}

// ============================================================
// TEMPLATE ENGINE (SIMPLIFIED)
// ============================================================

func (c *ProductConfig) ParseTemplate(template string, node *workflow_node.WorkflowNode) (interface{}, error) {
	log.Printf("[INFO] ParseTemplate called with template: %s", template)

	if !strings.HasPrefix(template, "{") || !strings.HasSuffix(template, "}") {
		log.Printf("[ERROR] Template validation failed - invalid format: %s", template)
		return nil, fmt.Errorf("invalid template format: %s", template)
	}
	log.Printf("[DEBUG] Template validation passed")

	// Handle {{...}} format
	trimmed := template
	if strings.HasPrefix(template, "{{") && strings.HasSuffix(template, "}}") {
		trimmed = strings.TrimSuffix(strings.TrimPrefix(template, "{{"), "}}")
	} else {
		trimmed = strings.TrimSuffix(strings.TrimPrefix(template, "{"), "}")
	}
	log.Printf("[DEBUG] Trimmed template: %s", trimmed)

	parts := strings.SplitN(trimmed, ".", 2)
	nodeID := parts[0]
	fieldPath := ""
	if len(parts) > 1 {
		fieldPath = parts[1]
	}

	log.Printf("[DEBUG] Split result - nodeID: %s, fieldPath: %s", nodeID, fieldPath)

	// ===== HANDLE KEYWORD "response" =====
	if nodeID == "response" {
		if node == nil {
			log.Printf("[ERROR] 'response' keyword requires node context")
			return nil, fmt.Errorf("'response' keyword requires node context")
		}

		if len(node.PreviousNodeIDs) == 0 {
			log.Printf("[ERROR] No previous node for %s", node.ID)
			return nil, fmt.Errorf("no previous node for %s", node.ID)
		}

		found := false
		for _, prevID := range node.PreviousNodeIDs {
			prevResult, exists := c.ExecutionResults[prevID]
			if !exists {
				continue
			}

			// ===== TAMBAH INI =====
			var resultMap map[string]interface{}

			// Convert []uint8 ke map
			if byteData, ok := prevResult.([]byte); ok {
				if err := json.Unmarshal(byteData, &resultMap); err != nil {
					log.Printf("[DEBUG] Failed to parse prevResult as JSON: %v", err)
					continue
				}
			} else if byteData, ok := prevResult.([]uint8); ok {
				if err := json.Unmarshal(byteData, &resultMap); err != nil {
					log.Printf("[DEBUG] Failed to parse prevResult as JSON: %v", err)
					continue
				}
			} else if mapData, ok := prevResult.(map[string]interface{}); ok {
				resultMap = mapData
			} else if strData, ok := prevResult.(string); ok {
				if err := json.Unmarshal([]byte(strData), &resultMap); err != nil {
					log.Printf("[DEBUG] Failed to parse prevResult string as JSON: %v", err)
					continue
				}
			} else {
				log.Printf("[DEBUG] Unknown prevResult type: %T", prevResult)
				continue
			}
			// ===== SAMPAI SINI =====

			// Sekarang cek field
			if _, fieldExists := resultMap[fieldPath]; fieldExists {
				nodeID = prevID
				found = true
				log.Printf("[DEBUG] 'response' resolved: %s -> %s (field: %s)",
					node.ID, nodeID, fieldPath)
				break
			}
		}

		if !found {
			log.Printf("[ERROR] Field '%s' not found in any previous node for %s", fieldPath, node.ID)
			return nil, fmt.Errorf("field '%s' not found in any previous node", fieldPath)
		}
	}
	// ===== END HANDLE KEYWORD "response" =====

	log.Printf("[INFO] Calling GetExecutionResultField with nodeID: %s, fieldPath: %s", nodeID, fieldPath)

	result, err := c.GetExecutionResultField(nodeID, fieldPath)
	if err != nil {
		log.Printf("[ERROR] GetExecutionResultField returned error: %v", err)
		return nil, err
	}

	log.Printf("[INFO] GetExecutionResultField returned value: %v (type: %T)", result, result)
	log.Printf("[SUCCESS] ParseTemplate completed successfully")
	return result, nil
}

func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
