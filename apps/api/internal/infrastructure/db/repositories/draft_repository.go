// internal/infrastructure/repository/draft/repository.go
package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"seo-backend/internal/domain/draft"
	"seo-backend/internal/helper"
)

type RepositoryImpl struct {
	db *sql.DB
}

func NewDraftRepository(db *sql.DB) draft.Repository {
	return &RepositoryImpl{db: db}
}

func (r *RepositoryImpl) GetByID(ctx context.Context, id string) (*draft.DraftData, error) {
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
		target_products 
	FROM drafts 
	WHERE id = $1
`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&d.ID, &d.Title, &d.Topic, &d.Article,
		&d.ImageURL, &d.ImagePrompt, &targetProductsJSON,
	)
	if err != nil {
		log.Printf("=============%v", err)
		return nil, err
	}

	json.Unmarshal(targetProductsJSON, &d.TargetProducts)
	return &d, nil
}

func (r *RepositoryImpl) GetAll(ctx context.Context, teamID string) (*[]draft.Draft, error) {
	var drafts []draft.Draft

	query := `
		SELECT 
			id,
			title, 
			topic, 
			article, 
			target_products, 
			team_id,
			status,
			COALESCE(image_prompt, '')
		FROM drafts
		WHERE team_id = $1 AND status != 'scheduled'
	`

	rows, err := r.db.QueryContext(ctx, query, teamID)
	if err != nil {
		log.Printf("===================,%v", err)
		return nil, err
	}
	defer rows.Close()

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
			&d.Status,
			&d.ImagePrompt,
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

	log.Printf("drafts: %+v", drafts)

	return &drafts, nil
}

func (r *RepositoryImpl) GetAllScheduled(ctx context.Context, teamID string) (*[]draft.Draft, error) {
	var drafts []draft.Draft

	query := `
		SELECT 
			id,
			title, 
			topic, 
			article, 
			target_products, 
			team_id,
			status,
			COALESCE(image_prompt, ''),
			scheduled_for
		FROM drafts
		WHERE team_id = $1 
		AND status = 'scheduled'
	`

	rows, err := r.db.QueryContext(ctx, query, teamID)
	if err != nil {
		log.Printf("query error: %v", err)
		return nil, err
	}
	defer rows.Close()

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

	return &drafts, nil
}

func (r *RepositoryImpl) Create(ctx context.Context, req draft.CreateDraftRequest, userID, teamID string) (string, error) {
	targetProductsJSON, _ := json.Marshal(req.TargetProducts)
	log.Printf("==============,%s", req.TargetProducts)

	query := `INSERT INTO drafts (
		title, topic, article, image_url, image_prompt,
		status, target_products, has_image, created_by, team_id, user_id
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING id`

	var draftID string
	err := r.db.QueryRowContext(
		ctx,
		query,
		req.Title, req.Topic, req.Article, req.ImageURL, req.ImagePrompt,
		"draft", targetProductsJSON, req.HasImage, nullIfEmpty(userID), nullIfEmpty(teamID), nullIfEmpty(userID),
	).Scan(&draftID)

	return draftID, err
}

func (r *RepositoryImpl) Update(ctx context.Context, id string, data map[string]interface{}) error {
	data["updated_at"] = helper.ParseWIBTime(time.Now().Format(time.RFC3339))

	query, args, err := r.buildUpdateQuery(id, data)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
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

func (r *RepositoryImpl) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM drafts WHERE id = $1", id)
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
	targetProductsJSON, _ := json.Marshal(req.TargetProducts)

	var status string
	status = "published"

	if action == "failed" {
		status = "failed"
	}

	query := `INSERT INTO histories (
		title, topic, content, image_url, target_products,
		status, action, published_at, created_by, team_id, created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.db.ExecContext(ctx, query,
		req.Title, req.Topic, req.Article, req.ImageURL,
		targetProductsJSON, status, action, helper.ParseWIBTime(time.Now().Format(time.RFC3339)),
		userID, teamID, helper.ParseWIBTime(time.Now().Format(time.RFC3339)),
	)
	return err
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
