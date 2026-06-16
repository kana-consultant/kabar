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

func (s *WorkflowNodeService) SaveBatch(ctx context.Context, workflowID string, req workflow_node.BatchSaveRequest) (*workflow_node.BatchSaveResult, error) {
	result := &workflow_node.BatchSaveResult{
		Created: []workflow_node.BatchCreateResult{},
		Updated: []string{},
		Deleted: []string{},
		Errors:  []workflow_node.BatchError{},
	}

	// 1. DELETE nodes
	for _, id := range req.ToDelete {
		if err := s.repo.Delete(ctx, id); err != nil {
			result.Errors = append(result.Errors, workflow_node.BatchError{
				Operation: "delete",
				ID:        id,
				Message:   err.Error(),
			})
		} else {
			result.Deleted = append(result.Deleted, id)
		}
	}

	// 2. UPDATE nodes
	for _, update := range req.ToUpdate {
		if stepOrder, ok := update.Updates["stepOrder"]; ok {
			if step, ok := stepOrder.(float64); ok && step <= 0 {
				result.Errors = append(result.Errors, workflow_node.BatchError{
					Operation: "update",
					ID:        update.ID,
					Message:   "step order must be greater than 0",
				})
				continue
			}
		}

		if err := s.repo.Update(ctx, update.ID, update.Updates); err != nil {
			result.Errors = append(result.Errors, workflow_node.BatchError{
				Operation: "update",
				ID:        update.ID,
				Message:   err.Error(),
			})
		} else {
			result.Updated = append(result.Updated, update.ID)
		}
	}

	// 3. CREATE nodes
	if len(req.ToCreate) > 0 {
		nodes := make([]workflow_node.WorkflowNode, len(req.ToCreate))
		tempIDMap := make(map[int]string)

		for i, n := range req.ToCreate {
			tempIDMap[i] = n.TempID

			nodes[i] = workflow_node.WorkflowNode{
				WorkflowID:      workflowID,
				AdapterConfigID: n.AdapterConfigID,
				StepOrder:       n.StepOrder,
				NextNodeID:      n.NextNodeID,
			}
		}

		validNodes := []workflow_node.WorkflowNode{}
		validIndices := []int{}
		for i, node := range nodes {
			if node.AdapterConfigID != "" && node.StepOrder > 0 {
				validNodes = append(validNodes, node)
				validIndices = append(validIndices, i)
			}
		}

		if len(validNodes) > 0 {
			createdNodes, err := s.repo.InsertBatch(ctx, validNodes)
			if err != nil {
				for idx, node := range validNodes {
					if err := s.repo.Insert(ctx, &node); err != nil {
						result.Errors = append(result.Errors, workflow_node.BatchError{
							Operation: "create",
							ID:        tempIDMap[validIndices[idx]],
							Message:   err.Error(),
						})
					} else {
						result.Created = append(result.Created, workflow_node.BatchCreateResult{
							TempID: tempIDMap[validIndices[idx]],
							Node:   node,
						})
					}
				}
			} else {
				for i, node := range createdNodes {
					if i < len(validIndices) {
						result.Created = append(result.Created, workflow_node.BatchCreateResult{
							TempID: tempIDMap[validIndices[i]],
							Node:   node,
						})
					}
				}
			}
		}
	}

	// 4. Update nextNodeID references (from temp IDs to real IDs)
	if len(result.Created) > 0 {
		tempToReal := make(map[string]string)
		for _, created := range result.Created {
			tempToReal[created.TempID] = created.Node.ID
		}

		for _, created := range result.Created {
			if created.Node.NextNodeID != nil {
				nextID := *created.Node.NextNodeID
				if realID, ok := tempToReal[nextID]; ok {
					if err := s.repo.Update(ctx, created.Node.ID, map[string]interface{}{
						"next_node_id": realID,
					}); err != nil {
						result.Errors = append(result.Errors, workflow_node.BatchError{
							Operation: "update",
							ID:        created.Node.ID,
							Message:   fmt.Sprintf("failed to update nextNodeID: %v", err),
						})
					}
				}
			}
		}
	}

	return result, nil
}
