package auth_handler

import (
	"encoding/json"
	"net/http"
)

// Get current user info
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	// Get user_id from context (set by JWT middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.getUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
