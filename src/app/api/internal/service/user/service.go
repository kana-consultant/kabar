package user

import (
	"database/sql"
	"fmt"

	"seo-backend/internal/database"
	"seo-backend/internal/models"
)

type Service struct {
	db           *sql.DB
	repo         *Repository
	queryBuilder *QueryBuilder
	auth         *Authorizer
	password     *PasswordService
	validator    *Validator
}

func NewService() *Service {
	db := database.GetDB()
	repo := NewRepository(db)

	return &Service{
		db:           db,
		repo:         repo,
		queryBuilder: NewQueryBuilder(),
		auth:         NewAuthorizer(db),
		password:     NewPasswordService(),
		validator:    NewValidator(repo),
	}
}

// User CRUD operations
func (s *Service) GetAll(ctx models.UserContext, filters UserFilters) ([]models.User, error) {
	query, args := s.queryBuilder.BuildListQuery(ctx, filters)
	return s.repo.GetAll(query, args)
}

func (s *Service) GetByID(id string, ctx models.UserContext) (*models.User, error) {
	if err := s.auth.ValidateAccess(id, ctx); err != nil {
		return nil, err
	}
	return s.repo.GetByID(id)
}

func (s *Service) GetCurrentUser(ctx models.UserContext) (*models.User, error) {
	if ctx.GetUserID() == "" {
		return nil, fmt.Errorf("unauthorized")
	}
	return s.repo.GetByID(ctx.GetUserID())
}

func (s *Service) Create(req CreateUserRequest) (*models.User, error) {
	// Validate request
	if err := s.validator.ValidateCreate(req); err != nil {
		return nil, err
	}

	// Check email uniqueness
	if err := s.validator.ValidateEmailUniqueness(req.Email); err != nil {
		return nil, err
	}

	// Hash password
	passwordHash, err := s.password.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	return s.repo.Create(req, passwordHash)
}

func (s *Service) Update(id string, req UpdateUserRequest, ctx models.UserContext) error {
	// Check access
	if err := s.auth.ValidateAccess(id, ctx); err != nil {
		return err
	}

	// Get existing user
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Build updates
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Email != nil {
		// Check email uniqueness if changing
		if err := s.validator.ValidateUpdateEmailUniqueness(*req.Email, user.Email); err != nil {
			return err
		}
		updates["email"] = *req.Email
	}
	if req.Role != nil {
		updates["role"] = *req.Role
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	return s.repo.Update(id, updates)
}

func (s *Service) Delete(id string, ctx models.UserContext) error {
	// Check access (cannot delete self)
	if !s.auth.CanDelete(id, ctx) {
		return fmt.Errorf("cannot delete this user")
	}

	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	return s.repo.Delete(id)
}

// Helper method for auth (used by other services)
func (s *Service) CanAccessUser(targetUserID string, ctx models.UserContext) bool {
	return s.auth.CanAccess(targetUserID, ctx)
}
