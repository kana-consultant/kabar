package auth_handler

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"seo-backend/internal/builder"
	"seo-backend/internal/config"
	"seo-backend/internal/database"
	"seo-backend/internal/models"
)

type AuthHandler struct {
	cfg *config.Config
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

// Helper: Get team ID from database
func (h *AuthHandler) getTeamIDFromDB(userID string) (string, error) {
	var teamID string

	qb := builder.NewQueryBuilder("team_members")
	query := qb.Select("team_id").WhereEq("user_id", userID).Limit(1)
	sqlQuery, args, err := query.Build()
	if err != nil {
		return "", err
	}

	err = database.GetDB().QueryRow(sqlQuery, args...).Scan(&teamID)
	if err != nil {
		return "", err
	}

	log.Printf("TEAM ID FROM DB: %s", teamID)
	return teamID, nil
}

// Helper: Generate JWT token
func (h *AuthHandler) generateToken(userID, teamID, email, name, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"team_id": teamID,
		"email":   email,
		"name":    name,
		"role":    role,
		"exp":     time.Now().Add(time.Duration(h.cfg.JWTExpiry) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.cfg.JWTSecret))
}

// Helper: Update last active
func (h *AuthHandler) updateLastActive(userID string) {
	updateQb := builder.NewQueryBuilder("users")
	updateQuery, updateArgs, err := updateQb.Update().
		Set("last_active", time.Now()).
		WhereEq("id", userID).
		Build()

	if err == nil {
		_, _ = database.GetDB().Exec(updateQuery, updateArgs...)
	}
}

// Helper: Get user by email
func (h *AuthHandler) getUserByEmail(email string) (*models.User, string, error) {
	var user models.User
	var passwordHash string

	qb := builder.NewQueryBuilder("users")
	sqlQuery, args, err := qb.Select(
		"id", "email", "name", "role", "avatar", "status", "last_active",
		"created_at", "updated_at", "password_hash",
	).WhereEq("email", email).Build()

	if err != nil {
		return nil, "", err
	}

	err = database.GetDB().QueryRow(sqlQuery, args...).Scan(
		&user.ID, &user.Email, &user.Name, &user.Role, &user.Avatar,
		&user.Status, &user.LastActive, &user.CreatedAt, &user.UpdatedAt,
		&passwordHash,
	)
	if err != nil {
		return nil, "", err
	}

	return &user, passwordHash, nil
}

// Helper: Get user by ID
func (h *AuthHandler) getUserByID(userID string) (*models.User, error) {
	qb := builder.NewQueryBuilder("users")
	sqlQuery, args, err := qb.Select(
		"id", "email", "name", "role", "avatar", "status", "last_active", "created_at", "updated_at",
	).WhereEq("id", userID).Build()

	if err != nil {
		return nil, err
	}

	var user models.User
	err = database.GetDB().QueryRow(sqlQuery, args...).Scan(
		&user.ID, &user.Email, &user.Name, &user.Role, &user.Avatar,
		&user.Status, &user.LastActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
