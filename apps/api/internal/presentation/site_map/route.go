// internal/domain/sitemap/route.go
package sitemap

import (
	"database/sql"
	sitemapService "seo-backend/internal/application/sitemap"
	baseRoutes "seo-backend/internal/domain/base"
	"seo-backend/internal/domain/history"
	rbacCache "seo-backend/internal/infrastructure/db/repositories/rbac"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	baseroute      baseRoutes.Route
	sitemapHandler SitemapHandler
	permCache      *rbacCache.PermissionCache
}

func NewSitemapRoute(
	db *sql.DB,
	chi chi.Router,
	permCache *rbacCache.PermissionCache,
	historyRepo history.HistoryRepository,
) *Route {
	// Initialize service
	service := sitemapService.NewSitemapService(historyRepo)

	// Initialize handler
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
		// Public endpoint - no auth required (sitemap should be accessible)
		// But we can add optional auth if needed
		router.Get("/", r.sitemapHandler.GenerateSitemap)

		// If you want to protect it with auth, uncomment below:
		// router.With(authmw.SitemapView(c)).Get("/", r.sitemapHandler.GenerateSitemap)
	})

	return router
}
