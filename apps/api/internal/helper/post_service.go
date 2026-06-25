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
	"time"
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
	var allErrors []string
	someFailed := false
	allFailed := true
	stopOnError := true // Flag to stop processing on error

	for idx, productID := range draft.TargetProducts {
		log.Printf("[DRAFT] Processing product %d/%d: %s", idx+1, len(draft.TargetProducts), productID)

		// Get product config
		cfg, err := s.productController.GetProductConfig(ctx, productID, draft, userCtx)
		if err != nil {
			errMsg := fmt.Sprintf("failed to get product config: %v", err)
			log.Printf("[ERROR] Product %s: %s", productID, errMsg)

			result := map[string]interface{}{
				"product":    productID,
				"success":    false,
				"error":      errMsg,
				"synced":     false,
				"statusCode": 0,
			}
			postResults = append(postResults, result)
			allErrors = append(allErrors, fmt.Sprintf("Product %s: %s", productID, errMsg))
			someFailed = true

			if stopOnError {
				log.Printf("[DRAFT] Stopping processing due to error on product %s", productID)
				break // Stop processing all products
			}
			continue
		}

		// Validate config
		if cfg.FullURL == "" {
			errMsg := "full URL is empty"
			log.Printf("[ERROR] Product %s: %s", productID, errMsg)

			result := map[string]interface{}{
				"product":    productID,
				"success":    false,
				"error":      errMsg,
				"synced":     false,
				"statusCode": 0,
			}
			postResults = append(postResults, result)
			allErrors = append(allErrors, fmt.Sprintf("Product %s: %s", productID, errMsg))
			someFailed = true

			if stopOnError {
				log.Printf("[DRAFT] Stopping processing due to error on product %s", productID)
				break
			}
			continue
		}

		if !cfg.HasWorkflowNodes() {
			errMsg := "no workflow nodes configured"
			log.Printf("[WARNING] Product %s: %s", productID, errMsg)

			result := map[string]interface{}{
				"product":    productID,
				"success":    false,
				"error":      errMsg,
				"synced":     false,
				"statusCode": 0,
			}
			postResults = append(postResults, result)
			allErrors = append(allErrors, fmt.Sprintf("Product %s: %s", productID, errMsg))
			someFailed = true

			if stopOnError {
				log.Printf("[DRAFT] Stopping processing due to error on product %s", productID)
				break
			}
			continue
		}

		// Initialize execution context
		cfg.ExecutionResults = make(map[string]interface{})
		cfg.ExecutionResults["product_id"] = cfg.ProductID
		cfg.ExecutionResults["variables"] = make(map[string]interface{})
		cfg.ExecutionResults["started_at"] = time.Now().Format(time.RFC3339)

		productFailed := false
		var productErrors []string

		// Process each node
		for nodeIdx, node := range cfg.WorkflowNodes {
			log.Printf("[DRAFT] Product %s: Processing node %d/%d: %s", productID, nodeIdx+1, len(cfg.WorkflowNodes), node.ID)

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
				errMsg := fmt.Sprintf("failed to build request: %v", err)
				log.Printf("[ERROR] Product %s, Node %s: %s", productID, node.ID, errMsg)

				result := map[string]interface{}{
					"product":    productID,
					"node":       node.ID,
					"success":    false,
					"error":      errMsg,
					"synced":     false,
					"statusCode": 0,
				}
				postResults = append(postResults, result)
				productErrors = append(productErrors, fmt.Sprintf("Node %s: %s", node.ID, errMsg))
				productFailed = true
				someFailed = true
				break // STOP processing remaining nodes for this product
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

				result := map[string]interface{}{
					"product":    productID,
					"node":       node.ID,
					"success":    false,
					"error":      errorMsg,
					"synced":     false,
					"statusCode": statusCode,
				}
				postResults = append(postResults, result)
				productErrors = append(productErrors, fmt.Sprintf("Node %s: %s", node.ID, errorMsg))
				productFailed = true
				someFailed = true
				break // STOP processing remaining nodes for this product
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

			// Success result
			result := map[string]interface{}{
				"product":    productID,
				"node":       node.ID,
				"success":    true,
				"response":   string(response),
				"synced":     false,
				"statusCode": statusCode,
			}
			postResults = append(postResults, result)

			log.Printf("[SUCCESS] Product %s, Node %s: Status %d", productID, node.ID, statusCode)
		}

		// If product had any failure, mark it and optionally stop
		if productFailed {
			someFailed = true

			// Add product summary to errors
			allErrors = append(allErrors, fmt.Sprintf("Product %s failed: %v", productID, productErrors))

			// Check if we should stop processing remaining products
			if stopOnError {
				log.Printf("[DRAFT] Stopping processing due to product failure: %s", productID)
				break // Stop processing remaining products
			}
		} else {
			log.Printf("[SUCCESS] Product %s: All nodes processed successfully", productID)
		}
	}

	// Calculate final statistics
	totalProcessed := len(postResults)
	successCount := 0
	failedCount := 0

	for _, result := range postResults {
		if success, ok := result["success"].(bool); ok {
			if success {
				successCount++
			} else {
				failedCount++
			}
		}
	}

	// Final summary
	log.Printf("[DRAFT] Processing completed - Total results: %d, Success: %d, Failed: %d",
		totalProcessed, successCount, failedCount)

	if someFailed {
		log.Printf("[DRAFT] Completed with errors: %v", allErrors)
	} else {
		log.Printf("[DRAFT] All products processed successfully")
	}

	// Return error if all failed
	if allFailed && len(draft.TargetProducts) > 0 {
		return postResults, someFailed, allFailed, fmt.Errorf("all products failed to process: %v", allErrors)
	}

	// Return partial success if some failed but not all
	if someFailed {
		return postResults, someFailed, allFailed, fmt.Errorf("some products failed to process: %v", allErrors)
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
