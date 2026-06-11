package workflow

import (
	"context"
)

// internal/service/workflow_definition_service.go (tambahan di atas)
type WorkflowDefinitionService interface {
	GetByID(ctx context.Context, id string) (WorkflowDefinition, error)
	GetByProductID(ctx context.Context, productID string) ([]WorkflowDefinition, error)
	Create(ctx context.Context, wf *WorkflowDefinition) error
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
}
