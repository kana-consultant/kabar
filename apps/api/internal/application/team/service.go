// internal/application/team/service.go
package team

import (
	"context"
	"fmt"
	"log"

	"seo-backend/internal/domain/team"
	"seo-backend/internal/models"
)

type ServiceImpl struct {
	repo         team.Repository
	memberRepo   team.MemberRepository
	queryBuilder *QueryBuilder
	authorizer   team.Authorizer
	validator    team.Validator
}

func NewService(
	repo team.Repository,
	memberRepo team.MemberRepository,
	queryBuilder *QueryBuilder,
	authorizer team.Authorizer,
	validator team.Validator,
) team.Service {
	return &ServiceImpl{
		repo:         repo,
		memberRepo:   memberRepo,
		queryBuilder: queryBuilder,
		authorizer:   authorizer,
		validator:    validator,
	}
}

// GetAll implements team.Service
func (s *ServiceImpl) GetAll(ctx context.Context, userCtx models.UserContext, filters team.TeamFilters) ([]team.Team, error) {
	query, args := s.queryBuilder.BuildListQuery(userCtx, filters)
	teams, err := s.repo.GetAll(ctx, query, args)
	if err != nil {
		return nil, err
	}

	// Load members for each team
	for i := range teams {
		members, _ := s.memberRepo.GetByTeamID(ctx, teams[i].ID, team.MemberFilters{})
		teams[i].Members = members
	}

	return teams, nil
}

// GetByID implements team.Service
func (s *ServiceImpl) GetByID(ctx context.Context, id string, userCtx models.UserContext) (*team.Team, error) {
	if err := s.authorizer.ValidateTeamAccess(id, userCtx); err != nil {
		return nil, err
	}

	teamData, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if teamData == nil {
		return nil, nil
	}

	// Load members
	members, err := s.memberRepo.GetByTeamID(ctx, id, team.MemberFilters{})
	if err != nil {
		log.Printf("Warning: failed to load team members: %v", err)
	}
	teamData.Members = members

	return teamData, nil
}

// Create implements team.Service
func (s *ServiceImpl) Create(ctx context.Context, req team.CreateTeamRequest, createdBy string) (*team.Team, error) {
	if createdBy == "" {
		createdBy = "system"
	}

	teamID, err := s.repo.Insert(ctx, req, createdBy)
	if err != nil {
		return nil, err
	}

	// Auto-add creator as team member
	if err := s.memberRepo.Add(ctx, nil, teamID, createdBy, team.RoleViewer); err != nil {
		log.Printf("Warning: failed to add creator as member: %v", err)
	}

	return s.repo.GetByID(ctx, teamID)
}

// Update implements team.Service
func (s *ServiceImpl) Update(ctx context.Context, id string, updates map[string]interface{}, userCtx models.UserContext) error {
	if err := s.authorizer.ValidateTeamAccess(id, userCtx); err != nil {
		return err
	}

	teamData, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if teamData == nil {
		return fmt.Errorf("team not found")
	}

	return s.repo.Update(ctx, id, updates)
}

// Delete implements team.Service
func (s *ServiceImpl) Delete(ctx context.Context, id string, userCtx models.UserContext) error {
	if err := s.authorizer.ValidateTeamAccess(id, userCtx); err != nil {
		return err
	}

	teamData, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if teamData == nil {
		return fmt.Errorf("team not found")
	}

	if err := s.validator.CheckDeleteTeam(ctx, id); err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}

// GetTeamMembers implements team.Service
func (s *ServiceImpl) GetTeamMembers(ctx context.Context, teamID string, filters team.MemberFilters, userCtx models.UserContext) ([]team.TeamMember, error) {
	if err := s.authorizer.ValidateTeamAccess(teamID, userCtx); err != nil {
		return nil, err
	}
	return s.memberRepo.GetByTeamID(ctx, teamID, filters)
}

// AddMember implements team.Service
func (s *ServiceImpl) AddMember(ctx context.Context, teamID string, req team.AddTeamMemberRequest, userCtx models.UserContext) (*team.Team, error) {
	if err := s.authorizer.ValidateTeamAccess(teamID, userCtx); err != nil {
		return nil, err
	}

	// Check if member already exists
	exists, err := s.memberRepo.Exists(ctx, nil, teamID, req.UserID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("member already in team")
	}

	// Check member limit
	if err := s.validator.CheckMemberLimit(ctx, teamID); err != nil {
		return nil, err
	}

	role := req.Role
	if role == "" {
		role = team.RoleMember
	}

	if err := s.memberRepo.Add(ctx, nil, teamID, req.UserID, role); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, teamID, userCtx)
}

// UpdateMemberRole implements team.Service
func (s *ServiceImpl) UpdateMemberRole(ctx context.Context, teamID, userID string, role team.TeamMemberRole, userCtx models.UserContext) (*team.Team, error) {
	if err := s.authorizer.ValidateTeamAccess(teamID, userCtx); err != nil {
		return nil, err
	}

	if err := s.memberRepo.UpdateRole(ctx, teamID, userID, role); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, teamID, userCtx)
}

// RemoveMember implements team.Service
func (s *ServiceImpl) RemoveMember(ctx context.Context, teamID, userID string, userCtx models.UserContext) (*team.Team, error) {
	if err := s.authorizer.ValidateTeamAccess(teamID, userCtx); err != nil {
		return nil, err
	}

	if err := s.memberRepo.Remove(ctx, teamID, userID); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, teamID, userCtx)
}

// GetUserTeams implements team.Service
func (s *ServiceImpl) GetUserTeams(ctx context.Context, userID string) ([]team.Team, error) {
	return s.repo.GetUserTeams(ctx, userID)
}
