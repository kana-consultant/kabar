package generate

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type HTTPClient struct {
	client *http.Client
}

func NewHTTPClient(timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *HTTPClient) SendRequest(config *ModelConfig, body []byte, timeout time.Duration) ([]byte, error) {
	// Build URL
	fullURL := strings.TrimRight(config.BaseURL, "/") + "/" + strings.TrimLeft(config.Endpoint, "/")
	fullURL = strings.ReplaceAll(fullURL, "{model}", config.ModelName)

	// Create request
	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(req, config)
	c.setProviderSpecificHeaders(req, config)

	// Send request
	client := &http.Client{Timeout: timeout}
	log.Printf("Sending request to: %s", fullURL)

	startTime := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		return nil, fmt.Errorf("request failed after %v: %w", duration, err)
	}
	defer resp.Body.Close()

	log.Printf("Response received in %v, status: %d", duration, resp.StatusCode)

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Log response sample for debugging
	c.logResponseSample(respBody)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (c *HTTPClient) setAuthHeader(req *http.Request, config *ModelConfig) {
	switch config.AuthType {
	case "bearer":
		authValue := strings.TrimSpace(config.AuthPrefix + " " + config.APIKey)
		req.Header.Set(config.AuthHeader, authValue)
		log.Printf("Auth: Bearer")
	case "api_key":
		req.Header.Set(config.AuthHeader, config.APIKey)
		log.Printf("Auth: API Key")
	default:
		req.Header.Set("Authorization", "Bearer "+config.APIKey)
		log.Printf("Auth: Bearer (default)")
	}
}

func (c *HTTPClient) setProviderSpecificHeaders(req *http.Request, config *ModelConfig) {
	// OpenRouter specific headers
	if strings.Contains(config.BaseURL, "openrouter.ai") {
		req.Header.Set("HTTP-Referer", "http://localhost:8080")
		req.Header.Set("X-Title", "SEO Generator")
	}
}

func (c *HTTPClient) logResponseSample(respBody []byte) {
	if len(respBody) > 200 {
		log.Printf("Response sample: %s", string(respBody[:200]))
	} else {
		log.Printf("Response body: %s", string(respBody))
	}
}
