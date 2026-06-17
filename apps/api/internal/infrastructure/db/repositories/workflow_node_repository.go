// internal/infrastructure/db/repositories/workflow_node_repository.go
package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"seo-backend/internal/domain/workflow_node"
	"strings"

	"github.com/google/uuid"
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
			input_mapping, next_node_id, previous_node_ids,http_method, created_at
		FROM workflow_nodes WHERE id = $1
	`

	var node workflow_node.WorkflowNode
	var inputMappingJSON sql.NullString
	var nextNodeID pq.StringArray
	var previousNodeIDs pq.StringArray

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&node.ID, &node.WorkflowID, &node.AdapterConfigID,
		&node.StepOrder, &inputMappingJSON, &nextNodeID,
		&previousNodeIDs, &node.AdapterConfig.HTTPMethod, &node.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get workflow node: %w", err)
	}

	if inputMappingJSON.Valid {
		node.AdapterConfig.EndpointPath = inputMappingJSON.String
	}

	node.NextNodeIDs = nextNodeID

	node.PreviousNodeIDs = previousNodeIDs

	return &node, nil
}

func (r *WorkflowNodeRepository) GetByWorkflowID(ctx context.Context, workflowID string) ([]*workflow_node.WorkflowNode, error) {
	query := `
    SELECT id, workflow_id, adapter_config_id, step_order, 
			input_mapping, next_node_id, previous_node_ids,
			created_at
		FROM workflow_nodes 
		WHERE workflow_id = ANY($1)
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
		var inputMappingJSON sql.NullString
		var nextNodeID pq.StringArray
		var previousNodeIDs pq.StringArray

		err := rows.Scan(
			&node.ID, &node.WorkflowID, &node.AdapterConfigID,
			&node.StepOrder, &inputMappingJSON, &nextNodeID,
			&previousNodeIDs, &node.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workflow node: %w", err)
		}

		if inputMappingJSON.Valid {
			node.AdapterConfig.EndpointPath = inputMappingJSON.String
		}

		node.NextNodeIDs = nextNodeID

		node.PreviousNodeIDs = previousNodeIDs

		nodes = append(nodes, &node)
	}

	return nodes, nil
}

func (r *WorkflowNodeRepository) GetByWorkflowIDs(ctx context.Context, workflowIDs []string) ([]*workflow_node.WorkflowNodeCreate, error) {
	query := `
		SELECT id, workflow_id, adapter_config_id, step_order, 
			input_mapping, next_node_id, previous_node_ids, created_at,endpoint_path,http_method
		FROM workflow_nodes 
		WHERE workflow_id = ANY($1)
		ORDER BY step_order ASC
	`

	rows, err := r.db.QueryContext(ctx, query, pq.Array(workflowIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow nodes: %w", err)
	}
	defer rows.Close()

	var nodes []*workflow_node.WorkflowNodeCreate
	for rows.Next() {
		var node workflow_node.WorkflowNodeCreate
		var inputMappingJSON []byte
		var previousNodeIDsJSON []byte
		var nextNodeID []byte
		var EndpointPath sql.NullString

		err := rows.Scan(
			&node.ID, &node.WorkflowID, &node.AdapterConfigID,
			&node.StepOrder, &inputMappingJSON, &nextNodeID,
			&previousNodeIDsJSON, &node.CreatedAt, &EndpointPath, &node.HTTPMethod,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workflow node: %w", err)
		}

		if len(inputMappingJSON) > 0 {
			json.Unmarshal(inputMappingJSON, &node.InputMapping)
		}

		if EndpointPath.Valid {
			node.EndpointPath = &EndpointPath.String
		}

		if len(nextNodeID) > 0 {
			json.Unmarshal(nextNodeID, &node.NextNodeIDs)
		} else {
			node.NextNodeIDs = []string{}
		}

		if len(previousNodeIDsJSON) > 0 {
			json.Unmarshal(previousNodeIDsJSON, &node.PreviousNodeIDs)
		} else {
			node.PreviousNodeIDs = []string{}
		}

		nodes = append(nodes, &node)
	}

	return nodes, nil
}

func (r *WorkflowNodeRepository) Insert(ctx context.Context, node *workflow_node.WorkflowNode) error {
	inputMappingJSON, err := json.Marshal(node.AdapterConfig.FieldMapping)
	if err != nil {
		return fmt.Errorf("failed to marshal input mapping: %w", err)
	}

	query := `
		INSERT INTO workflow_nodes (workflow_id, adapter_config_id, step_order, 
			input_mapping, next_node_id, previous_node_ids, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING id, created_at
	`

	var NextNodeIDs interface{}
	if len(node.PreviousNodeIDs) > 0 {
		NextNodeIDs = pq.Array(node.NextNodeIDs)
	} else {
		NextNodeIDs = pq.Array([]string{})
	}

	var previousNodeIDs interface{}
	if len(node.PreviousNodeIDs) > 0 {
		previousNodeIDs = pq.Array(node.PreviousNodeIDs)
	} else {
		previousNodeIDs = pq.Array([]string{})
	}

	err = r.db.QueryRowContext(ctx, query,
		node.WorkflowID, node.AdapterConfigID, node.StepOrder,
		inputMappingJSON, NextNodeIDs, previousNodeIDs,
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
		inputMappingJSON, err := json.Marshal(node.AdapterConfig.FieldMapping)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal input mapping: %w", err)
		}

		var nextNodeIDs interface{}
		if len(node.PreviousNodeIDs) > 0 {
			nextNodeIDs = pq.Array(node.NextNodeIDs)
		} else {
			nextNodeIDs = pq.Array([]string{})
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
		savedNode.AdapterConfig.FieldMapping = node.AdapterConfig.FieldMapping
		savedNode.NextNodeIDs = node.NextNodeIDs
		savedNode.PreviousNodeIDs = node.PreviousNodeIDs

		err = tx.QueryRowContext(ctx, query,
			node.WorkflowID, node.AdapterConfigID, node.StepOrder,
			inputMappingJSON, nextNodeIDs, previousNodeIDs,
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
func (r *WorkflowNodeRepository) InsertBatchWithTx(
	ctx context.Context,
	tx *sql.Tx,
	nodes []workflow_node.WorkflowNodeCreate,
) ([]workflow_node.WorkflowNodeCreate, error) {

	log.Println("========== InsertBatchWithTx ==========")
	log.Printf("Total nodes to insert: %d\n", len(nodes))

	// ── TAHAP 1: Insert semua node tanpa next_node_id ──────────────────────
	insertQuery := `
		INSERT INTO workflow_nodes (
			id,
			workflow_id,
			adapter_config_id,
			step_order,
			input_mapping,
			next_node_id,
			previous_node_ids,
			endpoint_path,
			http_method,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, NULL, $6,$7, NOW())
		RETURNING created_at
	`

	var savedNodes []workflow_node.WorkflowNodeCreate

	for idx, node := range nodes {
		log.Printf("--- Inserting node[%d/%d] ID=%s ---\n", idx+1, len(nodes), node.ID)

		// Handle input_mapping
		var inputMappingJSON []byte
		switch {
		case len(node.InputMapping) == 0,
			string(node.InputMapping) == "null",
			string(node.InputMapping) == "[]":
			inputMappingJSON = []byte(`{}`)
		default:
			inputMappingJSON = node.InputMapping
		}
		log.Printf("InputMapping JSON: %s\n", string(inputMappingJSON))

		// Handle previous_node_ids
		var previousNodeIDsJSON []byte
		if len(node.PreviousNodeIDs) == 0 {
			previousNodeIDsJSON = []byte(`[]`)
		} else {
			b, err := json.Marshal(node.PreviousNodeIDs)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal previous_node_ids for node %s: %w", node.ID, err)
			}
			previousNodeIDsJSON = b
		}
		log.Printf("PreviousNodeIDs JSON: %s\n", string(previousNodeIDsJSON))

		savedNode := node
		err := tx.QueryRowContext(
			ctx,
			insertQuery,
			savedNode.ID,
			savedNode.WorkflowID,
			savedNode.AdapterConfigID,
			savedNode.StepOrder,
			inputMappingJSON,
			previousNodeIDsJSON,
			savedNode.HTTPMethod,
			savedNode.EndpointPath,
		).Scan(&savedNode.CreatedAt)

		if err != nil {
			log.Printf("❌ Failed to insert node %s: %v\n", savedNode.ID, err)
			return nil, fmt.Errorf("failed to insert workflow node %s: %w", savedNode.ID, err)
		}

		log.Printf("✅ Node inserted (without next_node_id): %s\n", savedNode.ID)
		savedNodes = append(savedNodes, savedNode)
	}

	// ── TAHAP 2: Update next_node_id setelah semua node ter-insert ─────────
	updateQuery := `UPDATE workflow_nodes SET next_node_id = $1 WHERE id = $2`

	for _, node := range nodes {
		if node.NextNodeIDs == nil {
			continue
		}

		log.Printf("Updating next_node_id for node %s -> %s\n", node.ID, node.NextNodeIDs)

		_, err := tx.ExecContext(ctx, updateQuery, node.NextNodeIDs, node.ID)
		if err != nil {
			log.Printf("❌ Failed to update next_node_id for node %s: %v\n", node.ID, err)
			return nil, fmt.Errorf("failed to update next_node_id for node %s: %w", node.ID, err)
		}

		// Update savedNodes juga supaya return value akurat
		for i := range savedNodes {
			if savedNodes[i].ID == node.ID {
				savedNodes[i].NextNodeIDs = node.NextNodeIDs
				break
			}
		}

		log.Printf("✅ next_node_id updated for node %s\n", node.ID)
	}

	log.Printf("✅ All %d nodes inserted successfully\n", len(savedNodes))
	log.Println("=====================================")

	return savedNodes, nil
}

// InsertWithTx inserts a single workflow node within a transaction
func (r *WorkflowNodeRepository) InsertWithTx(ctx context.Context, tx *sql.Tx, node *workflow_node.WorkflowNode) error {

	// ✅ Tambahkan id ke query
	query := `
        INSERT INTO workflow_nodes (
            id, 
            workflow_id, 
            adapter_config_id, 
            step_order, 
            input_mapping, 
            next_node_id, 
            previous_node_ids, 
            created_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
        RETURNING created_at
    `

	var nextNodeID interface{}
	if node.NextNodeIDs != nil {
		nextNodeID = node.NextNodeIDs
	}

	var previousNodeIDs interface{}
	if len(node.PreviousNodeIDs) > 0 {
		previousNodeIDs = pq.Array(node.PreviousNodeIDs)
	} else {
		previousNodeIDs = pq.Array([]string{})
	}

	// ✅ Gunakan node.ID yang sudah digenerate
	err := tx.QueryRowContext(ctx, query,
		node.ID, // Tambahkan ini
		node.WorkflowID,
		node.AdapterConfigID,
		node.StepOrder,
		node.AdapterConfig.FieldMapping,
		nextNodeID,
		previousNodeIDs,
	).Scan(&node.CreatedAt)

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
		"http_method":     "http_method",
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

// UpsertWithTx upserts multiple workflow nodes within a transaction
func (r *WorkflowNodeRepository) UpsertWithTx(ctx context.Context, tx *sql.Tx, nodes []workflow_node.WorkflowNode) error {
	if len(nodes) == 0 {
		log.Println("[WorkflowNodeRepository.UpsertWithTx] no nodes to upsert, skipping")
		return nil
	}

	log.Printf("[WorkflowNodeRepository.UpsertWithTx] upserting %d nodes", len(nodes))

	const fieldsPerNode = 9

	placeholders := make([]string, 0, len(nodes))
	args := make([]interface{}, 0, len(nodes)*fieldsPerNode)
	argIndex := 1

	for _, node := range nodes {
		// Validate UUIDs first, before doing any marshaling work
		if !isValidUUID(node.ID) {
			return fmt.Errorf("invalid UUID format for node ID: %s", node.ID)
		}
		if !isValidUUID(node.WorkflowID) {
			return fmt.Errorf("invalid UUID format for workflow ID: %s", node.WorkflowID)
		}
		if !isValidUUID(node.AdapterConfigID) {
			return fmt.Errorf("invalid UUID format for adapter config ID: %s", node.AdapterConfigID)
		}

		// Marshal input mapping
		inputMappingJSON, err := sanitizeJSON(node.AdapterConfig.FieldMapping)
		if err != nil {
			log.Printf("[WorkflowNodeRepository.UpsertWithTx] failed to sanitize input_mapping for node %s: %v", node.ID, err)
			return fmt.Errorf("failed to sanitize input_mapping for node %s: %w", node.ID, err)
		}

		// Marshal previous node IDs
		previousNodeIDs := node.PreviousNodeIDs
		if previousNodeIDs == nil {
			previousNodeIDs = []string{}
		}
		previousNodeIDsJSON, err := json.Marshal(previousNodeIDs)
		if err != nil {
			return fmt.Errorf("failed to marshal previous_node_ids for node %s: %w", node.ID, err)
		}

		// 9 placeholders for 9 fields
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			argIndex, argIndex+1, argIndex+2, argIndex+3, argIndex+4, argIndex+5, argIndex+6, argIndex+7, argIndex+8,
		))

		args = append(args,
			node.ID,                           // uuid
			node.WorkflowID,                   // uuid
			node.AdapterConfigID,              // uuid
			previousNodeIDsJSON,               // jsonb
			node.StepOrder,                    // integer
			node.NextNodeIDs,                  // uuid (nullable)
			node.AdapterConfig.EndpointPath,   // text
			node.AdapterConfig.HTTPMethod,     // text
			json.RawMessage(inputMappingJSON), // jsonb
		)

		log.Printf("[WorkflowNodeRepository.UpsertWithTx] preparing node id=%s workflow_id=%s adapter_config_id=%s step_order=%d",
			node.ID, node.WorkflowID, node.AdapterConfigID, node.StepOrder)

		argIndex += fieldsPerNode // advance by 9, not 1
	}

	query := fmt.Sprintf(`
		INSERT INTO workflow_nodes (
			id, 
			workflow_id, 
			adapter_config_id, 
			previous_node_ids, 
			step_order, 
			next_node_id,
			endpoint_path,
			http_method,
			input_mapping
		)
		VALUES %s
		ON CONFLICT (id) DO UPDATE SET
			workflow_id = EXCLUDED.workflow_id,
			adapter_config_id = EXCLUDED.adapter_config_id,
			previous_node_ids = EXCLUDED.previous_node_ids,
			step_order = EXCLUDED.step_order,
			next_node_id = EXCLUDED.next_node_id,
			endpoint_path = EXCLUDED.endpoint_path,
			http_method = EXCLUDED.http_method,
			input_mapping = EXCLUDED.input_mapping
	`, strings.Join(placeholders, ", "))

	log.Printf("[WorkflowNodeRepository.UpsertWithTx] executing query with %d args", len(args))

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		log.Printf("[WorkflowNodeRepository.UpsertWithTx] failed to upsert: %v", err)
		return fmt.Errorf("failed to upsert workflow nodes: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("[WorkflowNodeRepository.UpsertWithTx] failed to get rows affected: %v", err)
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	log.Printf("[WorkflowNodeRepository.UpsertWithTx] successfully upserted %d rows", rowsAffected)

	return nil
}

// Validasi UUID sebelum insert
func isValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

func sanitizeJSON(raw string) (json.RawMessage, error) {
	if raw == "" {
		return json.RawMessage("{}"), nil
	}
	if !json.Valid([]byte(raw)) {
		return nil, fmt.Errorf("invalid json: %s", raw)
	}
	return json.RawMessage(raw), nil
}
