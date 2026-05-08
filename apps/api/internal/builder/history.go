// internal/builder/query.go
package builder

import (
	"encoding/json"
	"log"
	"seo-backend/internal/models"
	"time"
)

func BuildHistoryInsert(req models.DraftDataPost, userID string, teamID string) (
	string,
	[]interface{},
	error,
) {

	log.Printf(
		"[HISTORY INSERT] req=%+v userID=%s teamID=%s",
		req,
		userID,
		teamID,
	)

	createdAt := time.Now()
	publishedAt := time.Now()

	qb := NewQueryBuilder("histories")

	targetProductsJSON, _ := json.Marshal(req.TargetProducts)

	sqlQuery, args, err := qb.
		Insert().
		Columns(
			"title",
			"topic",
			"content",
			"image_url",
			"target_products",
			"status",
			"action",
			"error_message",
			"published_at",
			"scheduled_for",
			"created_by",
			"team_id",
			"created_at",
		).
		Values(
			req.Title,
			req.Topic,
			req.Article,
			req.ImageURL,
			targetProductsJSON,
			"published",
			"publish",
			nil,
			publishedAt,
			nil,
			userID,
			teamID,
			createdAt,
		).
		Build()

	if err != nil {
		log.Printf(
			"[HISTORY INSERT BUILD ERROR] %v",
			err,
		)
		return "", nil, err
	}

	log.Printf("[HISTORY SQL] %s", sqlQuery)
	log.Printf("[HISTORY ARGS] %#v", args)

	return sqlQuery, args, nil
}
