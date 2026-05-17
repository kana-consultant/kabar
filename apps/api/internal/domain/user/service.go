// internal/domain/user/service.go
package user

import (
	"context"

	"seo-backend/internal/models"
)

// UserService defines the user business logic interface
type UserService interface {
	// =======================
	// USER CRUD OPERATIONS
	// =======================

	// GetAll retrieves all users with filters
	GetAll(ctx context.Context, userCtx models.UserContext) ([]models.User, error)

	// GetByID retrieves a user by ID with access validation
	GetByID(ctx context.Context, id string, userCtx models.UserContext) (*models.User, error)

	// GetCurrentUser retrieves the currently authenticated user
	GetCurrentUser(ctx context.Context, userCtx models.UserContext) (*models.User, error)

	// Create creates a new user
	Create(ctx context.Context, req CreateUserRequest) (*models.User, error)

	// Update updates an existing user
	Update(ctx context.Context, id string, req UpdateUserRequest, userCtx models.UserContext) error

	// UpdatePassword updates user's password
	UpdatePassword(ctx context.Context, id string, oldPassword, newPassword string, userCtx models.UserContext) error

	// Delete deletes a user
	Delete(ctx context.Context, id string, userCtx models.UserContext) error

	// =======================
	// HELPER METHODS (for other services)
	// =======================

	// CanAccessUser checks if a user can access another user's data
	CanAccessUser(ctx context.Context, targetUserID string, userCtx models.UserContext) bool

	// GetUserByEmail retrieves a user by email (for auth)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)

	// UpdateLastActive updates user's last active timestamp
	UpdateLastActive(ctx context.Context, userID string) error
}
