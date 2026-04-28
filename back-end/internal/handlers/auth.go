package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

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

func GetTeamIDFromDB(userID string) (string, error) {
	var teamID string

	qb := builder.NewQueryBuilder("team_members")

	query := qb.
		Select("team_id").
		WhereEq("user_id", userID).
		Limit(1)

	sqlQuery, args, err := query.Build()
	if err != nil {
		return "", err
	}

	err = database.GetDB().
		QueryRow(sqlQuery, args...).
		Scan(&teamID)

	if err != nil {
		return "", err
	}

	log.Printf("TEAM ID FROM DB: %s", teamID)

	return teamID, nil
}

// Login user
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Printf("Login attempt: email=%s", req.Email)

	var user models.User
	var passwordHash string

	qb := builder.NewQueryBuilder("users")
	sqlQuery, args, err := qb.Select(
		"id", "email", "name", "role", "avatar", "status", "last_active",
		"created_at", "updated_at", "password_hash",
	).WhereEq("email", req.Email).Build()

	if err != nil {
		log.Printf("Failed to build query: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = database.GetDB().QueryRow(sqlQuery, args...).Scan(
		&user.ID, &user.Email, &user.Name, &user.Role, &user.Avatar,
		&user.Status, &user.LastActive, &user.CreatedAt, &user.UpdatedAt,
		&passwordHash,
	)
	if err != nil {
		log.Printf("User not found: %v", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		log.Printf("Password mismatch: %v", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	team_id, err := GetTeamIDFromDB(user.ID)

	log.Printf("Login successful: %s", user.Email)

	// Update last_active using builder
	updateQb := builder.NewQueryBuilder("users")
	updateQuery, updateArgs, err := updateQb.Update().
		Set("last_active", time.Now()).
		WhereEq("id", user.ID).
		Build()

	if err == nil {
		_, _ = database.GetDB().Exec(updateQuery, updateArgs...)
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"team_id": team_id,
		"email":   user.Email,
		"name":    user.Name,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Duration(h.cfg.JWTExpiry) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Don't return password hash
	user.PasswordHash = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.LoginResponse{
		Token: tokenString,
		User:  user,
	})
}

// Register new user - langsung jadi ADMIN
// Register new user - langsung jadi ADMIN dan buat tim sendiri
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Email == "" || req.Name == "" || req.Password == "" {
		http.Error(w, "Email, name, and password are required", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	var exists bool
	err := database.GetDB().QueryRow(`
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
	`, req.Email).Scan(&exists)
	if err != nil {
		log.Printf("Failed to check user existence: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if exists {
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

	// User yang daftar langsung jadi ADMIN
	role := models.RoleAdmin

	// Insert user
	var user models.User
	query := `
		INSERT INTO users (id, email, name, password_hash, role, status) 
		VALUES (gen_random_uuid(), $1, $2, $3, $4, 'active') 
		RETURNING id, email, name, role, avatar, status, last_active, created_at, updated_at
	`
	err = tx.QueryRow(query, req.Email, req.Name, passwordHash, role).Scan(
		&user.ID, &user.Email, &user.Name, &user.Role, &user.Avatar,
		&user.Status, &user.LastActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Buat team baru untuk user ini
	teamName := req.Name + "'s Team"
	var teamID string
	err = tx.QueryRow(`
		INSERT INTO teams (id, name, description, created_by, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, NOW(), NOW())
		RETURNING id
	`, teamName, "Team untuk "+req.Name, user.ID).Scan(&teamID)
	if err != nil {
		log.Printf("Failed to create team: %v", err)
		http.Error(w, "Failed to create team", http.StatusInternalServerError)
		return
	}

	// Tambahkan user sebagai member team dengan role manager
	_, err = tx.Exec(`
		INSERT INTO team_members (id, team_id, user_id, role, joined_at)
		VALUES (gen_random_uuid(), $1, $2, 'manager', NOW())
	`, teamID, user.ID)
	if err != nil {
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
	log.Printf("Created team: %s with ID: %s for user: %s", teamName, teamID, user.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Get current user info
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	// Get user_id from context (set by JWT middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	qb := builder.NewQueryBuilder("users")
	sqlQuery, args, err := qb.Select(
		"id", "email", "name", "role", "avatar", "status", "last_active", "created_at", "updated_at",
	).WhereEq("id", userID).Build()

	if err != nil {
		log.Printf("Failed to build query: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var user models.User
	err = database.GetDB().QueryRow(sqlQuery, args...).Scan(
		&user.ID, &user.Email, &user.Name, &user.Role, &user.Avatar,
		&user.Status, &user.LastActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Change password
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Get user_id from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get current password hash
	var currentHash string
	err := database.GetDB().QueryRow(`
		SELECT password_hash FROM users WHERE id = $1
	`, userID).Scan(&currentHash)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(req.OldPassword)); err != nil {
		http.Error(w, "Invalid old password", http.StatusUnauthorized)
		return
	}

	// Hash new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Update password using builder
	qb := builder.NewQueryBuilder("users")
	sqlQuery, args, err := qb.Update().
		Set("password_hash", newHash).
		Set("updated_at", time.Now()).
		WhereEq("id", userID).
		Build()

	if err != nil {
		log.Printf("Failed to build update query: %v", err)
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	_, err = database.GetDB().Exec(sqlQuery, args...)
	if err != nil {
		log.Printf("Failed to update password: %v", err)
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Password changed successfully"})
}

// Forgot password - request reset
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Check if user exists
	var exists bool
	err := database.GetDB().QueryRow(`
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
	`, req.Email).Scan(&exists)
	if err != nil || !exists {
		// Don't reveal if user exists or not for security
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "If your email is registered, you will receive a reset link"})
		return
	}

	// TODO: Send email with reset token
	// For now, just return success

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "If your email is registered, you will receive a reset link"})
}

// Logout (optional, mostly frontend handles token removal)
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// JWT is stateless, so logout is handled by frontend removing the token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}
