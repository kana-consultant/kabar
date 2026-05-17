// internal/domain/team/entity.go
package team

import (
	"time"
)

type Team struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Logo        *string      `json:"logo"`
	Status      string       `json:"status"`
	MaxMembers  int          `json:"max_members"`
	CreatedBy   string       `json:"created_by"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Members     []TeamMember `json:"members,omitempty"`
}

type TeamMember struct {
	ID       string         `json:"id"`
	TeamID   string         `json:"team_id"`
	UserID   string         `json:"user_id"`
	Role     TeamMemberRole `json:"role"`
	JoinedAt time.Time      `json:"joined_at"`
}

type TeamMemberRole string

const (
	RoleManager TeamMemberRole = "manager"
	RoleAdmin   TeamMemberRole = "admin"
	RoleMember  TeamMemberRole = "member"
	RoleViewer  TeamMemberRole = "viewer"
)

type TeamFilters struct {
	Status string
}

type MemberFilters struct {
	Role string
}

type CreateTeamRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UserInvitedCreate struct {
	Password string `json:"password"`
	Token    string `json:"token"`
}

type AddTeamMemberRequest struct {
	UserID string         `json:"user_id"`
	Role   TeamMemberRole `json:"role"`
}

type UserContext interface {
	IsAdmin() bool
	GetTeamID() string
	GetUserID() string
	GetUserRole() string
}
