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
	// Log sebelum set headers dengan konteks produk
	logJSON("set_headers_start", map[string]interface{}{
		"product_id":           cfg.ProductID,
		"custom_headers_count": len(cfg.CustomHeaders),
		"api_key_present":      cfg.APIKey != "",
		"current_node_id":      cfg.CurrentNodeID,
	})

	// Set default headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Log default headers
	logJSON("default_headers_set", map[string]interface{}{
		"product_id":   cfg.ProductID,
		"content_type": "application/json",
		"accept":       "application/json",
	})

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

	resolvedHeaders := make(map[string]string)

	// Proses custom headers
	if cfg.CustomHeaders != nil {
		for k, v := range cfg.CustomHeaders {
			kLower := strings.ToLower(strings.TrimSpace(k))

			if kLower == "" {
				logJSON("skip_empty_header_key", map[string]interface{}{
					"product_id":   cfg.ProductID,
					"original_key": k,
				})
				continue
			}

			// Skip jika key adalah content-type atau accept (sudah di-set)
			if kLower == "content-type" || kLower == "accept" {
				logJSON("skip_default_header", map[string]interface{}{
					"product_id": cfg.ProductID,
					"key":        k,
					"reason":     "already set by default",
				})
				continue
			}

			originalV := v

			// Resolve template
			v = resolveTemplate(v, templateData)

			if strings.TrimSpace(v) == "" {
				logJSON("skip_empty_header_value", map[string]interface{}{
					"product_id":     cfg.ProductID,
					"key":            k,
					"original_value": originalV,
				})
				continue
			}

			// Set header
			req.Header.Set(k, v)
			resolvedHeaders[k] = v

			logJSON("header_set", map[string]interface{}{
				"product_id":     cfg.ProductID,
				"key":            k,
				"original_value": originalV,
				"resolved_value": v,
				"template_used":  originalV != v,
				"context_keys":   len(templateData),
			})
		}
	} else {
		logJSON("no_custom_headers", map[string]interface{}{
			"product_id": cfg.ProductID,
			"message":    "No custom headers to set",
		})
	}

	// Log final headers
	finalHeaders := make(map[string]string)
	for k, v := range req.Header {
		finalHeaders[k] = strings.Join(v, ", ")
	}

	logJSON("set_headers_complete", map[string]interface{}{
		"product_id":        cfg.ProductID,
		"total_headers_set": len(resolvedHeaders),
		"all_headers":       finalHeaders,
		"custom_headers":    resolvedHeaders,
	})
}

func SendWithRetry(cfg *product.ProductConfig, body map[string]interface{}) ([]byte, error) {
	var lastErr error

	// Prepare execution context
	if cfg.ExecutionResults == nil {
		cfg.ExecutionResults = make(map[string]interface{})
	}

	// Get endpoint path from workflow node
	endpointPath := getEndpointPathFromConfig(cfg)

	// Log start
	logJSON("send_request_start", map[string]interface{}{
		"product_id":      cfg.ProductID,
		"endpoint_path":   endpointPath,
		"method":          cfg.HTTPMethod,
		"url":             cfg.FullURL,
		"retry_count":     cfg.RetryCount,
		"timeout":         cfg.Timeout,
		"current_node":    cfg.CurrentNodeID,
		"api_key_present": cfg.APIKey != "",
	})

	// Validasi URL
	if cfg.FullURL == "" {
		logJSON("error_empty_url", map[string]interface{}{
			"product_id": cfg.ProductID,
			"error":      "full URL is empty",
		})
		return nil, fmt.Errorf("full URL is empty for product %s", cfg.ProductID)
	}

	// Validasi RetryCount
	if cfg.RetryCount <= 0 {
		logJSON("warning_invalid_retry", map[string]interface{}{
			"product_id":  cfg.ProductID,
			"retry_count": cfg.RetryCount,
			"message":     "RetryCount is 0, setting to default 3",
		})
		cfg.RetryCount = 3
	}

	// Validasi Timeout
	if cfg.Timeout <= 0 {
		logJSON("warning_invalid_timeout", map[string]interface{}{
			"product_id": cfg.ProductID,
			"timeout":    cfg.Timeout,
			"message":    "Timeout is 0, setting to default 30",
		})
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

		logJSON("attempt", map[string]interface{}{
			"product_id":   cfg.ProductID,
			"attempt":      attempt,
			"total_retry":  cfg.RetryCount,
			"current_node": cfg.CurrentNodeID,
		})

		// Marshal body
		jsonBody, err := json.Marshal(enrichedBody)
		if err != nil {
			lastErr = fmt.Errorf("marshal error (attempt %d/%d): %w", attempt, cfg.RetryCount, err)
			logJSON("marshal_error", map[string]interface{}{
				"product_id": cfg.ProductID,
				"attempt":    attempt,
				"error":      err.Error(),
			})
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
			logJSON("request_error", map[string]interface{}{
				"product_id": cfg.ProductID,
				"attempt":    attempt,
				"error":      err.Error(),
			})
			if i < cfg.RetryCount-1 {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, lastErr
		}

		// Set headers
		setRequestHeaders(req, cfg)

		// Send request
		logJSON("sending_request", map[string]interface{}{
			"product_id": cfg.ProductID,
			"attempt":    attempt,
			"method":     cfg.HTTPMethod,
			"url":        cfg.FullURL,
		})

		startTime := time.Now()
		resp, err := client.Do(req)
		elapsed := time.Since(startTime)

		if err != nil {
			lastErr = fmt.Errorf("request failed (attempt %d/%d): %w", attempt, cfg.RetryCount, err)
			logJSON("http_error", map[string]interface{}{
				"product_id": cfg.ProductID,
				"attempt":    attempt,
				"error":      err.Error(),
				"elapsed_ms": elapsed.Milliseconds(),
			})
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
			logJSON("read_error", map[string]interface{}{
				"product_id":  cfg.ProductID,
				"attempt":     attempt,
				"status_code": resp.StatusCode,
				"error":       err.Error(),
			})
			if i < cfg.RetryCount-1 {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, lastErr
		}

		// Log response
		safeBody := safeJSON(bodyBytes)
		logJSON("response", map[string]interface{}{
			"product_id":  cfg.ProductID,
			"attempt":     attempt,
			"status_code": resp.StatusCode,
			"body_size":   len(bodyBytes),
			"elapsed_ms":  elapsed.Milliseconds(),
			"body":        json.RawMessage(safeBody),
		})

		// Store in execution results
		cfg.ExecutionResults["last_response"] = string(safeBody)
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

			logJSON("request_success", map[string]interface{}{
				"product_id":  cfg.ProductID,
				"attempt":     attempt,
				"status_code": resp.StatusCode,
				"elapsed_ms":  elapsed.Milliseconds(),
			})
			return bodyBytes, nil
		}

		// Handle non-2xx response
		lastErr = fmt.Errorf("HTTP %d (attempt %d/%d): %s", resp.StatusCode, attempt, cfg.RetryCount, string(bodyBytes))
		logJSON("http_error_status", map[string]interface{}{
			"product_id":  cfg.ProductID,
			"attempt":     attempt,
			"status_code": resp.StatusCode,
			"error":       string(safeBody),
		})

		if i < cfg.RetryCount-1 {
			time.Sleep(2 * time.Second)
		}
	}

	logJSON("all_retries_failed", map[string]interface{}{
		"product_id":     cfg.ProductID,
		"total_attempts": cfg.RetryCount,
		"last_error":     fmt.Sprintf("%v", lastErr),
	})

	return nil, fmt.Errorf("all retries exhausted (%d attempts): %w", cfg.RetryCount, lastErr)
}

// getEndpointPathFromConfig mengambil endpoint path dari workflow node berdasarkan current node ID
func getEndpointPathFromConfig(cfg *product.ProductConfig) string {
	logJSON("get_endpoint_path_start", map[string]interface{}{
		"product_id":           cfg.ProductID,
		"current_node_id":      cfg.CurrentNodeID,
		"workflow_nodes_count": len(cfg.WorkflowNodes),
		"adapter_endpoint":     cfg.AdapterEndpoint,
		"full_url":             cfg.FullURL,
	})

	// Jika tidak ada current node ID, return empty
	if cfg.CurrentNodeID == "" {
		logJSON("get_endpoint_path_no_node_id", map[string]interface{}{
			"product_id": cfg.ProductID,
			"message":    "CurrentNodeID is empty, cannot find endpoint path",
		})
		return ""
	}

	// Cari node dengan ID yang sesuai
	logJSON("get_endpoint_path_searching", map[string]interface{}{
		"product_id":      cfg.ProductID,
		"current_node_id": cfg.CurrentNodeID,
		"message":         "Searching for node in workflow nodes",
	})

	for _, node := range cfg.WorkflowNodes {
		if node.ID == cfg.CurrentNodeID {
			logJSON("get_endpoint_path_node_found", map[string]interface{}{
				"product_id":         cfg.ProductID,
				"current_node_id":    cfg.CurrentNodeID,
				"node_id":            node.ID,
				"has_adapter_config": node.AdapterConfig != nil,
			})

			// Jika node memiliki adapter config dan endpoint path
			if node.AdapterConfig != nil && node.AdapterConfig.EndpointPath != "" {
				endpointPath := node.AdapterConfig.EndpointPath
				logJSON("get_endpoint_path_success", map[string]interface{}{
					"product_id":      cfg.ProductID,
					"current_node_id": cfg.CurrentNodeID,
					"endpoint_path":   endpointPath,
					"source":          "node_adapter_config",
				})
				return endpointPath
			}

			logJSON("get_endpoint_path_node_has_no_adapter", map[string]interface{}{
				"product_id":      cfg.ProductID,
				"current_node_id": cfg.CurrentNodeID,
				"adapter_config":  node.AdapterConfig,
				"message":         "Node has no adapter config or endpoint path is empty",
			})
			break
		}
	}

	// Fallback 1: coba ambil dari adapter endpoint
	if cfg.AdapterEndpoint != "" {
		logJSON("get_endpoint_path_fallback_adapter", map[string]interface{}{
			"product_id":       cfg.ProductID,
			"current_node_id":  cfg.CurrentNodeID,
			"adapter_endpoint": cfg.AdapterEndpoint,
			"source":           "adapter_endpoint",
		})
		return cfg.AdapterEndpoint
	}

	// Fallback 2: extract dari FullURL
	if cfg.FullURL != "" {
		extractedPath := extractPathFromURL(cfg.FullURL)
		logJSON("get_endpoint_path_fallback_url", map[string]interface{}{
			"product_id":      cfg.ProductID,
			"current_node_id": cfg.CurrentNodeID,
			"full_url":        cfg.FullURL,
			"extracted_path":  extractedPath,
			"source":          "url_extraction",
		})
		return extractedPath
	}

	logJSON("get_endpoint_path_not_found", map[string]interface{}{
		"product_id":      cfg.ProductID,
		"current_node_id": cfg.CurrentNodeID,
		"message":         "No endpoint path found from any source",
	})

	return ""
}

// extractPathFromURL mengambil path dari URL
func extractPathFromURL(fullURL string) string {
	logJSON("extract_path_from_url_start", map[string]interface{}{
		"full_url": fullURL,
	})

	if fullURL == "" {
		logJSON("extract_path_from_url_empty", map[string]interface{}{
			"message": "URL is empty",
		})
		return ""
	}

	// Parse URL
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		logJSON("extract_path_from_url_parse_error", map[string]interface{}{
			"full_url": fullURL,
			"error":    err.Error(),
		})

		// Jika gagal parse, coba ambil setelah domain
		if strings.Contains(fullURL, "://") {
			parts := strings.SplitN(fullURL, "://", 2)
			if len(parts) == 2 {
				hostParts := strings.SplitN(parts[1], "/", 2)
				if len(hostParts) == 2 {
					path := "/" + hostParts[1]
					logJSON("extract_path_from_url_fallback", map[string]interface{}{
						"full_url":  fullURL,
						"extracted": path,
						"method":    "manual_parsing",
					})
					return path
				}
			}
		}

		logJSON("extract_path_from_url_failed", map[string]interface{}{
			"full_url": fullURL,
			"message":  "Failed to extract path",
		})
		return fullURL
	}

	// Return path + query string
	path := parsedURL.Path
	if parsedURL.RawQuery != "" {
		path = path + "?" + parsedURL.RawQuery
	}

	logJSON("extract_path_from_url_success", map[string]interface{}{
		"full_url":  fullURL,
		"path":      path,
		"raw_path":  parsedURL.Path,
		"raw_query": parsedURL.RawQuery,
		"has_query": parsedURL.RawQuery != "",
	})

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
		// Cari pattern setelah domain
		if strings.Contains(fullURL, "://") {
			parts := strings.SplitN(fullURL, "://", 2)
			if len(parts) == 2 {
				// Ambil setelah host
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
			// Hindari overwrite field yang sudah ada
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

func logJSON(event string, fields map[string]interface{}) {
	fields["event"] = event
	fields["timestamp"] = time.Now().Format(time.RFC3339Nano)
	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		log.Printf("[LOG_ERROR] failed to marshal log fields: %v", err)
		return
	}
	log.Println(string(b))
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
