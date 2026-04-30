package auth_handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"seo-backend/internal/builder"
	"seo-backend/internal/database"
)

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

	// Validate new password
	if len(req.NewPassword) < 6 {
		http.Error(w, "New password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	// Get user_id from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get current password hash
	currentHash, err := h.getPasswordHash(userID)
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

	// Update password
	if err := h.updatePassword(userID, newHash); err != nil {
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

// Helper: Get password hash
func (h *AuthHandler) getPasswordHash(userID string) (string, error) {
	var currentHash string
	err := database.GetDB().QueryRow(`
		SELECT password_hash FROM users WHERE id = $1
	`, userID).Scan(&currentHash)
	return currentHash, err
}

// Helper: Update password
func (h *AuthHandler) updatePassword(userID string, newHash []byte) error {
	qb := builder.NewQueryBuilder("users")
	sqlQuery, args, err := qb.Update().
		Set("password_hash", newHash).
		Set("updated_at", time.Now()).
		WhereEq("id", userID).
		Build()

	if err != nil {
		return err
	}

	_, err = database.GetDB().Exec(sqlQuery, args...)
	return err
}
