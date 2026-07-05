// internal/domain/draft/request.go
package draft

import "time"

type CreateDraftRequest struct {
	Title          string    `json:"title"`
	Topic          string    `json:"topic"`
	Article        string    `json:"article"`
	ImageURL       *string   `json:"image_url"`
	ImagePrompt    string    `json:"image_prompt"`
	TargetProducts []string  `json:"target_products"`
	HasImage       bool      `json:"has_image"`
	ScheduledFor   string    `json:"scheduled_for"`
	Slug           string    `json:"slug"`
	Keywords       []string  `json:"keywords"`
	SEOScore       int       `json:"seo_score"`
	Excerpt        string    `json:"excerpt"`
	UpdateAt       time.Time `json:"update_at"`

	// Tambahan yang dibutuhkan untuk update
	Status    string
	TeamID    string
	UserID    string
	CreatedBy string
}

type DraftDataPost struct {
	Id             string   `json:"id"`
	Title          string   `json:"title"`
	Topic          string   `json:"topic"`
	Article        string   `json:"article"`
	ImageURL       *string  `json:"image_url"`
	ImagePrompt    string   `json:"image_prompt"`
	TargetProducts []string `json:"target_products"`
	Slug           string   `json:"slug"`
	Keywords       []string `json:"keywords"`
	Excerpt        string   `json:"excerpt"`
}

type ScheduleRequest struct {
	Title          string   `json:"Title"`
	Topic          string   `json:"Topic"`
	Article        string   `json:"Article"`
	ImageURL       string   `json:"ImageURL"`
	ImagePrompt    string   `json:"ImagePrompt"`
	TargetProducts []string `json:"target_products"`
	HasImage       bool     `json:"HasImage"`
	ScheduledFor   string   `json:"scheduled_for"`
}

type PublishHistoryRequest struct {
	Id             string
	Title          string
	Topic          string
	Article        string
	ImageURL       *string
	SEOScore       int
	TargetProducts []string
	Keywords       []string
	Excerpt        string
	Slug           string
}
