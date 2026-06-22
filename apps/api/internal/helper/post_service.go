package helper

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"seo-backend/internal/database"
	"seo-backend/internal/domain/draft"
	"seo-backend/internal/domain/product"
	"seo-backend/internal/models"
)

type PostService struct {
	db                *sql.DB
	productController product.ProductService
}

func NewPostService(db *sql.DB, productController product.ProductService) *PostService {
	return &PostService{
		db:                db,
		productController: productController,
	}
}

func (s *PostService) markProductSynced(productID string) error {
	database.GetDB().Exec(`
		UPDATE products
		SET sync_status='connected',
			last_sync=NOW()
		WHERE id=$1
	`, productID)

	return nil
}

func (s *PostService) ProcessDraftProducts(ctx context.Context, draft draft.DraftDataPost, userCtx models.UserContext) ([]map[string]interface{}, bool, bool, error) {
	log.Printf("TargetProducts : %v", draft.TargetProducts)

	if len(draft.TargetProducts) == 0 {
		return nil, false, true, fmt.Errorf("target products is required and cannot be empty")
	}

	var postResults []map[string]interface{}
	someFailed := false
	allFailed := true

	log.Printf("============= PROCESSING DRAFT FOR PRODUCTS =============")
	log.Printf("Draft Data: %+v", draft)

	for _, productID := range draft.TargetProducts {
		log.Printf("======================= PROCESSING PRODUCT: %s =======================", productID)

		// Get product config
		cfg, err := s.productController.GetProductConfig(ctx, productID, draft, userCtx)
		if err != nil {
			errMsg := fmt.Sprintf("failed to get product config for %s: %v", productID, err)
			log.Printf("[ERROR] %s", errMsg)

			result := map[string]interface{}{
				"product": productID,
				"success": false,
				"error":   errMsg,
				"synced":  false,
			}
			postResults = append(postResults, result)
			someFailed = true
			continue
		}

		log.Printf("Product Config loaded successfully for %s", productID)
		log.Printf("Full URL: %s", cfg.FullURL)
		log.Printf("HTTP Method: %s", cfg.HTTPMethod)

		// Process workflow nodes in order
		if cfg.HasWorkflowNodes() {
			log.Printf("Processing %d workflow nodes for product %s", len(cfg.WorkflowNodes), productID)

			// 🔑 KEY: Initialize ExecutionResults dan Variables di config
			cfg.ExecutionResults = make(map[string]interface{})

			// Execute nodes in order (already reordered)
			for _, node := range cfg.WorkflowNodes {
				log.Printf("Executing node: %s (Step: %d)", node.ID, node.StepOrder)

				// 🔑 STEP 2: Build request body from node and ENRICHED draft
				requestBody, err := BuildRequestBody(&node, draft, cfg)
				if err != nil {
					errMsg := fmt.Sprintf("failed to build request body for node %s: %v", node.ID, err)
					log.Printf("[ERROR] %s", errMsg)

					result := map[string]interface{}{
						"product": productID,
						"node":    node.ID,
						"success": false,
						"error":   errMsg,
						"synced":  false,
					}
					postResults = append(postResults, result)
					someFailed = true
					continue
				}

				log.Printf("Request body built for node %s: %+v", node.ID, requestBody)

				// Send request with retry
				response, err := s.productController.SendWithRetry(*cfg, requestBody)
				if err != nil {
					errMsg := fmt.Sprintf("failed to send request for node %s: %v", node.ID, err)
					log.Printf("[ERROR] %s", errMsg)

					result := map[string]interface{}{
						"product": productID,
						"node":    node.ID,
						"success": false,
						"error":   errMsg,
						"synced":  false,
					}
					postResults = append(postResults, result)
					someFailed = true
					continue
				}

				log.Printf("Node %s executed successfully", node.ID)

				// 🔑 STEP 3: STORE response menggunakan SetExecutionResult
				cfg.SetExecutionResult(node.ID, response)

				// Store successful result
				result := map[string]interface{}{
					"product":  productID,
					"node":     node.ID,
					"success":  true,
					"response": response,
					"synced":   false,
				}
				postResults = append(postResults, result)
				allFailed = false
			}

			// Mark product as synced if all nodes succeeded

		} else {
			log.Printf("[WARN] No workflow nodes found for product %s", productID)
			result := map[string]interface{}{
				"product": productID,
				"success": false,
				"error":   "no workflow nodes configured",
				"synced":  false,
			}
			postResults = append(postResults, result)
			someFailed = true
		}
	}

	if allFailed && len(draft.TargetProducts) > 0 {
		return postResults, someFailed, allFailed, fmt.Errorf("all products failed to process")
	}

	log.Printf("============= DRAFT PROCESSING COMPLETED =============")
	log.Printf("Total results: %d, Some failed: %v, All failed: %v", len(postResults), someFailed, allFailed)

	return postResults, someFailed, allFailed, nil
}

// ============================================================
// 🔑 HELPER: Enrich draft dengan previous results dari config
// ============================================================
func (s *PostService) enrichDraftWithPreviousResults(
	draftData draft.DraftDataPost,
	cfg *product.ProductConfig,
) (draft.DraftDataPost, map[string]interface{}) {

	if cfg == nil {
		return draftData, nil
	}

	results := cfg.GetAllExecutionResults()
	if len(results) == 0 {
		return draftData, nil
	}

	return draftData, results
}
