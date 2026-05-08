package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"seo-backend/internal/domain/apikey"
)

type APIKeyRepository struct {
	db *sql.DB
}

// NewAPIKeyRepository - constructor
func NewAPIKeyRepository(db *sql.DB) apikey.Repository {
	return &APIKeyRepository{db: db}
}

// BeginTx - start transaction
func (r *APIKeyRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

// Create - insert new API key
func (r *APIKeyRepository) Create(ctx context.Context, tx *sql.Tx, key *apikey.APIKey) (string, error) {
	if tx == nil {
		return "", fmt.Errorf("transaction is required for Create")
	}

	query := `
		INSERT INTO api_keys (
			service, provider_id, model_id, key_encrypted,
			is_active, system_prompt, team_id, created_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()
		)
		RETURNING id
	`

	var id string
	err := tx.QueryRowContext(ctx, query,
		key.Service, key.ProviderID, key.ModelID, key.KeyEncrypted,
		key.IsActive, key.SystemPrompt, key.TeamID, key.CreatedBy,
	).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("failed to create API key: %w", err)
	}

	return id, nil
}

// GetByID - get API key by ID
func (r *APIKeyRepository) GetByID(ctx context.Context, id string) (*apikey.APIKey, error) {
	query := `
		SELECT id, service, provider_id, model_id, key_encrypted,
		       is_active, system_prompt, team_id, created_by, created_at, updated_at
		FROM api_keys WHERE id = $1
	`

	var key apikey.APIKey
	var teamID sql.NullString
	var systemPrompt sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&key.ID, &key.Service, &key.ProviderID, &key.ModelID, &key.KeyEncrypted,
		&key.IsActive, &systemPrompt, &teamID, &key.CreatedBy,
		&key.CreatedAt, &key.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	if teamID.Valid {
		key.TeamID = &teamID.String
	}
	key.SystemPrompt = systemPrompt.String

	return &key, nil
}

// GetAll - get all API keys with provider and model details
func (r *APIKeyRepository) GetAll(ctx context.Context, teamID *string, userID string, userRole string) ([]apikey.APIKeyDetail, error) {
	query := `
		SELECT 
			ak.id, ak.service, ak.provider_id, ak.model_id,
			ak.is_active, ak.system_prompt, ak.created_by,
			ak.created_at, ak.updated_at,
			COALESCE(p.name, '') as provider_name,
			COALESCE(p.display_name, '') as provider_display_name,
			COALESCE(m.name, '') as model_name,
			COALESCE(m.display_name, '') as model_display_name
		FROM api_keys ak
		LEFT JOIN api_providers p ON ak.provider_id = p.id
		LEFT JOIN ai_models m ON ak.model_id = m.id
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	// Apply role-based filtering
	switch userRole {
	case "admin":
		if teamID != nil && *teamID != "" {
			query += fmt.Sprintf(" AND ak.team_id = $%d", argIndex)
			args = append(args, *teamID)
			argIndex++
		} else {
			query += fmt.Sprintf(" AND ak.created_by = $%d", argIndex)
			args = append(args, userID)
			argIndex++
		}
	default:
		query += fmt.Sprintf(" AND ak.created_by = $%d", argIndex)
		args = append(args, userID)
		argIndex++
	}

	query += " ORDER BY ak.created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query API keys: %w", err)
	}
	defer rows.Close()

	var results []apikey.APIKeyDetail

	for rows.Next() {
		var item apikey.APIKeyDetail
		var systemPrompt sql.NullString
		var createdBy sql.NullString

		err := rows.Scan(
			&item.ID, &item.Service, &item.ProviderID, &item.ModelID,
			&item.IsActive, &systemPrompt, &createdBy,
			&item.CreatedAt, &item.UpdatedAt,
			&item.ProviderName, &item.ProviderDisplayName,
			&item.ModelName, &item.ModelDisplayName,
		)
		if err != nil {
			continue
		}

		item.SystemPrompt = systemPrompt.String
		if createdBy.Valid {
			item.CreatedBy = &createdBy.String
		}

		results = append(results, item)
	}

	return results, nil
}

// Update - update API key fields
func (r *APIKeyRepository) Update(ctx context.Context, tx *sql.Tx, id string, updates map[string]interface{}) error {
	if tx == nil {
		return fmt.Errorf("transaction is required for Update")
	}

	if len(updates) == 0 {
		return nil
	}

	setClauses := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1

	fieldMap := map[string]string{
		"service":      "service",
		"providerId":   "provider_id",
		"modelId":      "model_id",
		"isActive":     "is_active",
		"systemPrompt": "system_prompt",
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
	query := fmt.Sprintf("UPDATE api_keys SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "), argIndex)

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update API key: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("API key with id %s not found", id)
	}

	return nil
}

// UpdateKey - update encrypted key only
func (r *APIKeyRepository) UpdateKey(ctx context.Context, tx *sql.Tx, id string, encryptedKey string) error {
	if tx == nil {
		return fmt.Errorf("transaction is required for UpdateKey")
	}

	query := `UPDATE api_keys SET key_encrypted = $1, updated_at = NOW() WHERE id = $2`

	result, err := tx.ExecContext(ctx, query, encryptedKey, id)
	if err != nil {
		return fmt.Errorf("failed to update API key: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("API key with id %s not found", id)
	}

	return nil
}

// Delete - soft delete or hard delete
func (r *APIKeyRepository) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	if tx == nil {
		return fmt.Errorf("transaction is required for Delete")
	}

	query := `DELETE FROM api_keys WHERE id = $1`

	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("API key with id %s not found", id)
	}

	return nil
}
