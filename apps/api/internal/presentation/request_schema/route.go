package request_schema

import (
	"database/sql"
	requestSchemaService "seo-backend/internal/application/request_schema"

	"seo-backend/internal/infrastructure/db/repositories"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
)

type RequestSchemaRoute struct {
	db          *sql.DB
	chi         chi.Router
	redisClient *redis.Client
	handler     *RequestSchemaHandler
}

func NewRoute(db *sql.DB, chi chi.Router, redisClient *redis.Client) *RequestSchemaRoute {
	// Repository
	repo := repositories.NewRequestSchemaRepository(db, redisClient)

	// Service
	service := requestSchemaService.NewService(repo)

	// Handler
	handler := NewRequestSchemaHandler(service)

	return &RequestSchemaRoute{
		db:          db,
		chi:         chi,
		redisClient: redisClient,
		handler:     handler,
	}
}

func (r *RequestSchemaRoute) SetupRoutes() chi.Router {
	router := r.chi

	router.Route("/schemas", func(router chi.Router) {
		router.Post("/", r.handler.Create)
		router.Get("/", r.handler.GetAll)
		router.Get("/{id}", r.handler.GetByID)
		router.Put("/{id}", r.handler.Update)
		router.Delete("/{id}", r.handler.Delete)
	})

	// Get by provider
	router.Route("/providers/{provider_id}/schemas", func(router chi.Router) {
		router.Get("/", r.handler.GetByProvider)
	})

	return router
}
