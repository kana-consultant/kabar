package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"seo-backend/internal/domain/auth"
	"seo-backend/internal/models"
)

type AuthRepository struct {
	db *sql.DB
}

// NewAuthRepository creates a new auth repository
func NewAuthRepository(db *sql.DB) auth.Repository {
	return &AuthRepository{db: db}
}

// BeginTx starts a new transaction
func (r *AuthRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

// GetUserByEmail retrieves user by email with password hash
func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, string, error) {
	var user models.User
	var passwordHash string

	query := `
		SELECT id, email, name, password_hash, role, avatar, status, last_active, created_at, updated_at
		FROM users WHERE email = $1
	`

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Name, &passwordHash, &user.Role,
		&user.Avatar, &user.Status, &user.LastActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", nil
		}
		return nil, "", fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, passwordHash, nil
}

// GetUserByID retrieves user by ID
func (r *AuthRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User

	query := `
		SELECT id, email, name, role, avatar, status, last_active, created_at, updated_at
		FROM users WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.Role,
		&user.Avatar, &user.Status, &user.LastActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

// GetPasswordHash retrieves password hash for a user
func (r *AuthRepository) GetPasswordHash(ctx context.Context, userID string) (string, error) {
	var passwordHash string
	err := r.db.QueryRowContext(ctx, `
		SELECT password_hash FROM users WHERE id = $1
	`, userID).Scan(&passwordHash)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found")
		}
		return "", fmt.Errorf("failed to get password hash: %w", err)
	}

	return passwordHash, nil
}

// UserExists checks if a user exists by email
func (r *AuthRepository) UserExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
	`, email).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

// CreateUser creates a new user
func (r *AuthRepository) CreateUser(ctx context.Context, tx *sql.Tx, email, name string, passwordHash []byte, role models.UserRole) (*models.User, error) {
	if tx == nil {
		return nil, fmt.Errorf("transaction is required for CreateUser")
	}

	var user models.User
	query := `
		INSERT INTO users (id, email, name, password_hash, role, status) 
		VALUES (gen_random_uuid(), $1, $2, $3, $4, 'active') 
		RETURNING id, email, name, role, avatar, status, last_active, created_at, updated_at
	`

	err := tx.QueryRowContext(ctx, query, email, name, passwordHash, role).Scan(
		&user.ID, &user.Email, &user.Name, &user.Role, &user.Avatar,
		&user.Status, &user.LastActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// UpdatePassword updates user's password
func (r *AuthRepository) UpdatePassword(ctx context.Context, tx *sql.Tx, userID string, newHash []byte) error {
	if tx == nil {
		return fmt.Errorf("transaction is required for UpdatePassword")
	}

	_, err := tx.ExecContext(ctx, `
		UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2
	`, newHash, userID)

	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// UpdateLastActive updates user's last active timestamp
func (r *AuthRepository) UpdateLastActive(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET last_active = NOW() WHERE id = $1
	`, userID)

	if err != nil {
		return fmt.Errorf("failed to update last active: %w", err)
	}

	return nil
}

// CreateTeamForUser creates a team for a user
func (r *AuthRepository) CreateTeamForUser(ctx context.Context, tx *sql.Tx, userID, userName string) (string, error) {
	if tx == nil {
		return "", fmt.Errorf("transaction is required for CreateTeamForUser")
	}

	teamName := userName + "'s Team"
	var teamID string

	err := tx.QueryRowContext(ctx, `
		INSERT INTO teams (id, name, description, created_by, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, NOW(), NOW())
		RETURNING id
	`, teamName, "Team untuk "+userName, userID).Scan(&teamID)

	if err != nil {
		return "", fmt.Errorf("failed to create team: %w", err)
	}

	return teamID, nil
}

// AddUserToTeam adds a user to a team
func (r *AuthRepository) AddUserToTeam(ctx context.Context, tx *sql.Tx, teamID, userID string) error {
	if tx == nil {
		return fmt.Errorf("transaction is required for AddUserToTeam")
	}

	_, err := tx.ExecContext(ctx, `
		INSERT INTO team_members (team_id, user_id, role, joined_at)
		VALUES ($1, $2, 'manager', NOW())
	`, teamID, userID)

	if err != nil {
		return fmt.Errorf("failed to add user to team: %w", err)
	}

	return nil
}

// GetTeamIDByUserID retrieves team ID for a user
func (r *AuthRepository) GetTeamIDByUserID(ctx context.Context, userID string) (string, error) {
	var teamID string
	err := r.db.QueryRowContext(ctx, `
		SELECT team_id FROM team_members WHERE user_id = $1 LIMIT 1
	`, userID).Scan(&teamID)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("failed to get team ID: %w", err)
	}

	return teamID, nil
}
