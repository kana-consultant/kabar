package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"seo-backend/internal/domain/workflow"
	"seo-backend/internal/helper"
)

type WorkflowDefinitionRepository struct {
	db *sql.DB
}

func NewWorkflowDefinitionRepository(db *sql.DB) workflow.WorkflowDefinitionRepository {
	return &WorkflowDefinitionRepository{db: db}
}

func (r *WorkflowDefinitionRepository) GetByID(ctx context.Context, id string) (*workflow.WorkflowDefinition, error) {
	query := `
		SELECT id, product_id, name, created_at, updated_at
		FROM workflow_definitions WHERE id = $1
	`

	var wf workflow.WorkflowDefinition
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&wf.ID, &wf.ProductID, &wf.Name, &wf.CreatedAt, &wf.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get workflow definition: %w", err)
	}

	return &wf, nil
}

func (r *WorkflowDefinitionRepository) GetByProductID(ctx context.Context, productID string) ([]*workflow.WorkflowDefinition, error) {
	query := `
		SELECT id, product_id, name, created_at, updated_at
		FROM workflow_definitions WHERE product_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow definitions: %w", err)
	}
	defer rows.Close()

	var workflows []*workflow.WorkflowDefinition
	for rows.Next() {
		var wf workflow.WorkflowDefinition
		err := rows.Scan(&wf.ID, &wf.ProductID, &wf.Name, &wf.CreatedAt, &wf.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workflow definition: %w", err)
		}
		workflows = append(workflows, &wf)
	}

	return workflows, nil
}

func (r *WorkflowDefinitionRepository) Insert(ctx context.Context, wf *workflow.WorkflowDefinition) error {
	query := `
		INSERT INTO workflow_definitions (product_id, name, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query, wf.ProductID, wf.Name).Scan(
		&wf.ID, &wf.CreatedAt, &wf.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to insert workflow definition: %w", err)
	}

	return nil
}

func (r *WorkflowDefinitionRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	setClauses := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1

	fieldMap := map[string]string{
		"name": "name",
	}

	for key, value := range updates {
		if dbField, ok := fieldMap[key]; ok {
			setClauses = append(setClauses, fmt.Sprintf("%s = $%d", dbField, argIndex))
			args = append(args, value)
			argIndex++
		}
	}

	if len(setClauses) == 0 {
		return nil
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, helper.ParseWIBTime(time.Now().Format(time.RFC3339)))
	argIndex++

	args = append(args, id)
	query := fmt.Sprintf("UPDATE workflow_definitions SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "), argIndex)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update workflow definition: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("workflow definition %s not found", id)
	}

	return nil
}

func (r *WorkflowDefinitionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM workflow_definitions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete workflow definition: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("workflow definition %s not found", id)
	}

	return nil
}

func (r *WorkflowDefinitionRepository) InsertWithTx(ctx context.Context, tx *sql.Tx, wf *workflow.WorkflowDefinition) error {
	query := `
        INSERT INTO workflow_definitions (id, product_id, name, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `

	err := tx.QueryRowContext(ctx, query, wf.ID, wf.ProductID, wf.Name).Scan(
		&wf.ID, &wf.CreatedAt, &wf.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to insert workflow definition with tx: %w", err)
	}

	return nil
}
