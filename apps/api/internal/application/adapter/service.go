package adapter

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"seo-backend/internal/domain/adapter"
	"seo-backend/internal/domain/product"
)

// AdapterConfigService - Application Layer
type AdapterConfigService struct {
	adapterRepo adapter.AdapterConfigRepository
	db          *sql.DB
}

// NewAdapterConfigService - constructor
func NewAdapterConfigService(db *sql.DB, adapterRepo adapter.AdapterConfigRepository) adapter.AdapterConfigService {
	return &AdapterConfigService{
		db:          db,
		adapterRepo: adapterRepo,
	}
}

// GetAdapterConfig - get config with defaults
func (s *AdapterConfigService) GetAdapterConfig(ctx context.Context, productID string) (*product.AdapterConfig, error) {
	if productID == "" {
		return nil, errors.New("product id is required")
	}

	return s.adapterRepo.GetOrDefault(ctx, productID)
}

// UpdateAdapterConfig - update adapter configuration
func (s *AdapterConfigService) UpdateAdapterConfig(ctx context.Context, productID string, config *product.AdapterConfig) error {
	if productID == "" {
		return errors.New("product id is required")
	}

	if config == nil {
		return errors.New("config is required")
	}

	// Business validation
	if err := s.validateConfig(config); err != nil {
		return err
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if exists
	existingConfig, err := s.adapterRepo.GetByProductID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get existing config: %w", err)
	}

	if existingConfig == nil {
		// Insert new config
		if err := s.adapterRepo.InsertWithTx(ctx, tx, productID, config); err != nil {
			return fmt.Errorf("failed to insert adapter config: %w", err)
		}
	} else {
		// Update existing config - pass entire config object
		if err := s.adapterRepo.UpdateWithTx(ctx, tx, productID, *config); err != nil {
			return fmt.Errorf("failed to update adapter config: %w", err)
		}
	}

	return tx.Commit()
}

// CreateOrUpdateAdapterConfig - create or update full config
func (s *AdapterConfigService) CreateOrUpdateAdapterConfig(ctx context.Context, productID string, config *product.AdapterConfig) error {
	if productID == "" {
		return errors.New("product id is required")
	}

	if config == nil {
		return errors.New("config is required")
	}

	// Business validation
	if err := s.validateConfig(config); err != nil {
		return err
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if exists
	existingConfig, err := s.adapterRepo.GetByProductID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get existing config: %w", err)
	}

	if existingConfig == nil {
		// Insert new config
		if err := s.adapterRepo.InsertWithTx(ctx, tx, productID, config); err != nil {
			return fmt.Errorf("failed to insert adapter config: %w", err)
		}
	} else {
		// Update existing config - pass entire config object
		if err := s.adapterRepo.UpdateWithTx(ctx, tx, productID, *config); err != nil {
			return fmt.Errorf("failed to update adapter config: %w", err)
		}
	}

	return tx.Commit()
}

// LoadConfigForProduct - load and attach config to product
func (s *AdapterConfigService) LoadConfigForProduct(ctx context.Context, product *product.Product) error {
	if product == nil {
		return errors.New("product is required")
	}

	return s.adapterRepo.LoadForProduct(ctx, product)
}

// Helper methods
func (s *AdapterConfigService) validateConfig(config *product.AdapterConfig) error {
	if config.TimeoutSeconds < 1 || config.TimeoutSeconds > 300 {
		return errors.New("timeout seconds must be between 1 and 300")
	}

	if config.RetryCount < 0 || config.RetryCount > 10 {
		return errors.New("retry count must be between 0 and 10")
	}

	validMethods := map[string]bool{
		"GET": true, "POST": true, "PUT": true,
		"DELETE": true, "PATCH": true,
	}
	if config.HTTPMethod != "" && !validMethods[config.HTTPMethod] {
		return errors.New("invalid HTTP method")
	}

	return nil
}

func (s *AdapterConfigService) configToUpdates(config *product.AdapterConfig) (map[string]interface{}, error) {
	updates := make(map[string]interface{})

	// Basic fields
	if config.EndpointPath != "" {
		updates["endpointPath"] = config.EndpointPath
	}

	if config.HTTPMethod != "" {
		updates["httpMethod"] = config.HTTPMethod
	}

	if config.TimeoutSeconds > 0 {
		updates["timeoutSeconds"] = config.TimeoutSeconds
	}

	if config.RetryCount > 0 {
		updates["retryCount"] = config.RetryCount
	}

	// JSON fields - handle unmarshaling from string
	if config.CustomHeaders != "" {
		var headers interface{}
		if err := json.Unmarshal([]byte(config.CustomHeaders), &headers); err != nil {
			return nil, fmt.Errorf("invalid customHeaders JSON: %w", err)
		}
		updates["customHeaders"] = headers
	}

	if config.FieldMapping != "" {
		var mapping interface{}
		if err := json.Unmarshal([]byte(config.FieldMapping), &mapping); err != nil {
			return nil, fmt.Errorf("invalid fieldMapping JSON: %w", err)
		}
		updates["fieldMapping"] = mapping
	}

	if config.ResponseMapping != nil {
		updates["responseMapping"] = config.ResponseMapping
	}

	if config.MetaConfig != "" {
		var meta interface{}
		if err := json.Unmarshal([]byte(config.MetaConfig), &meta); err != nil {
			return nil, fmt.Errorf("invalid metaConfig JSON: %w", err)
		}
		updates["metaConfig"] = meta
	}

	if config.SitemapConfig != "" {
		var sitemap interface{}
		if err := json.Unmarshal([]byte(config.SitemapConfig), &sitemap); err != nil {
			return nil, fmt.Errorf("invalid sitemapConfig JSON: %w", err)
		}
		updates["sitemapConfig"] = sitemap
	}

	return updates, nil
}
