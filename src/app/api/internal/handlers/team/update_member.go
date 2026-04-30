package team

import (
	"encoding/json"
	"log"
	"net/http"

	"seo-backend/internal/models"

	"github.com/go-chi/chi/v5"
)

// UpdateMemberRole updates a team member's role
func (h *TeamHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)
	teamID := chi.URLParam(r, "id")
	userID := chi.URLParam(r, "userId")

	var role models.TeamMemberRole
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		h.writeError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if role == "" {
		h.writeError(w, "role is required", http.StatusBadRequest)
		return
	}

	team, err := h.teamService.UpdateMemberRole(teamID, userID, role, userCtx)
	if err != nil {
		log.Printf("Failed to update member role: %v", err)
		if err.Error() == "member not found" {
			h.writeError(w, "Member not found", http.StatusNotFound)
		} else {
			h.writeError(w, "Failed to update member role", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSON(w, team, http.StatusOK)
}
