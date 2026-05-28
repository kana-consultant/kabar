package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"seo-backend/internal/domain/history"
	"seo-backend/internal/domain/paginate"
	"seo-backend/internal/helper"
	userRole "seo-backend/internal/helper/filter"
	"seo-backend/internal/helper/keywords"
	"seo-backend/internal/models"
	"strings"
	"time"
)

type HistoryRepository struct {
	db *sql.DB
}

func NewHistoryRepository(db *sql.DB) history.HistoryRepository {
	return &HistoryRepository{db: db}
}

// Create inserts a new history record
func (r *HistoryRepository) Create(ctx context.Context, data *history.History) (string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	targetProductsJSON, _ := json.Marshal(data.TargetProducts)

	query := `
		INSERT INTO history (
			id, title, topic, content, image_url, target_products,
			status, action, error_message, published_at, scheduled_for,
			created_by, team_id, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW())
		RETURNING id
	`

	var id string
	err = r.db.QueryRowContext(ctx, query,
		data.ID, data.Title, data.Topic, data.Content, data.ImageURL, targetProductsJSON,
		data.Status, data.Action, data.ErrorMessage, data.PublishedAt, data.ScheduledFor,
		data.CreatedBy, data.TeamID,
	).Scan(&id)

	if err := keywords.UpdateKeywords(ctx, tx, keywords.HistorySource{HistoryID: id}, data.Keywords); err != nil {
		return "", fmt.Errorf("failed to update keywords: %w", err)
	}

	if err != nil {
		return "", fmt.Errorf("failed to create history: %w", err)
	}

	return id, nil
}

// GetByID retrieves history by ID
func (r *HistoryRepository) GetByID(ctx context.Context, id string) (*history.History, error) {
	query := `
		SELECT id, title, topic, content, image_url, target_products,
			status, action, error_message, published_at, scheduled_for,
			created_by, team_id, created_at
		FROM histories WHERE id = $1
	`

	var h history.History
	var targetProductsJSON []byte
	var errorMessage sql.NullString
	var createdBy sql.NullString
	var teamID sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&h.ID, &h.Title, &h.Topic, &h.Content, &h.ImageURL, &targetProductsJSON,
		&h.Status, &h.Action, &errorMessage, &h.PublishedAt, &h.ScheduledFor,
		&createdBy, &teamID, &h.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get history: %w", err)
	}

	// Unmarshal JSON
	if len(targetProductsJSON) > 0 {
		json.Unmarshal(targetProductsJSON, &h.TargetProducts)
	}
	if errorMessage.Valid {
		h.ErrorMessage = &errorMessage.String
	}
	if createdBy.Valid {
		h.CreatedBy = &createdBy.String
	}
	if teamID.Valid {
		h.TeamID = &teamID.String
	}

	keywords, err := keywords.GetKeywords(ctx, r.db, keywords.HistorySource{HistoryID: id})

	if err != nil {

		log.Printf(
			"Error getting keywords for draft %s: %v",
			id,
			err,
		)

		// fallback empty array
		h.Keywords = []string{}

	} else {

		h.Keywords = keywords
	}

	return &h, nil
}

func (r *HistoryRepository) GetAll(ctx context.Context, userCtx models.UserContext, params history.HistoryFilter) (*paginate.PaginatedResult[history.History], error) {
	// Build access filter
	whereClause, whereArgs := userRole.BuildAccessFilter(userCtx)

	// Count total + per status dalam satu query
	var totalItems, totalSuccess, totalFailed int
	countQuery := fmt.Sprintf(`
		SELECT
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE status = 'published') AS total_success,
			COUNT(*) FILTER (WHERE status = 'failed')  AS total_failed
		FROM histories
		WHERE %s
	`, whereClause)
	if err := r.db.QueryRowContext(ctx, countQuery, whereArgs...).
		Scan(&totalItems, &totalSuccess, &totalFailed); err != nil {
		return nil, fmt.Errorf("failed to count histories: %w", err)
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(params.Limit)))
	currentPage := (params.Offset / params.Limit) + 1

	// Append LIMIT & OFFSET
	args := append(whereArgs, params.Limit, params.Offset)
	query := fmt.Sprintf(`
		SELECT id, title, topic, content, image_url, target_products,
			status, action, error_message, published_at, scheduled_for,
			created_by, team_id, created_at
		FROM histories
		WHERE %s
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
		SELECT id, title, topic, content, image_url, target_products,
			status, action, error_message, published_at, scheduled_for,
			created_by, team_id, created_at
		FROM histories WHERE team_id = $1
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
		SELECT id, title, topic, content, image_url, target_products,
			status, action, error_message, published_at, scheduled_for,
			created_by, team_id, created_at
		FROM histories WHERE created_by = $1
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
		SELECT id, title, topic, content, image_url, target_products,
			status, action, error_message, published_at, scheduled_for,
			created_by, team_id, created_at
		FROM histories WHERE status = $1
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
		SELECT id, title, topic, content, image_url, target_products,
			status, action, error_message, published_at, scheduled_for,
			created_by, team_id, created_at
		FROM histories 
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
	fmt.Println("============== team_id", query.TeamID)
	countQuery := `SELECT COUNT(*) FROM histories WHERE 1=1`
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
		countQuery += fmt.Sprintf(" AND (title ILIKE $%d OR content ILIKE $%d)", argIndex, argIndex)
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
		FROM histories 
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

	for key, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, helper.ParseWIBTime(time.Now().Format(time.RFC3339)))
	argIndex++

	args = append(args, id)
	query := fmt.Sprintf("UPDATE history SET %s WHERE id = $%d",
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
	query := `DELETE FROM histories WHERE id = $1`

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
	query := `DELETE FROM histories WHERE team_id = $1`

	_, err := r.db.ExecContext(ctx, query, teamID)
	if err != nil {
		return fmt.Errorf("failed to delete history by team: %w", err)
	}

	return nil
}

// DeleteByStatus deletes history by status
func (r *HistoryRepository) DeleteByStatus(ctx context.Context, status string) error {
	query := `DELETE FROM histories WHERE status = $1`

	_, err := r.db.ExecContext(ctx, query, status)
	if err != nil {
		return fmt.Errorf("failed to delete history by status: %w", err)
	}

	return nil
}

// GetStats returns statistics
func (r *HistoryRepository) GetStats(ctx context.Context, query *history.HistoryFilter) (*history.HistoryStats, error) {
	stats := &history.HistoryStats{}

	// Total
	total, err := r.Count(ctx, *query)
	if err != nil {
		return nil, err
	}
	stats.Total = total

	// Count by status
	statusCount, err := r.GetCountByStatus(ctx, query.TeamID)
	if err != nil {
		return nil, err
	}

	// Mapping ke struct
	stats.SuccessCount = statusCount["success"]
	stats.FailedCount = statusCount["failed"]
	stats.PublishedCount = statusCount["published"]
	stats.ScheduledCount = statusCount["scheduled"]

	// Hitung success rate
	if stats.Total > 0 {
		stats.SuccessRate = float64(stats.SuccessCount) / float64(stats.Total) * 100
	}

	return stats, nil
}

// Helper functions
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
		var errorMessage sql.NullString
		var createdBy sql.NullString
		var teamID sql.NullString

		err := rows.Scan(
			&h.ID, &h.Title, &h.Topic, &h.Content, &h.ImageURL, &targetProductsJSON,
			&h.Status, &h.Action, &errorMessage, &h.PublishedAt, &h.ScheduledFor,
			&createdBy, &teamID, &h.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history: %w", err)
		}

		if len(targetProductsJSON) > 0 {
			json.Unmarshal(targetProductsJSON, &h.TargetProducts)
		}
		if errorMessage.Valid {
			h.ErrorMessage = &errorMessage.String
		}
		if createdBy.Valid {
			h.CreatedBy = &createdBy.String
		}
		if teamID.Valid {
			h.TeamID = &teamID.String
		}

		histories = append(histories, h)
	}

	return histories, nil
}
