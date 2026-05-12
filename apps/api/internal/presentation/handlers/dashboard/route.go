package dashboard

import (
	dashboardHandler "seo-backend/internal/presentation/handlers/dashboard"
	dashboardApp "seo-backend/internal/application/dashboard"
	"seo-backend/internal/infrastructure/db/repositories"
)
type Route struct{
	baseroute baseRoute.Route
	DashboardHandler DashboardHandler
}

func NewRoute(db *sql.DB, chi *chi.Mux) *Route {
	dashboardRepo := repositories.NewDashboardRepository(db)
	dashboardApp.NewDashboardService(dashboardRepo)
	dashboardHandler.NewDashboardHandler(dashboardService)
	return &Route{
		baseroute: baseRoutes.Route{
			CHI: chi,
			DB:  db,
		},
		DashboardHandler: *dashboardHandler,
	}
}

func (h *Route) SetupRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/api/dashboard/stats", container.DashboardHandler.GetStats)
	return r
}
