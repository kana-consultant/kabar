package team

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Update updates an existing team
func (h *TeamHandler) Update(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)
	id := chi.URLParam(r, "id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		h.writeError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.teamService.Update(id, updates, userCtx); err != nil {
		log.Printf("Failed to update team: %v", err)
		h.handleServiceError(w, err)
		return
	}

	h.writeJSON(w, map[string]string{"message": "Team updated successfully"}, http.StatusOK)
}
