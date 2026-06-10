package provider

import (
	"database/sql"
	providerService "seo-backend/internal/application/provider"
	BaseRoutes "seo-backend/internal/domain/base"
	"seo-backend/internal/infrastructure/db/repositories"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
)

type Route struct {
	baseroute       BaseRoutes.Route
	ProviderHandler ProviderHandler
}

func NewRoute(db *sql.DB, chi chi.Router, redisClient *redis.Client) *Route {
	providerRepo := repositories.NewProviderRepository(db, redisClient)
	familyModelRepo := repositories.NewModelFamiliesRepository(db, redisClient)
	SchemaRepo := repositories.NewRequestSchemaRepository(db, redisClient)
	providerService := providerService.NewService(db, providerRepo, familyModelRepo, SchemaRepo, redisClient)
	ProviderHandler := NewProviderHandler(providerService)
	return &Route{
		baseroute: BaseRoutes.Route{
			DB:  db,
			CHI: chi,
		},
		ProviderHandler: *ProviderHandler,
	}
}

func (h *Route) SetupRoutes() chi.Router {
	r := h.baseroute.CHI
	r.Route("/providers", func(r chi.Router) {
		r.Get("/", h.ProviderHandler.GetAll)
		r.Post("/", h.ProviderHandler.Create)
		r.Get("/{id}", h.ProviderHandler.GetByID)
		r.Put("/{id}", h.ProviderHandler.Update)
		r.Delete("/{id}", h.ProviderHandler.Delete)
	})
	return r
}
