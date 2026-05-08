// internal/infrastructure/ai/builder/request_builder.go
package builder

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"seo-backend/internal/domain/generate"
	"seo-backend/internal/helper"
)

type RequestBuilder struct{}

func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{}
}

func (b *RequestBuilder) BuildArticleRequestBody(
	config *generate.ModelConfig,
	prompt string,
) ([]byte, error) {

	log.Println("========== BUILD REQUEST BODY ==========")

	log.Println("[INFO] Model Name:", config.ModelName)
	log.Println("[INFO] Prompt Length:", len(prompt))

	if len(prompt) > 300 {
		log.Println("[INFO] Prompt Preview:")
		log.Println(prompt[:300] + "...")
	} else {
		log.Println("[INFO] Prompt Preview:")
		log.Println(prompt)
	}

	// =========================
	// TEMPLATE VARIABLES
	// =========================
	vars := map[string]string{
		"{model}":       config.ModelName,
		"{prompt}":      helper.EscapeJSON(prompt),
		"{temperature}": "0.7",
		"{max_tokens}":  "4000",
	}

	log.Println("[INFO] Template Variables:")
	for k, v := range vars {

		preview := v

		// prevent huge logs
		if len(preview) > 200 {
			preview = preview[:200] + "..."
		}

		log.Printf("  %s => %s\n", k, preview)
	}

	// =========================
	// TEMPLATE REPLACEMENT
	// =========================
	template := helper.ReplaceTemplate(config.Template, vars)

	log.Println("[INFO] Template After Replacement:")

	if len(template) > 1000 {
		log.Println(template[:1000] + "...")
	} else {
		log.Println(template)
	}

	// =========================
	// JSON VALIDATION
	// =========================
	var bodyMap map[string]interface{}

	if err := json.Unmarshal([]byte(template), &bodyMap); err != nil {

		log.Println("[ERROR] Invalid JSON Template")
		log.Println("[ERROR] JSON Unmarshal Error:", err)

		return nil, fmt.Errorf("invalid template: %w", err)
	}

	log.Println("[INFO] JSON template parsed successfully")

	// =========================
	// DEFAULT MODEL
	// =========================
	if _, ok := bodyMap["model"]; !ok {

		log.Println("[WARNING] 'model' field missing")
		log.Println("[INFO] Applying default model:", config.ModelName)

		bodyMap["model"] = config.ModelName
	}

	// =========================
	// DEFAULT MESSAGES
	// =========================
	if _, ok := bodyMap["messages"]; !ok {

		log.Println("[WARNING] 'messages' field missing")
		log.Println("[INFO] Applying default messages")

		bodyMap["messages"] = []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		}
	}

	// =========================
	// FINAL BODY
	// =========================
	finalBody, err := json.MarshalIndent(bodyMap, "", "  ")
	if err != nil {

		log.Println("[ERROR] Failed to marshal final request body:", err)

		return nil, err
	}

	log.Println("[INFO] Final Request Body:")

	if len(finalBody) > 2000 {
		log.Println(string(finalBody[:2000]) + "...")
	} else {
		log.Println(string(finalBody))
	}

	log.Println("[SUCCESS] Request body generated successfully")

	return finalBody, nil
}

func (b *RequestBuilder) BuildImageRequestBody(config *generate.ModelConfig, prompt string) ([]byte, error) {
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
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}
