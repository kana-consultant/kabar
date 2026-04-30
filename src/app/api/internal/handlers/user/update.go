package user

import (
	"encoding/json"
	"log"
	"net/http"

	"seo-backend/internal/models"
	"seo-backend/internal/service/user"

	"github.com/go-chi/chi/v5"
)

// Update updates an existing user
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)
	id := chi.URLParam(r, "id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		h.writeError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	updateReq := user.UpdateUserRequest{}

	if name, ok := updates["name"].(string); ok {
		updateReq.Name = &name
	}
	if email, ok := updates["email"].(string); ok {
		updateReq.Email = &email
	}
	if role, ok := updates["role"].(string); ok {
		roleVal := models.UserRole(role)
		updateReq.Role = &roleVal
	}
	if status, ok := updates["status"].(string); ok {
		updateReq.Status = &status
	}

	if err := h.userService.Update(id, updateReq, userCtx); err != nil {
		log.Printf("Failed to update user: %v", err)
		h.handleServiceError(w, err)
		return
	}

	h.writeJSON(w, map[string]string{"message": "User updated successfully"}, http.StatusOK)
}
