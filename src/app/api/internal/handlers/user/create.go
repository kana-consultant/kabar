package user

import (
	"encoding/json"
	"log"
	"net/http"

	"seo-backend/internal/models"
	"seo-backend/internal/service/user"
)

// Create creates a new user
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if err := validateCreateUserRequest(req.Email, req.Name, req.Password); err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	createReq := user.CreateUserRequest{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
		Role:     models.UserRole(req.Role),
	}

	user, err := h.userService.Create(createReq)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		h.handleServiceError(w, err)
		return
	}

	h.writeJSON(w, user, http.StatusCreated)
}
