package user

import (
	"fmt"

	"seo-backend/internal/models"
)

type QueryBuilder struct{}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

// BuildListQuery builds query for listing users with filters
// BuildListQuery builds query for listing users with filters
func (qb *QueryBuilder) BuildListQuery(ctx models.UserContext) (string, []interface{}) {
	query := `
		SELECT DISTINCT u.id, u.email, u.name, u.role, u.avatar, u.status, u.last_active, u.created_at, u.updated_at
		FROM users u
		INNER JOIN team_members tm ON u.id = tm.user_id 
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	userID := ctx.GetUserID()
	teamID := ctx.GetTeamID()

	// Filter berdasarkan team_id jika tersedia
	if teamID != "" {
		query += fmt.Sprintf(" AND tm.team_id = $%d", argIndex)
		args = append(args, teamID)
		argIndex++
	} else {
		// Jika tidak ada team_id spesifik, ambil semua team yang dimiliki user
		query += fmt.Sprintf(" AND tm.team_id IN (SELECT team_id FROM team_members WHERE user_id = $%d)", argIndex)
		args = append(args, userID)
		argIndex++
	}

	// Non-admin hanya bisa melihat active users
	query += fmt.Sprintf(" AND u.status = $%d", argIndex)
	args = append(args, "active")
	argIndex++

	orderBy := "u.created_at DESC"

	query += fmt.Sprintf(" ORDER BY %s", orderBy)

	return query, args
}
