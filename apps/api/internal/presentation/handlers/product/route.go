package product

import (
	"database/sql"
	productApp "seo-backend/internal/application/product"
	baseRoutes "seo-backend/internal/domain/base"
	"seo-backend/internal/infrastructure/db/repositories"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	baseroute      baseRoutes.Route
	ProductHandler ProductHandler
}

func NewRoute(db *sql.DB, chi chi.Router) *Route {
	productRepo := repositories.NewProductRepository(db)
	adapterConfigRepo := repositories.NewAdapterConfigRepository(db)
	productService := productApp.NewProductService(db, productRepo, adapterConfigRepo)
	ProductHandler := NewProductHandler(productService)
	return &Route{
		baseroute: baseRoutes.Route{
			CHI: chi,
			DB:  db,
		},
		ProductHandler: *ProductHandler,
	}
}

func (h *Route) SetupRoutes() chi.Router {
	r := h.baseroute.CHI
	r.Route("/products", func(r chi.Router) {
		r.Get("/", h.ProductHandler.GetAll)
		r.Post("/", h.ProductHandler.Create)
		r.Get("/{id}", h.ProductHandler.GetByID)
		r.Put("/{id}", h.ProductHandler.Update)
		r.Delete("/{id}", h.ProductHandler.Delete)
		r.Post("/{id}/test", h.ProductHandler.UpdateConnectionStatus)
	})

	return r
}
