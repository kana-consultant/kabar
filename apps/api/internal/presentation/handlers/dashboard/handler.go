package dashboard

import (
	"encoding/json"
	"log"
	"net/http"
	"seo-backend/internal/domain/dashboard"
	"seo-backend/internal/helper"
)

type DashboardHandler struct {
	service dashboard.DashboardService
}

func NewDashboardHandler(service dashboard.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		service: service,
	}
}

func (h *DashboardHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := helper.GetUserContext(r)
	userContext := dashboard.DashboardFilter{
		UserID: userCtx.GetUserID(),
		TeamID: userCtx.GetTeamID(),
		Role:   userCtx.GetRole(),
	}
	stats, _ := h.service.GetStats(ctx, userContext)

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
