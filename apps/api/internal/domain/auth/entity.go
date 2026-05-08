// internal/domain/auth/dto.go
package auth

import "seo-backend/internal/models"

// ========== REQUEST DTOs ==========

// LoginRequest DTO
type LoginRequest struct {
	Email    string
	Password string
}

// RegisterRequest DTO
type RegisterRequest struct {
	Email    string
	Name     string
	Password string
}

// GetMeRequest DTO
type GetMeRequest struct {
	UserID string
}

// ChangePasswordRequest DTO
type ChangePasswordRequest struct {
	UserID      string
	OldPassword string
	NewPassword string
}

// ForgotPasswordRequest DTO
type ForgotPasswordRequest struct {
	Email string
}

// ========== RESPONSE DTOs ==========

// LoginResponse DTO
type LoginResponse struct {
	Token string
	User  *models.User
}
