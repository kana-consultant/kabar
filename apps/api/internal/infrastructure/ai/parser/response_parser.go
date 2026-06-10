// internal/infrastructure/ai/parser/response_parser.go
package parser

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
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
		// If failed, extract fields manually with regex
		log.Printf("JSON parse failed, extracting with regex: %v", err)

		article.Title = extractJSONField(text, "title")
		article.Slug = extractJSONField(text, "slug")
		article.Excerpt = extractJSONField(text, "excerpt")
		article.Content = extractJSONField(text, "content")

		// If still empty, use raw text as content
		if article.Content == "" && article.Title == "" {
			article.Content = text
		}
	}

	return &article, nil
}

func extractJSONField(jsonStr, field string) string {
	// Cari posisi field
	searchFor := fmt.Sprintf(`"%s":"`, field)
	start := strings.Index(jsonStr, searchFor)
	if start == -1 {
		searchFor = fmt.Sprintf(`"%s": "`, field)
		start = strings.Index(jsonStr, searchFor)
	}
	if start == -1 {
		return ""
	}

	start += len(searchFor)

	// Ambil sampe akhir atau sampe ketemu ", (tapi kalo kepotong sampe akhir)
	end := strings.Index(jsonStr[start:], `"`)
	if end == -1 {
		// Kepotong, ambil sampe akhir
		return jsonStr[start:]
	}

	return jsonStr[start : start+end]
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
		return "", "", err
	}

	// Navigasi ke field yang mengandung Base64 (contoh: data[0].b64_json)
	current := helper.ExtractByPath(result, responsePath)

	base64Str := current

	// Coba tentukan content-type dari Base64 header
	contentType := detectContentTypeFromBase64(base64Str)

	return base64Str, contentType, nil
}

func detectContentTypeFromBase64(base64Str string) string {
	// Decode sedikit data untuk deteksi
	decoded, err := base64.StdEncoding.DecodeString(base64Str[:min(100, len(base64Str))])
	if err != nil {
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
	}

	return "image/png" // default
}
