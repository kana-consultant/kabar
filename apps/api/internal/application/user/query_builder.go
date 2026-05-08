package user

import (
	"fmt"

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
		SELECT id, email, name, role, avatar, status, last_active, created_at, updated_at
		FROM users
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	// Apply role-based filtering
	switch ctx.GetRole() {
	case "admin", "super_admin":
		// Admins can see all users
	default:
		// Non-admins can only see active users
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, "active")
		argIndex++
	}

	// Apply search filter
	if filters.Search != "" {
		query += fmt.Sprintf(" AND (email ILIKE $%d OR name ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+filters.Search+"%")
		argIndex++
	}

	// Apply role filter
	if filters.Role != "" {
		query += fmt.Sprintf(" AND role = $%d", argIndex)
		args = append(args, filters.Role)
		argIndex++
	}

	// Apply status filter
	if filters.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filters.Status)
		argIndex++
	}

	// Order by
	orderBy := "created_at DESC"
	if filters.OrderBy != "" {
		orderBy = filters.OrderBy
	}
	query += fmt.Sprintf(" ORDER BY %s", orderBy)

	// Pagination
	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filters.Limit)
		argIndex++
	}
	if filters.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filters.Offset)
		argIndex++
	}

	return query, args
}
