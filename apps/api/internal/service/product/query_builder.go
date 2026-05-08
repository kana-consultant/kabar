package product

import (
	"fmt"
	"strings"
)

type QueryBuilder struct{}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

func (qb *QueryBuilder) BuildListQuery(ctx UserContext, filters ProductFilters) (string, []interface{}) {
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1

	// Role-based conditions
	conditions = append(conditions, qb.buildAccessCondition(ctx, &argIndex, &args)...)

	// Additional filters
	if filters.Platform != "" {
		conditions = append(conditions, fmt.Sprintf("platform = $%d", argIndex))
		args = append(args, filters.Platform)
		argIndex++
	}
	if filters.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, filters.Status)
		argIndex++
	}
	if filters.SyncStatus != "" {
		conditions = append(conditions, fmt.Sprintf("sync_status = $%d", argIndex))
		args = append(args, filters.SyncStatus)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT id, name, platform, api_endpoint, status, sync_status, last_sync,
			created_by, team_id, user_id, created_at, updated_at
		FROM products %s
		ORDER BY created_at DESC
	`, whereClause)

	return query, args
}

func (qb *QueryBuilder) buildAccessCondition(ctx UserContext, argIndex *int, args *[]interface{}) []string {
	conditions := make([]string, 0)

	switch ctx.GetRole() {
	case "super_admin":
		// No filter
	case "admin":
		if ctx.GetTeamID() != "" && !isZeroUUID(ctx.GetTeamID()) {
			conditions = append(conditions, fmt.Sprintf("team_id = $%d", *argIndex))
			*args = append(*args, ctx.GetTeamID())
			*argIndex++
		} else {
			conditions = append(conditions, fmt.Sprintf("user_id = $%d", *argIndex))
			*args = append(*args, ctx.GetUserID())
			*argIndex++
		}
	default:
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", *argIndex))
		*args = append(*args, ctx.GetUserID())
		*argIndex++
	}

	return conditions
}

func isZeroUUID(id string) bool {
	return id == "00000000-0000-0000-0000-000000000000"
}
