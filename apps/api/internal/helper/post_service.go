package helper

import (
	"database/sql"
	"seo-backend/internal/database"
)

type PostService struct {
	db *sql.DB
}

func NewPostService(db *sql.DB) *PostService {
	return &PostService{
		db: db,
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
