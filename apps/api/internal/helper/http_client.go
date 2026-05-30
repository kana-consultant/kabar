package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func (s *PostService) sendWithRetry(cfg *ProductConfig, body map[string]interface{}) ([]byte, error) {
	var lastErr error

	if cfg.FullURL == "" {
		return nil, fmt.Errorf("full URL is empty")
	}

	client := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
	}

	for i := 0; i < cfg.RetryCount; i++ {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			lastErr = fmt.Errorf("failed to marshal request body (attempt %d/%d): %w", i+1, cfg.RetryCount, err)
			if i < cfg.RetryCount-1 {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, lastErr
		}

		log.Printf("[REQUEST BODY] %s", string(jsonBody))

		req, err := http.NewRequest(cfg.HTTPMethod, cfg.FullURL, bytes.NewReader(jsonBody))
		if err != nil {
			lastErr = fmt.Errorf("failed to create request (attempt %d/%d): %w", i+1, cfg.RetryCount, err)
			if i < cfg.RetryCount-1 {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, lastErr
		}

		// Set headers
		s.setRequestHeaders(req, cfg)

		log.Printf("[AUTH HEADER] %s", req.Header.Get("Authorization"))
		log.Printf("[REQUEST URL] %s %s", cfg.HTTPMethod, cfg.FullURL)

		resp, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed (attempt %d/%d): %w", i+1, cfg.RetryCount, err)
			if i < cfg.RetryCount-1 {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, lastErr
		}

		bodyBytes, err := func() ([]byte, error) {
			defer resp.Body.Close()
			return io.ReadAll(resp.Body)
		}()

		log.Printf(
			"[RESPONSE] status=%d body=%s",
			resp.StatusCode,
			string(bodyBytes),
		)

		if err != nil {
			lastErr = fmt.Errorf("failed to read response body (attempt %d/%d): %w", i+1, cfg.RetryCount, err)
			if i < cfg.RetryCount-1 {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, lastErr
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return bodyBytes, nil
		}

		lastErr = fmt.Errorf("HTTP %d %s (attempt %d/%d): %s", resp.StatusCode, cfg.FullURL, i+1, cfg.RetryCount, string(bodyBytes))
		log.Printf("[ERROR] %v", lastErr)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("all retries exhausted (%d attempts): %w", cfg.RetryCount, lastErr)
}

func (s *PostService) setRequestHeaders(req *http.Request, cfg *ProductConfig) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if cfg.CustomHeaders != nil {
		for k, v := range cfg.CustomHeaders {
			v = resolveTemplate(v, map[string]string{
				"api_key": cfg.APIKey,
			})

			kLower := strings.ToLower(strings.TrimSpace(k))
			vTrim := strings.TrimSpace(v)

			if kLower == "" {
				continue
			}

			// Handle API key header
			if kLower == "x-api-key" && (vTrim == "" || vTrim == "{{api_key}}") {
				v = cfg.APIKey
			}

			// Handle Bearer token
			if kLower == "authorization" && (vTrim == "" || vTrim == "Bearer" || vTrim == "Bearer {{api_key}}") {
				v = "Bearer " + cfg.APIKey
			}

			req.Header.Set(k, v)
		}
	}

	// Fallback auth
	if req.Header.Get("Authorization") == "" && cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	}
}
