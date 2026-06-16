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

func (s *ProductService) CreateProduct(
	ctx context.Context,
	req product.CreateProductRequest,
	userCtx models.UserContext,
) (string, error) {

	log.Println("========== CREATE PRODUCT ==========")
	log.Printf("Request Payload: %+v\n", req)

	// Business validation
	if err := s.validateCreateRequest(req); err != nil {
		log.Printf("Validation failed: %v\n", err)
		return "", err
	}

	// ========== 1. GENERATE ALL REAL IDs ==========
	log.Println("Step 1: Generating real IDs for all entities...")

	productID := generateID()
	log.Printf("Product ID: %s\n", productID)

	// Mapping adapter config temp ID -> real ID
	adapterConfigTempToReal := make(map[string]string)

	if req.AdapterConfig != nil {
		realID := generateID()
		if req.AdapterConfig.ID != "" {
			adapterConfigTempToReal[req.AdapterConfig.ID] = realID
			log.Printf("AdapterConfig mapping: temp=%s -> real=%s\n", req.AdapterConfig.ID, realID)
		}
		req.AdapterConfig.ID = realID
		log.Printf("AdapterConfig real ID: %s\n", req.AdapterConfig.ID)
	}

	// Mapping node temp ID -> real ID
	nodeTempToReal := make(map[string]string)

	for wfIdx := range req.Workflows {
		workflowRealID := generateID()
		tempWorkflowID := req.Workflows[wfIdx].ID
		log.Printf("Workflow[%d] mapping: temp=%s -> real=%s\n", wfIdx, tempWorkflowID, workflowRealID)
		req.Workflows[wfIdx].ID = workflowRealID

		for nodeIdx := range req.Workflows[wfIdx].Nodes {
			node := &req.Workflows[wfIdx].Nodes[nodeIdx]
			nodeRealID := generateID()

			// Case 1: frontend kirim temp node ID di field ID (benar)
			if node.ID != "" {
				nodeTempToReal[node.ID] = nodeRealID
				log.Printf("Node[%d:%d] mapping from node.ID: temp=%s -> real=%s\n",
					wfIdx, nodeIdx, node.ID, nodeRealID)
			}

			// Case 2: frontend salah kirim temp node ID di field WorkflowID
			if node.WorkflowID != "" && node.WorkflowID != tempWorkflowID {
				nodeTempToReal[node.WorkflowID] = nodeRealID
				log.Printf("Node[%d:%d] mapping from node.WorkflowID: temp=%s -> real=%s\n",
					wfIdx, nodeIdx, node.WorkflowID, nodeRealID)
			}

			// Assign real IDs
			node.ID = nodeRealID
			node.WorkflowID = workflowRealID

			log.Printf("Node[%d:%d] real ID assigned: %s\n", wfIdx, nodeIdx, nodeRealID)
		}
	}

	log.Printf("Node Temp to Real Mapping: %+v\n", nodeTempToReal)
	log.Printf("AdapterConfig Temp to Real Mapping: %+v\n", adapterConfigTempToReal)

	// ========== 2. RESOLVE ALL RELATIONSHIPS ==========
	log.Println("Step 2: Resolving relationships...")

	for wfIdx := range req.Workflows {
		for nodeIdx := range req.Workflows[wfIdx].Nodes {
			node := &req.Workflows[wfIdx].Nodes[nodeIdx]

			log.Printf("Resolving node[%d:%d] real ID: %s\n", wfIdx, nodeIdx, node.ID)

			// Resolve AdapterConfigID
			if req.AdapterConfig != nil {
				node.AdapterConfigID = req.AdapterConfig.ID
				log.Printf("  AdapterConfigID set to: %s\n", node.AdapterConfigID)
			}

			// Resolve NextNodeID
			if node.NextNodeID != nil && *node.NextNodeID != "" {
				tempNextID := *node.NextNodeID
				log.Printf("  Resolving NextNodeID: '%s'\n", tempNextID)

				if realNextID, ok := nodeTempToReal[tempNextID]; ok {
					node.NextNodeID = &realNextID
					log.Printf("  ✅ NextNodeID resolved: %s -> %s\n", tempNextID, realNextID)
				} else {
					log.Printf("  ⚠️ NextNodeID '%s' not found in mapping, setting to nil\n", tempNextID)
					node.NextNodeID = nil
				}
			}

			// Resolve PreviousNodeIDs
			if len(node.PreviousNodeIDs) > 0 {
				resolvedPrevIDs := make([]string, 0, len(node.PreviousNodeIDs))
				for _, tempPrevID := range node.PreviousNodeIDs {
					log.Printf("  Resolving PreviousNodeID: '%s'\n", tempPrevID)

					if realPrevID, ok := nodeTempToReal[tempPrevID]; ok {
						resolvedPrevIDs = append(resolvedPrevIDs, realPrevID)
						log.Printf("  ✅ PreviousNodeID resolved: %s -> %s\n", tempPrevID, realPrevID)
					} else {
						log.Printf("  ⚠️ PreviousNodeID '%s' not found in mapping, skipping\n", tempPrevID)
					}
				}
				node.PreviousNodeIDs = resolvedPrevIDs
			}
		}
	}

	// ========== 3. VALIDATE NODE RELATIONSHIPS ==========
	log.Println("Step 3: Validating node relationships...")

	for wfIdx, wf := range req.Workflows {
		nodeIDSet := make(map[string]bool)
		for _, node := range wf.Nodes {
			nodeIDSet[node.ID] = true
		}

		for _, node := range wf.Nodes {
			if node.NextNodeID != nil && *node.NextNodeID != "" {
				if !nodeIDSet[*node.NextNodeID] {
					return "", fmt.Errorf("workflow[%d] node[%s]: NextNodeID %s not found in same workflow",
						wfIdx, node.ID, *node.NextNodeID)
				}
				log.Printf("  ✅ Node %s -> NextNode: %s\n", node.ID, *node.NextNodeID)
			}

			for _, prevID := range node.PreviousNodeIDs {
				if !nodeIDSet[prevID] {
					return "", fmt.Errorf("workflow[%d] node[%s]: PreviousNodeID %s not found in same workflow",
						wfIdx, node.ID, prevID)
				}
				log.Printf("  ✅ Node %s <- PrevNode: %s\n", node.ID, prevID)
			}
		}
	}

	// ========== 4. DATABASE TRANSACTION ==========
	log.Println("Step 4: Starting database transaction...")

	tx, err := s.productRepo.BeginTx(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert product
	log.Println("Inserting product...")

	repoReq := product.CreateProductRequest{
		Name:        req.Name,
		APIEndpoint: req.APIEndpoint,
		APIKey:      req.APIKey,
	}

	err = s.productRepo.InsertProductWithTx(ctx, tx, productID, repoReq)
	if err != nil {
		return "", fmt.Errorf("failed to insert product: %w", err)
	}
	log.Printf("Product inserted: %s\n", productID)

	// Insert AdapterConfig
	log.Println("Inserting adapter config...")

	if req.AdapterConfig == nil {
		return "", fmt.Errorf("adapter config is required")
	}

	cfg := product.AdapterConfig{
		ID:             req.AdapterConfig.ID,
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

	if err := s.adapterConfigRepo.InsertWithTx(ctx, tx, productID, &cfg); err != nil {
		return "", fmt.Errorf("failed to insert adapter config: %w", err)
	}
	log.Printf("AdapterConfig inserted: ID=%s, Endpoint=%s\n", cfg.ID, cfg.EndpointPath)

	// ========== INSERT WORKFLOWS & NODES ==========
	log.Println("Inserting workflows and nodes...")

	for wfIdx, wfDef := range req.Workflows {
		log.Printf("Inserting workflow[%d]: ID=%s, Name=%s\n", wfIdx, wfDef.ID, wfDef.Name)

		newWorkflow := workflow.WorkflowDefinition{
			ID:        wfDef.ID,
			ProductID: productID,
			Name:      wfDef.Name,
		}

		if err := s.workflowRepo.InsertWithTx(ctx, tx, &newWorkflow); err != nil {
			return "", fmt.Errorf("failed to insert workflow[%d]: %w", wfIdx, err)
		}
		log.Printf("Workflow[%d] inserted: %s\n", wfIdx, newWorkflow.ID)

		if len(wfDef.Nodes) == 0 {
			log.Printf("Workflow[%d] has no nodes, skipping\n", wfIdx)
			continue
		}

		// Prepare nodes
		nodes := make([]workflow_node.WorkflowNodeCreate, 0, len(wfDef.Nodes))

		for nIdx, n := range wfDef.Nodes {
			if n.AdapterConfigID == "" {
				return "", fmt.Errorf("workflow[%d] node[%d] has empty AdapterConfigID after resolution", wfIdx, nIdx)
			}
			if n.StepOrder <= 0 {
				return "", fmt.Errorf("workflow[%d] node[%d] StepOrder must be > 0", wfIdx, nIdx)
			}

			log.Printf("  Preparing node[%d:%d]: ID=%s, StepOrder=%d, AdapterConfigID=%s, NextNodeID=%v\n",
				wfIdx, nIdx, n.ID, n.StepOrder, n.AdapterConfigID, n.NextNodeID)

			nodes = append(nodes, workflow_node.WorkflowNodeCreate{
				ID:              n.ID,
				WorkflowID:      newWorkflow.ID,
				AdapterConfigID: n.AdapterConfigID,
				EndpointPath:    &n.AdapterConfig.EndpointPath,
				StepOrder:       n.StepOrder,
				InputMapping:    json.RawMessage(n.AdapterConfig.FieldMapping),
				NextNodeID:      n.NextNodeID,
				PreviousNodeIDs: n.PreviousNodeIDs,
			})
		}

		// Insert nodes: batch dengan 2 tahap (insert dulu tanpa next_node_id, update setelahnya)
		log.Printf("Inserting %d nodes for workflow[%d]...\n", len(nodes), wfIdx)

		createdNodes, err := s.workflowNodeRepo.InsertBatchWithTx(ctx, tx, nodes)
		if err != nil {
			return "", fmt.Errorf("failed to insert nodes for workflow[%d]: %w", wfIdx, err)
		}

		log.Printf("Workflow[%d]: %d nodes inserted successfully\n", wfIdx, len(createdNodes))

		for _, node := range createdNodes {
			log.Printf("  ✅ Node: ID=%s, StepOrder=%d, AdapterConfigID=%s, NextNodeID=%v, PrevIDs=%v\n",
				node.ID, node.StepOrder, node.AdapterConfigID, node.NextNodeID, node.PreviousNodeIDs)
		}
	}

	// ========== 5. COMMIT ==========
	log.Println("Step 5: Committing transaction...")

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("✅ CreateProduct SUCCESS. ProductID: %s\n", productID)
	log.Println("====================================")

	return productID, nil
}

func (s *ProductService) validateCreateRequest(req product.CreateProductRequest) error {
	if req.Name == "" {
		return fmt.Errorf("product name is required")
	}

	if len(req.Workflows) == 0 {
		return fmt.Errorf("at least one workflow is required")
	}

	for i, wf := range req.Workflows {
		if wf.Name == "" {
			return fmt.Errorf("workflow[%d] name is required", i)
		}

		if len(wf.Nodes) == 0 {
			return fmt.Errorf("workflow[%d] must have at least one node", i)
		}

		for j, node := range wf.Nodes {
			if node.StepOrder <= 0 {
				return fmt.Errorf("workflow[%d] node[%d] step_order must be greater than 0", i, j)
			}
		}
	}

	return nil
}

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
func (s *ProductService) UpdateProduct(
	ctx context.Context,
	id string,
	req product.UpdateProductRequest,
	userCtx models.UserContext,
) error {

	log.Println("========== UPDATE PRODUCT ==========")
	log.Printf("Product ID: %s\n", id)
	log.Printf("Request Payload: %+v\n", req)

	// ========== BUSINESS VALIDATION ==========
	if id == "" {
		return errors.New("product id is required")
	}

	if err := s.validateUpdateRequest(req); err != nil {
		log.Printf("Validation failed: %v\n", err)
		return err
	}

	// ========== 1. GENERATE ALL REAL IDs ==========
	log.Println("Step 1: Generating real IDs for all entities...")

	// Mapping node temp ID -> real ID
	nodeTempToReal := make(map[string]string)

	if req.Workflows != nil {
		for wfIdx := range req.Workflows {
			for nodeIdx := range req.Workflows[wfIdx].Nodes {
				node := &req.Workflows[wfIdx].Nodes[nodeIdx]

				// Skip jika node sudah punya real ID (existing node)
				if node.ID != "" && !isTempID(node.ID) {
					log.Printf("Node[%d:%d] using existing ID: %s\n", wfIdx, nodeIdx, node.ID)
					continue
				}

				nodeRealID := generateID()
				tempWorkflowID := req.Workflows[wfIdx].ID

				// Case 1: frontend kirim temp node ID di field ID
				if node.ID != "" {
					nodeTempToReal[node.ID] = nodeRealID
					log.Printf("Node[%d:%d] mapping from node.ID: temp=%s -> real=%s\n",
						wfIdx, nodeIdx, node.ID, nodeRealID)
				}

				// Case 2: frontend salah kirim temp node ID di field WorkflowID
				if node.WorkflowID != "" && node.WorkflowID != tempWorkflowID {
					nodeTempToReal[node.WorkflowID] = nodeRealID
					log.Printf("Node[%d:%d] mapping from node.WorkflowID: temp=%s -> real=%s\n",
						wfIdx, nodeIdx, node.WorkflowID, nodeRealID)
				}

				node.ID = nodeRealID
				log.Printf("Node[%d:%d] real ID assigned: %s\n", wfIdx, nodeIdx, nodeRealID)
			}
		}
	}

	log.Printf("Node Temp to Real Mapping: %+v\n", nodeTempToReal)

	// ========== 2. RESOLVE ALL RELATIONSHIPS ==========
	log.Println("Step 2: Resolving relationships...")

	if req.Workflows != nil {
		for wfIdx := range req.Workflows {
			for nodeIdx := range req.Workflows[wfIdx].Nodes {
				node := &req.Workflows[wfIdx].Nodes[nodeIdx]

				log.Printf("Resolving node[%d:%d] real ID: %s\n", wfIdx, nodeIdx, node.ID)

				// Resolve AdapterConfigID dari existing product
				if req.AdapterConfig != nil && req.AdapterConfig.ID != "" {
					node.AdapterConfigID = req.AdapterConfig.ID
					log.Printf("  AdapterConfigID set to: %s\n", node.AdapterConfigID)
				}

				// Resolve NextNodeID
				if node.NextNodeID != nil && *node.NextNodeID != "" {
					tempNextID := *node.NextNodeID
					log.Printf("  Resolving NextNodeID: '%s'\n", tempNextID)

					if realNextID, ok := nodeTempToReal[tempNextID]; ok {
						node.NextNodeID = &realNextID
						log.Printf("  ✅ NextNodeID resolved: %s -> %s\n", tempNextID, realNextID)
					} else {
						log.Printf("  ⚠️ NextNodeID '%s' not found in mapping, keeping as-is\n", tempNextID)
					}
				}

				// Resolve PreviousNodeIDs
				if len(node.PreviousNodeIDs) > 0 {
					resolvedPrevIDs := make([]string, 0, len(node.PreviousNodeIDs))
					for _, tempPrevID := range node.PreviousNodeIDs {
						log.Printf("  Resolving PreviousNodeID: '%s'\n", tempPrevID)

						if realPrevID, ok := nodeTempToReal[tempPrevID]; ok {
							resolvedPrevIDs = append(resolvedPrevIDs, realPrevID)
							log.Printf("  ✅ PreviousNodeID resolved: %s -> %s\n", tempPrevID, realPrevID)
						} else {
							// Bisa jadi sudah real ID (existing node)
							resolvedPrevIDs = append(resolvedPrevIDs, tempPrevID)
							log.Printf("  ℹ️ PreviousNodeID '%s' kept as-is (assumed real ID)\n", tempPrevID)
						}
					}
					node.PreviousNodeIDs = resolvedPrevIDs
				}
			}
		}
	}

	// ========== 3. VALIDATE NODE RELATIONSHIPS ==========
	log.Println("Step 3: Validating node relationships...")

	if req.Workflows != nil {
		for wfIdx, wf := range req.Workflows {
			nodeIDSet := make(map[string]bool)
			for _, node := range wf.Nodes {
				nodeIDSet[node.ID] = true
			}

			for _, node := range wf.Nodes {
				if node.NextNodeID != nil && *node.NextNodeID != "" {
					if !nodeIDSet[*node.NextNodeID] {
						return fmt.Errorf("workflow[%d] node[%s]: NextNodeID %s not found in same workflow",
							wfIdx, node.ID, *node.NextNodeID)
					}
					log.Printf("  ✅ Node %s -> NextNode: %s\n", node.ID, *node.NextNodeID)
				}

				for _, prevID := range node.PreviousNodeIDs {
					if !nodeIDSet[prevID] {
						return fmt.Errorf("workflow[%d] node[%s]: PreviousNodeID %s not found in same workflow",
							wfIdx, node.ID, prevID)
					}
					log.Printf("  ✅ Node %s <- PrevNode: %s\n", node.ID, prevID)
				}
			}
		}
	}

	// ========== 4. DATABASE TRANSACTION ==========
	log.Println("Step 4: Starting database transaction...")

	tx, err := s.productRepo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get existing product with lock (FOR UPDATE)
	log.Println("Fetching existing product with lock...")

	existingProduct, err := s.productRepo.GetByIDWithTx(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}
	log.Printf("Existing product fetched: %+v\n", existingProduct)

	// Business rule: cannot update active product's critical fields
	if existingProduct.Status == "active" {
		if req.Status != "" {
			return errors.New("cannot update status of active product")
		}
		if req.Platform != "" {
			return errors.New("cannot update platform of active product")
		}
	}

	// Update product
	log.Println("Updating product...")

	if err := s.productRepo.UpdateProductWithTx(ctx, tx, id, req); err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	log.Printf("Product updated: %s\n", id)

	// Update AdapterConfig jika ada
	if req.AdapterConfig != nil {
		log.Println("Updating adapter config...")

		cfg := product.AdapterConfig{
			ID:             req.AdapterConfig.ID,
			ProductID:      id,
			EndpointPath:   req.AdapterConfig.EndpointPath,
			HTTPMethod:     req.AdapterConfig.HTTPMethod,
			CustomHeaders:  req.AdapterConfig.CustomHeaders,
			FieldMapping:   req.AdapterConfig.FieldMapping,
			MetaConfig:     req.AdapterConfig.MetaConfig,
			SitemapConfig:  req.AdapterConfig.SitemapConfig,
			TimeoutSeconds: req.AdapterConfig.TimeoutSeconds,
			RetryCount:     req.AdapterConfig.RetryCount,
		}

		if err := s.adapterConfigRepo.UpdateWithTx(ctx, tx, id, &cfg); err != nil {
			return fmt.Errorf("failed to update adapter config: %w", err)
		}
		log.Printf("AdapterConfig updated: ID=%s, Endpoint=%s\n", cfg.ID, cfg.EndpointPath)
	}

	// ========== UPDATE WORKFLOWS & NODES ==========
	if req.Workflows != nil {
		log.Println("Updating workflows and nodes...")

		for wfIdx, wfDef := range req.Workflows {
			log.Printf("Processing workflow[%d]: ID=%s, Name=%s\n", wfIdx, wfDef.ID, wfDef.Name)

			updatedWorkflow := workflow.WorkflowDefinition{
				ID:        wfDef.ID,
				ProductID: id,
				Name:      wfDef.Name,
			}

			if err := s.workflowRepo.UpsertWithTx(ctx, tx, &updatedWorkflow); err != nil {
				return fmt.Errorf("failed to upsert workflow[%d]: %w", wfIdx, err)
			}
			log.Printf("Workflow[%d] upserted: %s\n", wfIdx, updatedWorkflow.ID)

			if len(wfDef.Nodes) == 0 {
				log.Printf("Workflow[%d] has no nodes, skipping\n", wfIdx)
				continue
			}

			// Prepare nodes
			nodes := make([]workflow_node.WorkflowNodeCreate, 0, len(wfDef.Nodes))

			for nIdx, n := range wfDef.Nodes {
				if n.AdapterConfigID == "" {
					return fmt.Errorf("workflow[%d] node[%d] has empty AdapterConfigID after resolution", wfIdx, nIdx)
				}
				if n.StepOrder <= 0 {
					return fmt.Errorf("workflow[%d] node[%d] StepOrder must be > 0", wfIdx, nIdx)
				}

				log.Printf("  Preparing node[%d:%d]: ID=%s, StepOrder=%d, AdapterConfigID=%s, NextNodeID=%v\n",
					wfIdx, nIdx, n.ID, n.StepOrder, n.AdapterConfigID, n.NextNodeID)

				nodes = append(nodes, workflow_node.WorkflowNodeCreate{
					ID:              n.ID,
					WorkflowID:      updatedWorkflow.ID,
					AdapterConfigID: n.AdapterConfigID,
					EndpointPath:    &n.AdapterConfig.EndpointPath,
					StepOrder:       n.StepOrder,
					InputMapping:    json.RawMessage(n.AdapterConfig.FieldMapping),
					NextNodeID:      n.NextNodeID,
					PreviousNodeIDs: n.PreviousNodeIDs,
				})
			}

			// Upsert nodes: existing node di-update, node baru di-insert
			log.Printf("Upserting %d nodes for workflow[%d]...\n", len(nodes), wfIdx)

			upsertedNodes, err := s.workflowNodeRepo.UpsertBatchWithTx(ctx, tx, nodes)
			if err != nil {
				return fmt.Errorf("failed to upsert nodes for workflow[%d]: %w", wfIdx, err)
			}

			log.Printf("Workflow[%d]: %d nodes upserted successfully\n", wfIdx, len(upsertedNodes))

			for _, node := range upsertedNodes {
				log.Printf("  ✅ Node: ID=%s, StepOrder=%d, AdapterConfigID=%s, NextNodeID=%v, PrevIDs=%v\n",
					node.ID, node.StepOrder, node.AdapterConfigID, node.NextNodeID, node.PreviousNodeIDs)
			}
		}
	}

	// ========== 5. COMMIT ==========
	log.Println("Step 5: Committing transaction...")

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("✅ UpdateProduct SUCCESS. ProductID: %s\n", id)
	log.Println("====================================")

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

func (s *ProductService) GetByID(
	ctx context.Context,
	productID string,
	userCtx models.UserContext,
) (*product.CreateProductRequest, error) {

	log.Println("========== GET PRODUCT BY ID ==========")
	log.Printf("ProductID: %s\n", productID)

	// ========== 1. FETCH PRODUCT ==========
	log.Println("Step 1: Fetching product...")

	prod, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}
	if prod == nil {
		return nil, fmt.Errorf("product not found: %s", productID)
	}
	log.Printf("Product found: ID=%s, Name=%s\n", prod.ID, prod.Name)

	// ========== 2. FETCH ADAPTER CONFIG ==========
	log.Println("Step 2: Fetching adapter config...")

	adapterCfg, err := s.adapterConfigRepo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch adapter config: %w", err)
	}
	if adapterCfg == nil {
		return nil, fmt.Errorf("adapter config not found for product: %s", productID)
	}
	log.Printf("AdapterConfig found: ID=%s, Endpoint=%s\n", adapterCfg.ID, adapterCfg.EndpointPath)

	// ========== 3. FETCH WORKFLOWS ==========
	log.Println("Step 3: Fetching workflows...")

	workflows, err := s.workflowRepo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch workflows: %w", err)
	}
	log.Printf("Found %d workflow(s)\n", len(workflows))

	// ========== 4. FETCH ALL NODES (avoid N+1) ==========
	log.Println("Step 4: Fetching all nodes...")

	workflowIDs := make([]string, 0, len(workflows))
	for _, wf := range workflows {
		workflowIDs = append(workflowIDs, wf.ID)
	}

	allNodes, err := s.workflowNodeRepo.GetByWorkflowIDs(ctx, workflowIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch nodes: %w", err)
	}
	log.Printf("Total nodes fetched: %d\n", len(allNodes))

	nodesByWorkflowID := make(map[string][]workflow_node.WorkflowNodeCreate)
	for _, n := range allNodes {
		nodesByWorkflowID[n.WorkflowID] = append(nodesByWorkflowID[n.WorkflowID], *n)
	}

	// ========== 5. ASSEMBLE WORKFLOWS + NODES ==========
	log.Println("Step 5: Assembling workflows and nodes...")

	workflowDefs := make([]workflow.WorkflowDefinition, 0, len(workflows))
	for wfIdx, wf := range workflows {
		nodes := nodesByWorkflowID[wf.ID]
		log.Printf("Workflow[%d]: ID=%s, Name=%s, Nodes=%d\n", wfIdx, wf.ID, wf.Name, len(nodes))

		nodeDefs := make([]workflow_node.WorkflowNode, 0, len(nodes))
		for nIdx, n := range nodes {
			log.Printf("  Node[%d:%d]: ID=%s, StepOrder=%d, AdapterConfigID=%s, NextNodeID=%v, PrevIDs=%v\n",
				wfIdx, nIdx, n.ID, n.StepOrder, n.AdapterConfigID, n.NextNodeID, n.PreviousNodeIDs)

			nodeDefs = append(nodeDefs, workflow_node.WorkflowNode{
				ID:              n.ID,
				WorkflowID:      n.WorkflowID,
				AdapterConfigID: n.AdapterConfigID,
				PreviousNodeIDs: n.PreviousNodeIDs,
				StepOrder:       n.StepOrder,
				NextNodeID:      n.NextNodeID,
				CreatedAt:       n.CreatedAt,
				AdapterConfig: &workflow_node.NodeAdapterConfig{
					EndpointPath: adapterCfg.EndpointPath,
					HTTPMethod:   adapterCfg.HTTPMethod,
					FieldMapping: string(n.InputMapping), // dari input_mapping node
				},
			})
		}

		wf.Nodes = nodeDefs
		workflowDefs = append(workflowDefs, *wf)
	}

	// ========== 6. ASSEMBLE RESPONSE ==========
	log.Println("Step 6: Assembling response...")

	result := &product.CreateProductRequest{
		ID:            productID,
		Name:          prod.Name,
		Platform:      prod.Platform,
		APIEndpoint:   prod.APIEndpoint,
		APIKey:        prod.APIKeyEncrypted,
		AdapterConfig: adapterCfg,
		Workflows:     workflowDefs,
	}

	log.Printf("✅ GetProductByID SUCCESS. ProductID: %s, Workflows: %d\n", productID, len(workflowDefs))
	log.Println("=======================================")

	return result, nil
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

func nullIfEmpty(id string) interface{} {
	if id == "" || id == "00000000-0000-0000-0000-000000000000" {
		return nil
	}
	return id
}
