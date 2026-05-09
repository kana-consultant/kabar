// internal/infrastructure/http/client/client.go
package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"seo-backend/internal/domain/generate"
)

type HTTPClient struct {
	client *http.Client
}

func NewHTTPClient(timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{Timeout: timeout},
	}
}

func (c *HTTPClient) SendRequest(
	ctx context.Context,
	config *generate.ModelConfig,
	body []byte,
	timeout time.Duration,
) ([]byte, error) {

	url := config.BaseURL + config.Endpoint

	log.Println("========== HTTP REQUEST ==========")
	log.Println("[INFO] URL:", url)
	log.Println("[INFO] Method: POST")
	log.Println("[INFO] Timeout:", timeout)

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		url,
		bytes.NewBuffer(body),
	)

	if err != nil {
		log.Println("[ERROR] Failed to create request:", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// =========================
	// AUTH HEADER
	// =========================
	header := config.AuthHeader
	if header == "" {
		header = "Authorization"
	}

	authValue := config.APIKey

	if config.AuthPrefix != "" {
		authValue = config.AuthPrefix + " " + config.APIKey
	}

	req.Header.Set(header, authValue)

	// Mask API KEY
	maskedKey := config.APIKey
	if len(maskedKey) > 10 {
		maskedKey = maskedKey[:10] + "..."
	}

	log.Println("[INFO] Auth Type:", config.AuthType)
	log.Println("[INFO] Auth Header:", header)
	log.Println("[INFO] Auth Prefix:", config.AuthPrefix)
	log.Println("[INFO] API Key Exists:", config.APIKey != "")
	log.Println("[INFO] API Key Preview:", maskedKey)

	req.Header.Set("Content-Type", "application/json")

	log.Println("[INFO] Headers:")
	for k, v := range req.Header {
		log.Printf("  %s: %v\n", k, v)
	}

	log.Println("[INFO] Request Body:")
	log.Println(string(body))

	client := &http.Client{
		Timeout: timeout,
	}

	log.Println("[INFO] Sending request...")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("[ERROR] Failed to send request:", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	log.Println("[INFO] Response Status:", resp.Status)
	log.Println("[INFO] Response Status Code:", resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ERROR] Failed to read response body:", err)
		return nil, err
	}

	log.Println("[INFO] Response Body:")
	log.Println(string(respBody))

	if resp.StatusCode != http.StatusOK {
		log.Printf(
			"[ERROR] API returned non-200 status (%d)\n",
			resp.StatusCode,
		)

		return nil, fmt.Errorf(
			"API returned status %d: %s",
			resp.StatusCode,
			string(respBody),
		)
	}

	log.Println("[SUCCESS] Request completed successfully")

	return respBody, nil
}
