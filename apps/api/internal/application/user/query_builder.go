package user

import (
	"fmt"
	"strings"

	"seo-backend/internal/domain/user"
	"seo-backend/internal/models"
)

type QueryBuilder struct{}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

// BuildListQuery builds query for listing users with filters
func (qb *QueryBuilder) BuildListQuery(ctx models.UserContext, filters user.UserFilters) (string, []interface{}) {
	query := `
		SELECT DISTINCT u.id, u.email, u.name, u.role, u.avatar, u.status, u.last_active, u.created_at, u.updated_at
		FROM users u
		INNER JOIN team_members tm ON u.id = tm.user_id
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	// Role-based team filtering
	userRole := ctx.GetRole()
	userID := ctx.GetUserID()

	switch userRole {
	case "admin", "super_admin":
		// Admin bisa melihat user dari team tertentu (jika ada filter TeamID)
		if filters.TeamID != "" {
			query += fmt.Sprintf(" AND tm.team_id = $%d", argIndex)
			args = append(args, filters.TeamID)
			argIndex++
		}
		// Jika tidak ada filter TeamID, admin bisa melihat semua user dari semua team
	default:
		// Non-admin: hanya bisa melihat user dari team mereka sendiri
		// Ambil semua team_id yang dimiliki user (pakai IN)
		query += fmt.Sprintf(" AND tm.team_id IN (SELECT team_id FROM team_members WHERE user_id = $%d)", argIndex)
		args = append(args, userID)
		argIndex++

		// Non-admin hanya bisa melihat active users
		query += fmt.Sprintf(" AND u.status = $%d", argIndex)
		args = append(args, "active")
		argIndex++
	}

	// Apply search filter
	if filters.Search != "" {
		query += fmt.Sprintf(" AND (u.email ILIKE $%d OR u.name ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+filters.Search+"%")
		argIndex++
	}

	// Apply role filter (hanya untuk admin)
	if filters.Role != "" && filters.Role != "all" {
		if userRole == "admin" || userRole == "super_admin" {
			query += fmt.Sprintf(" AND u.role = $%d", argIndex)
			args = append(args, filters.Role)
			argIndex++
		}
	}

	// Apply status filter (hanya untuk admin)
	if filters.Status != "" && filters.Status != "all" {
		if userRole == "admin" || userRole == "super_admin" {
			query += fmt.Sprintf(" AND u.status = $%d", argIndex)
			args = append(args, filters.Status)
			argIndex++
		}
	}

	// Order by - dengan validasi untuk mencegah SQL injection
	allowedOrderFields := map[string]bool{
		"u.created_at":  true,
		"u.updated_at":  true,
		"u.email":       true,
		"u.name":        true,
		"u.role":        true,
		"u.status":      true,
		"u.last_active": true,
	}

	orderBy := "u.created_at DESC"
	if filters.OrderBy != "" {
		orderField := filters.OrderBy
		orderDir := "ASC"

		parts := strings.Fields(orderField)
		if len(parts) == 2 {
			orderField = parts[0]
			if strings.ToUpper(parts[1]) == "DESC" {
				orderDir = "DESC"
			}
		}

		if allowedOrderFields[orderField] {
			orderBy = fmt.Sprintf("%s %s", orderField, orderDir)
		}
	}
	query += fmt.Sprintf(" ORDER BY %s", orderBy)

	// Pagination
	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filters.Limit)
		argIndex++
	} else {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, 100)
		argIndex++
	}

	if filters.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filters.Offset)
		argIndex++
	}

	return query, args
}
