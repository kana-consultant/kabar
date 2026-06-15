package product

import (
	"database/sql"
	productApp "seo-backend/internal/application/product"
	baseRoutes "seo-backend/internal/domain/base"
	"seo-backend/internal/infrastructure/db/repositories"
	rbacCache "seo-backend/internal/infrastructure/db/repositories/rbac"
	authmw "seo-backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	baseroute      baseRoutes.Route
	ProductHandler ProductHandler
	permCache      *rbacCache.PermissionCache
}

func NewRoute(db *sql.DB, chi chi.Router, permCache *rbacCache.PermissionCache) *Route {
	productRepo := repositories.NewProductRepository(db)
	adapterConfigRepo := repositories.NewAdapterConfigRepository(db)
	WorkflowDefinitionRepository := repositories.NewWorkflowDefinitionRepository(db)
	WorkflowNodeRepository := repositories.NewWorkflowNodeRepository(db)
	productService := productApp.NewProductService(db, productRepo, adapterConfigRepo, WorkflowDefinitionRepository, WorkflowNodeRepository)
	ProductHandler := NewProductHandler(productService)

	return &Route{
		baseroute: baseRoutes.Route{
			CHI: chi,
			DB:  db,
		},
		ProductHandler: *ProductHandler,
		permCache:      permCache,
	}
}

func (h *Route) SetupRoutes() chi.Router {
	r := h.baseroute.CHI
	c := h.permCache

	r.Route("/products", func(r chi.Router) {
		r.With(authmw.ProductView(c)).Get("/", h.ProductHandler.GetAll)
		r.With(authmw.ProductCreate(c)).Post("/", h.ProductHandler.Create)
		r.With(authmw.ProductView(c)).Get("/{id}", h.ProductHandler.GetByID)
		r.With(authmw.ProductEdit(c)).Put("/{id}", h.ProductHandler.Update)
		r.With(authmw.ProductDelete(c)).Delete("/{id}", h.ProductHandler.Delete)
		r.With(authmw.ProductEdit(c)).Post("/{id}/test", h.ProductHandler.UpdateConnectionStatus)
	})

	return r
}
