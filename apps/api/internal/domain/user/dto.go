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
	TeamID  string `json:"team_id"` // Tambahkan field ini
	Search  string `json:"search"`
	Role    string `json:"role"`
	Status  string `json:"status"`
	OrderBy string `json:"order_by"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
}
