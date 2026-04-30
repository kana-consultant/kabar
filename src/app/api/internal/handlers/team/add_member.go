package team

import (
	"encoding/json"
	"log"
	"net/http"

	"seo-backend/internal/models"

	"github.com/go-chi/chi/v5"
)

// AddMember adds a user to a team
func (h *TeamHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	ctx := h.getUserContext(r)
	teamID := chi.URLParam(r, "id")

	var req models.AddTeamMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		h.writeError(w, "user_id is required", http.StatusBadRequest)
		return
	}

	team, err := h.teamService.AddMember(teamID, req, ctx)
	if err != nil {
		log.Printf("Failed to add member: %v", err)
		h.handleServiceError(w, err)
		return
	}

	h.writeJSON(w, team, http.StatusCreated)
}
