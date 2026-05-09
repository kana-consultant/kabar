package adapter

import (
	"context"
	"database/sql"
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
func (s *AdapterConfigService) UpdateAdapterConfig(ctx context.Context, productID string, updates map[string]interface{}) error {
	if productID == "" {
		return errors.New("product id is required")
	}

	if len(updates) == 0 {
		return errors.New("no updates provided")
	}

	// Business validation
	if err := s.validateUpdates(updates); err != nil {
		return err
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Update config
	if err := s.adapterRepo.UpdateWithTx(ctx, tx, productID, updates); err != nil {
		return fmt.Errorf("failed to update adapter config: %w", err)
	}

	// Commit transaction
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
		// Build updates map from config
		updates := s.configToUpdates(config)
		if err := s.adapterRepo.UpdateWithTx(ctx, tx, productID, updates); err != nil {
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
func (s *AdapterConfigService) validateUpdates(updates map[string]interface{}) error {
	// Validate timeout seconds
	if timeout, ok := updates["timeoutSeconds"]; ok {
		if t, ok := timeout.(int); ok {
			if t < 1 || t > 300 {
				return errors.New("timeout seconds must be between 1 and 300")
			}
		}
	}

	// Validate retry count
	if retry, ok := updates["retryCount"]; ok {
		if r, ok := retry.(int); ok {
			if r < 0 || r > 10 {
				return errors.New("retry count must be between 0 and 10")
			}
		}
	}

	// Validate HTTP method
	if method, ok := updates["httpMethod"]; ok {
		if m, ok := method.(string); ok {
			validMethods := map[string]bool{
				"GET": true, "POST": true, "PUT": true,
				"DELETE": true, "PATCH": true,
			}
			if !validMethods[m] {
				return errors.New("invalid HTTP method")
			}
		}
	}

	return nil
}

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
	if !validMethods[config.HTTPMethod] {
		return errors.New("invalid HTTP method")
	}

	return nil
}

func (s *AdapterConfigService) configToUpdates(config *product.AdapterConfig) map[string]interface{} {
	updates := make(map[string]interface{})

	if config.EndpointPath != "" {
		updates["endpointPath"] = config.EndpointPath
	}
	if config.HTTPMethod != "" {
		updates["httpMethod"] = config.HTTPMethod
	}
	if len(config.CustomHeaders) > 0 {
		updates["customHeaders"] = config.CustomHeaders
	}
	if len(config.FieldMapping) > 0 {
		updates["fieldMapping"] = config.FieldMapping
	}
	if config.TimeoutSeconds > 0 {
		updates["timeoutSeconds"] = config.TimeoutSeconds
	}
	if config.RetryCount > 0 {
		updates["retryCount"] = config.RetryCount
	}

	return updates
}
