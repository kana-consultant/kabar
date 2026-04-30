package product

import (
	"database/sql"
	"fmt"
	"log"

	"seo-backend/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(id string) (*models.Product, error) {
	query := `
		SELECT id, name, platform, api_key_encrypted, api_endpoint,
			status, sync_status, last_sync, created_by, team_id,
			user_id, created_at, updated_at
		FROM products WHERE id = $1
	`

	var product models.Product
	err := r.db.QueryRow(query, id).Scan(
		&product.ID, &product.Name, &product.Platform, &product.APIKeyEncrypted,
		&product.APIEndpoint, &product.Status, &product.SyncStatus, &product.LastSync,
		&product.CreatedBy, &product.TeamID, &product.UserID,
		&product.CreatedAt, &product.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &product, nil
}

func (r *Repository) GetAll(query string, args []interface{}) ([]models.Product, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}
	defer rows.Close()

	return r.scanProducts(rows)
}

func (r *Repository) Delete(id string) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

func (r *Repository) UpdateConnectionStatus(productID string, isConnected bool) {
	status := "connected"
	if !isConnected {
		status = "pending"
	}

	_, err := r.db.Exec(`
		UPDATE products 
		SET status = $1, last_sync = NOW(), updated_at = NOW()
		WHERE id = $2
	`, status, productID)

	if err != nil {
		log.Printf("ERROR update product status: %v", err)
	}
}

func (r *Repository) GetProductBasicInfo(id string) (*ProductBasicInfo, error) {
	var info ProductBasicInfo
	query := `
		SELECT id, name, platform, api_endpoint, COALESCE(api_key_encrypted, '')
		FROM products WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&info.ID, &info.Name, &info.Platform, &info.APIEndpoint, &info.APIKey,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product info: %w", err)
	}
	return &info, nil
}

func (r *Repository) scanProducts(rows *sql.Rows) ([]models.Product, error) {
	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(
			&p.ID, &p.Name, &p.Platform, &p.APIEndpoint,
			&p.Status, &p.SyncStatus, &p.LastSync,
			&p.CreatedBy, &p.TeamID, &p.UserID,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		products = append(products, p)
	}
	return products, nil
}
