package auth_handler

import (
	"encoding/json"
	"net/http"
)

// Logout (JWT is stateless, so logout is handled by frontend removing the token)
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}
