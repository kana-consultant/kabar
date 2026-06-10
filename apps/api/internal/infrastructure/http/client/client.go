// internal/infrastructure/http/client/client.go
package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"seo-backend/internal/domain/generate"
)

type HTTPClient struct {
	client *http.Client
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{}, // No default timeout, use context instead
	}
}

func (c *HTTPClient) SendRequest(
	ctx context.Context,
	config *generate.ModelConfig,
	body []byte,
	timeout time.Duration,
) ([]byte, error) {

	// =========================
	// BUILD URL WITH MODEL REPLACEMENT
	// =========================
	baseURL := config.BaseURL
	endpoint := config.Endpoint

	baseURL = strings.ReplaceAll(baseURL, "{model}", config.ModelName)
	endpoint = strings.ReplaceAll(endpoint, "{model}", config.ModelName)

	url := baseURL + endpoint

	log.Println("========== HTTP REQUEST ==========")
	log.Println("[INFO] Full URL:", url)
	log.Println("[INFO] Model Name:", config.ModelName)
	log.Println("[INFO] Timeout:", timeout)

	// =========================
	// AUTH HEADER - FIXED FOR sql.NullString
	// =========================
	header := "Authorization"
	if config.AuthHeader.Valid && config.AuthHeader.String != "" {
		header = config.AuthHeader.String
	}

	authValue := config.APIKey
	if config.AuthPrefix.Valid && config.AuthPrefix.String != "" {
		authValue = config.AuthPrefix.String + " " + config.APIKey
	}

	// Mask API KEY
	maskedKey := config.APIKey
	if len(maskedKey) > 10 {
		maskedKey = maskedKey[:10] + "..."
	}

	log.Println("[INFO] Auth Header:", header)
	log.Println("[INFO] API Key Preview:", maskedKey)

	// =========================
	// RETRY LOGIC
	// =========================
	maxRetries := 3
	backoff := 1 * time.Second

	var resp *http.Response
	var respBody []byte
	var err error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			log.Printf("[INFO] Retry attempt %d/%d after %v", attempt+1, maxRetries, backoff)
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
			backoff *= 2
			if backoff > 30*time.Second {
				backoff = 30 * time.Second
			}
		}

		// =========================
		// CREATE NEW CONTEXT WITH TIMEOUT
		// IMPORTANT: Use context.Background() to ignore parent context deadline
		// =========================
		reqCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		req, err := http.NewRequestWithContext(reqCtx, "POST", url, bytes.NewBuffer(body))
		if err != nil {
			log.Println("[ERROR] Failed to create request:", err)
			continue
		}

		// Set headers
		req.Header.Set(header, authValue)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
		req.Header.Set("Accept", "application/json")

		// Log request body (truncate if too long)
		bodyStr := string(body)
		if len(bodyStr) > 2000 {
			log.Println("[INFO] Request Body (truncated):", bodyStr[:2000]+"...")
		} else {
			log.Println("[INFO] Request Body:", bodyStr)
		}

		// Request dump for debug
		dump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			log.Println("[WARNING] Failed to dump request:", err)
		} else {
			log.Println("========== RAW REQUEST DUMP ==========")
			log.Println(string(dump))
			log.Println("========== END REQUEST DUMP ==========")
		}

		log.Println("[INFO] Sending request...")
		startTime := time.Now()

		resp, err = c.client.Do(req)
		elapsed := time.Since(startTime)

		if err != nil {
			log.Printf("[ERROR] Request attempt %d failed after %v: %v", attempt+1, elapsed, err)

			// Check if it's context error, don't retry
			if strings.Contains(err.Error(), "context canceled") || strings.Contains(err.Error(), "context deadline") {
				log.Println("[ERROR] Context error, stopping retries")
				return nil, fmt.Errorf("request timeout after %v: %w", timeout, err)
			}
			continue
		}

		log.Printf("[INFO] Request completed in %v", elapsed)
		log.Println("[INFO] Response Status:", resp.Status)

		// Read response body
		respBody, err = io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			log.Printf("[ERROR] Failed to read response body: %v", err)
			continue
		}

		// Log response body (truncate)
		respBodyStr := string(respBody)
		if len(respBodyStr) > 2000 {
			log.Println("[INFO] Response Body (truncated):", respBodyStr[:2000]+"...")
		} else {
			log.Println("[INFO] Response Body:", respBodyStr)
		}

		// Check if status code is retryable
		if resp.StatusCode == 503 || resp.StatusCode == 429 || resp.StatusCode >= 500 {
			log.Printf("[WARNING] Got status %d, will retry", resp.StatusCode)
			continue
		}

		// Success
		if resp.StatusCode == http.StatusOK {
			log.Println("[SUCCESS] Request completed successfully")
			log.Println("========== END HTTP REQUEST ==========")
			return respBody, nil
		}

		// Non-retryable error status
		log.Printf("[ERROR] API returned non-200 status (%d)", resp.StatusCode)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Out of retries
	if resp == nil {
		return nil, fmt.Errorf("failed after %d retries: %w", maxRetries, err)
	}

	return nil, fmt.Errorf("failed after %d retries: status %d", maxRetries, resp.StatusCode)
}
