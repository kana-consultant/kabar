// internal/presentation/workflow/route.go
package workflow

import (
	"database/sql"

	workflowService "seo-backend/internal/application/workflow"
	workflowNodeService "seo-backend/internal/application/workflow_node"
	BaseRoutes "seo-backend/internal/domain/base"
	"seo-backend/internal/infrastructure/db/repositories"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	BaseRoutes                BaseRoutes.Route
	WorkflowDefinitionHandler *WorkflowDefinitionHandler
	WorkflowNodeHandler       *WorkflowNodeHandler
}

func NewRoute(db *sql.DB, chiRouter chi.Router) *Route {
	workflowDefRepo := repositories.NewWorkflowDefinitionRepository(db)
	workflowDefSvc := workflowService.NewWorkflowDefinitionService(workflowDefRepo)
	workflowDefHandler := NewWorkflowDefinitionHandler(workflowDefSvc)

	workflowNodeRepo := repositories.NewWorkflowNodeRepository(db)
	workflowNodeSvc := workflowNodeService.NewWorkflowNodeService(workflowNodeRepo)
	workflowNodeHandler := NewWorkflowNodeHandler(workflowNodeSvc)

	return &Route{
		BaseRoutes: BaseRoutes.Route{
			DB:  db,
			CHI: chiRouter,
		},
		WorkflowDefinitionHandler: workflowDefHandler,
		WorkflowNodeHandler:       workflowNodeHandler,
	}
}

func (r *Route) SetupRoutes() chi.Router {
	router := r.BaseRoutes.CHI

	router.Route("/workflows", func(route chi.Router) {
		// Workflow Definitions
		route.Post("/", r.WorkflowDefinitionHandler.Create)
		route.Get("/{id}", r.WorkflowDefinitionHandler.GetByID)
		route.Put("/{id}", r.WorkflowDefinitionHandler.Update)
		route.Delete("/{id}", r.WorkflowDefinitionHandler.Delete)

		// Workflow Nodes (nested under workflow)
		route.Get("/{id}/nodes", r.WorkflowNodeHandler.GetByWorkflowID)
		route.Post("/{id}/nodes", r.WorkflowNodeHandler.SaveBatch)
		route.Put("/{id}/nodes/reorder", r.WorkflowNodeHandler.Reorder)
	})

	router.Route("/products/{productId}/workflows", func(route chi.Router) {
		route.Get("/", r.WorkflowDefinitionHandler.GetByProductID)
	})

	router.Route("/nodes", func(route chi.Router) {
		route.Get("/{id}", r.WorkflowNodeHandler.GetByID)
		route.Put("/{id}", r.WorkflowNodeHandler.Update)
		route.Delete("/{id}", r.WorkflowNodeHandler.Delete)
	})

	return router
}
