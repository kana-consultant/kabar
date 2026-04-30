package team

import (
	"log"
	"net/http"
	"seo-backend/internal/service/team"

	"github.com/go-chi/chi/v5"
)

// GetTeamMembers returns all members of a team
func (h *TeamHandler) GetTeamMembers(w http.ResponseWriter, r *http.Request) {
	useCtx := h.getUserContext(r)
	teamID := chi.URLParam(r, "id")

	filters := team.MemberFilters{
		Role: r.URL.Query().Get("role"),
	}

	members, err := h.teamService.GetTeamMembers(teamID, filters, useCtx)
	if err != nil {
		log.Printf("Failed to fetch team members: %v", err)
		h.writeError(w, "Failed to fetch team members", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, members, http.StatusOK)
}
