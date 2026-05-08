// internal/application/team/authorizer.go
package team

import (
	"fmt"

	"seo-backend/internal/domain/team"
	"seo-backend/internal/models"
)

type AuthorizerImpl struct{}

func NewAuthorizer() team.Authorizer {
	return &AuthorizerImpl{}
}

func (a *AuthorizerImpl) CanAccessTeam(teamID string, ctx models.UserContext) bool {
	if ctx.IsAdmin() {
		return true
	}
	return ctx.GetTeamID() == teamID
}

func (a *AuthorizerImpl) CanManageTeam(teamID string, ctx models.UserContext, userRole string) bool {
	if ctx.IsAdmin() {
		return true
	}
	if ctx.GetTeamID() != teamID {
		return false
	}
	return userRole == "manager" || userRole == "admin"
}

func (a *AuthorizerImpl) ValidateTeamAccess(teamID string, ctx models.UserContext) error {
	if !a.CanAccessTeam(teamID, ctx) {
		return fmt.Errorf("access denied")
	}
	return nil
}
