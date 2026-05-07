package draft

import (
	"encoding/json"
	"time"

	"seo-backend/internal/models"
)

func (s *Service) GetDraftByID(id string) (*models.Draft, error) {
	var draft models.Draft
	var targetProductsJSON []byte

	query := `SELECT id, title, topic, article, image_url, image_prompt,
		status, scheduled_for, target_products, has_image,
		team_id, user_id, created_at, updated_at
		FROM drafts WHERE id = $1`

	err := s.db.QueryRow(query, id).Scan(
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

func (s *Service) CreateDraftRecord(req models.CreateDraftRequest, userID, teamID string) (string, error) {
	targetProductsJSON, _ := json.Marshal(req.TargetProducts)

	createdBy := nullIfEmpty(userID)
	teamIDPtr := nullIfEmpty(teamID)

	query := `INSERT INTO drafts (
		title, topic, article, image_url, image_prompt,
		status, target_products, has_image, created_by, team_id, user_id
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING id`

	var draftID string
	err := s.db.QueryRow(
		query,
		req.Title, req.Topic, req.Article, req.ImageURL, req.ImagePrompt,
		"draft", targetProductsJSON, req.HasImage, createdBy, teamIDPtr, createdBy,
	).Scan(&draftID)

	return draftID, err
}

func (s *Service) UpdateDraftRecord(id string, data map[string]interface{}) error {
	data["updated_at"] = time.Now()

	query, args, err := buildUpdateQuery(id, data)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(query, args...)
	return err
}

func (s *Service) DeleteDraftRecord(id string) error {
	query := `DELETE FROM drafts WHERE id = $1`
	_, err := s.db.Exec(query, id)
	return err
}

func (s *Service) getDraftData(id string) (*DraftData, error) {
	var draft DraftData
	var targetProductsJSON []byte

	query := `SELECT title, topic, article, image_url, target_products, COALESCE(image_prompt, '')
		FROM drafts WHERE id = $1`

	err := s.db.QueryRow(query, id).Scan(
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

func (s *Service) updateDraftStatus(id string, status string, scheduledFor *time.Time) error {
	if scheduledFor != nil {
		_, err := s.db.Exec(`
			UPDATE drafts 
			SET status = $1, scheduled_for = $2 
			WHERE id = $3
		`, status, *scheduledFor, id)
		return err
	}

	_, err := s.db.Exec(`
		UPDATE drafts 
		SET status = $1, updated_at = $2
		WHERE id = $3 AND status = 'scheduled'
	`, status, time.Now(), id)
	return err
}

func (s *Service) insertScheduledDraftRecord(req models.ScheduleRequest, scheduledFor time.Time, teamID, userID string) (string, error) {
	targetProductsJSON, _ := json.Marshal(req.TargetProducts)
	now := time.Now()

	var draftID string
	err := s.db.QueryRow(`
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

func (s *Service) deleteDraftAfterPublish(id string) error {
	_, err := s.db.Exec("DELETE FROM drafts WHERE id = $1", id)
	return err
}
