package workflow_node

import (
	"database/sql"

	workflowNodeService "seo-backend/internal/application/workflow_node"
	BaseRoutes "seo-backend/internal/domain/base"
	"seo-backend/internal/infrastructure/db/repositories"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	BaseRoutes          BaseRoutes.Route
	WorkflowNodeHandler *WorkflowNodeHandler
}

func NewRoute(db *sql.DB, chiRouter chi.Router) *Route {
	workflowNodeRepo := repositories.NewWorkflowNodeRepository(db)
	workflowNodeSvc := workflowNodeService.NewWorkflowNodeService(workflowNodeRepo)
	workflowNodeHandler := NewWorkflowNodeHandler(workflowNodeSvc)

	return &Route{
		BaseRoutes: BaseRoutes.Route{
			DB:  db,
			CHI: chiRouter,
		},
		WorkflowNodeHandler: workflowNodeHandler,
	}
}

func (r *Route) SetupRoutes() chi.Router {
	router := r.BaseRoutes.CHI

	router.Route("/workflows", func(route chi.Router) {
		route.Get("/{id}/nodes", r.WorkflowNodeHandler.GetByWorkflowID)
		route.Post("/{id}/nodes", r.WorkflowNodeHandler.Create)
		route.Put("/{id}/nodes/reorder", r.WorkflowNodeHandler.Reorder)
	})

	router.Route("/nodes", func(route chi.Router) {
		route.Get("/{id}", r.WorkflowNodeHandler.GetByID)
		route.Put("/{id}", r.WorkflowNodeHandler.Update)
		route.Delete("/{id}", r.WorkflowNodeHandler.Delete)
	})

	return router
}
