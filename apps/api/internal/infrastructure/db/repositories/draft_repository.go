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
	"regexp"
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
		slug,
		excerpt,
		status
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
		&d.Slug,
		&d.Excerpt,
		&d.Status,
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
	fullWhere := whereClause + searchClause + " AND status != 'scheduled' AND status != 'failed' AND status != 'published'"

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
            image_url, status, COALESCE(image_prompt, '')
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
			&d.Status, &d.ImagePrompt,
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

// CREATE - tambahin SEO score calculation
func (r *RepositoryImpl) Create(ctx context.Context, req draft.CreateDraftRequest, userID, teamID string) (string, error) {
	targetProductsJSON, _ := json.Marshal(req.TargetProducts)

	// HITUNG SEO SCORE PAKE CalculateSEOScore
	seoScore := draft.CalculateSEOScore(req.Title, req.Article, req.Topic, req.Excerpt, req.Keywords).Total

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	slug := req.Slug
	if slug == "" {
		slug = generateSlug(req.Title)
		slug, err = r.generateUniqueSlug(ctx, slug, "")
		if err != nil {
			return "", fmt.Errorf("failed to generate unique slug: %w", err)
		}
	} else {
		slug = cleanSlug(slug)
		exists, err := r.slugExists(ctx, slug, "")
		if err != nil {
			return "", fmt.Errorf("failed to check slug existence: %w", err)
		}
		if exists {
			slug, err = r.generateUniqueSlug(ctx, slug, "")
			if err != nil {
				return "", fmt.Errorf("failed to generate unique slug: %w", err)
			}
		}
	}

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
		"draft", targetProductsJSON, req.HasImage, nullIfEmpty(userID), nullIfEmpty(teamID), nullIfEmpty(userID), slug, req.Excerpt, seoScore,
	).Scan(&draftID)

	if err != nil {
		log.Printf("Error creating draft: %v", err)
		return "", err
	}

	if len(req.Keywords) > 0 {
		uniqueKeywords := keywords.UniqueStrings(req.Keywords)
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

		err = keywords.InsertKeywordsWithDuplicateCheck(ctx, tx, draftID, keywordEntities)
		if err != nil {
			log.Printf("FAILED INSERT KEYWORDS => %v", err)
			return "", fmt.Errorf("failed to insert keywords: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.invalidateDashboardCache(ctx, teamID)

	if err := r.InvalidateDraftCacheByTeam(ctx, teamID); err != nil {
		log.Printf("failed to invalidate draft cache: %v", err)
	}

	return draftID, nil
}

// UPDATE - tambahin SEO score calculation
func (r *RepositoryImpl) Update(ctx context.Context, id string, TeamID string, data draft.CreateDraftRequest) error {
	// HITUNG SEO SCORE PAKE CalculateSEOScore
	data.SEOScore = draft.CalculateSEOScore(data.Title, data.Article, data.Topic, data.Excerpt, data.Keywords).Total

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	data.UpdateAt = helper.ParseWIBTime(time.Now().Format(time.RFC3339))

	query, args, err := r.buildUpdateQuery(id, data)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update draft: %w", err)
	}

	if data.Keywords != nil {
		if err := keywords.UpdateKeywords(ctx, tx, keywords.DraftSource{DraftID: id}, data.Keywords); err != nil {
			return fmt.Errorf("failed to update keywords: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.invalidateDashboardCache(ctx, TeamID)

	if err := r.InvalidateDraftCacheByTeam(ctx, TeamID); err != nil {
		log.Printf("failed to invalidate draft cache: %v", err)
	}

	return nil
}

// DASHBOARD - tetep baca dari DB, GAUSAH diitung ulang
func (r *RepositoryImpl) GetDashboardStats(ctx context.Context, filter models.UserContext) (*draft.DraftStats, error) {
	startTime := time.Now()

	whereClause, whereArgs := userRole.BuildAccessFilter(filter)

	cacheKey := fmt.Sprintf("dashboard_stats:v2:%s", filter.GetRole())
	if filter.GetTeamID() != "" {
		cacheKey = fmt.Sprintf("dashboard_stats:v2:%s:%s", filter.GetRole(), filter.GetTeamID())
	}

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
		ProductCoverage:      make(map[string]int),
		ProductStatus:        make(map[string]int),
		StatusBreakdown:      make(map[string]int),
		TopicBreakdown:       make(map[string]int),
		SEOScoreDistribution: make(map[string]int),
	}

	summaryQuery := fmt.Sprintf(`
		SELECT
			COUNT(DISTINCT d.id)                                        AS total_draft,
			COUNT(DISTINCT d.id) FILTER (WHERE d.has_image = true)     AS total_with_image,
			COUNT(DISTINCT d.id) FILTER (WHERE d.has_image = false)    AS total_without_image,
			COUNT(DISTINCT d.id) FILTER (WHERE d.status = 'scheduled') AS total_scheduled,
			COUNT(DISTINCT d.id) FILTER (WHERE d.status = 'published') AS total_published,
			COUNT(DISTINCT d.id) FILTER (WHERE k.id IS NOT NULL)       AS total_with_keywords,
			COUNT(DISTINCT d.id) FILTER (WHERE d.seo_score > 0)        AS total_with_seo,
			COALESCE(AVG(d.seo_score), 0)                               AS avg_seo_score
		FROM drafts d
		LEFT JOIN keywords k ON d.id = k.id_draft
		WHERE %s
	`, whereClause)

	if err := r.db.QueryRowContext(ctx, summaryQuery, whereArgs...).Scan(
		&stats.TotalDraft,
		&stats.TotalWithImage,
		&stats.TotalWithoutImage,
		&stats.TotalScheduled,
		&stats.TotalPublished,
		&stats.TotalWithKeywords,
		&stats.TotalWithSEO,
		&stats.SEOScoreAvg,
	); err != nil {
		log.Printf("[DashboardStats] failed to scan summary stats | filter=%+v | err=%v", filter, err)
		return nil, fmt.Errorf("failed to get summary stats: %w", err)
	}

	// Status breakdown
	statusQuery := fmt.Sprintf(`
		SELECT d.status, COUNT(DISTINCT d.id) as count
		FROM drafts d
		WHERE %s
		GROUP BY d.status
	`, whereClause)

	statusRows, err := r.db.QueryContext(ctx, statusQuery, whereArgs...)
	if err == nil {
		defer statusRows.Close()
		for statusRows.Next() {
			var status string
			var count int
			if err := statusRows.Scan(&status, &count); err == nil {
				stats.StatusBreakdown[status] = count
			}
		}
	}

	// Product coverage
	productQuery := fmt.Sprintf(`
		SELECT p.name, COUNT(DISTINCT filtered_drafts.id) AS count, p.status as product_status
		FROM (
			SELECT id, target_products
			FROM drafts
			WHERE %s
			AND target_products IS NOT NULL
			AND jsonb_typeof(target_products) = 'array'
		) filtered_drafts,
		jsonb_array_elements_text(filtered_drafts.target_products) AS product_id
		JOIN products p ON p.id = product_id::uuid
		GROUP BY p.id, p.name, p.status
		ORDER BY count DESC
		LIMIT 100
	`, whereClause)

	productRows, err := r.db.QueryContext(ctx, productQuery, whereArgs...)
	if err == nil {
		defer productRows.Close()
		for productRows.Next() {
			var productName string
			var count int
			var productStatus string
			if err := productRows.Scan(&productName, &count, &productStatus); err == nil {
				stats.ProductCoverage[productName] = count
				stats.ProductStatus[productStatus]++
			}
		}
	}

	// Topic breakdown
	topicQuery := fmt.Sprintf(`
		SELECT 
			COALESCE(d.topic, 'Uncategorized') as topic,
			COUNT(DISTINCT d.id) as count,
			COALESCE(AVG(d.seo_score), 0) as avg_seo
		FROM drafts d
		WHERE %s
		GROUP BY d.topic
		ORDER BY count DESC
		LIMIT 20
	`, whereClause)

	topicRows, err := r.db.QueryContext(ctx, topicQuery, whereArgs...)
	if err == nil {
		defer topicRows.Close()
		stats.TopTopics = make([]draft.TopicStats, 0)
		for topicRows.Next() {
			var topic draft.TopicStats
			if err := topicRows.Scan(&topic.Topic, &topic.Count, &topic.AvgSEO); err == nil {
				stats.TopicBreakdown[topic.Topic] = topic.Count
				stats.TopTopics = append(stats.TopTopics, topic)
			}
		}
	}

	// SEO score distribution
	seoQuery := fmt.Sprintf(`
		SELECT 
			CASE 
				WHEN d.seo_score = 0 THEN '0'
				WHEN d.seo_score BETWEEN 1 AND 20 THEN '1-20'
				WHEN d.seo_score BETWEEN 21 AND 40 THEN '21-40'
				WHEN d.seo_score BETWEEN 41 AND 60 THEN '41-60'
				WHEN d.seo_score BETWEEN 61 AND 80 THEN '61-80'
				WHEN d.seo_score BETWEEN 81 AND 100 THEN '81-100'
				ELSE 'unknown'
			END as score_range,
			COUNT(DISTINCT d.id) as count
		FROM drafts d
		WHERE %s
		GROUP BY score_range
		ORDER BY MIN(d.seo_score)
	`, whereClause)

	seoRows, err := r.db.QueryContext(ctx, seoQuery, whereArgs...)
	if err == nil {
		defer seoRows.Close()
		for seoRows.Next() {
			var scoreRange string
			var count int
			if err := seoRows.Scan(&scoreRange, &count); err == nil {
				stats.SEOScoreDistribution[scoreRange] = count
			}
		}
	}

	// Daily activity
	activityQuery := fmt.Sprintf(`
		SELECT
			TO_CHAR(d.created_at, 'YYYY-MM-DD') AS date,
			COUNT(DISTINCT d.id) AS total,
			COUNT(DISTINCT d.id) FILTER (WHERE d.status = 'scheduled') AS scheduled,
			COUNT(DISTINCT d.id) FILTER (WHERE d.status = 'published') AS published,
			COUNT(DISTINCT d.id) FILTER (WHERE d.has_image = true) AS with_image,
			COUNT(DISTINCT d.id) FILTER (WHERE k.id IS NOT NULL) AS with_keywords,
			COALESCE(AVG(d.seo_score), 0) AS avg_seo
		FROM drafts d
		LEFT JOIN keywords k ON d.id = k.id_draft
		WHERE %s
			AND d.created_at >= NOW() - INTERVAL '30 days'
		GROUP BY date
		ORDER BY date ASC
	`, whereClause)

	activityRows, err := r.db.QueryContext(ctx, activityQuery, whereArgs...)
	if err == nil {
		defer activityRows.Close()
		for activityRows.Next() {
			var a draft.DailyActivity
			if err := activityRows.Scan(&a.Date, &a.Count, &a.Scheduled, &a.Published, &a.WithImage, &a.WithKeywords, &a.AvgSEO); err == nil {
				stats.DailyActivity = append(stats.DailyActivity, a)
			}
		}
	}

	// Keywords stats
	keywordsQuery := fmt.Sprintf(`
		SELECT 
			k.name as keyword,
			COUNT(DISTINCT k.id_draft) as usage_count
		FROM keywords k
		INNER JOIN drafts d ON k.id_draft = d.id
		WHERE %s
		GROUP BY k.name
		ORDER BY usage_count DESC
		LIMIT 20
	`, whereClause)

	keywordsRows, err := r.db.QueryContext(ctx, keywordsQuery, whereArgs...)
	if err == nil {
		defer keywordsRows.Close()
		stats.TopKeywords = make([]draft.KeywordStats, 0)
		for keywordsRows.Next() {
			var kw draft.KeywordStats
			if err := keywordsRows.Scan(&kw.Keyword, &kw.Count); err == nil {
				stats.TopKeywords = append(stats.TopKeywords, kw)
			}
		}

		avgKeywordsQuery := fmt.Sprintf(`
			SELECT COALESCE(AVG(keyword_count), 0) as avg_keywords
			FROM (
				SELECT d.id, COUNT(k.id) as keyword_count
				FROM drafts d
				LEFT JOIN keywords k ON d.id = k.id_draft
				WHERE %s
				GROUP BY d.id
			) draft_keywords
		`, whereClause)

		r.db.QueryRowContext(ctx, avgKeywordsQuery, whereArgs...).Scan(&stats.KeywordsAvgCount)
	}

	// Scheduled upcoming
	scheduledQuery := fmt.Sprintf(`
		SELECT 
			d.id, d.title, d.scheduled_for, d.target_products,
			COUNT(k.id) as keyword_count
		FROM drafts d
		LEFT JOIN keywords k ON d.id = k.id_draft
		WHERE %s 
			AND d.status = 'scheduled'
			AND d.scheduled_for IS NOT NULL
			AND d.scheduled_for >= NOW()
		GROUP BY d.id, d.title, d.scheduled_for, d.target_products
		ORDER BY d.scheduled_for ASC
		LIMIT 10
	`, whereClause)

	scheduledRows, err := r.db.QueryContext(ctx, scheduledQuery, whereArgs...)
	if err == nil {
		defer scheduledRows.Close()
		stats.ScheduledUpcoming = make([]draft.ScheduledItem, 0)
		for scheduledRows.Next() {
			var item draft.ScheduledItem
			var products []string
			var keywordCount int
			if err := scheduledRows.Scan(&item.ID, &item.Title, &item.ScheduledFor, &products, &keywordCount); err == nil {
				item.Products = products
				stats.ScheduledUpcoming = append(stats.ScheduledUpcoming, item)
			}
		}
	}

	stats.CalculateDerivedMetrics()

	stats.CacheMetadata = draft.CacheMetadata{
		CachedAt:   time.Now(),
		TTL:        "5m",
		Generation: float64(time.Since(startTime).Milliseconds()),
	}

	statsBytes, err := json.Marshal(stats)
	if err == nil {
		ttl := 5 * time.Minute
		if stats.TotalDraft > 1000 {
			ttl = 10 * time.Minute
		}
		r.redisClient.Set(ctx, cacheKey, statsBytes, ttl)
	}

	log.Printf("[DashboardStats] completed | filter=%+v | duration=%v", filter, time.Since(startTime))
	return stats, nil
}

// Update cache invalidation di Create, Update, Delete functions
func (r *RepositoryImpl) invalidateDashboardCache(ctx context.Context, teamID string) {
	// Invalidate semua dashboard cache untuk team ini
	patterns := []string{
		fmt.Sprintf("dashboard_stats:v2:*:%s", teamID),
		"dashboard_stats:v2:*", // Invalidate semua role-based cache juga
	}

	for _, pattern := range patterns {
		iter := r.redisClient.Scan(ctx, 0, pattern, 0).Iterator()
		deleted := 0

		for iter.Next(ctx) {
			key := iter.Val()
			if err := r.redisClient.Del(ctx, key).Err(); err != nil {
				log.Printf("[Cache] failed to delete key | key=%s | err=%v", key, err)
				continue
			}
			deleted++
		}

		if err := iter.Err(); err != nil {
			log.Printf("[Cache] scan error | pattern=%s | err=%v", pattern, err)
		}

		if deleted > 0 {
			log.Printf("[Cache] dashboard cache invalidated | pattern=%s | deleted=%d", pattern, deleted)
		}
	}
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

// generateSlug generates a URL slug from title
func generateSlug(title string) string {
	if title == "" {
		return "untitled"
	}

	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters, keep only alphanumeric and dash
	var result strings.Builder
	for _, ch := range slug {
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-' {
			result.WriteRune(ch)
		}
	}

	// Get the cleaned string
	slug = result.String()

	// Remove multiple consecutive dashes
	re := regexp.MustCompile(`-+`)
	slug = re.ReplaceAllString(slug, "-")

	// Trim dashes from start and end
	slug = strings.Trim(slug, "-")

	// Limit length to 100 characters
	if len(slug) > 100 {
		slug = slug[:100]
		// Trim any trailing dash after truncation
		slug = strings.TrimSuffix(slug, "-")
	}

	return slug
}

// cleanSlug cleans and validates a slug
func cleanSlug(slug string) string {
	if slug == "" {
		return "untitled"
	}

	// Convert to lowercase
	slug = strings.ToLower(slug)

	// Replace spaces and special chars with hyphens
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug = re.ReplaceAllString(slug, "-")

	// Remove multiple consecutive dashes
	re = regexp.MustCompile(`-+`)
	slug = re.ReplaceAllString(slug, "-")

	// Trim dashes from start and end
	slug = strings.Trim(slug, "-")

	// Limit length
	if len(slug) > 100 {
		slug = slug[:100]
		slug = strings.TrimSuffix(slug, "-")
	}

	if slug == "" {
		return "untitled"
	}

	return slug
}

// generateUniqueSlug generates a unique slug by adding a number suffix if needed
func (r *RepositoryImpl) generateUniqueSlug(ctx context.Context, baseSlug, excludeID string) (string, error) {
	slug := baseSlug
	counter := 1

	for {
		exists, err := r.slugExists(ctx, slug, excludeID)
		if err != nil {
			return "", err
		}

		if !exists {
			return slug, nil
		}

		// Add counter suffix
		slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++

		// Safety limit to prevent infinite loop
		if counter > 1000 {
			return "", fmt.Errorf("unable to generate unique slug after 1000 attempts")
		}
	}
}

// slugExists checks if a slug already exists in drafts or published articles
func (r *RepositoryImpl) slugExists(ctx context.Context, slug, excludeID string) (bool, error) {
	var count int

	query := `
		SELECT COUNT(*) FROM drafts WHERE slug = $1
	`
	args := []interface{}{slug}

	if excludeID != "" {
		query += " AND id != $2"
		args = append(args, excludeID)
	}

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return count > 0, nil
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

// Helper methods
func (r *RepositoryImpl) buildUpdateQuery(id string, data draft.CreateDraftRequest) (string, []interface{}, error) {
	log.Println("========== BUILD UPDATE QUERY ==========")

	// Build update data directly
	setClauses := make([]string, 0)
	args := make([]interface{}, 0)
	i := 1

	// Helper function to add field
	addField := func(column string, value interface{}) {
		if value != nil {
			// Skip empty strings for optional fields
			if str, ok := value.(string); ok && str == "" {
				return
			}
			// Skip zero time
			if t, ok := value.(time.Time); ok && t.IsZero() {
				return
			}
			// Handle pointer to string
			if ptr, ok := value.(*string); ok && ptr != nil && *ptr != "" {
				value = *ptr
			}
			// Handle pointer to time
			if ptr, ok := value.(*time.Time); ok && ptr != nil && !(*ptr).IsZero() {
				value = *ptr
			}
			// Handle bool (always include)
			if _, ok := value.(bool); ok {
				// Booleans are always included
			}

			setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, i))
			args = append(args, value)
			i++
		}
	}

	// Add all fields from CreateDraftRequest
	addField("title", data.Title)
	addField("topic", data.Topic)
	addField("article", data.Article)
	addField("image_url", data.ImageURL)
	addField("slug", data.Slug)
	addField("excerpt", data.Excerpt)
	addField("status", data.Status)
	addField("team_id", data.TeamID)
	addField("user_id", data.UserID)
	addField("created_by", data.CreatedBy)
	addField("created_at", time.Now())

	// Handle target_products (JSONB)
	if len(data.TargetProducts) > 0 {
		jsonValue, err := json.Marshal(data.TargetProducts)
		if err != nil {
			log.Printf("[ERROR] Failed to marshal target_products: %v", err)
			return "", nil, fmt.Errorf("failed to marshal target_products: %w", err)
		}
		addField("target_products", jsonValue)
	}

	// Handle scheduled_for (string)
	if data.ScheduledFor != "" {
		addField("scheduled_for", data.ScheduledFor)
	}

	// Always update has_image, seo_score, and updated_at
	addField("has_image", data.HasImage)
	addField("seo_score", data.SEOScore)
	addField("updated_at", time.Now())

	// Check if we have any data to update
	if len(setClauses) == 0 {
		log.Println("[ERROR] No data to update")
		return "", nil, fmt.Errorf("no data to update")
	}

	// Add WHERE clause
	args = append(args, id)

	query := fmt.Sprintf("UPDATE drafts SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "), i)

	log.Printf("[INFO] Generated Query: %s", query)
	log.Printf("[INFO] Query Arguments: %+v", args)

	// Log tipe data setiap argumen
	log.Println("[INFO] Arguments type check:")
	for idx, arg := range args {
		log.Printf("  $%d => %#v (type: %T)", idx+1, arg, arg)
	}

	log.Println("[SUCCESS] buildUpdateQuery completed")
	return query, args, nil
}

// Helper function to validate column names (prevent SQL injection)
func isValidColumnName(column string) bool {
	// List of allowed column names
	allowedColumns := map[string]bool{
		"id":              true,
		"title":           true,
		"topic":           true,
		"article":         true,
		"image_url":       true,
		"image_prompt":    true,
		"slug":            true,
		"target_products": true,
		"status":          true,
		"scheduled_for":   true,
		"has_image":       true,
		"excerpt":         true,
		"team_id":         true,
		"user_id":         true,
		"created_by":      true,
		"created_at":      true,
		"updated_at":      true,
	}

	// Allow only alphanumeric characters and underscores
	if !regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`).MatchString(column) {
		return false
	}

	return allowedColumns[column]
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
