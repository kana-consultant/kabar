// repositories/team_invite_repository.go

package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"seo-backend/internal/domain/team"
	"time"
)

type TeamInviteRepository struct {
	db *sql.DB
}

func NewTeamInviteRepository(db *sql.DB) *TeamInviteRepository {
	return &TeamInviteRepository{db: db}
}

// Create - membuat invite baru
func (r *TeamInviteRepository) Create(ctx context.Context, invite *team.TeamInvite) error {
	query := `
        INSERT INTO team_invites (
            id, email, team_id, role, token, status, invited_by, expires_at, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `

	_, err := r.db.ExecContext(ctx, query,
		invite.ID,
		invite.Email,
		invite.TeamID,
		invite.Role,
		invite.Token,
		invite.Status,
		invite.InvitedBy,
		invite.ExpiresAt,
		invite.CreatedAt,
		invite.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create team invite: %w", err)
	}

	return nil
}

// GetByID - mendapatkan invite berdasarkan ID
func (r *TeamInviteRepository) GetByID(ctx context.Context, id string) (*team.TeamInvite, error) {
	query := `
        SELECT id, email, team_id, role, token, status, invited_by, expires_at, created_at, updated_at
        FROM team_invites
        WHERE id = $1
    `

	var invite team.TeamInvite
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&invite.ID,
		&invite.Email,
		&invite.TeamID,
		&invite.Role,
		&invite.Token,
		&invite.Status,
		&invite.InvitedBy,
		&invite.ExpiresAt,
		&invite.CreatedAt,
		&invite.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("team invite not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get team invite: %w", err)
	}

	return &invite, nil
}

// GetByToken - mendapatkan invite berdasarkan token
func (r *TeamInviteRepository) GetByToken(ctx context.Context, token string) (*team.TeamInvite, error) {
	query := `
        SELECT id, email, team_id, role, token, status, invited_by, expires_at, created_at, updated_at
        FROM team_invites
        WHERE token = $1
    `

	var invite team.TeamInvite
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&invite.ID,
		&invite.Email,
		&invite.TeamID,
		&invite.Role,
		&invite.Token,
		&invite.Status,
		&invite.InvitedBy,
		&invite.ExpiresAt,
		&invite.CreatedAt,
		&invite.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("team invite not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get team invite by token: %w", err)
	}

	return &invite, nil
}

// GetPendingByEmailAndTeam - mendapatkan pending invite berdasarkan email dan team
func (r *TeamInviteRepository) GetPendingByEmailAndTeam(ctx context.Context, email, teamID string) (*team.TeamInvite, error) {
	query := `
        SELECT id, email, team_id, role, token, status, invited_by, expires_at, created_at, updated_at
        FROM team_invites
        WHERE email = $1 AND team_id = $2 AND status = 'pending'
        ORDER BY created_at DESC
        LIMIT 1
    `

	var invite team.TeamInvite
	err := r.db.QueryRowContext(ctx, query, email, teamID).Scan(
		&invite.ID,
		&invite.Email,
		&invite.TeamID,
		&invite.Role,
		&invite.Token,
		&invite.Status,
		&invite.InvitedBy,
		&invite.ExpiresAt,
		&invite.CreatedAt,
		&invite.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Tidak ada invite pending
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get pending invite: %w", err)
	}

	return &invite, nil
}

// GetPendingByEmail - mendapatkan semua pending invite berdasarkan email
func (r *TeamInviteRepository) GetPendingByEmail(ctx context.Context, email string) ([]team.TeamInvite, error) {
	query := `
        SELECT id, email, team_id, role, token, status, invited_by, expires_at, created_at, updated_at
        FROM team_invites
        WHERE email = $1 AND status = 'pending' AND expires_at > $2
        ORDER BY created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, query, email, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to get pending invites: %w", err)
	}
	defer rows.Close()

	var invites []team.TeamInvite
	for rows.Next() {
		var invite team.TeamInvite
		err := rows.Scan(
			&invite.ID,
			&invite.Email,
			&invite.TeamID,
			&invite.Role,
			&invite.Token,
			&invite.Status,
			&invite.InvitedBy,
			&invite.ExpiresAt,
			&invite.CreatedAt,
			&invite.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan invite: %w", err)
		}
		invites = append(invites, invite)
	}

	return invites, nil
}

// GetByTeamID - mendapatkan semua invite untuk sebuah team
func (r *TeamInviteRepository) GetByTeamID(ctx context.Context, teamID string) ([]team.TeamInvite, error) {
	query := `
        SELECT id, email, team_id, role, token, status, invited_by, expires_at, created_at, updated_at
        FROM team_invites
        WHERE team_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, query, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team invites: %w", err)
	}
	defer rows.Close()

	var invites []team.TeamInvite
	for rows.Next() {
		var invite team.TeamInvite
		err := rows.Scan(
			&invite.ID,
			&invite.Email,
			&invite.TeamID,
			&invite.Role,
			&invite.Token,
			&invite.Status,
			&invite.InvitedBy,
			&invite.ExpiresAt,
			&invite.CreatedAt,
			&invite.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan invite: %w", err)
		}
		invites = append(invites, invite)
	}

	return invites, nil
}

// Update - mengupdate data invite
func (r *TeamInviteRepository) Update(ctx context.Context, invite *team.TeamInvite) error {
	query := `
        UPDATE team_invites 
        SET email = $1, 
            team_id = $2, 
            role = $3, 
            token = $4, 
            status = $5, 
            invited_by = $6, 
            expires_at = $7, 
            updated_at = $8
        WHERE id = $9
    `

	result, err := r.db.ExecContext(ctx, query,
		invite.Email,
		invite.TeamID,
		invite.Role,
		invite.Token,
		invite.Status,
		invite.InvitedBy,
		invite.ExpiresAt,
		invite.UpdatedAt,
		invite.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update team invite: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("team invite not found")
	}

	return nil
}

// UpdateStatus - mengupdate status invite saja
func (r *TeamInviteRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `
        UPDATE team_invites 
        SET status = $1, updated_at = $2
        WHERE id = $3
    `

	result, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update invite status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("team invite not found")
	}

	return nil
}

// Delete - menghapus invite
func (r *TeamInviteRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM team_invites WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete team invite: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("team invite not found")
	}

	return nil
}
