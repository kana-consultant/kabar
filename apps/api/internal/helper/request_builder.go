package helper

import (
	"encoding/json"
	"fmt"
	"log"
	"seo-backend/internal/domain/draft"
	"seo-backend/internal/domain/product"
	"seo-backend/internal/domain/workflow_node"
	"strings"

	"github.com/google/uuid"
)

func BuildRequestBody(
	node *workflow_node.WorkflowNode,
	draft draft.DraftDataPost,
	cfg *product.ProductConfig,
) (map[string]interface{}, error) {

	var fieldMapping map[string]interface{}

	// Parse field mapping dari node
	if node.AdapterConfig.FieldMapping != "" {
		trimmedMapping := strings.TrimSpace(node.AdapterConfig.FieldMapping)
		if trimmedMapping != "" && trimmedMapping != "{}" {
			if err := json.Unmarshal([]byte(node.AdapterConfig.FieldMapping), &fieldMapping); err != nil {
				log.Printf("[WARNING] Failed to parse node field mapping: %v", err)
				fieldMapping = nil
			}
		}
	}

	// Fallback ke cfg
	if fieldMapping == nil && cfg.FieldMappingStr != "" {
		trimmedMapping := strings.TrimSpace(cfg.FieldMappingStr)
		if trimmedMapping != "" && trimmedMapping != "{}" {
			if err := json.Unmarshal([]byte(cfg.FieldMappingStr), &fieldMapping); err != nil {
				log.Printf("[WARNING] Failed to parse config field mapping: %v", err)
				fieldMapping = nil
			}
		}
	}

	if fieldMapping == nil {
		fieldMapping = make(map[string]interface{})
	}

	// Validasi draft
	if strings.TrimSpace(draft.Title) == "" {
		return nil, fmt.Errorf("draft title is required")
	}
	if strings.TrimSpace(draft.Article) == "" {
		return nil, fmt.Errorf("draft article content is required")
	}

	// Build request body
	requestBody := buildPayload(fieldMapping, draft, cfg, node)

	// Fallback jika kosong
	if len(requestBody) == 0 {
		log.Printf("[WARNING] Empty payload, using default body for product %s", cfg.ProductID)
		requestBody = map[string]interface{}{
			"id":       draft.Id,
			"title":    draft.Title,
			"topic":    draft.Topic,
			"content":  draft.Article,
			"slug":     draft.Slug,
			"excerpt":  draft.Excerpt,
			"keywords": draft.Keywords,
		}
		if draft.ImageURL != nil {
			requestBody["image_url"] = *draft.ImageURL
		}
		if draft.ImagePrompt != "" {
			requestBody["image_prompt"] = draft.ImagePrompt
		}
	}

	return requestBody, nil
}

func buildPayload(
	fieldMapping map[string]interface{},
	draft draft.DraftDataPost,
	cfg *product.ProductConfig,
	node *workflow_node.WorkflowNode, // TAMBAH
) map[string]interface{} {

	payload := make(map[string]interface{})

	if len(fieldMapping) == 0 {
		return payload
	}

	draftMap := structToMap(draft)

	for targetField, sourceConfig := range fieldMapping {
		value, found := getValue(sourceConfig, draftMap, cfg, node) // TAMBAH node
		if found {
			payload[targetField] = value
		}
	}

	return payload
}

func getValue(
	source interface{},
	draftMap map[string]interface{},
	cfg *product.ProductConfig,
	node *workflow_node.WorkflowNode,
) (interface{}, bool) {

	switch v := source.(type) {
	case string:
		// AUTO-GENERATE UUID untuk key "id"
		if v == "id" || v == "{{id}}" || v == "{id}" {
			uuid := generateUUID()
			log.Printf("[SUCCESS] getValue: Auto-generated UUID for 'id' -> '%s'", uuid)
			return uuid, true
		}

		// Normalisasi template
		normalized := v
		if strings.HasPrefix(v, "{{") && strings.HasSuffix(v, "}}") {
			normalized = v[2 : len(v)-2]
		} else if strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}") {
			normalized = v[1 : len(v)-1]
		}

		// Cek lagi untuk normalized "id"
		if normalized == "id" {
			uuid := generateUUID()
			log.Printf("[SUCCESS] getValue: Auto-generated UUID for 'id' (normalized) -> '%s'", uuid)
			return uuid, true
		}

		// Cek template {{...}} atau {...}
		if (strings.HasPrefix(v, "{{") && strings.HasSuffix(v, "}}")) ||
			(strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}")) {

			// Coba resolve via ProductConfig (workflow execution results)
			if cfg != nil {
				val, err := cfg.ParseTemplate(v, node)
				if err == nil && val != nil {
					log.Printf("[SUCCESS] getValue: Template '%s' resolved via cfg -> '%v'", v, val)
					return val, true
				}
				if err != nil {
					log.Printf("[WARNING] getValue: Template parse failed for '%s': %v", normalized, err)
				}
			}
		}

		// Alias content -> article (untuk lookup di draftMap)
		lookupKey := normalized
		if normalized == "content" {
			lookupKey = "article"
		}

		// Lookup di draftMap
		val, exists := draftMap[lookupKey]
		if exists {
			// Khusus content: strip HTML tags
			if normalized == "content" {
				if strVal, ok := val.(string); ok {
					stripped := stripHTMLTags(strVal)
					log.Printf("[SUCCESS] getValue: DraftMap lookup '%s' (alias: '%s') -> content stripped (%d chars)",
						normalized, lookupKey, len(stripped))
					return stripped, true
				}
			}
			log.Printf("[SUCCESS] getValue: DraftMap lookup '%s' -> '%v'", normalized, val)
			return val, true
		}

		// Not found
		log.Printf("[DEBUG] getValue: Not found for '%s' (lookup key: '%s')", normalized, lookupKey)
		return nil, false

	case map[string]interface{}:
		log.Printf("[DEBUG] getValue: Processing map with %d keys", len(v))
		result := make(map[string]interface{})
		hasAny := false

		for key, sub := range v {
			if val, found := getValue(sub, draftMap, cfg, node); found {
				result[key] = val
				hasAny = true
				log.Printf("[SUCCESS] getValue: Map key '%s' resolved -> '%v'", key, val)
			} else {
				log.Printf("[DEBUG] getValue: Map key '%s' NOT resolved", key)
			}
		}

		if !hasAny {
			log.Printf("[DEBUG] getValue: Map returned empty, no keys resolved")
			return nil, false
		}
		log.Printf("[SUCCESS] getValue: Map resolved with %d/%d keys", len(result), len(v))
		return result, true

	case []interface{}:
		log.Printf("[DEBUG] getValue: Processing array with %d items", len(v))
		result := make([]interface{}, 0, len(v))
		hasAny := false

		for i, item := range v {
			if val, found := getValue(item, draftMap, cfg, node); found {
				result = append(result, val)
				hasAny = true
				log.Printf("[SUCCESS] getValue: Array[%d] resolved -> '%v'", i, val)
			} else {
				log.Printf("[DEBUG] getValue: Array[%d] NOT resolved", i)
			}
		}

		if !hasAny {
			log.Printf("[DEBUG] getValue: Array returned empty, no items resolved")
			return nil, false
		}
		log.Printf("[SUCCESS] getValue: Array resolved with %d/%d items", len(result), len(v))
		return result, true

	case nil:
		log.Printf("[DEBUG] getValue: Source is nil")
		return nil, false

	default:
		log.Printf("[SUCCESS] getValue: Default type '%T' -> '%v'", v, v)
		return v, true
	}
}

// Fungsi helper untuk generate UUID
func generateUUID() string {
	return uuid.New().String()
}

func structToMap(data interface{}) map[string]interface{} {
	bytes, _ := json.Marshal(data)
	var result map[string]interface{}
	json.Unmarshal(bytes, &result)
	return result
}

func parseFieldMapping(fieldMappingStr string, fieldMapping *map[string]interface{}) error {
	raw := strings.TrimSpace(fieldMappingStr)

	if err := json.Unmarshal([]byte(raw), fieldMapping); err != nil {
		var nested string
		if err2 := json.Unmarshal([]byte(raw), &nested); err2 == nil {
			if err3 := json.Unmarshal([]byte(nested), fieldMapping); err3 != nil {
				return fmt.Errorf("failed nested parse field mapping: %w", err3)
			}
		} else {
			return fmt.Errorf("failed to parse field mapping: %w", err)
		}
	}
	return nil
}

func buildFromFieldMapping(fieldMapping map[string]interface{}, draft draft.DraftDataPost) map[string]interface{} {
	requestBody := make(map[string]interface{})

	for key, value := range fieldMapping {
		switch v := value.(type) {
		case string:
			arrayPlaceholders := getArrayPlaceholders(draft)
			if arr, ok := arrayPlaceholders[v]; ok {
				requestBody[key] = arr
				break
			}
			requestBody[key] = replaceAllPlaceholders(v, draft)

		case []interface{}:
			result := make([]string, 0, len(v))
			for _, item := range v {
				if strVal, ok := item.(string); ok {
					result = append(result, replaceAllPlaceholders(strVal, draft))
				}
			}
			requestBody[key] = result

		case map[string]interface{}:
			nestedResult := make(map[string]interface{})
			for nestedKey, nestedValue := range v {
				if strVal, ok := nestedValue.(string); ok {
					arrayPlaceholders := getArrayPlaceholders(draft)
					if arr, ok := arrayPlaceholders[strVal]; ok {
						nestedResult[nestedKey] = arr
						continue
					}
					nestedResult[nestedKey] = replaceAllPlaceholders(strVal, draft)
				} else {
					nestedResult[nestedKey] = nestedValue
				}
			}
			requestBody[key] = nestedResult

		default:
			requestBody[key] = value
		}
	}

	return requestBody
}

func (s *PostService) addMetaTagsToBody(
	metaConfig map[string]interface{},
	sitemapConfig map[string]interface{},
	draft draft.DraftDataPost,
	baseURL string,
	requestBody map[string]interface{},
) {
	if metaConfig == nil {
		return
	}

	enabled, ok := metaConfig["enabled"].(bool)
	if !ok || !enabled {
		return
	}

	metaTags := generateMetaTags(metaConfig, draft, baseURL, sitemapConfig)
	if len(metaTags) > 0 {
		requestBody["seo"] = metaTags
	}
}

func (s *PostService) addSitemapInfoToBody(
	sitemapConfig map[string]interface{},
	draft draft.DraftDataPost,
	baseURL string,
	requestBody map[string]interface{},
) {
	if sitemapConfig == nil {
		return
	}

	enabled, ok := sitemapConfig["enabled"].(bool)
	if !ok || !enabled {
		return
	}

	sitemapInfo := generateSitemapInfo(sitemapConfig, draft, baseURL)
	if sitemapInfo != nil {
		requestBody["sitemap"] = sitemapInfo
	}
}

func parseConfig(configStr string) map[string]interface{} {
	if configStr == "" || configStr == "{}" {
		return nil
	}

	raw := strings.TrimSpace(configStr)
	var config map[string]interface{}

	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		var nested string
		if err2 := json.Unmarshal([]byte(raw), &nested); err2 == nil {
			if err3 := json.Unmarshal([]byte(nested), &config); err3 != nil {
				return nil
			}
		} else {
			return nil
		}
	}

	return config
}
