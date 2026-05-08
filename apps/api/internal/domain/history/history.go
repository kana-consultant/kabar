package history

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type HistoryStatus string

const (
	HistoryStatusSuccess HistoryStatus = "success"
	HistoryStatusFailed  HistoryStatus = "failed"
	HistoryStatusPending HistoryStatus = "pending"
)

type HistoryAction string

const (
	HistoryActionPublished  HistoryAction = "published"
	HistoryActionScheduled  HistoryAction = "scheduled"
	HistoryActionDraftSaved HistoryAction = "draft_saved"
)

type History struct {
	ID             string        `json:"id"`
	Title          string        `json:"title"`
	Topic          string        `json:"topic"`
	Content        string        `json:"content"`
	ImageURL       *string       `json:"imageUrl,omitempty"`
	TargetProducts []string      `json:"targetProducts"`
	Status         HistoryStatus `json:"status"`
	Action         HistoryAction `json:"action"`
	ErrorMessage   *string       `json:"errorMessage,omitempty"`
	PublishedAt    time.Time     `json:"publishedAt"`
	ScheduledFor   *string       `json:"scheduledFor,omitempty"`
	CreatedBy      *string       `json:"createdBy,omitempty"`
	TeamID         *string       `json:"teamId,omitempty"`
	CreatedAt      time.Time     `json:"createdAt"`
}

type HistoryFilter struct {
	TeamID  string
	Status  string
	Action  string
	Search  string
	Limit   int
	Offset  int
	Topic   string
	OrderBy string
}

type HistoryStats struct {
	Total          int     `json:"total"`
	SuccessCount   int     `json:"success"`
	FailedCount    int     `json:"failed"`
	PublishedCount int     `json:"published"`
	ScheduledCount int     `json:"scheduled"`
	SuccessRate    float64 `json:"successRate"`
}

type CreateHistoryRequest struct {
	Title          string
	Topic          string
	Content        string
	ImageURL       *string
	TargetProducts []string
	Status         string
	Action         string
	ErrorMessage   *string
	ScheduledFor   *string
	CreatedBy      string
	TeamID         string
}

type UserStatus string

const (
	StatusActive   UserStatus = "active"
	StatusInactive UserStatus = "inactive"
)

type SimpleUserContext struct {
	UserID string
	TeamID string
	Role   string
}

func (c *SimpleUserContext) GetUserID() string { return c.UserID }
func (c *SimpleUserContext) GetTeamID() string { return c.TeamID }
func (c *SimpleUserContext) GetRole() string   { return c.Role }
func (c *SimpleUserContext) IsAdmin() bool     { return c.Role == "admin" || c.Role == "super_admin" }

// Helper: Parse pagination from query params
func ParsePagination(r *http.Request) (limit, offset int, err error) {
	const (
		DefaultLimit = 10
		MaxLimit     = 100
	)
	// default
	limit = DefaultLimit
	offset = 0

	// limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		l, e := strconv.Atoi(limitStr)
		if e != nil {
			return 0, 0, fmt.Errorf("invalid limit")
		}
		if l < 1 {
			return 0, 0, fmt.Errorf("limit must be >= 1")
		}
		if l > MaxLimit {
			l = MaxLimit
		}
		limit = l
	}

	// offset
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		o, e := strconv.Atoi(offsetStr)
		if e != nil {
			return 0, 0, fmt.Errorf("invalid offset")
		}
		if o < 0 {
			return 0, 0, fmt.Errorf("offset must be >= 0")
		}
		offset = o
	}

	return limit, offset, nil
}
