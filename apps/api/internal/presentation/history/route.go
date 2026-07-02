package history

import (
	"database/sql"
	historyService "seo-backend/internal/application/history"
	baseRoutes "seo-backend/internal/domain/base"
	"seo-backend/internal/domain/history"
	historyBuilder "seo-backend/internal/infrastructure/db/query_builder"
	rbacCache "seo-backend/internal/infrastructure/db/repositories/rbac"
	authmw "seo-backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	baseroute      baseRoutes.Route
	HistoryHandler HistoryHandler
	permCache      *rbacCache.PermissionCache
}

func NewHistoryRoute(db *sql.DB, chi chi.Router, permCache *rbacCache.PermissionCache, historyRepo history.HistoryRepository) *Route {
	qb_history := historyBuilder.NewQueryBuilder()
	HistoryService := historyService.NewService(historyRepo, *qb_history)
	historyHandler := NewHistoryHandler(HistoryService)

	return &Route{
		baseroute: baseRoutes.Route{
			DB:  db,
			CHI: chi,
		},
		HistoryHandler: *historyHandler,
		permCache:      permCache,
	}
}

func (h *Route) SetupRoute() chi.Router {
	r := h.baseroute.CHI
	c := h.permCache

	r.Route("/history", func(r chi.Router) {
		r.With(authmw.HistoryView(c)).Get("/", h.HistoryHandler.GetAll)
		r.With(authmw.HistoryView(c)).Get("/recently", h.HistoryHandler.GetRecentActivity)
		r.With(authmw.HistoryView(c)).Get("/{id}", h.HistoryHandler.GetByID)
		r.With(authmw.HistoryDeleteGlobal(c)).Delete("/{id}", h.HistoryHandler.Delete)

		// endpoint ini tidak ada di schema permission history,
		// tapi tetap dijaga minimal dengan HistoryView
		r.With(authmw.HistoryView(c)).Put("/{id}", h.HistoryHandler.Update)
	})

	return r
}
