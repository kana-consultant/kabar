package user

import (
	"context"
	"database/sql"
	"errors"

	"seo-backend/internal/models"
)

type Authorizer struct {
	db *sql.DB
}

func NewAuthorizer(db *sql.DB) *Authorizer {
	return &Authorizer{db: db}
}

// ValidateAccess checks if user can access another user's data
func (a *Authorizer) ValidateAccess(ctx context.Context, targetUserID string, userCtx models.UserContext) error {
	if a.CanAccess(ctx, targetUserID, userCtx) {
		return nil
	}
	return errors.New("access denied")
}

// CanAccess checks permission to access user data
func (a *Authorizer) CanAccess(ctx context.Context, targetUserID string, userCtx models.UserContext) bool {
	userRole := userCtx.GetRole()
	currentUserID := userCtx.GetUserID()

	// User can access their own data
	if currentUserID == targetUserID {
		return true
	}

	// Admins can access any user
	if userRole == "admin" || userRole == "super_admin" {
		return true
	}

	return false
}

// CanDelete checks if user can delete another user
func (a *Authorizer) CanDelete(ctx context.Context, targetUserID string, userCtx models.UserContext) bool {
	currentUserID := userCtx.GetUserID()
	userRole := userCtx.GetRole()

	// Cannot delete self
	if currentUserID == targetUserID {
		return false
	}

	// Only admins can delete users
	if userRole == "admin" || userRole == "super_admin" {
		return true
	}

	return false
}
