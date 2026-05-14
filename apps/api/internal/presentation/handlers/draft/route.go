package draft

import (
	"database/sql"
	DraftService "seo-backend/internal/application/draft"
	baseRoutes "seo-backend/internal/domain/base"
	"seo-backend/internal/helper"
	"seo-backend/internal/infrastructure/db/repositories"
	"seo-backend/internal/scheduler"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	baseRoute      baseRoutes.Route
	DraftHandler   DraftHandler
	RedisScheduler *scheduler.RedisScheduler
}

func NewRoute(db *sql.DB, chi chi.Router, redisScheduler *scheduler.RedisScheduler) *Route {
	postService := helper.NewPostService(db)
	productRepo := repositories.NewProductRepository(db)
	DraftRepo := repositories.NewDraftRepository(db)
	DraftService := DraftService.NewService(DraftRepo, redisScheduler, postService, productRepo)
	DraftHandler := NewDraftHandler(DraftService)

	return &Route{
		baseRoute: baseRoutes.Route{
			DB:  db,
			CHI: chi,
		},
		DraftHandler: *DraftHandler,
	}
}

func (h *Route) SetupRoutes() chi.Router {
	r := h.baseRoute.CHI
	r.Route("/drafts", func(r chi.Router) {

		r.Get("/", h.DraftHandler.GetAll)
		r.Get("/scheduled", h.DraftHandler.GetAllScheduled)
		r.Post("/", h.DraftHandler.Create)
		r.Get("/{id}", h.DraftHandler.GetByID)
		r.Put("/{id}", h.DraftHandler.Update)
		r.Delete("/{id}", h.DraftHandler.Delete)

		r.Post("/{id}/publish", h.DraftHandler.Publish)
		r.Post("/publish", h.DraftHandler.PublishContent)

		r.Post("/schedule", h.DraftHandler.ScheduleDraft)
		r.Post("/schedule/cancel", h.DraftHandler.CancelScheduledDraft)

		r.Get("/{id}/seo-score", h.DraftHandler.GetSEOScore)
		r.Get("/{id}/check-similarity", h.DraftHandler.CheckSimilarity)
	})
	return r
}
