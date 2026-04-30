package models

import (
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
	ID             string         `json:"id"`
	Title          string         `json:"title"`
	Topic          string         `json:"topic"`
	Content        string         `json:"content"`
	ImageURL       *string        `json:"imageUrl,omitempty"`
	TargetProducts []string       `json:"targetProducts"`
	Status         HistoryStatus  `json:"status"`
	Action         HistoryAction  `json:"action"`
	ErrorMessage   *string        `json:"errorMessage,omitempty"`
	PublishedAt    time.Time      `json:"publishedAt"`
	ScheduledFor   *time.Time     `json:"scheduledFor,omitempty"`
	CreatedBy      *string        `json:"createdBy,omitempty"`
	TeamID         *string        `json:"teamId,omitempty"`
	CreatedAt      time.Time      `json:"createdAt"`
}