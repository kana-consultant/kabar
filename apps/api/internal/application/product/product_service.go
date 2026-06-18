package product

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"

	"seo-backend/internal/domain/adapter"
	"seo-backend/internal/domain/draft"
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
	req product.ProductRequest,
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

			// Resolve NextNodeIDs
			if len(node.NextNodeIDs) > 0 {
				resolvedNextIDs := make([]string, 0, len(node.NextNodeIDs))
				for _, tempNextID := range node.NextNodeIDs {
					log.Printf("  Resolving NextNodeID: '%s'\n", tempNextID)

					if realNextID, ok := nodeTempToReal[tempNextID]; ok {
						resolvedNextIDs = append(resolvedNextIDs, realNextID)
						log.Printf("  ✅ NextNodeID resolved: %s -> %s\n", tempNextID, realNextID)
					} else {
						log.Printf("  ⚠️ NextNodeID '%s' not found in mapping, skipping\n", tempNextID)
					}
				}
				node.NextNodeIDs = resolvedNextIDs
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
			for _, nextID := range node.NextNodeIDs {
				if !nodeIDSet[nextID] {
					return "", fmt.Errorf("workflow[%d] node[%s]: NextNodeID %s not found in same workflow",
						wfIdx, node.ID, nextID)
				}
				log.Printf("  ✅ Node %s -> NextNode: %s\n", node.ID, nextID)
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

	repoReq := product.ProductRequest{
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

			log.Printf("  Preparing node[%d:%d]: ID=%s, StepOrder=%d, AdapterConfigID=%s, NextNodeIDs=%v\n",
				wfIdx, nIdx, n.ID, n.StepOrder, n.AdapterConfigID, n.NextNodeIDs)

			nodes = append(nodes, workflow_node.WorkflowNodeCreate{
				ID:              n.ID,
				WorkflowID:      newWorkflow.ID,
				AdapterConfigID: n.AdapterConfigID,
				EndpointPath:    &n.AdapterConfig.EndpointPath,
				StepOrder:       n.StepOrder,
				InputMapping:    json.RawMessage(n.AdapterConfig.FieldMapping),
				NextNodeIDs:     n.NextNodeIDs,
				PreviousNodeIDs: n.PreviousNodeIDs,
				HTTPMethod:      n.AdapterConfig.HTTPMethod,
			})
		}

		// Insert nodes: batch dengan 2 tahap (insert dulu tanpa next_node_ids, update setelahnya)
		log.Printf("Inserting %d nodes for workflow[%d]...\n", len(nodes), wfIdx)

		createdNodes, err := s.workflowNodeRepo.InsertBatchWithTx(ctx, tx, nodes)
		if err != nil {
			return "", fmt.Errorf("failed to insert nodes for workflow[%d]: %w", wfIdx, err)
		}

		log.Printf("Workflow[%d]: %d nodes inserted successfully\n", wfIdx, len(createdNodes))

		for _, node := range createdNodes {
			log.Printf("  ✅ Node: ID=%s, StepOrder=%d, AdapterConfigID=%s, NextNodeIDs=%v, PrevIDs=%v\n",
				node.ID, node.StepOrder, node.AdapterConfigID, node.NextNodeIDs, node.PreviousNodeIDs)
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

// UpdateProduct - Application mengelola transaction
func (s *ProductService) UpdateProduct(
	ctx context.Context,
	id string,
	req product.ProductRequest,
	userCtx models.UserContext,
) error {

	log.Println("========== UPDATE PRODUCT ==========")
	log.Printf("Product ID: %s\n", id)
	log.Printf("Request Payload: %+v\n", req)

	if id == "" {
		return errors.New("product id is required")
	}

	// ========== 1. GENERATE ALL REAL IDs ==========
	log.Println("Step 1: Generating real IDs for all entities...")

	if req.AdapterConfig != nil {
		if req.AdapterConfig.ID != "" && !isTempID(req.AdapterConfig.ID) {
			// Frontend kirim real ID yang valid — pakai langsung
			log.Printf("AdapterConfig ID from request (real): %s\n", req.AdapterConfig.ID)
		} else {
			// Kosong atau temp — tandai, resolve dari DB nanti
			log.Printf("AdapterConfig ID is empty/temp (%s), will resolve from DB\n", req.AdapterConfig.ID)
			req.AdapterConfig.ID = "" // kosongkan agar tahu harus fetch
		}
	}
	nodeTempToReal := make(map[string]string)

	if req.Workflows != nil {
		for wfIdx := range req.Workflows {
			// [FIX Gap 1] Generate real workflow ID jika temp atau kosong
			if req.Workflows[wfIdx].ID == "" || isTempID(req.Workflows[wfIdx].ID) {
				realWfID := generateID()
				log.Printf("Workflow[%d] mapping: temp=%s -> real=%s\n", wfIdx, req.Workflows[wfIdx].ID, realWfID)
				req.Workflows[wfIdx].ID = realWfID
			}
			log.Printf("Workflow[%d] ID: %s\n", wfIdx, req.Workflows[wfIdx].ID)

			for nodeIdx := range req.Workflows[wfIdx].Nodes {
				node := &req.Workflows[wfIdx].Nodes[nodeIdx]

				// Skip jika node sudah punya real ID (existing node)
				if node.ID != "" && !isTempID(node.ID) {
					log.Printf("Node[%d:%d] using existing ID: %s\n", wfIdx, nodeIdx, node.ID)
					continue
				}

				nodeRealID := generateID()
				tempWorkflowID := req.Workflows[wfIdx].ID // sudah real di sini

				// Case 1: frontend kirim temp node ID di field ID
				if node.ID != "" {
					nodeTempToReal[node.ID] = nodeRealID
					log.Printf("Node[%d:%d] mapping from node.ID: temp=%s -> real=%s\n",
						wfIdx, nodeIdx, node.ID, nodeRealID)
				}

				// [FIX Gap 5] Case 2: frontend salah kirim temp node ID di field WorkflowID
				if node.WorkflowID != "" && node.WorkflowID != tempWorkflowID && isTempID(node.WorkflowID) {
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

				// [FIX Gap 3] Set WorkflowID ke real workflow ID, konsisten dengan CreateProduct
				node.WorkflowID = req.Workflows[wfIdx].ID
				log.Printf("  WorkflowID set to: %s\n", node.WorkflowID)

				// Resolve AdapterConfigID dari existing product
				if req.AdapterConfig != nil && req.AdapterConfig.ID != "" {
					node.AdapterConfigID = req.AdapterConfig.ID
					log.Printf("  AdapterConfigID set to: %s\n", node.AdapterConfigID)
				}

				// Resolve NextNodeIDs
				if len(node.NextNodeIDs) > 0 {
					resolvedNextIDs := make([]string, 0, len(node.NextNodeIDs))
					for _, tempNextID := range node.NextNodeIDs {
						log.Printf("  Resolving NextNodeID: '%s'\n", tempNextID)
						if realNextID, ok := nodeTempToReal[tempNextID]; ok {
							resolvedNextIDs = append(resolvedNextIDs, realNextID)
							log.Printf("  ✅ NextNodeID resolved: %s -> %s\n", tempNextID, realNextID)
						} else {
							// Bisa jadi sudah real ID (existing node)
							resolvedNextIDs = append(resolvedNextIDs, tempNextID)
							log.Printf("  ℹ️ NextNodeID '%s' kept as-is (assumed real ID)\n", tempNextID)
						}
					}
					node.NextNodeIDs = resolvedNextIDs
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
	// [FIX Gap 4] Validasi strict: semua next/prev ID harus ada di workflow yang sama
	log.Println("Step 3: Validating node relationships...")

	if req.Workflows != nil {
		for wfIdx, wf := range req.Workflows {
			nodeIDSet := make(map[string]bool)
			for _, node := range wf.Nodes {
				nodeIDSet[node.ID] = true
			}

			for _, node := range wf.Nodes {
				for _, nextID := range node.NextNodeIDs {
					if !nodeIDSet[nextID] {
						// [FIX Gap 4] Strict seperti CreateProduct — return error, bukan skip
						return fmt.Errorf("workflow[%d] node[%s]: NextNodeID %s not found in same workflow",
							wfIdx, node.ID, nextID)
					}
					log.Printf("  ✅ Node %s -> NextNode: %s\n", node.ID, nextID)
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

	log.Println("Fetching existing product with lock...")
	existingProduct, err := s.productRepo.GetByIDWithTx(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}
	log.Printf("Existing product fetched: %+v\n", existingProduct)

	if existingProduct.Status == "active" {
		if req.Platform != "" {
			return errors.New("cannot update platform of active product")
		}
	}

	// [FIX] Resolve AdapterConfig ID dari DB jika tidak ada di request
	if req.AdapterConfig != nil && req.AdapterConfig.ID == "" {
		existingCfg, err := s.adapterConfigRepo.GetByProductID(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to get existing adapter config: %w", err)
		}
		req.AdapterConfig.ID = existingCfg.ID
		log.Printf("AdapterConfig ID resolved from DB: %s\n", req.AdapterConfig.ID)
	}

	// [FIX] Backfill AdapterConfigID ke semua node setelah ID diketahui
	if req.AdapterConfig != nil && req.Workflows != nil {
		for wfIdx := range req.Workflows {
			for nodeIdx := range req.Workflows[wfIdx].Nodes {
				node := &req.Workflows[wfIdx].Nodes[nodeIdx]
				if node.AdapterConfigID == "" {
					node.AdapterConfigID = req.AdapterConfig.ID
					log.Printf("Node[%d:%d] AdapterConfigID backfilled: %s\n", wfIdx, nodeIdx, node.AdapterConfigID)
				}
			}
		}
	}

	log.Println("Updating product...")
	if err := s.productRepo.UpdateProductWithTx(ctx, tx, id, req); err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	log.Printf("Product updated: %s\n", id)

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
		if err := s.adapterConfigRepo.UpdateWithTx(ctx, tx, id, cfg); err != nil {
			return fmt.Errorf("failed to update adapter config: %w", err)
		}
		log.Printf("AdapterConfig updated: ID=%s, Endpoint=%s\n", cfg.ID, cfg.EndpointPath)
	}

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

			nodes := make([]workflow_node.WorkflowNode, 0, len(wfDef.Nodes))
			for nIdx, n := range wfDef.Nodes {

				if n.StepOrder <= 0 {
					return fmt.Errorf("workflow[%d] node[%d] StepOrder must be > 0", wfIdx, nIdx)
				}

				log.Printf("  Preparing node[%d:%d]: ID=%s, WorkflowID=%s, StepOrder=%d, AdapterConfigID=%s, NextNodeIDs=%v\n",
					wfIdx, nIdx, n.ID, n.WorkflowID, n.StepOrder, req.AdapterConfig.ID, n.NextNodeIDs)

				nodes = append(nodes, workflow_node.WorkflowNode{
					ID:              n.ID,
					WorkflowID:      updatedWorkflow.ID,
					AdapterConfigID: req.AdapterConfig.ID,
					StepOrder:       n.StepOrder,
					AdapterConfig:   n.AdapterConfig,
					NextNodeIDs:     n.NextNodeIDs,
					PreviousNodeIDs: n.PreviousNodeIDs,
				})
			}

			log.Printf("Upserting %d nodes for workflow[%d]...\n", len(nodes), wfIdx)
			if err := s.workflowNodeRepo.UpsertWithTx(ctx, tx, nodes); err != nil {
				return fmt.Errorf("failed to upsert nodes for workflow[%d]: %w", wfIdx, err)
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

func (s *ProductService) validateCreateRequest(req product.ProductRequest) error {
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

// Helper untuk deteksi temp ID dari frontend
func isTempID(id string) bool {
	// Sesuaikan dengan format temp ID frontend kamu
	// Contoh: return strings.HasPrefix(id, "temp-")
	return strings.HasPrefix(id, "temp-")
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
) (*product.ProductRequest, error) {

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
				wfIdx, nIdx, n.ID, n.StepOrder, n.AdapterConfigID, n.NextNodeIDs, n.PreviousNodeIDs)

			nodeDefs = append(nodeDefs, workflow_node.WorkflowNode{
				ID:              n.ID,
				WorkflowID:      n.WorkflowID,
				AdapterConfigID: n.AdapterConfigID,
				PreviousNodeIDs: n.PreviousNodeIDs,
				StepOrder:       n.StepOrder,
				NextNodeIDs:     n.NextNodeIDs,
				CreatedAt:       n.CreatedAt,
				AdapterConfig: &workflow_node.NodeAdapterConfig{
					EndpointPath: *n.EndpointPath,
					HTTPMethod:   n.HTTPMethod,
					FieldMapping: string(n.InputMapping), // dari input_mapping node
				},
			})
		}

		wf.Nodes = nodeDefs
		workflowDefs = append(workflowDefs, *wf)
	}

	// ========== 6. ASSEMBLE RESPONSE ==========
	log.Println("Step 6: Assembling response...")

	result := &product.ProductRequest{
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

// services/product_service.go
func (s *ProductService) GetProductConfig(ctx context.Context, productID string, draft draft.DraftDataPost) (*product.ProductConfig, error) {
	var cfg product.ProductConfig

	log.Printf("========== GET PRODUCT CONFIG ==========")
	log.Printf("PRODUCT ID: %s", productID)

	// 1. Get product data from repository
	productData, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		log.Printf("[ERROR] Failed to get product: %v", err)
		return nil, fmt.Errorf("failed to get product %s: %w", productID, err)
	}

	// 2. Set basic product config
	cfg.ProductID = productData.ID
	cfg.APIEndpoint = productData.APIEndpoint
	cfg.APIKey = productData.APIKeyEncrypted

	log.Printf("PRODUCT FOUND: %s", cfg.ProductID)
	log.Printf("API ENDPOINT: %s", cfg.APIEndpoint)

	if cfg.APIEndpoint == "" {
		log.Printf("[ERROR] EMPTY API ENDPOINT")
		return nil, fmt.Errorf("product %s has empty API endpoint", productID)
	}

	// 3. Set default config values
	if productData.AdapterConfig != nil {
		cfg.Timeout = productData.AdapterConfig.TimeoutSeconds
		cfg.RetryCount = productData.AdapterConfig.RetryCount
	} else {
		cfg.Timeout = 30   // default timeout
		cfg.RetryCount = 3 // default retry
	}
	cfg.CustomHeaders = make(map[string]string)
	cfg.MetaConfigStr = "{}"
	cfg.SitemapConfigStr = "{}"

	// 4. Load adapter config if API key exists
	if cfg.APIKey != "" {
		if err := s.LoadAdapterConfig(ctx, &cfg); err != nil {
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

	// 5. Build full URL
	cfg.FullURL = strings.TrimRight(cfg.APIEndpoint, "/") + "/" + strings.TrimLeft(cfg.AdapterEndpoint, "/")
	log.Printf("FULL URL: %s", cfg.FullURL)

	// 6. Reorder workflow nodes and store them in config
	if len(productData.Workflows.Nodes) > 0 {
		log.Printf("REORDERING WORKFLOW NODES (Total: %d)", len(productData.Workflows.Nodes))

		// Convert to pointer slice for reordering
		nodes := make([]*workflow_node.WorkflowNode, len(productData.Workflows.Nodes))
		for i := range productData.Workflows.Nodes {
			nodes[i] = &productData.Workflows.Nodes[i]
		}

		// Reorder nodes with batch
		sortedNodes, levelMap := s.ReorderNodesWithBatch(nodes)

		// Store reordered nodes in config
		cfg.WorkflowNodes = make([]workflow_node.WorkflowNode, len(sortedNodes))
		for i, node := range sortedNodes {
			if node != nil {
				cfg.WorkflowNodes[i] = *node
			}
		}

		// Also store level map for parallel execution if needed
		cfg.WorkflowLevelMap = levelMap

		// Log batch information
		log.Printf("NODES REORDERED INTO %d BATCHES", len(levelMap))
		for level, batch := range levelMap {
			nodeIDs := make([]string, len(batch))
			for i, node := range batch {
				if node != nil {
					nodeIDs[i] = node.ID
				}
			}
			log.Printf("BATCH %d: %v", level, nodeIDs)
		}
	} else {
		log.Printf("[INFO] NO WORKFLOW NODES TO REORDER")
		cfg.WorkflowNodes = []workflow_node.WorkflowNode{}
		cfg.WorkflowLevelMap = make(map[int][]*workflow_node.WorkflowNode)
	}

	log.Printf("========== PRODUCT CONFIG READY ==========")
	return &cfg, nil
}

// Fixed ReorderNodesWithBatch with correct parameter type
func (s *ProductService) ReorderNodesWithBatch(nodes []*workflow_node.WorkflowNode) ([]*workflow_node.WorkflowNode, map[int][]*workflow_node.WorkflowNode) {
	if len(nodes) == 0 {
		return nodes, make(map[int][]*workflow_node.WorkflowNode)
	}

	// Filter out nil nodes
	validNodes := make([]*workflow_node.WorkflowNode, 0, len(nodes))
	for _, node := range nodes {
		if node != nil {
			validNodes = append(validNodes, node)
		}
	}

	if len(validNodes) == 0 {
		return nodes, make(map[int][]*workflow_node.WorkflowNode)
	}

	// Reorder nodes based on level
	sortedNodes := s.ReorderNodesByLevel(validNodes)
	if sortedNodes == nil {
		return nodes, make(map[int][]*workflow_node.WorkflowNode)
	}

	// Create map for quick access
	nodeMap := make(map[string]*workflow_node.WorkflowNode)
	for _, node := range sortedNodes {
		if node != nil {
			nodeMap[node.ID] = node
		}
	}

	// Calculate levels for each node
	levels := make(map[string]int)
	var calculateLevel func(nodeID string) int
	calculateLevel = func(nodeID string) int {
		// Check cache
		if val, ok := levels[nodeID]; ok {
			return val
		}

		node := nodeMap[nodeID]
		if node == nil {
			levels[nodeID] = 0
			return 0
		}

		// Root node (no previous nodes)
		if len(node.PreviousNodeIDs) == 0 {
			levels[nodeID] = 0
			return 0
		}

		// Calculate maximum level from previous nodes
		maxLevel := 0
		for _, prevID := range node.PreviousNodeIDs {
			if prevID == "" {
				continue
			}
			level := calculateLevel(prevID) + 1
			if level > maxLevel {
				maxLevel = level
			}
		}
		levels[nodeID] = maxLevel
		return maxLevel
	}

	// Calculate levels for all nodes
	for _, node := range sortedNodes {
		if node != nil {
			calculateLevel(node.ID)
		}
	}

	// Group by level
	levelMap := make(map[int][]*workflow_node.WorkflowNode)
	for _, node := range sortedNodes {
		if node != nil {
			level := levels[node.ID]
			levelMap[level] = append(levelMap[level], node)
		}
	}

	// Update step_order based on sorted order
	for i, node := range sortedNodes {
		if node != nil {
			node.StepOrder = i
		}
	}

	return sortedNodes, levelMap
}

// Helper function: ReorderNodesByLevel
func (s *ProductService) ReorderNodesByLevel(nodes []*workflow_node.WorkflowNode) []*workflow_node.WorkflowNode {
	if len(nodes) == 0 {
		return nodes
	}

	// Filter nil nodes
	validNodes := make([]*workflow_node.WorkflowNode, 0, len(nodes))
	for _, node := range nodes {
		if node != nil {
			validNodes = append(validNodes, node)
		}
	}

	if len(validNodes) == 0 {
		return nodes
	}

	// Create map for quick access
	nodeMap := make(map[string]*workflow_node.WorkflowNode)
	for _, node := range validNodes {
		if node != nil {
			nodeMap[node.ID] = node
		}
	}

	// Calculate level for each node
	levels := make(map[string]int)
	var calculateLevel func(nodeID string) int
	calculateLevel = func(nodeID string) int {
		if val, ok := levels[nodeID]; ok {
			return val
		}

		node := nodeMap[nodeID]
		if node == nil || len(node.PreviousNodeIDs) == 0 {
			levels[nodeID] = 0
			return 0
		}

		maxLevel := 0
		for _, prevID := range node.PreviousNodeIDs {
			if prevID == "" {
				continue
			}
			level := calculateLevel(prevID) + 1
			if level > maxLevel {
				maxLevel = level
			}
		}
		levels[nodeID] = maxLevel
		return maxLevel
	}

	// Calculate levels for all nodes
	for _, node := range validNodes {
		if node != nil {
			calculateLevel(node.ID)
		}
	}

	// Sort by level, if same level sort by step_order
	sortedNodes := make([]*workflow_node.WorkflowNode, len(validNodes))
	copy(sortedNodes, validNodes)

	sort.Slice(sortedNodes, func(i, j int) bool {
		if sortedNodes[i] == nil || sortedNodes[j] == nil {
			return false
		}
		// Sort by level first, then by step_order
		if levels[sortedNodes[i].ID] != levels[sortedNodes[j].ID] {
			return levels[sortedNodes[i].ID] < levels[sortedNodes[j].ID]
		}
		return sortedNodes[i].StepOrder < sortedNodes[j].StepOrder
	})

	return sortedNodes
}

// Add these helper methods if not exist
func (s *ProductService) LoadAdapterConfig(ctx context.Context, cfg *product.ProductConfig) error {
	log.Printf("API KEY EXISTS => LOADING ADAPTER CONFIG")

	// Get data from repository
	adapterData, err := s.adapterConfigRepo.GetByProductID(ctx, cfg.ProductID)
	if err != nil {
		log.Printf("[ERROR] FAILED QUERY ADAPTER CONFIG: %v", err)
		return fmt.Errorf("failed to query adapter config for product %s: %w", cfg.ProductID, err)
	}

	// Map data to config if found
	if adapterData != nil {
		cfg.AdapterEndpoint = adapterData.EndpointPath
		cfg.HTTPMethod = adapterData.HTTPMethod
		cfg.CustomHeadersStr = adapterData.CustomHeaders
		cfg.FieldMappingStr = adapterData.FieldMapping // Add this field
		cfg.MetaConfigStr = adapterData.MetaConfig
		cfg.SitemapConfigStr = adapterData.SitemapConfig
		if adapterData.TimeoutSeconds > 0 {
			cfg.Timeout = adapterData.TimeoutSeconds
		}
		if adapterData.RetryCount > 0 {
			cfg.RetryCount = adapterData.RetryCount
		}
	}

	log.Printf("ADAPTER CONFIG LOADED")
	log.Printf("HTTP METHOD: %s", cfg.HTTPMethod)
	log.Printf("FIELD MAPPING: %s", cfg.FieldMappingStr)
	log.Printf("META CONFIG: %s", cfg.MetaConfigStr)
	log.Printf("SITEMAP CONFIG: %s", cfg.SitemapConfigStr)

	// Parse custom headers
	if err := s.ParseCustomHeaders(cfg); err != nil {
		log.Printf("[WARN] Failed to parse custom headers: %v", err)
		cfg.CustomHeaders = make(map[string]string)
	}

	return nil
}

func (s *ProductService) ParseCustomHeaders(cfg *product.ProductConfig) error {
	if cfg.CustomHeadersStr == "" || cfg.CustomHeadersStr == "{}" {
		cfg.CustomHeaders = make(map[string]string)
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

// sendWithRetry method (if not exists)
func (s *ProductService) SendWithRetry(cfg product.ProductConfig, requestBody interface{}) (interface{}, error) {
	// Implement your retry logic here
	// This is a placeholder
	return nil, nil
}

// markProductSynced method (if not exists)
func (s *ProductService) markProductSynced(productID string) error {
	// Implement your sync marking logic here
	// This is a placeholder
	return nil
}
