package team

import (
	"database/sql"
	"fmt"
	"log"

	"seo-backend/internal/database"
	"seo-backend/internal/models"
)

type Service struct {
	db           *sql.DB
	repo         *Repository
	memberRepo   *MemberRepository
	queryBuilder *QueryBuilder
	auth         *Authorizer
	validator    *Validator
}

func NewService() *Service {
	db := database.GetDB()
	memberRepo := NewMemberRepository(db)

	return &Service{
		db:           db,
		repo:         NewRepository(db),
		memberRepo:   memberRepo,
		queryBuilder: NewQueryBuilder(),
		auth:         NewAuthorizer(),
		validator:    NewValidator(memberRepo),
	}
}

// Team CRUD operations
func (s *Service) GetAll(ctx models.UserContext, filters TeamFilters) ([]models.Team, error) {
	query, args := s.queryBuilder.BuildListQuery(ctx, filters)
	teams, err := s.repo.GetAll(query, args)
	if err != nil {
		return nil, err
	}

	// Load members for each team
	for i := range teams {
		members, _ := s.memberRepo.GetByTeamID(teams[i].ID, MemberFilters{})
		teams[i].Members = members
	}

	return teams, nil
}

func (s *Service) GetByID(id string, ctx models.UserContext) (*models.Team, error) {
	if err := s.auth.ValidateTeamAccess(id, ctx); err != nil {
		return nil, err
	}

	team, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, nil
	}

	// Load members
	members, err := s.memberRepo.GetByTeamID(id, MemberFilters{})
	if err != nil {
		log.Printf("Warning: failed to load team members: %v", err)
	}
	team.Members = members

	return team, nil
}

func (s *Service) Create(req models.CreateTeamRequest, createdBy string) (*models.Team, error) {
	if createdBy == "" {
		createdBy = "system"
	}

	teamID, err := s.repo.Insert(req, createdBy)
	if err != nil {
		return nil, err
	}

	// Auto-add creator as team member
	if err := s.memberRepo.Add(nil, teamID, createdBy, "viewer"); err != nil {
		log.Printf("Warning: failed to add creator as member: %v", err)
	}

	return s.repo.GetByID(teamID)
}

func (s *Service) Update(id string, updates map[string]interface{}, ctx models.UserContext) error {
	if err := s.auth.ValidateTeamAccess(id, ctx); err != nil {
		return err
	}

	team, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if team == nil {
		return fmt.Errorf("team not found")
	}

	return s.repo.Update(id, updates)
}

func (s *Service) Delete(id string, ctx models.UserContext) error {
	if err := s.auth.ValidateTeamAccess(id, ctx); err != nil {
		return err
	}

	team, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if team == nil {
		return fmt.Errorf("team not found")
	}

	if err := s.validator.CheckDeleteTeam(id); err != nil {
		return err
	}

	return s.repo.Delete(id)
}

// Member management
func (s *Service) GetTeamMembers(teamID string, filters MemberFilters, ctx models.UserContext) ([]models.TeamMember, error) {
	if err := s.auth.ValidateTeamAccess(teamID, ctx); err != nil {
		return nil, err
	}
	return s.memberRepo.GetByTeamID(teamID, filters)
}

func (s *Service) AddMember(teamID string, req models.AddTeamMemberRequest, ctx models.UserContext) (*models.Team, error) {
	if err := s.auth.ValidateTeamAccess(teamID, ctx); err != nil {
		return nil, err
	}

	// Check if member already exists
	exists, err := s.memberRepo.Exists(nil, teamID, req.UserID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("member already in team")
	}

	// Check member limit
	if err := s.validator.CheckMemberLimit(teamID); err != nil {
		return nil, err
	}

	role := req.Role
	if role == "" {
		role = "member"
	}

	if err := s.memberRepo.Add(nil, teamID, req.UserID, role); err != nil {
		return nil, err
	}

	return s.GetByID(teamID, ctx)
}

func (s *Service) UpdateMemberRole(teamID, userID string, role models.TeamMemberRole, ctx models.UserContext) (*models.Team, error) {
	if err := s.auth.ValidateTeamAccess(teamID, ctx); err != nil {
		return nil, err
	}

	if err := s.memberRepo.UpdateRole(teamID, userID, role); err != nil {
		return nil, err
	}

	return s.GetByID(teamID, ctx)
}

func (s *Service) RemoveMember(teamID, userID string, ctx models.UserContext) (*models.Team, error) {
	if err := s.auth.ValidateTeamAccess(teamID, ctx); err != nil {
		return nil, err
	}

	if err := s.memberRepo.Remove(teamID, userID); err != nil {
		return nil, err
	}

	return s.GetByID(teamID, ctx)
}

func (s *Service) GetUserTeams(userID string) ([]models.Team, error) {
	return s.repo.GetUserTeams(userID)
}
