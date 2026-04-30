package draft

import (
	"database/sql"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"seo-backend/internal/database"
	"seo-backend/internal/models"
)

type DraftQueryBuilder struct {
	teamID   string
	userRole string
	userID   string
	status   string
	search   string
}

func NewDraftQueryBuilder() *DraftQueryBuilder {
	return &DraftQueryBuilder{}
}

func (q *DraftQueryBuilder) SetUserContext(role, teamID, userID string) *DraftQueryBuilder {
	q.userRole = role
	q.teamID = teamID
	q.userID = userID
	return q
}

func (q *DraftQueryBuilder) WithStatus(status string) *DraftQueryBuilder {
	q.status = status
	return q
}

func (q *DraftQueryBuilder) WithSearch(search string) *DraftQueryBuilder {
	q.search = search
	return q
}

func (q *DraftQueryBuilder) Execute() ([]models.Draft, error) {
	query, args := q.buildQuery()
	rows, err := database.GetDB().Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return q.scanRows(rows)
}

func (q *DraftQueryBuilder) buildQuery() (string, []interface{}) {
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1

	// Role-based conditions
	switch q.userRole {
	case "super_admin":
		// No filter for super admin
	case "admin":
		if q.teamID != "" && q.teamID != "00000000-0000-0000-0000-000000000000" {
			conditions = append(conditions, "team_id = $"+strconv.Itoa(argIndex))
			args = append(args, q.teamID)
			argIndex++
		} else {
			conditions = append(conditions, "user_id = $"+strconv.Itoa(argIndex))
			args = append(args, q.userID)
			argIndex++
		}
	default:
		conditions = append(conditions, "user_id = $"+strconv.Itoa(argIndex))
		args = append(args, q.userID)
		argIndex++
	}

	// Status filter
	if q.status != "" {
		conditions = append(conditions, "status = $"+strconv.Itoa(argIndex))
		args = append(args, q.status)
		argIndex++
	}

	// Search filter
	if q.search != "" {
		conditions = append(conditions, "title ILIKE $"+strconv.Itoa(argIndex))
		args = append(args, "%"+q.search+"%")
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := `SELECT id, title, topic, article, image_url, image_prompt,
		status, scheduled_for, target_products, has_image,
		team_id, user_id, created_at, updated_at
		FROM drafts ` + whereClause + ` ORDER BY created_at DESC`

	return query, args
}

func (q *DraftQueryBuilder) scanRows(rows *sql.Rows) ([]models.Draft, error) {
	var drafts []models.Draft
	for rows.Next() {
		var d models.Draft
		var targetProductsJSON []byte

		err := rows.Scan(
			&d.ID, &d.Title, &d.Topic, &d.Article, &d.ImageURL, &d.ImagePrompt,
			&d.Status, &d.ScheduledFor, &targetProductsJSON, &d.HasImage,
			&d.TeamID, &d.UserID, &d.CreatedAt, &d.UpdatedAt,
		)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}

		json.Unmarshal(targetProductsJSON, &d.TargetProducts)
		drafts = append(drafts, d)
	}
	return drafts, nil
}
