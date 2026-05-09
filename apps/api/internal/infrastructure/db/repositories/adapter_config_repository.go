package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"seo-backend/internal/domain/adapter"
	"seo-backend/internal/domain/product"
)

type AdapterConfigRepository struct {
	db *sql.DB
}

// Constructor - implements interface
func NewAdapterConfigRepository(db *sql.DB) adapter.AdapterConfigRepository {
	return &AdapterConfigRepository{db: db}
}

// LoadForProduct - load config and attach to product
func (r *AdapterConfigRepository) LoadForProduct(ctx context.Context, data *product.Product) error {
	query := `
		SELECT id, product_id, endpoint_path, http_method,
			custom_headers, field_mapping, timeout_seconds,
			retry_count, created_at, updated_at
		FROM adapter_configs WHERE product_id = $1
	`

	var config product.AdapterConfig
	var customHeadersJSON []byte
	var fieldMappingJSON []byte

	err := r.db.QueryRowContext(ctx, query, data.ID).Scan(
		&config.ID, &config.ProductID, &config.EndpointPath,
		&config.HTTPMethod, &customHeadersJSON, &fieldMappingJSON,
		&config.TimeoutSeconds, &config.RetryCount,
		&config.CreatedAt, &config.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil // No config is fine
		}
		return fmt.Errorf("failed to load adapter config: %w", err)
	}

	// Unmarshal JSON fields
	if len(customHeadersJSON) > 0 {
		if err := json.Unmarshal(customHeadersJSON, &config.CustomHeaders); err != nil {
			return fmt.Errorf("failed to unmarshal custom headers: %w", err)
		}
	}

	if len(fieldMappingJSON) > 0 {
		if err := json.Unmarshal(fieldMappingJSON, &config.FieldMapping); err != nil {
			return fmt.Errorf("failed to unmarshal field mapping: %w", err)
		}
	}

	data.AdapterConfig = &config
	return nil
}

// GetOrDefault - get config with default values
func (r *AdapterConfigRepository) GetOrDefault(ctx context.Context, productID string) (*product.AdapterConfig, error) {
	config := &product.AdapterConfig{
		HTTPMethod:     "GET",
		TimeoutSeconds: 10,
		RetryCount:     1,
		CustomHeaders:  "",
		FieldMapping:   "",
	}

	query := `
		SELECT 
			COALESCE(endpoint_path, ''),
			COALESCE(http_method, 'GET'),
			COALESCE(custom_headers, '{}'),
			COALESCE(field_mapping, '{}'),
			COALESCE(timeout_seconds, 10),
			COALESCE(retry_count, 1)
		FROM adapter_configs 
		WHERE product_id = $1
	`

	var customHeadersJSON []byte
	var fieldMappingJSON []byte

	err := r.db.QueryRowContext(ctx, query, productID).Scan(
		&config.EndpointPath,
		&config.HTTPMethod,
		&customHeadersJSON,
		&fieldMappingJSON,
		&config.TimeoutSeconds,
		&config.RetryCount,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get adapter config: %w", err)
	}

	if err == nil {
		if len(customHeadersJSON) > 0 {
			json.Unmarshal(customHeadersJSON, &config.CustomHeaders)
		}
		if len(fieldMappingJSON) > 0 {
			json.Unmarshal(fieldMappingJSON, &config.FieldMapping)
		}
	}

	return config, nil
}

// GetByProductID - get config by product ID (without defaults)
func (r *AdapterConfigRepository) GetByProductID(ctx context.Context, productID string) (*product.AdapterConfig, error) {
	query := `
		SELECT id, product_id, endpoint_path, http_method,
			custom_headers, field_mapping, timeout_seconds,
			retry_count, created_at, updated_at
		FROM adapter_configs WHERE product_id = $1
	`

	var config product.AdapterConfig
	var customHeadersJSON []byte
	var fieldMappingJSON []byte

	err := r.db.QueryRowContext(ctx, query, productID).Scan(
		&config.ID, &config.ProductID, &config.EndpointPath,
		&config.HTTPMethod, &customHeadersJSON, &fieldMappingJSON,
		&config.TimeoutSeconds, &config.RetryCount,
		&config.CreatedAt, &config.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get adapter config: %w", err)
	}

	// Unmarshal JSON fields
	if len(customHeadersJSON) > 0 {
		json.Unmarshal(customHeadersJSON, &config.CustomHeaders)
	}
	if len(fieldMappingJSON) > 0 {
		json.Unmarshal(fieldMappingJSON, &config.FieldMapping)
	}

	return &config, nil
}

// InsertWithTx - insert config with transaction
func (r *AdapterConfigRepository) InsertWithTx(ctx context.Context, tx *sql.Tx, productID string, config *product.AdapterConfig) error {
	if tx == nil {
		return fmt.Errorf("transaction is required for InsertWithTx")
	}

	customHeadersJSON, _ := json.Marshal(config.CustomHeaders)
	fieldMappingJSON, _ := json.Marshal(config.FieldMapping)

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
			custom_headers, field_mapping, timeout_seconds, retry_count,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	`

	_, err := tx.ExecContext(ctx, query,
		productID, config.EndpointPath, config.HTTPMethod,
		customHeadersJSON, fieldMappingJSON, timeoutSeconds, retryCount,
	)

	if err != nil {
		return fmt.Errorf("failed to insert adapter config: %w", err)
	}

	return nil
}

// UpdateWithTx - update config with transaction
func (r *AdapterConfigRepository) UpdateWithTx(ctx context.Context, tx *sql.Tx, productID string, updates map[string]interface{}) error {
	if tx == nil {
		return fmt.Errorf("transaction is required for UpdateWithTx")
	}

	if len(updates) == 0 {
		return nil
	}

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
			switch key {
			case "customHeaders":
				jsonValue, err := json.Marshal(value)
				if err != nil {
					return fmt.Errorf("failed to marshal custom headers: %w", err)
				}
				setClauses = append(setClauses, fmt.Sprintf("%s = $%d", dbField, argIndex))
				args = append(args, jsonValue)

			case "fieldMapping":
				jsonValue, err := json.Marshal(value)
				if err != nil {
					return fmt.Errorf("failed to marshal field mapping: %w", err)
				}
				setClauses = append(setClauses, fmt.Sprintf("%s = $%d", dbField, argIndex))
				args = append(args, jsonValue)

			default:
				setClauses = append(setClauses, fmt.Sprintf("%s = $%d", dbField, argIndex))
				args = append(args, value)
			}
			argIndex++
		}
	}

	if len(setClauses) == 0 {
		return nil
	}

	// Always update timestamp
	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	// Add WHERE clause
	args = append(args, productID)
	query := fmt.Sprintf("UPDATE adapter_configs SET %s WHERE product_id = $%d",
		strings.Join(setClauses, ", "), argIndex)

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update adapter config: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("adapter config for product %s not found", productID)
	}

	return nil
}

// internal/repositories/adapter_config_repository.go

// DeleteWithTx - delete adapter config by product ID with transaction
func (r *AdapterConfigRepository) DeleteWithTx(ctx context.Context, tx *sql.Tx, productID string) error {
	if tx == nil {
		return fmt.Errorf("transaction is required for DeleteWithTx")
	}

	query := `DELETE FROM adapter_configs WHERE product_id = $1`

	result, err := tx.ExecContext(ctx, query, productID)
	if err != nil {
		return fmt.Errorf("failed to delete adapter config: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	// Return nil even if no rows deleted (config might not exist)
	if rowsAffected == 0 {
		return nil // No config found, not an error
	}

	return nil
}
