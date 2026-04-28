package handlers

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"seo-backend/internal/database"
)

type GenerateHandler struct{}

func NewGenerateHandler() *GenerateHandler {
	return &GenerateHandler{}
}

type GenerateArticleRequest struct {
	Topic             string `json:"topic"`
	ModelID           string `json:"modelId"`
	Tone              string `json:"tone"`
	Length            string `json:"length"`
	Language          string `json:"language"`
	AutoGenerateImage bool   `json:"autoGenerateImage"`
}

type GenerateArticleResponse struct {
	Title            string   `json:"title"`
	Content          string   `json:"content"`
	Excerpt          string   `json:"excerpt"`
	Keywords         []string `json:"keywords"`
	ImagePrompt      string   `json:"imagePrompt"`
	ImageURL         string   `json:"imageUrl"`
	WordCount        int      `json:"wordCount"`
	ReadabilityScore int      `json:"readabilityScore"`
	SeoScore         int      `json:"seoScore"`
}

type GenerateImageRequest struct {
	Prompt  string `json:"prompt"`
	ModelID string `json:"modelId"`
}

type GenerateImageResponse struct {
	ImageURL    string `json:"imageUrl"`
	Prompt      string `json:"prompt"`
	GeneratedAt string `json:"generatedAt"`
	Model       string `json:"model"`
}

// Tambahkan fungsi decrypt di GenerateHandler
func (h *GenerateHandler) decrypt(cipherText string) (string, error) {
	key := os.Getenv("ENCRYPTION_KEY")
	if key == "" {
		key = "01234567890123456789012345678901"
	}

	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("invalid ciphertext")
	}

	nonce, encrypted := data[:nonceSize], data[nonceSize:]
	plain, err := gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return "", err
	}

	return string(plain), nil
}

func (h *GenerateHandler) GenerateArticle(w http.ResponseWriter, r *http.Request) {
	log.Println("========== GENERATE ARTICLE ==========")

	var req GenerateArticleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR decode request: %v", err)
		http.Error(w, "Invalid request", 400)
		return
	}

	log.Printf(
		"Request: Topic=%s ModelID=%s Tone=%s Length=%s Language=%s",
		req.Topic,
		req.ModelID,
		req.Tone,
		req.Length,
		req.Language,
	)

	if req.Topic == "" || req.ModelID == "" {
		http.Error(w, "topic and modelId required", 400)
		return
	}

	// ----------------------------------
	// helper dynamic template replacer
	// ----------------------------------
	replaceTemplate := func(tpl string, vars map[string]string) string {
		for k, v := range vars {
			tpl = strings.ReplaceAll(
				tpl,
				k,
				v,
			)
		}
		return tpl
	}

	// ----------------------------------
	// DB config
	// ----------------------------------

	var (
		encryptedKey string
		modelName    string
		templateStr  string
		responsePath string
		baseURL      string
		authType     string
		authHeader   string
		authPrefix   string
		endpoint     string
	)

	log.Printf("Querying DB for modelId=%s", req.ModelID)

	err := database.GetDB().QueryRow(`
		SELECT
			ak.key_encrypted,
			m.name,
			m.request_template,
			m.response_text_path,
			p.base_url,
			p.auth_type,
			p.auth_header,
			p.auth_prefix,
			p.text_endpoint
		FROM api_keys ak
		JOIN ai_models m
			ON ak.model_id = m.id
		JOIN api_providers p
			ON m.provider_id = p.id
		WHERE ak.id = $1
		AND ak.is_active = true
		AND ak.service = 'text'
		LIMIT 1
	`, req.ModelID).Scan(
		&encryptedKey,
		&modelName,
		&templateStr,
		&responsePath,
		&baseURL,
		&authType,
		&authHeader,
		&authPrefix,
		&endpoint,
	)

	if err != nil {
		log.Printf("DB ERROR: %v", err)
		http.Error(
			w,
			"model not found or no active api key",
			400,
		)
		return
	}

	// ----------------------------------
	// decrypt api key
	// ----------------------------------

	apiKey, err := h.decrypt(encryptedKey)
	if err != nil {
		log.Printf(
			"decrypt failed, fallback plaintext: %v",
			err,
		)
		apiKey = encryptedKey
	}

	log.Printf(
		"Provider=%s URL=%s Endpoint=%s",
		modelName,
		baseURL,
		endpoint,
	)

	url := baseURL + strings.ReplaceAll(
		endpoint,
		"{model}",
		modelName,
	)

	log.Printf("Full URL: %s", url)

	// ----------------------------------
	// build prompt
	// ----------------------------------

	prompt := fmt.Sprintf(`
Write a high-quality SEO-friendly article in %s about "%s".

Requirements:
- Tone: %s
- Length: %s
- Content must be valid HTML
- SEO optimized

Return ONLY valid JSON.

{
"title":"string",
"content":"<h1>...</h1>",
"excerpt":"string",
"seoScore":number,
"readabilityScore":number,
"wordCount":number
}
`,
		req.Language,
		req.Topic,
		req.Tone,
		req.Length,
	)

	log.Printf(
		"Prompt length: %d",
		len(prompt),
	)

	// ----------------------------------
	// dynamic template variables
	// ----------------------------------

	vars := map[string]string{
		"{model}":       modelName,
		"{prompt}":      escapeJSON(prompt),
		"{temperature}": "0.7",
		"{max_tokens}":  "4000",
	}

	templateStr = replaceTemplate(
		templateStr,
		vars,
	)

	log.Printf(
		"Template after replace: %s",
		templateStr,
	)

	// ----------------------------------
	// validate template JSON
	// ----------------------------------

	var bodyMap map[string]interface{}

	if err := json.Unmarshal(
		[]byte(templateStr),
		&bodyMap,
	); err != nil {

		log.Printf(
			"Template parse error: %v",
			err,
		)

		log.Printf(
			"Broken template=%s",
			templateStr,
		)

		http.Error(
			w,
			"invalid template",
			500,
		)
		return
	}

	// fallback defaults
	if _, ok := bodyMap["model"]; !ok {
		bodyMap["model"] = modelName
	}

	if _, ok := bodyMap["messages"]; !ok {
		bodyMap["messages"] = []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		}
	}

	bodyBytes, err := json.Marshal(bodyMap)
	if err != nil {
		http.Error(
			w,
			"marshal request failed",
			500,
		)
		return
	}

	log.Printf(
		"Request body: %s",
		string(bodyBytes),
	)

	// ----------------------------------
	// outbound request
	// ----------------------------------

	httpReq, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(bodyBytes),
	)

	if err != nil {
		http.Error(
			w,
			"failed create request",
			500,
		)
		return
	}

	httpReq.Header.Set(
		"Content-Type",
		"application/json",
	)

	switch authType {

	case "bearer":
		httpReq.Header.Set(
			authHeader,
			authPrefix+" "+apiKey,
		)

	case "api_key":
		httpReq.Header.Set(
			authHeader,
			apiKey,
		)

	default:
		httpReq.Header.Set(
			"Authorization",
			"Bearer "+apiKey,
		)
	}

	// OpenRouter extras
	if strings.Contains(
		baseURL,
		"openrouter.ai",
	) {
		httpReq.Header.Set(
			"HTTP-Referer",
			"http://localhost:8080",
		)

		httpReq.Header.Set(
			"X-Title",
			"SEO Generator",
		)
	}

	client := &http.Client{
		Timeout: 90 * time.Second,
	}

	log.Println("Sending request...")

	start := time.Now()

	resp, err := client.Do(httpReq)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			500,
		)
		return
	}

	defer resp.Body.Close()

	log.Printf(
		"Response in %v status=%d",
		time.Since(start),
		resp.StatusCode,
	)

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.Printf(
			"API ERROR: %s",
			string(respBody),
		)

		http.Error(
			w,
			string(respBody),
			500,
		)
		return
	}

	// ----------------------------------
	// parse provider response
	// ----------------------------------

	var result interface{}

	if err := json.Unmarshal(
		respBody,
		&result,
	); err != nil {
		http.Error(
			w,
			"failed parse provider response",
			500,
		)
		return
	}

	text := extractByPath(
		result,
		responsePath,
	)

	jsonStr := extractJSON(text)

	var article GenerateArticleResponse

	if err := json.Unmarshal(
		[]byte(jsonStr),
		&article,
	); err != nil {

		log.Printf(
			"Article parse error=%v",
			err,
		)

		log.Printf(
			"Raw JSON=%s",
			jsonStr,
		)

		http.Error(
			w,
			"failed parse article response",
			500,
		)
		return
	}

	log.Printf(
		"SUCCESS title=%s words=%d",
		article.Title,
		article.WordCount,
	)

	log.Println("========== END GENERATE ARTICLE ==========")

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(article)
}

func (h *GenerateHandler) GenerateImage(w http.ResponseWriter, r *http.Request) {
	log.Println("========== GENERATE IMAGE ==========")

	var req GenerateImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: Failed to decode request: %v", err)
		http.Error(w, "Invalid request", 400)
		return
	}

	log.Printf("Request: Prompt=%s, ModelID=%s", req.Prompt, req.ModelID)

	if req.Prompt == "" || req.ModelID == "" {
		log.Printf("ERROR: Missing required fields")
		http.Error(w, "prompt and modelId required", 400)
		return
	}

	// Query database
	var encryptedKey, modelName, template, responsePath, baseURL, authType, authHeader, authPrefix, endpoint string
	var isActive bool

	log.Printf("Querying database for modelId: %s", req.ModelID)

	err := database.GetDB().QueryRow(`
		SELECT ak.key_encrypted, m.name, m.request_template, m.response_image_path,
		       p.base_url, p.auth_type, p.auth_header, p.auth_prefix, p.text_endpoint,
		       m.is_active
		FROM api_keys ak
		JOIN ai_models m ON ak.model_id = m.id
		JOIN api_providers p ON m.provider_id = p.id
		WHERE ak.id = $1 AND ak.is_active = true
	`, req.ModelID).Scan(
		&encryptedKey, &modelName, &template, &responsePath,
		&baseURL, &authType, &authHeader, &authPrefix, &endpoint, &isActive,
	)

	if err != nil {
		log.Printf("ERROR: Database query failed: %v", err)
		http.Error(w, "model not found or no active api key", 400)
		return
	}

	log.Printf("Found config: model=%s, baseURL=%s, endpoint=%s", modelName, baseURL, endpoint)

	// Decrypt API key
	apiKey, err := h.decrypt(encryptedKey)
	if err != nil {
		log.Printf("Decrypt failed, using as plain text: %v", err)
		apiKey = encryptedKey
	}

	// Build URL dengan benar
	fullURL := strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(endpoint, "/")
	// Replace {model} di URL jika ada
	fullURL = strings.ReplaceAll(fullURL, "{model}", modelName)
	log.Printf("Full URL: %s", fullURL)

	// Prepare request body (jangan replace {model} di template dulu)
	requestBody := template
	requestBody = strings.ReplaceAll(requestBody, "{prompt}", escapeJSON(req.Prompt))
	requestBody = strings.ReplaceAll(requestBody, "{model}", modelName)

	log.Printf("Request body: %s", requestBody)

	var bodyMap map[string]interface{}
	if err := json.Unmarshal([]byte(requestBody), &bodyMap); err != nil {
		log.Printf("ERROR: Template parse error: %v", err)
		http.Error(w, "invalid template", 500)
		return
	}

	bodyBytes, err := json.Marshal(bodyMap)
	if err != nil {
		log.Printf("ERROR: Failed to marshal body: %v", err)
		http.Error(w, "failed to marshal body", 500)
		return
	}

	httpReq, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("ERROR: Failed to create request: %v", err)
		http.Error(w, "failed to create request", 500)
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// Set auth header
	switch authType {
	case "bearer":
		authValue := strings.TrimSpace(authPrefix + " " + apiKey)
		httpReq.Header.Set(authHeader, authValue)
		log.Printf("Auth: Bearer %s...", apiKey[:min(10, len(apiKey))])
	case "api_key":
		httpReq.Header.Set(authHeader, apiKey)
		log.Printf("Auth: API Key")
	default:
		httpReq.Header.Set("Authorization", "Bearer "+apiKey)
		log.Printf("Auth: Bearer (default)")
	}

	// Send request
	client := &http.Client{Timeout: 120 * time.Second}
	log.Printf("Sending request to provider...")

	startTime := time.Now()
	resp, err := client.Do(httpReq)
	duration := time.Since(startTime)

	if err != nil {
		log.Printf("ERROR: Request failed after %v: %v", duration, err)
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	log.Printf("Response received in %v, status: %d", duration, resp.StatusCode)

	// Read response body ONCE
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read response body: %v", err)
		http.Error(w, "failed to read response", 500)
		return
	}

	// Log sample response
	if len(respBody) > 200 {
		log.Printf("Response sample (first 200 chars): %s", string(respBody[:200]))
	} else {
		log.Printf("Response body: %s", string(respBody))
	}

	if resp.StatusCode != 200 {
		log.Printf("ERROR: API error - status=%d, body=%s", resp.StatusCode, string(respBody))
		http.Error(w, fmt.Sprintf("API error: %s", string(respBody)), 500)
		return
	}

	// Parse JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		log.Printf("ERROR: Failed to parse response JSON: %v", err)
		log.Printf("Raw response: %s", string(respBody))
		http.Error(w, "failed to parse response", 500)
		return
	}

	// Extract image URL
	imageURL := extractByPath(result, responsePath)
	log.Printf("Extracted image URL by path '%s': %s", responsePath, imageURL)

	if imageURL == "" {
		log.Printf("ERROR: No image URL found in response")
		http.Error(w, "no image generated", 500)
		return
	}

	response := GenerateImageResponse{
		ImageURL:    imageURL,
		Prompt:      req.Prompt,
		GeneratedAt: time.Now().Format(time.RFC3339),
		Model:       modelName,
	}

	log.Printf("SUCCESS: Image generated")
	log.Println("========== END GENERATE IMAGE ==========")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ==================== HELPER FUNCTIONS ====================

func escapeJSON(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\r", `\r`)
	s = strings.ReplaceAll(s, "\t", `\t`)
	return s
}

func extractJSON(text string) string {
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start != -1 && end != -1 && end > start {
		return text[start : end+1]
	}
	return "{}"
}

func extractByPath(data interface{}, path string) string {
	parts := strings.Split(path, ".")
	var current interface{} = data

	for _, part := range parts {

		if current == nil {
			return ""
		}

		// support candidates[0]
		if strings.Contains(part, "[") {
			idxStart := strings.Index(part, "[")
			idxEnd := strings.Index(part, "]")

			arrName := part[:idxStart]

			var idx int
			fmt.Sscanf(part[idxStart+1:idxEnd], "%d", &idx)

			m, ok := current.(map[string]interface{})
			if !ok {
				return ""
			}

			arr, ok := m[arrName].([]interface{})
			if !ok || idx >= len(arr) {
				return ""
			}

			current = arr[idx]
			continue
		}

		// support plain numeric segment "0"
		if idx, err := strconv.Atoi(part); err == nil {
			arr, ok := current.([]interface{})
			if !ok || idx >= len(arr) {
				return ""
			}

			current = arr[idx]
			continue
		}

		// normal map access
		m, ok := current.(map[string]interface{})
		if !ok {
			return ""
		}

		current = m[part]
	}

	if s, ok := current.(string); ok {
		return s
	}

	return ""
}
