// internal/infrastructure/repository/draft/repository.go
package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"reflect"
	"strings"
	"time"

	"seo-backend/internal/domain/draft"
	"seo-backend/internal/domain/paginate"
	"seo-backend/internal/helper"
	userRole "seo-backend/internal/helper/filter"
	"seo-backend/internal/helper/keywords"
	"seo-backend/internal/models"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type RepositoryImpl struct {
	db          *sql.DB
	redisClient *redis.Client
}

func NewDraftRepository(db *sql.DB, redisClient *redis.Client) draft.Repository {
	return &RepositoryImpl{
		db:          db,
		redisClient: redisClient,
	}
}
func (r *RepositoryImpl) GetByID(
	ctx context.Context,
	id string,
) (*draft.DraftData, error) {

	var d draft.DraftData
	var targetProductsJSON []byte

	// ✅ Fix: hapus trailing comma setelah target_products
	query := `
	SELECT 
		id, 
		title, 
		topic, 
		article, 
		image_url, 
		COALESCE(image_prompt, ''),
		target_products,
		seo_score
	FROM drafts 
	WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&d.ID,
		&d.Title,
		&d.Topic,
		&d.Article,
		&d.ImageURL,
		&d.ImagePrompt,
		&targetProductsJSON,
		&d.SEOScore,
	)

	if err != nil {
		log.Printf("Error getting draft by ID %s: %v", id, err)
		return nil, err
	}

	log.Printf("DRAFT DATA => %+v", d)

	if len(targetProductsJSON) > 0 {
		if err := json.Unmarshal(targetProductsJSON, &d.TargetProducts); err != nil {
			log.Printf("Error unmarshaling target_products: %v", err)
		}
	}

	log.Printf("TARGET PRODUCTS => %+v", d.TargetProducts)

	keywords, err := keywords.GetKeywords(ctx, r.db, keywords.DraftSource{DraftID: id})
	if err != nil {
		log.Printf("Error getting keywords for draft %s: %v", id, err)
		d.Keywords = []string{}
	} else {
		d.Keywords = keywords
	}

	log.Printf("KEYWORDS => %+v", d.Keywords)
	log.Printf("FINAL DRAFT RESPONSE => %+v", d)

	return &d, nil
}
func (r *RepositoryImpl) GetAll(ctx context.Context, filter models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[draft.Draft], error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}

	// Build access filter
	whereClause, whereArgs := userRole.BuildAccessFilter(filter)

	// Cache key
	cacheKey := fmt.Sprintf("draft_list:%s:%s:limit=%d:offset=%d:search=%s", filter.GetRole(), filter.GetTeamID(), params.Limit, params.Offset, params.Search)

	// Cek Redis cache
	cached, err := r.redisClient.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var result paginate.PaginatedResult[draft.Draft]
		if err := json.Unmarshal(cached, &result); err == nil {
			log.Printf("[GetAll] cache hit | filter=%+v | key=%s", filter, cacheKey)
			return &result, nil
		}
	}

	log.Printf("[GetAll] cache miss | filter=%+v | key=%s", filter, cacheKey)

	// Build search clause
	args := append([]any{}, whereArgs...)
	searchClause := ""
	if params.Search != "" {
		args = append(args, "%"+params.Search+"%")
		searchClause = fmt.Sprintf(" AND (title ILIKE $%d OR topic ILIKE $%d)", len(args), len(args))
	}

	// Combine WHERE
	fullWhere := whereClause + searchClause + " AND status != 'scheduled'"

	// Count query
	var totalItems int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM drafts WHERE %s", fullWhere)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(params.Limit)))
	currentPage := (params.Offset / params.Limit) + 1

	// Append LIMIT & OFFSET
	args = append(args, params.Limit, params.Offset)
	query := fmt.Sprintf(`
        SELECT 
            id, title, topic, article, target_products, team_id,
            image_url, status, seo_score, COALESCE(image_prompt, '')
        FROM drafts
        WHERE %s
        ORDER BY created_at DESC
        LIMIT $%d OFFSET $%d
    `, fullWhere, len(args)-1, len(args))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("[GetAll] query error | filter=%+v | err=%v", filter, err)
		return nil, err
	}
	defer rows.Close()

	var drafts []draft.Draft
	for rows.Next() {
		var d draft.Draft
		var targetProductsJSON []byte

		err := rows.Scan(
			&d.ID, &d.Title, &d.Topic, &d.Article,
			&targetProductsJSON, &d.TeamID, &d.ImageURL,
			&d.Status, &d.SeoScore, &d.ImagePrompt,
		)
		if err != nil {
			log.Printf("[GetAll] scan error | filter=%+v | err=%v", filter, err)
			return nil, err
		}

		if err := json.Unmarshal(targetProductsJSON, &d.TargetProducts); err != nil {
			return nil, err
		}

		drafts = append(drafts, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := &paginate.PaginatedResult[draft.Draft]{
		Data:        drafts,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		Limit:       params.Limit,
		Offset:      params.Offset,
	}

	// Simpan ke Redis, TTL 2 menit
	resultBytes, err := json.Marshal(result)
	if err != nil {
		log.Printf("[GetAll] failed to marshal result | filter=%+v | err=%v", filter, err)
	} else {
		if err := r.redisClient.Set(ctx, cacheKey, resultBytes, 2*time.Minute).Err(); err != nil {
			log.Printf("[GetAll] failed to set cache | filter=%+v | err=%v", filter, err)
		} else {
			log.Printf("[GetAll] cache saved | filter=%+v | key=%s | ttl=2m", filter, cacheKey)
		}
	}

	return result, nil
}

func (r *RepositoryImpl) GetDashboardStats(ctx context.Context, filter models.UserContext) (*draft.DraftStats, error) {
	// Build access filter
	whereClause, whereArgs := userRole.BuildAccessFilter(filter)

	cacheKey := fmt.Sprintf("dashboard_stats:%s", filter.GetRole())
	if filter.GetTeamID() != "" {
		cacheKey = fmt.Sprintf("dashboard_stats:%s:%s", filter.GetRole(), filter.GetTeamID())
	}

	// 1. Cek Redis cache dulu
	cached, err := r.redisClient.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var stats draft.DraftStats
		if err := json.Unmarshal(cached, &stats); err == nil {
			log.Printf("[DashboardStats] cache hit | filter=%+v", filter)
			return &stats, nil
		}
	}

	log.Printf("[DashboardStats] cache miss | filter=%+v", filter)

	stats := &draft.DraftStats{
		ProductCoverage: make(map[string]int),
	}

	// 2. Total draft, with/without image, scheduled
	summaryQuery := fmt.Sprintf(`
		SELECT
			COUNT(*)                                        AS total_draft,
			COUNT(*) FILTER (WHERE has_image = true)        AS total_with_image,
			COUNT(*) FILTER (WHERE has_image = false)       AS total_without_image,
			COUNT(*) FILTER (WHERE status = 'scheduled')    AS total_scheduled
		FROM drafts
		WHERE %s
	`, whereClause)
	if err := r.db.QueryRowContext(ctx, summaryQuery, whereArgs...).Scan(
		&stats.TotalDraft,
		&stats.TotalWithImage,
		&stats.TotalWithoutImage,
		&stats.TotalScheduled,
	); err != nil {
		log.Printf("[DashboardStats] failed to scan summary stats | filter=%+v | err=%v", filter, err)
		return nil, fmt.Errorf("failed to get summary stats: %w", err)
	}
	log.Printf("[DashboardStats] summary | filter=%+v | total=%d with_image=%d without_image=%d scheduled=%d",
		filter, stats.TotalDraft, stats.TotalWithImage, stats.TotalWithoutImage, stats.TotalScheduled)

	// 3. Product coverage
	productQuery := fmt.Sprintf(`
		SELECT p.name, COUNT(*) AS count
		FROM (
			SELECT id, target_products
			FROM drafts
			WHERE %s
			AND target_products IS NOT NULL
			AND jsonb_typeof(target_products) = 'array'
		) filtered_drafts,
		jsonb_array_elements_text(filtered_drafts.target_products) AS product_id
		JOIN products p ON p.id = product_id::uuid
		GROUP BY p.id, p.name
		ORDER BY count DESC
	`, whereClause)

	productRows, err := r.db.QueryContext(ctx, productQuery, whereArgs...)
	if err != nil {
		log.Printf("[DashboardStats] failed to query product coverage | filter=%+v | err=%v", filter, err)
		return nil, fmt.Errorf("failed to get product coverage: %w", err)
	}
	defer productRows.Close()

	for productRows.Next() {
		var product string
		var count int
		if err := productRows.Scan(&product, &count); err != nil {
			log.Printf("[DashboardStats] failed to scan product coverage | filter=%+v | err=%v", filter, err)
			return nil, fmt.Errorf("failed to scan product coverage: %w", err)
		}
		stats.ProductCoverage[product] = count
	}
	if err := productRows.Err(); err != nil {
		log.Printf("[DashboardStats] product rows iteration error | filter=%+v | err=%v", filter, err)
		return nil, err
	}
	log.Printf("[DashboardStats] product coverage | filter=%+v | total_products=%d", filter, len(stats.ProductCoverage))

	// 4. Daily activity — 30 hari terakhir
	activityQuery := fmt.Sprintf(`
		SELECT
			TO_CHAR(created_at, 'YYYY-MM-DD') AS date,
			COUNT(*) AS count
		FROM drafts
		WHERE %s
			AND created_at >= NOW() - INTERVAL '30 days'
		GROUP BY date
		ORDER BY date ASC
	`, whereClause)
	activityRows, err := r.db.QueryContext(ctx, activityQuery, whereArgs...)
	if err != nil {
		log.Printf("[DashboardStats] failed to query daily activity | filter=%+v | err=%v", filter, err)
		return nil, fmt.Errorf("failed to get daily activity: %w", err)
	}
	defer activityRows.Close()

	for activityRows.Next() {
		var a draft.DailyActivity
		if err := activityRows.Scan(&a.Date, &a.Count); err != nil {
			log.Printf("[DashboardStats] failed to scan daily activity | filter=%+v | err=%v", filter, err)
			return nil, fmt.Errorf("failed to scan daily activity: %w", err)
		}
		stats.DailyActivity = append(stats.DailyActivity, a)
	}
	if err := activityRows.Err(); err != nil {
		log.Printf("[DashboardStats] activity rows iteration error | filter=%+v | err=%v", filter, err)
		return nil, err
	}
	log.Printf("[DashboardStats] daily activity | filter=%+v | total_days=%d", filter, len(stats.DailyActivity))

	// 5. Simpan ke Redis cache, TTL 5 menit
	statsBytes, err := json.Marshal(stats)
	if err != nil {
		log.Printf("[DashboardStats] failed to marshal stats for cache | filter=%+v | err=%v", filter, err)
	} else {
		if err := r.redisClient.Set(ctx, cacheKey, statsBytes, 5*time.Minute).Err(); err != nil {
			log.Printf("[DashboardStats] failed to set cache | filter=%+v | err=%v", filter, err)
		} else {
			log.Printf("[DashboardStats] cache saved | filter=%+v | ttl=5m", filter)
		}
	}

	log.Printf("[DashboardStats] completed | filter=%+v", filter)
	return stats, nil
}

func (r *RepositoryImpl) GetAllScheduled(ctx context.Context, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[draft.Draft], error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}

	// Build access filter
	whereClause, whereArgs := userRole.BuildAccessFilter(userCtx)

	var totalItems int
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM drafts
		WHERE %s AND status = 'scheduled'
	`, whereClause)
	if err := r.db.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalItems); err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(params.Limit)))
	currentPage := (params.Offset / params.Limit) + 1

	// Append LIMIT & OFFSET
	args := append(whereArgs, params.Limit, params.Offset)
	query := fmt.Sprintf(`
		SELECT 
			id,
			title, 
			topic, 
			article, 
			target_products, 
			team_id,
			image_url,
			status,
			COALESCE(image_prompt, ''),
			scheduled_for
		FROM drafts
		WHERE %s AND status = 'scheduled'
		ORDER BY scheduled_for ASC
		LIMIT $%d OFFSET $%d
	`, whereClause, len(args)-1, len(args))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var drafts []draft.Draft

	for rows.Next() {
		var d draft.Draft
		var targetProductsJSON []byte

		err := rows.Scan(
			&d.ID,
			&d.Title,
			&d.Topic,
			&d.Article,
			&targetProductsJSON,
			&d.TeamID,
			&d.ImageURL,
			&d.Status,
			&d.ImagePrompt,
			&d.ScheduledFor,
		)
		if err != nil {
			log.Printf("scan error: %v", err)
			return nil, err
		}

		if err := json.Unmarshal(targetProductsJSON, &d.TargetProducts); err != nil {
			return nil, err
		}

		drafts = append(drafts, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &paginate.PaginatedResult[draft.Draft]{
		Data:        drafts,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		Limit:       params.Limit,
		Offset:      params.Offset,
	}, nil
}

func (r *RepositoryImpl) Create(ctx context.Context, req draft.CreateDraftRequest, userID, teamID string) (string, error) {
	targetProductsJSON, _ := json.Marshal(req.TargetProducts)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// ✅ Fix: hapus duplikat excerpt, sesuaikan jumlah kolom & placeholder ($1-$14)
	query := `INSERT INTO drafts (
		title, topic, article, image_url, image_prompt,
		status, target_products, has_image, created_by, team_id, user_id, slug, excerpt, seo_score
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	RETURNING id`

	var draftID string
	err = tx.QueryRowContext(
		ctx,
		query,
		req.Title, req.Topic, req.Article, req.ImageURL, req.ImagePrompt,
		"draft", targetProductsJSON, req.HasImage, nullIfEmpty(userID), nullIfEmpty(teamID), nullIfEmpty(userID), req.Slug, req.Excerpt, req.SEOScore,
	).Scan(&draftID)

	if err != nil {
		log.Printf("Error creating draft: %v", err)
		return "", err
	}

	if len(req.Keywords) > 0 {
		log.Printf("RAW KEYWORDS => %+v", req.Keywords)

		uniqueKeywords := keywords.UniqueStrings(req.Keywords)
		log.Printf("UNIQUE KEYWORDS => %+v", uniqueKeywords)

		var keywordEntities []draft.Keywords
		now := time.Now()

		for _, keyword := range uniqueKeywords {
			keywordEntities = append(keywordEntities, draft.Keywords{
				ID:        uuid.NewString(),
				IDDraft:   draftID,
				Name:      keyword,
				CreatedAt: now,
				UpdatedAt: now,
			})
		}

		log.Printf("KEYWORD ENTITIES => %+v", keywordEntities)

		err = keywords.InsertKeywordsWithDuplicateCheck(ctx, tx, draftID, keywordEntities)
		if err != nil {
			log.Printf("FAILED INSERT KEYWORDS => %v", err)
			return "", fmt.Errorf("failed to insert keywords: %w", err)
		}

		log.Println("INSERT KEYWORDS SUCCESS")
	}

	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.invalidateCache(ctx,
		fmt.Sprintf("dashboard_stats:%s", teamID))

	if err := r.InvalidateDraftCacheByTeam(ctx, teamID); err != nil {
		log.Printf("failed to invalidate draft cache: %v", err)
	}

	return draftID, nil
}

func (r *RepositoryImpl) Update(ctx context.Context, id string, TeamID string, data map[string]interface{}) error {
	// Start transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Update draft
	data["updated_at"] = helper.ParseWIBTime(time.Now().Format(time.RFC3339))

	query, args, err := r.buildUpdateQuery(id, data)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update draft: %w", err)
	}

	// Handle keywords update if present in data
	if kw, ok := data["keywords"]; ok {
		if err := keywords.UpdateKeywords(ctx, tx, keywords.DraftSource{DraftID: id}, kw); err != nil {
			return fmt.Errorf("failed to update keywords: %w", err)
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.invalidateCache(ctx,
		fmt.Sprintf("dashboard_stats:%s", TeamID))

	if err := r.InvalidateDraftCacheByTeam(ctx, TeamID); err != nil {
		log.Printf("failed to invalidate draft cache: %v", err)
	}

	return nil
}

// updateKeywords handles keywords update strategies
func (r *RepositoryImpl) updateKeywords(ctx context.Context, tx *sql.Tx, draftID string, kw interface{}) error {
	switch kw := kw.(type) {
	case []string:
		// Jika keywords berupa slice of strings
		return keywords.ReplaceKeywords(ctx, tx, keywords.DraftSource{DraftID: draftID}, kw)

	case []draft.Keywords:
		// Jika keywords berupa slice of Keyword struct
		names := make([]string, len(kw))
		for i, k := range kw {
			names[i] = k.Name
		}
		return keywords.ReplaceKeywords(ctx, tx, keywords.DraftSource{DraftID: draftID}, names)

	default:
		return fmt.Errorf("unsupported keywords type: %T", kw)
	}

}

func (r *RepositoryImpl) UpdateStatus(ctx context.Context, id string, status string, scheduledFor *time.Time) error {
	if scheduledFor != nil {
		_, err := r.db.ExecContext(ctx, `
			UPDATE drafts 
			SET status = $1, scheduled_for = $2 
			WHERE id = $3
		`, status, *scheduledFor, id)
		return err
	}

	_, err := r.db.ExecContext(ctx, `
		UPDATE drafts 
		SET status = $1, updated_at = $2
		WHERE id = $3 AND status = 'scheduled'
	`, status, helper.ParseWIBTime(time.Now().Format(time.RFC3339)), id)
	return err
}

func (r *RepositoryImpl) Delete(ctx context.Context, TeamID string, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM drafts WHERE id = $1", id)

	r.invalidateCache(ctx,
		fmt.Sprintf("dashboard_stats:%s", TeamID))

	if err := r.InvalidateDraftCacheByTeam(ctx, TeamID); err != nil {
		log.Printf("failed to invalidate draft cache: %v", err)
	}
	return err
}

func (r *RepositoryImpl) InsertScheduledDraft(ctx context.Context, req draft.ScheduleRequest, scheduledFor time.Time, teamID, userID string) (string, error) {
	targetProductsJSON, _ := json.Marshal(req.TargetProducts)
	now := helper.ParseWIBTime(time.Now().Format(time.RFC3339))

	var draftID string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO drafts (
			title, topic, article, image_url, image_prompt,
			target_products, has_image, status, scheduled_for,
			created_at, created_by, team_id, user_id
		) VALUES ($1,$2,$3,$4,$5,$6::jsonb,$7,$8,$9,$10,$11,$12,$13)
		RETURNING id
	`,
		req.Title, req.Topic, req.Article, req.ImageURL, req.ImagePrompt,
		string(targetProductsJSON), req.HasImage, "scheduled",
		scheduledFor, now, userID, teamID, userID,
	).Scan(&draftID)

	return draftID, err
}

func (r *RepositoryImpl) InsertHistory(ctx context.Context, req draft.PublishHistoryRequest, userID, teamID, action string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	targetProductsJSON, _ := json.Marshal(req.TargetProducts)

	status := "published"
	if action == "failed" {
		status = "failed"
	}

	query := `
		INSERT INTO histories (
			title, topic, content, image_url, target_products,
			status, action, published_at, created_by, team_id, created_at,seo_score
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`

	now := helper.ParseWIBTime(time.Now().Format(time.RFC3339))

	var historyID string
	err = tx.QueryRowContext(ctx, query,
		req.Title, req.Topic, req.Article, req.ImageURL,
		targetProductsJSON, status, action, now,
		userID, teamID, now, req.SEOScore,
	).Scan(&historyID)
	if err != nil {
		return fmt.Errorf("failed to insert history: %w", err)
	}

	if err := keywords.UpdateKeywords(ctx, tx, keywords.HistorySource{HistoryID: historyID}, req.Keywords); err != nil {
		return fmt.Errorf("failed to update keywords: %w", err)
	}

	r.invalidateCache(ctx,
		fmt.Sprintf("dashboard_stats:%s", teamID))

	if err := r.InvalidateDraftCacheByTeam(ctx, teamID); err != nil {
		log.Printf("failed to invalidate draft cache: %v", err)
	}

	return tx.Commit()
}

// Helper methods
func (r *RepositoryImpl) buildUpdateQuery(id string, data map[string]interface{}) (string, []interface{}, error) {
	log.Println("========== BUILD UPDATE QUERY ==========")

	if len(data) == 0 {
		log.Println("[ERROR] No data to update")
		return "", nil, fmt.Errorf("no data to update")
	}

	log.Printf("[INFO] Building UPDATE query for draft ID: %s", id)
	log.Printf("[INFO] Number of fields to update: %d", len(data))

	// Log semua data yang akan diupdate
	log.Println("[INFO] Data to update:")
	for column, value := range data {
		log.Printf("  %s => %#v (type: %T)", column, value, value)
	}

	setClauses := make([]string, 0, len(data))
	args := make([]interface{}, 0, len(data)+1)
	i := 1

	for column, value := range data {
		log.Printf("[INFO] Processing column: %s, value: %#v (type: %T)", column, value, value)

		// 🔥 CEK KHUSUS: Jika value adalah slice, log warning
		if isSlice(value) {
			log.Printf("[WARNING] Column '%s' contains a slice/array! This might cause SQL error if column doesn't support array type.", column)
			log.Printf("[WARNING] Slice content: %#v", value)
		}

		// 🔥 CEK KHUSUS: Jika value adalah JSON bytes
		if isJSONBytes(value) {
			log.Printf("[INFO] Column '%s' contains JSON bytes, will be stored as JSONB", column)
		}

		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, i))
		args = append(args, value)
		log.Printf("[INFO] Added to SET clause: %s = $%d", column, i)
		i++
	}

	args = append(args, id)
	log.Printf("[INFO] Added WHERE clause: id = $%d", i)

	query := fmt.Sprintf("UPDATE drafts SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "), i)

	log.Println("==================================================")
	log.Printf("[INFO] Generated Query: %s", query)
	log.Printf("[INFO] Query Arguments: %+v", args)

	// Log tipe data setiap argumen
	log.Println("[INFO] Arguments type check:")
	for idx, arg := range args {
		log.Printf("  $%d => %#v (type: %T)", idx+1, arg, arg)
	}

	log.Println("[SUCCESS] buildUpdateQuery completed")
	log.Println("==================================================")

	return query, args, nil
}

// Helper function untuk cek apakah value adalah slice/array
func isSlice(value interface{}) bool {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Slice, reflect.Array:
		return true
	default:
		return false
	}
}

// Helper function untuk cek apakah value adalah JSON bytes
func isJSONBytes(value interface{}) bool {
	if b, ok := value.([]byte); ok {
		// Cek apakah ini JSON valid
		var js json.RawMessage
		if json.Unmarshal(b, &js) == nil {
			return true
		}
	}
	return false
}

func nullIfEmpty(id string) interface{} {
	if id == "" || id == "00000000-0000-0000-0000-000000000000" {
		return nil
	}
	return id
}

func (r *RepositoryImpl) invalidateCache(ctx context.Context, keys ...string) {
	if len(keys) == 0 {
		return
	}

	if err := r.redisClient.Del(ctx, keys...).Err(); err != nil {
		log.Printf("[Cache] failed to delete keys | keys=%v | err=%v", keys, err)
		return
	}

	log.Printf("[Cache] keys deleted | keys=%v", keys)
}

func (r *RepositoryImpl) InvalidateDraftCacheByTeam(
	ctx context.Context,
	teamID string,
) error {
	pattern := fmt.Sprintf("draft_list:*:%s:*", teamID)

	iter := r.redisClient.Scan(ctx, 0, pattern, 0).Iterator()

	deleted := 0

	for iter.Next(ctx) {
		key := iter.Val()

		if err := r.redisClient.Del(ctx, key).Err(); err != nil {
			log.Printf(
				"[InvalidateDraftCacheByTeam] failed to delete cache | team_id=%s | key=%s | err=%v",
				teamID,
				key,
				err,
			)
			continue
		}

		deleted++
	}

	if err := iter.Err(); err != nil {
		return err
	}

	log.Printf(
		"[InvalidateDraftCacheByTeam] cache invalidated | team_id=%s | deleted=%d",
		teamID,
		deleted,
	)

	return nil
}
