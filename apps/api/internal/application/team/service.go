// internal/application/team/service.go
package team

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"seo-backend/internal/domain/team"
	"seo-backend/internal/domain/user"
	"seo-backend/internal/models"
	services "seo-backend/internal/service"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ServiceImpl struct {
	repo         team.Repository
	memberRepo   team.MemberRepository
	queryBuilder *QueryBuilder
	authorizer   team.Authorizer
	validator    team.Validator
	inviteRepo   team.TeamInviteRepository
	userRepo     user.Repository
	emailService services.SMTPEmailService
}

func NewService(
	repo team.Repository,
	memberRepo team.MemberRepository,
	queryBuilder *QueryBuilder,
	authorizer team.Authorizer,
	validator team.Validator,
	inviteRepo team.TeamInviteRepository,
	userRepo user.Repository,
	emailService services.SMTPEmailService,

) team.Service {
	return &ServiceImpl{
		repo:         repo,
		memberRepo:   memberRepo,
		queryBuilder: queryBuilder,
		authorizer:   authorizer,
		validator:    validator,
		userRepo:     userRepo,
		inviteRepo:   inviteRepo,
		emailService: emailService,
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
	log.Printf("ID : , %v", id)
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

	role := req.Role
	if role == "" {
		role = team.RoleMember
	}

	if err := s.memberRepo.Add(ctx, nil, teamID, req.UserID, role); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, teamID, userCtx)
}

// InviteMember - mengirim undangan ke user yang belum terdaftar atau sudah terdaftar
func (s *ServiceImpl) InviteMember(ctx context.Context, teamID string, req team.InviteTeamMemberRequest, userCtx models.UserContext) (*team.TeamInvite, error) {
	// Validate access
	if err := s.authorizer.ValidateTeamAccess(teamID, userCtx); err != nil {
		return nil, err
	}

	// Check if email is already a member of the team
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		// User exists, check if already member
		exists, err := s.memberRepo.Exists(ctx, nil, teamID, existingUser.ID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("user is already a member of this team")
		}
	}

	// Check for existing pending invite
	// existingInvite, err := s.inviteRepo.GetPendingByEmailAndTeam(ctx, req.Email, teamID)
	// if err == nil && existingInvite != nil {
	// 	return nil, fmt.Errorf("pending invitation already exists for this email")
	// }

	// Generate unique token
	token := generateToken()

	// Set expiry (7 days from now)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	// Create invite
	invite := &team.TeamInvite{
		ID:        uuid.New().String(),
		Email:     req.Email,
		TeamID:    teamID,
		Role:      req.Role,
		Token:     token,
		Status:    team.InviteStatusPending,
		InvitedBy: userCtx.GetUserID(),
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Send invitation email
	if err := s.emailService.SendInvitation(ctx, req.Email, token, teamID); err != nil {
		// Log error but don't fail the request
		return nil, fmt.Errorf("Failed to send invitation email: %v\n", err)

	}

	if err := s.inviteRepo.Create(ctx, invite); err != nil {
		return nil, fmt.Errorf("failed to create invitation: %w", err)
	}

	return invite, nil
}

// AcceptInvite - menerima undangan dan langsung menambahkan user ke team
func (s *ServiceImpl) AcceptInvite(
	ctx context.Context,
	userCtx models.UserContext,
	req team.UserInvitedCreate,
) (*team.Team, error) {

	log.Println("========== ACCEPT INVITE ==========")
	log.Printf("Incoming Token : %s\n", req.Token)
	log.Printf("Requester User : %s\n", userCtx.GetUserID())

	// Get invite by token
	invite, err := s.inviteRepo.GetByToken(ctx, req.Token)
	if err != nil {
		log.Printf("FAILED GET INVITE BY TOKEN: %v\n", err)
		return nil, fmt.Errorf("invalid or expired invitation")
	}

	log.Println("SUCCESS GET INVITE")
	log.Printf("Invite ID       : %s\n", invite.ID)
	log.Printf("Invite Email    : %s\n", invite.Email)
	log.Printf("Invite Team ID  : %s\n", invite.TeamID)
	log.Printf("Invite Role     : %s\n", invite.Role)
	log.Printf("Invite Status   : %s\n", invite.Status)
	log.Printf("Invite Expires  : %v\n", invite.ExpiresAt)

	// Check if invite is still pending
	if invite.Status != team.InviteStatusPending {
		log.Printf("INVITE INVALID STATUS: %s\n", invite.Status)
		return nil, fmt.Errorf("invitation is no longer valid")
	}

	// Check if invite has expired
	if time.Now().After(invite.ExpiresAt) {

		log.Println("INVITATION EXPIRED")
		log.Println("UPDATING STATUS TO EXPIRED")

		_ = s.inviteRepo.UpdateStatus(
			ctx,
			invite.ID,
			team.InviteStatusExpired,
		)

		return nil, fmt.Errorf("invitation has expired")
	}

	// Get existing user
	log.Printf("CHECKING USER BY EMAIL: %s\n", invite.Email)

	user, err := s.userRepo.GetByEmail(ctx, invite.Email)

	log.Printf("USER RESULT : %+v\n", user)
	log.Printf("USER ERROR  : %v\n", err)

	// User not found -> create new user
	if err == nil {

		log.Println("USER NOT FOUND, CREATING NEW USER")

		user, err = s.createUserFromInvite(ctx, invite, userCtx, req)
		if err != nil {

			log.Printf("FAILED CREATE USER: %v\n", err)

			return nil, fmt.Errorf(
				"failed to create user: %w",
				err,
			)
		}

		log.Println("SUCCESS CREATE USER")
		log.Printf("NEW USER ID : %s\n", user.ID)
	}

	// Check existing membership
	log.Println("CHECKING MEMBER EXISTENCE")

	exists, err := s.memberRepo.Exists(
		ctx,
		nil,
		invite.TeamID,
		user.ID,
	)

	if err != nil {

		log.Printf("FAILED CHECK MEMBER EXISTENCE: %v\n", err)

		return nil, fmt.Errorf(
			"failed to check membership: %w",
			err,
		)
	}

	log.Printf("MEMBER EXISTS : %v\n", exists)

	if exists {

		log.Println("USER ALREADY MEMBER")
		log.Println("UPDATING INVITE STATUS TO ACCEPTED")

		_ = s.inviteRepo.UpdateStatus(
			ctx,
			invite.ID,
			team.InviteStatusAccepted,
		)

		log.Println("FETCHING TEAM DETAIL")

		return s.GetByID(ctx, invite.TeamID, userCtx)
	}

	// Add member
	role := invite.Role

	log.Println("ADDING MEMBER TO TEAM")
	log.Printf("TEAM ID : %s\n", invite.TeamID)
	log.Printf("USER ID : %s\n", user.ID)
	log.Printf("ROLE    : %s\n", role)

	err = s.memberRepo.Add(
		ctx,
		nil,
		invite.TeamID,
		user.ID,
		role,
	)

	if err != nil {

		log.Printf("FAILED ADD MEMBER: %v\n", err)

		return nil, fmt.Errorf(
			"failed to add member to team: %w",
			err,
		)
	}

	log.Println("SUCCESS ADD MEMBER")

	// Update invite status
	log.Println("UPDATING INVITE STATUS TO ACCEPTED")

	err = s.inviteRepo.UpdateStatus(
		ctx,
		invite.ID,
		team.InviteStatusAccepted,
	)

	if err != nil {

		log.Printf(
			"WARNING FAILED UPDATE INVITE STATUS: %v\n",
			err,
		)

	} else {

		log.Println("SUCCESS UPDATE INVITE STATUS")
	}

	log.Println("FETCHING FINAL TEAM DETAIL")

	data, err := s.GetByID(ctx, invite.TeamID, userCtx)
	if err != nil {

		log.Printf("FAILED GET TEAM DETAIL: %v\n", err)

		return nil, err
	}

	log.Println("========== ACCEPT INVITE SUCCESS ==========")

	return data, nil
}

// VerificationInvite - hanya verifikasi token tanpa auto-join (untuk cek validitas)
func (s *ServiceImpl) VerificationInvite(ctx context.Context, token string) (*team.TeamInvite, error) {
	// Get invite by token
	invite, err := s.inviteRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired invitation")
	}

	// Check if invite is still pending
	if invite.Status != team.InviteStatusPending {
		return nil, fmt.Errorf("invitation is no longer valid")
	}

	// Check if invite has expired
	if time.Now().After(invite.ExpiresAt) {
		// Update status to expired
		_ = s.inviteRepo.UpdateStatus(ctx, invite.ID, team.InviteStatusExpired)
		return nil, fmt.Errorf("invitation has expired")
	}

	return invite, nil
}

// Helper function to create user from invite
func (s *ServiceImpl) createUserFromInvite(
	ctx context.Context,
	invite *team.TeamInvite,
	userCtx models.UserContext,
	req team.UserInvitedCreate,
) (*models.User, error) {

	log.Println("========== CREATE USER FROM INVITE ==========")
	log.Printf("Invite Email      : %s\n", invite.Email)
	log.Printf("Invite Team ID    : %s\n", invite.TeamID)
	log.Printf("Invite Role       : %s\n", invite.Role)
	log.Printf("Invite Token      : %s\n", invite.Token)
	log.Printf("Triggered By User : %s\n", userCtx.GetUserID())

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		log.Printf("FAILED HASH PASSWORD: %v\n", err)
		return nil, err
	}

	log.Println("Password hashing success")

	user := &user.CreateUserRequest{
		Email:    invite.Email,
		Password: req.Password,
		Name:     extractNameFromEmail(invite.Email),
	}

	log.Println("Prepared create user payload")
	log.Printf("User Email : %s\n", user.Email)
	log.Printf("User Name  : %s\n", user.Name)

	data, err := s.userRepo.Create(ctx, *user, hashedPassword)
	if err != nil {
		log.Printf("FAILED CREATE USER: %v\n", err)
		return nil, err
	}

	log.Println("SUCCESS CREATE USER")
	log.Printf("Created User ID : %s\n", data.ID)

	err = s.emailService.SendWelcomeEmail(
		ctx,
		invite.Email,
		req.Password,
	)

	if err != nil {
		log.Printf("FAILED SEND WELCOME EMAIL: %v\n", err)
	} else {
		log.Println("SUCCESS SEND WELCOME EMAIL")
	}

	log.Println("========== FINISH CREATE USER FROM INVITE ==========")

	return data, nil
}

// Helper functions
func generateRandomPassword() string {
	// Generate random 12 character password
	const charset = "newmember"
	return string(charset)
}

func extractNameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// AcceptInvite - menerima undangan (auto-join)

// CancelInvite - membatalkan undangan
func (s *ServiceImpl) CancelInvite(ctx context.Context, inviteID string, userCtx models.UserContext) error {
	// Get invite
	invite, err := s.inviteRepo.GetByID(ctx, inviteID)
	if err != nil {
		return fmt.Errorf("invitation not found")
	}

	// Validate access to the team
	if err := s.authorizer.ValidateTeamAccess(invite.TeamID, userCtx); err != nil {
		return err
	}

	// Only pending invites can be cancelled
	if invite.Status != team.InviteStatusPending {
		return fmt.Errorf("only pending invitations can be cancelled")
	}

	return s.inviteRepo.UpdateStatus(ctx, inviteID, team.InviteStatusCancelled)
}

// ResendInvite - mengirim ulang undangan
func (s *ServiceImpl) ResendInvite(ctx context.Context, inviteID string, userCtx models.UserContext) (*team.TeamInvite, error) {
	// Get existing invite
	invite, err := s.inviteRepo.GetByID(ctx, inviteID)
	if err != nil {
		return nil, fmt.Errorf("invitation not found")
	}

	// Validate access
	if err := s.authorizer.ValidateTeamAccess(invite.TeamID, userCtx); err != nil {
		return nil, err
	}

	// Generate new token and extend expiry
	newToken := generateToken()
	newExpiry := time.Now().Add(7 * 24 * time.Hour)

	// Update invite
	invite.Token = newToken
	invite.ExpiresAt = newExpiry
	invite.Status = team.InviteStatusPending
	invite.UpdatedAt = time.Now()

	if err := s.inviteRepo.Update(ctx, invite); err != nil {
		return nil, fmt.Errorf("failed to resend invitation: %w", err)
	}

	// Resend email
	if err := s.emailService.SendInvitation(ctx, invite.Email, newToken, invite.TeamID); err != nil {
		fmt.Printf("Failed to resend invitation email: %v\n", err)
	}

	return invite, nil
}

// GetTeamInvites - mendapatkan semua undangan untuk sebuah team
func (s *ServiceImpl) GetTeamInvites(ctx context.Context, teamID string, userCtx models.UserContext) ([]team.TeamInvite, error) {
	if err := s.authorizer.ValidateTeamAccess(teamID, userCtx); err != nil {
		return nil, err
	}

	return s.inviteRepo.GetByTeamID(ctx, teamID)
}

// AutoJoinOnRegister - auto-join saat user register (Hybrid approach)
func (s *ServiceImpl) AutoJoinOnRegister(ctx context.Context, email string, userID string) error {
	// Get all pending invites for this email
	invites, err := s.inviteRepo.GetPendingByEmail(ctx, email)
	if err != nil {
		return err
	}

	if len(invites) == 0 {
		return nil
	}

	for _, invite := range invites {
		// Check if invite is expired
		if time.Now().After(invite.ExpiresAt) {
			// Update to expired status
			_ = s.inviteRepo.UpdateStatus(ctx, invite.ID, team.InviteStatusExpired)
			continue
		}

		// Check if user is already a member
		exists, err := s.memberRepo.Exists(ctx, nil, invite.TeamID, userID)
		if err != nil {
			// Log error but continue with other invites
			fmt.Printf("Failed to check member existence: %v\n", err)
			continue
		}

		if !exists {
			// Add to team
			if err := s.memberRepo.Add(ctx, nil, invite.TeamID, userID, team.TeamMemberRole(invite.Role)); err != nil {
				fmt.Printf("Failed to auto-join team %s: %v\n", invite.TeamID, err)
				continue
			}

			// Update invite status
			_ = s.inviteRepo.UpdateStatus(ctx, invite.ID, team.InviteStatusAccepted)
		}
	}

	return nil
}

// Helper function to generate unique token
func generateToken() string {
	return fmt.Sprintf("%s-%d", uuid.New().String(), time.Now().UnixNano())
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
