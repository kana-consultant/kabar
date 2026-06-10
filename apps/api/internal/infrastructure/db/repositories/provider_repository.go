package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"seo-backend/internal/domain/paginate"
	"seo-backend/internal/domain/provider"
	apiprovider "seo-backend/internal/domain/provider"
	userRole "seo-backend/internal/helper/filter"
	"seo-backend/internal/models"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type provider_repositories struct {
	db          *sql.DB
	redisClient *redis.Client
}

// NewRepository creates a new APIProvider repository implementation
func NewProviderRepository(db *sql.DB, redisClient *redis.Client) apiprovider.Repository {
	return &provider_repositories{
		db:          db,
		redisClient: redisClient,
	}
}

// ==================== Cache Helper ====================

func (r *provider_repositories) deleteKeysByPattern(ctx context.Context, pattern string) {
	iter := r.redisClient.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := r.redisClient.Del(ctx, iter.Val()).Err(); err != nil {
			log.Printf("[Cache] failed to delete key %s: %v", iter.Val(), err)
		}
	}
	if err := iter.Err(); err != nil {
		log.Printf("[Cache] scan error: %v", err)
	}
}

func (r *provider_repositories) clearCacheByTeam(ctx context.Context, teamID string) {
	pattern := fmt.Sprintf("api_providers:*:%s:*", teamID)
	r.deleteKeysByPattern(ctx, pattern)
}

func (r *provider_repositories) clearAllProviderCache(ctx context.Context) {
	r.deleteKeysByPattern(ctx, "api_providers:*")
}

func (r *provider_repositories) clearProviderCacheByID(ctx context.Context, id string) {
	// Hapus cache individual provider
	key1 := fmt.Sprintf("api_provider:id:%s", id)
	r.redisClient.Del(ctx, key1)

	// Hapus cache berdasarkan nama
	key2 := fmt.Sprintf("api_provider:name:*")
	r.deleteKeysByPattern(ctx, key2)

	// Hapus cache list (broad)
	r.deleteKeysByPattern(ctx, "api_providers:list:*")
	r.deleteKeysByPattern(ctx, "api_providers:active")
}

// ClearAllCache - public method to clear all cache
func (r *provider_repositories) ClearAllCache(ctx context.Context) {
	r.clearAllProviderCache(ctx)
}

// ClearCacheByID - public method to clear cache by provider ID
func (r *provider_repositories) ClearCacheByID(ctx context.Context, id string) {
	r.clearProviderCacheByID(ctx, id)
}

// ==================== Transaction Management ====================

func (r *provider_repositories) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

// ==================== Create Operations ====================

func (r *provider_repositories) Create(ctx context.Context, provider *apiprovider.APIProvider) error {
	query := `
        INSERT INTO api_providers (
            id, name, display_name, description, base_url,
            auth_type, auth_header, auth_prefix, default_headers,
            is_active, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
    `

	id := uuid.New()
	now := time.Now()

	authType := "bearer"
	if provider.AuthType != nil {
		authType = *provider.AuthType
	}

	authHeader := "Authorization"
	if provider.AuthHeader != nil {
		authHeader = *provider.AuthHeader
	}

	authPrefix := "Bearer"
	if provider.AuthPrefix != nil {
		authPrefix = *provider.AuthPrefix
	}

	defaultHeaders := []byte("{}")
	if provider.DefaultHeaders != nil {
		defaultHeaders = provider.DefaultHeaders
	}

	isActive := true
	if provider.IsActive != nil {
		isActive = *provider.IsActive
	}

	_, err := r.db.ExecContext(ctx, query,
		id,
		provider.Name,
		provider.DisplayName,
		provider.Description,
		provider.BaseURL,
		authType,
		authHeader,
		authPrefix,
		defaultHeaders,
		isActive,
		now,
		now,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") ||
			strings.Contains(err.Error(), "UNIQUE") ||
			strings.Contains(err.Error(), "23505") {
			return apiprovider.ErrDuplicate
		}
		return fmt.Errorf("%w: %v", apiprovider.ErrDatabase, err)
	}

	provider.ID = id
	provider.CreatedAt = now
	provider.UpdatedAt = now
	provider.AuthType = &authType
	provider.AuthHeader = &authHeader
	provider.AuthPrefix = &authPrefix
	provider.IsActive = &isActive

	// Clear cache after create
	go r.clearAllProviderCache(context.Background())

	return nil
}

func (r *provider_repositories) CreateWithTx(ctx context.Context, tx *sql.Tx, provider *apiprovider.APIProvider) error {
	query := `
        INSERT INTO api_providers (
            id, name, display_name, description, base_url,
            auth_type, auth_header, auth_prefix, default_headers,
            is_active, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
    `

	id := uuid.New()
	now := time.Now()

	authType := "bearer"
	if provider.AuthType != nil {
		authType = *provider.AuthType
	}

	authHeader := "Authorization"
	if provider.AuthHeader != nil {
		authHeader = *provider.AuthHeader
	}

	authPrefix := "Bearer"
	if provider.AuthPrefix != nil {
		authPrefix = *provider.AuthPrefix
	}

	defaultHeaders := []byte("{}")
	if provider.DefaultHeaders != nil {
		defaultHeaders = provider.DefaultHeaders
	}

	isActive := true
	if provider.IsActive != nil {
		isActive = *provider.IsActive
	}

	_, err := tx.ExecContext(ctx, query,
		id,
		provider.Name,
		provider.DisplayName,
		provider.Description,
		provider.BaseURL,
		authType,
		authHeader,
		authPrefix,
		defaultHeaders,
		isActive,
		now,
		now,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") ||
			strings.Contains(err.Error(), "UNIQUE") ||
			strings.Contains(err.Error(), "23505") {
			return apiprovider.ErrDuplicate
		}
		return fmt.Errorf("%w: %v", apiprovider.ErrDatabase, err)
	}

	provider.ID = id
	provider.CreatedAt = now
	provider.UpdatedAt = now
	provider.AuthType = &authType
	provider.AuthHeader = &authHeader
	provider.AuthPrefix = &authPrefix
	provider.IsActive = &isActive

	return nil
}

// ==================== Read Operations ====================

func (r *provider_repositories) GetByID(ctx context.Context, id string) (*provider.APIProvider, error) {
	cacheKey := fmt.Sprintf("api_provider:id:%s", id)

	log.Printf("[Provider:GetByID] lookup cache key=%s", cacheKey)

	// Try cache first
	cached, err := r.redisClient.Get(ctx, cacheKey).Bytes()
	if err == nil {
		log.Printf("[Provider:GetByID] cache HIT key=%s", cacheKey)

		var provider provider.APIProvider
		if err := json.Unmarshal(cached, &provider); err == nil {
			log.Printf("[Provider:GetByID] cache unmarshaled successfully id=%s", id)
			return &provider, nil
		}

		log.Printf("[Provider:GetByID] cache unmarshal failed key=%s err=%v", cacheKey, err)
	} else {
		log.Printf("[Provider:GetByID] cache MISS key=%s err=%v", cacheKey, err)
	}

	log.Printf("[Provider:GetByID] querying database id=%s", id)

	query := `
        SELECT id, name, display_name, description, base_url,
               auth_type, auth_header, auth_prefix, default_headers,
               is_active, created_at, updated_at
        FROM api_providers
        WHERE id = $1
    `

	var provider provider.APIProvider
	var defaultHeaders []byte
	var description sql.NullString
	var authType, authHeader, authPrefix sql.NullString

	err = r.db.QueryRowContext(ctx, query, id).Scan(
		&provider.ID,
		&provider.Name,
		&provider.DisplayName,
		&description,
		&provider.BaseURL,
		&authType,
		&authHeader,
		&authPrefix,
		&defaultHeaders,
		&provider.IsActive,
		&provider.CreatedAt,
		&provider.UpdatedAt,
	)

	if err != nil {
		log.Printf("[Provider:GetByID] database error id=%s err=%v", id, err)

		if err == sql.ErrNoRows {
			return nil, apiprovider.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", apiprovider.ErrDatabase, err)
	}

	log.Printf("[Provider:GetByID] database HIT id=%s", provider.ID)

	if description.Valid {
		provider.Description = &description.String
	}
	if authType.Valid {
		provider.AuthType = &authType.String
	}
	if authHeader.Valid {
		provider.AuthHeader = &authHeader.String
	}
	if authPrefix.Valid {
		provider.AuthPrefix = &authPrefix.String
	}

	provider.DefaultHeaders = json.RawMessage(defaultHeaders)

	// Store in cache
	go func() {
		data, err := json.Marshal(provider)
		if err != nil {
			log.Printf("[Provider:GetByID] cache marshal failed id=%s err=%v", id, err)
			return
		}

		if err := r.redisClient.Set(
			context.Background(),
			cacheKey,
			data,
			5*time.Minute,
		).Err(); err != nil {
			log.Printf("[Provider:GetByID] cache SET failed key=%s err=%v", cacheKey, err)
			return
		}

		log.Printf("[Provider:GetByID] cache SET success key=%s ttl=5m", cacheKey)
	}()

	log.Printf(
		"[Provider:GetByID] raw DB values authType={Valid:%v Value:%q} authHeader={Valid:%v Value:%q} authPrefix={Valid:%v Value:%q}",
		authType.Valid, authType.String,
		authHeader.Valid, authHeader.String,
		authPrefix.Valid, authPrefix.String,
	)

	return &provider, nil
}

func (r *provider_repositories) GetByName(ctx context.Context, name string) (*provider.APIProvider, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("api_provider:name:%s", name)
	cached, err := r.redisClient.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var provider provider.APIProvider
		if err := json.Unmarshal(cached, &provider); err == nil {
			return &provider, nil
		}
	}

	query := `
        SELECT id, name, display_name, description, base_url,
               auth_type, auth_header, auth_prefix, default_headers,
               is_active, created_at, updated_at
        FROM api_providers
        WHERE name = $1
    `

	var provider provider.APIProvider
	var defaultHeaders []byte
	var description sql.NullString
	var authType, authHeader, authPrefix sql.NullString

	err = r.db.QueryRowContext(ctx, query, name).Scan(
		&provider.ID,
		&provider.Name,
		&provider.DisplayName,
		&description,
		&provider.BaseURL,
		&authType,
		&authHeader,
		&authPrefix,
		&defaultHeaders,
		&provider.IsActive,
		&provider.CreatedAt,
		&provider.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apiprovider.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", apiprovider.ErrDatabase, err)
	}

	if description.Valid {
		provider.Description = &description.String
	}
	if authType.Valid {
		provider.AuthType = &authType.String
	}
	if authHeader.Valid {
		provider.AuthHeader = &authHeader.String
	}
	if authPrefix.Valid {
		provider.AuthPrefix = &authPrefix.String
	}
	provider.DefaultHeaders = json.RawMessage(defaultHeaders)

	// Store in cache
	go func() {
		data, _ := json.Marshal(provider)
		r.redisClient.Set(context.Background(), cacheKey, data, 5*time.Minute)
	}()

	return &provider, nil
}

func (r *provider_repositories) GetAll(ctx context.Context, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[provider.APIProvider], error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	// Cache key berdasarkan user context
	cacheKey := fmt.Sprintf("api_providers:list:%s:%s:limit=%d:offset=%d:search=%s",
		userCtx.GetRole(), userCtx.GetTeamID(), params.Limit, params.Offset, params.Search)

	// Try cache
	cached, err := r.redisClient.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var result paginate.PaginatedResult[provider.APIProvider]
		if err := json.Unmarshal(cached, &result); err == nil {
			return &result, nil
		}
	}

	// Build access filter
	whereClause, whereArgs := userRole.BuildAccessFilter(userCtx)

	// Build search clause
	args := append([]any{}, whereArgs...)
	searchClause := ""
	if params.Search != "" {
		args = append(args, "%"+params.Search+"%")
		searchClause = fmt.Sprintf(" AND (name ILIKE $%d OR display_name ILIKE $%d)", len(args), len(args))
	}

	fullWhere := whereClause + searchClause

	// Count query
	var totalItems int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM api_providers WHERE %s", fullWhere)
	err = r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalItems)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", apiprovider.ErrDatabase, err)
	}

	totalPages := totalItems / params.Limit
	if totalItems%params.Limit > 0 {
		totalPages++
	}
	currentPage := (params.Offset / params.Limit) + 1

	// Append LIMIT & OFFSET
	args = append(args, params.Limit, params.Offset)
	query := fmt.Sprintf(`
		SELECT id, name, display_name, description, base_url,
		       auth_type, auth_header, auth_prefix, default_headers,
		       is_active, created_at, updated_at
		FROM api_providers
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, fullWhere, len(args)-1, len(args))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", apiprovider.ErrDatabase, err)
	}
	defer rows.Close()

	var providers []provider.APIProvider
	for rows.Next() {
		var p provider.APIProvider
		var defaultHeaders []byte
		var description sql.NullString
		var authType, authHeader, authPrefix sql.NullString

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.DisplayName,
			&description,
			&p.BaseURL,
			&authType,
			&authHeader,
			&authPrefix,
			&defaultHeaders,
			&p.IsActive,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", apiprovider.ErrDatabase, err)
		}

		if description.Valid {
			p.Description = &description.String
		}
		if authType.Valid {
			p.AuthType = &authType.String
		}
		if authHeader.Valid {
			p.AuthHeader = &authHeader.String
		}
		if authPrefix.Valid {
			p.AuthPrefix = &authPrefix.String
		}
		p.DefaultHeaders = json.RawMessage(defaultHeaders)

		providers = append(providers, p)
	}

	result := &paginate.PaginatedResult[provider.APIProvider]{
		Data:        providers,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		Limit:       params.Limit,
		Offset:      params.Offset,
	}

	// Store in cache
	go func() {
		data, _ := json.Marshal(result)
		r.redisClient.Set(context.Background(), cacheKey, data, 2*time.Minute)
	}()

	return result, nil
}

func (r *provider_repositories) GetActive(ctx context.Context) ([]provider.APIProvider, error) {
	cacheKey := "api_providers:active"
	cached, err := r.redisClient.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var providers []provider.APIProvider
		if err := json.Unmarshal(cached, &providers); err == nil {
			return providers, nil
		}
	}

	query := `
        SELECT id, name, display_name, description, base_url,
               auth_type, auth_header, auth_prefix, default_headers,
               is_active, created_at, updated_at
        FROM api_providers
        WHERE is_active = true
        ORDER BY name ASC
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", apiprovider.ErrDatabase, err)
	}
	defer rows.Close()

	var providers []provider.APIProvider
	for rows.Next() {
		var p provider.APIProvider
		var defaultHeaders []byte
		var description sql.NullString
		var authType, authHeader, authPrefix sql.NullString

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.DisplayName,
			&description,
			&p.BaseURL,
			&authType,
			&authHeader,
			&authPrefix,
			&defaultHeaders,
			&p.IsActive,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", apiprovider.ErrDatabase, err)
		}

		if description.Valid {
			p.Description = &description.String
		}
		if authType.Valid {
			p.AuthType = &authType.String
		}
		if authHeader.Valid {
			p.AuthHeader = &authHeader.String
		}
		if authPrefix.Valid {
			p.AuthPrefix = &authPrefix.String
		}
		p.DefaultHeaders = json.RawMessage(defaultHeaders)

		providers = append(providers, p)
	}

	// Store in cache
	go func() {
		data, _ := json.Marshal(providers)
		r.redisClient.Set(context.Background(), cacheKey, data, 2*time.Minute)
	}()

	return providers, nil
}

// ==================== Update Operations ====================

func (r *provider_repositories) Update(ctx context.Context, provider *apiprovider.APIProvider) error {
	query := `
        UPDATE api_providers
        SET name = $1, display_name = $2, description = $3, base_url = $4,
            auth_type = $5, auth_header = $6, auth_prefix = $7, 
            default_headers = $8, is_active = $9, updated_at = $10
        WHERE id = $11
    `

	now := time.Now()

	authType := "bearer"
	if provider.AuthType != nil {
		authType = *provider.AuthType
	}

	authHeader := "Authorization"
	if provider.AuthHeader != nil {
		authHeader = *provider.AuthHeader
	}

	authPrefix := "Bearer"
	if provider.AuthPrefix != nil {
		authPrefix = *provider.AuthPrefix
	}

	defaultHeaders := []byte("{}")
	if provider.DefaultHeaders != nil {
		defaultHeaders = provider.DefaultHeaders
	}

	isActive := true
	if provider.IsActive != nil {
		isActive = *provider.IsActive
	}

	result, err := r.db.ExecContext(ctx, query,
		provider.Name,
		provider.DisplayName,
		provider.Description,
		provider.BaseURL,
		authType,
		authHeader,
		authPrefix,
		defaultHeaders,
		isActive,
		now,
		provider.ID,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") ||
			strings.Contains(err.Error(), "UNIQUE") ||
			strings.Contains(err.Error(), "23505") {
			return apiprovider.ErrDuplicate
		}
		return fmt.Errorf("%w: %v", apiprovider.ErrDatabase, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apiprovider.ErrNotFound
	}

	provider.UpdatedAt = now

	// Clear cache after update
	go r.clearAllProviderCache(context.Background())

	return nil
}

func (r *provider_repositories) UpdateWithTx(ctx context.Context, tx *sql.Tx, id string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	setClause := ""
	args := []interface{}{}
	i := 1

	for key, value := range updates {
		if setClause != "" {
			setClause += ", "
		}
		setClause += fmt.Sprintf("%s = $%d", key, i)
		args = append(args, value)
		i++
	}

	setClause += fmt.Sprintf(", updated_at = $%d", i)
	args = append(args, time.Now())
	i++

	args = append(args, id)

	query := fmt.Sprintf("UPDATE api_providers SET %s WHERE id = $%d", setClause, i)

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") ||
			strings.Contains(err.Error(), "UNIQUE") ||
			strings.Contains(err.Error(), "23505") {
			return apiprovider.ErrDuplicate
		}
		return fmt.Errorf("%w: %v", apiprovider.ErrDatabase, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apiprovider.ErrNotFound
	}

	return nil
}

func (r *provider_repositories) UpdateDefaultHeaders(ctx context.Context, id string, headers json.RawMessage) error {
	query := `UPDATE api_providers SET default_headers = $1, updated_at = $2 WHERE id = $3`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, headers, now, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apiprovider.ErrNotFound
	}

	// Clear cache after update
	go r.clearAllProviderCache(context.Background())

	return nil
}

func (r *provider_repositories) UpdateDefaultHeadersWithTx(ctx context.Context, tx *sql.Tx, id string, headers json.RawMessage) error {
	query := `UPDATE api_providers SET default_headers = $1, updated_at = $2 WHERE id = $3`

	now := time.Now()
	result, err := tx.ExecContext(ctx, query, headers, now, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apiprovider.ErrNotFound
	}

	return nil
}

func (r *provider_repositories) ToggleActive(ctx context.Context, id string, isActive bool) error {
	query := `UPDATE api_providers SET is_active = $1, updated_at = $2 WHERE id = $3`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, isActive, now, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apiprovider.ErrNotFound
	}

	// Clear cache after toggle
	go r.clearAllProviderCache(context.Background())

	return nil
}

func (r *provider_repositories) ToggleActiveWithTx(ctx context.Context, tx *sql.Tx, id string, isActive bool) error {
	query := `UPDATE api_providers SET is_active = $1, updated_at = $2 WHERE id = $3`

	now := time.Now()
	result, err := tx.ExecContext(ctx, query, isActive, now, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apiprovider.ErrNotFound
	}

	return nil
}

// ==================== Delete Operations ====================

func (r *provider_repositories) Delete(ctx context.Context, id string) error {
	query := `UPDATE api_providers SET is_active = false, updated_at = $1 WHERE id = $2`

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
		return apiprovider.ErrNotFound
	}

	// Clear cache after delete
	go r.clearAllProviderCache(context.Background())

	return nil
}

func (r *provider_repositories) DeleteWithTx(ctx context.Context, tx *sql.Tx, id string) error {
	query := `UPDATE api_providers SET is_active = false, updated_at = $1 WHERE id = $2`

	now := time.Now()
	result, err := tx.ExecContext(ctx, query, now, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apiprovider.ErrNotFound
	}

	return nil
}

func (r *provider_repositories) HardDelete(ctx context.Context, id string) error {
	query := `DELETE FROM api_providers WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apiprovider.ErrNotFound
	}

	// Clear cache after hard delete
	go r.clearAllProviderCache(context.Background())

	return nil
}

func (r *provider_repositories) HardDeleteWithTx(ctx context.Context, tx *sql.Tx, id string) error {
	query := `DELETE FROM api_providers WHERE id = $1`

	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apiprovider.ErrNotFound
	}

	return nil
}

// ==================== Utility Operations ====================

func (r *provider_repositories) Exists(ctx context.Context, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM api_providers WHERE name = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%w: %v", apiprovider.ErrDatabase, err)
	}

	return exists, nil
}

func (r *provider_repositories) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM api_providers`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", apiprovider.ErrDatabase, err)
	}

	return count, nil
}

func (r *provider_repositories) CheckProviderUsage(ctx context.Context, providerID string) (int, error) {
	query := `
        SELECT COUNT(*) FROM api_keys WHERE provider_id = $1
    `

	var count int
	err := r.db.QueryRowContext(ctx, query, providerID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", apiprovider.ErrDatabase, err)
	}

	return count, nil
}
