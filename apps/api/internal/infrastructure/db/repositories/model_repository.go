package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"seo-backend/internal/domain/ai_model"
	model_family "seo-backend/internal/domain/modelfamily"
	"seo-backend/internal/domain/paginate"
	userRole "seo-backend/internal/helper/filter"
	"seo-backend/internal/models"

	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type ModelRepositoryImpl struct {
	db          *sql.DB
	redisClient *redis.Client
}

func ModelRepository(db *sql.DB, redisClient *redis.Client) ai_model.Repository {
	return &ModelRepositoryImpl{
		db:          db,
		redisClient: redisClient,
	}
}

// invalidateCache menghapus semua cache yang berkaitan dengan team_id tertentu
func (r *ModelRepositoryImpl) invalidateCache(ctx context.Context, teamID string) {
	if teamID == "" {
		return
	}

	// Pattern cache keys untuk team ini
	patterns := []string{
		fmt.Sprintf("ai_models_all:*:%s:*", teamID),
		fmt.Sprintf("ai_models_with_status:*:%s:*", teamID),
	}

	for _, pattern := range patterns {
		iter := r.redisClient.Scan(ctx, 0, pattern, 0).Iterator()
		for iter.Next(ctx) {
			key := iter.Val()
			if err := r.redisClient.Del(ctx, key).Err(); err != nil {
				log.Printf("[Cache] Failed to delete key %s: %v", key, err)
			} else {
				log.Printf("[Cache] Deleted key: %s", key)
			}
		}
		if err := iter.Err(); err != nil {
			log.Printf("[Cache] Scan error for pattern %s: %v", pattern, err)
		}
	}
}

// invalidateAllCache menghapus semua cache ai_models
func (r *ModelRepositoryImpl) invalidateAllCache(ctx context.Context) {
	patterns := []string{
		"ai_models_all:*",
		"ai_models_with_status:*",
	}

	for _, pattern := range patterns {
		iter := r.redisClient.Scan(ctx, 0, pattern, 0).Iterator()
		for iter.Next(ctx) {
			key := iter.Val()
			if err := r.redisClient.Del(ctx, key).Err(); err != nil {
				log.Printf("[Cache] Failed to delete key %s: %v", key, err)
			} else {
				log.Printf("[Cache] Deleted key: %s", key)
			}
		}
		if err := iter.Err(); err != nil {
			log.Printf("[Cache] Scan error for pattern %s: %v", pattern, err)
		}
	}
}

// Create inserts a new AI model
func (r *ModelRepositoryImpl) Create(ctx context.Context, model *ai_model.AIModel) error {
	query := `
        INSERT INTO ai_models (
            id, family_id, provider_id, team_id, display_name,
            description, system_prompt, max_tokens, temperature, context_window,
            is_active, is_default, created_by, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
    `

	if model.ID == "" {
		model.ID = uuid.New().String()
	}

	now := time.Now()

	isActive := true
	if model.IsActive != nil {
		isActive = *model.IsActive
	}

	isDefault := false
	if model.IsDefault != nil {
		isDefault = *model.IsDefault
	}

	maxTokens := 4096
	if model.MaxTokens != nil {
		maxTokens = *model.MaxTokens
	}

	temperature := 0.7
	if model.Temperature != nil {
		temperature = *model.Temperature
	}

	var teamID *string
	if model.TeamID != nil && *model.TeamID != "" {
		teamID = model.TeamID
	}

	displayName := ""
	if model.DisplayName != nil && *model.DisplayName != "" {
		displayName = *model.DisplayName
	}

	description := ""
	if model.Description != nil {
		description = *model.Description
	}

	systemPrompt := ""
	if model.SystemPrompt != nil {
		systemPrompt = *model.SystemPrompt
	}

	contextWindow := 0
	if model.ContextWindow != nil {
		contextWindow = *model.ContextWindow
	}

	createdBy := ""
	if model.CreatedBy != nil {
		createdBy = *model.CreatedBy
	}

	_, err := r.db.ExecContext(ctx, query,
		model.ID,
		model.FamilyID,
		model.ProviderID,
		teamID,
		displayName,
		description,
		systemPrompt,
		maxTokens,
		temperature,
		contextWindow,
		isActive,
		isDefault,
		createdBy,
		now,
		now,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") ||
			strings.Contains(err.Error(), "UNIQUE") ||
			strings.Contains(err.Error(), "23505") {
			return ai_model.ErrDuplicate
		}
		return fmt.Errorf("%w: %v", ai_model.ErrDatabase, err)
	}

	model.CreatedAt = now
	model.UpdatedAt = now
	model.IsActive = &isActive
	model.IsDefault = &isDefault
	model.MaxTokens = &maxTokens
	model.Temperature = &temperature

	// Invalidate cache berdasarkan team_id
	if teamID != nil && *teamID != "" {
		go r.invalidateCache(context.Background(), *teamID)
	} else {
		go r.invalidateAllCache(context.Background())
	}

	return nil
}

// GetByID retrieves an AI model by ID
func (r *ModelRepositoryImpl) GetByID(ctx context.Context, id string) (*ai_model.AIModel, error) {
	query := `
        SELECT id, family_id, provider_id, team_id, display_name,
               description, system_prompt, max_tokens, temperature, context_window,
               is_active, is_default, created_by, created_at, updated_at
        FROM ai_models 
        WHERE id = $1
    `

	var m ai_model.AIModel
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID,
		&m.FamilyID,
		&m.ProviderID,
		&m.TeamID,
		&m.DisplayName,
		&m.Description,
		&m.SystemPrompt,
		&m.MaxTokens,
		&m.Temperature,
		&m.ContextWindow,
		&m.IsActive,
		&m.IsDefault,
		&m.CreatedBy,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ai_model.ErrNotFound
		}
		return nil, err
	}

	return &m, nil
}

// GetSchemaByModelID retrieves schema from model_families
func (r *ModelRepositoryImpl) GetSchemaByModelID(ctx context.Context, modelID string) (*ai_model.RequestSchema, error) {
	query := `
        SELECT 
            rs.id, rs.provider_id, rs.endpoint_path,
            rs.request_template, rs.response_text_path, rs.response_image_path,
            rs.supports_temperature, rs.supports_streaming,
            rs.created_at, rs.updated_at
        FROM request_schemas rs
        JOIN model_families mf ON mf.schema_id = rs.id
        JOIN ai_models am ON am.family_id = mf.id
        WHERE am.id = $1
    `

	var schema ai_model.RequestSchema
	err := r.db.QueryRowContext(ctx, query, modelID).Scan(
		&schema.ID,
		&schema.ProviderID,
		&schema.EndpointPath,
		&schema.RequestTemplate,
		&schema.ResponseTextPath,
		&schema.ResponseImagePath,
		&schema.SupportsTemperature,
		&schema.SupportsStreaming,
		&schema.CreatedAt,
		&schema.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ai_model.ErrNotFound
		}
		return nil, err
	}

	return &schema, nil
}

// GetAll retrieves all AI models with pagination and user context
func (r *ModelRepositoryImpl) GetAll(ctx context.Context, userCtx models.UserContext, params paginate.PaginationParams) ([]ai_model.AIModel, error) {
	whereClause, whereArgs := userRole.BuildAccessFilter(userCtx)

	cacheKey := fmt.Sprintf("ai_models_all:%s:%s:search=%s", userCtx.GetRole(), userCtx.GetTeamID(), params.Search)

	cached, err := r.redisClient.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var models []ai_model.AIModel
		if err := json.Unmarshal(cached, &models); err == nil {
			log.Printf("[GetAll] cache hit | key=%s", cacheKey)
			return models, nil
		}
	}

	log.Printf("[GetAll] cache miss | key=%s", cacheKey)

	args := append([]any{}, whereArgs...)
	searchClause := ""
	if params.Search != "" {
		args = append(args, "%"+params.Search+"%")
		searchClause = fmt.Sprintf(" AND display_name ILIKE $%d", len(args))
	}

	fullWhere := whereClause + searchClause

	query := fmt.Sprintf(`
		SELECT id, family_id, provider_id, team_id, display_name,
		       description, system_prompt, max_tokens, temperature, context_window,
		       is_active, is_default, created_by, created_at, updated_at
		FROM ai_models
		WHERE %s
		ORDER BY created_at DESC
	`, fullWhere)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []ai_model.AIModel
	for rows.Next() {
		var m ai_model.AIModel
		err := rows.Scan(
			&m.ID,
			&m.FamilyID,
			&m.ProviderID,
			&m.TeamID,
			&m.DisplayName,
			&m.Description,
			&m.SystemPrompt,
			&m.MaxTokens,
			&m.Temperature,
			&m.ContextWindow,
			&m.IsActive,
			&m.IsDefault,
			&m.CreatedBy,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		models = append(models, m)
	}

	modelsBytes, err := json.Marshal(models)
	if err != nil {
		log.Printf("[GetAll] failed to marshal models | err=%v", err)
	} else {
		if err := r.redisClient.Set(ctx, cacheKey, modelsBytes, 2*time.Minute).Err(); err != nil {
			log.Printf("[GetAll] failed to set cache | err=%v", err)
		}
	}

	return models, nil
}

// GetByFamily retrieves all AI models in a family with pagination
func (r *ModelRepositoryImpl) GetByFamily(ctx context.Context, familyID string, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[ai_model.AIModel], error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}

	whereClause, whereArgs := userRole.BuildAccessFilter(userCtx)
	whereClause += " AND family_id = $" + fmt.Sprintf("%d", len(whereArgs)+1)
	whereArgs = append(whereArgs, familyID)

	args := append([]any{}, whereArgs...)
	searchClause := ""
	if params.Search != "" {
		args = append(args, "%"+params.Search+"%")
		searchClause = fmt.Sprintf(" AND display_name ILIKE $%d", len(args))
	}

	fullWhere := whereClause + searchClause

	var totalItems int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM ai_models WHERE %s", fullWhere)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(params.Limit)))
	currentPage := (params.Offset / params.Limit) + 1

	args = append(args, params.Limit, params.Offset)
	query := fmt.Sprintf(`
		SELECT id, family_id, provider_id, team_id, display_name,
		       description, system_prompt, max_tokens, temperature, context_window,
		       is_active, is_default, created_by, created_at, updated_at
		FROM ai_models
		WHERE %s
		ORDER BY display_name ASC
		LIMIT $%d OFFSET $%d
	`, fullWhere, len(args)-1, len(args))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []ai_model.AIModel
	for rows.Next() {
		var m ai_model.AIModel
		err := rows.Scan(
			&m.ID,
			&m.FamilyID,
			&m.ProviderID,
			&m.TeamID,
			&m.DisplayName,
			&m.Description,
			&m.SystemPrompt,
			&m.MaxTokens,
			&m.Temperature,
			&m.ContextWindow,
			&m.IsActive,
			&m.IsDefault,
			&m.CreatedBy,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		models = append(models, m)
	}

	return &paginate.PaginatedResult[ai_model.AIModel]{
		Data:        models,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		Limit:       params.Limit,
		Offset:      params.Offset,
	}, nil
}

// GetByProvider retrieves all AI models for a provider with pagination
func (r *ModelRepositoryImpl) GetByProvider(ctx context.Context, providerID string, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[ai_model.AIModel], error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}

	whereClause, whereArgs := userRole.BuildAccessFilter(userCtx)
	whereClause += " AND provider_id = $" + fmt.Sprintf("%d", len(whereArgs)+1)
	whereArgs = append(whereArgs, providerID)

	args := append([]any{}, whereArgs...)
	searchClause := ""
	if params.Search != "" {
		args = append(args, "%"+params.Search+"%")
		searchClause = fmt.Sprintf(" AND display_name ILIKE $%d", len(args))
	}

	fullWhere := whereClause + searchClause

	var totalItems int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM ai_models WHERE %s", fullWhere)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(params.Limit)))
	currentPage := (params.Offset / params.Limit) + 1

	args = append(args, params.Limit, params.Offset)
	query := fmt.Sprintf(`
		SELECT id, family_id, provider_id, team_id, display_name,
		       description, system_prompt, max_tokens, temperature, context_window,
		       is_active, is_default, created_by, created_at, updated_at
		FROM ai_models
		WHERE %s
		ORDER BY display_name ASC
		LIMIT $%d OFFSET $%d
	`, fullWhere, len(args)-1, len(args))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []ai_model.AIModel
	for rows.Next() {
		var m ai_model.AIModel
		err := rows.Scan(
			&m.ID,
			&m.FamilyID,
			&m.ProviderID,
			&m.TeamID,
			&m.DisplayName,
			&m.Description,
			&m.SystemPrompt,
			&m.MaxTokens,
			&m.Temperature,
			&m.ContextWindow,
			&m.IsActive,
			&m.IsDefault,
			&m.CreatedBy,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		models = append(models, m)
	}

	return &paginate.PaginatedResult[ai_model.AIModel]{
		Data:        models,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		Limit:       params.Limit,
		Offset:      params.Offset,
	}, nil
}

// GetByTeam retrieves all AI models for a team with pagination
func (r *ModelRepositoryImpl) GetByTeam(ctx context.Context, teamID string, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[ai_model.AIModel], error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}

	whereClause, whereArgs := userRole.BuildAccessFilter(userCtx)
	whereClause += " AND team_id = $" + fmt.Sprintf("%d", len(whereArgs)+1)
	whereArgs = append(whereArgs, teamID)

	args := append([]any{}, whereArgs...)
	searchClause := ""
	if params.Search != "" {
		args = append(args, "%"+params.Search+"%")
		searchClause = fmt.Sprintf(" AND display_name ILIKE $%d", len(args))
	}

	fullWhere := whereClause + searchClause

	var totalItems int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM ai_models WHERE %s", fullWhere)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(params.Limit)))
	currentPage := (params.Offset / params.Limit) + 1

	args = append(args, params.Limit, params.Offset)
	query := fmt.Sprintf(`
		SELECT id, family_id, provider_id, team_id, display_name,
		       description, system_prompt, max_tokens, temperature, context_window,
		       is_active, is_default, created_by, created_at, updated_at
		FROM ai_models
		WHERE %s
		ORDER BY display_name ASC
		LIMIT $%d OFFSET $%d
	`, fullWhere, len(args)-1, len(args))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []ai_model.AIModel
	for rows.Next() {
		var m ai_model.AIModel
		err := rows.Scan(
			&m.ID,
			&m.FamilyID,
			&m.ProviderID,
			&m.TeamID,
			&m.DisplayName,
			&m.Description,
			&m.SystemPrompt,
			&m.MaxTokens,
			&m.Temperature,
			&m.ContextWindow,
			&m.IsActive,
			&m.IsDefault,
			&m.CreatedBy,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		models = append(models, m)
	}

	return &paginate.PaginatedResult[ai_model.AIModel]{
		Data:        models,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		Limit:       params.Limit,
		Offset:      params.Offset,
	}, nil
}

// GetDefault retrieves default AI models
func (r *ModelRepositoryImpl) GetDefault(ctx context.Context, userCtx models.UserContext) ([]ai_model.AIModel, error) {
	query := `
        SELECT id, family_id, provider_id, team_id, display_name,
               description, system_prompt, max_tokens, temperature, context_window,
               is_active, is_default, created_by, created_at, updated_at
        FROM ai_models
        WHERE is_default = true AND is_active = true
        ORDER BY display_name ASC
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []ai_model.AIModel
	for rows.Next() {
		var m ai_model.AIModel
		err := rows.Scan(
			&m.ID,
			&m.FamilyID,
			&m.ProviderID,
			&m.TeamID,
			&m.DisplayName,
			&m.Description,
			&m.SystemPrompt,
			&m.MaxTokens,
			&m.Temperature,
			&m.ContextWindow,
			&m.IsActive,
			&m.IsDefault,
			&m.CreatedBy,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		models = append(models, m)
	}

	return models, nil
}

// GetAllWithStatus retrieves models with status based on user role
func (r *ModelRepositoryImpl) GetAllWithStatus(ctx context.Context, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[ai_model.ModelWithStatus], error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}

	cacheKey := fmt.Sprintf("ai_models_with_status:%s:%s:limit=%d:offset=%d:search=%s", userCtx.GetRole(), userCtx.GetTeamID(), params.Limit, params.Offset, params.Search)

	cached, err := r.redisClient.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var result paginate.PaginatedResult[ai_model.ModelWithStatus]
		if err := json.Unmarshal(cached, &result); err == nil {
			log.Printf("[GetAllWithStatus] cache hit | key=%s", cacheKey)
			return &result, nil
		}
	}

	log.Printf("[GetAllWithStatus] cache miss | key=%s", cacheKey)

	query := `
        SELECT DISTINCT
            m.id, 
            m.provider_id, 
            m.display_name,
            m.is_active,
            m.is_default
        FROM ai_models m
        INNER JOIN api_keys ak ON ak.model_id = m.id
        WHERE ak.is_active = true AND ak.service = 'text'
    `

	args := []interface{}{}

	if params.Search != "" {
		query += " AND m.display_name ILIKE $" + fmt.Sprintf("%d", len(args)+1)
		args = append(args, "%"+params.Search+"%")
	}

	if userCtx.GetRole() != "superadmin" {
		query += " AND m.is_active = true"
	}

	var totalItems int
	countQuery := "SELECT COUNT(*) FROM (" + query + ") AS subq"
	err = r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalItems)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(params.Limit)))
	currentPage := (params.Offset / params.Limit) + 1

	query += " ORDER BY m.display_name ASC"
	query += " LIMIT $" + fmt.Sprintf("%d", len(args)+1) + " OFFSET $" + fmt.Sprintf("%d", len(args)+2)
	args = append(args, params.Limit, params.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("Query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var models []ai_model.ModelWithStatus
	for rows.Next() {
		var m ai_model.ModelWithStatus
		err := rows.Scan(&m.ID, &m.ProviderID, &m.DisplayName, &m.IsActive, &m.IsDefault)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		models = append(models, m)
	}

	result := &paginate.PaginatedResult[ai_model.ModelWithStatus]{
		Data:        models,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		Limit:       params.Limit,
		Offset:      params.Offset,
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		log.Printf("[GetAllWithStatus] failed to marshal result | err=%v", err)
	} else {
		if err := r.redisClient.Set(ctx, cacheKey, resultBytes, 2*time.Minute).Err(); err != nil {
			log.Printf("[GetAllWithStatus] failed to set cache | err=%v", err)
		}
	}

	return result, nil
}

// Update updates an existing AI model
func (r *ModelRepositoryImpl) Update(ctx context.Context, model *ai_model.AIModel) error {
	query := `
        UPDATE ai_models
        SET family_id = $1, provider_id = $2, team_id = $3,
            display_name = $4, description = $5, system_prompt = $6,
            max_tokens = $7, temperature = $8, context_window = $9,
            is_active = $10, is_default = $11, created_by = $12, updated_at = $13
        WHERE id = $14
    `

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		model.FamilyID,
		model.ProviderID,
		model.TeamID,
		model.DisplayName,
		model.Description,
		model.SystemPrompt,
		model.MaxTokens,
		model.Temperature,
		model.ContextWindow,
		model.IsActive,
		model.IsDefault,
		model.CreatedBy,
		now,
		model.ID,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") ||
			strings.Contains(err.Error(), "UNIQUE") ||
			strings.Contains(err.Error(), "23505") {
			return ai_model.ErrDuplicate
		}
		return fmt.Errorf("%w: %v", ai_model.ErrDatabase, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ai_model.ErrNotFound
	}

	model.UpdatedAt = now

	// Invalidate cache berdasarkan team_id
	if model.TeamID != nil && *model.TeamID != "" {
		go r.invalidateCache(context.Background(), *model.TeamID)
	} else {
		go r.invalidateAllCache(context.Background())
	}

	return nil
}

// Delete soft deletes an AI model
func (r *ModelRepositoryImpl) Delete(ctx context.Context, id string) error {
	// Dapatkan team_id sebelum delete
	var teamID sql.NullString
	getTeamQuery := `SELECT team_id FROM ai_models WHERE id = $1`
	err := r.db.QueryRowContext(ctx, getTeamQuery, id).Scan(&teamID)

	query := `UPDATE ai_models SET is_active = false, updated_at = $1 WHERE id = $2`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ai_model.ErrNotFound
	}

	// Invalidate cache berdasarkan team_id
	if teamID.Valid && teamID.String != "" {
		go r.invalidateCache(context.Background(), teamID.String)
	} else {
		go r.invalidateAllCache(context.Background())
	}

	return nil
}

// HardDelete permanently deletes an AI model
func (r *ModelRepositoryImpl) HardDelete(ctx context.Context, id string) error {
	// Dapatkan team_id sebelum delete
	var teamID sql.NullString
	getTeamQuery := `SELECT team_id FROM ai_models WHERE id = $1`
	err := r.db.QueryRowContext(ctx, getTeamQuery, id).Scan(&teamID)

	query := `DELETE FROM ai_models WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ai_model.ErrNotFound
	}

	// Invalidate cache berdasarkan team_id
	if teamID.Valid && teamID.String != "" {
		go r.invalidateCache(context.Background(), teamID.String)
	} else {
		go r.invalidateAllCache(context.Background())
	}

	return nil
}

// Exists checks if an AI model exists by display_name
func (r *ModelRepositoryImpl) Exists(ctx context.Context, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM ai_models WHERE display_name = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, name).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// ExistsByFamilyAndName checks if a model exists in a family
func (r *ModelRepositoryImpl) ExistsByFamilyAndName(ctx context.Context, familyID, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM ai_models WHERE family_id = $1 AND display_name = $2)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, familyID, name).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// Count returns total number of AI models
func (r *ModelRepositoryImpl) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM ai_models`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CountByFamily returns number of models in a family
func (r *ModelRepositoryImpl) CountByFamily(ctx context.Context, familyID string) (int64, error) {
	query := `SELECT COUNT(*) FROM ai_models WHERE family_id = $1`

	var count int64
	err := r.db.QueryRowContext(ctx, query, familyID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// SetDefaultForProvider sets a model as default for its provider
func (r *ModelRepositoryImpl) SetDefaultForProvider(ctx context.Context, providerID, modelID string) error {
	resetQuery := `
        UPDATE ai_models
        SET is_default = false
        WHERE provider_id = $1 AND id != $2
    `
	_, err := r.db.ExecContext(ctx, resetQuery, providerID, modelID)
	if err != nil {
		return err
	}

	updateQuery := `UPDATE ai_models SET is_default = true WHERE id = $1`
	_, err = r.db.ExecContext(ctx, updateQuery, modelID)
	if err != nil {
		return err
	}

	// Invalidate all cache karena default model berubah
	go r.invalidateAllCache(context.Background())

	return nil
}

// GetByIDWithSchema retrieves model with its schema from family
func (r *ModelRepositoryImpl) GetByIDWithSchema(ctx context.Context, id string) (*model_family.ModelFamilyWithSchema, error) {
	query := `
        SELECT 
            mf.id, 
            mf.provider_id, 
            mf.schema_id, 
            mf.name, 
            mf.display_name, 
            mf.description, 
            mf.created_at, 
            mf.updated_at,
            rs.id, 
            rs.provider_id, 
            rs.endpoint_path,
            rs.request_template, 
            rs.response_text_path, 
            rs.response_image_path,
            rs.supports_temperature, 
            rs.supports_streaming,
            rs.created_at, 
            rs.updated_at
        FROM model_families mf
        LEFT JOIN request_schemas rs ON mf.schema_id = rs.id
        WHERE mf.id = $1
    `

	var result model_family.ModelFamilyWithSchema
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&result.ID,
		&result.ProviderID,
		&result.SchemaID,
		&result.Name,
		&result.DisplayName,
		&result.Description,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.Schema.ID,
		&result.Schema.ProviderID,
		&result.Schema.EndpointPath,
		&result.Schema.RequestTemplate,
		&result.Schema.ResponseTextPath,
		&result.Schema.ResponseImagePath,
		&result.Schema.SupportsTemperature,
		&result.Schema.SupportsStreaming,
		&result.Schema.CreatedAt,
		&result.Schema.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model_family.ErrNotFound
		}
		return nil, err
	}

	return &result, nil
}
