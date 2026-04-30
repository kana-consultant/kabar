package history

import (
	"fmt"
	"seo-backend/internal/models"
)

type Authorizer struct{}

func NewAuthorizer() *Authorizer {
	return &Authorizer{}
}

func (a *Authorizer) CanAccess(history *models.History, ctx models.UserContext) bool {
	if history == nil {
		return false
	}

	switch ctx.GetRole() {
	case "super_admin":
		return true
	case "admin":
		// Admin can access if same team or created by them
		if history.TeamID != nil && ctx.GetTeamID() != "" && *history.TeamID == ctx.GetTeamID() {
			return true
		}
		if history.CreatedBy != nil && *history.CreatedBy == ctx.GetUserID() {
			return true
		}
		return false
	default:
		// Regular user can only access their own history
		return history.CreatedBy != nil && *history.CreatedBy == ctx.GetUserID()
	}
}

func (a *Authorizer) ValidateAdminAccess(ctx models.UserContext) error {
	if !ctx.IsAdmin() {
		return fmt.Errorf("admin access required")
	}
	return nil
}
