package helper

import (
	"encoding/json"
	"fmt"
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

	fmt.Println("========== BUILD REQUEST BODY ==========")
	fmt.Printf("Node ID: %s, Step: %d\n", node.ID, node.StepOrder)

	// 1. PARSE FIELD MAPPING
	var fieldMapping map[string]interface{}

	// Dari node
	if node.AdapterConfig.FieldMapping != "" && node.AdapterConfig.FieldMapping != "{}" {
		if err := json.Unmarshal([]byte(node.AdapterConfig.FieldMapping), &fieldMapping); err != nil {
			fmt.Printf("[WARN] Failed to parse node field mapping: %v\n", err)
			fieldMapping = nil
		} else {
			fmt.Println("✅ Using node field mapping")
		}
	}

	// Dari cfg kalo node kosong
	if len(fieldMapping) == 0 && cfg.FieldMappingStr != "" && cfg.FieldMappingStr != "{}" {
		if err := json.Unmarshal([]byte(cfg.FieldMappingStr), &fieldMapping); err != nil {
			fmt.Printf("[WARN] Failed to parse config field mapping: %v\n", err)
			fieldMapping = nil
		} else {
			fmt.Println("✅ Using config field mapping")
		}
	}

	fmt.Printf("FIELD MAPPING KEYS: %v\n", getMapKeys(fieldMapping))

	// 2. VALIDASI DRAFT
	if strings.TrimSpace(draft.Title) == "" {
		return nil, fmt.Errorf("draft title is required")
	}
	if strings.TrimSpace(draft.Article) == "" {
		return nil, fmt.Errorf("draft article content is required")
	}

	fmt.Printf("Draft Title: %s\n", draft.Title)
	fmt.Printf("Draft Topic: %s\n", draft.Topic)
	fmt.Printf("Draft Article Length: %d\n", len(draft.Article))

	// 3. BUILD REQUEST BODY
	requestBody := buildPayload(fieldMapping, draft, cfg)

	// 4. FALLBACK KALO KOSONG
	if len(requestBody) == 0 {
		fmt.Println("⚠️  FIELD MAPPING EMPTY, USING DEFAULT BODY")
		requestBody = map[string]interface{}{
			"title":   draft.Title,
			"topic":   draft.Topic,
			"content": draft.Article,
		}
		if draft.ImageURL != nil {
			requestBody["image_url"] = *draft.ImageURL
		}
	}

	printJSON("✅ FINAL REQUEST BODY", requestBody)

	return requestBody, nil
}

// ============================================================
// buildPayload - BUILD PAYLOAD DARI FIELD MAPPING
// ============================================================
func buildPayload(
	fieldMapping map[string]interface{},
	draft draft.DraftDataPost,
	cfg *product.ProductConfig,
) map[string]interface{} {

	payload := make(map[string]interface{})

	if len(fieldMapping) == 0 {
		return payload
	}

	// Convert draft ke map
	draftMap := structToMap(draft)

	// Tambahin hasil node sebelumnya
	if len(cfg.ExecutionResults) > 0 {
		draftMap["previous_results"] = cfg.ExecutionResults
	}
	if len(cfg.Variables) > 0 {
		draftMap["variables"] = cfg.Variables
	}

	// Loop mapping
	for targetField, sourceConfig := range fieldMapping {
		value := getValue(sourceConfig, draftMap, cfg)
		if value != nil {
			payload[targetField] = value
		}
	}

	return payload
}

// ============================================================
// getValue - AMBIL NILAI DARI SOURCE (SUPPORT TEMPLATE)
// ============================================================
func getValue(
	source interface{},
	draftMap map[string]interface{},
	cfg *product.ProductConfig,
) interface{} {

	switch v := source.(type) {
	case string:
		// TEMPLATE: {{node-id.field}}
		if strings.HasPrefix(v, "{{") && strings.HasSuffix(v, "}}") {
			val, err := cfg.ParseTemplate(v)
			if err != nil {
				fmt.Printf("[WARN] Template error '%s': %v\n", v, err)
				return nil
			}
			return val
		}

		// DARI DRAFT
		if val, exists := draftMap[v]; exists {
			return val
		}
		return nil

	case map[string]interface{}:
		// NESTED OBJECT
		result := make(map[string]interface{})
		for key, sub := range v {
			if val := getValue(sub, draftMap, cfg); val != nil {
				result[key] = val
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result

	case []interface{}:
		// ARRAY
		result := make([]interface{}, 0)
		for _, item := range v {
			if val := getValue(item, draftMap, cfg); val != nil {
				result = append(result, val)
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result

	default:
		return v
	}
}

func structToMap(data interface{}) map[string]interface{} {
	bytes, _ := json.Marshal(data)
	var result map[string]interface{}
	json.Unmarshal(bytes, &result)
	return result
}

// ─────────────────────────────────────────────
// Parse field mapping (handle double encoded)
// ─────────────────────────────────────────────

func parseFieldMapping(fieldMappingStr string, fieldMapping *map[string]interface{}) error {
	raw := strings.TrimSpace(fieldMappingStr)
	fmt.Println("RAW FIELD MAPPING (length):", len(raw))

	if err := json.Unmarshal([]byte(raw), fieldMapping); err != nil {
		fmt.Println("DIRECT PARSE FAILED:", err)

		var nested string
		if err2 := json.Unmarshal([]byte(raw), &nested); err2 == nil {
			fmt.Println("DOUBLE ENCODED DETECTED")
			if err3 := json.Unmarshal([]byte(nested), fieldMapping); err3 != nil {
				fmt.Println("NESTED PARSE FAILED:", err3)
				return fmt.Errorf("failed nested parse field mapping: %w", err3)
			}
		} else {
			fmt.Println("FIELD MAPPING PARSE FAILED:", err)
			return fmt.Errorf("failed to parse field mapping: %w", err)
		}
	}
	return nil
}

// ─────────────────────────────────────────────
// Build request body dari field mapping
// ─────────────────────────────────────────────

func buildFromFieldMapping(fieldMapping map[string]interface{}, draft draft.DraftDataPost) map[string]interface{} {
	requestBody := make(map[string]interface{})

	for key, value := range fieldMapping {
		fmt.Println("PROCESS FIELD:", key)

		switch v := value.(type) {
		case string:
			arrayPlaceholders := getArrayPlaceholders(draft)
			if arr, ok := arrayPlaceholders[v]; ok {
				fmt.Printf("ARRAY FIELD => %s : %v\n", key, arr)
				requestBody[key] = arr
				break
			}

			replaced := replaceAllPlaceholders(v, draft)
			if isSensitiveField(key) {
				fmt.Printf("STRING FIELD (sensitive) => %s : [REDACTED]\n", key)
			} else {
				fmt.Printf("STRING FIELD => %s : %v\n", key, truncateForLog(replaced, 100))
			}
			requestBody[key] = replaced

		case []interface{}:
			fmt.Printf("ARRAY FIELD => %s\n", key)
			result := make([]string, 0, len(v))
			for _, item := range v {
				if strVal, ok := item.(string); ok {
					result = append(result, replaceAllPlaceholders(strVal, draft))
				}
			}
			requestBody[key] = result

		case map[string]interface{}:
			fmt.Printf("NESTED FIELD => %s\n", key)
			nestedResult := make(map[string]interface{})

			for nestedKey, nestedValue := range v {
				if strVal, ok := nestedValue.(string); ok {
					arrayPlaceholders := getArrayPlaceholders(draft)
					if arr, ok := arrayPlaceholders[strVal]; ok {
						fmt.Printf("NESTED ARRAY => %s.%s : %v\n", key, nestedKey, arr)
						nestedResult[nestedKey] = arr
						continue
					}

					replaced := replaceAllPlaceholders(strVal, draft)
					if isSensitiveField(nestedKey) {
						fmt.Printf("NESTED STRING (sensitive) => %s.%s : [REDACTED]\n", key, nestedKey)
					} else {
						fmt.Printf("NESTED STRING => %s.%s : %v\n", key, nestedKey, truncateForLog(replaced, 100))
					}
					nestedResult[nestedKey] = replaced
				} else {
					fmt.Printf("NESTED RAW => %s.%s : %v\n", key, nestedKey, truncateForLog(fmt.Sprintf("%v", nestedValue), 100))
					nestedResult[nestedKey] = nestedValue
				}
			}
			requestBody[key] = nestedResult

		default:
			fmt.Printf("RAW FIELD => %s : %v\n", key, truncateForLog(fmt.Sprintf("%v", value), 100))
			requestBody[key] = value
		}
	}

	return requestBody
}

// ─────────────────────────────────────────────
// Add meta tags (seo) ke request body
// ─────────────────────────────────────────────

func (s *PostService) addMetaTagsToBody(
	metaConfig map[string]interface{},
	sitemapConfig map[string]interface{},
	draft draft.DraftDataPost,
	baseURL string,
	requestBody map[string]interface{},
) {
	if metaConfig == nil {
		fmt.Println("SKIP META TAG GENERATION: metaConfig nil")
		return
	}

	enabled, ok := metaConfig["enabled"].(bool)
	if !ok || !enabled {
		fmt.Println("SKIP META TAG GENERATION: disabled")
		return
	}

	fmt.Println("META CONFIG ENABLED")
	metaTags := generateMetaTags(metaConfig, draft, baseURL, sitemapConfig)
	fmt.Println("GENERATED META TAGS (count):", len(metaTags))

	if len(metaTags) > 0 {
		fmt.Println("INSERT seo TO REQUEST BODY")
		requestBody["seo"] = metaTags
	}
}

// ─────────────────────────────────────────────
// Add sitemap info ke request body
// ─────────────────────────────────────────────

func (s *PostService) addSitemapInfoToBody(
	sitemapConfig map[string]interface{},
	draft draft.DraftDataPost,
	baseURL string,
	requestBody map[string]interface{},
) {
	if sitemapConfig == nil {
		fmt.Println("SKIP SITEMAP GENERATION: sitemapConfig nil")
		return
	}

	enabled, ok := sitemapConfig["enabled"].(bool)
	if !ok || !enabled {
		fmt.Println("SKIP SITEMAP GENERATION: disabled")
		return
	}

	fmt.Println("SITEMAP CONFIG ENABLED")
	sitemapInfo := generateSitemapInfo(sitemapConfig, draft, baseURL)
	if sitemapInfo != nil {
		fmt.Println("GENERATED SITEMAP INFO (non-nil)")
		fmt.Println("INSERT sitemap TO REQUEST BODY")
		requestBody["sitemap"] = sitemapInfo
	} else {
		fmt.Println("GENERATED SITEMAP INFO IS NIL")
	}
}

// ─────────────────────────────────────────────
// Parse config JSON (handle double encoded)
// ─────────────────────────────────────────────

func parseConfig(configStr string) map[string]interface{} {
	if configStr == "" || configStr == "{}" {
		return nil
	}

	raw := strings.TrimSpace(configStr)
	var config map[string]interface{}

	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		var nested string
		if err2 := json.Unmarshal([]byte(raw), &nested); err2 == nil {
			fmt.Println("DOUBLE ENCODED CONFIG DETECTED")
			if err3 := json.Unmarshal([]byte(nested), &config); err3 != nil {
				fmt.Println("NESTED PARSE FAILED:", err3)
				return nil
			}
		} else {
			fmt.Println("CONFIG PARSE FAILED:", err)
			return nil
		}
	}

	return config
}
