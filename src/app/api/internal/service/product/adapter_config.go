package product

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"seo-backend/internal/models"
	"strings"
	"time"
)

type AdapterConfigRepo struct {
	db *sql.DB
}

func NewAdapterConfigRepo(db *sql.DB) *AdapterConfigRepo {
	return &AdapterConfigRepo{db: db}
}

func (r *AdapterConfigRepo) LoadForProduct(product *models.Product) error {
	query := `
		SELECT id, product_id, endpoint_path, http_method,
			custom_headers, field_mapping, timeout_seconds,
			retry_count, created_at, updated_at
		FROM adapter_configs WHERE product_id = $1
	`

	var config models.AdapterConfig
	var customHeadersJSON []byte

	err := r.db.QueryRow(query, product.ID).Scan(
		&config.ID, &config.ProductID, &config.EndpointPath,
		&config.HTTPMethod, &customHeadersJSON, &config.FieldMapping,
		&config.TimeoutSeconds, &config.RetryCount,
		&config.CreatedAt, &config.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil // No config is fine
		}
		return err
	}

	json.Unmarshal(customHeadersJSON, &config.CustomHeaders)
	product.AdapterConfig = &config
	return nil
}

func (r *AdapterConfigRepo) GetOrDefault(productID string) *AdapterConfig {
	config := &AdapterConfig{
		HTTPMethod:     "GET",
		TimeoutSeconds: 10,
		RetryCount:     1,
		CustomHeaders:  make(map[string]string),
	}

	query := `
		SELECT 
			COALESCE(endpoint_path, ''), 
			COALESCE(http_method, 'GET'),
			COALESCE(custom_headers, '{}'),
			COALESCE(timeout_seconds, 10),
			COALESCE(retry_count, 1)
		FROM adapter_configs 
		WHERE product_id = $1
	`

	var customHeadersJSON []byte
	err := r.db.QueryRow(query, productID).Scan(
		&config.EndpointPath,
		&config.HTTPMethod,
		&customHeadersJSON,
		&config.TimeoutSeconds,
		&config.RetryCount,
	)

	if err == nil {
		json.Unmarshal(customHeadersJSON, &config.CustomHeaders)
	}

	return config
}

func (r *AdapterConfigRepo) Insert(tx *sql.Tx, productID string, config *models.AdapterConfig) error {
	customHeadersJSON, _ := json.Marshal(config.CustomHeaders)

	timeoutSeconds := config.TimeoutSeconds
	if timeoutSeconds == 0 {
		timeoutSeconds = 30
	}
	retryCount := config.RetryCount
	if retryCount == 0 {
		retryCount = 3
	}

	query := `
		INSERT INTO adapter_configs (
			product_id, endpoint_path, http_method,
			custom_headers, field_mapping, timeout_seconds, retry_count
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := tx.Exec(
		query,
		productID, config.EndpointPath, config.HTTPMethod,
		customHeadersJSON, config.FieldMapping, timeoutSeconds, retryCount,
	)
	return err
}

func (r *AdapterConfigRepo) Update(tx *sql.Tx, productID string, updates map[string]interface{}) error {
	setClauses := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1

	fieldMap := map[string]string{
		"endpointPath":   "endpoint_path",
		"httpMethod":     "http_method",
		"customHeaders":  "custom_headers",
		"fieldMapping":   "field_mapping",
		"timeoutSeconds": "timeout_seconds",
		"retryCount":     "retry_count",
	}

	for key, value := range updates {
		if dbField, ok := fieldMap[key]; ok {
			if key == "customHeaders" {
				jsonValue, _ := json.Marshal(value)
				setClauses = append(setClauses, fmt.Sprintf("%s = $%d", dbField, argIndex))
				args = append(args, jsonValue)
			} else {
				setClauses = append(setClauses, fmt.Sprintf("%s = $%d", dbField, argIndex))
				args = append(args, value)
			}
			argIndex++
		}
	}

	if len(setClauses) == 0 {
		return nil
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	args = append(args, productID)
	query := fmt.Sprintf("UPDATE adapter_configs SET %s WHERE product_id = $%d",
		strings.Join(setClauses, ", "), argIndex)

	_, err := tx.Exec(query, args...)
	return err
}
