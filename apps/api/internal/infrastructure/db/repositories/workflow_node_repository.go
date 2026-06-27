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
			input_mapping, next_node_ids, previous_node_ids,http_method, created_at
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
			input_mapping, next_node_ids, previous_node_ids,
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
			input_mapping, next_node_ids, previous_node_ids, created_at,endpoint_path,http_method
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
		var nextNodeIDs []byte
		var EndpointPath sql.NullString

		err := rows.Scan(
			&node.ID, &node.WorkflowID, &node.AdapterConfigID,
			&node.StepOrder, &inputMappingJSON, &nextNodeIDs,
			&previousNodeIDsJSON, &node.CreatedAt, &EndpointPath, &node.HTTPMethod,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workflow node: %w", err)
		}

		if len(inputMappingJSON) > 0 {
			node.InputMapping = inputMappingJSON
		}

		if EndpointPath.Valid {
			node.EndpointPath = &EndpointPath.String
		}

		if len(nextNodeIDs) > 0 {
			json.Unmarshal(nextNodeIDs, &node.NextNodeIDs)
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

// InsertBatchWithTx inserts multiple workflow nodes within a transaction.
//
// Karena setiap node bisa mereferensikan node lain sebagai "next node",
// dan node tersebut belum tentu sudah ada di DB saat insert pertama kali,
// proses dilakukan dalam 2 tahap:
//
//	Tahap 1: insert semua node dengan next_node_ids = NULL (atau '[]').
//	Tahap 2: setelah semua node ter-insert (ID-nya pasti sudah ada),
//	         baru UPDATE next_node_ids untuk masing-masing node.
func (r *WorkflowNodeRepository) InsertBatchWithTx(
	ctx context.Context,
	tx *sql.Tx,
	nodes []workflow_node.WorkflowNode,
) ([]workflow_node.WorkflowNode, error) {

	if len(nodes) == 0 {
		log.Println("[WorkflowNodeRepository.InsertBatchWithTx] no nodes to insert, skipping")
		return nil, nil
	}

	log.Println("========== InsertBatchWithTx ==========")
	log.Printf("Total nodes to insert: %d\n", len(nodes))

	// ── Validasi awal: pastikan semua UUID valid sebelum menyentuh DB ───────
	for _, node := range nodes {
		if !isValidUUID(node.ID) {
			return nil, fmt.Errorf("invalid UUID format for node ID: %s", node.ID)
		}
		if !isValidUUID(node.WorkflowID) {
			return nil, fmt.Errorf("invalid UUID format for workflow ID: %s", node.WorkflowID)
		}
		if !isValidUUID(node.AdapterConfigID) {
			return nil, fmt.Errorf("invalid UUID format for adapter config ID: %s", node.AdapterConfigID)
		}
	}

	// ── TAHAP 1: Insert semua node tanpa next_node_ids ──────────────────────
	insertQuery := `
		INSERT INTO workflow_nodes (
			id,
			workflow_id,
			adapter_config_id,
			step_order,
			input_mapping,
			next_node_ids,
			previous_node_ids,
			endpoint_path,
			http_method,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, NULL, $6, $7, $8, NOW())
		RETURNING created_at
	`

	savedNodes := make([]workflow_node.WorkflowNode, 0, len(nodes))

	for idx, node := range nodes {
		log.Printf("--- Inserting node[%d/%d] ID=%s ---\n", idx+1, len(nodes), node.ID)

		inputMappingJSON, err := sanitizeJSON(node.AdapterConfig.FieldMapping)
		if err != nil {
			log.Printf("[WorkflowNodeRepository.InsertBatchWithTx] failed to sanitize input_mapping for node %s: %v", node.ID, err)
			return nil, fmt.Errorf("failed to sanitize input_mapping for node %s: %w", node.ID, err)
		}

		// Handle previous_node_ids
		previousNodeIDs := node.PreviousNodeIDs
		if previousNodeIDs == nil {
			previousNodeIDs = []string{}
		}
		previousNodeIDsJSON, err := json.Marshal(previousNodeIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal previous_node_ids for node %s: %w", node.ID, err)
		}

		savedNode := node
		err = tx.QueryRowContext(
			ctx,
			insertQuery,
			savedNode.ID,
			savedNode.WorkflowID,
			savedNode.AdapterConfigID,
			savedNode.StepOrder,
			json.RawMessage(inputMappingJSON),
			json.RawMessage(previousNodeIDsJSON),
			savedNode.AdapterConfig.EndpointPath,
			savedNode.AdapterConfig.HTTPMethod,
		).Scan(&savedNode.CreatedAt)

		if err != nil {
			log.Printf("❌ Failed to insert node %s: %v\n", savedNode.ID, err)
			return nil, fmt.Errorf("failed to insert workflow node %s: %w", savedNode.ID, err)
		}

		log.Printf("✅ Node inserted (without next_node_ids): %s\n", savedNode.ID)
		savedNodes = append(savedNodes, savedNode)
	}

	// ── TAHAP 2: Update next_node_ids setelah semua node ter-insert ─────────
	updateQuery := `UPDATE workflow_nodes SET next_node_ids = $1 WHERE id = $2`

	for _, node := range nodes {
		nextNodeIDs := node.NextNodeIDs
		if nextNodeIDs == nil {
			nextNodeIDs = []string{}
		}

		nextNodeIDsJSON, err := json.Marshal(nextNodeIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal next_node_ids for node %s: %w", node.ID, err)
		}

		log.Printf("[WorkflowNodeRepository.InsertBatchWithTx] updating next_node_ids for node id=%s -> %s",
			node.ID, string(nextNodeIDsJSON))

		result, err := tx.ExecContext(ctx, updateQuery, json.RawMessage(nextNodeIDsJSON), node.ID)
		if err != nil {
			log.Printf("❌ Failed to update next_node_ids for node %s: %v\n", node.ID, err)
			return nil, fmt.Errorf("failed to update next_node_ids for node %s: %w", node.ID, err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("failed to check rows affected for node %s: %w", node.ID, err)
		}
		if rowsAffected == 0 {
			return nil, fmt.Errorf("no rows updated for next_node_ids on node %s (node not found?)", node.ID)
		}
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
            next_node_ids, 
            previous_node_ids, 
            created_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
        RETURNING created_at
    `

	var nextNodeIDs interface{}
	if len(node.NextNodeIDs) > 0 {
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

	// ✅ Gunakan node.ID yang sudah digenerate
	err := tx.QueryRowContext(ctx, query,
		node.ID, // Tambahkan ini
		node.WorkflowID,
		node.AdapterConfigID,
		node.StepOrder,
		node.AdapterConfig.FieldMapping,
		nextNodeIDs,
		previousNodeIDs,
	).Scan(&node.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert workflow node with tx: %w", err)
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

		// Marshal next node IDs
		nextNodeIDs := node.NextNodeIDs
		if nextNodeIDs == nil {
			nextNodeIDs = []string{}
		}
		nextNodeIDsJSON, err := json.Marshal(nextNodeIDs)
		if err != nil {
			return fmt.Errorf("failed to marshal next_node_ids for node %s: %w", node.ID, err)
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
			json.RawMessage(nextNodeIDsJSON),  // jsonb
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
			next_node_ids,
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
			next_node_ids = EXCLUDED.next_node_ids,
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
