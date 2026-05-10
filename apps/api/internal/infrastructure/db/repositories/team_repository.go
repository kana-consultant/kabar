// internal/infrastructure/repository/team/repository.go
package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"seo-backend/internal/domain/team"
	"seo-backend/internal/helper"
)

type TeamRepositoryImpl struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) team.Repository {
	return &TeamRepositoryImpl{db: db}
}

func (r *TeamRepositoryImpl) GetByID(ctx context.Context, id string) (*team.Team, error) {
	query := `
		SELECT id, name, description,  status, max_members,
			created_by, created_at, updated_at
		FROM teams WHERE id = $1
	`

	var t team.Team
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.Name, &t.Description, &t.Status,
		&t.MaxMembers, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	return &t, nil
}

func (r *TeamRepositoryImpl) GetAll(ctx context.Context, query string, args []interface{}) ([]team.Team, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch teams: %w", err)
	}
	defer rows.Close()

	return r.scanTeams(rows)
}

func (r *TeamRepositoryImpl) Insert(ctx context.Context, req team.CreateTeamRequest, createdBy string) (string, error) {
	query := `
		INSERT INTO teams (name, description, status, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var teamID string
	err := r.db.QueryRowContext(
		ctx,
		query,
		req.Name, req.Description, "active", createdBy,
	).Scan(&teamID)

	return teamID, err
}

func (r *TeamRepositoryImpl) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	setClauses := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1

	fieldMap := map[string]string{
		"name":        "name",
		"description": "description",
		"":            "",
		"status":      "status",
		"maxMembers":  "max_members",
	}

	for key, value := range updates {
		if dbField, ok := fieldMap[key]; ok && value != nil {
			setClauses = append(setClauses, fmt.Sprintf("%s = $%d", dbField, argIndex))
			args = append(args, value)
			argIndex++
		}
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, helper.ParseWIBTime(time.Now().Format(time.RFC3339)))
	argIndex++

	args = append(args, id)
	query := fmt.Sprintf("UPDATE teams SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "), argIndex)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update team: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("team not found")
	}

	return nil
}

func (r *TeamRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM teams WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("team not found")
	}

	return nil
}

func (r *TeamRepositoryImpl) GetUserTeams(ctx context.Context, userID string) ([]team.Team, error) {
	query := `
		SELECT t.id, t.name, t.description,  t.status, t.max_members, 
			t.created_by, t.created_at, t.updated_at
		FROM teams t
		INNER JOIN team_members tm ON t.id = tm.team_id
		WHERE tm.user_id = $1
		ORDER BY tm.joined_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user teams: %w", err)
	}
	defer rows.Close()

	return r.scanTeams(rows)
}

func (r *TeamRepositoryImpl) scanTeams(rows *sql.Rows) ([]team.Team, error) {
	var teams []team.Team

	for rows.Next() {
		var t team.Team
		err := rows.Scan(
			&t.ID, &t.Name, &t.Description, &t.Status,
			&t.MaxMembers, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			continue
		}
		teams = append(teams, t)
	}

	return teams, nil
}
