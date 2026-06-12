// internal/infrastructure/db/repositories/workflow_node_repository.go
package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"seo-backend/internal/domain/workflow_node"
	"strings"

	"github.com/lib/pq"
)

type WorkflowNodeRepository struct {
	db *sql.DB
}

func NewWorkflowNodeRepository(db *sql.DB) workflow_node.WorkflowNodeRepository {
	return &WorkflowNodeRepository{db: db}
}

func (r *WorkflowNodeRepository) GetByID(ctx context.Context, id string) (*workflow_node.WorkflowNode, error) {
	query := `
		SELECT id, workflow_id, adapter_config_id, step_order, 
			input_mapping, next_node_id, created_at
		FROM workflow_nodes WHERE id = $1
	`

	var node workflow_node.WorkflowNode
	var inputMappingJSON []byte
	var nextNodeID sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&node.ID, &node.WorkflowID, &node.AdapterConfigID,
		&node.StepOrder, &inputMappingJSON, &nextNodeID, &node.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get workflow node: %w", err)
	}

	if len(inputMappingJSON) > 0 {
		json.Unmarshal(inputMappingJSON, &node.InputMapping)
	}

	if nextNodeID.Valid {
		node.NextNodeID = &nextNodeID.String
	}

	return &node, nil
}

func (r *WorkflowNodeRepository) GetByWorkflowID(ctx context.Context, workflowID string) ([]*workflow_node.WorkflowNode, error) {
	query := `
		SELECT id, workflow_id, adapter_config_id, step_order, 
			input_mapping, next_node_id, created_at
		FROM workflow_nodes 
		WHERE workflow_id = $1
		ORDER BY step_order ASC
	`

	rows, err := r.db.QueryContext(ctx, query, workflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow nodes: %w", err)
	}
	defer rows.Close()

	var nodes []*workflow_node.WorkflowNode
	for rows.Next() {
		var node workflow_node.WorkflowNode
		var inputMappingJSON []byte
		var nextNodeID sql.NullString

		err := rows.Scan(
			&node.ID, &node.WorkflowID, &node.AdapterConfigID,
			&node.StepOrder, &inputMappingJSON, &nextNodeID, &node.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workflow node: %w", err)
		}

		if len(inputMappingJSON) > 0 {
			json.Unmarshal(inputMappingJSON, &node.InputMapping)
		}

		if nextNodeID.Valid {
			node.NextNodeID = &nextNodeID.String
		}

		nodes = append(nodes, &node)
	}

	return nodes, nil
}

func (r *WorkflowNodeRepository) Insert(ctx context.Context, node *workflow_node.WorkflowNode) error {
	inputMappingJSON, _ := json.Marshal(node.InputMapping)

	query := `
		INSERT INTO workflow_nodes (workflow_id, adapter_config_id, step_order, 
			input_mapping, next_node_id, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, created_at
	`

	var nextNodeID interface{}
	if node.NextNodeID != nil {
		nextNodeID = *node.NextNodeID
	}

	err := r.db.QueryRowContext(ctx, query,
		node.WorkflowID, node.AdapterConfigID, node.StepOrder,
		inputMappingJSON, nextNodeID,
	).Scan(&node.ID, &node.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert workflow node: %w", err)
	}

	return nil
}

func (r *WorkflowNodeRepository) InsertBatch(ctx context.Context, nodes []workflow_node.WorkflowNode) ([]workflow_node.WorkflowNode, error) {
	query := `
		WITH data AS (
			SELECT 
				unnest($1::uuid[]) as workflow_id,
				unnest($2::uuid[]) as adapter_config_id,
				unnest($3::int[]) as step_order,
				unnest($4::jsonb[]) as input_mapping,
				unnest($5::uuid[]) as next_node_id
		)
		INSERT INTO workflow_nodes (workflow_id, adapter_config_id, step_order, input_mapping, next_node_id, created_at)
		SELECT workflow_id, adapter_config_id, step_order, input_mapping, next_node_id, NOW()
		FROM data
		RETURNING id, workflow_id, adapter_config_id, step_order, input_mapping, next_node_id, created_at
	`

	workflowIDs := make([]string, len(nodes))
	adapterConfigIDs := make([]string, len(nodes))
	stepOrders := make([]int, len(nodes))
	inputMappings := make([]json.RawMessage, len(nodes))
	nextNodeIDs := make([]interface{}, len(nodes))

	for i, node := range nodes {
		workflowIDs[i] = node.WorkflowID
		adapterConfigIDs[i] = node.AdapterConfigID
		stepOrders[i] = node.StepOrder
		inputMappings[i], _ = json.Marshal(node.InputMapping)
		if node.NextNodeID != nil {
			nextNodeIDs[i] = *node.NextNodeID
		} else {
			nextNodeIDs[i] = nil
		}
	}

	rows, err := r.db.QueryContext(ctx, query,
		pq.Array(workflowIDs), pq.Array(adapterConfigIDs),
		pq.Array(stepOrders), pq.Array(inputMappings), pq.Array(nextNodeIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var savedNodes []workflow_node.WorkflowNode
	for rows.Next() {
		var node workflow_node.WorkflowNode
		err := rows.Scan(&node.ID, &node.WorkflowID, &node.AdapterConfigID,
			&node.StepOrder, &node.InputMapping, &node.NextNodeID, &node.CreatedAt)
		if err != nil {
			return nil, err
		}
		savedNodes = append(savedNodes, node)
	}

	return savedNodes, nil
}

func (r *WorkflowNodeRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	setClauses := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1

	fieldMap := map[string]string{
		"adapterConfigId": "adapter_config_id",
		"stepOrder":       "step_order",
		"inputMapping":    "input_mapping",
		"nextNodeId":      "next_node_id",
	}

	for key, value := range updates {
		if dbField, ok := fieldMap[key]; ok {
			if key == "inputMapping" {
				jsonValue, err := json.Marshal(value)
				if err != nil {
					return fmt.Errorf("failed to marshal input mapping: %w", err)
				}
				setClauses = append(setClauses, fmt.Sprintf("%s = $%d", dbField, argIndex))
				args = append(args, jsonValue)
			} else {
				setClauses = append(setClauses, fmt.Sprintf("%s = $%d", dbField, argIndex))
				args = append(args, value)
			}
			argIndex++
		}
	}

	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE workflow_nodes SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "), argIndex)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update workflow node: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("workflow node %s not found", id)
	}

	return nil
}

func (r *WorkflowNodeRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM workflow_nodes WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete workflow node: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("workflow node %s not found", id)
	}

	return nil
}

func (r *WorkflowNodeRepository) DeleteByWorkflowID(ctx context.Context, workflowID string) error {
	query := `DELETE FROM workflow_nodes WHERE workflow_id = $1`

	_, err := r.db.ExecContext(ctx, query, workflowID)
	if err != nil {
		return fmt.Errorf("failed to delete workflow nodes: %w", err)
	}

	return nil
}

func (r *WorkflowNodeRepository) ReorderNodes(ctx context.Context, workflowID string, nodeIDs []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for i, nodeID := range nodeIDs {
		query := `UPDATE workflow_nodes SET step_order = $1 WHERE id = $2 AND workflow_id = $3`
		_, err := tx.ExecContext(ctx, query, i+1, nodeID, workflowID)
		if err != nil {
			return fmt.Errorf("failed to reorder node %s: %w", nodeID, err)
		}
	}

	return tx.Commit()
}
