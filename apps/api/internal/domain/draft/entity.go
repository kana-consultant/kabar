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

type Keywords struct {
	ID        string    `json:"id"`
	IDDraft   string    `json:"id_draft"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DraftData struct {
	ID             string     `json:"id"`
	Title          string     `json:"title"`
	Topic          string     `json:"topic"`
	Article        string     `json:"article"`
	ImageURL       *string    `json:"image_url"`
	ImagePrompt    string     `json:"image_prompt"`
	TargetProducts []string   `json:"target_products"`
	Keywords       []Keywords `json:"keywords,omitempty"` // Tambahkan ini
}

type PublishResult struct {
	Results      interface{} `json:"results"`
	SomeFailed   bool        `json:"some_failed"`
	AllFailed    bool        `json:"all_failed"`
	Status       string      `json:"status"`
	ScheduledFor *time.Time  `json:"scheduled_for"`
}

type SEOScore struct {
	Total       int            `json:"total"`
	Details     map[string]int `json:"details"`
	Suggestions []string       `json:"suggestions"`
}

type SimilarityResult struct {
	DraftID    string  `json:"draft_id"`
	Title      string  `json:"title"`
	Similarity float64 `json:"similarity"`
}
