package history

import (
	"fmt"
	"seo-backend/internal/models"
	"strings"
)

type QueryBuilder struct{}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

func (qb *QueryBuilder) BuildListQuery(ctx models.UserContext, filter HistoryFilter) (string, []interface{}) {
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1

	// Role-based conditions
	conditions = append(conditions, qb.buildAccessCondition(ctx, &argIndex, &args)...)

	// Additional filters
	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, filter.Status)
		argIndex++
	}

	if filter.Action != "" {
		conditions = append(conditions, fmt.Sprintf("action = $%d", argIndex))
		args = append(args, filter.Action)
		argIndex++
	}

	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		conditions = append(conditions, fmt.Sprintf("(title ILIKE $%d OR topic ILIKE $%d)", argIndex, argIndex+1))
		args = append(args, searchPattern, searchPattern)
		argIndex += 2
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT id, title, topic, content, image_url, target_products,
			status, action, error_message, published_at, scheduled_for,
			created_by, team_id, created_at
		FROM histories %s
		ORDER BY published_at DESC
	`, whereClause)

	// Add pagination
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}
	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
	}

	return query, args
}

func (qb *QueryBuilder) BuildStatsQuery(ctx models.UserContext) (string, []interface{}) {
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1

	// Role-based conditions
	conditions = append(conditions, qb.buildAccessCondition(ctx, &argIndex, &args)...)

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'success' THEN 1 END) as success_count,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_count,
			COUNT(CASE WHEN action = 'published' THEN 1 END) as published_count,
			COUNT(CASE WHEN action = 'scheduled' THEN 1 END) as scheduled_count
		FROM histories %s
	`, whereClause)

	return query, args
}

func (qb *QueryBuilder) buildAccessCondition(ctx models.UserContext, argIndex *int, args *[]interface{}) []string {
	conditions := make([]string, 0)

	switch ctx.GetRole() {
	case "super_admin":
		// No filter - access all
	case "admin":
		if ctx.GetTeamID() != "" && !isZeroUUID(ctx.GetTeamID()) {
			conditions = append(conditions, fmt.Sprintf("team_id = $%d", *argIndex))
			*args = append(*args, ctx.GetTeamID())
			*argIndex++
		} else {
			conditions = append(conditions, fmt.Sprintf("created_by = $%d", *argIndex))
			*args = append(*args, ctx.GetUserID())
			*argIndex++
		}
	default:
		conditions = append(conditions, fmt.Sprintf("created_by = $%d", *argIndex))
		*args = append(*args, ctx.GetUserID())
		*argIndex++
	}

	return conditions
}

func isZeroUUID(id string) bool {
	return id == "00000000-0000-0000-0000-000000000000"
}
