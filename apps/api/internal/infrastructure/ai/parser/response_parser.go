// internal/infrastructure/ai/parser/response_parser.go
package parser

import (
	"encoding/json"
	"fmt"
	"log"

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
	jsonStr := helper.ExtractJSON(text)

	var article generate.ArticleResult
	if err := json.Unmarshal([]byte(jsonStr), &article); err != nil {
		log.Printf("Raw JSON: %s", jsonStr)
		return nil, fmt.Errorf("failed to parse article response: %w", err)
	}

	return &article, nil
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
