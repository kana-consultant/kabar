package generate

import (
	"encoding/json"
	"fmt"
	"log"

	"seo-backend/internal/helper"
)

func parseArticleResponse(response []byte, responsePath string) (*ArticleResult, error) {
	var result interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to parse provider response: %w", err)
	}

	text := helper.ExtractByPath(result, responsePath)
	jsonStr := helper.ExtractJSON(text)

	var article ArticleResult
	if err := json.Unmarshal([]byte(jsonStr), &article); err != nil {
		log.Printf("Raw JSON: %s", jsonStr)
		return nil, fmt.Errorf("failed to parse article response: %w", err)
	}

	return &article, nil
}

func parseImageResponse(response []byte, responsePath string) (string, error) {
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
