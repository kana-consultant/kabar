package product

import (
	"fmt"

	"seo-backend/internal/models"
)

type Authorizer struct{}

func NewAuthorizer() *Authorizer {
	return &Authorizer{}
}

func (a *Authorizer) CanAccess(product *models.Product, ctx UserContext) bool {
	if product == nil {
		return false
	}

	switch ctx.GetRole() {
	case "super_admin":
		return true
	case "admin":
		if product.TeamID != nil && ctx.GetTeamID() != "" && *product.TeamID == ctx.GetTeamID() {
			return true
		}
		if product.UserID != nil && *product.UserID == ctx.GetUserID() {
			return true
		}
		return false
	default:
		return product.UserID != nil && *product.UserID == ctx.GetUserID()
	}
}

func (a *Authorizer) ValidateAdminAccess(ctx UserContext) error {
	if !ctx.IsAdmin() {
		return fmt.Errorf("admin access required")
	}
	return nil
}
