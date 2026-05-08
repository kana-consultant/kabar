// internal/application/team/query_builder.go
package team

import (
	"fmt"
	"strings"

	"seo-backend/internal/domain/team"
	"seo-backend/internal/models"
)

type QueryBuilder struct{}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

func (qb *QueryBuilder) BuildListQuery(ctx models.UserContext, filters team.TeamFilters) (string, []interface{}) {
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1

	if !ctx.IsAdmin() {
		conditions = append(conditions, fmt.Sprintf("id IN (SELECT team_id FROM team_members WHERE user_id = $%d)", argIndex))
		args = append(args, ctx.GetUserID())
		argIndex++
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
