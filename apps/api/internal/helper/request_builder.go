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

	logJSON("build_request_start", map[string]interface{}{
		"node_id":    node.ID,
		"step":       node.StepOrder,
		"product_id": cfg.ProductID,
	})

	// ========== LOG ALL DRAFT FIELDS ==========
	logJSON("draft_all_fields", map[string]interface{}{
		"id":             draft.Id,
		"title":          draft.Title,
		"topic":          draft.Topic,
		"article":        draft.Article,
		"article_length": len(draft.Article),
		"image_url": func() interface{} {
			if draft.ImageURL != nil {
				return *draft.ImageURL
			}
			return nil
		}(),
		"image_prompt":          draft.ImagePrompt,
		"target_products":       draft.TargetProducts,
		"target_products_count": len(draft.TargetProducts),
		"slug":                  draft.Slug,
		"keywords":              draft.Keywords,
		"keywords_count":        len(draft.Keywords),
		"excerpt":               draft.Excerpt,
		"excerpt_length":        len(draft.Excerpt),
	})

	// ========== LOG DRAFT SUMMARY ==========
	logJSON("draft_summary", map[string]interface{}{
		"has_title":           draft.Title != "",
		"has_topic":           draft.Topic != "",
		"has_article":         draft.Article != "",
		"has_image":           draft.ImageURL != nil,
		"has_image_prompt":    draft.ImagePrompt != "",
		"has_slug":            draft.Slug != "",
		"has_keywords":        len(draft.Keywords) > 0,
		"has_excerpt":         draft.Excerpt != "",
		"has_target_products": len(draft.TargetProducts) > 0,
	})

	// ========== LOG DRAFT DETAILS ==========
	if len(draft.Keywords) > 0 {
		logJSON("draft_keywords", map[string]interface{}{
			"keywords": draft.Keywords,
			"count":    len(draft.Keywords),
		})
	}

	if draft.ImageURL != nil {
		logJSON("draft_image", map[string]interface{}{
			"image_url":    *draft.ImageURL,
			"image_prompt": draft.ImagePrompt,
		})
	}

	if draft.Excerpt != "" {
		logJSON("draft_excerpt", map[string]interface{}{
			"excerpt": draft.Excerpt,
			"length":  len(draft.Excerpt),
		})
	}

	if len(draft.TargetProducts) > 0 {
		logJSON("draft_target_products", map[string]interface{}{
			"target_products": draft.TargetProducts,
			"count":           len(draft.TargetProducts),
		})
	}

	// ========== LOG FIELD MAPPING ==========
	var fieldMapping map[string]interface{}

	// Dari node
	if node.AdapterConfig.FieldMapping != "" {
		trimmedMapping := strings.TrimSpace(node.AdapterConfig.FieldMapping)
		if trimmedMapping != "" && trimmedMapping != "{}" {
			if err := json.Unmarshal([]byte(node.AdapterConfig.FieldMapping), &fieldMapping); err != nil {
				logJSON("field_mapping_parse_error", map[string]interface{}{
					"source": "node",
					"error":  err.Error(),
					"raw":    node.AdapterConfig.FieldMapping,
				})
				fieldMapping = nil
			} else {
				logJSON("field_mapping_from_node", map[string]interface{}{
					"mapping": fieldMapping,
					"keys":    getMapKeys(fieldMapping),
				})
			}
		} else {
			logJSON("field_mapping_empty", map[string]interface{}{
				"source": "node",
				"reason": "empty or '{}'",
			})
		}
	}

	// Dari cfg kalo node kosong atau nil
	if fieldMapping == nil && cfg.FieldMappingStr != "" {
		trimmedMapping := strings.TrimSpace(cfg.FieldMappingStr)
		if trimmedMapping != "" && trimmedMapping != "{}" {
			if err := json.Unmarshal([]byte(cfg.FieldMappingStr), &fieldMapping); err != nil {
				logJSON("field_mapping_parse_error", map[string]interface{}{
					"source": "config",
					"error":  err.Error(),
					"raw":    cfg.FieldMappingStr,
				})
				fieldMapping = nil
			} else {
				logJSON("field_mapping_from_config", map[string]interface{}{
					"mapping": fieldMapping,
					"keys":    getMapKeys(fieldMapping),
				})
			}
		} else {
			logJSON("field_mapping_empty", map[string]interface{}{
				"source": "config",
				"reason": "empty or '{}'",
			})
		}
	}

	// Final check
	if fieldMapping == nil {
		logJSON("field_mapping_not_found", map[string]interface{}{
			"message": "No field mapping available, using default",
		})
		fieldMapping = make(map[string]interface{})
	}

	// 2. VALIDASI DRAFT
	if strings.TrimSpace(draft.Title) == "" {
		logJSON("validation_error", map[string]interface{}{
			"field": "title",
			"error": "draft title is required",
			"draft": draft,
		})
		return nil, fmt.Errorf("draft title is required")
	}
	if strings.TrimSpace(draft.Article) == "" {
		logJSON("validation_error", map[string]interface{}{
			"field": "article",
			"error": "draft article content is required",
			"draft": draft,
		})
		return nil, fmt.Errorf("draft article content is required")
	}

	// ========== BUILD REQUEST BODY ==========
	requestBody := buildPayload(fieldMapping, draft, cfg)

	// Log hasil build
	logJSON("request_body_built", map[string]interface{}{
		"fields_count": len(requestBody),
		"keys":         getMapKeys(requestBody),
		"body":         requestBody,
	})

	// 4. FALLBACK KALO KOSONG
	if len(requestBody) == 0 {
		logJSON("fallback_used", map[string]interface{}{
			"reason": "field mapping empty, using default body",
			"draft": map[string]interface{}{
				"id":       draft.Id,
				"title":    draft.Title,
				"topic":    draft.Topic,
				"content":  draft.Article,
				"slug":     draft.Slug,
				"excerpt":  draft.Excerpt,
				"keywords": draft.Keywords,
			},
		})
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

	// ========== LOG FINAL REQUEST BODY ==========
	logJSON("final_request_body", map[string]interface{}{
		"node_id":     node.ID,
		"product_id":  cfg.ProductID,
		"body_fields": len(requestBody),
		"body_keys":   getMapKeys(requestBody),
		"body":        requestBody,
	})

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

	log.Printf("=== BUILD PAYLOAD START ===")
	log.Printf("Field Mapping: %+v", fieldMapping)

	if len(fieldMapping) == 0 {
		log.Printf("Field Mapping kosong")
		return payload
	}

	// Convert draft ke map
	draftMap := structToMap(draft)

	log.Printf("Draft Map: %+v", draftMap)

	// Loop mapping
	for targetField, sourceConfig := range fieldMapping {

		log.Printf("Processing targetField=%s", targetField)
		log.Printf("Source Config=%+v", sourceConfig)

		value, found := getValue(sourceConfig, draftMap, cfg)

		log.Printf(
			"Result targetField=%s found=%v value=%#v",
			targetField,
			found,
			value,
		)

		if found {
			payload[targetField] = value

			log.Printf(
				"Added payload[%s] = %#v",
				targetField,
				value,
			)
		}
	}

	log.Printf("Final Payload: %+v", payload)
	log.Printf("=== BUILD PAYLOAD END ===")

	return payload
}

// ============================================================
// getValue - AMBIL NILAI DARI SOURCE (SUPPORT TEMPLATE)
// Return (value, found):
//   - found=false artinya "tidak ada nilai untuk dimasukkan ke payload"
//     (key tidak ditemukan, template error, atau nested map/array kosong)
//   - found=true artinya nilai valid untuk dimasukkan, TERMASUK kalau
//     nilainya nil/false/0/"" secara eksplisit.
//
// ============================================================
func getValue(
	source interface{},
	draftMap map[string]interface{},
	cfg *product.ProductConfig,
) (interface{}, bool) {

	logJSON("get_value_start", map[string]interface{}{
		"source_type":  fmt.Sprintf("%T", source),
		"source_value": fmt.Sprintf("%#v", source),
	})

	switch v := source.(type) {
	case string:
		logJSON("get_value_string", map[string]interface{}{
			"value": v,
		})

		// NORMALISASI: strip {} atau {{}}
		normalized := v
		if strings.HasPrefix(v, "{{") && strings.HasSuffix(v, "}}") {
			normalized = v[2 : len(v)-2]
			logJSON("get_value_normalized", map[string]interface{}{
				"original":   v,
				"normalized": normalized,
				"type":       "double_curly",
			})
		} else if strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}") {
			normalized = v[1 : len(v)-1]
			logJSON("get_value_normalized", map[string]interface{}{
				"original":   v,
				"normalized": normalized,
				"type":       "single_curly",
			})
		}

		// CEK APAKAH TEMPLATE
		if (strings.HasPrefix(v, "{{") && strings.HasSuffix(v, "}}")) || (strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}")) {
			logJSON("template_detected", map[string]interface{}{
				"template": v,
				"key":      normalized,
			})

			if cfg == nil {
				logJSON("template_error", map[string]interface{}{
					"error": "cfg is nil",
				})
				return nil, false
			}

			// PARSE TEMPLATE
			val, err := cfg.ParseTemplate(v)
			logJSON("template_parse_result", map[string]interface{}{
				"key":     normalized,
				"value":   val,
				"error":   err,
				"success": err == nil && val != nil,
			})

			// JIKA BERHASIL, KEMBALIKAN HASIL PARSE
			if err == nil && val != nil {
				logJSON("template_success", map[string]interface{}{
					"key":   normalized,
					"value": val,
				})
				return val, true
			}

			if err != nil {
				logJSON("template_error_fallback", map[string]interface{}{
					"key":   normalized,
					"error": err.Error(),
				})
			} else {
				logJSON("template_nil_fallback", map[string]interface{}{
					"key": normalized,
				})
			}
		}

		// ✅ ALIAS MAPPING: content -> article
		// Jika key adalah "content", kita akan mengambil dari "article"
		lookupKey := normalized
		if normalized == "content" {
			lookupKey = "article"
			logJSON("alias_mapping", map[string]interface{}{
				"alias":  "content",
				"target": "article",
				"reason": "content is alias for article (HTML content)",
			})
		}

		// ✅ FALLBACK KE DRAFT MAP (pakai lookupKey)
		logJSON("draft_lookup", map[string]interface{}{
			"original_key": normalized,
			"lookup_key":   lookupKey,
		})

		val, exists := draftMap[lookupKey]
		logJSON("draft_result", map[string]interface{}{
			"original_key": normalized,
			"lookup_key":   lookupKey,
			"exists":       exists,
			"value_type":   fmt.Sprintf("%T", val),
		})

		if exists {
			// ✅ KHUSUS UNTUK {content} - Ambil plain text tanpa HTML
			if normalized == "content" {
				if strVal, ok := val.(string); ok {
					// Hapus HTML tags untuk content
					plainText := stripHTMLTags(strVal)
					logJSON("content_plain_text", map[string]interface{}{
						"original_key":    normalized,
						"source_key":      lookupKey,
						"original_length": len(strVal),
						"plain_length":    len(plainText),
						"plain_preview":   plainText[:min(len(plainText), 100)],
					})
					return plainText, true
				}
			}

			// ✅ KHUSUS UNTUK {article} - Kembalikan HTML content asli
			if normalized == "article" {
				if strVal, ok := val.(string); ok {
					logJSON("article_html", map[string]interface{}{
						"html_length":  len(strVal),
						"html_preview": strVal[:min(len(strVal), 100)],
					})
					return strVal, true
				}
			}

			return val, true
		}

		logJSON("key_not_found", map[string]interface{}{
			"original_key": normalized,
			"lookup_key":   lookupKey,
		})
		return nil, false

	case map[string]interface{}:
		logJSON("map_detected", map[string]interface{}{
			"keys": getMapKeys(v),
			"len":  len(v),
		})

		result := make(map[string]interface{})
		hasAny := false

		for key, sub := range v {
			logJSON("map_iteration", map[string]interface{}{
				"key": key,
			})

			if val, found := getValue(sub, draftMap, cfg); found {
				result[key] = val
				hasAny = true
			}
		}

		logJSON("map_result", map[string]interface{}{
			"has_any": hasAny,
			"result":  result,
			"keys":    getMapKeys(result),
		})

		if !hasAny {
			return nil, false
		}
		return result, true

	case []interface{}:
		logJSON("array_detected", map[string]interface{}{
			"len": len(v),
		})

		result := make([]interface{}, 0, len(v))
		hasAny := false

		for i, item := range v {
			logJSON("array_iteration", map[string]interface{}{
				"index": i,
			})

			if val, found := getValue(item, draftMap, cfg); found {
				result = append(result, val)
				hasAny = true
			}
		}

		logJSON("array_result", map[string]interface{}{
			"has_any": hasAny,
			"result":  result,
			"len":     len(result),
		})

		if !hasAny {
			return nil, false
		}
		return result, true

	case nil:
		logJSON("nil_source", map[string]interface{}{
			"message": "source is nil",
		})
		return nil, false

	default:
		logJSON("literal_value", map[string]interface{}{
			"type":  fmt.Sprintf("%T", v),
			"value": v,
		})
		return v, true
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
