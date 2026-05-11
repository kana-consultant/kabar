package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

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
func (s *ProductService) CreateProduct(
	ctx context.Context,
	req product.CreateProductRequest,
	userCtx models.UserContext,
) (string, error) {

	log.Println("========== CREATE PRODUCT ==========")

	// REQUEST LOG
	log.Printf("Request Payload: %+v\n", req)

	// Business validation
	log.Println("Validating request...")
	if err := s.validateCreateRequest(req); err != nil {
		log.Printf("Validation failed: %v\n", err)
		return "", err
	}
	log.Println("Validation success")

	// Begin transaction
	log.Println("Starting database transaction...")
	tx, err := s.productRepo.BeginTx(ctx)
	if err != nil {
		log.Printf("Failed begin transaction: %v\n", err)
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		log.Println("Rollback transaction (if not committed)...")
		tx.Rollback()
	}()

	TeamId := userCtx.GetTeamID()

	// Prepare repository request
	repoReq := product.CreateProductRequest{
		Name:        req.Name,
		Platform:    req.Platform,
		APIEndpoint: req.APIEndpoint,
		APIKey:      req.APIKey,
		TeamID:      TeamId,
	}

	log.Printf("Repository Request: %+v\n", repoReq)

	// Insert product via repository
	log.Println("Inserting product into database...")

	productID, err := s.productRepo.InsertProductWithTx(ctx, tx, repoReq)
	if err != nil {
		log.Printf("Failed insert product: %v\n", err)
		return "", fmt.Errorf("failed to insert product: %w", err)
	}

	log.Printf("Product inserted successfully. ProductID: %s\n", productID)

	// INSERT ADAPTER CONFIG IF PROVIDED IN PAYLOAD
	if req.AdapterConfig != nil {

		log.Println("AdapterConfig detected")

		log.Printf("AdapterConfig Payload: %+v\n", req.AdapterConfig)

		cfg := product.AdapterConfig{
			ProductID:      productID,
			EndpointPath:   req.AdapterConfig.EndpointPath,
			HTTPMethod:     req.AdapterConfig.HTTPMethod,
			CustomHeaders:  req.AdapterConfig.CustomHeaders,
			FieldMapping:   req.AdapterConfig.FieldMapping,
			TimeoutSeconds: req.AdapterConfig.TimeoutSeconds,
			RetryCount:     req.AdapterConfig.RetryCount,
		}

		log.Printf("Prepared AdapterConfig: %+v\n", cfg)

		log.Println("Inserting adapter config...")

		if err := s.adapterConfigRepo.InsertWithTx(ctx, tx, productID, &cfg); err != nil {
			log.Printf("Failed insert adapter config: %v\n", err)
			return "", fmt.Errorf("failed to insert adapter config: %w", err)
		}

		log.Println("Adapter config inserted successfully")

	} else {
		log.Println("AdapterConfig is nil, skipping insert")
	}

	// Commit transaction
	log.Println("Committing transaction...")

	if err := tx.Commit(); err != nil {
		log.Printf("Failed commit transaction: %v\n", err)
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Println("Transaction committed successfully")
	log.Printf("CreateProduct SUCCESS. ProductID: %s\n", productID)
	log.Println("====================================")

	return productID, nil
}

// setDefaultAdapterConfigValues - set default values for missing fields
func (s *ProductService) setDefaultAdapterConfigValues(config *product.AdapterConfig) {
	if config.HTTPMethod == "" {
		config.HTTPMethod = "GET"
	}
	if config.TimeoutSeconds == 0 {
		config.TimeoutSeconds = 30
	}
	if config.RetryCount == 0 {
		config.RetryCount = 3
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

func (s *ProductService) GetByID(ctx context.Context, id string) (*product.Product, error) {
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
func (s *ProductService) GetAllProducts(ctx context.Context, TeamId string) ([]product.Product, error) {
	return s.productRepo.GetProductsByTeamID(ctx, TeamId)
}

// UpdateConnectionStatus - dengan transaction optional
func (s *ProductService) UpdateConnectionStatus(
	ctx context.Context,
	productID string,
	isConnected bool,
) error {

	log.Println("========== UPDATE CONNECTION STATUS ==========")

	log.Printf("ProductID: %s\n", productID)
	log.Printf("IsConnected: %v\n", isConnected)

	if productID == "" {
		log.Println("Validation failed: product id is empty")
		return errors.New("product id is required")
	}

	// Begin transaction
	log.Println("Starting transaction...")

	tx, err := s.productRepo.BeginTx(ctx)
	if err != nil {
		log.Printf("Failed begin transaction: %v\n", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		log.Println("Rollback transaction (if not committed)...")
		tx.Rollback()
	}()

	// Get product
	log.Println("Fetching product basic info...")

	product, err := s.productRepo.GetProductBasicInfo(ctx, productID)
	if err != nil {
		log.Printf("Failed get product basic info: %v\n", err)
		return err
	}

	if product == nil {
		log.Printf("Product not found. ProductID: %s\n", productID)
		return errors.New("product not found")
	}

	log.Printf("Product found: %+v\n", product)

	// Update status
	log.Println("Updating connection status...")

	if err := s.productRepo.UpdateConnectionStatusWithTx(
		ctx,
		tx,
		productID,
		isConnected,
	); err != nil {

		log.Printf("Failed update connection status: %v\n", err)

		return fmt.Errorf(
			"failed to update connection status: %w",
			err,
		)
	}

	log.Println("Connection status updated successfully")

	// Commit transaction
	log.Println("Committing transaction...")

	if err := tx.Commit(); err != nil {
		log.Printf("Failed commit transaction: %v\n", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Println("Transaction committed successfully")
	log.Println("========== END UPDATE CONNECTION STATUS ==========")

	return nil
}

// Helper methods
func (s *ProductService) validateCreateRequest(req product.CreateProductRequest) error {
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
