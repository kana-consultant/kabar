package user

import "fmt"

type Validator struct {
	repo *Repository
}

func NewValidator(repo *Repository) *Validator {
	return &Validator{repo: repo}
}

func (v *Validator) ValidateCreate(req CreateUserRequest) error {
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}
	if len(req.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	return nil
}

func (v *Validator) ValidateEmailUniqueness(email string) error {
	exists, err := v.repo.EmailExists(email)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("email already exists")
	}
	return nil
}

func (v *Validator) ValidateUpdateEmailUniqueness(email, currentEmail string) error {
	if email != currentEmail {
		return v.ValidateEmailUniqueness(email)
	}
	return nil
}
