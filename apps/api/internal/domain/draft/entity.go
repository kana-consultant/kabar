// internal/domain/draft/entity.go
package draft

import "time"

type Draft struct {
	ID             string     `json:"id"`
	Title          string     `json:"title"`
	Topic          string     `json:"topic"`
	Article        string     `json:"article"`
	ImageURL       *string    `json:"imageUrl"`
	ImagePrompt    string     `json:"imagePrompt"`
	Status         string     `json:"status"`
	ScheduledFor   *time.Time `json:"scheduledFor"`
	TargetProducts []string   `json:"targetProducts"`
	HasImage       bool       `json:"hasImage"`
	TeamID         *string    `json:"teamId"`
	UserID         *string    `json:"userId"`
	CreatedBy      *string    `json:"createdBy"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

type DraftData struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Topic          string   `json:"topic"`
	Article        string   `json:"article"`
	ImageURL       *string  `json:"imageUrl"`
	ImagePrompt    string   `json:"imagePrompt"`
	TargetProducts []string `json:"targetProducts"`
}

type PublishResult struct {
	Results      interface{}
	SomeFailed   bool
	AllFailed    bool
	Status       string
	ScheduledFor *time.Time
}
