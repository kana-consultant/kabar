package history

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"seo-backend/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(id string) (*models.History, error) {
	query := `
		SELECT id, title, topic, content, image_url, target_products,
			status, action, error_message, published_at, scheduled_for,
			created_by, team_id, created_at
		FROM histories WHERE id = $1
	`

	var history models.History
	var targetProductsJSON []byte

	err := r.db.QueryRow(query, id).Scan(
		&history.ID, &history.Title, &history.Topic, &history.Content, &history.ImageURL,
		&targetProductsJSON, &history.Status, &history.Action, &history.ErrorMessage,
		&history.PublishedAt, &history.ScheduledFor, &history.CreatedBy, &history.TeamID, &history.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch history: %w", err)
	}

	json.Unmarshal(targetProductsJSON, &history.TargetProducts)
	return &history, nil
}

func (r *Repository) GetAll(query string, args []interface{}) ([]models.History, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch history: %w", err)
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *Repository) Create(req CreateHistoryRequest) (string, error) {
	targetProductsJSON, _ := json.Marshal(req.TargetProducts)

	query := `
		INSERT INTO histories (
			title, topic, content, image_url, target_products,
			status, action, error_message, scheduled_for, published_at,
			created_by, team_id, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), $10, $11, NOW())
		RETURNING id
	`

	var historyID string
	err := r.db.QueryRow(
		query,
		req.Title, req.Topic, req.Content, req.ImageURL, targetProductsJSON,
		req.Status, req.Action, req.ErrorMessage, req.ScheduledFor,
		nullIfEmpty(req.CreatedBy), nullIfEmpty(req.TeamID),
	).Scan(&historyID)

	if err != nil {
		return "", fmt.Errorf("failed to create history: %w", err)
	}

	return historyID, nil
}

func (r *Repository) Delete(id string) error {
	query := "DELETE FROM histories WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete history: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("history not found")
	}

	return nil
}

func (r *Repository) ClearAll() error {
	query := "DELETE FROM histories"
	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to clear history: %w", err)
	}
	return nil
}

func (r *Repository) GetStats(query string, args []interface{}) (*HistoryStats, error) {
	var stats HistoryStats
	err := r.db.QueryRow(query, args...).Scan(
		&stats.Total, &stats.SuccessCount, &stats.FailedCount,
		&stats.PublishedCount, &stats.ScheduledCount,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stats: %w", err)
	}
	return &stats, nil
}

func (r *Repository) scanRows(rows *sql.Rows) ([]models.History, error) {
	var histories []models.History

	for rows.Next() {
		var h models.History
		var targetProductsJSON []byte

		err := rows.Scan(
			&h.ID, &h.Title, &h.Topic, &h.Content, &h.ImageURL,
			&targetProductsJSON, &h.Status, &h.Action, &h.ErrorMessage,
			&h.PublishedAt, &h.ScheduledFor, &h.CreatedBy, &h.TeamID, &h.CreatedAt,
		)
		if err != nil {
			continue
		}

		json.Unmarshal(targetProductsJSON, &h.TargetProducts)
		histories = append(histories, h)
	}

	return histories, nil
}

func nullIfEmpty(id string) interface{} {
	if id == "" || id == "00000000-0000-0000-0000-000000000000" {
		return nil
	}
	return id
}
