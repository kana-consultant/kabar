package dashboard

import (
	"encoding/json"
	"log"
	"net/http"
	"seo-backend/internal/domain/dashboard"
	auth "seo-backend/internal/middleware"
)

type DashboardHandler struct {
	service dashboard.DashboardService
}

func NewDashboardHandler(service dashboard.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		service: service,
	}
}

// GetStats godoc
// @Summary Get dashboard statistics
// @Description Get comprehensive dashboard statistics for authenticated user/team
// @Tags Dashboard
// @Accept json
// @Produce json
// @Success 200 {object} dashboard.DashboardStats "Dashboard statistics"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/dashboard/stats [get]
func (h *DashboardHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := auth.GetUserContext(r)
	userContext := dashboard.DashboardFilter{
		UserID: userCtx.GetUserID(),
		TeamID: userCtx.GetTeamID(),
		Role:   userCtx.GetRole(),
	}
	stats, err := h.service.GetStats(ctx, userContext)
	if err != nil {
		log.Printf("Failed to get dashboard stats: %v", err)
		h.writeError(w, "Failed to retrieve dashboard statistics", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, stats, http.StatusOK)
}

// Helper: Write JSON response
func (h *DashboardHandler) writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

// Helper: Write error response
func (h *DashboardHandler) writeError(w http.ResponseWriter, message string, status int) {
	h.writeJSON(w, map[string]string{"error": message}, status)
}
