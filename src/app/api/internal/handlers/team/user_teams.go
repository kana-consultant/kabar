package team

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetUserTeams returns all teams a user belongs to
func (h *TeamHandler) GetUserTeams(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")

	if userID == "" {
		h.writeError(w, "user_id is required", http.StatusBadRequest)
		return
	}

	teams, err := h.teamService.GetUserTeams(userID)
	if err != nil {
		log.Printf("Failed to fetch user teams: %v", err)
		h.writeError(w, "Failed to fetch user teams", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, teams, http.StatusOK)
}
