package models

import (
	"time"
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleManager UserRole = "manager"
	RoleEditor  UserRole = "editor"
	RoleViewer  UserRole = "viewer"
)

type UserStatus string

const (
	StatusActive   UserStatus = "active"
	StatusInactive UserStatus = "inactive"
)

type User struct {
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	Name         string     `json:"name"`
	PasswordHash string     `json:"-"`
	Role         UserRole   `json:"role"`
	Avatar       *string    `json:"avatar,omitempty"`
	Status       UserStatus `json:"status"`
	LastActive   *time.Time `json:"lastActive,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

func isZeroUUID(id string) bool {
	return id == "" || id == "00000000-0000-0000-0000-000000000000"
}

func nullIfEmpty(id string) interface{} {
	if id == "" || isZeroUUID(id) {
		return nil
	}
	return id
}

type UserContext interface {
	GetUserID() string
	GetTeamID() string
	GetRole() string
	IsAdmin() bool
}

type TeamFilters struct {
	Status string
}

type MemberFilters struct {
	Role string
}

type SimpleUserContext struct {
	UserID string
	TeamID string
	Role   string
}

func (c *SimpleUserContext) GetUserID() string { return c.UserID }
func (c *SimpleUserContext) GetTeamID() string { return c.TeamID }
func (c *SimpleUserContext) GetRole() string   { return c.Role }
func (c *SimpleUserContext) IsAdmin() bool     { return c.Role == "admin" || c.Role == "super_admin" }
