package helper

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"seo-backend/internal/database"
	"strings"
)

type ProductConfig struct {
	ProductID       string
	APIEndpoint     string
	APIKey          string
	AdapterEndpoint string
	HTTPMethod      string
	FullURL         string

	FieldMappingStr  string
	MetaConfigStr    string
	SitemapConfigStr string
	BaseURL          string
	CustomHeadersStr string
	CustomHeaders    map[string]string

	Timeout    int
	RetryCount int
}

func (s *PostService) getProductConfig(productID string) (*ProductConfig, error) {
	var cfg ProductConfig

	log.Printf("========== GET PRODUCT CONFIG ==========")
	log.Printf("PRODUCT ID: %s", productID)

	// Get basic product info
	err := database.GetDB().QueryRow(`
		SELECT id, api_endpoint, COALESCE(api_key_encrypted, '')
		FROM products
		WHERE id = $1
	`, productID).Scan(&cfg.ProductID, &cfg.APIEndpoint, &cfg.APIKey)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[ERROR] PRODUCT NOT FOUND: %s", productID)
			return nil, fmt.Errorf("product with ID %s not found", productID)
		}
		log.Printf("[ERROR] FAILED QUERY PRODUCT: %v", err)
		return nil, fmt.Errorf("failed to query product %s: %w", productID, err)
	}

	log.Printf("PRODUCT FOUND")
	log.Printf("API ENDPOINT: %s", cfg.APIEndpoint)

	if cfg.APIEndpoint == "" {
		log.Printf("[ERROR] EMPTY API ENDPOINT")
		return nil, fmt.Errorf("product %s has empty API endpoint", productID)
	}

	// Set default values
	cfg.HTTPMethod = DefaultHTTPMethod
	cfg.Timeout = DefaultTimeout
	cfg.RetryCount = DefaultRetryCount
	cfg.CustomHeaders = make(map[string]string)
	cfg.MetaConfigStr = "{}"
	cfg.SitemapConfigStr = "{}"

	// Load adapter config if API key exists
	if cfg.APIKey != "" {
		if err := s.loadAdapterConfig(&cfg); err != nil {
			return nil, err
		}
	} else {
		log.Printf("[INFO] SKIPPING ADAPTER CONFIGS BECAUSE API KEY EMPTY")
		cfg.AdapterEndpoint = ""
		cfg.FieldMappingStr = "{}"
		cfg.CustomHeadersStr = "{}"
		cfg.MetaConfigStr = "{}"
		cfg.SitemapConfigStr = "{}"
	}

	// Validate HTTP method
	if !validHTTPMethods[strings.ToUpper(cfg.HTTPMethod)] {
		log.Printf("[WARN] INVALID HTTP METHOD: %s", cfg.HTTPMethod)
		cfg.HTTPMethod = DefaultHTTPMethod
		log.Printf("[INFO] DEFAULTING TO POST")
	}

	// Build full URL
	cfg.FullURL = strings.TrimRight(cfg.APIEndpoint, "/") + "/" + strings.TrimLeft(cfg.AdapterEndpoint, "/")
	log.Printf("FULL URL: %s", cfg.FullURL)
	log.Printf("========== PRODUCT CONFIG READY ==========")

	return &cfg, nil
}

func (s *PostService) loadAdapterConfig(cfg *ProductConfig) error {
	log.Printf("API KEY EXISTS => LOADING ADAPTER CONFIG")

	err := database.GetDB().QueryRow(`
		SELECT
			COALESCE(endpoint_path, ''),
			COALESCE(http_method, 'POST'),
			COALESCE(field_mapping, '{}'),
			COALESCE(custom_headers, '{}'),
			COALESCE(meta_config, '{}'),
			COALESCE(sitemap_config, '{}'),
			COALESCE(timeout_seconds, 60),
			COALESCE(retry_count, 3)
		FROM adapter_configs
		WHERE product_id = $1
	`, cfg.ProductID).Scan(
		&cfg.AdapterEndpoint,
		&cfg.HTTPMethod,
		&cfg.FieldMappingStr,
		&cfg.CustomHeadersStr,
		&cfg.MetaConfigStr,
		&cfg.SitemapConfigStr,
		&cfg.Timeout,
		&cfg.RetryCount,
	)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("[ERROR] FAILED QUERY ADAPTER CONFIG: %v", err)
		return fmt.Errorf("failed to query adapter config for product %s: %w", cfg.ProductID, err)
	}

	log.Printf("ADAPTER CONFIG LOADED")
	log.Printf("HTTP METHOD: %s", cfg.HTTPMethod)
	log.Printf("FIELD MAPPING: %s", cfg.FieldMappingStr)
	log.Printf("META CONFIG: %s", cfg.MetaConfigStr)
	log.Printf("SITEMAP CONFIG: %s", cfg.SitemapConfigStr)

	// Parse custom headers
	if err := s.parseCustomHeaders(cfg); err != nil {
		log.Printf("[WARN] Failed to parse custom headers: %v", err)
		cfg.CustomHeaders = make(map[string]string)
	}

	return nil
}

func (s *PostService) parseCustomHeaders(cfg *ProductConfig) error {
	if cfg.CustomHeadersStr == "" || cfg.CustomHeadersStr == "{}" {
		return nil
	}

	raw := cfg.CustomHeadersStr
	log.Printf("RAW CUSTOM HEADERS: %s", raw)

	// Try direct object
	err := json.Unmarshal([]byte(raw), &cfg.CustomHeaders)
	if err != nil {
		log.Printf("[WARN] DIRECT PARSE FAILED: %v", err)

		// Try nested string
		var nested string
		if err2 := json.Unmarshal([]byte(raw), &nested); err2 == nil {
			log.Printf("DOUBLE ENCODED CUSTOM HEADERS DETECTED")
			if err3 := json.Unmarshal([]byte(nested), &cfg.CustomHeaders); err3 != nil {
				log.Printf("[WARN] FAILED NESTED PARSE CUSTOM HEADERS: %v", err3)
				return err3
			}
		} else {
			log.Printf("[WARN] FAILED PARSE CUSTOM HEADERS: %v", err)
			return err
		}
	}

	return nil
}
