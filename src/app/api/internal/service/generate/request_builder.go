package generate

import (
	"encoding/json"
	"fmt"
	"strings"

	"seo-backend/internal/helper"
)

func buildArticleRequestBody(config *ModelConfig, prompt string) ([]byte, error) {
	vars := map[string]string{
		"{model}":       config.ModelName,
		"{prompt}":      helper.EscapeJSON(prompt),
		"{temperature}": "0.7",
		"{max_tokens}":  "4000",
	}

	template := helper.ReplaceTemplate(config.Template, vars)

	var bodyMap map[string]interface{}
	if err := json.Unmarshal([]byte(template), &bodyMap); err != nil {
		return nil, fmt.Errorf("invalid template: %w", err)
	}

	// Set defaults if missing
	if _, ok := bodyMap["model"]; !ok {
		bodyMap["model"] = config.ModelName
	}
	if _, ok := bodyMap["messages"]; !ok {
		bodyMap["messages"] = []map[string]string{
			{"role": "user", "content": prompt},
		}
	}

	return json.Marshal(bodyMap)
}

func buildImageRequestBody(config *ModelConfig, prompt string) ([]byte, error) {
	template := config.Template
	template = strings.ReplaceAll(template, "{prompt}", escapeJSON(prompt))
	template = strings.ReplaceAll(template, "{model}", config.ModelName)

	var bodyMap map[string]interface{}
	if err := json.Unmarshal([]byte(template), &bodyMap); err != nil {
		return nil, fmt.Errorf("invalid image template: %w", err)
	}

	return json.Marshal(bodyMap)
}

func escapeJSON(s string) string {
	// Escape special characters for JSON
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}
