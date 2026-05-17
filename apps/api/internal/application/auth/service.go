package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"seo-backend/internal/domain/auth"
	"seo-backend/internal/models"
)

// Service handles authentication business logic
type Service struct {
	repo     auth.Repository
	tokenGen auth.TokenGenerator
	db       *sql.DB
}

// NewService creates a new auth service
func NewService(db *sql.DB, repo auth.Repository, tokenGen auth.TokenGenerator) auth.AuthService {
	return &Service{
		db:       db,
		repo:     repo,
		tokenGen: tokenGen,
	}
}

// Login authenticates a user
func (s *Service) Login(ctx context.Context, req auth.LoginRequest) (*auth.LoginResponse, error) {
	// Get user by email
	user, passwordHash, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Get team ID
	teamID, _ := s.repo.GetTeamIDByUserID(ctx, user.ID)

	// Update last active
	_ = s.repo.UpdateLastActive(ctx, user.ID)

	// Generate token
	token, err := s.tokenGen.GenerateToken(user.ID, teamID, user.Email, user.Name, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &auth.LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// Register creates a new user with team
func (s *Service) Register(ctx context.Context, req auth.RegisterRequest) (*models.User, error) {
	// Validate request
	if err := s.validateRegisterRequest(req); err != nil {
		return nil, err
	}

	// Check if user exists
	exists, err := s.repo.UserExists(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return nil, errors.New("user already exists")
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Begin transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create user
	user, err := s.repo.CreateUser(ctx, tx, req.Email, req.Name, passwordHash, models.RoleManager)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create team for user
	teamID, err := s.repo.CreateTeamForUser(ctx, tx, user.ID, req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	// Add user to team
	if err := s.repo.AddUserToTeam(ctx, tx, teamID, user.ID); err != nil {
		return nil, fmt.Errorf("failed to add user to team: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return user, nil
}

// GetMeRequest DTO
type GetMeRequest struct {
	UserID string
}

// GetMe returns current user info
func (s *Service) GetMe(ctx context.Context, userID string) (*models.User, error) {
	if userID == "" {
		return nil, errors.New("unauthorized")
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// ChangePassword updates user's password
func (s *Service) ChangePassword(ctx context.Context, req auth.ChangePasswordRequest) error {
	if req.UserID == "" {
		return errors.New("unauthorized")
	}

	// Validate new password
	if len(req.NewPassword) < 6 {
		return errors.New("new password must be at least 6 characters")
	}

	// Get current password hash
	currentHash, err := s.repo.GetPasswordHash(ctx, req.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(req.OldPassword)); err != nil {
		return errors.New("invalid old password")
	}

	// Hash new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Begin transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update password
	if err := s.repo.UpdatePassword(ctx, tx, req.UserID, newHash); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ForgotPassword sends reset link (placeholder)
func (s *Service) ForgotPassword(ctx context.Context, req auth.ForgotPasswordRequest) error {
	// Check if user exists (don't reveal for security)
	exists, err := s.repo.UserExists(ctx, req.Email)
	if err != nil || !exists {
		// Return success even if user doesn't exist for security
		return nil
	}

	// TODO: Send email with reset token
	// For now, just return success

	return nil
}

// Helper: Validate register request
func (s *Service) validateRegisterRequest(req auth.RegisterRequest) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}
