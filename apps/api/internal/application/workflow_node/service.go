// internal/application/workflow_node/workflow_node_service.go
package workflow_node

import (
	"context"

	"fmt"
	"seo-backend/internal/domain/workflow_node"
)

type WorkflowNodeService struct {
	repo workflow_node.WorkflowNodeRepository
}

func NewWorkflowNodeService(repo workflow_node.WorkflowNodeRepository) workflow_node.WorkflowNodeService {
	return &WorkflowNodeService{repo: repo}
}

func (s *WorkflowNodeService) GetByID(ctx context.Context, id string) (*workflow_node.WorkflowNode, error) {
	node, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow node: %w", err)
	}
	if node == nil {
		return nil, fmt.Errorf("workflow node not found")
	}
	return node, nil
}

func (s *WorkflowNodeService) GetByWorkflowID(ctx context.Context, workflowID string) ([]*workflow_node.WorkflowNode, error) {
	nodes, err := s.repo.GetByWorkflowID(ctx, workflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow nodes: %w", err)
	}
	return nodes, nil
}

func (s *WorkflowNodeService) Create(ctx context.Context, node *workflow_node.WorkflowNode) error {
	if node.WorkflowID == "" {
		return fmt.Errorf("workflow ID is required")
	}
	if node.AdapterConfigID == "" {
		return fmt.Errorf("adapter config ID is required")
	}
	if node.StepOrder <= 0 {
		return fmt.Errorf("step order must be greater than 0")
	}

	return s.repo.Insert(ctx, node)
}

func (s *WorkflowNodeService) CreateBatch(ctx context.Context, nodes []workflow_node.WorkflowNode) ([]workflow_node.WorkflowNode, error) {
	if len(nodes) == 0 {
		return nil, fmt.Errorf("no nodes to create")
	}

	for i := range nodes {
		if nodes[i].WorkflowID == "" {
			return nil, fmt.Errorf("workflow ID is required for node %d", i)
		}
		if nodes[i].AdapterConfigID == "" {
			return nil, fmt.Errorf("adapter config ID is required for node %d", i)
		}
		if nodes[i].StepOrder <= 0 {
			return nil, fmt.Errorf("step order must be greater than 0 for node %d", i)
		}
	}

	return s.repo.InsertBatch(ctx, nodes)
}

func (s *WorkflowNodeService) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}

	if stepOrder, ok := updates["stepOrder"]; ok {
		if step, ok := stepOrder.(float64); ok && step <= 0 {
			return fmt.Errorf("step order must be greater than 0")
		}
	}

	return s.repo.Update(ctx, id, updates)
}

func (s *WorkflowNodeService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *WorkflowNodeService) ReorderNodes(ctx context.Context, workflowID string, nodeIDs []string) error {
	if len(nodeIDs) == 0 {
		return fmt.Errorf("node IDs cannot be empty")
	}
	if workflowID == "" {
		return fmt.Errorf("workflow ID is required")
	}

	return s.repo.ReorderNodes(ctx, workflowID, nodeIDs)
}
