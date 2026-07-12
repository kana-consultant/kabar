package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	model_family "seo-backend/internal/domain/modelfamily"
	"seo-backend/internal/domain/paginate"
	userRole "seo-backend/internal/helper/filter"
	"seo-backend/internal/models"

	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type repositoryImpl struct {
	db          *sql.DB
	redisClient *redis.Client
}

// NewRepository creates a new ModelFamily repository implementation
func NewModelFamiliesRepository(db *sql.DB, redisClient *redis.Client) model_family.Repository {
	return &repositoryImpl{
		db:          db,
		redisClient: redisClient,
	}
}

// Create creates a new model family (single item)
func (r *repositoryImpl) Create(ctx context.Context, mf *model_family.ModelFamily) error {
	query := `
        INSERT INTO model_families (
            id, provider_id, schema_id, name, display_name, 
            description, max_tokens, temperature, system_prompt,
            created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `

	id := uuid.New()
	now := time.Now()

	_, err := r.db.ExecContext(ctx, query,
		id,
		mf.ProviderID,
		mf.SchemaID,
		mf.Name,
		mf.DisplayName,
		mf.Description,
		mf.MaxTokens,
		mf.Temperature,
		mf.SystemPrompt,
		now,
		now,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") ||
			strings.Contains(err.Error(), "UNIQUE") ||
			strings.Contains(err.Error(), "23505") { // PostgreSQL unique violation code
			return model_family.ErrDuplicate
		}
		return fmt.Errorf("%w: %v", model_family.ErrDatabase, err)
	}

	mf.ID = id.String()
	mf.CreatedAt = now
	mf.UpdatedAt = now

	return nil
}

// CreateBatch creates multiple model families at once
func (r *repositoryImpl) CreateBatchWithTx(
	ctx context.Context,
	tx *sql.Tx,
	families []model_family.ModelFamilyWithSchema,
) error {

	query := `
		INSERT INTO model_families (
			id,
			provider_id,
			schema_id,
			name,
			display_name,
			description,
			max_tokens,
			temperature,
			system_prompt,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	now := time.Now()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("%w: %v", model_family.ErrDatabase, err)
	}
	defer stmt.Close()

	for i := range families {
		id := uuid.New()

		_, err := stmt.ExecContext(
			ctx,
			id,
			families[i].ProviderID,
			families[i].SchemaID,
			families[i].Name,
			families[i].DisplayName,
			families[i].Description,
			families[i].MaxTokens,
			families[i].Temperature,
			families[i].SystemPrompt,
			now,
			now,
		)

		if err != nil {
			if strings.Contains(err.Error(), "duplicate") ||
				strings.Contains(err.Error(), "UNIQUE") ||
				strings.Contains(err.Error(), "23505") {
				return model_family.ErrDuplicate
			}

			return fmt.Errorf("%w: %v", model_family.ErrDatabase, err)
		}

		families[i].ID = id.String()
		families[i].CreatedAt = now
		families[i].UpdatedAt = now
	}

	return nil
}

func (r *repositoryImpl) GetByID(ctx context.Context, id string) (*model_family.ModelFamily, error) {
	query := `
        SELECT id, provider_id, schema_id, name, display_name, 
               description, max_tokens, temperature, system_prompt,
               created_at, updated_at
        FROM model_families
        WHERE id = $1
    `

	var mf model_family.ModelFamily
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&mf.ID,
		&mf.ProviderID,
		&mf.SchemaID,
		&mf.Name,
		&mf.DisplayName,
		&mf.Description,
		&mf.MaxTokens,
		&mf.Temperature,
		&mf.SystemPrompt,
		&mf.CreatedAt,
		&mf.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model_family.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", model_family.ErrDatabase, err)
	}

	return &mf, nil
}

func (r *repositoryImpl) GetByProviderAndName(ctx context.Context, providerID string, name string) (*model_family.ModelFamily, error) {
	query := `
        SELECT id, provider_id, schema_id, name, display_name, 
               description, max_tokens, temperature, system_prompt,
               created_at, updated_at
        FROM model_families
        WHERE provider_id = $1 AND name = $2
    `

	var mf model_family.ModelFamily
	err := r.db.QueryRowContext(ctx, query, providerID, name).Scan(
		&mf.ID,
		&mf.ProviderID,
		&mf.SchemaID,
		&mf.Name,
		&mf.DisplayName,
		&mf.Description,
		&mf.MaxTokens,
		&mf.Temperature,
		&mf.SystemPrompt,
		&mf.CreatedAt,
		&mf.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model_family.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", model_family.ErrDatabase, err)
	}

	return &mf, nil
}

func (r *repositoryImpl) GetAll(ctx context.Context, userCtx models.UserContext, params paginate.PaginationParams) ([]model_family.ModelFamilyWithProvider, error) {
	// Build access filter
	whereClause, whereArgs := userRole.BuildAccessFilter(userCtx)

	// Cache key
	cacheKey := fmt.Sprintf("model_families_list:%s:%s:search=%s", userCtx.GetRole(), userCtx.GetTeamID(), params.Search)

	// Cek cache
	cached, err := r.redisClient.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var families []model_family.ModelFamilyWithProvider
		if err := json.Unmarshal(cached, &families); err == nil {
			log.Printf("[GetAll] cache hit | key=%s", cacheKey)
			return families, nil
		}
	}

	log.Printf("[GetAll] cache miss | key=%s", cacheKey)

	// Build search clause
	args := append([]any{}, whereArgs...)
	searchClause := ""
	if params.Search != "" {
		args = append(args, "%"+params.Search+"%")
		searchClause = fmt.Sprintf(" AND (mf.name ILIKE $%d OR mf.display_name ILIKE $%d OR p.name ILIKE $%d OR p.display_name ILIKE $%d)",
			len(args), len(args), len(args), len(args))
	}

	// Combine WHERE
	fullWhere := whereClause + searchClause

	// Query dengan JOIN ke api_providers
	query := fmt.Sprintf(`
		SELECT 
			mf.id, 
			mf.provider_id, 
			mf.schema_id, 
			mf.name, 
			mf.display_name, 
			mf.description, 
			mf.max_tokens,
			mf.temperature,
			mf.system_prompt,
			mf.created_at, 
			mf.updated_at,
			COALESCE(p.name, '') as provider_name,
			COALESCE(p.display_name, '') as provider_display_name
		FROM model_families mf
		LEFT JOIN api_providers p ON mf.provider_id = p.id
		WHERE %s
		ORDER BY mf.created_at DESC
	`, fullWhere)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", model_family.ErrDatabase, err)
	}
	defer rows.Close()

	var families []model_family.ModelFamilyWithProvider
	for rows.Next() {
		var mf model_family.ModelFamilyWithProvider
		err := rows.Scan(
			&mf.ID,
			&mf.ProviderID,
			&mf.SchemaID,
			&mf.Name,
			&mf.DisplayName,
			&mf.Description,
			&mf.MaxTokens,
			&mf.Temperature,
			&mf.SystemPrompt,
			&mf.CreatedAt,
			&mf.UpdatedAt,
			&mf.ProviderName,
			&mf.ProviderDisplayName,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", model_family.ErrDatabase, err)
		}
		families = append(families, mf)
	}

	// Simpan ke cache
	resultBytes, err := json.Marshal(families)
	if err != nil {
		log.Printf("[GetAll] failed to marshal result | err=%v", err)
	} else {
		if err := r.redisClient.Set(ctx, cacheKey, resultBytes, 2*time.Minute).Err(); err != nil {
			log.Printf("[GetAll] failed to set cache | err=%v", err)
		} else {
			log.Printf("[GetAll] cache saved | key=%s | ttl=2m", cacheKey)
		}
	}

	return families, nil
}

func (r *repositoryImpl) isAdmin(role string) bool {
	return role == "admin" || role == "superadmin"
}

func (r *repositoryImpl) GetByProviderID(ctx context.Context, providerID string) ([]model_family.ModelFamilyWithSchema, error) {
	query := `
		SELECT
			mf.id,
			mf.provider_id,
			mf.schema_id,
			mf.name,
			mf.display_name,
			mf.description,
			mf.max_tokens,
			mf.temperature,
			mf.system_prompt,
			mf.created_at,
			mf.updated_at,

			rs.id,
			rs.provider_id,
			rs.name,
			rs.endpoint_path,
			rs.max_tokens_key,
			rs.system_role_key,
			rs.response_text_path,
			rs.response_image_path,
			rs.request_template,
			rs.supports_temperature,
			rs.supports_streaming,
			rs.created_at,
			rs.updated_at
		FROM model_families mf
		INNER JOIN request_schemas rs
			ON rs.id = mf.schema_id
		WHERE mf.provider_id = $1
		ORDER BY mf.name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, providerID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", model_family.ErrDatabase, err)
	}
	defer rows.Close()

	var result []model_family.ModelFamilyWithSchema

	for rows.Next() {
		var item model_family.ModelFamilyWithSchema

		err := rows.Scan(
			&item.ID,
			&item.ProviderID,
			&item.SchemaID,
			&item.Name,
			&item.DisplayName,
			&item.Description,
			&item.MaxTokens,
			&item.Temperature,
			&item.SystemPrompt,
			&item.CreatedAt,
			&item.UpdatedAt,

			&item.Schema.ID,
			&item.Schema.ProviderID,
			&item.Schema.Name,
			&item.Schema.EndpointPath,
			&item.Schema.MaxTokensKey,
			&item.Schema.SystemRoleKey,
			&item.Schema.ResponseTextPath,
			&item.Schema.ResponseImagePath,
			&item.Schema.RequestTemplate,
			&item.Schema.SupportsTemperature,
			&item.Schema.SupportsStreaming,
			&item.Schema.CreatedAt,
			&item.Schema.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", model_family.ErrDatabase, err)
		}

		result = append(result, item)
	}

	return result, nil
}

func (r *repositoryImpl) GetBySchema(ctx context.Context, schemaID string) ([]model_family.ModelFamily, error) {
	query := `
        SELECT id, provider_id, schema_id, name, display_name, 
               description, max_tokens, temperature, system_prompt,
               created_at, updated_at
        FROM model_families
        WHERE schema_id = $1
        ORDER BY name ASC
    `

	rows, err := r.db.QueryContext(ctx, query, schemaID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", model_family.ErrDatabase, err)
	}
	defer rows.Close()

	var families []model_family.ModelFamily
	for rows.Next() {
		var mf model_family.ModelFamily
		err := rows.Scan(
			&mf.ID,
			&mf.ProviderID,
			&mf.SchemaID,
			&mf.Name,
			&mf.DisplayName,
			&mf.Description,
			&mf.MaxTokens,
			&mf.Temperature,
			&mf.SystemPrompt,
			&mf.CreatedAt,
			&mf.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", model_family.ErrDatabase, err)
		}
		families = append(families, mf)
	}

	return families, nil
}

func (r *repositoryImpl) Update(ctx context.Context, mf *model_family.ModelFamily) error {
	query := `
        UPDATE model_families
        SET schema_id = $1, name = $2, display_name = $3, 
            description = $4, max_tokens = $5, temperature = $6, 
            system_prompt = $7, updated_at = $8
        WHERE id = $9
    `

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		mf.SchemaID,
		mf.Name,
		mf.DisplayName,
		mf.Description,
		mf.MaxTokens,
		mf.Temperature,
		mf.SystemPrompt,
		now,
		mf.ID,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") ||
			strings.Contains(err.Error(), "UNIQUE") ||
			strings.Contains(err.Error(), "23505") {
			return model_family.ErrDuplicate
		}
		return fmt.Errorf("%w: %v", model_family.ErrDatabase, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", model_family.ErrDatabase, err)
	}

	if rowsAffected == 0 {
		return model_family.ErrNotFound
	}

	mf.UpdatedAt = now
	return nil
}

func (r *repositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM model_families WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", model_family.ErrDatabase, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", model_family.ErrDatabase, err)
	}

	if rowsAffected == 0 {
		return model_family.ErrNotFound
	}

	return nil
}

func (r *repositoryImpl) Exists(ctx context.Context, providerID string, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM model_families WHERE provider_id = $1 AND name = $2)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, providerID, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%w: %v", model_family.ErrDatabase, err)
	}

	return exists, nil
}

func (r *repositoryImpl) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM model_families`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", model_family.ErrDatabase, err)
	}

	return count, nil
}

func (r *repositoryImpl) CountByProvider(ctx context.Context, providerID string) (int64, error) {
	query := `SELECT COUNT(*) FROM model_families WHERE provider_id = $1`

	var count int64
	err := r.db.QueryRowContext(ctx, query, providerID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", model_family.ErrDatabase, err)
	}

	return count, nil
}

// GetWithSchemaByID retrieves a model family with its schema by ID
func (r *repositoryImpl) GetWithSchemaByID(ctx context.Context, id string) (*model_family.ModelFamilyWithSchema, error) {
	query := `
		SELECT
			mf.id,
			mf.provider_id,
			mf.schema_id,
			mf.name,
			mf.display_name,
			mf.description,
			mf.max_tokens,
			mf.temperature,
			mf.system_prompt,
			mf.created_at,
			mf.updated_at,

			rs.id,
			rs.provider_id,
			rs.name,
			rs.endpoint_path,
			rs.max_tokens_key,
			rs.system_role_key,
			rs.response_text_path,
			rs.response_image_path,
			rs.request_template,
			rs.supports_temperature,
			rs.supports_streaming,
			rs.created_at,
			rs.updated_at
		FROM model_families mf
		INNER JOIN request_schemas rs
			ON rs.id = mf.schema_id
		WHERE mf.id = $1
	`

	var item model_family.ModelFamilyWithSchema
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.ProviderID,
		&item.SchemaID,
		&item.Name,
		&item.DisplayName,
		&item.Description,
		&item.MaxTokens,
		&item.Temperature,
		&item.SystemPrompt,
		&item.CreatedAt,
		&item.UpdatedAt,

		&item.Schema.ID,
		&item.Schema.ProviderID,
		&item.Schema.Name,
		&item.Schema.EndpointPath,
		&item.Schema.MaxTokensKey,
		&item.Schema.SystemRoleKey,
		&item.Schema.ResponseTextPath,
		&item.Schema.ResponseImagePath,
		&item.Schema.RequestTemplate,
		&item.Schema.SupportsTemperature,
		&item.Schema.SupportsStreaming,
		&item.Schema.CreatedAt,
		&item.Schema.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model_family.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", model_family.ErrDatabase, err)
	}

	return &item, nil
}

func (r *repositoryImpl) UpdateWithTx(ctx context.Context, tx *sql.Tx, mf *model_family.ModelFamilyWithSchema) error {
	query := `
        UPDATE model_families
        SET schema_id = $1, name = $2, display_name = $3, 
            description = $4, max_tokens = $5, temperature = $6,
            system_prompt = $7, updated_at = $8
        WHERE id = $9
    `

	now := time.Now()
	result, err := tx.ExecContext(ctx, query,
		mf.SchemaID,
		mf.Name,
		mf.DisplayName,
		mf.Description,
		mf.MaxTokens,
		mf.Temperature,
		mf.SystemPrompt,
		now,
		mf.ID,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") ||
			strings.Contains(err.Error(), "UNIQUE") ||
			strings.Contains(err.Error(), "23505") {
			return model_family.ErrDuplicate
		}
		return fmt.Errorf("%w: %v", model_family.ErrDatabase, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", model_family.ErrDatabase, err)
	}

	if rowsAffected == 0 {
		return model_family.ErrNotFound
	}

	mf.UpdatedAt = now
	return nil
}

func (r *repositoryImpl) CreateWithTx(ctx context.Context, tx *sql.Tx, mf *model_family.ModelFamily) error {
	query := `
        INSERT INTO model_families (
            id, provider_id, schema_id, name, display_name, 
            description, max_tokens, temperature, system_prompt,
            created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `

	id := uuid.New()
	now := time.Now()

	_, err := tx.ExecContext(ctx, query,
		id,
		mf.ProviderID,
		mf.SchemaID,
		mf.Name,
		mf.DisplayName,
		mf.Description,
		mf.MaxTokens,
		mf.Temperature,
		mf.SystemPrompt,
		now,
		now,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") ||
			strings.Contains(err.Error(), "UNIQUE") ||
			strings.Contains(err.Error(), "23505") {
			return model_family.ErrDuplicate
		}
		return fmt.Errorf("%w: %v", model_family.ErrDatabase, err)
	}

	mf.ID = id.String()
	mf.CreatedAt = now
	mf.UpdatedAt = now

	return nil
}

func (r *repositoryImpl) DeleteWithTx(ctx context.Context, tx *sql.Tx, id string) error {
	query := `DELETE FROM model_families WHERE id = $1`

	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", model_family.ErrDatabase, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", model_family.ErrDatabase, err)
	}

	if rowsAffected == 0 {
		return model_family.ErrNotFound
	}

	return nil
}
