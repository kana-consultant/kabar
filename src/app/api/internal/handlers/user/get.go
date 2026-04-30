package user

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetByID returns a specific user by ID
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)
	id := chi.URLParam(r, "id")

	user, err := h.userService.GetByID(id, userCtx)
	if err != nil {
		log.Printf("Failed to fetch user: %v", err)
		h.handleServiceError(w, err)
		return
	}

	if user == nil {
		h.writeError(w, "User not found", http.StatusNotFound)
		return
	}

	h.writeJSON(w, user, http.StatusOK)
}
