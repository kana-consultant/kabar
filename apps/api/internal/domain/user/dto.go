package user

// CreateUserRequest DTO for user creation
type CreateUserRequest struct {
	Email    string
	Name     string
	Password string
	Role     string
}

// UpdateUserRequest DTO for user update
type UpdateUserRequest struct {
	Email  *string
	Name   *string
	Role   *string
	Status *string
	Avatar *string
}

// UserFilters for filtering users
type UserFilters struct {
	Search  string
	Role    string
	Status  string
	OrderBy string
	Limit   int
	Offset  int
}
