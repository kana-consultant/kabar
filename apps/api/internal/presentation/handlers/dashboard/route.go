package dashboard

import (
	"database/sql"
	dashboardApp "seo-backend/internal/application/dashboard"
	baseRoutes "seo-backend/internal/domain/base"
	"seo-backend/internal/infrastructure/db/repositories"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	baseroute        baseRoutes.Route
	DashboardHandler DashboardHandler
}

func NewRoute(db *sql.DB, chi chi.Router) *Route {
	dashboardRepo := repositories.NewDashboardRepository(db)
	dashboardService := dashboardApp.NewDashboardService(dashboardRepo)
	dashboardHandler := NewDashboardHandler(dashboardService)
	return &Route{
		baseroute: baseRoutes.Route{
			CHI: chi,
			DB:  db,
		},
		DashboardHandler: *dashboardHandler,
	}
}

func (h *Route) SetupRoutes() chi.Router {
	r := h.baseroute.CHI
	r.Get("/dashboard/stats", h.DashboardHandler.GetStats)
	return r
}
