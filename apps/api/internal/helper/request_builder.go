package helper

import (
	"encoding/json"
	"fmt"
	"seo-backend/internal/domain/draft"
	"strings"
)

func (s *PostService) buildRequestBody(cfg *ProductConfig, draft draft.DraftDataPost) (map[string]interface{}, error) {
	fmt.Println("========== BUILD REQUEST BODY ==========")

	var fieldMapping map[string]interface{}

	// Parse field mapping
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

	// Build request body from field mapping
	requestBody := s.buildFromFieldMapping(fieldMapping, draft)

	// Add meta tags
	s.addMetaTagsToBody(cfg, draft, requestBody)

	// Add sitemap info
	s.addSitemapInfoToBody(cfg, draft, requestBody)

	// Add configs
	if cfg.MetaConfigStr != "" && cfg.MetaConfigStr != "{}" {
		requestBody["meta_config"] = removeSensitiveFields(parseConfig(cfg.MetaConfigStr))
	}
	if cfg.SitemapConfigStr != "" && cfg.SitemapConfigStr != "{}" {
		requestBody["sitemap_config"] = removeSensitiveFields(parseConfig(cfg.SitemapConfigStr))
	}

	// Default request body if empty
	if len(requestBody) == 0 {
		requestBody = map[string]interface{}{
			"title":   draft.Title,
			"topic":   draft.Topic,
			"content": draft.Article,
		}
		if draft.ImageURL != nil {
			requestBody["image_url"] = *draft.ImageURL
		}
	}

	// Log final body (redacted)
	redactedBody := redactSensitiveFields(requestBody)
	printJSON("FINAL REQUEST BODY (REDACTED)", redactedBody)

	return requestBody, nil
}

func parseFieldMapping(fieldMappingStr string, fieldMapping *map[string]interface{}) error {
	raw := strings.TrimSpace(fieldMappingStr)
	fmt.Println("RAW FIELD MAPPING (length):", len(raw))

	// Try direct JSON object
	if err := json.Unmarshal([]byte(raw), fieldMapping); err != nil {
		fmt.Println("DIRECT PARSE FAILED:", err)

		// Handle double encoded JSON string
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

func (s *PostService) buildFromFieldMapping(fieldMapping map[string]interface{}, draft draft.DraftDataPost) map[string]interface{} {
	requestBody := make(map[string]interface{})

	for key, value := range fieldMapping {
		fmt.Println("PROCESS FIELD:", key)

		switch v := value.(type) {
		case string:
			replaced := replaceAllPlaceholders(v, draft)
			if isSensitiveField(key) {
				fmt.Printf("STRING FIELD (sensitive) => %s : [REDACTED]\n", key)
			} else {
				fmt.Printf("STRING FIELD => %s : %v\n", key, truncateForLog(replaced, 100))
			}
			requestBody[key] = replaced

		case map[string]interface{}:
			fmt.Printf("NESTED FIELD => %s\n", key)
			nestedResult := make(map[string]interface{})

			for nestedKey, nestedValue := range v {
				if strVal, ok := nestedValue.(string); ok {
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

func (s *PostService) addMetaTagsToBody(cfg *ProductConfig, draft draft.DraftDataPost, requestBody map[string]interface{}) {
	metaConfig := parseConfig(cfg.MetaConfigStr)
	metaEnabled := false

	if metaConfig != nil {
		if enabled, ok := metaConfig["enabled"].(bool); ok && enabled {
			metaEnabled = true
			fmt.Println("META CONFIG ENABLED")
		}
	}

	if metaEnabled {
		sitemapConfig := parseConfig(cfg.SitemapConfigStr)
		metaTags := generateMetaTags(metaConfig, draft, cfg.BaseURL, sitemapConfig)
		fmt.Println("GENERATED META TAGS (count):", len(metaTags))

		if len(metaTags) > 0 {
			fmt.Println("INSERT meta_tags TO REQUEST BODY")
			requestBody["meta_tags"] = metaTags
		}
	} else {
		fmt.Println("SKIP META TAG GENERATION")
	}
}

func (s *PostService) addSitemapInfoToBody(cfg *ProductConfig, draft draft.DraftDataPost, requestBody map[string]interface{}) {
	sitemapConfig := parseConfig(cfg.SitemapConfigStr)
	sitemapEnabled := false

	if sitemapConfig != nil {
		if enabled, ok := sitemapConfig["enabled"].(bool); ok && enabled {
			sitemapEnabled = true
			fmt.Println("SITEMAP CONFIG ENABLED")
		}
	}

	if sitemapEnabled {
		sitemapInfo := generateSitemapInfo(sitemapConfig, draft, cfg.BaseURL)
		if sitemapInfo != nil {
			fmt.Println("GENERATED SITEMAP INFO (non-nil)")
			fmt.Println("INSERT sitemap_info TO REQUEST BODY")
			requestBody["sitemap_info"] = sitemapInfo
		} else {
			fmt.Println("GENERATED SITEMAP INFO IS NIL")
		}
	} else {
		fmt.Println("SKIP SITEMAP GENERATION")
	}
}

func parseConfig(configStr string) map[string]interface{} {
	if configStr == "" || configStr == "{}" {
		return nil
	}

	raw := strings.TrimSpace(configStr)
	var config map[string]interface{}

	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		// Try double encoded
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
