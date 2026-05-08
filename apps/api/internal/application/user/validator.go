package user

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"seo-backend/internal/domain/user"
)

type Validator struct {
	repo user.Repository
}

func NewValidator(repo user.Repository) *Validator {
	return &Validator{repo: repo}
}

// ValidateCreate validates user creation request
func (v *Validator) ValidateCreate(req user.CreateUserRequest) error {
	if req.Email == "" {
		return errors.New("email is required")
	}

	if !v.isValidEmail(req.Email) {
		return errors.New("invalid email format")
	}

	if req.Name == "" {
		return errors.New("name is required")
	}

	if len(req.Name) < 2 {
		return errors.New("name must be at least 2 characters")
	}

	if len(req.Name) > 100 {
		return errors.New("name must not exceed 100 characters")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	if req.Role != "" {
		validRoles := map[string]bool{
			"admin": true, "viewer": true, "editor": true, "manager": true,
		}
		if !validRoles[req.Role] {
			return errors.New("invalid role")
		}
	}

	return nil
}

// ValidateUpdate validates user update request
func (v *Validator) ValidateUpdate(req user.UpdateUserRequest) error {
	if req.Email != nil && *req.Email != "" {
		if !v.isValidEmail(*req.Email) {
			return errors.New("invalid email format")
		}
	}

	if req.Name != nil && *req.Name != "" {
		if len(*req.Name) < 2 {
			return errors.New("name must be at least 2 characters")
		}
		if len(*req.Name) > 100 {
			return errors.New("name must not exceed 100 characters")
		}
	}

	if req.Role != nil && *req.Role != "" {
		validRoles := map[string]bool{
			"admin": true, "viewer": true, "editor": true, "manager": true,
		}
		if !validRoles[*req.Role] {
			return errors.New("invalid role")
		}
	}

	if req.Status != nil && *req.Status != "" {
		validStatuses := map[string]bool{
			"active": true, "inactive": true, "suspended": true,
		}
		if !validStatuses[*req.Status] {
			return errors.New("invalid status")
		}
	}

	return nil
}

// ValidateEmailUniqueness checks if email is unique
func (v *Validator) ValidateEmailUniqueness(ctx context.Context, email string) error {
	exists, err := v.repo.EmailExists(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already exists")
	}
	return nil
}

// ValidateUpdateEmailUniqueness checks if email is unique for update
func (v *Validator) ValidateUpdateEmailUniqueness(ctx context.Context, newEmail, currentEmail string) error {
	if strings.EqualFold(newEmail, currentEmail) {
		return nil // Same email, no need to check
	}
	return v.ValidateEmailUniqueness(ctx, newEmail)
}

// isValidEmail validates email format
func (v *Validator) isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
