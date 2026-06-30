package builder

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"seo-backend/internal/domain/generate"
	"seo-backend/internal/helper"
)

// RequestType defines the type of request being built
type RequestType string

const (
	ArticleRequest RequestType = "article"
	ImageRequest   RequestType = "image"
)

// RequestBuilder handles building request bodies for various AI model APIs
type RequestBuilder struct {
	logger *log.Logger
}

// NewRequestBuilder creates a new RequestBuilder instance
func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{
		logger: log.Default(),
	}
}

// SetLogger allows injecting a custom logger
func (b *RequestBuilder) SetLogger(logger *log.Logger) {
	if logger != nil {
		b.logger = logger
	}
}

// BuildRequestBody builds request body for both article and image generation
func (b *RequestBuilder) BuildRequestBody(
	config *generate.ModelConfig,
	prompt string,
	requestType RequestType,
) ([]byte, error) {

	b.logRequestInfo(config, prompt, requestType)

	// Validate input parameters
	if err := b.validateInput(config, prompt, requestType); err != nil {
		return nil, err
	}

	// Clean the template for image requests
	template := config.Template
	if requestType == ImageRequest {
		template = b.cleanImageTemplate(template)
	}

	// Build template variables with proper escaping
	vars := b.buildTemplateVariables(config, prompt, requestType)

	// Replace template placeholders
	replacedTemplate, err := b.replaceTemplate(template, vars)
	if err != nil {
		return nil, err
	}

	// Parse and validate JSON structure
	bodyMap, err := b.parseAndValidateJSON(replacedTemplate, requestType)
	if err != nil {
		return nil, err
	}

	// Perform request-type specific validation
	if requestType == ImageRequest {
		if err := b.validateImagePrompt(bodyMap, prompt); err != nil {
			return nil, err
		}
	}

	// Marshal the final request body
	return b.marshalFinalBody(bodyMap, requestType)
}

// BuildArticleRequestBody builds request body for article generation (backward compatibility)
func (b *RequestBuilder) BuildArticleRequestBody(
	config *generate.ModelConfig,
	prompt string,
) ([]byte, error) {
	return b.BuildRequestBody(config, prompt, ArticleRequest)
}

// BuildImageRequestBody builds request body for image generation (backward compatibility)
func (b *RequestBuilder) BuildImageRequestBody(
	config *generate.ModelConfig,
	prompt string,
) ([]byte, error) {
	return b.BuildRequestBody(config, prompt, ImageRequest)
}

// cleanImageTemplate ensures the template has the {prompt} placeholder
func (b *RequestBuilder) cleanImageTemplate(template string) string {
	// Check if template already has {prompt} placeholder
	if strings.Contains(template, "{prompt}") {
		return template
	}

	// Try to parse and fix the template
	var templateMap map[string]interface{}
	if err := json.Unmarshal([]byte(template), &templateMap); err != nil {
		b.logger.Printf("[WARNING] Could not parse template for cleaning: %v", err)
		return template
	}

	// Replace hardcoded prompt with placeholder
	if _, ok := templateMap["prompt"]; ok {
		templateMap["prompt"] = "{prompt}"
		b.logger.Println("[INFO] Replaced hardcoded prompt with {prompt} placeholder")

		// Marshal back to JSON
		if cleaned, err := json.MarshalIndent(templateMap, "", "  "); err == nil {
			return string(cleaned)
		}
	}

	return template
}

// logRequestInfo logs detailed request information
func (b *RequestBuilder) logRequestInfo(config *generate.ModelConfig, prompt string, requestType RequestType) {
	title := "BUILD REQUEST BODY"
	if requestType == ImageRequest {
		title = "BUILD IMAGE REQUEST BODY"
	}

	b.logger.Printf("========== %s ==========", title)
	b.logger.Printf("[INFO] Request Type: %s", requestType)
	b.logger.Printf("[INFO] Model Name: %s", config.ModelName)
	b.logger.Printf("[INFO] MaxTokens: %d", config.MaxTokens)
	b.logger.Printf("[INFO] Temperature: %.2f", config.Temperature)
	b.logger.Printf("[INFO] Prompt Length: %d", len(prompt))
	b.logger.Printf("[INFO] System Prompt Length: %d", len(config.SystemPrompt))

	if len(prompt) > 300 {
		b.logger.Printf("[INFO] Prompt Preview:\n%s...", prompt[:300])
	} else {
		b.logger.Printf("[INFO] Prompt Preview:\n%s", prompt)
	}
}

// validateInput validates all input parameters
func (b *RequestBuilder) validateInput(config *generate.ModelConfig, prompt string, requestType RequestType) error {
	// Validate prompt for image requests
	if requestType == ImageRequest && prompt == "" {
		err := fmt.Errorf("prompt cannot be empty for image generation")
		b.logger.Printf("[ERROR] %v", err)
		return err
	}

	// Validate temperature range
	if config.Temperature < 0 || config.Temperature > 2 {
		b.logger.Printf("[WARNING] Temperature %.2f is outside recommended range (0-2)", config.Temperature)
	}

	// Validate max tokens
	if config.MaxTokens < 1 {
		b.logger.Printf("[WARNING] MaxTokens %d is invalid, will use default", config.MaxTokens)
	}

	// Check if template requires system_prompt but it's empty
	if strings.Contains(config.Template, "{system_prompt}") && config.SystemPrompt == "" {
		b.logger.Println("[WARNING] Template requires {system_prompt} but SystemPrompt is empty")
	}

	// For image requests, check if template contains prompt after cleaning
	if requestType == ImageRequest {
		cleanedTemplate := b.cleanImageTemplate(config.Template)
		if !strings.Contains(cleanedTemplate, "{prompt}") {
			err := fmt.Errorf("image template must contain {prompt} placeholder")
			b.logger.Printf("[ERROR] %v", err)
			templatePreview := cleanedTemplate
			if len(templatePreview) > 200 {
				templatePreview = templatePreview[:200] + "..."
			}
			b.logger.Printf("[ERROR] Template content: %s", templatePreview)
			return err
		}
	}

	return nil
}

// buildTemplateVariables creates the template variables map with proper escaping
func (b *RequestBuilder) buildTemplateVariables(config *generate.ModelConfig, prompt string, requestType RequestType) map[string]string {
	maxTokens := config.MaxTokens
	if maxTokens < 1 {
		maxTokens = 1024
	}

	// For image requests, use only the prompt without system prompt prepended
	imagePrompt := prompt
	if requestType == ImageRequest {
		// If system prompt is being prepended elsewhere, we might want to clean it
		// The actual prompt should be the clean user prompt
		imagePrompt = prompt
	}

	vars := map[string]string{
		"{model}":         config.ModelName,
		"{prompt}":        escapeJSON(imagePrompt),
		"{system_prompt}": escapeJSON(config.SystemPrompt),
		"{temperature}":   fmt.Sprintf("%.2f", config.Temperature),
		"{max_tokens}":    fmt.Sprintf("%d", maxTokens),
	}

	b.logger.Println("[INFO] Template Variables:")
	for k, v := range vars {
		preview := v
		if len(preview) > 200 {
			preview = preview[:200] + "..."
		}
		b.logger.Printf("  %s => %s\n", k, preview)
	}

	return vars
}

// replaceTemplate replaces placeholders in the template
func (b *RequestBuilder) replaceTemplate(template string, vars map[string]string) (string, error) {
	replacedTemplate := helper.ReplaceTemplate(template, vars)

	b.logger.Println("[INFO] Template After Replacement:")
	if len(replacedTemplate) > 1000 {
		b.logger.Printf("%s...", replacedTemplate[:1000])
	} else {
		b.logger.Println(replacedTemplate)
	}

	return replacedTemplate, nil
}

// parseAndValidateJSON parses and validates the JSON structure
func (b *RequestBuilder) parseAndValidateJSON(replacedTemplate string, requestType RequestType) (map[string]interface{}, error) {
	var bodyMap map[string]interface{}
	if err := json.Unmarshal([]byte(replacedTemplate), &bodyMap); err != nil {
		errorMsg := "Invalid JSON Template"
		if requestType == ImageRequest {
			errorMsg = "Invalid image template"
		}
		b.logger.Printf("[ERROR] %s: %v", errorMsg, err)
		return nil, fmt.Errorf("%s: %w", strings.ToLower(errorMsg), err)
	}

	b.logger.Println("[INFO] JSON template parsed successfully")
	return bodyMap, nil
}

// validateImagePrompt validates the prompt in image request body
func (b *RequestBuilder) validateImagePrompt(bodyMap map[string]interface{}, originalPrompt string) error {
	promptInBody, ok := bodyMap["prompt"].(string)
	if !ok {
		b.logger.Println("[WARNING] No 'prompt' field found in final body map")
		return nil
	}

	// Unescape the prompt for proper comparison
	unescapedPrompt, err := unescapeJSON(promptInBody)
	if err != nil {
		b.logger.Printf("[WARNING] Failed to unescape prompt: %v", err)
		unescapedPrompt = promptInBody
	}

	// For image prompts, only compare the core prompt (not the system prompt part)
	// If the prompt contains system prompt prepended, we might want to trim it
	if !strings.Contains(unescapedPrompt, originalPrompt) && !strings.Contains(originalPrompt, unescapedPrompt) {
		b.logger.Println("[WARNING] Prompt in final body may differ from input prompt!")
		b.logger.Printf("[WARNING] Input prompt (first 100 chars): %s...",
			originalPrompt[:min(100, len(originalPrompt))])
		b.logger.Printf("[WARNING] Body prompt (first 100 chars): %s...",
			promptInBody[:min(100, len(promptInBody))])
		// Don't return error for image prompts as they might be enriched
	} else {
		b.logger.Println("[INFO] Prompt in final body matches input prompt ✓")
	}

	return nil
}

// marshalFinalBody marshals the final JSON body
func (b *RequestBuilder) marshalFinalBody(bodyMap map[string]interface{}, requestType RequestType) ([]byte, error) {
	finalBody, err := json.MarshalIndent(bodyMap, "", "  ")
	if err != nil {
		b.logger.Printf("[ERROR] Failed to marshal final request body: %v", err)
		return nil, err
	}

	b.logger.Println("[INFO] Final Request Body:")
	if len(finalBody) > 2000 {
		b.logger.Printf("%s...", string(finalBody[:2000]))
	} else {
		b.logger.Println(string(finalBody))
	}

	title := "Request body"
	if requestType == ImageRequest {
		title = "Image request body"
	}
	b.logger.Printf("[SUCCESS] %s generated successfully", title)
	b.logger.Println("========== END BUILD REQUEST BODY ==========")

	return finalBody, nil
}

// Helper functions

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// escapeJSON properly escapes special characters for JSON string using json.Marshal
func escapeJSON(s string) string {
	if s == "" {
		return ""
	}

	// Use json.Marshal for proper escaping
	b, err := json.Marshal(s)
	if err != nil {
		return s
	}

	// Remove the surrounding quotes from marshaled JSON string
	if len(b) >= 2 {
		return string(b[1 : len(b)-1])
	}
	return s
}

// unescapeJSON unescapes a JSON-escaped string
func unescapeJSON(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	var result string
	// Add quotes around the string to make it a valid JSON string
	err := json.Unmarshal([]byte(`"`+s+`"`), &result)
	if err != nil {
		return s, err
	}
	return result, nil
}
