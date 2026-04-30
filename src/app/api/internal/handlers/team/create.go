package team

import (
	"encoding/json"
	"log"
	"net/http"

	"seo-backend/internal/models"
)

// Create creates a new team
func (h *TeamHandler) Create(w http.ResponseWriter, r *http.Request) {
	useCtx := h.getUserContext(r)

	var req models.CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := validateTeamRequest(req); err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	team, err := h.teamService.Create(req, useCtx.GetUserID())
	if err != nil {
		log.Printf("Failed to create team: %v", err)
		h.writeError(w, "Failed to create team", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, team, http.StatusCreated)
}
