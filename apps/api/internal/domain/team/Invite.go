package team

import (
	"context"
	"time"
)

type TeamInviteRepository interface {
	Create(ctx context.Context, invite *TeamInvite) error
	GetByID(ctx context.Context, id string) (*TeamInvite, error)
	GetByToken(ctx context.Context, token string) (*TeamInvite, error)
	GetPendingByEmailAndTeam(ctx context.Context, email, teamID string) (*TeamInvite, error)
	GetPendingByEmail(ctx context.Context, email string) ([]TeamInvite, error)
	GetByTeamID(ctx context.Context, teamID string) ([]TeamInvite, error)
	Update(ctx context.Context, invite *TeamInvite) error
	UpdateStatus(ctx context.Context, id string, status string) error
	Delete(ctx context.Context, id string) error
}

// team/models.go
// models.go
type TeamInvite struct {
	ID        string         `db:"id" json:"id"`
	Email     string         `db:"email" json:"email"`
	TeamID    string         `db:"team_id" json:"team_id"`
	TeamName  string         `db:"team_name" json:"team_name"` // Tambah field ini
	Role      TeamMemberRole `db:"role" json:"role"`
	Token     string         `db:"token" json:"token"`
	Status    string         `db:"status" json:"status"`
	InvitedBy string         `db:"invited_by" json:"invited_by"`
	ExpiresAt time.Time      `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at"`
}

type InviteTeamMemberRequest struct {
	Email string         `json:"email" validate:"required,email"`
	Role  TeamMemberRole `json:"role" validate:"required,oneof=admin member viewer"`
}

const (
	InviteStatusPending   = "pending"
	InviteStatusAccepted  = "accepted"
	InviteStatusExpired   = "expired"
	InviteStatusCancelled = "cancelled"
)
