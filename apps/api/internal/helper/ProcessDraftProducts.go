package helper

import (
	"bytes"
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
}

func NewPostService() *PostService {
	return &PostService{}
}

func (s *PostService) ProcessDraftProducts(draft draft.DraftDataPost) (
	[]map[string]interface{},
	bool,
	bool,
	error,
) {

	var postResults []map[string]interface{}

	someFailed := false
	allFailed := true

	for _, productID := range draft.TargetProducts {

		fmt.Printf("=======================,%v", productID)
		result, err := s.processSingleProduct(draft, productID)
		postResults = append(postResults, result)

		if err != nil {
			someFailed = true
		} else {
			allFailed = false
		}
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
		return map[string]interface{}{
			"product": productID,
			"success": false,
			"error":   "Product configuration not found",
		}, err
	}

	requestBody, err := s.buildRequestBody(cfg, draft)
	if err != nil {
		return map[string]interface{}{
			"product": productID,
			"success": false,
			"error":   err.Error(),
		}, err
	}

	response, err := s.sendWithRetry(cfg, requestBody)

	if err != nil {
		return map[string]interface{}{
			"product": productID,
			"success": false,
			"error":   err.Error(),
		}, err
	}

	s.markProductSynced(cfg.ProductID)

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
		log.Printf("================ERRRRRR %s", err)
		log.Printf("================ERRRRRR %s", productID)
		return nil, err
	}

	// 2. Default values adapter config
	cfg.HTTPMethod = "POST"
	cfg.Timeout = 30
	cfg.RetryCount = 3

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

		if err != nil {
			// fallback aman kalau config tidak ada
			cfg.AdapterEndpoint = ""
			cfg.FieldMappingStr = "{}"
			cfg.CustomHeadersStr = "{}"
			cfg.HTTPMethod = "POST"
			cfg.Timeout = 30
			cfg.RetryCount = 3
		}

	} else {
		log.Printf("[INFO] Skipping adapter_configs because API key is empty (product: %s)", cfg.ProductID)

		// safe default
		cfg.AdapterEndpoint = ""
		cfg.FieldMappingStr = "{}"
		cfg.CustomHeadersStr = "{}"
	}

	// 4. Build full URL
	cfg.FullURL = cfg.APIEndpoint + cfg.AdapterEndpoint

	return &cfg, nil
}

func (s *PostService) buildRequestBody(cfg *ProductConfig, draft draft.DraftDataPost) (map[string]interface{}, error) {

	var fieldMapping map[string]interface{}

	if cfg.FieldMappingStr != "" {
		json.Unmarshal([]byte(cfg.FieldMappingStr), &fieldMapping)
	}

	requestBody := make(map[string]interface{})

	for key, value := range fieldMapping {

		if str, ok := value.(string); ok {

			str = strings.ReplaceAll(str, "{title}", draft.Title)
			str = strings.ReplaceAll(str, "{topic}", draft.Topic)
			str = strings.ReplaceAll(str, "{content}", draft.Article)
			str = strings.ReplaceAll(str, "{image_url}", *draft.ImageURL)

			requestBody[key] = str
		} else {
			requestBody[key] = value
		}
	}

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

	client := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
	}

	for i := 0; i < cfg.RetryCount; i++ {

		// 1. Marshal body
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		log.Printf("[REQUEST BODY] %s", string(jsonBody))

		// 2. Create request
		req, err := http.NewRequest(
			cfg.HTTPMethod,
			cfg.FullURL,
			bytes.NewReader(jsonBody),
		)
		if err != nil {
			lastErr = err
			continue
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

		// 3. Execute request
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(2 * time.Second)
			continue
		}

		bodyBytes, err := func() ([]byte, error) {
			defer resp.Body.Close()
			return io.ReadAll(resp.Body)
		}()

		if err != nil {
			lastErr = err
			time.Sleep(2 * time.Second)
			continue
		}

		// success
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return bodyBytes, nil
		}

		lastErr = fmt.Errorf(
			"HTTP %d %s: %s",
			resp.StatusCode,
			cfg.FullURL,
			string(bodyBytes),
		)

		time.Sleep(2 * time.Second)
	}

	return nil, lastErr
}

func (s *PostService) markProductSynced(productID string) {

	_, err := database.GetDB().Exec(`
		UPDATE products
		SET sync_status='connected',
			last_sync=NOW()
		WHERE id=$1
	`, productID)

	if err != nil {
		log.Printf("[WARN] sync update failed: %v", err)
	}
}

func resolveTemplate(v string, data map[string]string) string {
	for k, val := range data {
		v = strings.ReplaceAll(v, "{{"+k+"}}", val)
	}
	return v
}
