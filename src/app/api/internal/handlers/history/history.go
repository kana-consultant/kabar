package history

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"seo-backend/internal/middleware/auth"
	"seo-backend/internal/models"
	services "seo-backend/internal/service/history"
)

type HistoryHandler struct {
	historyService *services.Service
}

func NewHistoryHandler() *HistoryHandler {
	return &HistoryHandler{
		historyService: services.NewService(),
	}
}

// Helper: Get user context from request
func (h *HistoryHandler) getUserContext(r *http.Request) models.UserContext {
	ctx := r.Context()
	return &models.SimpleUserContext{
		UserID: auth.GetUserID(ctx),
		TeamID: auth.GetTeamID(ctx),
		Role:   auth.GetUserRole(ctx),
	}
}

// Helper: Write JSON response
func (h *HistoryHandler) writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

// Helper: Write error response
func (h *HistoryHandler) writeError(w http.ResponseWriter, message string, status int) {
	h.writeJSON(w, map[string]string{"error": message}, status)
}

// Helper: Handle service errors
func (h *HistoryHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch err.Error() {
	case "access denied":
		h.writeError(w, "Forbidden", http.StatusForbidden)
	case "history not found":
		h.writeError(w, "History not found", http.StatusNotFound)
	case "admin access required":
		h.writeError(w, "Forbidden", http.StatusForbidden)
	default:
		log.Printf("Unexpected error: %v", err)
		h.writeError(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Helper: Parse pagination from query params
func parsePagination(r *http.Request) (limit, offset int) {
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}
	return
}
