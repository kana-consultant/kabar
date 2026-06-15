package product

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"seo-backend/internal/domain/adapter"
	"seo-backend/internal/domain/product"
	"seo-backend/internal/domain/workflow"
	"seo-backend/internal/domain/workflow_node"
	"seo-backend/internal/models"

	"github.com/google/uuid"
)

type ProductService struct {
	productRepo       product.ProductRepository
	adapterConfigRepo adapter.AdapterConfigRepository
	workflowRepo      workflow.WorkflowDefinitionRepository
	workflowNodeRepo  workflow_node.WorkflowNodeRepository
}

func NewProductService(db *sql.DB, productRepo product.ProductRepository,
	adapterConfigRepo adapter.AdapterConfigRepository,
	WorkflowDefinitionRepository workflow.WorkflowDefinitionRepository,
	WorkflowNodeRepository workflow_node.WorkflowNodeRepository) product.ProductService {
	return &ProductService{
		productRepo:       productRepo,
		adapterConfigRepo: adapterConfigRepo,
		workflowRepo:      WorkflowDefinitionRepository,
		workflowNodeRepo:  WorkflowNodeRepository,
	}
}

// CreateProduct - Application mengelola transaction
func (s *ProductService) CreateProduct(
	ctx context.Context,
	req product.CreateProductRequest,
	userCtx models.UserContext,
) (string, error) {

	log.Println("========== CREATE PRODUCT ==========")
	log.Printf("Request Payload: %+v\n", req)

	// Business validation
	log.Println("Validating request...")
	if err := s.validateCreateRequest(req); err != nil {
		log.Printf("Validation failed: %v\n", err)
		return "", err
	}
	log.Println("Validation success")

	// Generate IDs untuk semua entity
	log.Println("Generating IDs...")

	productID := generateID()
	log.Printf("Generated ProductID: %s\n", productID)

	// Generate IDs untuk adapter configs
	adapterConfigTempIDs := make(map[string]string)     // tempID (from request) -> real ID
	adapterConfigEndpointMap := make(map[string]string) // endpoint -> real ID

	// Single adapter config
	if req.AdapterConfig != nil {
		realID := generateID()
		if req.AdapterConfig.ID != "" {
			adapterConfigTempIDs[req.AdapterConfig.ID] = realID
		}
		adapterConfigEndpointMap[req.AdapterConfig.EndpointPath] = realID
		req.AdapterConfig.ID = realID // set real ID ke request
		log.Printf("Generated AdapterConfig ID: %s for endpoint: %s\n", realID, req.AdapterConfig.EndpointPath)
	}

	// Batch adapter configs
	for i := range req.AdapterConfigs {
		realID := generateID()
		if req.AdapterConfigs[i].ID != "" {
			adapterConfigTempIDs[req.AdapterConfigs[i].ID] = realID
		}
		adapterConfigEndpointMap[req.AdapterConfigs[i].EndpointPath] = realID
		req.AdapterConfigs[i].ID = realID // set real ID ke request
		log.Printf("Generated AdapterConfig ID: %s for endpoint: %s\n", realID, req.AdapterConfigs[i].EndpointPath)
	}

	// Generate IDs untuk workflows dan nodes
	workflowTempIDs := make(map[string]string) // tempID -> real ID
	for wfIdx := range req.Workflows {
		workflowRealID := generateID()
		if req.Workflows[wfIdx].ID != "" {
			workflowTempIDs[req.Workflows[wfIdx].ID] = workflowRealID
		}
		req.Workflows[wfIdx].ID = workflowRealID
		log.Printf("Generated Workflow ID: %s for workflow: %s\n", workflowRealID, req.Workflows[wfIdx].Name)

		// Generate IDs untuk nodes dalam workflow ini
		for nodeIdx := range req.Workflows[wfIdx].Nodes {
			nodeRealID := generateID()
			if req.Workflows[wfIdx].Nodes[nodeIdx].ID != "" {
				// Simpan mapping temp ID ke real ID untuk resolve next_node_id nanti
				// Bisa disimpan di context atau map terpisah
			}
			req.Workflows[wfIdx].Nodes[nodeIdx].ID = nodeRealID
			log.Printf("Generated Node ID: %s for workflow %s step %d\n",
				nodeRealID, workflowRealID, req.Workflows[wfIdx].Nodes[nodeIdx].StepOrder)
		}
	}

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

	// Insert product
	repoReq := product.CreateProductRequest{
		Name:        req.Name,
		APIEndpoint: req.APIEndpoint,
		APIKey:      req.APIKey,
	}

	log.Printf("Repository Request: %+v\n", repoReq)
	log.Println("Inserting product into database...")

	// Gunakan productID yang sudah digenerate
	err = s.productRepo.InsertProductWithTx(ctx, tx, productID, repoReq)
	if err != nil {
		log.Printf("Failed insert product: %v\n", err)
		return "", fmt.Errorf("failed to insert product: %w", err)
	}

	log.Printf("Product inserted successfully. ProductID: %s\n", productID)

	// ─── INSERT ADAPTER CONFIGS ───────────────────────────────────────────────
	// Map: EndpointPath → real AdapterConfig ID (untuk resolve workflow nodes)
	adapterConfigIDMap := make(map[string]string)

	// Single adapter config (legacy field)
	if req.AdapterConfig != nil {
		log.Println("Single AdapterConfig detected")
		log.Printf("AdapterConfig Payload: %+v\n", req.AdapterConfig)

		cfg := product.AdapterConfig{
			ID:             req.AdapterConfig.ID, // Gunakan ID yang sudah digenerate
			ProductID:      productID,
			EndpointPath:   req.AdapterConfig.EndpointPath,
			HTTPMethod:     req.AdapterConfig.HTTPMethod,
			CustomHeaders:  req.AdapterConfig.CustomHeaders,
			FieldMapping:   req.AdapterConfig.FieldMapping,
			MetaConfig:     req.AdapterConfig.MetaConfig,
			SitemapConfig:  req.AdapterConfig.SitemapConfig,
			TimeoutSeconds: req.AdapterConfig.TimeoutSeconds,
			RetryCount:     req.AdapterConfig.RetryCount,
		}

		log.Printf("Prepared AdapterConfig: %+v\n", cfg)
		log.Println("Inserting single adapter config...")

		if err := s.adapterConfigRepo.InsertWithTx(ctx, tx, cfg.ProductID, &cfg); err != nil {
			log.Printf("Failed insert adapter config: %v\n", err)
			return "", fmt.Errorf("failed to insert adapter config: %w", err)
		}

		log.Printf("Single adapter config inserted. ID: %s\n", cfg.ID)
		adapterConfigIDMap[cfg.EndpointPath] = cfg.ID
		if req.AdapterConfig.ID != "" {
			adapterConfigIDMap[req.AdapterConfig.ID] = cfg.ID
		}

	} else {
		log.Println("Single AdapterConfig is nil, skipping")
	}

	// Batch adapter configs
	if len(req.AdapterConfigs) > 0 {
		log.Printf("Batch AdapterConfigs detected: %d configs\n", len(req.AdapterConfigs))

		for i, ac := range req.AdapterConfigs {
			cfg := product.AdapterConfig{
				ID:             ac.ID, // Gunakan ID yang sudah digenerate
				ProductID:      productID,
				EndpointPath:   ac.EndpointPath,
				HTTPMethod:     ac.HTTPMethod,
				CustomHeaders:  ac.CustomHeaders,
				FieldMapping:   ac.FieldMapping,
				MetaConfig:     ac.MetaConfig,
				SitemapConfig:  ac.SitemapConfig,
				TimeoutSeconds: ac.TimeoutSeconds,
				RetryCount:     ac.RetryCount,
			}

			log.Printf("Inserting AdapterConfig[%d]: %+v\n", i, cfg)

			if err := s.adapterConfigRepo.InsertWithTx(ctx, tx, cfg.ProductID, &cfg); err != nil {
				log.Printf("Failed insert adapter config[%d]: %v\n", i, err)
				return "", fmt.Errorf("failed to insert adapter config[%d]: %w", i, err)
			}

			log.Printf("AdapterConfig[%d] inserted. ID: %s, EndpointPath: %s\n", i, cfg.ID, cfg.EndpointPath)

			// Map original ID (dari payload) dan EndpointPath → real ID
			if ac.ID != "" {
				adapterConfigIDMap[ac.ID] = cfg.ID
			}
			adapterConfigIDMap[cfg.EndpointPath] = cfg.ID
		}
	} else {
		log.Println("No batch AdapterConfigs, skipping")
	}

	log.Printf("AdapterConfig ID Map: %+v\n", adapterConfigIDMap)

	// ─── INSERT WORKFLOWS & NODES ─────────────────────────────────────────────
	if len(req.Workflows) > 0 {
		log.Printf("Workflows detected: %d workflows\n", len(req.Workflows))

		for wfIdx, wfDef := range req.Workflows {
			log.Printf("Processing Workflow[%d]: ID=%s, Name=%s\n", wfIdx, wfDef.ID, wfDef.Name)

			// Insert workflow definition dengan ID yang sudah digenerate
			newWorkflow := workflow.WorkflowDefinition{
				ID:        wfDef.ID, // Gunakan ID yang sudah digenerate
				ProductID: productID,
				Name:      wfDef.Name,
			}

			err := s.workflowRepo.InsertWithTx(ctx, tx, &newWorkflow)
			if err != nil {
				log.Printf("Failed insert workflow[%d]: %v\n", wfIdx, err)
				return "", fmt.Errorf("failed to insert workflow[%d] %q: %w", wfIdx, wfDef.Name, err)
			}

			log.Printf("Workflow[%d] inserted. WorkflowID: %s\n", wfIdx, newWorkflow.ID)

			if len(wfDef.Nodes) == 0 {
				log.Printf("Workflow[%d] has no nodes, skipping node insert\n", wfIdx)
				continue
			}

			// ── Prepare nodes ──────────────────────────────────────────────
			// Map: tempID (dari payload) → real ID (sudah digenerate)
			tempIDToRealID := make(map[string]string)
			nodes := make([]workflow_node.WorkflowNode, 0, len(wfDef.Nodes))

			for nIdx, n := range wfDef.Nodes {
				log.Printf("Preparing Node[%d]: ID=%s, AdapterConfigID=%s, StepOrder=%d\n",
					nIdx, n.ID, n.AdapterConfigID, n.StepOrder)

				// Resolve AdapterConfigID: payload ID → real DB ID
				realAdapterCfgID := n.AdapterConfigID
				if mapped, ok := adapterConfigIDMap[n.AdapterConfigID]; ok {
					realAdapterCfgID = mapped
					log.Printf("Node[%d] AdapterConfigID resolved: %s → %s\n", nIdx, n.AdapterConfigID, realAdapterCfgID)
				}

				// Validasi wajib
				if realAdapterCfgID == "" {
					log.Printf("Node[%d] skipped: AdapterConfigID is empty\n", nIdx)
					continue
				}
				if n.StepOrder <= 0 {
					log.Printf("Node[%d] skipped: StepOrder must be > 0 (got %d)\n", nIdx, n.StepOrder)
					continue
				}

				inputMappingJSON, err := json.Marshal(n.InputMapping)
				if err != nil {
					log.Printf("Node[%d] failed to marshal InputMapping: %v\n", nIdx, err)
					return "", fmt.Errorf("workflow[%d] node[%d]: failed to marshal input mapping: %w", wfIdx, nIdx, err)
				}

				// Simpan mapping temp ID ke real ID (node ID sudah digenerate)
				if n.ID != "" {
					tempIDToRealID[n.ID] = n.ID // ID sudah real karena digenerate di awal
				}

				// Resolve NextNodeID jika menggunakan temp ID
				var nextNodeIDPtr *string
				if n.NextNodeID != nil && *n.NextNodeID != "" {
					// Cek apakah NextNodeID adalah temp ID yang perlu di-resolve
					// Karena semua ID sudah digenerate, NextNodeID harusnya sudah real
					nextNodeIDPtr = n.NextNodeID
				}

				nodes = append(nodes, workflow_node.WorkflowNode{
					ID:              n.ID, // Gunakan ID yang sudah digenerate
					WorkflowID:      newWorkflow.ID,
					AdapterConfigID: realAdapterCfgID,
					StepOrder:       n.StepOrder,
					InputMapping:    inputMappingJSON,
					NextNodeID:      nextNodeIDPtr,
					PreviousNodeIDs: n.PreviousNodeIDs,
				})
			}

			if len(nodes) == 0 {
				log.Printf("Workflow[%d] has no valid nodes after filtering\n", wfIdx)
				continue
			}

			// ── Insert nodes (batch, fallback ke satu-satu) ────────────────
			log.Printf("Inserting %d nodes for Workflow[%d]...\n", len(nodes), wfIdx)

			createdNodes, err := s.workflowNodeRepo.InsertBatchWithTx(ctx, tx, nodes)
			if err != nil {
				log.Printf("Batch insert failed for Workflow[%d], falling back to single insert: %v\n", wfIdx, err)

				createdNodes = make([]workflow_node.WorkflowNode, 0, len(nodes))
				for nIdx, node := range nodes {
					nodeCopy := node
					if err := s.workflowNodeRepo.InsertWithTx(ctx, tx, &nodeCopy); err != nil {
						log.Printf("Failed insert single node[%d] for Workflow[%d]: %v\n", nIdx, wfIdx, err)
						return "", fmt.Errorf("workflow[%d] node[%d]: failed to insert: %w", wfIdx, nIdx, err)
					}
					createdNodes = append(createdNodes, nodeCopy)
				}
			}

			log.Printf("Workflow[%d]: %d nodes inserted\n", wfIdx, len(createdNodes))

			// ── Resolve nextNodeID jika masih menggunakan temp ID ───────────
			for _, created := range createdNodes {
				if created.NextNodeID == nil {
					continue
				}

				nextID := *created.NextNodeID
				// Cek apakah nextID adalah temp ID yang perlu di-resolve
				if realID, ok := tempIDToRealID[nextID]; ok && realID != nextID {
					log.Printf("Updating NextNodeID for node %s: %s → %s\n", created.ID, nextID, realID)

					if err := s.workflowNodeRepo.UpdateWithTx(ctx, tx, created.ID, map[string]interface{}{
						"next_node_id": realID,
					}); err != nil {
						log.Printf("Failed update nextNodeID for node %s: %v\n", created.ID, err)
						return "", fmt.Errorf("workflow[%d]: failed to update nextNodeID for node %s: %w", wfIdx, created.ID, err)
					}
				}
			}
		}
	} else {
		log.Println("No Workflows in payload, skipping")
	}

	// ─── COMMIT ───────────────────────────────────────────────────────────────
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

// Helper function untuk generate ID
func generateID() string {
	return uuid.New().String()
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
		log.Printf("Warning: failed to delete adapter config: %v\n", err)
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
func (s *ProductService) GetAllProducts(ctx context.Context, userCtx models.UserContext) ([]product.Product, error) {
	return s.productRepo.GetProductsByTeamID(ctx, userCtx)
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
