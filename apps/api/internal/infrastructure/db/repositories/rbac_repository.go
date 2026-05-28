// internal/infrastructure/db/repositories/rbac_repository.go
package repositories

import (
	"context"
	"database/sql"
)

type Permission struct {
	Module string
	Action string
	Scope  string // "global" | "team"
}

type RBACRepository struct {
	db *sql.DB
}

func NewRBACRepository(db *sql.DB) *RBACRepository {
	return &RBACRepository{db: db}
}

// GetPermissionsByRole → query permission berdasarkan role name
func (r *RBACRepository) GetPermissionsByRole(ctx context.Context, roleName string) ([]Permission, error) {
	query := `
        SELECT p.module, p.action, p.scope
        FROM permissions p
        JOIN role_permissions rp ON rp.permission_id = p.id
        JOIN roles ro            ON ro.id = rp.role_id
        WHERE ro.name = $1
    `

	rows, err := r.db.QueryContext(ctx, query, roleName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var p Permission
		if err := rows.Scan(&p.Module, &p.Action, &p.Scope); err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}
	return permissions, rows.Err()
}
