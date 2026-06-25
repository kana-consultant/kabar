package helper

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"seo-backend/internal/database"
	"seo-backend/internal/domain/draft"
	"seo-backend/internal/domain/product"
	"seo-backend/internal/models"
	"strings"
)

type PostService struct {
	db                *sql.DB
	productController product.ProductService
}

func (s *PostService) ProcessDraftProducts(ctx context.Context, draft draft.DraftDataPost, userCtx models.UserContext) ([]map[string]interface{}, bool, bool, error) {
	log.Printf("[DRAFT] Processing %d products: %v", len(draft.TargetProducts), draft.TargetProducts)

	if len(draft.TargetProducts) == 0 {
		return nil, false, true, fmt.Errorf("target products is required and cannot be empty")
	}

	var postResults []map[string]interface{}
	someFailed := false
	allFailed := true

	for _, productID := range draft.TargetProducts {
		// Get product config
		cfg, err := s.productController.GetProductConfig(ctx, productID, draft, userCtx)
		if err != nil {
			log.Printf("[ERROR] Product %s: failed to get config - %v", productID, err)
			postResults = append(postResults, map[string]interface{}{
				"product": productID,
				"success": false,
				"error":   fmt.Sprintf("failed to get product config: %v", err),
				"synced":  false,
			})
			someFailed = true
			continue
		}

		// Validate config
		if cfg.FullURL == "" {
			log.Printf("[ERROR] Product %s: empty URL", productID)
			postResults = append(postResults, map[string]interface{}{
				"product": productID,
				"success": false,
				"error":   "full URL is empty",
				"synced":  false,
			})
			someFailed = true
			continue
		}

		if !cfg.HasWorkflowNodes() {
			log.Printf("[WARNING] Product %s: no workflow nodes", productID)
			postResults = append(postResults, map[string]interface{}{
				"product": productID,
				"success": false,
				"error":   "no workflow nodes configured",
				"synced":  false,
			})
			someFailed = true
			continue
		}

		// Initialize execution context
		cfg.ExecutionResults = make(map[string]interface{})
		cfg.ExecutionResults["product_id"] = cfg.ProductID
		cfg.ExecutionResults["variables"] = make(map[string]interface{})

		productFailed := false

		// Process each node
		for _, node := range cfg.WorkflowNodes {
			// Set current node
			cfg.CurrentNodeID = node.ID

			// Update config from node if available
			if node.AdapterConfig != nil {
				if node.AdapterConfig.HTTPMethod != "" {
					cfg.HTTPMethod = node.AdapterConfig.HTTPMethod
				}
				if node.AdapterConfig.EndpointPath != "" {
					cfg.FullURL = strings.TrimRight(cfg.APIEndpoint, "/") + "/" +
						strings.TrimLeft(node.AdapterConfig.EndpointPath, "/")
					cfg.ExecutionResults["endpoint_path"] = node.AdapterConfig.EndpointPath
				}
			}

			// Build request body
			requestBody, err := BuildRequestBody(&node, draft, cfg)
			if err != nil {
				log.Printf("[ERROR] Product %s, Node %s: build request failed - %v", productID, node.ID, err)
				postResults = append(postResults, map[string]interface{}{
					"product":    productID,
					"node":       node.ID,
					"success":    false,
					"error":      fmt.Sprintf("failed to build request: %v", err),
					"synced":     false,
					"statusCode": 0,
				})
				someFailed = true
				productFailed = true
				continue
			}

			// Send request with retry and get status code
			response, statusCode, err := SendWithRetry(cfg, requestBody)

			// Check if request failed or status code is not 2xx
			if err != nil || statusCode < 200 || statusCode >= 300 {
				errorMsg := fmt.Sprintf("request failed with status %d: %v", statusCode, err)
				if err == nil {
					errorMsg = fmt.Sprintf("request returned non-success status code: %d", statusCode)
				}

				log.Printf("[ERROR] Product %s, Node %s: %s", productID, node.ID, errorMsg)
				postResults = append(postResults, map[string]interface{}{
					"product":    productID,
					"node":       node.ID,
					"success":    false,
					"error":      errorMsg,
					"synced":     false,
					"statusCode": statusCode,
				})
				someFailed = true
				productFailed = true

				// Optional: Stop processing further nodes for this product on failure
				// break
				continue
			}

			// Store result (success case)
			cfg.SetExecutionResult(node.ID, response)
			allFailed = false

			// Parse & store response
			if len(response) > 0 {
				var responseData map[string]interface{}
				if err := json.Unmarshal(response, &responseData); err == nil {
					for key, value := range responseData {
						cfg.ExecutionResults["response_"+key] = value
						if variables, ok := cfg.ExecutionResults["variables"].(map[string]interface{}); ok {
							variables[key] = value
						}
					}
				}
			}

			postResults = append(postResults, map[string]interface{}{
				"product":    productID,
				"node":       node.ID,
				"success":    true,
				"response":   string(response),
				"synced":     false,
				"statusCode": statusCode,
			})
		}

		// If product had any failure, mark it
		if productFailed {
			someFailed = true
		}
	}

	// Final summary
	if someFailed {
		log.Printf("[DRAFT] Completed with errors - Total: %d, Failed: %v", len(postResults), someFailed)
	} else {
		log.Printf("[DRAFT] All products processed successfully - Total: %d", len(postResults))
	}

	if allFailed && len(draft.TargetProducts) > 0 {
		return postResults, someFailed, allFailed, fmt.Errorf("all products failed to process")
	}

	return postResults, someFailed, allFailed, nil
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
