package draft

import (
	"database/sql"
	baseRoutes "seo-backend/internal/domain/base"
	"seo-backend/internal/domain/draft"
	rbacCache "seo-backend/internal/infrastructure/db/repositories/rbac"
	authmw "seo-backend/internal/middleware"
	"seo-backend/internal/scheduler"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
)

type Route struct {
	baseRoute      baseRoutes.Route
	DraftHandler   DraftHandler
	RedisScheduler *scheduler.RedisScheduler
	permCache      *rbacCache.PermissionCache // ← tambah ini
}

func NewRoute(
	db *sql.DB,
	chi chi.Router,
	redisClient *redis.Client,
	redisScheduler *scheduler.RedisScheduler,
	permCache *rbacCache.PermissionCache,
	DraftService draft.Service,
) *Route {

	DraftHandler := NewDraftHandler(DraftService)

	return &Route{
		baseRoute: baseRoutes.Route{
			DB:  db,
			CHI: chi,
		},
		DraftHandler:   *DraftHandler,
		RedisScheduler: redisScheduler,
		permCache:      permCache, // ← simpan di sini, bukan di baseRoute
	}
}

func (h *Route) SetupRoutes() chi.Router {
	r := h.baseRoute.CHI
	c := h.permCache // shortcut

	r.Route("/drafts", func(r chi.Router) {

		r.With(authmw.DraftView(c)).Get("/", h.DraftHandler.GetAll)
		r.With(authmw.DraftView(c)).Get("/scheduled", h.DraftHandler.GetAllScheduled)
		r.With(authmw.DraftCreate(c)).Post("/", h.DraftHandler.Create)

		r.With(authmw.DraftPublish(c)).Post("/publish", h.DraftHandler.PublishContent)
		r.With(authmw.DraftCreate(c)).Post("/schedule", h.DraftHandler.ScheduleDraft)
		r.With(authmw.DraftCreate(c)).Post("/schedule/cancel", h.DraftHandler.CancelScheduledDraft)

		r.With(authmw.DraftView(c)).Get("/{id}", h.DraftHandler.GetByID)
		r.With(authmw.DraftEdit(c)).Put("/{id}", h.DraftHandler.Update)
		r.With(authmw.DraftDelete(c)).Delete("/{id}", h.DraftHandler.Delete)
		r.With(authmw.DraftPublish(c)).Post("/{id}/publish", h.DraftHandler.Publish)
		r.With(authmw.DraftView(c)).Get("/{id}/seo-score", h.DraftHandler.GetSEOScore)
		r.With(authmw.DraftView(c)).Get("/{id}/check-similarity", h.DraftHandler.CheckSimilarity)
	})
	return r
}
