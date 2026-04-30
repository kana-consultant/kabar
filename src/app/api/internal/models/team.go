package models

import (
	"time"
)

type TeamStatus string

const (
	TeamStatusActive   TeamStatus = "active"
	TeamStatusInactive TeamStatus = "inactive"
	TeamStatusArchived TeamStatus = "archived"
)

type Team struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	Logo        *string    `json:"logo,omitempty"`
	Status      TeamStatus `json:"status"`
	MaxMembers  int        `json:"maxMembers"`
	CreatedBy   *string    `json:"createdBy,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	Members     []TeamMember `json:"members,omitempty"`
}

type TeamMemberRole string

const (
	TeamMemberRoleManager TeamMemberRole = "manager"
	TeamMemberRoleEditor  TeamMemberRole = "editor"
	TeamMemberRoleViewer  TeamMemberRole = "viewer"
)

type TeamMemberStatus string

const (
	TeamMemberStatusActive  TeamMemberStatus = "active"
	TeamMemberStatusInactive TeamMemberStatus = "inactive"
	TeamMemberStatusPending TeamMemberStatus = "pending"
	TeamMemberStatusLeft    TeamMemberStatus = "left"
)

type TeamMember struct {
	ID         string            `json:"id"`
	TeamID     string            `json:"teamId"`
	UserID     string            `json:"userId"`
	UserEmail  string            `json:"userEmail"`
	UserName   string            `json:"userName"`
	UserAvatar *string           `json:"userAvatar,omitempty"`
	Role       TeamMemberRole    `json:"role"`
	Status     TeamMemberStatus  `json:"status"`
	JoinedAt   time.Time         `json:"joinedAt"`
	LeftAt     *time.Time        `json:"leftAt,omitempty"`
}

type CreateTeamRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Logo        *string `json:"logo,omitempty"`
}

type UpdateTeamRequest struct {
	Name        *string     `json:"name,omitempty"`
	Description *string     `json:"description,omitempty"`
	Logo        *string     `json:"logo,omitempty"`
	Status      *TeamStatus `json:"status,omitempty"`
	MaxMembers  *int        `json:"maxMembers,omitempty"`
}

type AddTeamMemberRequest struct {
	UserID string          `json:"userId"`
	Role   TeamMemberRole `json:"role"`
}