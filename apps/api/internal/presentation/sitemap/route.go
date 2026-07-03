// internal/domain/sitemap/route.go
package sitemap

import (
	"database/sql"
	sitemapService "seo-backend/internal/application/sitemap"
	baseRoutes "seo-backend/internal/domain/base"
	"seo-backend/internal/domain/history"
	"seo-backend/internal/domain/product"
	rbacCache "seo-backend/internal/infrastructure/db/repositories/rbac"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	baseroute      baseRoutes.Route
	sitemapHandler SitemapHandler
	permCache      *rbacCache.PermissionCache
}

func NewRoute(
	db *sql.DB,
	chi chi.Router,
	permCache *rbacCache.PermissionCache,
	historyRepo history.HistoryRepository,
	productRepo product.ProductRepository,
) *Route {
	service := sitemapService.NewService(historyRepo, productRepo)
	sitemapHandler := NewSitemapHandler(service)

	return &Route{
		baseroute: baseRoutes.Route{
			DB:  db,
			CHI: chi,
		},
		sitemapHandler: *sitemapHandler,
		permCache:      permCache,
	}
}

func (r *Route) SetupRoute() chi.Router {
	router := r.baseroute.CHI

	router.Route("/sitemap", func(router chi.Router) {
		// Public endpoint - generate sitemap
		router.Get("/", r.sitemapHandler.GenerateSitemap)

		// Protected endpoint - get sitemap history
		router.Get("/history", r.sitemapHandler.GetSitemapHistory)
	})

	return router
}
