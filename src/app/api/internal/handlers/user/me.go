package user

import (
	"log"
	"net/http"
)

// GetCurrentUser returns the currently authenticated user
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)

	user, err := h.userService.GetCurrentUser(userCtx)
	if err != nil {
		log.Printf("Failed to fetch current user: %v", err)
		h.handleServiceError(w, err)
		return
	}

	h.writeJSON(w, user, http.StatusOK)
}
