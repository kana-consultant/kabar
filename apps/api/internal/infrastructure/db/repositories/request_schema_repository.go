package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"seo-backend/internal/domain/request_schema"

	"github.com/go-redis/redis/v8"
)

type requestSchemaRepositoryImpl struct {
	db          *sql.DB
	redisClient *redis.Client
}

func NewRequestSchemaRepository(db *sql.DB, redisClient *redis.Client) request_schema.Repository {
	return &requestSchemaRepositoryImpl{
		db:          db,
		redisClient: redisClient,
	}
}

// ============================================================
// CREATE METHODS
// ============================================================

func (r *requestSchemaRepositoryImpl) Create(ctx context.Context, rs *request_schema.RequestSchema) error {
	query := `
		INSERT INTO request_schemas (
			id, provider_id, name, endpoint_path, max_tokens_key,
			system_role_key, response_text_path, response_image_path,
			request_template, supports_temperature, supports_streaming,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.db.ExecContext(ctx, query,
		rs.ID, rs.ProviderID, rs.Name, rs.EndpointPath, rs.MaxTokensKey,
		rs.SystemRoleKey, rs.ResponseTextPath, rs.ResponseImagePath,
		rs.RequestTemplate, rs.SupportsTemperature, rs.SupportsStreaming,
		rs.CreatedAt, rs.UpdatedAt,
	)

	if err != nil {
		r.invalidateProviderCache(ctx, rs.ProviderID)
		return err
	}

	r.invalidateProviderCache(ctx, rs.ProviderID)
	return nil
}

func (r *requestSchemaRepositoryImpl) CreateWithTx(ctx context.Context, tx *sql.Tx, rs *request_schema.RequestSchema) error {
	query := `
		INSERT INTO request_schemas (
			id, provider_id, name, endpoint_path, max_tokens_key,
			system_role_key, response_text_path, response_image_path,
			request_template, supports_temperature, supports_streaming,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := tx.ExecContext(ctx, query,
		rs.ID, rs.ProviderID, rs.Name, rs.EndpointPath, rs.MaxTokensKey,
		rs.SystemRoleKey, rs.ResponseTextPath, rs.ResponseImagePath,
		rs.RequestTemplate, rs.SupportsTemperature, rs.SupportsStreaming,
		rs.CreatedAt, rs.UpdatedAt,
	)

	if err != nil {
		return err
	}

	go r.invalidateProviderCache(context.Background(), rs.ProviderID)
	return nil
}

// ============================================================
// READ METHODS
// ============================================================

func (r *requestSchemaRepositoryImpl) GetByID(ctx context.Context, id string) (*request_schema.RequestSchema, error) {
	cacheKey := fmt.Sprintf("request_schema:id:%s", id)
	cached, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var rs request_schema.RequestSchema
		if err := json.Unmarshal([]byte(cached), &rs); err == nil {
			return &rs, nil
		}
	}

	query := `
		SELECT id, provider_id, name, endpoint_path, max_tokens_key,
			system_role_key, response_text_path, response_image_path,
			request_template, supports_temperature, supports_streaming,
			created_at, updated_at
		FROM request_schemas WHERE id = $1
	`

	var rs request_schema.RequestSchema
	var maxTokensKey, systemRoleKey, responseTextPath, responseImagePath, requestTemplate sql.NullString
	var supportsTemperature, supportsStreaming sql.NullBool

	err = r.db.QueryRowContext(ctx, query, id).Scan(
		&rs.ID, &rs.ProviderID, &rs.Name, &rs.EndpointPath, &maxTokensKey,
		&systemRoleKey, &responseTextPath, &responseImagePath,
		&requestTemplate, &supportsTemperature, &supportsStreaming,
		&rs.CreatedAt, &rs.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, request_schema.ErrNotFound
		}
		return nil, err
	}

	if maxTokensKey.Valid {
		rs.MaxTokensKey = &maxTokensKey.String
	}
	if systemRoleKey.Valid {
		rs.SystemRoleKey = &systemRoleKey.String
	}
	if responseTextPath.Valid {
		rs.ResponseTextPath = &responseTextPath.String
	}
	if responseImagePath.Valid {
		rs.ResponseImagePath = &responseImagePath.String
	}
	if requestTemplate.Valid {
		rs.RequestTemplate = &requestTemplate.String
	}
	if supportsTemperature.Valid {
		rs.SupportsTemperature = &supportsTemperature.Bool
	}
	if supportsStreaming.Valid {
		rs.SupportsStreaming = &supportsStreaming.Bool
	}

	data, _ := json.Marshal(rs)
	r.redisClient.Set(ctx, cacheKey, data, 5*time.Minute)

	return &rs, nil
}

func (r *requestSchemaRepositoryImpl) GetByProviderAndName(ctx context.Context, providerID string, name string) (*request_schema.RequestSchema, error) {
	query := `
		SELECT id, provider_id, name, endpoint_path, max_tokens_key,
			system_role_key, response_text_path, response_image_path,
			request_template, supports_temperature, supports_streaming,
			created_at, updated_at
		FROM request_schemas 
		WHERE provider_id = $1 AND name = $2
	`

	var rs request_schema.RequestSchema
	var maxTokensKey, systemRoleKey, responseTextPath, responseImagePath, requestTemplate sql.NullString
	var supportsTemperature, supportsStreaming sql.NullBool

	err := r.db.QueryRowContext(ctx, query, providerID, name).Scan(
		&rs.ID, &rs.ProviderID, &rs.Name, &rs.EndpointPath, &maxTokensKey,
		&systemRoleKey, &responseTextPath, &responseImagePath,
		&requestTemplate, &supportsTemperature, &supportsStreaming,
		&rs.CreatedAt, &rs.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, request_schema.ErrNotFound
		}
		return nil, err
	}

	if maxTokensKey.Valid {
		rs.MaxTokensKey = &maxTokensKey.String
	}
	if systemRoleKey.Valid {
		rs.SystemRoleKey = &systemRoleKey.String
	}
	if responseTextPath.Valid {
		rs.ResponseTextPath = &responseTextPath.String
	}
	if responseImagePath.Valid {
		rs.ResponseImagePath = &responseImagePath.String
	}
	if requestTemplate.Valid {
		rs.RequestTemplate = &requestTemplate.String
	}
	if supportsTemperature.Valid {
		rs.SupportsTemperature = &supportsTemperature.Bool
	}
	if supportsStreaming.Valid {
		rs.SupportsStreaming = &supportsStreaming.Bool
	}

	return &rs, nil
}

func (r *requestSchemaRepositoryImpl) GetAll(ctx context.Context, limit, offset int) ([]request_schema.RequestSchema, error) {
	query := `
		SELECT id, provider_id, name, endpoint_path, max_tokens_key,
			system_role_key, response_text_path, response_image_path,
			request_template, supports_temperature, supports_streaming,
			created_at, updated_at
		FROM request_schemas
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schemas []request_schema.RequestSchema
	for rows.Next() {
		var rs request_schema.RequestSchema
		var maxTokensKey, systemRoleKey, responseTextPath, responseImagePath, requestTemplate sql.NullString
		var supportsTemperature, supportsStreaming sql.NullBool

		err := rows.Scan(
			&rs.ID, &rs.ProviderID, &rs.Name, &rs.EndpointPath, &maxTokensKey,
			&systemRoleKey, &responseTextPath, &responseImagePath,
			&requestTemplate, &supportsTemperature, &supportsStreaming,
			&rs.CreatedAt, &rs.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if maxTokensKey.Valid {
			rs.MaxTokensKey = &maxTokensKey.String
		}
		if systemRoleKey.Valid {
			rs.SystemRoleKey = &systemRoleKey.String
		}
		if responseTextPath.Valid {
			rs.ResponseTextPath = &responseTextPath.String
		}
		if responseImagePath.Valid {
			rs.ResponseImagePath = &responseImagePath.String
		}
		if requestTemplate.Valid {
			rs.RequestTemplate = &requestTemplate.String
		}
		if supportsTemperature.Valid {
			rs.SupportsTemperature = &supportsTemperature.Bool
		}
		if supportsStreaming.Valid {
			rs.SupportsStreaming = &supportsStreaming.Bool
		}

		schemas = append(schemas, rs)
	}

	return schemas, nil
}

func (r *requestSchemaRepositoryImpl) GetByProvider(ctx context.Context, providerID string) ([]request_schema.RequestSchema, error) {
	cacheKey := fmt.Sprintf("request_schemas:provider:%s", providerID)
	cached, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var schemas []request_schema.RequestSchema
		if err := json.Unmarshal([]byte(cached), &schemas); err == nil {
			return schemas, nil
		}
	}

	query := `
		SELECT id, provider_id, name, endpoint_path, max_tokens_key,
			system_role_key, response_text_path, response_image_path,
			request_template, supports_temperature, supports_streaming,
			created_at, updated_at
		FROM request_schemas 
		WHERE provider_id = $1
		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, providerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schemas []request_schema.RequestSchema
	for rows.Next() {
		var rs request_schema.RequestSchema
		var maxTokensKey, systemRoleKey, responseTextPath, responseImagePath, requestTemplate sql.NullString
		var supportsTemperature, supportsStreaming sql.NullBool

		err := rows.Scan(
			&rs.ID, &rs.ProviderID, &rs.Name, &rs.EndpointPath, &maxTokensKey,
			&systemRoleKey, &responseTextPath, &responseImagePath,
			&requestTemplate, &supportsTemperature, &supportsStreaming,
			&rs.CreatedAt, &rs.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if maxTokensKey.Valid {
			rs.MaxTokensKey = &maxTokensKey.String
		}
		if systemRoleKey.Valid {
			rs.SystemRoleKey = &systemRoleKey.String
		}
		if responseTextPath.Valid {
			rs.ResponseTextPath = &responseTextPath.String
		}
		if responseImagePath.Valid {
			rs.ResponseImagePath = &responseImagePath.String
		}
		if requestTemplate.Valid {
			rs.RequestTemplate = &requestTemplate.String
		}
		if supportsTemperature.Valid {
			rs.SupportsTemperature = &supportsTemperature.Bool
		}
		if supportsStreaming.Valid {
			rs.SupportsStreaming = &supportsStreaming.Bool
		}

		schemas = append(schemas, rs)
	}

	data, _ := json.Marshal(schemas)
	r.redisClient.Set(ctx, cacheKey, data, 5*time.Minute)

	return schemas, nil
}

func (r *requestSchemaRepositoryImpl) GetByProviderSingle(ctx context.Context, providerID string) (*request_schema.RequestSchema, error) {
	cacheKey := fmt.Sprintf("request_schema:provider:%s", providerID)
	cached, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var schema request_schema.RequestSchema
		if err := json.Unmarshal([]byte(cached), &schema); err == nil {
			return &schema, nil
		}
	}

	query := `
		SELECT id, provider_id, name, endpoint_path, max_tokens_key,
			system_role_key, response_text_path, response_image_path,
			request_template, supports_temperature, supports_streaming,
			created_at, updated_at
		FROM request_schemas 
		WHERE provider_id = $1
		LIMIT 1
	`

	var rs request_schema.RequestSchema
	var maxTokensKey, systemRoleKey, responseTextPath, responseImagePath, requestTemplate sql.NullString
	var supportsTemperature, supportsStreaming sql.NullBool

	err = r.db.QueryRowContext(ctx, query, providerID).Scan(
		&rs.ID, &rs.ProviderID, &rs.Name, &rs.EndpointPath, &maxTokensKey,
		&systemRoleKey, &responseTextPath, &responseImagePath,
		&requestTemplate, &supportsTemperature, &supportsStreaming,
		&rs.CreatedAt, &rs.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, request_schema.ErrNotFound
		}
		return nil, err
	}

	if maxTokensKey.Valid {
		rs.MaxTokensKey = &maxTokensKey.String
	}
	if systemRoleKey.Valid {
		rs.SystemRoleKey = &systemRoleKey.String
	}
	if responseTextPath.Valid {
		rs.ResponseTextPath = &responseTextPath.String
	}
	if responseImagePath.Valid {
		rs.ResponseImagePath = &responseImagePath.String
	}
	if requestTemplate.Valid {
		rs.RequestTemplate = &requestTemplate.String
	}
	if supportsTemperature.Valid {
		rs.SupportsTemperature = &supportsTemperature.Bool
	}
	if supportsStreaming.Valid {
		rs.SupportsStreaming = &supportsStreaming.Bool
	}

	data, _ := json.Marshal(rs)
	r.redisClient.Set(ctx, cacheKey, data, 5*time.Minute)

	return &rs, nil
}

// ============================================================
// UPDATE METHODS
// ============================================================

func (r *requestSchemaRepositoryImpl) Update(ctx context.Context, rs *request_schema.RequestSchema) error {
	query := `
		UPDATE request_schemas
		SET name = $1, endpoint_path = $2, max_tokens_key = $3,
			system_role_key = $4, response_text_path = $5, response_image_path = $6,
			request_template = $7, supports_temperature = $8, supports_streaming = $9,
			updated_at = $10
		WHERE id = $11
	`

	result, err := r.db.ExecContext(ctx, query,
		rs.Name, rs.EndpointPath, rs.MaxTokensKey,
		rs.SystemRoleKey, rs.ResponseTextPath, rs.ResponseImagePath,
		rs.RequestTemplate, rs.SupportsTemperature, rs.SupportsStreaming,
		rs.UpdatedAt, rs.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return request_schema.ErrNotFound
	}

	r.invalidateCache(ctx, rs.ID)
	r.invalidateProviderCache(ctx, rs.ProviderID)

	return nil
}

func (r *requestSchemaRepositoryImpl) UpdateWithTx(ctx context.Context, tx *sql.Tx, rs *request_schema.RequestSchema) error {
	query := `
		UPDATE request_schemas
		SET name = $1, endpoint_path = $2, max_tokens_key = $3,
			system_role_key = $4, response_text_path = $5, response_image_path = $6,
			request_template = $7, supports_temperature = $8, supports_streaming = $9,
			updated_at = $10
		WHERE id = $11
	`

	result, err := tx.ExecContext(ctx, query,
		rs.Name, rs.EndpointPath, rs.MaxTokensKey,
		rs.SystemRoleKey, rs.ResponseTextPath, rs.ResponseImagePath,
		rs.RequestTemplate, rs.SupportsTemperature, rs.SupportsStreaming,
		rs.UpdatedAt, rs.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return request_schema.ErrNotFound
	}

	go r.invalidateCache(context.Background(), rs.ID)
	go r.invalidateProviderCache(context.Background(), rs.ProviderID)

	return nil
}

// ============================================================
// DELETE METHODS
// ============================================================

func (r *requestSchemaRepositoryImpl) Delete(ctx context.Context, id string) error {
	var providerID string
	query := `SELECT provider_id FROM request_schemas WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&providerID)
	if err != nil {
		return err
	}

	deleteQuery := `DELETE FROM request_schemas WHERE id = $1`
	_, err = r.db.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return err
	}

	r.invalidateCache(ctx, id)
	r.invalidateProviderCache(ctx, providerID)

	return nil
}

func (r *requestSchemaRepositoryImpl) DeleteWithTx(ctx context.Context, tx *sql.Tx, id string) error {
	var providerID string
	query := `SELECT provider_id FROM request_schemas WHERE id = $1`
	err := tx.QueryRowContext(ctx, query, id).Scan(&providerID)
	if err != nil {
		return err
	}

	deleteQuery := `DELETE FROM request_schemas WHERE id = $1`
	_, err = tx.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return err
	}

	go r.invalidateCache(context.Background(), id)
	go r.invalidateProviderCache(context.Background(), providerID)

	return nil
}

func (r *requestSchemaRepositoryImpl) DeleteByTeam(ctx context.Context, teamID string) error {
	query := `SELECT id FROM request_schemas WHERE provider_id = $1`
	rows, err := r.db.QueryContext(ctx, query, teamID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return err
		}
		ids = append(ids, id)
	}

	deleteQuery := `DELETE FROM request_schemas WHERE provider_id = $1`
	_, err = r.db.ExecContext(ctx, deleteQuery, teamID)
	if err != nil {
		return err
	}

	for _, id := range ids {
		r.invalidateCache(ctx, id)
	}
	r.invalidateProviderCache(ctx, teamID)

	return nil
}

// ============================================================
// UTILITY METHODS
// ============================================================

func (r *requestSchemaRepositoryImpl) Exists(ctx context.Context, providerID string, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM request_schemas WHERE provider_id = $1 AND name = $2)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, providerID, name).Scan(&exists)
	return exists, err
}

func (r *requestSchemaRepositoryImpl) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM request_schemas`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (r *requestSchemaRepositoryImpl) CountByProvider(ctx context.Context, providerID string) (int64, error) {
	query := `SELECT COUNT(*) FROM request_schemas WHERE provider_id = $1`
	var count int64
	err := r.db.QueryRowContext(ctx, query, providerID).Scan(&count)
	return count, err
}

// ============================================================
// CACHE HELPERS
// ============================================================

func (r *requestSchemaRepositoryImpl) invalidateCache(ctx context.Context, id string) {
	if r.redisClient == nil {
		return
	}
	cacheKey := fmt.Sprintf("request_schema:id:%s", id)
	r.redisClient.Del(ctx, cacheKey)
}

func (r *requestSchemaRepositoryImpl) invalidateProviderCache(ctx context.Context, providerID string) {
	if r.redisClient == nil {
		return
	}
	cacheKey1 := fmt.Sprintf("request_schemas:provider:%s", providerID)
	cacheKey2 := fmt.Sprintf("request_schema:provider:%s", providerID)
	r.redisClient.Del(ctx, cacheKey1, cacheKey2)
}
