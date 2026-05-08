package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"seo-backend/internal/domain/adapter"
	"seo-backend/internal/domain/product"
	"seo-backend/internal/models"
)

type ProductService struct {
	productRepo       product.ProductRepository
	adapterConfigRepo adapter.AdapterConfigRepository
}

func NewProductService(db *sql.DB, productRepo product.ProductRepository, adapterConfigRepo adapter.AdapterConfigRepository) product.ProductService {
	return &ProductService{
		productRepo:       productRepo,
		adapterConfigRepo: adapterConfigRepo,
	}
}

// CreateProduct - Application mengelola transaction
// CreateProduct - Application mengelola transaction
func (s *ProductService) CreateProduct(ctx context.Context, req models.CreateProductRequest, userCtx models.UserContext) (string, error) {
	// Business validation
	if err := s.validateCreateRequest(req); err != nil {
		return "", err
	}

	// Begin transaction
	tx, err := s.productRepo.BeginTx(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Prepare repository request
	repoReq := product.CreateProductRequest{
		Name:        req.Name,
		Platform:    req.Platform,
		APIEndpoint: req.APIEndpoint,
		APIKey:      req.APIKey,
		Status:      "pending",
		SyncStatus:  "idle",
		CreatedBy:   nullIfEmpty(userCtx.GetUserID()),
		TeamID:      nullIfEmpty(userCtx.GetTeamID()),
		UserID:      nullIfEmpty(userCtx.GetUserID()),
	}

	// Insert product via repository
	productID, err := s.productRepo.InsertProductWithTx(ctx, tx, repoReq)
	if err != nil {
		return "", fmt.Errorf("failed to insert product: %w", err)
	}

	// INSERT ADAPTER CONFIG IF PROVIDED IN PAYLOAD

	if req.AdapterConfig != nil {
		cfg := models.AdapterConfig{
			ProductID:      productID,
			EndpointPath:   req.AdapterConfig.EndpointPath,
			HTTPMethod:     req.AdapterConfig.HTTPMethod,
			CustomHeaders:  req.AdapterConfig.CustomHeaders,
			FieldMapping:   req.AdapterConfig.FieldMapping,
			TimeoutSeconds: req.AdapterConfig.TimeoutSeconds,
			RetryCount:     req.AdapterConfig.RetryCount,
		}

		if err := s.adapterConfigRepo.InsertWithTx(ctx, tx, productID, &cfg); err != nil {
			return "", fmt.Errorf("failed to insert adapter config: %w", err)
		}
	}

	// Commit transaction (both product and adapter config if exists)
	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return productID, nil
}

// setDefaultAdapterConfigValues - set default values for missing fields
func (s *ProductService) setDefaultAdapterConfigValues(config *models.AdapterConfig) {
	if config.HTTPMethod == "" {
		config.HTTPMethod = "GET"
	}
	if config.TimeoutSeconds == 0 {
		config.TimeoutSeconds = 30
	}
	if config.RetryCount == 0 {
		config.RetryCount = 3
	}
	if config.CustomHeaders == nil {
		config.CustomHeaders = make(map[string]string)
	}
	if config.FieldMapping == nil {
		config.FieldMapping = make(map[string]string)
	}
}

// UpdateProduct - Application mengelola transaction
func (s *ProductService) UpdateProduct(ctx context.Context, id string, updates map[string]interface{}, userCtx models.UserContext) error {
	// Business validation
	if id == "" {
		return errors.New("product id is required")
	}
	if len(updates) == 0 {
		return errors.New("no updates provided")
	}

	// Begin transaction
	tx, err := s.productRepo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get existing product with lock (FOR UPDATE)
	existingProduct, err := s.productRepo.GetByIDWithTx(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	// Business rule: cannot update active product's critical fields
	if existingProduct.Status == "active" {
		if _, ok := updates["status"]; ok {
			return errors.New("cannot update status of active product")
		}
		if _, ok := updates["platform"]; ok {
			return errors.New("cannot update platform of active product")
		}
	}

	// UPDATE PRODUCT VIA REPOSITORY

	if err := s.productRepo.UpdateProductWithTx(ctx, tx, id, updates); err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	// UPDATE ADAPTER CONFIG IF PROVIDED IN PAYLOAD

	if adapterUpdates, ok := updates["adapterConfig"]; ok {
		if adapterUpdatesMap, ok := adapterUpdates.(map[string]interface{}); ok && len(adapterUpdatesMap) > 0 {
			if err := s.adapterConfigRepo.UpdateWithTx(ctx, tx, id, adapterUpdatesMap); err != nil {
				return fmt.Errorf("failed to update adapter config: %w", err)
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// DeleteProduct - Application mengelola transaction
// DeleteProduct - Standard pattern (mirip Create)
func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	// Business validation
	if id == "" {
		return errors.New("product id is required")
	}

	// Begin transaction
	tx, err := s.productRepo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get existing product with lock
	existingProduct, err := s.productRepo.GetByIDWithTx(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	// Business rule
	if existingProduct.Status == "active" {
		return errors.New("cannot delete active product, deactivate it first")
	}

	// Delete adapter config first (if exists)
	if err := s.adapterConfigRepo.DeleteWithTx(ctx, tx, id); err != nil {
		// Log but don't fail - config might not exist
		// Continue with product deletion
	}

	// Delete product
	if err := s.productRepo.DeleteProductWithTx(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *ProductService) GetByID(ctx context.Context, id string) (*models.Product, error) {
	tx, err := s.productRepo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. ambil product
	product, err := s.productRepo.GetByIDWithTx(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	// 2. ambil adapter config
	adapterConfig, err := s.adapterConfigRepo.GetByProductID(ctx, id)
	if err != nil {
		return nil, err
	}
	// 3. inject ke product
	product.AdapterConfig = adapterConfig

	return product, nil
}

// GetProduct - tanpa transaction (read-only)
func (s *ProductService) GetAllProducts(ctx context.Context) ([]models.Product, int, error) {
	return s.productRepo.GetAllWithFilters(ctx)
}

// UpdateConnectionStatus - dengan transaction optional
func (s *ProductService) UpdateConnectionStatus(ctx context.Context, productID string, isConnected bool) error {
	if productID == "" {
		return errors.New("product id is required")
	}

	// Begin transaction
	tx, err := s.productRepo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get product with lock - wajib pake tx
	product, err := s.productRepo.GetProductBasicInfo(ctx, productID)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("product not found")
	}

	-, err := s.adapterConfigRepo.GetByProductID(ctx, productID)

	// Update status - pake tx
	if err := s.productRepo.UpdateConnectionStatusWithTx(ctx, tx, productID, isConnected); err != nil {
		return fmt.Errorf("failed to update connection status: %w", err)
	}

	return tx.Commit()
}

// Helper methods
func (s *ProductService) validateCreateRequest(req models.CreateProductRequest) error {
	if req.Name == "" {
		return errors.New("product name is required")
	}
	if req.Platform == "" {
		return errors.New("platform is required")
	}
	if req.APIEndpoint == "" {
		return errors.New("API endpoint is required")
	}
	if req.APIKey == "" {
		return errors.New("API key is required")
	}
	return nil
}

func nullIfEmpty(id string) interface{} {
	if id == "" || id == "00000000-0000-0000-0000-000000000000" {
		return nil
	}
	return id
}
