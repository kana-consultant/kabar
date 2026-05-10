package draft

import (
	"encoding/json"
	"time"

	"seo-backend/internal/helper"
	"seo-backend/internal/models"
)

func (s *Service) InsertHistory(req models.DraftDataPost, userID, teamID, action string) error {
	targetProductsJSON, _ := json.Marshal(req.TargetProducts)

	query := `INSERT INTO histories (
		title, topic, content, image_url, target_products,
		status, action, published_at, created_by, team_id, created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := s.db.Exec(query,
		req.Title, req.Topic, req.Article, req.ImageURL,
		targetProductsJSON, "published", action, helper.ParseWIBTime(time.Now().Format(time.RFC3339)),
		userID, teamID, helper.ParseWIBTime(time.Now().Format(time.RFC3339)),
	)
	return err
}
