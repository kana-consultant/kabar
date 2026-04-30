package team

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Delete removes a team
func (h *TeamHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)
	id := chi.URLParam(r, "id")

	if err := h.teamService.Delete(id, userCtx); err != nil {
		log.Printf("Failed to delete team: %v", err)
		h.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
