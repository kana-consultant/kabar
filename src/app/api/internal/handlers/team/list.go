package team

import (
	"log"
	"net/http"
	"seo-backend/internal/service/team"
)

// GetAll returns all teams based on user permissions
func (h *TeamHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)

	filters := team.TeamFilters{
		Status: r.URL.Query().Get("status"),
	}

	teams, err := h.teamService.GetAll(userCtx, filters)
	if err != nil {
		log.Printf("Failed to fetch teams: %v", err)
		h.writeError(w, "Failed to fetch teams", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, teams, http.StatusOK)
}
