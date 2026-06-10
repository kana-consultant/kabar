// internal/infrastructure/ai/parser/response_parser.go
package parser

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"seo-backend/internal/domain/generate"
	"seo-backend/internal/helper"
)

type ResponseParser struct{}

func NewResponseParser() *ResponseParser {
	return &ResponseParser{}
}

func (p *ResponseParser) ParseArticleResponse(response []byte, responsePath string) (*generate.ArticleResult, error) {
	var result interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to parse provider response: %w", err)
	}

	text := helper.ExtractByPath(result, responsePath)

	// Clean the text
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	var article generate.ArticleResult

	// Try to parse as JSON
	if err := json.Unmarshal([]byte(text), &article); err != nil {
		log.Printf("[WARNING] JSON parse failed: %v", err)
		log.Printf("[DEBUG] Attempting to extract fields with regex...")

		// Extract fields manually with regex
		article.Title = extractJSONField(text, "title")
		article.Slug = extractJSONField(text, "slug")
		article.Excerpt = extractJSONField(text, "excerpt")
		article.Content = extractJSONField(text, "content")
		article.ImagePrompt = extractJSONField(text, "imagePrompt")
		article.WordCount = extractJSONNumber(text, "wordCount")

		// Extract keywords array
		article.Keywords = extractJSONArray(text, "keywords")

		// Log extracted fields
		log.Printf("[DEBUG] Extracted - Title: %s, WordCount: %d, HasContent: %v",
			article.Title, article.WordCount, article.Content != "")

		// If content still empty, use raw text
		if article.Content == "" && article.Title == "" {
			article.Content = text
		}
	}

	// Calculate word count if still 0
	if article.WordCount == 0 && article.Content != "" {
		article.WordCount = len(strings.Fields(article.Content))
		log.Printf("[INFO] Calculated word count: %d", article.WordCount)
	}

	return &article, nil
}

// extractJSONField extracts string value from JSON field
func extractJSONField(jsonStr, field string) string {
	// Multiple patterns to handle different JSON formats
	patterns := []string{
		// "field": "value"
		fmt.Sprintf(`"%s"\s*:\s*"((?:\\.|[^"\\])*)"`, field),
		// "field": "value with escaped quotes
		fmt.Sprintf(`"%s"\s*:\s*"((?:[^"\\]|\\.)*)"`, field),
		// "field": "value (unclosed - truncated)
		fmt.Sprintf(`"%s"\s*:\s*"([^"]*(?:\\.[^"]*)*)`, field),
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(jsonStr)
		if len(matches) > 1 {
			// Unescape JSON string
			value := matches[1]
			value = strings.ReplaceAll(value, `\"`, `"`)
			value = strings.ReplaceAll(value, `\\`, `\`)
			value = strings.ReplaceAll(value, `\/`, `/`)
			value = strings.ReplaceAll(value, `\n`, "\n")
			value = strings.ReplaceAll(value, `\t`, "\t")
			return value
		}
	}
	return ""
}

// extractJSONNumber extracts number value from JSON field
func extractJSONNumber(jsonStr, field string) int {
	patterns := []string{
		fmt.Sprintf(`"%s"\s*:\s*([0-9]+)`, field),
		fmt.Sprintf(`"%s"\s*:\s*([0-9]+)`, field),
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(jsonStr)
		if len(matches) > 1 {
			val, err := strconv.Atoi(matches[1])
			if err == nil {
				return val
			}
		}
	}
	return 0
}

// extractJSONArray extracts array of strings from JSON field
func extractJSONArray(jsonStr, field string) []string {
	patterns := []string{
		fmt.Sprintf(`"%s"\s*:\s*\[(.*?)\]`, field),
		fmt.Sprintf(`"%s"\s*:\s*\[(.*?)(?:\]|$)`, field), // truncated
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(jsonStr)
		if len(matches) > 1 {
			arrayContent := matches[1]
			if arrayContent == "" {
				return []string{}
			}

			// Split by comma and clean
			items := strings.Split(arrayContent, ",")
			result := make([]string, 0, len(items))
			for _, item := range items {
				item = strings.TrimSpace(item)
				item = strings.Trim(item, `"`)
				item = strings.TrimSpace(item)
				if item != "" {
					result = append(result, item)
				}
			}
			return result
		}
	}
	return []string{}
}

func (p *ResponseParser) ParseImageResponse(response []byte, responsePath string) (string, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		return "", fmt.Errorf("failed to parse response JSON: %w", err)
	}

	imageURL := helper.ExtractByPath(result, responsePath)
	if imageURL == "" {
		return "", fmt.Errorf("no image URL found in response")
	}

	return imageURL, nil
}

func (p *ResponseParser) ParseImageResponseBase64(response []byte, responsePath string) (string, string, error) {
	// Parsing JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		log.Printf("[ERROR] Failed to parse image response JSON: %v", err)
		log.Printf("[DEBUG] Raw response (first 500 chars): %s", string(response[:min(500, len(response))]))
		return "", "", err
	}

	// Navigasi ke field yang mengandung Base64
	base64Str := helper.ExtractByPath(result, responsePath)

	if base64Str == "" {
		return "", "", fmt.Errorf("no base64 data found at path: %s", responsePath)
	}

	// Clean base64 string
	base64Str = cleanBase64String(base64Str)

	// Tentukan content-type dari Base64 header
	contentType := detectContentTypeFromBase64(base64Str)

	return base64Str, contentType, nil
}

func cleanBase64String(base64Str string) string {
	// Remove data URL prefix if present
	if strings.Contains(base64Str, "base64,") {
		parts := strings.Split(base64Str, "base64,")
		base64Str = parts[len(parts)-1]
	}

	// Remove common prefixes
	base64Str = strings.TrimPrefix(base64Str, "data:image/png;base64,")
	base64Str = strings.TrimPrefix(base64Str, "data:image/jpeg;base64,")
	base64Str = strings.TrimPrefix(base64Str, "data:image/jpg;base64,")
	base64Str = strings.TrimPrefix(base64Str, "data:image/webp;base64,")
	base64Str = strings.TrimPrefix(base64Str, "data:image/gif;base64,")

	// Remove whitespace
	base64Str = strings.TrimSpace(base64Str)

	return base64Str
}

func detectContentTypeFromBase64(base64Str string) string {
	// Decode sedikit data untuk deteksi
	decodedLen := min(100, len(base64Str))
	decoded, err := base64.StdEncoding.DecodeString(base64Str[:decodedLen])
	if err != nil {
		log.Printf("[WARNING] Failed to decode base64 for content type detection: %v", err)
		return "image/png"
	}

	// Deteksi magic bytes
	if len(decoded) >= 8 {
		// PNG: 89 50 4E 47
		if decoded[0] == 0x89 && decoded[1] == 0x50 && decoded[2] == 0x4E && decoded[3] == 0x47 {
			return "image/png"
		}
		// JPEG: FF D8 FF
		if decoded[0] == 0xFF && decoded[1] == 0xD8 && decoded[2] == 0xFF {
			return "image/jpeg"
		}
		// WEBP: 52 49 46 46
		if decoded[0] == 0x52 && decoded[1] == 0x49 && decoded[2] == 0x46 && decoded[3] == 0x46 {
			return "image/webp"
		}
		// GIF: 47 49 46
		if decoded[0] == 0x47 && decoded[1] == 0x49 && decoded[2] == 0x46 {
			return "image/gif"
		}
		// BMP: 42 4D
		if len(decoded) >= 2 && decoded[0] == 0x42 && decoded[1] == 0x4D {
			return "image/bmp"
		}
	}

	return "image/png" // default
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
