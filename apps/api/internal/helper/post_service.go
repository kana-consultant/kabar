package helper

import (
	"database/sql"
	"fmt"
	"log"
	"seo-backend/internal/database"
	"seo-backend/internal/domain/draft"
)

type PostService struct {
	db *sql.DB
}

func NewPostService(db *sql.DB) *PostService {
	return &PostService{
		db: db,
	}
}

func (s *PostService) ProcessDraftProducts(draft draft.DraftDataPost) ([]map[string]interface{}, bool, bool, error) {
	log.Printf("TargetProducts : %v", draft.TargetProducts)

	if len(draft.TargetProducts) == 0 {
		return nil, false, true, fmt.Errorf("%s", "target products is required and cannot be empty")
	}

	var postResults []map[string]interface{}
	someFailed := false
	allFailed := true

	log.Printf("=============%v", draft)

	for _, productID := range draft.TargetProducts {
		fmt.Printf("=======================,%v", productID)

		result, err := s.processSingleProduct(draft, productID)

		// Update status regardless of error
		if updateErr := s.updateConnectionStatus(productID, err == nil); updateErr != nil {
			log.Printf("failed update product status for %s: %v", productID, updateErr)
		}

		postResults = append(postResults, result)

		if err != nil {
			someFailed = true
			log.Printf("failed to process product %s: %v", productID, err)
		} else {
			allFailed = false
		}
	}

	if allFailed && len(draft.TargetProducts) > 0 {
		return postResults, someFailed, allFailed, fmt.Errorf("%s", "all products failed to process")
	}

	return postResults, someFailed, allFailed, nil
}

func (s *PostService) processSingleProduct(draft draft.DraftDataPost, productID string) (map[string]interface{}, error) {
	log.Printf("[START] Processing product: %s", productID)

	cfg, err := s.getProductConfig(productID)
	if err != nil {
		errMsg := fmt.Sprintf("failed to get product config for %s: %v", productID, err)
		log.Printf("[ERROR] %s", errMsg)
		return map[string]interface{}{
			"product": productID,
			"success": false,
			"error":   errMsg,
		}, fmt.Errorf("%s", errMsg)
	}

	requestBody, err := s.buildRequestBody(cfg, draft)
	if err != nil {
		errMsg := fmt.Sprintf("failed to build request body for %s: %v", productID, err)
		log.Printf("[ERROR] %s", errMsg)
		return map[string]interface{}{
			"product": productID,
			"success": false,
			"error":   errMsg,
		}, fmt.Errorf("%s", errMsg)
	}

	response, err := s.sendWithRetry(cfg, requestBody)
	if err != nil {
		errMsg := fmt.Sprintf("failed to send request for %s after retries: %v", productID, err)
		log.Printf("[ERROR] %s", errMsg)
		return map[string]interface{}{
			"product": productID,
			"success": false,
			"error":   errMsg,
		}, fmt.Errorf("%s", errMsg)
	}

	if err := s.markProductSynced(cfg.ProductID); err != nil {
		log.Printf("[WARN] failed to mark product %s as synced: %v", cfg.ProductID, err)
	}

	return map[string]interface{}{
		"product":    productID,
		"success":    true,
		"response":   string(response),
		"product_id": cfg.ProductID,
	}, nil
}

func (s *PostService) updateConnectionStatus(productID string, isConnected bool) error {
	status := "connected"
	if !isConnected {
		status = "pending"
	}

	_, err := s.db.Exec(`
		UPDATE products 
		SET status = $1, last_sync = NOW(), updated_at = NOW()
		WHERE id = $2
	`, status, productID)

	if err != nil {
		return fmt.Errorf("failed to update product status for product %s: %w", productID, err)
	}

	return nil
}

func (s *PostService) markProductSynced(productID string) error {
	result, err := database.GetDB().Exec(`
		UPDATE products
		SET sync_status='connected',
			last_sync=NOW()
		WHERE id=$1
	`, productID)

	if err != nil {
		return fmt.Errorf("failed to update sync status for product %s: %w", productID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected for product %s: %w", productID, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product %s not found when updating sync status", productID)
	}

	log.Printf("[INFO] Product %s marked as synced successfully", productID)
	return nil
}
