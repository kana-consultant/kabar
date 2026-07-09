package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"seo-backend/internal/domain/draft"
	"seo-backend/internal/domain/history"
	"seo-backend/internal/domain/paginate"
	"seo-backend/internal/helper"
	userRole "seo-backend/internal/helper/filter"
	"seo-backend/internal/helper/keywords"
	"seo-backend/internal/models"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type HistoryRepository struct {
	db          *sql.DB
	redisClient redis.Client
}

func NewHistoryRepository(db *sql.DB, redisClient redis.Client) history.HistoryRepository {
	return &HistoryRepository{db: db, redisClient: redisClient}
}

// Create inserts a new history record
func (r *HistoryRepository) Create(ctx context.Context, req draft.PublishHistoryRequest, userID, teamID, action string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	targetProductsJSON, err := json.Marshal(req.TargetProducts)
	if err != nil {
		return fmt.Errorf("failed to marshal target products: %w", err)
	}

	// Hitung SEO score
	seoScore := draft.CalculateSEOScore(req.Title, req.Article, req.Topic, req.Excerpt, req.Keywords).Total

	// Tentukan has_image
	hasImage := req.ImageURL != nil && *req.ImageURL != ""

	query := `
        INSERT INTO drafts (
            id, title, topic, article, image_url, target_products,
            status, published_at, has_image, seo_score,
            created_by, team_id, user_id, created_at, slug, excerpt
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW(), $14, $15)
        RETURNING id
    `

	now := helper.ParseWIBTime(time.Now().Format(time.RFC3339))
	draftID := uuid.New().String()

	err = tx.QueryRowContext(ctx, query,
		draftID,
		req.Title,
		req.Topic,
		req.Article,
		req.ImageURL,
		targetProductsJSON,
		action,
		now,
		hasImage,
		seoScore,
		userID,
		teamID,
		userID,
		req.Slug,
		req.Excerpt,
	).Scan(&draftID)
	if err != nil {
		return fmt.Errorf("failed to insert history: %w", err)
	}

	if len(req.Keywords) > 0 {
		if err := keywords.UpdateKeywords(ctx, tx, keywords.DraftSource{DraftID: draftID}, req.Keywords); err != nil {
			return fmt.Errorf("failed to update keywords: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.invalidateDashboardCache(ctx, teamID)
	if err := r.InvalidateDraftCacheByTeam(ctx, teamID); err != nil {
		log.Printf("failed to invalidate draft cache: %v", err)
	}

	return nil
}

func (r *HistoryRepository) invalidateDashboardCache(ctx context.Context, teamID string) {
	patterns := []string{
		fmt.Sprintf("dashboard_stats:v2:*:%s", teamID),
		"dashboard_stats:v2:*",
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

func (r *HistoryRepository) InvalidateDraftCacheByTeam(
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

// GetByID retrieves history by ID
func (r *HistoryRepository) GetByID(ctx context.Context, id string) (*history.History, error) {
	query := `
		SELECT id, title, slug, topic, article, excerpt, image_url, target_products,
			status, published_at, scheduled_for, has_image,
			created_by, team_id, created_at
		FROM drafts WHERE id = $1
	`

	var h history.History
	var targetProductsJSON []byte
	var createdBy sql.NullString
	var teamID sql.NullString
	var hasImage sql.NullBool
	var excerpt sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&h.ID, &h.Title, &h.Slug, &h.Topic, &h.Content, &excerpt, &h.ImageURL, &targetProductsJSON,
		&h.Status, &h.PublishedAt, &h.ScheduledFor, &hasImage,
		&createdBy, &teamID, &h.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get history: %w", err)
	}

	if len(targetProductsJSON) > 0 {
		json.Unmarshal(targetProductsJSON, &h.TargetProducts)
	}
	if createdBy.Valid {
		h.CreatedBy = &createdBy.String
	}
	if teamID.Valid {
		h.TeamID = &teamID.String
	}
	if hasImage.Valid {
		h.HasImage = hasImage.Bool
	}
	if excerpt.Valid {
		h.Excerpt = excerpt.String
	}

	keywords, err := keywords.GetKeywords(ctx, r.db, keywords.DraftSource{DraftID: id})
	if err != nil {
		h.Keywords = []string{}
	} else {
		h.Keywords = keywords
	}

	return &h, nil
}

func (r *HistoryRepository) GetAll(ctx context.Context, userCtx models.UserContext, params history.HistoryFilter) (*paginate.PaginatedResult[history.History], error) {
	whereClause, whereArgs := userRole.BuildAccessFilter(userCtx)

	var totalItems, totalSuccess, totalFailed int

	countQuery := fmt.Sprintf(`
	SELECT
		COUNT(*) FILTER (WHERE status = 'published') AS total_success,
		COUNT(*) FILTER (WHERE status IN ('published', 'failed')) AS total,
		COUNT(*) FILTER (WHERE status = 'failed') AS total_failed
	FROM drafts
	WHERE %s
`, whereClause)

	if err := r.db.QueryRowContext(ctx, countQuery, whereArgs...).
		Scan(&totalSuccess, &totalItems, &totalFailed); err != nil {
		return nil, fmt.Errorf("failed to count histories: %w", err)
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(params.Limit)))
	currentPage := (params.Offset / params.Limit) + 1

	args := append(whereArgs, params.Limit, params.Offset)
	query := fmt.Sprintf(`
		SELECT id, title, slug, topic, article, excerpt, image_url, target_products,
			status, published_at, scheduled_for, has_image,
			created_by, team_id, created_at
		FROM drafts
		WHERE %s
		AND status IN ('published', 'failed')
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, len(args)-1, len(args))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query history: %w", err)
	}
	defer rows.Close()

	historyData, err := r.scanHistory(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan history: %w", err)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return &paginate.PaginatedResult[history.History]{
		Data:         historyData,
		TotalItems:   totalItems,
		TotalSuccess: totalSuccess,
		TotalFailed:  totalFailed,
		TotalPages:   totalPages,
		CurrentPage:  currentPage,
		Limit:        params.Limit,
		Offset:       params.Offset,
	}, nil
}

// GetAllWithQuery retrieves history with custom query
func (r *HistoryRepository) GetAllWithQuery(ctx context.Context, query string, args []interface{}) ([]history.History, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query history: %w", err)
	}
	defer rows.Close()

	return r.scanHistory(rows)
}

// GetByTeamID retrieves history by team ID
func (r *HistoryRepository) GetByTeamID(ctx context.Context, teamID string) ([]history.History, error) {
	query := `
		SELECT id, title, slug, topic, article, excerpt, image_url, target_products,
			status, published_at, scheduled_for, has_image,
			created_by, team_id, created_at
		FROM drafts WHERE team_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get history by team: %w", err)
	}
	defer rows.Close()

	return r.scanHistory(rows)
}

// GetByCreatedBy retrieves history by creator
func (r *HistoryRepository) GetByCreatedBy(ctx context.Context, createdBy string) ([]history.History, error) {
	query := `
		SELECT id, title, slug, topic, article, excerpt, image_url, target_products,
			status, published_at, scheduled_for, has_image,
			created_by, team_id, created_at
		FROM drafts WHERE created_by = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, createdBy)
	if err != nil {
		return nil, fmt.Errorf("failed to get history by creator: %w", err)
	}
	defer rows.Close()

	return r.scanHistory(rows)
}

// GetByStatus retrieves history by status
func (r *HistoryRepository) GetByStatus(ctx context.Context, status string) ([]history.History, error) {
	query := `
		SELECT id, title, slug, topic, article, excerpt, image_url, target_products,
			status, published_at, scheduled_for, has_image,
			created_by, team_id, created_at
		FROM drafts WHERE status = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get history by status: %w", err)
	}
	defer rows.Close()

	return r.scanHistory(rows)
}

// GetRecentActivity retrieves recent history activity
func (r *HistoryRepository) GetRecentActivity(ctx context.Context, teamID string, limit int) ([]history.History, error) {
	query := `
		SELECT id, title, slug, topic, article, excerpt, image_url, target_products,
			status, published_at, scheduled_for, has_image,
			created_by, team_id, created_at
		FROM drafts 
		WHERE team_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, teamID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent activity: %w", err)
	}
	defer rows.Close()

	return r.scanHistory(rows)
}

// Count returns total count based on filters
func (r *HistoryRepository) Count(ctx context.Context, query history.HistoryFilter) (int, error) {
	countQuery := `SELECT COUNT(*) FROM drafts WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if query.TeamID != "" {
		countQuery += fmt.Sprintf(" AND team_id = $%d", argIndex)
		args = append(args, query.TeamID)
		argIndex++
	}
	if query.Status != "" {
		countQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, query.Status)
		argIndex++
	}
	if query.Topic != "" {
		countQuery += fmt.Sprintf(" AND topic = $%d", argIndex)
		args = append(args, query.Topic)
		argIndex++
	}
	if query.Search != "" {
		searchPattern := "%" + query.Search + "%"
		countQuery += fmt.Sprintf(" AND (title ILIKE $%d OR article ILIKE $%d)", argIndex, argIndex)
		args = append(args, searchPattern)
		argIndex++
	}

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to count history: %w", err)
	}

	return total, nil
}

// GetCountByStatus returns count by status
func (r *HistoryRepository) GetCountByStatus(ctx context.Context, teamID string) (map[string]int, error) {
	query := `
		SELECT status, COUNT(*) 
		FROM drafts 
		WHERE team_id = $1 
		GROUP BY status
	`

	rows, err := r.db.QueryContext(ctx, query, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get count by status: %w", err)
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			continue
		}
		result[status] = count
	}

	return result, nil
}

// Update updates a history record
func (r *HistoryRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	setClauses := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1

	// Auto-set has_image if image_url updated
	if imageURL, ok := updates["image_url"]; ok {
		hasImage := false
		if imageURL != nil {
			if imgStr, ok := imageURL.(string); ok && imgStr != "" {
				hasImage = true
			}
		}
		setClauses = append(setClauses, fmt.Sprintf("has_image = $%d", argIndex))
		args = append(args, hasImage)
		argIndex++
	}

	for key, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, helper.ParseWIBTime(time.Now().Format(time.RFC3339)))
	argIndex++

	args = append(args, id)
	query := fmt.Sprintf("UPDATE drafts SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "), argIndex)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update history: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("history with id %s not found", id)
	}

	return nil
}

// Delete deletes a history record
func (r *HistoryRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM drafts WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete history: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("history with id %s not found", id)
	}

	return nil
}

// DeleteByTeamID deletes all history for a team
func (r *HistoryRepository) DeleteByTeamID(ctx context.Context, teamID string) error {
	query := `DELETE FROM drafts WHERE team_id = $1`

	_, err := r.db.ExecContext(ctx, query, teamID)
	if err != nil {
		return fmt.Errorf("failed to delete history by team: %w", err)
	}

	return nil
}

// DeleteByStatus deletes history by status
func (r *HistoryRepository) DeleteByStatus(ctx context.Context, status string) error {
	query := `DELETE FROM drafts WHERE status = $1`

	_, err := r.db.ExecContext(ctx, query, status)
	if err != nil {
		return fmt.Errorf("failed to delete history by status: %w", err)
	}

	return nil
}

// GetStats returns statistics
func (r *HistoryRepository) GetStats(ctx context.Context, query *history.HistoryFilter) (*history.HistoryStats, error) {
	stats := &history.HistoryStats{}

	total, err := r.Count(ctx, *query)
	if err != nil {
		return nil, err
	}
	stats.Total = total

	statusCount, err := r.GetCountByStatus(ctx, query.TeamID)
	if err != nil {
		return nil, err
	}

	stats.SuccessCount = statusCount["success"]
	stats.FailedCount = statusCount["failed"]
	stats.PublishedCount = statusCount["published"]
	stats.ScheduledCount = statusCount["scheduled"]

	if stats.Total > 0 {
		stats.SuccessRate = float64(stats.SuccessCount) / float64(stats.Total) * 100
	}

	return stats, nil
}

func (r *HistoryRepository) queryHistory(ctx context.Context, query string) ([]history.History, error) {
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query history: %w", err)
	}
	defer rows.Close()

	return r.scanHistory(rows)
}

func (r *HistoryRepository) scanHistory(rows *sql.Rows) ([]history.History, error) {
	var histories []history.History

	for rows.Next() {
		var h history.History
		var targetProductsJSON []byte
		var createdBy sql.NullString
		var teamID sql.NullString
		var hasImage sql.NullBool
		var excerpt sql.NullString

		err := rows.Scan(
			&h.ID, &h.Title, &h.Slug, &h.Topic, &h.Content, &excerpt, &h.ImageURL, &targetProductsJSON,
			&h.Status, &h.PublishedAt, &h.ScheduledFor, &hasImage,
			&createdBy, &teamID, &h.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history: %w", err)
		}

		if len(targetProductsJSON) > 0 {
			json.Unmarshal(targetProductsJSON, &h.TargetProducts)
		}
		if createdBy.Valid {
			h.CreatedBy = &createdBy.String
		}
		if teamID.Valid {
			h.TeamID = &teamID.String
		}
		if hasImage.Valid {
			h.HasImage = hasImage.Bool
		}
		if excerpt.Valid {
			h.Excerpt = excerpt.String
		}

		histories = append(histories, h)
	}

	return histories, nil
}

func (r *HistoryRepository) GetAllPublished(
	ctx context.Context,
	filter history.HistoryFilter,
) (*paginate.PaginatedResult[history.History], error) {

	// Set default pagination
	filter = r.normalizeFilter(filter)

	// Build where clause
	whereClause, args := r.buildWhereClause(filter)

	// Build query with pagination - PASS argCount!
	query := r.buildPublishedQuery(whereClause, len(args))

	// Combine WHERE args with pagination args
	allArgs := append(args, filter.Limit, filter.Offset)

	rows, err := r.db.QueryContext(ctx, query, allArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to query published histories: %w", err)
	}
	defer rows.Close()

	histories, err := r.scanHistoryPublished(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan history: %w", err)
	}

	// Get total count - use original WHERE args (without pagination)
	totalItems, err := r.countPublished(ctx, whereClause, args)
	if err != nil {
		return nil, err
	}

	return r.buildPaginatedResult(histories, totalItems, filter), nil
}

// Helper functions
func (r *HistoryRepository) normalizeFilter(filter history.HistoryFilter) history.HistoryFilter {
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	return filter
}

func (r *HistoryRepository) buildWhereClause(filter history.HistoryFilter) (string, []interface{}) {
	whereClause := "status = 'published'"
	var args []interface{}
	argIndex := 1

	if filter.ProductID != "" {
		whereClause += fmt.Sprintf(` AND target_products @> $%d::jsonb`, argIndex)
		jsonValue, _ := json.Marshal([]string{filter.ProductID})
		args = append(args, string(jsonValue))
		argIndex++
	}

	if filter.Search != "" {
		whereClause += fmt.Sprintf(` AND (title ILIKE $%d OR topic ILIKE $%d)`, argIndex, argIndex+1)
		searchTerm := "%" + filter.Search + "%"
		args = append(args, searchTerm, searchTerm)
		argIndex += 2
	}

	if filter.TeamID != "" {
		whereClause += fmt.Sprintf(` AND team_id = $%d`, argIndex)
		args = append(args, filter.TeamID)
		argIndex++
	}

	return whereClause, args
}

func (r *HistoryRepository) buildPublishedQuery(whereClause string, argCount int) string {
	return fmt.Sprintf(`
		SELECT 
			id, title, slug, topic, article, excerpt,
			image_url, image_prompt, target_products, 
			status, published_at, scheduled_for,
			created_by, team_id, user_id,
			has_image, seo_score, created_at
		FROM drafts
		WHERE %s
		ORDER BY 
			CASE 
				WHEN published_at IS NOT NULL THEN published_at 
				ELSE created_at 
			END DESC,
			created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argCount+1, argCount+2)
}

func (r *HistoryRepository) countPublished(ctx context.Context, whereClause string, args []interface{}) (int, error) {
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM drafts WHERE %s`, whereClause)
	var totalItems int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalItems)
	return totalItems, err
}

func (r *HistoryRepository) buildPaginatedResult(
	data []history.History,
	totalItems int,
	filter history.HistoryFilter,
) *paginate.PaginatedResult[history.History] {
	totalPages := 0
	if filter.Limit > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(filter.Limit)))
	}

	currentPage := 1
	if filter.Limit > 0 {
		currentPage = (filter.Offset / filter.Limit) + 1
	}

	return &paginate.PaginatedResult[history.History]{
		Data:        data,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		Limit:       filter.Limit,
		Offset:      filter.Offset,
	}
}

func (r *HistoryRepository) scanHistoryPublished(rows *sql.Rows) ([]history.History, error) {
	var histories []history.History

	for rows.Next() {
		var h history.History
		var imageURL, imagePrompt sql.NullString
		var publishedAt, scheduledFor sql.NullTime
		var targetProductsJSON []byte
		var excerpt sql.NullString
		var seoScore sql.NullInt32
		var hasImage sql.NullBool
		var userID sql.NullString
		var teamID sql.NullString

		err := rows.Scan(
			&h.ID,
			&h.Title,
			&h.Slug,
			&h.Topic,
			&h.Content,
			&excerpt,
			&imageURL,
			&imagePrompt,
			&targetProductsJSON,
			&h.Status,
			&publishedAt,
			&scheduledFor,
			&h.CreatedBy,
			&teamID,
			&userID,
			&hasImage,
			&seoScore,
			&h.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if imageURL.Valid {
			h.ImageURL = &imageURL.String
		}

		if publishedAt.Valid {
			h.PublishedAt = &publishedAt.Time
		}

		if hasImage.Valid {
			h.HasImage = hasImage.Bool
		}
		if excerpt.Valid {
			h.Excerpt = excerpt.String
		}
		if seoScore.Valid {
			h.SeoScore = int(seoScore.Int32)
		}

		if teamID.Valid {
			h.TeamID = &teamID.String
		}

		if len(targetProductsJSON) > 0 {
			var targetProducts []string
			if err := json.Unmarshal(targetProductsJSON, &targetProducts); err != nil {
				log.Printf("Warning: failed to parse target_products JSON: %v", err)
			} else {
				h.TargetProducts = targetProducts
			}
		}

		histories = append(histories, h)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return histories, nil
}
