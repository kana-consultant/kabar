package user

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Delete removes a user
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)
	id := chi.URLParam(r, "id")

	if err := h.userService.Delete(id, userCtx); err != nil {
		log.Printf("Failed to delete user: %v", err)
		h.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
