package product

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"seo-backend/internal/models"
)

type TransactionManager struct {
	db *sql.DB
}

func NewTransactionManager(db *sql.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

func (tm *TransactionManager) CreateProduct(req models.CreateProductRequest, ctx UserContext) (string, error) {
	tx, err := tm.db.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	productID, err := tm.insertProduct(tx, req, ctx)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return productID, nil
}

func (tm *TransactionManager) UpdateProduct(id string, updates map[string]interface{}, ctx UserContext) error {
	tx, err := tm.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := tm.updateProduct(tx, id, updates); err != nil {
		return err
	}

	if adapterUpdates, ok := updates["adapterConfig"].(map[string]interface{}); ok {
		adapterRepo := NewAdapterConfigRepo(tm.db)
		if err := adapterRepo.Update(tx, id, adapterUpdates); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (tm *TransactionManager) insertProduct(tx *sql.Tx, req models.CreateProductRequest, ctx UserContext) (string, error) {
	query := `
		INSERT INTO products (
			name, platform, api_endpoint, api_key_encrypted,
			status, sync_status, created_by, team_id, user_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	var productID string
	err := tx.QueryRow(
		query,
		req.Name, req.Platform, req.APIEndpoint, req.APIKey,
		"pending", "idle", nullIfEmpty(ctx.GetUserID()),
		nullIfEmpty(ctx.GetTeamID()), nullIfEmpty(ctx.GetUserID()),
	).Scan(&productID)

	return productID, err
}

func (tm *TransactionManager) updateProduct(tx *sql.Tx, id string, updates map[string]interface{}) error {
	setClauses := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1

	fieldMap := map[string]string{
		"name":        "name",
		"platform":    "platform",
		"apiEndpoint": "api_endpoint",
		"status":      "status",
		"apiKey":      "api_key_encrypted",
	}

	for key, value := range updates {
		if dbField, ok := fieldMap[key]; ok {
			setClauses = append(setClauses, fmt.Sprintf("%s = $%d", dbField, argIndex))
			args = append(args, value)
			argIndex++
		}
	}

	if len(setClauses) == 0 {
		return nil
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	args = append(args, id)
	query := fmt.Sprintf("UPDATE products SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "), argIndex)

	_, err := tx.Exec(query, args...)
	return err
}

func nullIfEmpty(id string) interface{} {
	if id == "" || id == "00000000-0000-0000-0000-000000000000" {
		return nil
	}
	return id
}
