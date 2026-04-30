package team

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"seo-backend/internal/models"
)

type MemberRepository struct {
	db *sql.DB
}

func NewMemberRepository(db *sql.DB) *MemberRepository {
	return &MemberRepository{db: db}
}

func (r *MemberRepository) GetByTeamID(teamID string, filters MemberFilters) ([]models.TeamMember, error) {
	query, args := r.buildMemberListQuery(teamID, filters)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch team members: %w", err)
	}
	defer rows.Close()

	return r.scanMembers(rows)
}

func (r *MemberRepository) Add(tx *sql.Tx, teamID, userID string, role models.TeamMemberRole) error {
	query := `
		INSERT INTO team_members (team_id, user_id, role, joined_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := tx.Exec(query, teamID, userID, role, time.Now())
	return err
}

func (r *MemberRepository) UpdateRole(teamID, userID string, role models.TeamMemberRole) error {
	query := `UPDATE team_members SET role = $1 WHERE team_id = $2 AND user_id = $3`

	result, err := r.db.Exec(query, role, teamID, userID)
	if err != nil {
		return fmt.Errorf("failed to update member role: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("member not found")
	}

	return nil
}

func (r *MemberRepository) Remove(teamID, userID string) error {
	query := `DELETE FROM team_members WHERE team_id = $1 AND user_id = $2`

	result, err := r.db.Exec(query, teamID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove member: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("member not found")
	}

	return nil
}

func (r *MemberRepository) Exists(tx *sql.Tx, teamID, userID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM team_members WHERE team_id = $1 AND user_id = $2)`
	err := tx.QueryRow(query, teamID, userID).Scan(&exists)
	return exists, err
}

func (r *MemberRepository) GetCount(teamID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM team_members WHERE team_id = $1`
	err := r.db.QueryRow(query, teamID).Scan(&count)
	return count, err
}

func (r *MemberRepository) GetMaxMembers(teamID string) (int, error) {
	var maxMembers int
	query := `SELECT COALESCE(max_members, 10) FROM teams WHERE id = $1`
	err := r.db.QueryRow(query, teamID).Scan(&maxMembers)
	return maxMembers, err
}

func (r *MemberRepository) buildMemberListQuery(teamID string, filters MemberFilters) (string, []interface{}) {
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1

	conditions = append(conditions, fmt.Sprintf("team_id = $%d", argIndex))
	args = append(args, teamID)
	argIndex++

	if filters.Role != "" {
		conditions = append(conditions, fmt.Sprintf("role = $%d", argIndex))
		args = append(args, filters.Role)
		argIndex++
	}

	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	query := fmt.Sprintf(`
		SELECT id, team_id, user_id, role, joined_at
		FROM team_members %s
		ORDER BY 
			CASE WHEN role = 'manager' THEN 1 ELSE 2 END,
			joined_at ASC
	`, whereClause)

	return query, args
}

func (r *MemberRepository) scanMembers(rows *sql.Rows) ([]models.TeamMember, error) {
	var members []models.TeamMember

	for rows.Next() {
		var m models.TeamMember
		err := rows.Scan(
			&m.ID, &m.TeamID, &m.UserID, &m.Role, &m.JoinedAt,
		)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		members = append(members, m)
	}

	return members, nil
}
