package team

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetByID returns a specific team by ID
func (h *TeamHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)
	id := chi.URLParam(r, "id")

	team, err := h.teamService.GetByID(id, userCtx)
	if err != nil {
		log.Printf("Failed to fetch team: %v", err)
		if err.Error() == "access denied" {
			h.writeError(w, "Forbidden", http.StatusForbidden)
		} else {
			h.writeError(w, "Team not found", http.StatusNotFound)
		}
		return
	}

	if team == nil {
		h.writeError(w, "Team not found", http.StatusNotFound)
		return
	}

	h.writeJSON(w, team, http.StatusOK)
}
