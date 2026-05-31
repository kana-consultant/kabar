package generate

import (
	"database/sql"
	"log"
	generateService "seo-backend/internal/application/generate"
	"seo-backend/internal/config"
	BaseRoutes "seo-backend/internal/domain/base"
	aiBuilder "seo-backend/internal/infrastructure/ai/builder"
	aiParser "seo-backend/internal/infrastructure/ai/parser"
	"seo-backend/internal/infrastructure/db/repositories"
	"seo-backend/internal/infrastructure/http/client"
	"seo-backend/internal/infrastructure/http/minio"
	"time"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	baseroute       BaseRoutes.Route
	GenerateHandler GenerateHandler
}

func NewRoute(db *sql.DB, chi chi.Router, cfg *config.Config) *Route {

	promptBuilder := aiBuilder.NewPromptBuilder()
	requestBuilder := aiBuilder.NewRequestBuilder()
	responseParser := aiParser.NewResponseParser()
	generateRepo := repositories.NewRepository(db)
	var timeOut time.Duration
	timeOut = 30
	client := client.NewHTTPClient(timeOut)
	minioStorage, err := minio.NewMinioService(
		cfg.MinioEndpoint,
		cfg.MinioAccessKey,
		cfg.MinioSecretKey,
		cfg.MinioBucket,
	)

	if err != nil {
		log.Fatal(err)
	}

	generateService := generateService.NewService(generateRepo, client, minioStorage, promptBuilder, requestBuilder, responseParser)
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
