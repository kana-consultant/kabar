package auth

import (
	"context"
	"database/sql"
	"seo-backend/internal/models"
)

// Repository interface for auth data access
type Repository interface {
	// User operations
	GetUserByEmail(ctx context.Context, email string) (*models.User, string, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetPasswordHash(ctx context.Context, userID string) (string, error)
	UserExists(ctx context.Context, email string) (bool, error)

	// User creation
	CreateUser(ctx context.Context, tx *sql.Tx, email, name string, passwordHash []byte, role models.UserRole) (*models.User, error)

	// User update
	UpdatePassword(ctx context.Context, tx *sql.Tx, userID string, newHash []byte) error
	UpdateLastActive(ctx context.Context, userID string) error

	// Team operations
	CreateTeamForUser(ctx context.Context, tx *sql.Tx, userID, userName string) (string, error)
	AddUserToTeam(ctx context.Context, tx *sql.Tx, teamID, userID string) error
	GetTeamIDByUserID(ctx context.Context, userID string) (string, error)

	// Transaction
	BeginTx(ctx context.Context) (*sql.Tx, error)
}

// TokenGenerator interface for JWT operations
type TokenGenerator interface {
	GenerateToken(userID, teamID, email, name, role string, perm []string) (string, error)
}
