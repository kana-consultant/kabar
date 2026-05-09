package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"seo-backend/internal/domain/product"
)

type ProductRepository struct {
	db *sql.DB
}

// Constructor - memastikan implementasi interface
func NewProductRepository(db *sql.DB) product.ProductRepository {
	return &ProductRepository{db: db}
}

// =======================
// TRANSACTION MANAGEMENT
// =======================
func (r *ProductRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

// =======================
// BASIC CRUD (dengan optional tx)
// =======================
func (r *ProductRepository) GetByID(ctx context.Context, id string) (*product.Product, error) {
	return r.GetByIDWithTx(ctx, nil, id)
}

func (r *ProductRepository) GetByIDWithTx(ctx context.Context, tx *sql.Tx, id string) (*product.Product, error) {
	query := `
		SELECT id, name, platform, api_key_encrypted, api_endpoint,
			status, sync_status, last_sync, created_by, team_id,
			user_id, created_at, updated_at
		FROM products WHERE id = $1
	`

	var p product.Product
	var err error

	if tx != nil {
		err = tx.QueryRowContext(ctx, query, id).Scan(
			&p.ID, &p.Name, &p.Platform, &p.APIKeyEncrypted,
			&p.APIEndpoint, &p.Status, &p.SyncStatus, &p.LastSync,
			&p.CreatedBy, &p.TeamID, &p.UserID,
			&p.CreatedAt, &p.UpdatedAt,
		)
	} else {
		err = r.db.QueryRowContext(ctx, query, id).Scan(
			&p.ID, &p.Name, &p.Platform, &p.APIKeyEncrypted,
			&p.APIEndpoint, &p.Status, &p.SyncStatus, &p.LastSync,
			&p.CreatedBy, &p.TeamID, &p.UserID,
			&p.CreatedAt, &p.UpdatedAt,
		)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &p, nil
}

// GetAllWithFilters - Build query sendiri dengan filters
func (r *ProductRepository) GetAllWithFilters(ctx context.Context) ([]product.Product, int, error) {
	// Base query
	baseQuery := `
		SELECT id, name, platform, api_endpoint, status, sync_status, 
			last_sync, created_by, team_id, user_id, created_at, updated_at
		FROM products
		WHERE 1=1
	`

	// Count query
	countQuery := `
		SELECT COUNT(*)
		FROM products
		WHERE 1=1
	`

	args := []interface{}{}

	// Get total count
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Execute query
	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch products: %w", err)
	}
	defer rows.Close()

	products, err := r.scanProducts(rows)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// GetAll - legacy method, tetap menerima query dari luar untuk kompatibilitas
func (r *ProductRepository) GetAll(ctx context.Context, query string, args []interface{}) ([]product.Product, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}
	defer rows.Close()

	return r.scanProducts(rows)
}

// Delete - delete product
func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	return r.DeleteWithTx(ctx, nil, id)
}

func (r *ProductRepository) DeleteWithTx(ctx context.Context, tx *sql.Tx, id string) error {
	query := "DELETE FROM products WHERE id = $1"

	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.ExecContext(ctx, query, id)
	} else {
		result, err = r.db.ExecContext(ctx, query, id)
	}

	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// =======================
// BUSINESS SPECIFIC OPERATIONS (dengan optional tx)
// =======================
func (r *ProductRepository) UpdateConnectionStatus(ctx context.Context, productID string, isConnected bool) error {
	return r.UpdateConnectionStatusWithTx(ctx, nil, productID, isConnected)
}

func (r *ProductRepository) UpdateConnectionStatusWithTx(ctx context.Context, tx *sql.Tx, productID string, isConnected bool) error {
	status := "connected"
	if !isConnected {
		status = "pending"
	}

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, `
			UPDATE products 
			SET status = $1, last_sync = NOW(), updated_at = NOW()
			WHERE id = $2
		`, status, productID)
	} else {
		_, err = r.db.ExecContext(ctx, `
			UPDATE products 
			SET status = $1, last_sync = NOW(), updated_at = NOW()
			WHERE id = $2
		`, status, productID)
	}

	if err != nil {
		return fmt.Errorf("failed to update product status: %w", err)
	}

	return nil
}

func (r *ProductRepository) GetProductBasicInfo(ctx context.Context, id string) (*product.ProductBasicInfo, error) {
	return r.GetProductBasicInfoWithTx(ctx, nil, id)
}

func (r *ProductRepository) GetProductBasicInfoWithTx(ctx context.Context, tx *sql.Tx, id string) (*product.ProductBasicInfo, error) {
	query := `
		SELECT id, name,custom_headers, platform, api_endpoint, COALESCE(api_key_encrypted, '')
		FROM products WHERE id = $1
	`

	var info product.ProductBasicInfo
	var err error

	if tx != nil {
		err = tx.QueryRowContext(ctx, query, id).Scan(
			&info.ID,
			&info.Name,
			&info.Platform,
			&info.APIEndpoint,
			&info.APIKey,
		)
	} else {
		err = r.db.QueryRowContext(ctx, query, id).Scan(
			&info.ID,
			&info.Name,
			&info.Platform,
			&info.APIEndpoint,
			&info.APIKey,
		)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product info: %w", err)
	}

	return &info, nil
}

func (r *ProductRepository) GetProductsByTeamID(ctx context.Context, teamID string) ([]product.Product, error) {
	query := `
		SELECT id, name, platform, api_endpoint, status, sync_status, 
			last_sync, created_by, team_id, user_id, created_at, updated_at
		FROM products WHERE team_id = $1 ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products by team: %w", err)
	}
	defer rows.Close()

	return r.scanProducts(rows)
}

func (r *ProductRepository) GetProductsByUserID(ctx context.Context, userID string) ([]product.Product, error) {
	query := `
		SELECT id, name, platform, api_endpoint, status, sync_status, 
			last_sync, created_by, team_id, user_id, created_at, updated_at
		FROM products WHERE user_id = $1 ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products by user: %w", err)
	}
	defer rows.Close()

	return r.scanProducts(rows)
}

// =======================
// TRANSACTION-AWARE OPERATIONS (wajib pake tx)
// =======================
func (r *ProductRepository) InsertProductWithTx(ctx context.Context, tx *sql.Tx, req product.CreateProductRequest) (string, error) {
	if tx == nil {
		return "", fmt.Errorf("transaction is required for InsertProductWithTx")
	}

	query := `
		INSERT INTO products (
			name, platform, api_endpoint, api_key_encrypted,
			team_id, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id
	`

	var productID string
	err := tx.QueryRowContext(ctx, query,
		req.Name, req.Platform, req.APIEndpoint, req.APIKey, req.TeamID,
	).Scan(&productID)

	if err != nil {
		return "", fmt.Errorf("failed to insert product: %w", err)
	}

	return productID, nil
}

func (r *ProductRepository) UpdateProductWithTx(ctx context.Context, tx *sql.Tx, id string, updates map[string]interface{}) error {
	if tx == nil {
		return fmt.Errorf("transaction is required for UpdateProductWithTx")
	}

	if len(updates) == 0 {
		return nil
	}

	setClauses := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1

	// Field mapping: application field -> database column
	fieldMap := map[string]string{
		"name":        "name",
		"platform":    "platform",
		"apiEndpoint": "api_endpoint",
		"status":      "status",
		"syncStatus":  "sync_status",
		"apiKey":      "api_key_encrypted",
	}

	for key, value := range updates {
		if dbField, ok := fieldMap[key]; ok {
			setClauses = append(setClauses, fmt.Sprintf("%s = $%d", dbField, argIndex))
			args = append(args, value)
			argIndex++
		}
	}

	// Always update timestamp
	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	// Add WHERE clause
	args = append(args, id)
	query := fmt.Sprintf("UPDATE products SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "), argIndex)

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with id %s not found", id)
	}

	return nil
}

func (r *ProductRepository) GetProductByIDWithTx(ctx context.Context, tx *sql.Tx, id string) (*product.Product, error) {
	if tx == nil {
		return nil, fmt.Errorf("transaction is required for GetProductByIDWithTx")
	}

	query := `SELECT id, name, platform, api_endpoint, status, sync_status, created_at, updated_at 
	          FROM products WHERE id = $1 FOR UPDATE`

	var product product.Product
	err := tx.QueryRowContext(ctx, query, id).Scan(
		&product.ID, &product.Name, &product.Platform, &product.APIEndpoint,
		&product.Status, &product.SyncStatus, &product.CreatedAt, &product.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &product, nil
}

func (r *ProductRepository) DeleteProductWithTx(ctx context.Context, tx *sql.Tx, id string) error {
	if tx == nil {
		return fmt.Errorf("transaction is required for DeleteProductWithTx")
	}

	query := `DELETE FROM products WHERE id = $1`

	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with id %s not found", id)
	}

	return nil
}

// =======================
// HELPER FUNCTIONS
// =======================
func (r *ProductRepository) scanProducts(rows *sql.Rows) ([]product.Product, error) {
	var products []product.Product

	for rows.Next() {
		var p product.Product

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Platform,
			&p.APIEndpoint,
			&p.Status,
			&p.SyncStatus,
			&p.LastSync,
			&p.CreatedBy,
			&p.TeamID,
			&p.UserID,
			&p.CreatedAt,
			&p.UpdatedAt,
		)

		if err != nil {
			log.Printf("error scanning product row: %v", err)
			continue
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return products, nil
}
