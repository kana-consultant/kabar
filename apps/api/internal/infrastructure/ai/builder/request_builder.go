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
	log.Println("[INFO] MaxTokens:", config.MaxTokens)
	log.Println("[INFO] Temperature:", config.Temperature)
	log.Println("[INFO] Prompt Length:", len(prompt))
	log.Println("[INFO] System Prompt Length:", len(config.SystemPrompt))

	if len(prompt) > 300 {
		log.Println("[INFO] Prompt Preview:")
		log.Println(prompt[:300] + "...")
	} else {
		log.Println("[INFO] Prompt Preview:")
		log.Println(prompt)
	}

	// Validasi & normalisasi MaxTokens SEBELUM dipakai
	maxTokens := config.MaxTokens
	if maxTokens < 1 {
		log.Printf("[WARNING] MaxTokens %d is invalid, using default 1024", maxTokens)
		maxTokens = 1024
	}

	// =========================
	// TEMPLATE VARIABLES - AMBIL DARI CONFIG
	// =========================
	vars := map[string]string{
		"{model}":         config.ModelName,
		"{prompt}":        escapeJSON(prompt),
		"{system_prompt}": escapeJSON(config.SystemPrompt),
		"{temperature}":   fmt.Sprintf("%.2f", config.Temperature),
		"{max_token}":     fmt.Sprintf("%d", maxTokens),
		"{max_tokens}":    fmt.Sprintf("%d", maxTokens), // alias, biar typo lama di template tetap kebaca
	}

	log.Println("[INFO] Template Variables:")
	for k, v := range vars {
		preview := v
		if len(preview) > 200 {
			preview = preview[:200] + "..."
		}
		log.Printf("  %s => %s\n", k, preview)
	}

	// =========================
	// VALIDASI PLACEHOLDER
	// =========================
	template := config.Template

	if strings.Contains(template, "{system_prompt}") && config.SystemPrompt == "" {
		log.Println("[WARNING] Template requires {system_prompt} but SystemPrompt is empty")
	}

	if config.Temperature < 0 || config.Temperature > 2 {
		log.Printf("[WARNING] Temperature %.2f is outside recommended range (0-2)", config.Temperature)
	}

	// =========================
	// TEMPLATE REPLACEMENT
	// =========================
	replacedTemplate := helper.ReplaceTemplate(template, vars)

	log.Println("[INFO] Template After Replacement:")
	if len(replacedTemplate) > 1000 {
		log.Println(replacedTemplate[:1000] + "...")
	} else {
		log.Println(replacedTemplate)
	}

	// =========================
	// JSON VALIDATION
	// =========================
	var bodyMap map[string]interface{}
	if err := json.Unmarshal([]byte(replacedTemplate), &bodyMap); err != nil {
		log.Println("[ERROR] Invalid JSON Template")
		log.Println("[ERROR] JSON Unmarshal Error:", err)
		return nil, fmt.Errorf("invalid template: %w", err)
	}

	log.Println("[INFO] JSON template parsed successfully")

	// =========================
	// FORCE FIELD NUMERIK & TAMBAHAN THINKING CONFIG
	// Ini langkah krusial: walaupun template menulis angka
	// sebagai string (misal "3000" atau "1.00"), kita overwrite
	// di level Go pakai tipe numerik asli (int/float64) supaya
	// hasil json.Marshal nanti berupa JSON number, bukan string.
	// =========================
	if genConfig, ok := bodyMap["generationConfig"].(map[string]interface{}); ok {
		// --- Gemini style: nested generationConfig ---
		genConfig["maxOutputTokens"] = maxTokens      // int asli
		genConfig["temperature"] = config.Temperature // float64 asli
		genConfig["thinkingConfig"] = map[string]interface{}{
			"thinkingBudget": 0, // matikan thinking biar token output nggak kepotong
		}
		log.Println("[INFO] Format terdeteksi: Gemini (generationConfig nested)")
	} else if _, hasMaxTokens := bodyMap["max_tokens"]; hasMaxTokens {
		// --- OpenAI-compatible style: max_tokens & temperature di root ---
		// Dipakai oleh gpt-oss, Groq, OpenRouter, dan provider sejenis.
		bodyMap["max_tokens"] = maxTokens           // int asli, bukan string "1024"
		bodyMap["temperature"] = config.Temperature // float64 asli, bukan string "1.00"
		log.Println("[INFO] Format terdeteksi: OpenAI-compatible (root-level max_tokens)")
	} else {
		log.Println("[WARNING] Tidak ditemukan generationConfig maupun max_tokens di root, skip normalisasi numerik")
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
	log.Println("========== END BUILD REQUEST BODY ==========")

	return finalBody, nil
}

func (b *RequestBuilder) BuildImageRequestBody(
	config *generate.ModelConfig,
	prompt string,
) ([]byte, error) {

	log.Println("========== BUILD IMAGE REQUEST BODY ==========")
	log.Println("[INFO] Model Name:", config.ModelName)
	log.Println("[INFO] MaxTokens:", config.MaxTokens)
	log.Println("[INFO] Temperature:", config.Temperature)
	log.Println("[INFO] System Prompt Length:", len(config.SystemPrompt))

	maxTokens := config.MaxTokens
	if maxTokens < 1 {
		log.Printf("[WARNING] MaxTokens %d is invalid, using default 1024", maxTokens)
		maxTokens = 1024
	}

	vars := map[string]string{
		"{model}":         config.ModelName,
		"{prompt}":        escapeJSON(prompt),
		"{system_prompt}": escapeJSON(config.SystemPrompt),
		"{temperature}":   fmt.Sprintf("%.2f", config.Temperature),
		"{max_tokens}":    fmt.Sprintf("%d", maxTokens),
		"{max_token}":     fmt.Sprintf("%d", maxTokens), // alias safety
	}

	log.Println("[INFO] Template Variables:")
	for k, v := range vars {
		preview := v
		if len(preview) > 200 {
			preview = preview[:200] + "..."
		}
		log.Printf("  %s => %s\n", k, preview)
	}

	template := config.Template

	if config.Temperature < 0 || config.Temperature > 2 {
		log.Printf("[WARNING] Temperature %.2f is outside recommended range (0-2)", config.Temperature)
	}

	replacedTemplate := helper.ReplaceTemplate(template, vars)

	log.Println("[INFO] Template After Replacement:")
	if len(replacedTemplate) > 1000 {
		log.Println(replacedTemplate[:1000] + "...")
	} else {
		log.Println(replacedTemplate)
	}

	var bodyMap map[string]interface{}
	if err := json.Unmarshal([]byte(replacedTemplate), &bodyMap); err != nil {
		log.Println("[ERROR] Invalid image template:", err)
		return nil, fmt.Errorf("invalid image template: %w", err)
	}

	log.Println("[INFO] JSON template parsed successfully")

	// Force numerik: cover format Gemini (nested) maupun OpenAI-compatible (root-level)
	if genConfig, ok := bodyMap["generationConfig"].(map[string]interface{}); ok {
		genConfig["maxOutputTokens"] = maxTokens
		genConfig["temperature"] = config.Temperature
		log.Println("[INFO] Format terdeteksi: Gemini (generationConfig nested)")
	} else if _, hasMaxTokens := bodyMap["max_tokens"]; hasMaxTokens {
		bodyMap["max_tokens"] = maxTokens
		bodyMap["temperature"] = config.Temperature
		log.Println("[INFO] Format terdeteksi: OpenAI-compatible (root-level max_tokens)")
	} else {
		log.Println("[WARNING] Tidak ditemukan generationConfig maupun max_tokens di root, skip normalisasi numerik")
	}

	finalBody, err := json.MarshalIndent(bodyMap, "", "  ")
	if err != nil {
		log.Println("[ERROR] Failed to marshal image request body:", err)
		return nil, err
	}

	log.Println("[INFO] Final Image Request Body:")
	if len(finalBody) > 2000 {
		log.Println(string(finalBody[:2000]) + "...")
	} else {
		log.Println(string(finalBody))
	}

	log.Println("[SUCCESS] Image request body generated successfully")
	log.Println("========== END BUILD IMAGE REQUEST BODY ==========")

	return finalBody, nil
}

// escapeJSON escapes special characters for JSON string
func escapeJSON(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}
