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
			input_mapping, next_node_id, previous_node_ids, created_at
		FROM workflow_nodes WHERE id = $1
	`

	var node workflow_node.WorkflowNode
	var inputMappingJSON []byte
	var nextNodeID sql.NullString
	var previousNodeIDs pq.StringArray

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&node.ID, &node.WorkflowID, &node.AdapterConfigID,
		&node.StepOrder, &inputMappingJSON, &nextNodeID,
		&previousNodeIDs, &node.CreatedAt,
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

	node.PreviousNodeIDs = previousNodeIDs

	return &node, nil
}

func (r *WorkflowNodeRepository) GetByWorkflowID(ctx context.Context, workflowID string) ([]*workflow_node.WorkflowNode, error) {
	query := `
		SELECT id, workflow_id, adapter_config_id, step_order, 
			input_mapping, next_node_id, previous_node_ids, created_at
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
		var previousNodeIDs pq.StringArray

		err := rows.Scan(
			&node.ID, &node.WorkflowID, &node.AdapterConfigID,
			&node.StepOrder, &inputMappingJSON, &nextNodeID,
			&previousNodeIDs, &node.CreatedAt,
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

		node.PreviousNodeIDs = previousNodeIDs

		nodes = append(nodes, &node)
	}

	return nodes, nil
}

func (r *WorkflowNodeRepository) Insert(ctx context.Context, node *workflow_node.WorkflowNode) error {
	inputMappingJSON, err := json.Marshal(node.InputMapping)
	if err != nil {
		return fmt.Errorf("failed to marshal input mapping: %w", err)
	}

	query := `
		INSERT INTO workflow_nodes (workflow_id, adapter_config_id, step_order, 
			input_mapping, next_node_id, previous_node_ids, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING id, created_at
	`

	var nextNodeID interface{}
	if node.NextNodeID != nil {
		nextNodeID = *node.NextNodeID
	}

	var previousNodeIDs interface{}
	if len(node.PreviousNodeIDs) > 0 {
		previousNodeIDs = pq.Array(node.PreviousNodeIDs)
	} else {
		previousNodeIDs = pq.Array([]string{})
	}

	err = r.db.QueryRowContext(ctx, query,
		node.WorkflowID, node.AdapterConfigID, node.StepOrder,
		inputMappingJSON, nextNodeID, previousNodeIDs,
	).Scan(&node.ID, &node.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert workflow node: %w", err)
	}

	return nil
}

func (r *WorkflowNodeRepository) InsertBatch(ctx context.Context, nodes []workflow_node.WorkflowNode) ([]workflow_node.WorkflowNode, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO workflow_nodes (workflow_id, adapter_config_id, step_order, 
			input_mapping, next_node_id, previous_node_ids, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING id, created_at
	`

	var savedNodes []workflow_node.WorkflowNode

	for _, node := range nodes {
		inputMappingJSON, err := json.Marshal(node.InputMapping)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal input mapping: %w", err)
		}

		var nextNodeID interface{}
		if node.NextNodeID != nil {
			nextNodeID = *node.NextNodeID
		}

		var previousNodeIDs interface{}
		if len(node.PreviousNodeIDs) > 0 {
			previousNodeIDs = pq.Array(node.PreviousNodeIDs)
		} else {
			previousNodeIDs = pq.Array([]string{})
		}

		var savedNode workflow_node.WorkflowNode
		savedNode.WorkflowID = node.WorkflowID
		savedNode.AdapterConfigID = node.AdapterConfigID
		savedNode.StepOrder = node.StepOrder
		savedNode.InputMapping = node.InputMapping
		savedNode.NextNodeID = node.NextNodeID
		savedNode.PreviousNodeIDs = node.PreviousNodeIDs

		err = tx.QueryRowContext(ctx, query,
			node.WorkflowID, node.AdapterConfigID, node.StepOrder,
			inputMappingJSON, nextNodeID, previousNodeIDs,
		).Scan(&savedNode.ID, &savedNode.CreatedAt)

		if err != nil {
			return nil, fmt.Errorf("failed to insert workflow node: %w", err)
		}

		savedNodes = append(savedNodes, savedNode)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
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
		"previousNodeIds": "previous_node_ids",
	}

	for key, value := range updates {
		if dbField, ok := fieldMap[key]; ok {
			switch key {
			case "inputMapping":
				jsonValue, err := json.Marshal(value)
				if err != nil {
					return fmt.Errorf("failed to marshal input mapping: %w", err)
				}
				setClauses = append(setClauses, fmt.Sprintf("%s = $%d", dbField, argIndex))
				args = append(args, jsonValue)
			case "previousNodeIds":
				if prevIDs, ok := value.([]string); ok {
					setClauses = append(setClauses, fmt.Sprintf("%s = $%d", dbField, argIndex))
					args = append(args, pq.Array(prevIDs))
				} else {
					return fmt.Errorf("previousNodeIds must be []string")
				}
			default:
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

// ============ TRANSACTION METHODS ============

// InsertBatchWithTx inserts multiple workflow nodes within a transaction
func (r *WorkflowNodeRepository) InsertBatchWithTx(ctx context.Context, tx *sql.Tx, nodes []workflow_node.WorkflowNode) ([]workflow_node.WorkflowNode, error) {
	query := `
		INSERT INTO workflow_nodes (workflow_id, adapter_config_id, step_order, 
			input_mapping, next_node_id, previous_node_ids, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING id, created_at
	`

	var savedNodes []workflow_node.WorkflowNode

	for _, node := range nodes {
		inputMappingJSON, err := json.Marshal(node.InputMapping)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal input mapping: %w", err)
		}

		var nextNodeID interface{}
		if node.NextNodeID != nil {
			nextNodeID = *node.NextNodeID
		}

		var previousNodeIDs interface{}
		if len(node.PreviousNodeIDs) > 0 {
			previousNodeIDs = pq.Array(node.PreviousNodeIDs)
		} else {
			previousNodeIDs = pq.Array([]string{})
		}

		var savedNode workflow_node.WorkflowNode
		savedNode.WorkflowID = node.WorkflowID
		savedNode.AdapterConfigID = node.AdapterConfigID
		savedNode.StepOrder = node.StepOrder
		savedNode.InputMapping = node.InputMapping
		savedNode.NextNodeID = node.NextNodeID
		savedNode.PreviousNodeIDs = node.PreviousNodeIDs

		err = tx.QueryRowContext(ctx, query,
			node.WorkflowID, node.AdapterConfigID, node.StepOrder,
			inputMappingJSON, nextNodeID, previousNodeIDs,
		).Scan(&savedNode.ID, &savedNode.CreatedAt)

		if err != nil {
			return nil, fmt.Errorf("failed to insert workflow node: %w", err)
		}

		savedNodes = append(savedNodes, savedNode)
	}

	return savedNodes, nil
}

// InsertWithTx inserts a single workflow node within a transaction
func (r *WorkflowNodeRepository) InsertWithTx(ctx context.Context, tx *sql.Tx, node *workflow_node.WorkflowNode) error {
	inputMappingJSON, err := json.Marshal(node.InputMapping)
	if err != nil {
		return fmt.Errorf("failed to marshal input mapping: %w", err)
	}

	query := `
		INSERT INTO workflow_nodes (workflow_id, adapter_config_id, step_order, 
			input_mapping, next_node_id, previous_node_ids, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING id, created_at
	`

	var nextNodeID interface{}
	if node.NextNodeID != nil {
		nextNodeID = *node.NextNodeID
	}

	var previousNodeIDs interface{}
	if len(node.PreviousNodeIDs) > 0 {
		previousNodeIDs = pq.Array(node.PreviousNodeIDs)
	} else {
		previousNodeIDs = pq.Array([]string{})
	}

	err = tx.QueryRowContext(ctx, query,
		node.WorkflowID, node.AdapterConfigID, node.StepOrder,
		inputMappingJSON, nextNodeID, previousNodeIDs,
	).Scan(&node.ID, &node.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert workflow node with tx: %w", err)
	}

	return nil
}

// UpdateWithTx updates a workflow node within a transaction
func (r *WorkflowNodeRepository) UpdateWithTx(ctx context.Context, tx *sql.Tx, id string, updates map[string]interface{}) error {
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
		"previousNodeIds": "previous_node_ids",
	}

	for key, value := range updates {
		if dbField, ok := fieldMap[key]; ok {
			switch key {
			case "inputMapping":
				jsonValue, err := json.Marshal(value)
				if err != nil {
					return fmt.Errorf("failed to marshal input mapping: %w", err)
				}
				setClauses = append(setClauses, fmt.Sprintf("%s = $%d", dbField, argIndex))
				args = append(args, jsonValue)
			case "previousNodeIds":
				if prevIDs, ok := value.([]string); ok {
					setClauses = append(setClauses, fmt.Sprintf("%s = $%d", dbField, argIndex))
					args = append(args, pq.Array(prevIDs))
				} else {
					return fmt.Errorf("previousNodeIds must be []string")
				}
			default:
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

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update workflow node with tx: %w", err)
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
