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
		attempt := i + 1

		jsonBody, err := json.Marshal(body)
		if err != nil {
			lastErr = fmt.Errorf("failed to marshal request body (attempt %d/%d): %w", attempt, cfg.RetryCount, err)
			logJSON("marshal_error", map[string]interface{}{
				"attempt":     attempt,
				"retry_count": cfg.RetryCount,
				"error":       err.Error(),
			})
			if i < cfg.RetryCount-1 {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, lastErr
		}

		logJSON("request", map[string]interface{}{
			"attempt":     attempt,
			"retry_count": cfg.RetryCount,
			"method":      cfg.HTTPMethod,
			"url":         cfg.FullURL,
			"body":        json.RawMessage(jsonBody),
		})

		req, err := http.NewRequest(cfg.HTTPMethod, cfg.FullURL, bytes.NewReader(jsonBody))
		if err != nil {
			lastErr = fmt.Errorf("failed to create request (attempt %d/%d): %w", attempt, cfg.RetryCount, err)
			logJSON("create_request_error", map[string]interface{}{
				"attempt": attempt,
				"error":   err.Error(),
			})
			if i < cfg.RetryCount-1 {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, lastErr
		}

		s.setRequestHeaders(req, cfg)

		logJSON("request_headers", map[string]interface{}{
			"attempt":       attempt,
			"authorization": req.Header.Get("Authorization"),
			"content_type":  req.Header.Get("Content-Type"),
		})

		resp, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed (attempt %d/%d): %w", attempt, cfg.RetryCount, err)
			logJSON("http_error", map[string]interface{}{
				"attempt": attempt,
				"error":   err.Error(),
			})
			if i < cfg.RetryCount-1 {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, lastErr
		}

		bodyBytes, readErr := func() ([]byte, error) {
			defer resp.Body.Close()
			return io.ReadAll(resp.Body)
		}()

		// ✅ Cek readErr dulu sebelum log, agar body valid saat di-log
		if readErr != nil {
			lastErr = fmt.Errorf("failed to read response body (attempt %d/%d): %w", attempt, cfg.RetryCount, readErr)
			logJSON("read_body_error", map[string]interface{}{
				"attempt":     attempt,
				"status_code": resp.StatusCode,
				"error":       readErr.Error(),
			})
			if i < cfg.RetryCount-1 {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, lastErr
		}

		// ✅ Log response hanya setelah body berhasil dibaca
		logJSON("response", map[string]interface{}{
			"attempt":     attempt,
			"status_code": resp.StatusCode,
			"body":        json.RawMessage(safeJSON(bodyBytes)),
		})

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return bodyBytes, nil
		}

		lastErr = fmt.Errorf("HTTP %d %s (attempt %d/%d): %s", resp.StatusCode, cfg.FullURL, attempt, cfg.RetryCount, string(bodyBytes))
		logJSON("http_status_error", map[string]interface{}{
			"attempt":     attempt,
			"status_code": resp.StatusCode,
			"url":         cfg.FullURL,
			"error":       string(bodyBytes),
		})
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

func logJSON(event string, fields map[string]interface{}) {
	fields["event"] = event
	fields["time"] = time.Now().Format(time.RFC3339)
	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		log.Printf("[LOG_ERROR] failed to marshal log fields: %v", err)
		return
	}
	log.Println(string(b))
}

// safeJSON memastikan bodyBytes valid sebagai JSON; jika bukan, wrap sebagai string
func safeJSON(b []byte) []byte {
	if json.Valid(b) {
		return b
	}
	escaped, _ := json.Marshal(string(b))
	return escaped
}
