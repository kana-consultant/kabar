package user

import (
	"database/sql"
	"fmt"
	"seo-backend/internal/models"
)

type Authorizer struct {
	db *sql.DB
}

func NewAuthorizer(db *sql.DB) *Authorizer {
	return &Authorizer{db: db}
}

func (a *Authorizer) CanAccess(targetUserID string, ctx models.UserContext) bool {
	switch ctx.GetRole() {
	case "super_admin":
		return true
	case "admin":
		// Check if user is in same team
		var inSameTeam bool
		query := `
			SELECT EXISTS(
				SELECT 1 FROM team_members tm1
				JOIN team_members tm2 ON tm1.team_id = tm2.team_id
				WHERE tm1.user_id = $1 AND tm2.user_id = $2
			)
		`
		err := a.db.QueryRow(query, ctx.GetUserID(), targetUserID).Scan(&inSameTeam)
		if err == nil && inSameTeam {
			return true
		}
		// Fallback to self access
		return targetUserID == ctx.GetUserID()
	default:
		return targetUserID == ctx.GetUserID()
	}
}

func (a *Authorizer) ValidateAccess(targetUserID string, ctx models.UserContext) error {
	if !a.CanAccess(targetUserID, ctx) {
		return fmt.Errorf("access denied")
	}
	return nil
}

func (a *Authorizer) CanDelete(targetUserID string, ctx models.UserContext) bool {
	// Prevent deleting self
	if targetUserID == ctx.GetUserID() {
		return false
	}
	return a.CanAccess(targetUserID, ctx)
}
