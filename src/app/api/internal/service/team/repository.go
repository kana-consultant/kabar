package team

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"seo-backend/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(id string) (*models.Team, error) {
	query := `
		SELECT id, name, description, logo, status, max_members,
			created_by, created_at, updated_at
		FROM teams WHERE id = $1
	`

	var team models.Team
	err := r.db.QueryRow(query, id).Scan(
		&team.ID, &team.Name, &team.Description, &team.Logo, &team.Status,
		&team.MaxMembers, &team.CreatedBy, &team.CreatedAt, &team.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	return &team, nil
}

func (r *Repository) GetAll(query string, args []interface{}) ([]models.Team, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch teams: %w", err)
	}
	defer rows.Close()

	return r.scanTeams(rows)
}

func (r *Repository) Insert(req models.CreateTeamRequest, createdBy string) (string, error) {
	query := `
		INSERT INTO teams (name, description, status, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var teamID string
	err := r.db.QueryRow(
		query,
		req.Name, req.Description, "active", createdBy,
	).Scan(&teamID)

	return teamID, err
}

func (r *Repository) Update(id string, updates map[string]interface{}) error {
	setClauses := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1

	fieldMap := map[string]string{
		"name":        "name",
		"description": "description",
		"logo":        "logo",
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
	args = append(args, time.Now())
	argIndex++

	args = append(args, id)
	query := fmt.Sprintf("UPDATE teams SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "), argIndex)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update team: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("team not found")
	}

	return nil
}

func (r *Repository) Delete(id string) error {
	query := "DELETE FROM teams WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("team not found")
	}

	return nil
}

func (r *Repository) GetUserTeams(userID string) ([]models.Team, error) {
	query := `
		SELECT t.id, t.name, t.description, t.logo, t.status, t.max_members, 
			t.created_by, t.created_at, t.updated_at
		FROM teams t
		INNER JOIN team_members tm ON t.id = tm.team_id
		WHERE tm.user_id = $1
		ORDER BY tm.joined_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user teams: %w", err)
	}
	defer rows.Close()

	return r.scanTeams(rows)
}

func (r *Repository) scanTeams(rows *sql.Rows) ([]models.Team, error) {
	var teams []models.Team

	for rows.Next() {
		var t models.Team
		err := rows.Scan(
			&t.ID, &t.Name, &t.Description, &t.Logo, &t.Status,
			&t.MaxMembers, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		teams = append(teams, t)
	}

	return teams, nil
}
