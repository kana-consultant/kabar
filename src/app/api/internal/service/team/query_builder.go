package team

import (
	"fmt"
	"seo-backend/internal/models"
	"strings"
)

type QueryBuilder struct{}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

func (qb *QueryBuilder) BuildListQuery(ctx models.UserContext, filters TeamFilters) (string, []interface{}) {
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1

	// Role-based authorization
	if !ctx.IsAdmin() {
		if ctx.GetTeamID() != "" && !isZeroUUID(ctx.GetTeamID()) {
			conditions = append(conditions, fmt.Sprintf("id = $%d", argIndex))
			args = append(args, ctx.GetTeamID())
			argIndex++
		} else {
			conditions = append(conditions, "1=0")
		}
	}

	if filters.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, filters.Status)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT id, name, description,
			created_by, created_at, updated_at
		FROM teams %s
		ORDER BY created_at DESC
	`, whereClause)

	return query, args
}

func isZeroUUID(id string) bool {
	return id == "" || id == "00000000-0000-0000-0000-000000000000"
}
