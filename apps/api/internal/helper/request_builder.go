package helper

import (
	"encoding/json"
	"fmt"
	"seo-backend/internal/domain/draft"
	"strings"
)

func (s *PostService) buildRequestBody(cfg *ProductConfig, draft draft.DraftDataPost) (map[string]interface{}, error) {
	fmt.Println("========== BUILD REQUEST BODY ==========")

	// Parse field mapping
	var fieldMapping map[string]interface{}
	if cfg.FieldMappingStr != "" && cfg.FieldMappingStr != "{}" {
		if err := parseFieldMapping(cfg.FieldMappingStr, &fieldMapping); err != nil {
			return nil, err
		}
	}
	fmt.Println("PARSED FIELD MAPPING KEYS:", getMapKeys(fieldMapping))

	// Validate draft
	if strings.TrimSpace(draft.Title) == "" {
		fmt.Println("ERROR: draft title empty")
		return nil, fmt.Errorf("draft title is required")
	}
	if strings.TrimSpace(draft.Article) == "" {
		fmt.Println("ERROR: draft article empty")
		return nil, fmt.Errorf("draft article content is required")
	}

	fmt.Println("DRAFT TITLE (length):", len(draft.Title), "characters")
	fmt.Println("DRAFT TOPIC:", draft.Topic)
	fmt.Println("DRAFT ARTICLE (length):", len(draft.Article), "characters")

	// Parse config sekali di sini, reuse ke semua fungsi
	metaConfig := parseConfig(cfg.MetaConfigStr)
	sitemapConfig := parseConfig(cfg.SitemapConfigStr)

	// Build request body from field mapping
	requestBody := s.buildFromFieldMapping(fieldMapping, draft)

	// ✅ Fallback SEBELUM tambah meta/sitemap
	// Supaya title/topic/content tetap masuk kalau field mapping kosong
	if len(requestBody) == 0 {
		fmt.Println("FIELD MAPPING EMPTY, USING DEFAULT BODY")
		requestBody = map[string]interface{}{
			"title":   draft.Title,
			"topic":   draft.Topic,
			"content": draft.Article,
		}
		if draft.ImageURL != nil {
			requestBody["image_url"] = *draft.ImageURL
		}
	}

	// Add meta tags (seo)
	s.addMetaTagsToBody(metaConfig, sitemapConfig, draft, cfg.BaseURL, requestBody)

	// Add sitemap info
	s.addSitemapInfoToBody(sitemapConfig, draft, cfg.BaseURL, requestBody)

	// Log final body (redacted)
	redactedBody := redactSensitiveFields(requestBody)
	printJSON("FINAL REQUEST BODY (REDACTED)", redactedBody)

	return requestBody, nil
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

func (s *PostService) buildFromFieldMapping(fieldMapping map[string]interface{}, draft draft.DraftDataPost) map[string]interface{} {
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
