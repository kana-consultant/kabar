package auth_handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"seo-backend/internal/database"
	"seo-backend/internal/models"
)

// Register new user - langsung jadi ADMIN dan buat tim sendiri
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate input
	if err := h.validateRegisterRequest(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user already exists
	if exists, err := h.userExists(req.Email); err != nil || exists {
		if err != nil {
			log.Printf("Failed to check user existence: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Start transaction
	tx, err := database.GetDB().Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Create user
	user, err := h.insertUser(tx, req, passwordHash)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Create team for user
	teamID, err := h.createTeamForUser(tx, user.ID, req.Name)
	if err != nil {
		log.Printf("Failed to create team: %v", err)
		http.Error(w, "Failed to create team", http.StatusInternalServerError)
		return
	}

	// Add user as team member
	if err := h.addUserToTeam(tx, teamID, user.ID); err != nil {
		log.Printf("Failed to add user to team: %v", err)
		http.Error(w, "Failed to add user to team", http.StatusInternalServerError)
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		http.Error(w, "Failed to complete registration", http.StatusInternalServerError)
		return
	}

	log.Printf("User registered as ADMIN: %s (%s)", user.Email, user.Role)
	log.Printf("Created team: %s with ID: %s for user: %s", req.Name+"'s Team", teamID, user.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Helper: Validate register request
func (h *AuthHandler) validateRegisterRequest(req models.RegisterRequest) error {
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}
	if len(req.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	return nil
}

// Helper: Check if user exists
func (h *AuthHandler) userExists(email string) (bool, error) {
	var exists bool
	err := database.GetDB().QueryRow(`
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
	`, email).Scan(&exists)
	return exists, err
}

// Helper: Insert user
func (h *AuthHandler) insertUser(tx interface {
	QueryRow(string, ...interface{}) *sql.Row
}, req models.RegisterRequest, passwordHash []byte) (*models.User, error) {
	var user models.User
	query := `
		INSERT INTO users (id, email, name, password_hash, role, status) 
		VALUES (gen_random_uuid(), $1, $2, $3, $4, 'active') 
		RETURNING id, email, name, role, avatar, status, last_active, created_at, updated_at
	`
	err := tx.QueryRow(query, req.Email, req.Name, passwordHash, models.RoleAdmin).Scan(
		&user.ID, &user.Email, &user.Name, &user.Role, &user.Avatar,
		&user.Status, &user.LastActive, &user.CreatedAt, &user.UpdatedAt,
	)
	return &user, err
}

// Helper: Create team for user
func (h *AuthHandler) createTeamForUser(tx interface {
	QueryRow(string, ...interface{}) *sql.Row
}, userID, userName string) (string, error) {
	teamName := userName + "'s Team"
	var teamID string
	err := tx.QueryRow(`
		INSERT INTO teams (id, name, description, created_by, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, NOW(), NOW())
		RETURNING id
	`, teamName, "Team untuk "+userName, userID).Scan(&teamID)
	return teamID, err
}

// Helper: Add user to team
func (h *AuthHandler) addUserToTeam(tx interface {
	Exec(string, ...interface{}) (sql.Result, error)
}, teamID, userID string) error {
	_, err := tx.Exec(`
		INSERT INTO team_members (id, team_id, user_id, role, joined_at)
		VALUES (gen_random_uuid(), $1, $2, 'manager', NOW())
	`, teamID, userID)
	return err
}
