package draft

import (
	"database/sql"
	"encoding/json"
	"time"

	"seo-backend/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(id string) (*models.Draft, error) {
	var draft models.Draft
	var targetProductsJSON []byte

	query := `SELECT id, title, topic, article, image_url, image_prompt,
		status, scheduled_for, target_products, has_image,
		team_id, user_id, created_at, updated_at
		FROM drafts WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&draft.ID, &draft.Title, &draft.Topic, &draft.Article,
		&draft.ImageURL, &draft.ImagePrompt, &draft.Status,
		&draft.ScheduledFor, &targetProductsJSON, &draft.HasImage,
		&draft.TeamID, &draft.UserID, &draft.CreatedAt, &draft.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(targetProductsJSON, &draft.TargetProducts)
	return &draft, nil
}

func (r *Repository) Create(req models.CreateDraftRequest, userID, teamID string) (string, error) {
	targetProductsJSON, _ := json.Marshal(req.TargetProducts)

	query := `INSERT INTO drafts (
		title, topic, article, image_url, image_prompt,
		status, target_products, has_image, created_by, team_id, user_id
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING id`

	var draftID string
	err := r.db.QueryRow(
		query,
		req.Title, req.Topic, req.Article, req.ImageURL, req.ImagePrompt,
		"draft", targetProductsJSON, req.HasImage, nullIfEmpty(userID), nullIfEmpty(teamID), nullIfEmpty(userID),
	).Scan(&draftID)

	return draftID, err
}

func (r *Repository) Update(id string, data map[string]interface{}) error {
	data["updated_at"] = time.Now()

	query, args, err := buildUpdateQuery(id, data)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(query, args...)
	return err
}

func (r *Repository) Delete(id string) error {
	query := `DELETE FROM drafts WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *Repository) UpdateStatus(id, status string, scheduledFor *time.Time) error {
	if scheduledFor != nil {
		_, err := r.db.Exec(`
			UPDATE drafts 
			SET status = $1, scheduled_for = $2 
			WHERE id = $3
		`, status, *scheduledFor, id)
		return err
	}

	_, err := r.db.Exec(`
		UPDATE drafts 
		SET status = $1, updated_at = $2
		WHERE id = $3 AND status = 'scheduled'
	`, status, time.Now(), id)
	return err
}

func (r *Repository) GetDraftData(id string) (*DraftData, error) {
	var draft DraftData
	var targetProductsJSON []byte

	query := `SELECT title, topic, article, image_url, target_products, COALESCE(image_prompt, '')
		FROM drafts WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&draft.Title, &draft.Topic, &draft.Article,
		&draft.ImageURL, &targetProductsJSON, &draft.ImagePrompt,
	)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(targetProductsJSON, &draft.TargetProducts)
	draft.ID = id
	return &draft, nil
}

func (r *Repository) InsertScheduledDraft(req models.ScheduleRequest, scheduledFor time.Time, teamID, userID string) (string, error) {
	targetProductsJSON, _ := json.Marshal(req.TargetProducts)
	now := time.Now()

	var draftID string
	err := r.db.QueryRow(`
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

func (r *Repository) InsertHistory(req models.DraftDataPost, userID, teamID, action string) error {
	targetProductsJSON, _ := json.Marshal(req.TargetProducts)

	query := `INSERT INTO histories (
		title, topic, content, image_url, target_products,
		status, action, published_at, created_by, team_id, created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.db.Exec(query,
		req.Title, req.Topic, req.Article, req.ImageURL,
		targetProductsJSON, "published", action, time.Now(),
		userID, teamID, time.Now(),
	)
	return err
}
