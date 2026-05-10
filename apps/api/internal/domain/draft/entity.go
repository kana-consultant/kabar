// internal/domain/draft/entity.go
package draft

import "time"

type Draft struct {
	ID             string     `json:"id"`
	Title          string     `json:"title"`
	Topic          string     `json:"topic"`
	Article        string     `json:"article"`
	ImageURL       *string    `json:"image_url"`
	ImagePrompt    string     `json:"image_prompt"`
	Status         string     `json:"status"`
	ScheduledFor   *time.Time `json:"scheduled_for"`
	TargetProducts []string   `json:"target_products"`
	HasImage       bool       `json:"has_image"`
	TeamID         *string    `json:"team_id"`
	UserID         *string    `json:"user_id"`
	CreatedBy      *string    `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type DraftData struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Topic          string   `json:"topic"`
	Article        string   `json:"article"`
	ImageURL       *string  `json:"image_url"`
	ImagePrompt    string   `json:"image_prompt"`
	TargetProducts []string `json:"target_products"`
}

type PublishResult struct {
	Results      interface{} `json:"results"`
	SomeFailed   bool        `json:"some_failed"`
	AllFailed    bool        `json:"all_failed"`
	Status       string      `json:"status"`
	ScheduledFor *time.Time  `json:"scheduled_for"`
}
