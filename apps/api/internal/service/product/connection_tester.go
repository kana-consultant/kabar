package product

import (
	"io"
	"log"
	"net/http"
	"time"
)

type ConnectionTester struct{}

func NewConnectionTester() *ConnectionTester {
	return &ConnectionTester{}
}

func (t *ConnectionTester) Test(product *ProductBasicInfo, config *AdapterConfig) *ConnectionTestResult {
	fullURL := product.APIEndpoint + config.EndpointPath
	log.Printf("Testing connection to: %s [%s]", fullURL, config.HTTPMethod)

	var lastResp *http.Response
	var lastErr error

	for i := 0; i < config.RetryCount; i++ {
		req, err := http.NewRequest(config.HTTPMethod, fullURL, nil)
		if err != nil {
			lastErr = err
			continue
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")
		if product.APIKey != "" {
			req.Header.Set("Authorization", "Bearer "+product.APIKey)
		}
		for k, v := range config.CustomHeaders {
			req.Header.Set(k, v)
		}

		client := &http.Client{
			Timeout: time.Duration(config.TimeoutSeconds) * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			if i < config.RetryCount-1 {
				time.Sleep(1 * time.Second)
				continue
			}
			break
		}

		lastResp = resp
		lastErr = nil
		break
	}

	if lastErr != nil {
		return &ConnectionTestResult{
			Success:  false,
			Message:  "Failed to connect: " + lastErr.Error(),
			TestedAt: time.Now().Format(time.RFC3339),
		}
	}
	defer lastResp.Body.Close()

	bodyBytes, _ := io.ReadAll(io.LimitReader(lastResp.Body, 1024))
	isConnected := lastResp.StatusCode >= 200 && lastResp.StatusCode < 300

	return &ConnectionTestResult{
		Success:     isConnected,
		StatusCode:  lastResp.StatusCode,
		StatusText:  lastResp.Status,
		Message:     getStatusMessage(lastResp.StatusCode),
		ProductName: product.Name,
		Endpoint:    fullURL,
		Method:      config.HTTPMethod,
		Response:    truncateString(string(bodyBytes), 500),
		TestedAt:    time.Now().Format(time.RFC3339),
	}
}

func getStatusMessage(statusCode int) string {
	switch statusCode {
	case 200:
		return "OK - Connection successful"
	case 201:
		return "Created - Resource created successfully"
	case 400:
		return "Bad Request - Check your request parameters"
	case 401:
		return "Unauthorized - Invalid API key or credentials"
	case 403:
		return "Forbidden - Access denied"
	case 404:
		return "Not Found - Endpoint not found"
	case 500:
		return "Internal Server Error - Server side issue"
	default:
		return "Connection test completed"
	}
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
