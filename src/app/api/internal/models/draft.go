package models

import (
	"time"
)

type Draft struct {
	ID               string     `json:"id"`
	Title            string     `json:"title"`
	Topic            string     `json:"topic"`
	Article          string     `json:"article"`
	ImageURL         *string    `json:"imageUrl,omitempty"`
	ImagePrompt      string     `json:"imagePrompt,omitempty"`
	Status           string     `json:"status"`
	ScheduledFor     *time.Time `json:"scheduledFor,omitempty"`
	ScheduledPattern string     `json:"scheduledPattern,omitempty"`
	TargetProducts   []string   `json:"targetProducts"`
	HasImage         bool       `json:"hasImage"`
	TeamID           *string    `json:"teamId,omitempty"`
	UserID           *string    `json:"userId,omitempty"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	PublishedAt      *time.Time `json:"publishedAt,omitempty"`
}

type CreateDraftRequest struct {
	Title          string   `json:"title"`
	Topic          string   `json:"topic"`
	Article        string   `json:"article"`
	ImageURL       string   `json:"imageUrl,omitempty"`
	ImagePrompt    string   `json:"imagePrompt,omitempty"`
	TargetProducts []string `json:"targetProducts"`
	HasImage       bool     `json:"hasImage"`
	TeamID         string   `json:"teamId,omitempty"`
	UserID         string   `json:"userId,omitempty"`
}

type DraftDataPost struct {
	Title          string
	Topic          string
	Article        string
	ImageURL       *string
	TargetProducts []string
	ImagePrompt    string
}

type UpdateDraftRequest struct {
	Title          *string    `json:"title,omitempty"`
	Topic          *string    `json:"topic,omitempty"`
	Article        *string    `json:"article,omitempty"`
	ImageURL       *string    `json:"imageUrl,omitempty"`
	ImagePrompt    *string    `json:"imagePrompt,omitempty"`
	Status         *string    `json:"status,omitempty"`
	ScheduledFor   *time.Time `json:"scheduledFor,omitempty"`
	TargetProducts *[]string  `json:"targetProducts,omitempty"`
	HasImage       *bool      `json:"hasImage,omitempty"`
}

type PublishRequest struct {
	ScheduledFor string `json:"scheduledFor,omitempty"`
}

type PublishResult struct {
	Product   string `json:"product"`
	Success   bool   `json:"success"`
	Error     string `json:"error,omitempty"`
	Response  string `json:"response,omitempty"`
	ProductID string `json:"productId,omitempty"`
}
