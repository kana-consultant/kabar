// internal/service/workflow_definition_service.go
package service

import (
	"context"
	"fmt"

	"seo-backend/internal/domain/workflow"
)

type WorkflowDefinitionService struct {
	repo workflow.WorkflowDefinitionRepository
}

func NewWorkflowDefinitionService(repo workflow.WorkflowDefinitionRepository) workflow.WorkflowDefinitionService {
	return &WorkflowDefinitionService{repo: repo}
}

func (s *WorkflowDefinitionService) GetByID(ctx context.Context, id string) (workflow.WorkflowDefinition, error) {
	wf, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return workflow.WorkflowDefinition{}, fmt.Errorf("failed to get workflow definition: %w", err)
	}
	if wf == nil {
		return workflow.WorkflowDefinition{}, fmt.Errorf("workflow definition not found")
	}
	return *wf, nil
}

func (s *WorkflowDefinitionService) GetByProductID(ctx context.Context, productID string) ([]workflow.WorkflowDefinition, error) {
	workflows, err := s.repo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow definitions: %w", err)
	}

	result := make([]workflow.WorkflowDefinition, 0, len(workflows))
	for _, wf := range workflows {
		result = append(result, *wf)
	}
	return result, nil
}

func (s *WorkflowDefinitionService) Create(ctx context.Context, wf *workflow.WorkflowDefinition) error {
	if wf.ProductID == "" {
		return fmt.Errorf("product ID is required")
	}
	if wf.Name == "" {
		return fmt.Errorf("name is required")
	}
	return s.repo.Insert(ctx, wf)
}

func (s *WorkflowDefinitionService) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}
	return s.repo.Update(ctx, id, updates)
}

func (s *WorkflowDefinitionService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
