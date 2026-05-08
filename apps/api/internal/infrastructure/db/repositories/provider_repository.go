package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"seo-backend/internal/domain/provider"
)

type ProviderRepository struct {
	db *sql.DB
}

// NewProviderRepository creates a new provider repository
func NewProviderRepository(db *sql.DB) provider.Repository {
	return &ProviderRepository{db: db}
}

// BeginTx starts a new transaction
func (r *ProviderRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

// Create inserts a new provider
func (r *ProviderRepository) Create(ctx context.Context, tx *sql.Tx, p *provider.APIProvider) (string, error) {
	if tx == nil {
		return "", fmt.Errorf("transaction is required for Create")
	}

	// Marshal JSON fields
	defaultHeadersJSON, err := json.Marshal(p.DefaultHeaders)
	if err != nil {
		return "", fmt.Errorf("failed to marshal default headers: %w", err)
	}

	requestTemplateJSON, err := json.Marshal(p.RequestTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request template: %w", err)
	}

	query := `
		INSERT INTO api_providers (
			id, name, display_name, description, base_url,
			auth_type, auth_header, auth_prefix,
			text_endpoint, image_endpoint,
			default_headers, request_template,
			response_text_path, response_image_path,
			is_active, created_at, updated_at
		) VALUES (
			COALESCE($1, gen_random_uuid()), $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14,
			$15, NOW(), NOW()
		)
		RETURNING id
	`

	var id string
	err = tx.QueryRowContext(ctx, query,
		p.ID, p.Name, p.DisplayName, p.Description, p.BaseURL,
		p.AuthType, p.AuthHeader, p.AuthPrefix,
		p.TextEndpoint, p.ImageEndpoint,
		defaultHeadersJSON, requestTemplateJSON,
		p.ResponseTextPath, p.ResponseImagePath,
		p.IsActive,
	).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("failed to create provider: %w", err)
	}

	return id, nil
}

// GetByID retrieves a provider by ID
func (r *ProviderRepository) GetByID(ctx context.Context, id string) (*provider.APIProvider, error) {
	query := `
		SELECT 
			id, name, display_name, description, base_url,
			auth_type, auth_header, auth_prefix,
			text_endpoint, image_endpoint,
			default_headers, request_template,
			response_text_path, response_image_path,
			is_active, created_at, updated_at
		FROM api_providers 
		WHERE id = $1
	`

	var p provider.APIProvider
	var defaultHeadersJSON []byte
	var requestTemplateJSON []byte
	var imageEndpoint sql.NullString
	var responseImagePath sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.DisplayName, &p.Description, &p.BaseURL,
		&p.AuthType, &p.AuthHeader, &p.AuthPrefix,
		&p.TextEndpoint, &imageEndpoint,
		&defaultHeadersJSON, &requestTemplateJSON,
		&p.ResponseTextPath, &responseImagePath,
		&p.IsActive, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	// Unmarshal JSON fields
	if err := json.Unmarshal(defaultHeadersJSON, &p.DefaultHeaders); err != nil {
		p.DefaultHeaders = make(map[string]string)
	}
	if err := json.Unmarshal(requestTemplateJSON, &p.RequestTemplate); err != nil {
		p.RequestTemplate = make(map[string]interface{})
	}

	if imageEndpoint.Valid {
		p.ImageEndpoint = &imageEndpoint.String
	}
	if responseImagePath.Valid {
		p.ResponseImagePath = &responseImagePath.String
	}

	return &p, nil
}

// GetAll retrieves all providers based on user role
func (r *ProviderRepository) GetAll(ctx context.Context, userRole string) ([]provider.APIProvider, error) {
	query := `
		SELECT 
			id, name, display_name, description, base_url,
			auth_type, auth_header, auth_prefix,
			text_endpoint, image_endpoint,
			default_headers, request_template,
			response_text_path, response_image_path,
			is_active, created_at, updated_at
		FROM api_providers
	`

	// Apply role-based filtering
	if userRole != "admin" && userRole != "super_admin" {
		query += " WHERE is_active = true"
	}

	query += " ORDER BY name ASC"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query providers: %w", err)
	}
	defer rows.Close()

	var providers []provider.APIProvider

	for rows.Next() {
		var p provider.APIProvider
		var defaultHeadersJSON []byte
		var requestTemplateJSON []byte
		var imageEndpoint sql.NullString
		var responseImagePath sql.NullString

		err := rows.Scan(
			&p.ID, &p.Name, &p.DisplayName, &p.Description, &p.BaseURL,
			&p.AuthType, &p.AuthHeader, &p.AuthPrefix,
			&p.TextEndpoint, &imageEndpoint,
			&defaultHeadersJSON, &requestTemplateJSON,
			&p.ResponseTextPath, &responseImagePath,
			&p.IsActive, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			continue
		}

		// Unmarshal JSON fields
		json.Unmarshal(defaultHeadersJSON, &p.DefaultHeaders)
		json.Unmarshal(requestTemplateJSON, &p.RequestTemplate)

		if imageEndpoint.Valid {
			p.ImageEndpoint = &imageEndpoint.String
		}
		if responseImagePath.Valid {
			p.ResponseImagePath = &responseImagePath.String
		}

		if p.DefaultHeaders == nil {
			p.DefaultHeaders = make(map[string]string)
		}
		if p.RequestTemplate == nil {
			p.RequestTemplate = make(map[string]interface{})
		}

		providers = append(providers, p)
	}

	return providers, nil
}

// GetActiveProviders retrieves only active providers
func (r *ProviderRepository) GetActiveProviders(ctx context.Context) ([]provider.APIProvider, error) {
	return r.GetAll(ctx, "viewer") // viewer role only sees active
}

// Update updates a provider
func (r *ProviderRepository) Update(ctx context.Context, tx *sql.Tx, id string, updates map[string]interface{}) error {
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
		"name":              "name",
		"displayName":       "display_name",
		"description":       "description",
		"baseUrl":           "base_url",
		"authType":          "auth_type",
		"authHeader":        "auth_header",
		"authPrefix":        "auth_prefix",
		"textEndpoint":      "text_endpoint",
		"imageEndpoint":     "image_endpoint",
		"defaultHeaders":    "default_headers",
		"requestTemplate":   "request_template",
		"responseTextPath":  "response_text_path",
		"responseImagePath": "response_image_path",
		"isActive":          "is_active",
	}

	for key, value := range updates {
		if dbField, ok := fieldMap[key]; ok && value != nil {
			// Handle JSON fields
			if key == "defaultHeaders" || key == "requestTemplate" {
				jsonValue, err := json.Marshal(value)
				if err != nil {
					return fmt.Errorf("failed to marshal %s: %w", key, err)
				}
				setClauses = append(setClauses, fmt.Sprintf("%s = $%d", dbField, argIndex))
				args = append(args, jsonValue)
			} else {
				setClauses = append(setClauses, fmt.Sprintf("%s = $%d", dbField, argIndex))
				args = append(args, value)
			}
			argIndex++
		}
	}

	if len(setClauses) == 0 {
		return nil
	}

	// Always update timestamp
	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	// Add WHERE clause
	args = append(args, id)
	query := fmt.Sprintf("UPDATE api_providers SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "), argIndex)

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update provider: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("provider with id %s not found", id)
	}

	return nil
}

// Delete deletes a provider
func (r *ProviderRepository) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	if tx == nil {
		return fmt.Errorf("transaction is required for Delete")
	}

	query := `DELETE FROM api_providers WHERE id = $1`

	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete provider: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("provider with id %s not found", id)
	}

	return nil
}

// CheckProviderUsage checks if provider is used by any models
func (r *ProviderRepository) CheckProviderUsage(ctx context.Context, providerID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM ai_models WHERE provider_id = $1
	`, providerID).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to check provider usage: %w", err)
	}

	return count, nil
}

// ExistsByName checks if a provider with given name exists
func (r *ProviderRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS(SELECT 1 FROM api_providers WHERE name = $1)
	`, name).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check provider existence: %w", err)
	}

	return exists, nil
}
