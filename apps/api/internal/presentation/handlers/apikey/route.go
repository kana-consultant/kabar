package apikey

import (
	"database/sql"
	BaseRoutes "seo-backend/internal/domain/base"
	"seo-backend/internal/infrastructure/db/repositories"
	"seo-backend/internal/pkg/crypto"

	apiKeyService "seo-backend/internal/application/apikey"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	baseRoutes    BaseRoutes.Route
	APIKeyHandler APIKeyHandler
}

func NewRoute(db *sql.DB, chi chi.Router) *Route {
	repoApi := repositories.NewAPIKeyRepository(db)
	encrypt := crypto.NewAESEncryptor()
	apiKeyService := apiKeyService.NewService(repoApi, encrypt)
	APIKeyHandler := NewAPIKeyHandler(apiKeyService)
	return &Route{
		baseRoutes: BaseRoutes.Route{
			CHI: chi,
			DB:  db,
		},
		APIKeyHandler: *APIKeyHandler,
	}
}

func (h *Route) SetupRoutes() chi.Router {
	r := h.baseRoutes.CHI
	r.Route("/api-keys", func(r chi.Router) {
		r.Get("/", h.APIKeyHandler.GetAll)
		r.Post("/", h.APIKeyHandler.Create)
		r.Get("/{id}", h.APIKeyHandler.GetByID)
		r.Put("/{id}", h.APIKeyHandler.Update)
		r.Delete("/{id}", h.APIKeyHandler.Delete)
	})
	return r
}
