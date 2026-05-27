// internal/infrastructure/repository/draft/repository.go
package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"seo-backend/internal/domain/draft"
	"seo-backend/internal/domain/paginate"
	"seo-backend/internal/helper"
	"seo-backend/internal/helper/keywords"

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

	// Parse target_products JSON
	if len(targetProductsJSON) > 0 {

		if err := json.Unmarshal(
			targetProductsJSON,
			&d.TargetProducts,
		); err != nil {

			log.Printf(
				"Error unmarshaling target_products: %v",
				err,
			)
		}
	}

	log.Printf("TARGET PRODUCTS => %+v", d.TargetProducts)

	// Get keywords
	keywords, err := keywords.GetKeywords(ctx, r.db, keywords.DraftSource{DraftID: id})

	if err != nil {

		log.Printf(
			"Error getting keywords for draft %s: %v",
			id,
			err,
		)

		// fallback empty array
		d.Keywords = []string{}

	} else {

		d.Keywords = keywords
	}

	log.Printf("KEYWORDS => %+v", d.Keywords)

	log.Printf("FINAL DRAFT RESPONSE => %+v", d)

	return &d, nil
}

func (r *RepositoryImpl) GetAll(ctx context.Context, teamID string, params paginate.PaginationParams) (*paginate.PaginatedResult[draft.Draft], error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}

	// Cache key berdasarkan semua parameter
	cacheKey := fmt.Sprintf("draft_list:%s:limit=%d:offset=%d:search=%s", teamID, params.Limit, params.Offset, params.Search)

	// Cek Redis cache
	cached, err := r.redisClient.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var result paginate.PaginatedResult[draft.Draft]
		if err := json.Unmarshal(cached, &result); err == nil {
			log.Printf("[GetAll] cache hit | teamID=%s | key=%s", teamID, cacheKey)
			return &result, nil
		}
	}

	log.Printf("[GetAll] cache miss | teamID=%s | key=%s", teamID, cacheKey)

	// Build dynamic WHERE clause
	args := []any{teamID}
	searchClause := ""
	if params.Search != "" {
		args = append(args, "%"+params.Search+"%")
		searchClause = fmt.Sprintf(" AND (title ILIKE $%d OR topic ILIKE $%d)", len(args), len(args))
	}

	// Count query
	var totalItems int
	countQuery := fmt.Sprintf(`
        SELECT COUNT(*) 
        FROM drafts
        WHERE team_id = $1 AND status != 'scheduled'%s
    `, searchClause)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(params.Limit)))
	currentPage := (params.Offset / params.Limit) + 1

	// Append LIMIT & OFFSET
	args = append(args, params.Limit, params.Offset)
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
			seo_score,
            COALESCE(image_prompt, '')
        FROM drafts
        WHERE team_id = $1 AND status != 'scheduled'%s
        ORDER BY created_at DESC
        LIMIT $%d OFFSET $%d
    `, searchClause, len(args)-1, len(args))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("[GetAll] query error | teamID=%s | err=%v", teamID, err)
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
			&d.SeoScore,
			&d.ImagePrompt,
		)
		if err != nil {
			log.Printf("[GetAll] scan error | teamID=%s | err=%v", teamID, err)
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
		log.Printf("[GetAll] failed to marshal result | teamID=%s | err=%v", teamID, err)
	} else {
		if err := r.redisClient.Set(ctx, cacheKey, resultBytes, 2*time.Minute).Err(); err != nil {
			log.Printf("[GetAll] failed to set cache | teamID=%s | err=%v", teamID, err)
		} else {
			log.Printf("[GetAll] cache saved | teamID=%s | key=%s | ttl=2m", teamID, cacheKey)
		}
	}

	return result, nil
}

func (r *RepositoryImpl) GetDashboardStats(ctx context.Context, teamID string) (*draft.DraftStats, error) {
	cacheKey := fmt.Sprintf("dashboard_stats:%s", teamID)

	// 1. Cek Redis cache dulu
	cached, err := r.redisClient.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var stats draft.DraftStats
		if err := json.Unmarshal(cached, &stats); err == nil {
			log.Printf("[DashboardStats] cache hit | teamID=%s", teamID)
			return &stats, nil
		}
	}

	log.Printf("[DashboardStats] cache miss | teamID=%s", teamID)

	stats := &draft.DraftStats{
		ProductCoverage: make(map[string]int),
	}

	// 2. Total draft, with/without image, scheduled
	summaryQuery := `
		SELECT
			COUNT(*)                                        AS total_draft,
			COUNT(*) FILTER (WHERE has_image = true)        AS total_with_image,
			COUNT(*) FILTER (WHERE has_image = false)       AS total_without_image,
			COUNT(*) FILTER (WHERE status = 'scheduled')    AS total_scheduled
		FROM drafts
		WHERE team_id = $1
	`
	if err := r.db.QueryRowContext(ctx, summaryQuery, teamID).Scan(
		&stats.TotalDraft,
		&stats.TotalWithImage,
		&stats.TotalWithoutImage,
		&stats.TotalScheduled,
	); err != nil {
		log.Printf("[DashboardStats] failed to scan summary stats | teamID=%s | err=%v", teamID, err)
		return nil, fmt.Errorf("failed to get summary stats: %w", err)
	}
	log.Printf("[DashboardStats] summary | teamID=%s | total=%d with_image=%d without_image=%d scheduled=%d",
		teamID, stats.TotalDraft, stats.TotalWithImage, stats.TotalWithoutImage, stats.TotalScheduled)

	// 3. Product coverage
	productQuery := `
		SELECT p.name, COUNT(*) AS count
		FROM (
			SELECT id, target_products
			FROM drafts
			WHERE team_id = $1
			AND target_products IS NOT NULL
			AND jsonb_typeof(target_products) = 'array'
		) filtered_drafts,
		jsonb_array_elements_text(filtered_drafts.target_products) AS product_id
		JOIN products p ON p.id = product_id::uuid
		WHERE p.team_id = $1
		GROUP BY p.id, p.name
		ORDER BY count DESC
	`
	productRows, err := r.db.QueryContext(ctx, productQuery, teamID)
	if err != nil {
		log.Printf("[DashboardStats] failed to query product coverage | teamID=%s | err=%v", teamID, err)
		return nil, fmt.Errorf("failed to get product coverage: %w", err)
	}
	defer productRows.Close()

	for productRows.Next() {
		var product string
		var count int
		if err := productRows.Scan(&product, &count); err != nil {
			log.Printf("[DashboardStats] failed to scan product coverage | teamID=%s | err=%v", teamID, err)
			return nil, fmt.Errorf("failed to scan product coverage: %w", err)
		}
		stats.ProductCoverage[product] = count
	}
	if err := productRows.Err(); err != nil {
		log.Printf("[DashboardStats] product rows iteration error | teamID=%s | err=%v", teamID, err)
		return nil, err
	}
	log.Printf("[DashboardStats] product coverage | teamID=%s | total_products=%d", teamID, len(stats.ProductCoverage))

	// 4. Daily activity — 30 hari terakhir
	activityQuery := `
		SELECT
			TO_CHAR(created_at, 'YYYY-MM-DD') AS date,
			COUNT(*) AS count
		FROM drafts
		WHERE team_id = $1
			AND created_at >= NOW() - INTERVAL '30 days'
		GROUP BY date
		ORDER BY date ASC
	`
	activityRows, err := r.db.QueryContext(ctx, activityQuery, teamID)
	if err != nil {
		log.Printf("[DashboardStats] failed to query daily activity | teamID=%s | err=%v", teamID, err)
		return nil, fmt.Errorf("failed to get daily activity: %w", err)
	}
	defer activityRows.Close()

	for activityRows.Next() {
		var a draft.DailyActivity
		if err := activityRows.Scan(&a.Date, &a.Count); err != nil {
			log.Printf("[DashboardStats] failed to scan daily activity | teamID=%s | err=%v", teamID, err)
			return nil, fmt.Errorf("failed to scan daily activity: %w", err)
		}
		stats.DailyActivity = append(stats.DailyActivity, a)
	}
	if err := activityRows.Err(); err != nil {
		log.Printf("[DashboardStats] activity rows iteration error | teamID=%s | err=%v", teamID, err)
		return nil, err
	}
	log.Printf("[DashboardStats] daily activity | teamID=%s | total_days=%d", teamID, len(stats.DailyActivity))

	// 5. Simpan ke Redis cache, TTL 5 menit
	statsBytes, err := json.Marshal(stats)
	if err != nil {
		log.Printf("[DashboardStats] failed to marshal stats for cache | teamID=%s | err=%v", teamID, err)
	} else {
		if err := r.redisClient.Set(ctx, cacheKey, statsBytes, 5*time.Minute).Err(); err != nil {
			log.Printf("[DashboardStats] failed to set cache | teamID=%s | err=%v", teamID, err)
		} else {
			log.Printf("[DashboardStats] cache saved | teamID=%s | ttl=5m", teamID)
		}
	}

	log.Printf("[DashboardStats] completed | teamID=%s", teamID)
	return stats, nil
}

func (r *RepositoryImpl) GetAllScheduled(ctx context.Context, teamID string, params paginate.PaginationParams) (*paginate.PaginatedResult[draft.Draft], error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}

	var totalItems int
	countQuery := `
		SELECT COUNT(*) 
		FROM drafts
		WHERE team_id = $1 AND status = 'scheduled'
	`
	if err := r.db.QueryRowContext(ctx, countQuery, teamID).Scan(&totalItems); err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(params.Limit)))
	currentPage := (params.Offset / params.Limit) + 1

	query := `
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
		WHERE team_id = $1 AND status = 'scheduled'
		ORDER BY scheduled_for ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, teamID, params.Limit, params.Offset)
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

	// Start transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert draft
	query := `INSERT INTO drafts (
		title, topic, article, image_url, image_prompt,
		status, target_products, has_image, created_by, team_id, user_id, slug
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	RETURNING id`

	var draftID string
	err = tx.QueryRowContext(
		ctx,
		query,
		req.Title, req.Topic, req.Article, req.ImageURL, req.ImagePrompt,
		"draft", targetProductsJSON, req.HasImage, nullIfEmpty(userID), nullIfEmpty(teamID), nullIfEmpty(userID), req.Slug,
	).Scan(&draftID)

	if err != nil {
		log.Printf("Error creating draft: %v", err)
		return "", err
	}

	// Insert keywords with duplicate checking
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

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.invalidateCache(ctx,
		fmt.Sprintf("dashboard_stats:%s", teamID),
		fmt.Sprintf("draft_list:%s", teamID),
	)

	return draftID, nil
}

func (r *RepositoryImpl) Update(ctx context.Context, TeamID string, id string, data map[string]interface{}) error {
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
		fmt.Sprintf("dashboard_stats:%s", TeamID),
		fmt.Sprintf("draft_list:%s", TeamID),
	)

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
		fmt.Sprintf("dashboard_stats:%s", TeamID),
		fmt.Sprintf("draft_list:%s", TeamID),
	)
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
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
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

	return tx.Commit()
}

// Helper methods
func (r *RepositoryImpl) buildUpdateQuery(id string, data map[string]interface{}) (string, []interface{}, error) {
	if len(data) == 0 {
		return "", nil, fmt.Errorf("no data to update")
	}

	setClauses := make([]string, 0, len(data))
	args := make([]interface{}, 0, len(data)+1)
	i := 1

	for column, value := range data {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, i))
		args = append(args, value)
		i++
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE drafts SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "), i)

	return query, args, nil
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
