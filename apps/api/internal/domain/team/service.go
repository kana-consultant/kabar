// internal/domain/team/service.go
package team

import (
	"context"
	"seo-backend/internal/models"
)

type Service interface {
	// Team CRUD
	GetAll(ctx context.Context, userCtx models.UserContext, filters TeamFilters) ([]Team, error)
	GetByID(ctx context.Context, id string, userCtx models.UserContext) (*Team, error)
	Create(ctx context.Context, req CreateTeamRequest, createdBy string) (*Team, error)
	Update(ctx context.Context, id string, updates map[string]interface{}, userCtx models.UserContext) error
	Delete(ctx context.Context, id string, userCtx models.UserContext) error

	// Member management
	GetTeamMembers(ctx context.Context, teamID string, filters MemberFilters, userCtx models.UserContext) ([]TeamMember, error)
	AddMember(ctx context.Context, teamID string, req AddTeamMemberRequest, userCtx models.UserContext) (*Team, error)
	UpdateMemberRole(ctx context.Context, teamID, userID string, role TeamMemberRole, userCtx models.UserContext) (*Team, error)
	RemoveMember(ctx context.Context, teamID, userID string, userCtx models.UserContext) (*Team, error)

	// User teams
	GetUserTeams(ctx context.Context, userID string) ([]Team, error)
}

type Authorizer interface {
	CanAccessTeam(teamID string, ctx models.UserContext) bool
	CanManageTeam(teamID string, ctx models.UserContext, userRole string) bool
	ValidateTeamAccess(teamID string, ctx models.UserContext) error
}

type Validator interface {
	CheckMemberLimit(ctx context.Context, teamID string) error
	CheckDeleteTeam(ctx context.Context, teamID string) error
}
