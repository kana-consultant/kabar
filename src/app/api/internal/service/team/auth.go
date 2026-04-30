package team

import (
	"fmt"
	"seo-backend/internal/models"
)

type Authorizer struct{}

func NewAuthorizer() *Authorizer {
	return &Authorizer{}
}

func (a *Authorizer) CanAccessTeam(teamID string, ctx models.UserContext) bool {
	if ctx.IsAdmin() {
		return true
	}
	return ctx.GetTeamID() == teamID
}

func (a *Authorizer) CanManageTeam(teamID string, ctx models.UserContext, userRole string) bool {
	if ctx.IsAdmin() {
		return true
	}
	if ctx.GetTeamID() != teamID {
		return false
	}
	// Check if user is manager or admin of the team
	return userRole == "manager" || userRole == "admin"
}

func (a *Authorizer) ValidateTeamAccess(teamID string, ctx models.UserContext) error {
	if !a.CanAccessTeam(teamID, ctx) {
		return fmt.Errorf("access denied")
	}
	return nil
}
