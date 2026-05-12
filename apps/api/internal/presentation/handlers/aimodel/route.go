package aimodel

import (
	"database/sql"
	ModelApp "seo-backend/internal/application/model"
	baseRoutes "seo-backend/internal/domain/base"
	"seo-backend/internal/infrastructure/db/repositories"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	baseroute      baseRoutes.Route
	AIModelHandler AIModelHandler
}

func NewRoute(db *sql.DB, chi *chi.Mux) *Route {
	AiModelRepo := repositories.ModelRepository(db)
	AiModelService := ModelApp.NewService(AiModelRepo)
	AIModelHandler := NewAIModelHandler(AiModelService)
	return &Route{
		baseroute: baseRoutes.Route{
			DB:  db,
			CHI: chi,
		},
		AIModelHandler: *AIModelHandler,
	}
}

func (h *Route) SetupRoute() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/api/models", func(r chi.Router) {
		r.Get("/", h.AIModelHandler.GetAll)
		r.Get("/with-status", h.AIModelHandler.GetAllWithStatus)
		r.Get("/{id}", h.AIModelHandler.GetByID)
		r.Get("/provider/{providerId}", h.AIModelHandler.GetByProvider)
	})
	return r
}
