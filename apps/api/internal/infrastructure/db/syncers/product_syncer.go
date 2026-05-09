package syncers

import (
	"database/sql"
	"log"
)

type ProductSyncer struct {
	db *sql.DB
}

func NewProductSyncer(db *sql.DB) *ProductSyncer {
	return &ProductSyncer{db: db}
}

func (s *ProductSyncer) MarkProductSynced(productID string) {
	_, err := s.db.Exec(`
		UPDATE products
		SET sync_status = 'connected',
			last_sync = NOW()
		WHERE id = $1
	`, productID)

	if err != nil {
		log.Printf("[WARN] sync update failed: %v", err)
	}
}
