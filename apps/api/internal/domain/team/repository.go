// internal/domain/team/repository.go
package team

import (
	"context"
	"database/sql"
)

type Repository interface {
	// Team operations
	GetByID(ctx context.Context, id string) (*Team, error)
	GetAll(ctx context.Context, query string, args []interface{}) ([]Team, error)
	Insert(ctx context.Context, req CreateTeamRequest, createdBy string) (string, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	GetUserTeams(ctx context.Context, userID string) ([]Team, error)
}

type MemberRepository interface {
	GetByTeamID(ctx context.Context, teamID string, filters MemberFilters) ([]TeamMember, error)
	Add(ctx context.Context, tx *sql.Tx, teamID, userID string, role TeamMemberRole) error
	UpdateRole(ctx context.Context, teamID, userID string, role TeamMemberRole) error
	Remove(ctx context.Context, teamID, userID string) error
	Exists(ctx context.Context, tx *sql.Tx, teamID, userID string) (bool, error)
	GetCount(ctx context.Context, teamID string) (int, error)
	GetMaxMembers(ctx context.Context, teamID string) (int, error)
}
