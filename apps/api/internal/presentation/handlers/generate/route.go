package generate

import (
	"database/sql"
	generateService "seo-backend/internal/application/generate"
	BaseRoutes "seo-backend/internal/domain/base"
	aiBuilder "seo-backend/internal/infrastructure/ai/builder"
	aiParser "seo-backend/internal/infrastructure/ai/parser"
	"seo-backend/internal/infrastructure/db/repositories"
	"seo-backend/internal/infrastructure/http/client"

	"time"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	baseroute       BaseRoutes.Route
	GenerateHandler GenerateHandler
}

func NewRoute(db *sql.DB, chi chi.Router) *Route {

	promptBuilder := aiBuilder.NewPromptBuilder()
	requestBuilder := aiBuilder.NewRequestBuilder()
	responseParser := aiParser.NewResponseParser()
	generateRepo := repositories.NewRepository(db)
	var timeOut time.Duration
	timeOut = 30
	client := client.NewHTTPClient(timeOut)
	generateService := generateService.NewService(generateRepo, client, promptBuilder, requestBuilder, responseParser)
	GenerateHandler := NewGenerateHandler(generateService)

	return &Route{
		baseroute: BaseRoutes.Route{
			DB:  db,
			CHI: chi,
		},
		GenerateHandler: *GenerateHandler,
	}
}

func (h *Route) SetupRoutes() chi.Router {
	r := h.baseroute.CHI
	r.Route("/generate", func(r chi.Router) {
		r.Post("/article", h.GenerateHandler.GenerateArticle)
		r.Post("/image", h.GenerateHandler.GenerateImage)
	})
	return r

}
