package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"seo-backend/internal/domain/user"
	"seo-backend/internal/helper"
	"seo-backend/internal/models"
)

// Service handles user business logic
type Service struct {
	db           *sql.DB
	repo         user.Repository
	queryBuilder *QueryBuilder
	authorizer   *Authorizer
	passwordSvc  *PasswordService
	validator    *Validator
}

// NewService creates a new user service
func NewService(
	db *sql.DB,
	repo user.Repository,
	queryBuilder *QueryBuilder,
	authorizer *Authorizer,
	passwordSvc *PasswordService,
	validator *Validator,
) user.UserService {
	return &Service{
		db:           db,
		repo:         repo,
		queryBuilder: queryBuilder,
		authorizer:   authorizer,
		passwordSvc:  passwordSvc,
		validator:    validator,
	}
}

// =======================
// USER CRUD OPERATIONS
// =======================

// GetAll retrieves all users with filters
func (s *Service) GetAll(ctx context.Context, userCtx models.UserContext) ([]models.User, error) {
	// Build query based on user context and filters
	query, args := s.queryBuilder.BuildListQuery(userCtx)

	users, err := s.repo.GetAll(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

// GetByID retrieves a user by ID with access validation
func (s *Service) GetByID(ctx context.Context, id string, userCtx models.UserContext) (*models.User, error) {
	// Validate access
	if err := s.authorizer.ValidateAccess(ctx, id, userCtx); err != nil {
		return nil, err
	}

	// Get user from repository
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// GetCurrentUser retrieves the currently authenticated user
func (s *Service) GetCurrentUser(ctx context.Context, userCtx models.UserContext) (*models.User, error) {
	if userCtx.GetUserID() == "" {
		return nil, errors.New("unauthorized: no user context")
	}

	user, err := s.repo.GetByID(ctx, userCtx.GetUserID())
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// Create creates a new user
func (s *Service) Create(ctx context.Context, req user.CreateUserRequest) (*models.User, error) {
	// Validate request
	if err := s.validator.ValidateCreate(req); err != nil {
		return nil, err
	}

	// Check email uniqueness
	if err := s.validator.ValidateEmailUniqueness(ctx, req.Email); err != nil {
		return nil, err
	}

	// Hash password
	passwordHash, err := s.passwordSvc.Hash(req.Password)
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
	user, err := s.repo.Create(ctx, req, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return user, nil
}

// Update updates an existing user
func (s *Service) Update(ctx context.Context, id string, req user.UpdateUserRequest, userCtx models.UserContext) error {
	// Check access
	if err := s.authorizer.ValidateAccess(ctx, id, userCtx); err != nil {
		return err
	}

	// Get existing user
	existingUser, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	// Build updates map
	updates := make(map[string]interface{})

	if req.Name != nil && *req.Name != "" {
		updates["name"] = *req.Name
	}

	if req.Email != nil && *req.Email != "" {
		// Check email uniqueness if changing
		if err := s.validator.ValidateUpdateEmailUniqueness(ctx, *req.Email, existingUser.Email); err != nil {
			return err
		}
		updates["email"] = *req.Email
	}

	if req.Role != nil && *req.Role != "" {
		updates["role"] = *req.Role
	}

	if req.Status != nil && *req.Status != "" {
		updates["status"] = *req.Status
	}

	if req.Avatar != nil {
		updates["avatar"] = *req.Avatar
	}

	if len(updates) == 0 {
		return errors.New("no fields to update")
	}

	// Begin transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update user
	if err := s.repo.Update(ctx, id, updates); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// UpdatePassword updates user's password
func (s *Service) UpdatePassword(ctx context.Context, id string, oldPassword, newPassword string, userCtx models.UserContext) error {
	// Check access
	if err := s.authorizer.ValidateAccess(ctx, id, userCtx); err != nil {
		return err
	}

	// Get existing user
	existingUser, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	// Validate old password
	if err := s.passwordSvc.Verify(oldPassword, existingUser.PasswordHash); err != nil {
		return errors.New("invalid current password")
	}

	// Hash new password
	newPasswordHash, err := s.passwordSvc.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Begin transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update password
	updates := map[string]interface{}{
		"password_hash": newPasswordHash,
	}
	if err := s.repo.Update(ctx, id, updates); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Delete deletes a user
func (s *Service) Delete(ctx context.Context, id string, userCtx models.UserContext) error {
	// Check if can delete (cannot delete self)
	if !s.authorizer.CanDelete(ctx, id, userCtx) {
		return errors.New("cannot delete this user")
	}

	// Get user
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Begin transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete user
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// =======================
// HELPER METHODS (for other services)
// =======================

// CanAccessUser checks if a user can access another user's data
func (s *Service) CanAccessUser(ctx context.Context, targetUserID string, userCtx models.UserContext) bool {
	return s.authorizer.CanAccess(ctx, targetUserID, userCtx)
}

// GetUserByEmail retrieves a user by email (for auth)
func (s *Service) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	// Build query for email lookup
	query := "SELECT id, email, name, password_hash, role, avatar, status, last_active, created_at, updated_at FROM users WHERE email = $1"

	var user models.User
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.Role,
		&user.Avatar, &user.Status, &user.LastActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// UpdateLastActive updates user's last active timestamp
func (s *Service) UpdateLastActive(ctx context.Context, userID string) error {
	updates := map[string]interface{}{
		"last_active": helper.ParseWIBTime(time.Now().Format(time.RFC3339)),
	}
	return s.repo.Update(ctx, userID, updates)
}
