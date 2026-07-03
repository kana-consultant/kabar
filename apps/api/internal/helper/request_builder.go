package helper

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"seo-backend/internal/domain/draft"
	"seo-backend/internal/domain/product"
	"seo-backend/internal/domain/workflow_node"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

// ExtractFirstImageFromHTML mengambil URL gambar pertama dari HTML menggunakan goquery
func ExtractFirstImageFromHTML(htmlContent string) string {
	if htmlContent == "" {
		return ""
	}

	// Parse HTML dengan goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Printf("[ERROR] Failed to parse HTML: %v", err)
		return ""
	}

	// Cari tag img pertama
	var imageURL string
	doc.Find("img").First().Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists && src != "" {
			imageURL = src
			return
		}
	})

	// Jika tidak ada img tag, coba cari markdown image
	if imageURL == "" {
		// goquery tidak support markdown, kita tetap pakai regex untuk markdown
		reMarkdown := regexp.MustCompile(`!\[.*?\]\(([^)]+)\)`)
		matches := reMarkdown.FindStringSubmatch(htmlContent)
		if len(matches) > 1 {
			imageURL = matches[1]
		}
	}

	return imageURL
}

// ExtractAllImagesFromHTML mengambil semua URL gambar dari HTML
func ExtractAllImagesFromHTML(htmlContent string) []string {
	if htmlContent == "" {
		return []string{}
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Printf("[ERROR] Failed to parse HTML: %v", err)
		return []string{}
	}

	var images []string
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists && src != "" {
			images = append(images, src)
		}
	})

	// Jika tidak ada img tag, coba markdown
	if len(images) == 0 {
		reMarkdown := regexp.MustCompile(`!\[.*?\]\(([^)]+)\)`)
		matches := reMarkdown.FindAllStringSubmatch(htmlContent, -1)
		for _, match := range matches {
			if len(match) > 1 {
				images = append(images, match[1])
			}
		}
	}

	return images
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

		// 🔥 NEW: Handle image_url extraction from article using goquery
		if normalized == "image_url" || normalized == "image" || normalized == "img" {
			// Cek apakah ada di draftMap langsung
			if val, exists := draftMap["image_url"]; exists && val != nil {
				if strVal, ok := val.(string); ok && strVal != "" {
					log.Printf("[SUCCESS] getValue: image_url found directly -> '%s'", strVal)
					return strVal, true
				}
			}

			// 🔥 Ambil dari article dan extract gambar dengan goquery
			if articleVal, exists := draftMap["article"]; exists {
				if articleStr, ok := articleVal.(string); ok && articleStr != "" {
					imageURL := ExtractFirstImageFromHTML(articleStr)
					if imageURL != "" {
						log.Printf("[SUCCESS] getValue: image_url extracted from article with goquery -> '%s'", imageURL)
						return imageURL, true
					}
					log.Printf("[DEBUG] getValue: No image found in article for image_url")
				}
			}

			// 🔥 Ambil dari content (alias article)
			if contentVal, exists := draftMap["content"]; exists {
				if contentStr, ok := contentVal.(string); ok && contentStr != "" {
					imageURL := ExtractFirstImageFromHTML(contentStr)
					if imageURL != "" {
						log.Printf("[SUCCESS] getValue: image_url extracted from content with goquery -> '%s'", imageURL)
						return imageURL, true
					}
				}
			}

			log.Printf("[DEBUG] getValue: No image found for image_url")
			return nil, false
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

		// Alias content -> article
		lookupKey := normalized
		if normalized == "content" {
			lookupKey = "article"
		}

		// Lookup di draftMap
		val, exists := draftMap[lookupKey]
		if exists {
			if normalized == "content" {
				if strVal, ok := val.(string); ok {
					stripped := stripHTMLTags(strVal)
					log.Printf("[SUCCESS] getValue: DraftMap lookup '%s' -> content stripped (%d chars)",
						normalized, len(stripped))
					return stripped, true
				}
			}
			log.Printf("[SUCCESS] getValue: DraftMap lookup '%s' -> '%v'", normalized, val)
			return val, true
		}

		log.Printf("[DEBUG] getValue: Not found for '%s'", normalized)
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
			log.Printf("[DEBUG] getValue: Map returned empty")
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
			log.Printf("[DEBUG] getValue: Array returned empty")
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
