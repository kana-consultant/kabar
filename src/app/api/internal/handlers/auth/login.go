package auth_handler

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"seo-backend/internal/models"
)

// Login user
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Printf("Login attempt: email=%s", req.Email)

	// Get user by email
	user, passwordHash, err := h.getUserByEmail(req.Email)
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

	// Get team ID
	teamID, err := h.getTeamIDFromDB(user.ID)
	if err != nil {
		log.Printf("Failed to get team ID: %v", err)
		teamID = "" // Allow login even without team
	}

	log.Printf("Login successful: %s", user.Email)

	// Update last active
	h.updateLastActive(user.ID)

	// Generate JWT token
	tokenString, err := h.generateToken(user.ID, teamID, user.Email, user.Name, string(user.Role))
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
		User:  *user,
	})
}
