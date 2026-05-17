package history

import (
	"database/sql"
	historyService "seo-backend/internal/application/history"
	baseRoutes "seo-backend/internal/domain/base"
	historyBuilder "seo-backend/internal/infrastructure/db/query_builder"
	"seo-backend/internal/infrastructure/db/repositories"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	baseroute      baseRoutes.Route
	HistoryHandler HistoryHandler
}

func NewHistoryRoute(db *sql.DB, chi chi.Router) *Route {
	repoHistory := repositories.NewHistoryRepository(db)
	qb_history := historyBuilder.NewQueryBuilder()
	HistoryRepo := historyService.NewService(repoHistory, *qb_history)
	historyHandler := NewHistoryHandler(HistoryRepo)
	return &Route{
		baseroute: baseRoutes.Route{
			DB:  db,
			CHI: chi,
		},
		HistoryHandler: *historyHandler,
	}
}

func (h *Route) SetupRoute() chi.Router {
	r := h.baseroute.CHI
	r.Route("/history", func(r chi.Router) {
		r.Get("/", h.HistoryHandler.GetAll)
		r.Post("/", h.HistoryHandler.Create)
		r.Get("/{id}", h.HistoryHandler.GetByID)
		r.Put("/{id}", h.HistoryHandler.Update)
		r.Delete("/{id}", h.HistoryHandler.Delete)
	})

	return r
}
