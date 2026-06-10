package aimodel

import (
	"database/sql"
	ModelApp "seo-backend/internal/application/aimodel"
	baseRoutes "seo-backend/internal/domain/base"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
)

type Route struct {
	baseroute      baseRoutes.Route
	AIModelHandler AIModelHandler
}

func NewRoute(db *sql.DB, chi chi.Router, redisClient *redis.Client) *Route {
	AiModelService := ModelApp.NewService(db, redisClient)
	AIModelHandler := NewAIModelHandler(AiModelService)
	return &Route{
		baseroute: baseRoutes.Route{
			DB:  db,
			CHI: chi,
		},
		AIModelHandler: *AIModelHandler,
	}
}

func (h *Route) SetupRoute() chi.Router {
	r := h.baseroute.CHI
	r.Route("/models", func(r chi.Router) {
		// GET routes
		r.Get("/", h.AIModelHandler.GetAll)
		r.Get("/with-status", h.AIModelHandler.GetAllWithStatus)
		r.Get("/default", h.AIModelHandler.GetDefault)
		r.Get("/{id}", h.AIModelHandler.GetByID)
		r.Get("/{id}/schema", h.AIModelHandler.GetSchemaByID)
		r.Get("/provider/{providerId}", h.AIModelHandler.GetByProvider)
		r.Get("/family/{familyId}", h.AIModelHandler.GetByFamily)
		r.Get("/team/{teamId}", h.AIModelHandler.GetByTeam)

		// POST routes
		r.Post("/", h.AIModelHandler.Create)

		// PUT routes
		r.Put("/{id}", h.AIModelHandler.Update)
		r.Put("/{id}/default", h.AIModelHandler.SetAsDefault)

		// DELETE routes
		r.Delete("/{id}", h.AIModelHandler.Delete)
	})
	return r
}
