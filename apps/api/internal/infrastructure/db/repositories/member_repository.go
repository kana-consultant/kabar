package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"seo-backend/internal/domain/team"
	"seo-backend/internal/helper"
)

type MemberRepositoryImpl struct {
	db *sql.DB
}

func NewMemberRepository(db *sql.DB) team.MemberRepository {
	return &MemberRepositoryImpl{db: db}
}

func (r *MemberRepositoryImpl) GetByTeamID(ctx context.Context, teamID string, filters team.MemberFilters) ([]team.TeamMember, error) {
	query, args := r.buildMemberListQuery(teamID, filters)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch team members: %w", err)
	}
	defer rows.Close()

	return r.scanMembers(rows)
}

func (r *MemberRepositoryImpl) Add(ctx context.Context, tx *sql.Tx, teamID, userID string, role team.TeamMemberRole) error {
	query := `
		INSERT INTO team_members (team_id, user_id, role, joined_at)
		VALUES ($1, $2, $3, $4)
	`

	if tx != nil {
		_, err := tx.ExecContext(ctx, query, teamID, userID, role, helper.ParseWIBTime(time.Now().Format(time.RFC3339)))
		return err
	}

	_, err := r.db.ExecContext(ctx, query, teamID, userID, role, helper.ParseWIBTime(time.Now().Format(time.RFC3339)))
	return err
}

func (r *MemberRepositoryImpl) UpdateRole(ctx context.Context, teamID, userID string, role team.TeamMemberRole) error {
	query := `UPDATE team_members SET role = $1 WHERE team_id = $2 AND user_id = $3`

	result, err := r.db.ExecContext(ctx, query, role, teamID, userID)
	if err != nil {
		return fmt.Errorf("failed to update member role: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("member not found")
	}

	return nil
}

func (r *MemberRepositoryImpl) Remove(ctx context.Context, teamID, userID string) error {
	query := `DELETE FROM team_members WHERE team_id = $1 AND user_id = $2`

	result, err := r.db.ExecContext(ctx, query, teamID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove member: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("member not found")
	}

	return nil
}

func (r *MemberRepositoryImpl) Exists(ctx context.Context, tx *sql.Tx, teamID, userID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM team_members WHERE team_id = $1 AND user_id = $2)`

	if tx != nil {
		err := tx.QueryRowContext(ctx, query, teamID, userID).Scan(&exists)
		return exists, err
	}

	err := r.db.QueryRowContext(ctx, query, teamID, userID).Scan(&exists)
	return exists, err
}

func (r *MemberRepositoryImpl) GetCount(ctx context.Context, teamID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM team_members WHERE team_id = $1`
	err := r.db.QueryRowContext(ctx, query, teamID).Scan(&count)
	return count, err
}

func (r *MemberRepositoryImpl) GetMaxMembers(ctx context.Context, teamID string) (int, error) {
	var maxMembers int
	query := `SELECT COALESCE(max_members, 10) FROM teams WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, teamID).Scan(&maxMembers)
	return maxMembers, err
}

func (r *MemberRepositoryImpl) buildMemberListQuery(teamID string, filters team.MemberFilters) (string, []interface{}) {
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

func (r *MemberRepositoryImpl) scanMembers(rows *sql.Rows) ([]team.TeamMember, error) {
	var members []team.TeamMember

	for rows.Next() {
		var m team.TeamMember
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
