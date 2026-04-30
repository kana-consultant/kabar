// internal/services/api_caller.go
package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"seo-backend/internal/models"
)

type APICaller struct{}

func NewAPICaller() *APICaller {
	return &APICaller{}
}

type APIRequest struct {
	Provider     models.APIProvider
	APIKey       string
	Model        string
	SystemPrompt string
	UserPrompt   string
	IsImage      bool
}

func (c *APICaller) Call(req APIRequest) (string, error) {
	// Build URL
	endpoint := req.Provider.TextEndpoint
	if req.IsImage && req.Provider.ImageEndpoint.Valid {
		endpoint = req.Provider.ImageEndpoint.String
	}

	url := req.Provider.BaseURL + strings.ReplaceAll(endpoint, "{model}", req.Model)

	// Build request body from template
	var requestBody map[string]interface{}

	templateStr := req.Provider.RequestTemplate.String
	templateStr = strings.ReplaceAll(templateStr, "{model}", req.Model)
	templateStr = strings.ReplaceAll(templateStr, "{system_prompt}", escapeJSON(req.SystemPrompt))
	templateStr = strings.ReplaceAll(templateStr, "{prompt}", escapeJSON(req.UserPrompt))

	json.Unmarshal([]byte(templateStr), &requestBody)

	// Add headers
	httpReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody(requestBody)))

	// Parse default headers
	if req.Provider.DefaultHeaders.Valid {
		var headers map[string]string
		json.Unmarshal([]byte(req.Provider.DefaultHeaders.String), &headers)
		for k, v := range headers {
			httpReq.Header.Set(k, v)
		}
	}

	// Add auth header based on auth_type
	switch req.Provider.AuthType {
	case "bearer":
		httpReq.Header.Set(req.Provider.AuthHeader, req.Provider.AuthPrefix+" "+req.APIKey)
	case "api_key":
		httpReq.Header.Set(req.Provider.AuthHeader, req.APIKey)
	case "x-api-key":
		httpReq.Header.Set(req.Provider.AuthHeader, req.APIKey)
	default:
		httpReq.Header.Set("Authorization", "Bearer "+req.APIKey)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Parse response using JSON path
	return c.extractResponse(body, req.Provider, req.IsImage)
}

func (c *APICaller) extractResponse(body []byte, provider models.APIProvider, isImage bool) (string, error) {
	var data interface{}
	json.Unmarshal(body, &data)

	path := provider.ResponseTextPath
	if isImage && provider.ResponseImagePath.Valid {
		path = provider.ResponseImagePath.String
	}

	result := getJSONPath(data, path)
	if result == nil {
		return "", fmt.Errorf("failed to extract response from path: %s", path)
	}

	switch v := result.(type) {
	case string:
		if isImage && strings.HasPrefix(v, "data:image") {
			return v, nil
		}
		if isImage {
			return "data:image/png;base64," + v, nil
		}
		return v, nil
	default:
		b, _ := json.Marshal(v)
		return string(b), nil
	}
}

// Helper functions
func jsonBody(data map[string]interface{}) []byte {
	b, _ := json.Marshal(data)
	return b
}

func escapeJSON(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func getJSONPath(data interface{}, path string) interface{} {
	// Simple JSON path implementation
	// Supports: $.field.subfield[0].field
	parts := strings.Split(strings.TrimPrefix(path, "$."), ".")
	current := data

	for _, part := range parts {
		if strings.Contains(part, "[") {
			// Handle array index
			arrPart := strings.Split(part, "[")
			key := arrPart[0]
			idx := 0
			fmt.Sscanf(arrPart[1], "%d]", &idx)

			if m, ok := current.(map[string]interface{}); ok {
				if arr, ok := m[key].([]interface{}); ok && idx < len(arr) {
					current = arr[idx]
				} else {
					return nil
				}
			}
		} else {
			if m, ok := current.(map[string]interface{}); ok {
				current = m[part]
			} else {
				return nil
			}
		}
	}
	return current
}
