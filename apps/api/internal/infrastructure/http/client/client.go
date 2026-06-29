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
	maxBackoff := 30 * time.Second

	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			log.Printf("[INFO] Retry attempt %d/%d after %v", attempt+1, maxRetries, backoff)

			select {
			case <-time.After(backoff):
				// Lanjut retry
			case <-ctx.Done():
				log.Printf("[ERROR] Context cancelled during backoff: %v", ctx.Err())
				return nil, fmt.Errorf("request cancelled: %w", ctx.Err())
			}

			// Exponential backoff
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		// =========================
		// CREATE REQUEST WITH TIMEOUT
		// =========================
		reqCtx, cancel := context.WithTimeout(context.Background(), timeout)
		// JANGAN PAKAI defer cancel() di sini!

		req, err := http.NewRequestWithContext(reqCtx, "POST", url, bytes.NewBuffer(body))
		if err != nil {
			cancel() // Cancel context karena gagal buat request
			log.Printf("[ERROR] Failed to create request: %v", err)
			lastErr = fmt.Errorf("failed to create request: %w", err)
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

		log.Printf("[INFO] Sending request (attempt %d/%d)...", attempt+1, maxRetries)
		startTime := time.Now()

		resp, err := c.client.Do(req)
		elapsed := time.Since(startTime)

		if err != nil {
			cancel() // Cancel context jika request gagal
			log.Printf("[ERROR] Request attempt %d failed after %v: %v", attempt+1, elapsed, err)
			lastErr = fmt.Errorf("request failed: %w", err)

			// Check if it's timeout error
			if strings.Contains(err.Error(), "context deadline exceeded") ||
				strings.Contains(err.Error(), "timeout") ||
				strings.Contains(err.Error(), "deadline") {
				log.Printf("[ERROR] Request timeout after %v", timeout)
				continue
			}

			// Untuk error lain, tetap coba retry
			continue
		}

		log.Printf("[INFO] Request completed in %v", elapsed)
		log.Println("[INFO] Response Status:", resp.Status)

		// PENTING: Baca response body SEBELUM cancel context
		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		// BARU cancel context setelah selesai baca body
		cancel()

		if err != nil {
			log.Printf("[ERROR] Failed to read response body: %v", err)
			lastErr = fmt.Errorf("failed to read response: %w", err)
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
			lastErr = fmt.Errorf("server error with status %d: %s", resp.StatusCode, respBodyStr)
			continue
		}

		// Success
		if resp.StatusCode == http.StatusOK {
			log.Println("[SUCCESS] Request completed successfully")
			log.Println("========== END HTTP REQUEST ==========")
			return respBody, nil
		}

		// Non-retryable error status (400, 401, 403, 404, etc.)
		log.Printf("[ERROR] API returned non-retryable status (%d)", resp.StatusCode)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, respBodyStr)
	}

	// Out of retries
	log.Printf("[ERROR] All %d retries exhausted", maxRetries)
	if lastErr != nil {
		return nil, fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
	}
	return nil, fmt.Errorf("failed after %d retries with unknown error", maxRetries)
}
