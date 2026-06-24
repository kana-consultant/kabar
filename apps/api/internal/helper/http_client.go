package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"seo-backend/internal/domain/product"
	"strings"
	"time"

	"github.com/google/uuid"
)

func setRequestHeaders(req *http.Request, cfg *product.ProductConfig) {
	// Set default headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Persiapkan data untuk template
	templateData := map[string]string{
		"api_key":    cfg.APIKey,
		"id":         uuid.New().String(),
		"product_id": cfg.ProductID,
	}

	// Tambahkan data dari ExecutionResults untuk template
	if cfg.ExecutionResults != nil {
		for key, value := range cfg.ExecutionResults {
			if strVal, ok := value.(string); ok {
				templateData[key] = strVal
			}
		}
	}

	// Proses custom headers
	if cfg.CustomHeaders != nil {
		for k, v := range cfg.CustomHeaders {
			kLower := strings.ToLower(strings.TrimSpace(k))

			if kLower == "" {
				continue
			}

			// Skip default headers
			if kLower == "content-type" || kLower == "accept" {
				continue
			}

			// Resolve template
			v = resolveTemplate(v, templateData)

			if strings.TrimSpace(v) == "" {
				continue
			}

			// Set header
			req.Header.Set(k, v)
		}
	}
}

func SendWithRetry(cfg *product.ProductConfig, body map[string]interface{}) ([]byte, error) {
	var lastErr error

	// Prepare execution context
	if cfg.ExecutionResults == nil {
		cfg.ExecutionResults = make(map[string]interface{})
	}

	// Validasi URL
	if cfg.FullURL == "" {
		return nil, fmt.Errorf("full URL is empty for product %s", cfg.ProductID)
	}

	// Set default values
	if cfg.RetryCount <= 0 {
		cfg.RetryCount = 3
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30
	}

	// Trim URL
	cfg.FullURL = strings.TrimSpace(cfg.FullURL)

	// HTTP Client
	client := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
	}

	// Enrich body with context
	enrichedBody := enrichBodyWithContext(body, cfg)

	// Retry loop
	for i := 0; i < cfg.RetryCount; i++ {
		attempt := i + 1

		// Marshal body
		jsonBody, err := json.Marshal(enrichedBody)
		if err != nil {
			lastErr = fmt.Errorf("marshal error (attempt %d/%d): %w", attempt, cfg.RetryCount, err)
			if i < cfg.RetryCount-1 {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, lastErr
		}

		// Create request
		req, err := http.NewRequest(cfg.HTTPMethod, cfg.FullURL, bytes.NewReader(jsonBody))
		if err != nil {
			lastErr = fmt.Errorf("request creation error (attempt %d/%d): %w", attempt, cfg.RetryCount, err)
			if i < cfg.RetryCount-1 {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, lastErr
		}

		// Set headers
		setRequestHeaders(req, cfg)

		// Send request
		startTime := time.Now()
		resp, err := client.Do(req)
		elapsed := time.Since(startTime)

		if err != nil {
			lastErr = fmt.Errorf("request failed (attempt %d/%d): %w", attempt, cfg.RetryCount, err)
			log.Printf("[ERROR] Product %s, Node %s, Attempt %d/%d: %v (elapsed: %dms)",
				cfg.ProductID, cfg.CurrentNodeID, attempt, cfg.RetryCount, err, elapsed.Milliseconds())
			if i < cfg.RetryCount-1 {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, lastErr
		}
		defer resp.Body.Close()

		// Read response
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("read response error (attempt %d/%d): %w", attempt, cfg.RetryCount, err)
			log.Printf("[ERROR] Product %s, Node %s, Attempt %d/%d: failed to read response - %v",
				cfg.ProductID, cfg.CurrentNodeID, attempt, cfg.RetryCount, err)
			if i < cfg.RetryCount-1 {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, lastErr
		}

		// Store in execution results
		cfg.ExecutionResults["last_response"] = string(bodyBytes)
		cfg.ExecutionResults["last_status_code"] = resp.StatusCode
		cfg.ExecutionResults["last_attempt"] = attempt
		cfg.ExecutionResults["last_elapsed_ms"] = elapsed.Milliseconds()

		// Check success
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// Parse and store response data
			var responseData map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &responseData); err == nil {
				for key, value := range responseData {
					cfg.ExecutionResults["response_"+key] = value
				}
			}

			log.Printf("[SUCCESS] Product %s, Node %s, Status: %d (attempt %d, elapsed: %dms)",
				cfg.ProductID, cfg.CurrentNodeID, resp.StatusCode, attempt, elapsed.Milliseconds())
			return bodyBytes, nil
		}

		// Handle non-2xx response
		lastErr = fmt.Errorf("HTTP %d (attempt %d/%d): %s", resp.StatusCode, attempt, cfg.RetryCount, string(bodyBytes))
		log.Printf("[ERROR] Product %s, Node %s, Status: %d (attempt %d/%d)",
			cfg.ProductID, cfg.CurrentNodeID, resp.StatusCode, attempt, cfg.RetryCount)

		if i < cfg.RetryCount-1 {
			time.Sleep(2 * time.Second)
		}
	}

	log.Printf("[ERROR] Product %s, Node %s: all %d retries failed - %v",
		cfg.ProductID, cfg.CurrentNodeID, cfg.RetryCount, lastErr)

	return nil, fmt.Errorf("all retries exhausted (%d attempts): %w", cfg.RetryCount, lastErr)
}

// getEndpointPathFromConfig mengambil endpoint path dari workflow node berdasarkan current node ID
func getEndpointPathFromConfig(cfg *product.ProductConfig) string {
	// Jika tidak ada current node ID, return empty
	if cfg.CurrentNodeID == "" {
		return ""
	}

	// Cari node dengan ID yang sesuai
	for _, node := range cfg.WorkflowNodes {
		if node.ID == cfg.CurrentNodeID {
			if node.AdapterConfig != nil && node.AdapterConfig.EndpointPath != "" {
				return node.AdapterConfig.EndpointPath
			}
			break
		}
	}

	// Fallback 1: coba ambil dari adapter endpoint
	if cfg.AdapterEndpoint != "" {
		return cfg.AdapterEndpoint
	}

	// Fallback 2: extract dari FullURL
	if cfg.FullURL != "" {
		return extractPathFromURL(cfg.FullURL)
	}

	return ""
}

// extractPathFromURL mengambil path dari URL
func extractPathFromURL(fullURL string) string {
	if fullURL == "" {
		return ""
	}

	// Parse URL
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		// Jika gagal parse, coba ambil setelah domain
		if strings.Contains(fullURL, "://") {
			parts := strings.SplitN(fullURL, "://", 2)
			if len(parts) == 2 {
				hostParts := strings.SplitN(parts[1], "/", 2)
				if len(hostParts) == 2 {
					return "/" + hostParts[1]
				}
			}
		}
		return fullURL
	}

	// Return path + query string
	path := parsedURL.Path
	if parsedURL.RawQuery != "" {
		path = path + "?" + parsedURL.RawQuery
	}

	return path
}

// extractEndpointPath mengambil path dari URL
func extractEndpointPath(fullURL string) string {
	if fullURL == "" {
		return ""
	}

	// Parse URL
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		// Jika gagal parse, coba ambil setelah domain
		if strings.Contains(fullURL, "://") {
			parts := strings.SplitN(fullURL, "://", 2)
			if len(parts) == 2 {
				hostParts := strings.SplitN(parts[1], "/", 2)
				if len(hostParts) == 2 {
					return "/" + hostParts[1]
				}
			}
		}
		return fullURL
	}

	// Return path + query string
	path := parsedURL.Path
	if parsedURL.RawQuery != "" {
		path = path + "?" + parsedURL.RawQuery
	}
	return path
}

// enrichBodyWithContext menambahkan data dari execution results ke body
func enrichBodyWithContext(body map[string]interface{}, cfg *product.ProductConfig) map[string]interface{} {
	enriched := make(map[string]interface{})

	// Copy original body
	for k, v := range body {
		enriched[k] = v
	}

	// Tambahkan context dari execution results
	if cfg.ExecutionResults != nil {
		for key, value := range cfg.ExecutionResults {
			if _, exists := enriched[key]; !exists {
				enriched[key] = value
			}
		}
	}

	// Tambahkan metadata produk
	enriched["product_id"] = cfg.ProductID
	enriched["current_node"] = cfg.CurrentNodeID

	if cfg.MetaConfig != nil {
		for key, value := range cfg.MetaConfig {
			if _, exists := enriched[key]; !exists {
				enriched[key] = value
			}
		}
	}

	return enriched
}

// getWorkflowLevel mendapatkan level workflow dari node ID
func getWorkflowLevel(cfg *product.ProductConfig, nodeID string) int {
	if cfg.WorkflowLevelMap == nil {
		return 0
	}

	for level, nodes := range cfg.WorkflowLevelMap {
		for _, node := range nodes {
			if node.ID == nodeID {
				return level
			}
		}
	}
	return 0
}

func safeJSON(b []byte) []byte {
	if json.Valid(b) {
		return b
	}
	escaped, _ := json.Marshal(string(b))
	return escaped
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
