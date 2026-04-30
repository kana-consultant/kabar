package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"seo-backend/internal/middleware/auth"
	"seo-backend/internal/models"
	services "seo-backend/internal/service/user"
)

type UserHandler struct {
	userService *services.Service
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: services.NewService(),
	}
}

// Helper: Get user context from request
func (h *UserHandler) getUserContext(r *http.Request) models.UserContext {
	ctx := r.Context()
	return &models.SimpleUserContext{
		UserID: auth.GetUserID(ctx),
		TeamID: auth.GetTeamID(ctx),
		Role:   auth.GetUserRole(ctx),
	}
}

// Helper: Write JSON response
func (h *UserHandler) writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

// Helper: Write error response
func (h *UserHandler) writeError(w http.ResponseWriter, message string, status int) {
	h.writeJSON(w, map[string]string{"error": message}, status)
}

// Helper: Handle service errors
func (h *UserHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch err.Error() {
	case "access denied":
		h.writeError(w, "Forbidden", http.StatusForbidden)
	case "user not found":
		h.writeError(w, "User not found", http.StatusNotFound)
	case "unauthorized":
		h.writeError(w, "Unauthorized", http.StatusUnauthorized)
	case "email already exists":
		h.writeError(w, "Email already exists", http.StatusConflict)
	case "cannot delete your own account":
		h.writeError(w, err.Error(), http.StatusBadRequest)
	default:
		log.Printf("Unexpected error: %v", err)
		h.writeError(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Helper: Validate create user request
func validateCreateUserRequest(email, name, password string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if name == "" {
		return fmt.Errorf("name is required")
	}
	if password == "" {
		return fmt.Errorf("password is required")
	}
	if len(password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	return nil
}
