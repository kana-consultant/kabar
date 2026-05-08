// internal/domain/auth/service.go
package auth

import (
	"context"

	"seo-backend/internal/models"
)

// AuthService defines the authentication business logic interface
type AuthService interface {
	ChangePassword(ctx context.Context, req ChangePasswordRequest) error
	ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error
	GetMe(ctx context.Context, userID string) (*models.User, error)
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	Register(ctx context.Context, req RegisterRequest) (*models.User, error)
}
