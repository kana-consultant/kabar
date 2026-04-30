package draft

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"seo-backend/internal/builder"
	"seo-backend/internal/database"
	"seo-backend/internal/models"
)

type repositoryImpl struct {
	db *sql.DB
}

func NewRepository() Repository {
	return &repositoryImpl{
		db: database.GetDB(),
	}
}

func (r *repositoryImpl) GetAll(ctx context.Context, filters GetAllFilters) ([]models.Draft, error) {
	qb := builder.NewQueryBuilder("drafts")
	query := qb.Select(
		"id", "title", "topic", "article", "image_url", "image_prompt",
		"status", "scheduled_for", "target_products", "has_image",
		"team_id", "user_id", "created_at", "updated_at",
	).OrderBy("created_at DESC")

	// Apply filters
	if filters.TeamID != nil && *filters.TeamID != "" {
		query = query.WhereEq("team_id", *filters.TeamID)
	}
	if filters.UserID != nil && *filters.UserID != "" && filters.Role == "user" {
		query = query.WhereEq("user_id", *filters.UserID)
	}
	if filters.Status != "" {
		query = query.WhereEq("status", filters.Status)
	}
	if filters.Search != "" {
		query = query.WhereLike("title", filters.Search)
	}

	sqlQuery, args, err := query.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var drafts []models.Draft
	for rows.Next() {
		var d models.Draft
		var targetProductsJSON []byte
		var scheduledFor sql.NullTime

		err := rows.Scan(
			&d.ID, &d.Title, &d.Topic, &d.Article, &d.ImageURL, &d.ImagePrompt,
			&d.Status, &scheduledFor, &targetProductsJSON, &d.HasImage,
			&d.TeamID, &d.UserID, &d.CreatedAt, &d.UpdatedAt,
		)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}

		if scheduledFor.Valid {
			d.ScheduledFor = &scheduledFor.Time
		}

		json.Unmarshal(targetProductsJSON, &d.TargetProducts)
		drafts = append(drafts, d)
	}

	return drafts, nil
}

func (r *repositoryImpl) GetByID(ctx context.Context, id string) (*models.Draft, error) {
	qb := builder.NewQueryBuilder("drafts")
	sqlQuery, args, err := qb.Select(
		"id", "title", "topic", "article", "image_url", "image_prompt",
		"status", "scheduled_for", "target_products", "has_image",
		"team_id", "user_id", "created_at", "updated_at",
	).WhereEq("id", id).Build()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var draft models.Draft
	var targetProductsJSON []byte
	var scheduledFor sql.NullTime

	err = r.db.QueryRowContext(ctx, sqlQuery, args...).Scan(
		&draft.ID, &draft.Title, &draft.Topic, &draft.Article,
		&draft.ImageURL, &draft.ImagePrompt, &draft.Status,
		&scheduledFor, &targetProductsJSON, &draft.HasImage,
		&draft.TeamID, &draft.UserID, &draft.CreatedAt, &draft.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if scheduledFor.Valid {
		draft.ScheduledFor = &scheduledFor.Time
	}

	json.Unmarshal(targetProductsJSON, &draft.TargetProducts)
	return &draft, nil
}

func (r *repositoryImpl) Create(ctx context.Context, draft *models.Draft) (string, error) {
	targetProductsJSON, _ := json.Marshal(draft.TargetProducts)

	qb := builder.NewQueryBuilder("drafts")
	sqlQuery, args, err := qb.Insert().
		Columns(
			"title", "topic", "article", "image_url", "image_prompt",
			"status", "target_products", "has_image", "created_by", "team_id", "user_id",
		).
		Values(
			draft.Title, draft.Topic, draft.Article, draft.ImageURL, draft.ImagePrompt,
			draft.Status, targetProductsJSON, draft.HasImage, draft.UserID, draft.TeamID, draft.UserID,
		).
		Returning("id").
		Build()

	if err != nil {
		return "", fmt.Errorf("failed to build insert query: %w", err)
	}

	var id string
	err = r.db.QueryRowContext(ctx, sqlQuery, args...).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to create draft: %w", err)
	}

	return id, nil
}

func (r *repositoryImpl) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	updates["updated_at"] = time.Now()

	qb := builder.NewQueryBuilder("drafts")
	sqlQuery, args, err := qb.Update().
		SetMap(updates).
		WhereEq("id", id).
		Build()

	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to update draft: %w", err)
	}

	return nil
}

func (r *repositoryImpl) Delete(ctx context.Context, id string) error {
	qb := builder.NewQueryBuilder("drafts")
	sqlQuery, args, err := qb.Delete().
		WhereEq("id", id).
		Build()

	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to delete draft: %w", err)
	}

	return nil
}

func (r *repositoryImpl) GetDraftForPublish(ctx context.Context, id string) (*PublishData, error) {
	var data PublishData
	var targetProductsJSON []byte

	err := r.db.QueryRowContext(ctx, `
		SELECT id, title, topic, article, image_url, target_products, COALESCE(image_prompt, ''), team_id
		FROM drafts WHERE id = $1
	`, id).Scan(&data.ID, &data.Title, &data.Topic, &data.Article, &data.ImageURL, &targetProductsJSON, &data.ImagePrompt, &data.TeamID)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(targetProductsJSON, &data.TargetProducts)
	return &data, nil
}

func (r *repositoryImpl) UpdatePublishStatus(ctx context.Context, id string, status string, scheduledFor *time.Time, pattern string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	if scheduledFor != nil {
		updates["scheduled_for"] = scheduledFor
	}
	if pattern != "" {
		updates["scheduled_pattern"] = pattern
	}
	if status == "published" && scheduledFor == nil {
		updates["published_at"] = time.Now()
	}

	return r.Update(ctx, id, updates)
}

func (r *repositoryImpl) UpdateProductSyncStatus(ctx context.Context, productID string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE products
		SET sync_status = 'synced', last_sync = NOW()
		WHERE id = $1
	`, productID)
	return err
}

func (r *repositoryImpl) SaveToHistory(ctx context.Context, history *PublishHistory) error {
	productsJSON, _ := json.Marshal(history.Products)

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO histories (title, topic, content, image_url, target_products, 
		                       status, action, error_message, published_at, created_by, team_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, history.Title, history.Topic, history.Content, history.ImageURL, productsJSON,
		history.Status, history.Action, history.Error, history.PublishedAt, history.CreatedBy, history.TeamID)

	return err
}
