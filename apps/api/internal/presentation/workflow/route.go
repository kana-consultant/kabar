package workflow

import (
	"database/sql"
	service "seo-backend/internal/application/workflow"
	BaseRoutes "seo-backend/internal/domain/base"
	"seo-backend/internal/infrastructure/db/repositories"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	BaseRoutes                BaseRoutes.Route
	WorkflowDefinitionHandler *WorkflowDefinitionHandler
}

func NewRoute(db *sql.DB, chiRouter chi.Router) *Route {
	workflowRepo := repositories.NewWorkflowDefinitionRepository(db)
	workflowService := service.NewWorkflowDefinitionService(workflowRepo)
	workflowHandler := NewWorkflowDefinitionHandler(workflowService)

	return &Route{
		BaseRoutes: BaseRoutes.Route{
			DB:  db,
			CHI: chiRouter,
		},
		WorkflowDefinitionHandler: workflowHandler,
	}
}

func (r *Route) SetupRoutes() chi.Router {
	router := r.BaseRoutes.CHI
	router.Route("/workflows", func(route chi.Router) {
		route.Post("/", r.WorkflowDefinitionHandler.Create)
		route.Get("/{id}", r.WorkflowDefinitionHandler.GetByID)
		route.Put("/{id}", r.WorkflowDefinitionHandler.Update)
		route.Delete("/{id}", r.WorkflowDefinitionHandler.Delete)
	})

	router.Route("/products/{productId}/workflows", func(route chi.Router) {
		route.Get("/", r.WorkflowDefinitionHandler.GetByProductID)
	})

	return router
}
