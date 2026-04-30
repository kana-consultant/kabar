package team

import (
	"log"
	"net/http"
)

// RemoveMember removes a user from a team
func (h *TeamHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)

	team, err := h.teamService.RemoveMember(userCtx.GetTeamID(), userCtx.GetUserID(), userCtx)
	if err != nil {
		log.Printf("Failed to remove member: %v", err)
		if err.Error() == "member not found" {
			h.writeError(w, "Member not found", http.StatusNotFound)
		} else {
			h.writeError(w, "Failed to remove member", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSON(w, team, http.StatusOK)
}
