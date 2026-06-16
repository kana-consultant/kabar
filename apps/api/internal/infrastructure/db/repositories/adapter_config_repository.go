package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"seo-backend/internal/domain/adapter"
	"seo-backend/internal/domain/product"
	"seo-backend/internal/helper"
)

type AdapterConfigRepository struct {
	db *sql.DB
}

func NewAdapterConfigRepository(db *sql.DB) adapter.AdapterConfigRepository {
	return &AdapterConfigRepository{db: db}
}

func (r *AdapterConfigRepository) LoadForProduct(ctx context.Context, data *product.Product) error {
	query := `
		SELECT id, product_id, endpoint_path, http_method,
			custom_headers, field_mapping, response_mapping, meta_config, sitemap_config, timeout_seconds,
			retry_count, created_at, updated_at
		FROM adapter_configs WHERE product_id = $1
	`

	var config product.AdapterConfig
	var customHeadersJSON []byte
	var fieldMappingJSON []byte
	var responseMappingJSON []byte
	var metaConfigJSON []byte
	var sitemapConfigJSON []byte

	err := r.db.QueryRowContext(ctx, query, data.ID).Scan(
		&config.ID, &config.ProductID, &config.EndpointPath,
		&config.HTTPMethod, &customHeadersJSON, &fieldMappingJSON,
		&responseMappingJSON,
		&metaConfigJSON, &sitemapConfigJSON,
		&config.TimeoutSeconds, &config.RetryCount,
		&config.CreatedAt, &config.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return fmt.Errorf("failed to load adapter config: %w", err)
	}

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

	if len(responseMappingJSON) > 0 {
		if err := json.Unmarshal(responseMappingJSON, &config.ResponseMapping); err != nil {
			return fmt.Errorf("failed to unmarshal response mapping: %w", err)
		}
	}

	if len(metaConfigJSON) > 0 {
		if err := json.Unmarshal(metaConfigJSON, &config.MetaConfig); err != nil {
			return fmt.Errorf("failed to unmarshal meta config: %w", err)
		}
	}

	if len(sitemapConfigJSON) > 0 {
		if err := json.Unmarshal(sitemapConfigJSON, &config.SitemapConfig); err != nil {
			return fmt.Errorf("failed to unmarshal sitemap config: %w", err)
		}
	}

	data.AdapterConfig = &config
	return nil
}

func (r *AdapterConfigRepository) GetOrDefault(ctx context.Context, productID string) (*product.AdapterConfig, error) {
	config := &product.AdapterConfig{
		HTTPMethod:      "GET",
		TimeoutSeconds:  10,
		RetryCount:      1,
		CustomHeaders:   "",
		FieldMapping:    "",
		ResponseMapping: "",
		MetaConfig:      "",
		SitemapConfig:   "",
	}

	query := `
		SELECT 
			COALESCE(endpoint_path, ''),
			COALESCE(http_method, 'GET'),
			COALESCE(custom_headers, '{}'),
			COALESCE(field_mapping, '{}'),
			COALESCE(response_mapping, '{}'),
			COALESCE(meta_config, '{}'),
			COALESCE(sitemap_config, '{}'),
			COALESCE(timeout_seconds, 10),
			COALESCE(retry_count, 1)
		FROM adapter_configs 
		WHERE product_id = $1
	`

	var customHeadersJSON []byte
	var fieldMappingJSON []byte
	var responseMappingJSON []byte
	var metaConfigJSON []byte
	var sitemapConfigJSON []byte

	err := r.db.QueryRowContext(ctx, query, productID).Scan(
		&config.EndpointPath,
		&config.HTTPMethod,
		&customHeadersJSON,
		&fieldMappingJSON,
		&responseMappingJSON,
		&metaConfigJSON,
		&sitemapConfigJSON,
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
		if len(responseMappingJSON) > 0 {
			json.Unmarshal(responseMappingJSON, &config.ResponseMapping)
		}
		if len(metaConfigJSON) > 0 {
			json.Unmarshal(metaConfigJSON, &config.MetaConfig)
		}
		if len(sitemapConfigJSON) > 0 {
			json.Unmarshal(sitemapConfigJSON, &config.SitemapConfig)
		}
	}

	return config, nil
}

func (r *AdapterConfigRepository) GetByProductID(ctx context.Context, productID string) (*product.AdapterConfig, error) {
	query := `
		SELECT id, product_id, endpoint_path, http_method,
			custom_headers, field_mapping, response_mapping, meta_config, sitemap_config, timeout_seconds,
			retry_count, created_at, updated_at
		FROM adapter_configs WHERE product_id = $1
	`

	var config product.AdapterConfig
	var customHeadersJSON []byte
	var fieldMappingJSON []byte
	var responseMappingJSON []byte
	var metaConfigJSON []byte
	var sitemapConfigJSON []byte

	err := r.db.QueryRowContext(ctx, query, productID).Scan(
		&config.ID, &config.ProductID, &config.EndpointPath,
		&config.HTTPMethod, &customHeadersJSON, &fieldMappingJSON,
		&responseMappingJSON,
		&metaConfigJSON, &sitemapConfigJSON,
		&config.TimeoutSeconds, &config.RetryCount,
		&config.CreatedAt, &config.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get adapter config: %w", err)
	}

	if len(customHeadersJSON) > 0 {
		json.Unmarshal(customHeadersJSON, &config.CustomHeaders)
	}
	if len(fieldMappingJSON) > 0 {
		json.Unmarshal(fieldMappingJSON, &config.FieldMapping)
	}
	if len(responseMappingJSON) > 0 {
		json.Unmarshal(responseMappingJSON, &config.ResponseMapping)
	}
	if len(metaConfigJSON) > 0 {
		json.Unmarshal(metaConfigJSON, &config.MetaConfig)
	}
	if len(sitemapConfigJSON) > 0 {
		json.Unmarshal(sitemapConfigJSON, &config.SitemapConfig)
	}

	return &config, nil
}

func (r *AdapterConfigRepository) InsertWithTx(ctx context.Context, tx *sql.Tx, productID string, config *product.AdapterConfig) error {
	if tx == nil {
		return fmt.Errorf("transaction is required for InsertWithTx")
	}

	customHeadersJSON, _ := json.Marshal(config.CustomHeaders)
	fieldMappingJSON, _ := json.Marshal(config.FieldMapping)
	responseMappingJSON, _ := json.Marshal(config.ResponseMapping)
	metaConfigJSON, _ := json.Marshal(config.MetaConfig)
	sitemapConfigJSON, _ := json.Marshal(config.SitemapConfig)

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
			id,
			product_id, endpoint_path, http_method,
			custom_headers, field_mapping, response_mapping, meta_config, sitemap_config, timeout_seconds, retry_count,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
		RETURNING id
	`

	err := tx.QueryRowContext(ctx, query,
		config.ID,
		productID, config.EndpointPath, config.HTTPMethod,
		customHeadersJSON, fieldMappingJSON, responseMappingJSON, metaConfigJSON, sitemapConfigJSON,
		timeoutSeconds, retryCount,
	).Scan(&config.ID)

	if err != nil {
		return fmt.Errorf("failed to insert adapter config: %w", err)
	}

	return nil
}
func (r *AdapterConfigRepository) UpdateWithTx(ctx context.Context, tx *sql.Tx, productID string, cfg product.AdapterConfig) error {
	if tx == nil {
		return fmt.Errorf("transaction is required for UpdateWithTx")
	}

	query := `
		UPDATE adapter_configs 
		SET 
			endpoint_path = $1,
			http_method = $2,
			custom_headers = $3,
			field_mapping = $4,
			meta_config = $5,
			sitemap_config = $6,
			timeout_seconds = $7,
			retry_count = $8,
			updated_at = $9
		WHERE product_id = $10 
	`

	// Marshal JSON fields
	customHeadersJSON, err := json.Marshal(cfg.CustomHeaders)
	if err != nil {
		return fmt.Errorf("failed to marshal custom_headers: %w", err)
	}

	fieldMappingJSON, err := json.Marshal(cfg.FieldMapping)
	if err != nil {
		return fmt.Errorf("failed to marshal field_mapping: %w", err)
	}

	metaConfigJSON, err := json.Marshal(cfg.MetaConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal meta_config: %w", err)
	}

	sitemapConfigJSON, err := json.Marshal(cfg.SitemapConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal sitemap_config: %w", err)
	}

	// Prepare arguments
	args := []interface{}{
		cfg.EndpointPath,
		cfg.HTTPMethod,
		customHeadersJSON,
		fieldMappingJSON,
		metaConfigJSON,
		sitemapConfigJSON,
		cfg.TimeoutSeconds,
		cfg.RetryCount,
		helper.ParseWIBTime(time.Now().Format(time.RFC3339)),
		productID,
	}

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update adapter config: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("adapter config for product %s with ID %s not found", productID, cfg.ID)
	}

	return nil
}
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

	if rowsAffected == 0 {
		return nil
	}

	return nil
}
