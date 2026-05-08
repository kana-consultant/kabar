package user

import (
	"fmt"
	"seo-backend/internal/models"
	"strings"
)

type QueryBuilder struct{}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

func (qb *QueryBuilder) BuildListQuery(ctx models.UserContext, filters UserFilters) (string, []interface{}) {
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1

	// Role-based authorization
	conditions = append(conditions, qb.buildAccessCondition(ctx, &argIndex, &args)...)

	// Additional filters
	if filters.Role != "" {
		conditions = append(conditions, fmt.Sprintf("role = $%d", argIndex))
		args = append(args, filters.Role)
		argIndex++
	}
	if filters.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, filters.Status)
		argIndex++
	}
	if filters.Search != "" {
		searchPattern := "%" + filters.Search + "%"
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR email ILIKE $%d)", argIndex, argIndex+1))
		args = append(args, searchPattern, searchPattern)
		argIndex += 2
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT id, email, name, role, avatar, status, last_active, created_at, updated_at
		FROM users %s
		ORDER BY created_at DESC
	`, whereClause)

	return query, args
}

func (qb *QueryBuilder) buildAccessCondition(ctx models.UserContext, argIndex *int, args *[]interface{}) []string {
	conditions := make([]string, 0)

	switch ctx.GetRole() {
	case "super_admin":
		// No filter for super admin
	case "admin":
		if ctx.GetTeamID() != "" && !isZeroUUID(ctx.GetTeamID()) {
			conditions = append(conditions, fmt.Sprintf(`
				id IN (
					SELECT user_id FROM team_members WHERE team_id = $%d
				)
			`, *argIndex))
			*args = append(*args, ctx.GetTeamID())
			*argIndex++
		} else {
			conditions = append(conditions, fmt.Sprintf("id = $%d", *argIndex))
			*args = append(*args, ctx.GetUserID())
			*argIndex++
		}
	default:
		conditions = append(conditions, fmt.Sprintf("id = $%d", *argIndex))
		*args = append(*args, ctx.GetUserID())
		*argIndex++
	}

	return conditions
}

func isZeroUUID(id string) bool {
	return id == "" || id == "00000000-0000-0000-0000-000000000000"
}
