package generate

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"
	"time"

	"seo-backend/internal/domain/generate"
)

// Service handles content generation business logic
type Service struct {
	repo       generate.Repository
	httpClient *http.Client
}

// NewService creates a new generate service
func NewService(repo generate.Repository) *Service {
	return &Service{
		repo: repo,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// GenerateContent generates content based on the request
func (s *Service) GenerateContent(ctx context.Context, req generate.GenerateRequest) (*generate.GenerateResponse, error) {
	// Validation
	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	// Get model configuration
	config, err := s.repo.GetModelConfig(ctx, req.ModelID, req.ServiceType)
	if err != nil {
		return nil, fmt.Errorf("failed to get model config: %w", err)
	}

	// Prepare request body from template
	requestBody, err := s.prepareRequestBody(config.Template, req.Prompt, req.Parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare request body: %w", err)
	}

	// Build full URL
	fullURL := s.buildURL(config.BaseURL, config.Endpoint)

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", fullURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authentication headers
	s.setAuthHeaders(httpReq, config)

	// Execute request
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	response, err := s.parseResponse(body, config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return response, nil
}

// validateRequest validates the generation request
func (s *Service) validateRequest(req generate.GenerateRequest) error {
	if req.ModelID == "" {
		return errors.New("model ID is required")
	}
	if req.ServiceType == "" {
		return errors.New("service type is required")
	}
	if req.Prompt == "" {
		return errors.New("prompt is required")
	}
	return nil
}

// prepareRequestBody prepares the request body using template
func (s *Service) prepareRequestBody(templateStr, prompt string, params map[string]interface{}) ([]byte, error) {
	// Parse template
	tmpl, err := template.New("request").Parse(templateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	// Prepare data for template
	data := map[string]interface{}{
		"prompt":     prompt,
		"parameters": params,
	}

	// Merge additional parameters
	for k, v := range params {
		data[k] = v
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.Bytes(), nil
}

// buildURL builds the full URL
func (s *Service) buildURL(baseURL, endpoint string) string {
	baseURL = strings.TrimSuffix(baseURL, "/")
	endpoint = strings.TrimPrefix(endpoint, "/")
	return fmt.Sprintf("%s/%s", baseURL, endpoint)
}

// setAuthHeaders sets authentication headers
func (s *Service) setAuthHeaders(req *http.Request, config *generate.ModelConfig) {
	req.Header.Set("Content-Type", "application/json")

	switch config.AuthType {
	case "bearer":
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.APIKey))
	case "header":
		headerName := config.AuthHeader
		if headerName == "" {
			headerName = "Authorization"
		}
		prefix := config.AuthPrefix
		if prefix != "" {
			req.Header.Set(headerName, fmt.Sprintf("%s %s", prefix, config.APIKey))
		} else {
			req.Header.Set(headerName, config.APIKey)
		}
	case "api_key":
		// Handle API key in query parameter or header
		headerName := config.AuthHeader
		if headerName != "" {
			req.Header.Set(headerName, config.APIKey)
		} else {
			// Add to query params
			q := req.URL.Query()
			q.Add("api_key", config.APIKey)
			req.URL.RawQuery = q.Encode()
		}
	default:
		// Default to bearer
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.APIKey))
	}
}

// parseResponse parses the API response
func (s *Service) parseResponse(body []byte, config *generate.ModelConfig) (*generate.GenerateResponse, error) {
	var response generate.GenerateResponse
	var jsonData interface{}

	// Parse JSON
	if err := json.Unmarshal(body, &jsonData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Extract text using path
	if config.ResponsePath != "" {
		text, err := s.extractByPath(jsonData, config.ResponsePath)
		if err == nil {
			response.Text = text
		}
	}

	// Extract image using path
	if config.ResponseImagePath != "" {
		image, err := s.extractByPath(jsonData, config.ResponseImagePath)
		if err == nil {
			response.Image = image
		}
	}

	// If no text found, try to use the whole response
	if response.Text == "" {
		response.Text = string(body)
	}

	response.Metadata = map[string]interface{}{
		"model": config.ModelName,
	}

	return &response, nil
}

// extractByPath extracts value from JSON by dot notation path
func (s *Service) extractByPath(data interface{}, path string) (string, error) {
	parts := strings.Split(path, ".")
	current := data

	for _, part := range parts {
		switch v := current.(type) {
		case map[string]interface{}:
			if val, ok := v[part]; ok {
				current = val
			} else {
				return "", fmt.Errorf("path %s not found", path)
			}
		default:
			return "", fmt.Errorf("cannot traverse into non-object at %s", part)
		}
	}

	// Convert to string
	switch v := current.(type) {
	case string:
		return v, nil
	case float64:
		return fmt.Sprintf("%g", v), nil
	case bool:
		return fmt.Sprintf("%t", v), nil
	default:
		// Try to marshal to JSON string
		bytes, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}
}
