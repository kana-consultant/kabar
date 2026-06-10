package modelfamily

import (
	"database/sql"

	model_family "seo-backend/internal/application/model_families"
	"seo-backend/internal/infrastructure/db/repositories"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
)

type Route struct {
	db          *sql.DB
	chi         chi.Router
	redisClient *redis.Client
	handler     *ModelFamilyHandler
}

func NewRoute(db *sql.DB, chi chi.Router, redisClient *redis.Client) *Route {
	// Init repository
	repoFamilies := repositories.NewModelFamiliesRepository(db, redisClient)
	service := model_family.NewService(repoFamilies)

	// Init handler
	handler := NewModelFamilyHandler(service)

	return &Route{
		db:          db,
		chi:         chi,
		redisClient: redisClient,
		handler:     handler,
	}
}

func (r *Route) SetupRoute() chi.Router {
	router := r.chi

	router.Route("/families", func(router chi.Router) {
		// Create route
		router.Post("/", r.handler.Create)

		// GET routes
		router.Get("/", r.handler.GetAll)
		router.Get("/exists", r.handler.CheckExists)
		router.Get("/{id}", r.handler.GetByID)

		// Update & Delete
		router.Put("/{id}", r.handler.Update)
		router.Delete("/{id}", r.handler.Delete)
	})

	// Additional routes for provider and schema
	router.Route("/families/providers", func(router chi.Router) {
		router.Get("/{provider_id}/families", r.handler.GetByProvider)
	})

	router.Route("/families/schemas", func(router chi.Router) {
		router.Get("/{schema_id}/families", r.handler.GetBySchema)
	})

	return router
}
