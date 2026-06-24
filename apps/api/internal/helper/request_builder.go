package helper

import (
	"encoding/json"
	"fmt"
	"log"
	"seo-backend/internal/domain/draft"
	"seo-backend/internal/domain/product"
	"seo-backend/internal/domain/workflow_node"
	"strings"
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
	requestBody := buildPayload(fieldMapping, draft, cfg)

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
) map[string]interface{} {

	payload := make(map[string]interface{})

	if len(fieldMapping) == 0 {
		return payload
	}

	draftMap := structToMap(draft)

	for targetField, sourceConfig := range fieldMapping {
		value, found := getValue(sourceConfig, draftMap, cfg)
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
) (interface{}, bool) {

	switch v := source.(type) {
	case string:
		// Normalisasi template
		normalized := v
		if strings.HasPrefix(v, "{{") && strings.HasSuffix(v, "}}") {
			normalized = v[2 : len(v)-2]
		} else if strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}") {
			normalized = v[1 : len(v)-1]
		}

		// Cek template
		if (strings.HasPrefix(v, "{{") && strings.HasSuffix(v, "}}")) ||
			(strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}")) {

			if cfg != nil {
				val, err := cfg.ParseTemplate(v)
				if err == nil && val != nil {
					return val, true
				}
				if err != nil {
					log.Printf("[WARNING] Template parse failed for '%s': %v", normalized, err)
				}
			}
		}

		// Alias content -> article
		lookupKey := normalized
		if normalized == "content" {
			lookupKey = "article"
		}

		// Lookup di draft
		val, exists := draftMap[lookupKey]
		if exists {
			// Khusus content: strip HTML
			if normalized == "content" {
				if strVal, ok := val.(string); ok {
					return stripHTMLTags(strVal), true
				}
			}
			return val, true
		}

		return nil, false

	case map[string]interface{}:
		result := make(map[string]interface{})
		hasAny := false

		for key, sub := range v {
			if val, found := getValue(sub, draftMap, cfg); found {
				result[key] = val
				hasAny = true
			}
		}

		if !hasAny {
			return nil, false
		}
		return result, true

	case []interface{}:
		result := make([]interface{}, 0, len(v))
		hasAny := false

		for _, item := range v {
			if val, found := getValue(item, draftMap, cfg); found {
				result = append(result, val)
				hasAny = true
			}
		}

		if !hasAny {
			return nil, false
		}
		return result, true

	case nil:
		return nil, false

	default:
		return v, true
	}
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
