package user

import "seo-backend/internal/models"

type UserFilters struct {
	Role   string
	Status string
	Search string
}

type CreateUserRequest struct {
	Email    string
	Name     string
	Password string
	Role     models.UserRole
}

type UpdateUserRequest struct {
	Name   *string
	Email  *string
	Role   *models.UserRole
	Status *string
}
