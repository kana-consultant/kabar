package team

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"seo-backend/internal/middleware/auth"
	"seo-backend/internal/models"
	services "seo-backend/internal/service/team"
)

type TeamHandler struct {
	teamService *services.Service
}

func NewTeamHandler() *TeamHandler {
	return &TeamHandler{
		teamService: services.NewService(),
	}
}

// Helper: Get user context from request
func (h *TeamHandler) getUserContext(r *http.Request) models.UserContext {
	ctx := r.Context()
	return &models.SimpleUserContext{
		UserID: auth.GetUserID(ctx),
		TeamID: auth.GetTeamID(ctx),
		Role:   auth.GetUserRole(ctx),
	}
}

// Helper: Write JSON response
func (h *TeamHandler) writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

// Helper: Write error response
func (h *TeamHandler) writeError(w http.ResponseWriter, message string, status int) {
	h.writeJSON(w, map[string]string{"error": message}, status)
}

// Helper: Handle service errors
func (h *TeamHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch err.Error() {
	case "access denied":
		h.writeError(w, "Forbidden", http.StatusForbidden)
	case "team not found":
		h.writeError(w, "Team not found", http.StatusNotFound)
	case "cannot delete team with active members":
		h.writeError(w, err.Error(), http.StatusBadRequest)
	case "member already in team":
		h.writeError(w, err.Error(), http.StatusConflict)
	case "member not found":
		h.writeError(w, "Member not found", http.StatusNotFound)
	default:
		if strings.Contains(err.Error(), "maximum member limit") {
			h.writeError(w, err.Error(), http.StatusBadRequest)
		} else {
			log.Printf("Unexpected error: %v", err)
			h.writeError(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func validateTeamRequest(req models.CreateTeamRequest) error {
	if req.Name == "" {
		return fmt.Errorf("team name is required")
	}
	if len(req.Name) > 100 {
		return fmt.Errorf("team name too long (max 100 characters)")
	}
	return nil
}
