package helper

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"seo-backend/internal/database"
	"seo-backend/internal/domain/draft"
	"strings"
	"time"
)

type ProductConfig struct {
	ProductID       string
	APIEndpoint     string
	APIKey          string
	AdapterEndpoint string
	HTTPMethod      string
	FullURL         string

	FieldMappingStr  string
	CustomHeadersStr string
	CustomHeaders    map[string]string

	Timeout    int
	RetryCount int
}

type PostService struct {
	db *sql.DB
}

func NewPostService(db *sql.DB) *PostService {
	return &PostService{
		db: db,
	}
}

func (s *PostService) updateConnectionStatus(
	productID string,
	isConnected bool,
) error {

	status := "connected"
	if !isConnected {
		status = "pending"
	}

	_, err := s.db.Exec(`
		UPDATE products 
		SET status = $1, last_sync = NOW(), updated_at = NOW()
		WHERE id = $2
	`, status, productID)

	if err != nil {
		return fmt.Errorf("failed to update product status for product %s: %w", productID, err)
	}

	return nil
}

func (s *PostService) ProcessDraftProducts(draft draft.DraftDataPost) (
	[]map[string]interface{},
	bool,
	bool,
	error,
) {

	// Validasi target products
	log.Printf("TargetProducts : %v", draft.TargetProducts)
	if len(draft.TargetProducts) == 0 {
		return nil, false, true, fmt.Errorf("target products is required and cannot be empty")
	}

	var postResults []map[string]interface{}

	someFailed := false
	allFailed := true

	log.Printf("=============%v", draft)

	for _, productID := range draft.TargetProducts {

		fmt.Printf("=======================,%v", productID)

		result, err := s.processSingleProduct(draft, productID)

		// update status regardless of error
		if updateErr := s.updateConnectionStatus(productID, err == nil); updateErr != nil {
			log.Printf("failed update product status for %s: %v", productID, updateErr)
			// Don't return this error, just log it
		}

		postResults = append(postResults, result)

		if err != nil {
			someFailed = true
			log.Printf("failed to process product %s: %v", productID, err)
		} else {
			allFailed = false
		}
	}

	// Return error if all products failed
	if allFailed && len(draft.TargetProducts) > 0 {
		return postResults, someFailed, allFailed, fmt.Errorf("all products failed to process")
	}

	return postResults, someFailed, allFailed, nil
}

func (s *PostService) processSingleProduct(
	draft draft.DraftDataPost,
	productID string,
) (map[string]interface{}, error) {

	log.Printf("[START] Processing product: %s", productID)

	cfg, err := s.getProductConfig(productID)
	if err != nil {
		errMsg := fmt.Sprintf("failed to get product config for %s: %v", productID, err)
		log.Printf("[ERROR] %s", errMsg)
		return map[string]interface{}{
			"product": productID,
			"success": false,
			"error":   errMsg,
		}, fmt.Errorf(errMsg)
	}

	requestBody, err := s.buildRequestBody(cfg, draft)
	if err != nil {
		errMsg := fmt.Sprintf("failed to build request body for %s: %v", productID, err)
		log.Printf("[ERROR] %s", errMsg)
		return map[string]interface{}{
			"product": productID,
			"success": false,
			"error":   errMsg,
		}, fmt.Errorf(errMsg)
	}

	response, err := s.sendWithRetry(cfg, requestBody)
	if err != nil {
		errMsg := fmt.Sprintf("failed to send request for %s after retries: %v", productID, err)
		log.Printf("[ERROR] %s", errMsg)
		return map[string]interface{}{
			"product": productID,
			"success": false,
			"error":   errMsg,
		}, fmt.Errorf(errMsg)
	}

	if err := s.markProductSynced(cfg.ProductID); err != nil {
		log.Printf("[WARN] failed to mark product %s as synced: %v", cfg.ProductID, err)
		// Don't return error for this, as the main operation succeeded
	}

	return map[string]interface{}{
		"product":    productID,
		"success":    true,
		"response":   string(response),
		"product_id": cfg.ProductID,
	}, nil
}

func (s *PostService) getProductConfig(productID string) (*ProductConfig, error) {

	var cfg ProductConfig

	// 1. Ambil basic product dulu
	err := database.GetDB().QueryRow(`
		SELECT
			id,
			api_endpoint,
			COALESCE(api_key_encrypted, '')
		FROM products
		WHERE id = $1
	`, productID).Scan(
		&cfg.ProductID,
		&cfg.APIEndpoint,
		&cfg.APIKey,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with ID %s not found", productID)
		}
		log.Printf("================ERRRRRR %s", err)
		log.Printf("================ERRRRRR %s", productID)
		return nil, fmt.Errorf("failed to query product %s: %w", productID, err)
	}

	// Validate required fields
	if cfg.APIEndpoint == "" {
		return nil, fmt.Errorf("product %s has empty API endpoint", productID)
	}

	// 2. Default values adapter config
	cfg.HTTPMethod = "POST"
	cfg.Timeout = 30
	cfg.RetryCount = 3
	cfg.CustomHeaders = make(map[string]string)

	// 3. HANYA ambil adapter_configs kalau API key ADA
	if cfg.APIKey != "" {

		err = database.GetDB().QueryRow(`
			SELECT
				COALESCE(endpoint_path, ''),
				COALESCE(http_method, 'POST'),
				COALESCE(field_mapping, '{}'),
				COALESCE(custom_headers, '{}'),
				COALESCE(timeout_seconds, 30),
				COALESCE(retry_count, 3)
			FROM adapter_configs
			WHERE product_id = $1
		`, cfg.ProductID).Scan(
			&cfg.AdapterEndpoint,
			&cfg.HTTPMethod,
			&cfg.FieldMappingStr,
			&cfg.CustomHeadersStr,
			&cfg.Timeout,
			&cfg.RetryCount,
		)

		if err != nil && err != sql.ErrNoRows {
			// Real error, not just no rows
			return nil, fmt.Errorf("failed to query adapter config for product %s: %w", productID, err)
		}

		if cfg.CustomHeadersStr != "" && cfg.CustomHeadersStr != "{}" {

			raw := cfg.CustomHeadersStr

			log.Printf("RAW CUSTOM HEADERS: %s", raw)

			// first try direct object
			err := json.Unmarshal([]byte(raw), &cfg.CustomHeaders)

			if err != nil {

				// try nested string
				var nested string

				if err2 := json.Unmarshal([]byte(raw), &nested); err2 == nil {

					// parse nested json
					if err3 := json.Unmarshal([]byte(nested), &cfg.CustomHeaders); err3 != nil {
						log.Printf(
							"[WARN] failed nested parse custom headers for product %s: %v",
							productID,
							err3,
						)

						cfg.CustomHeaders = make(map[string]string)
					}

				} else {

					log.Printf(
						"[WARN] failed to parse custom headers for product %s: %v",
						productID,
						err,
					)

					cfg.CustomHeaders = make(map[string]string)
				}
			}
		}

	} else {
		log.Printf("[INFO] Skipping adapter_configs because API key is empty (product: %s)", cfg.ProductID)

		// safe default
		cfg.AdapterEndpoint = ""
		cfg.FieldMappingStr = "{}"
		cfg.CustomHeadersStr = "{}"
	}

	// Validate HTTP method
	validMethods := map[string]bool{"GET": true, "POST": true, "PUT": true, "PATCH": true, "DELETE": true}
	if !validMethods[strings.ToUpper(cfg.HTTPMethod)] {
		cfg.HTTPMethod = "POST"
		log.Printf("[WARN] invalid HTTP method for product %s, defaulting to POST", productID)
	}

	// 4. Build full URL
	cfg.FullURL = strings.TrimRight(cfg.APIEndpoint, "/") + "/" + strings.TrimLeft(cfg.AdapterEndpoint, "/")

	return &cfg, nil
}

func (s *PostService) buildRequestBody(
	cfg *ProductConfig,
	draft draft.DraftDataPost,
) (map[string]interface{}, error) {

	var fieldMapping map[string]interface{}

	// Parse field mapping
	if cfg.FieldMappingStr != "" && cfg.FieldMappingStr != "{}" {

		raw := strings.TrimSpace(cfg.FieldMappingStr)

		fmt.Println("RAW FIELD MAPPING:", raw)

		// Try direct JSON object
		if err := json.Unmarshal([]byte(raw), &fieldMapping); err != nil {

			// Handle double encoded JSON string
			var nested string

			if err2 := json.Unmarshal([]byte(raw), &nested); err2 == nil {

				if err3 := json.Unmarshal([]byte(nested), &fieldMapping); err3 != nil {
					return nil, fmt.Errorf(
						"failed nested parse field mapping: %w",
						err3,
					)
				}

			} else {
				return nil, fmt.Errorf(
					"failed to parse field mapping: %w",
					err,
				)
			}
		}
	}

	requestBody := make(map[string]interface{})

	// Validate draft
	if strings.TrimSpace(draft.Title) == "" {
		return nil, fmt.Errorf("draft title is required")
	}

	if strings.TrimSpace(draft.Article) == "" {
		return nil, fmt.Errorf("draft article content is required")
	}

	// Replace placeholders
	for key, value := range fieldMapping {

		switch v := value.(type) {

		case string:

			v = strings.ReplaceAll(v, "{title}", draft.Title)
			v = strings.ReplaceAll(v, "{topic}", draft.Topic)
			v = strings.ReplaceAll(v, "{content}", draft.Article)

			if draft.ImageURL != nil {
				v = strings.ReplaceAll(v, "{image_url}", *draft.ImageURL)
			}

			requestBody[key] = v

		default:
			requestBody[key] = value
		}
	}

	// Default request body if field mapping empty
	if len(requestBody) == 0 {

		requestBody = map[string]interface{}{
			"title":   draft.Title,
			"content": draft.Article,
			"topic":   draft.Topic,
		}

		if draft.ImageURL != nil {
			requestBody["image_url"] = *draft.ImageURL
		}
	}

	return requestBody, nil
}

func (s *PostService) sendWithRetry(cfg *ProductConfig, body map[string]interface{}) ([]byte, error) {

	var lastErr error

	// Validate config
	if cfg.FullURL == "" {
		return nil, fmt.Errorf("full URL is empty")
	}

	client := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
	}

	for i := 0; i < cfg.RetryCount; i++ {

		// 1. Marshal body
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

		// 2. Create request
		req, err := http.NewRequest(
			cfg.HTTPMethod,
			cfg.FullURL,
			bytes.NewReader(jsonBody),
		)
		if err != nil {
			lastErr = fmt.Errorf("failed to create request (attempt %d/%d): %w", i+1, cfg.RetryCount, err)
			if i < cfg.RetryCount-1 {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, lastErr
		}

		// default headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		// =========================
		// CUSTOM HEADERS PROCESS
		// =========================
		if cfg.CustomHeaders != nil {

			obj := cfg.CustomHeaders

			for k, v := range obj {

				v = resolveTemplate(v, map[string]string{
					"api_key": cfg.APIKey,
				})

				kLower := strings.ToLower(strings.TrimSpace(k))
				vTrim := strings.TrimSpace(v)

				if kLower == "" {
					continue
				}

				// =========================
				// AUTH HANDLING
				// =========================

				// API KEY header
				if kLower == "x-api-key" {
					if vTrim == "" || vTrim == "{{api_key}}" {
						v = cfg.APIKey
					}
				}

				// BEARER TOKEN
				if kLower == "authorization" {

					if vTrim == "" ||
						vTrim == "Bearer" ||
						vTrim == "Bearer {{api_key}}" {

						v = "Bearer " + cfg.APIKey
					}
				}

				req.Header.Set(k, v)
			}
		}

		// =========================
		// FALLBACK AUTH (SAFEGUARD)
		// =========================
		if req.Header.Get("Authorization") == "" && cfg.APIKey != "" {
			req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
		}

		// debug header penting
		log.Printf("[AUTH HEADER] %s", req.Header.Get("Authorization"))
		log.Printf("[REQUEST URL] %s %s", cfg.HTTPMethod, cfg.FullURL)

		// 3. Execute request
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

		if err != nil {
			lastErr = fmt.Errorf("failed to read response body (attempt %d/%d): %w", i+1, cfg.RetryCount, err)
			if i < cfg.RetryCount-1 {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, lastErr
		}

		// success
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return bodyBytes, nil
		}

		lastErr = fmt.Errorf(
			"HTTP %d %s (attempt %d/%d): %s",
			resp.StatusCode,
			cfg.FullURL,
			i+1,
			cfg.RetryCount,
			string(bodyBytes),
		)

		log.Printf("[ERROR] %v", lastErr)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("all retries exhausted (%d attempts): %w", cfg.RetryCount, lastErr)
}

func (s *PostService) markProductSynced(productID string) error {

	result, err := database.GetDB().Exec(`
		UPDATE products
		SET sync_status='connected',
			last_sync=NOW()
		WHERE id=$1
	`, productID)

	if err != nil {
		return fmt.Errorf("failed to update sync status for product %s: %w", productID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected for product %s: %w", productID, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product %s not found when updating sync status", productID)
	}

	log.Printf("[INFO] Product %s marked as synced successfully", productID)
	return nil
}

func resolveTemplate(v string, data map[string]string) string {
	result := v
	for k, val := range data {
		result = strings.ReplaceAll(result, "{{"+k+"}}", val)
	}
	return result
}
