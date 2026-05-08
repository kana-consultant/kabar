// internal/infrastructure/repository/aimodel/repository.go
package repositories

import (
	"context"
	"database/sql"
	"log"

	aimodel "seo-backend/internal/domain/model"
)

type ModelRepositoryImpl struct {
	db *sql.DB
}

func ModelRepository(db *sql.DB) aimodel.Repository {
	return &ModelRepositoryImpl{db: db}
}

func (r *ModelRepositoryImpl) GetAll(ctx context.Context) ([]aimodel.AIModel, error) {
	query := `
		SELECT id, name, display_name, provider_id, is_active
		FROM ai_models 
		WHERE is_active = true
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []aimodel.AIModel
	for rows.Next() {
		var m aimodel.AIModel
		err := rows.Scan(&m.ID, &m.Name, &m.DisplayName, &m.ProviderID, &m.IsActive)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		models = append(models, m)
	}

	return models, nil
}

func (r *ModelRepositoryImpl) GetAllWithStatus(ctx context.Context, userRole string) ([]aimodel.ModelWithStatus, error) {
	query := `
		SELECT DISTINCT
			m.id, 
			m.name, 
			m.provider_id, 
			m.display_name
		FROM ai_models m
		INNER JOIN api_keys ak ON ak.model_id = m.id
		WHERE ak.is_active = true AND ak.service = 'text'
	`

	// Filter berdasarkan role
	if userRole != "super_admin" {
		query += " AND m.is_active = true"
	}

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("Query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var models []aimodel.ModelWithStatus
	for rows.Next() {
		var m aimodel.ModelWithStatus
		err := rows.Scan(&m.ID, &m.Name, &m.ProviderID, &m.DisplayName)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		models = append(models, m)
	}

	log.Printf("Found %d models", len(models))
	return models, nil
}

func (r *ModelRepositoryImpl) GetByID(ctx context.Context, id string) (*aimodel.AIModel, error) {
	query := `
		SELECT id, name, display_name, provider_id, is_active
		FROM ai_models 
		WHERE id = $1
	`

	var m aimodel.AIModel
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.Name, &m.DisplayName, &m.ProviderID, &m.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &m, nil
}

func (r *ModelRepositoryImpl) GetByProvider(ctx context.Context, providerID string) ([]aimodel.AIModel, error) {
	query := `
		SELECT id, name, display_name, provider_id, is_active
		FROM ai_models 
		WHERE provider_id = $1 AND is_active = true
	`

	rows, err := r.db.QueryContext(ctx, query, providerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []aimodel.AIModel
	for rows.Next() {
		var m aimodel.AIModel
		err := rows.Scan(&m.ID, &m.Name, &m.DisplayName, &m.ProviderID, &m.IsActive)
		if err != nil {
			continue
		}
		models = append(models, m)
	}

	return models, nil
}
