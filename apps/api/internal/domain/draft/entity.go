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
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Topic          string   `json:"topic"`
	Article        string   `json:"article"`
	ImageURL       *string  `json:"image_url"`
	ImagePrompt    string   `json:"image_prompt"`
	TargetProducts []string `json:"target_products"`
	Keywords       []string `json:"keywords,omitempty"` // Tambahkan ini
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

type DraftStats struct {
	TotalDraft        int             `json:"total_draft"`
	TotalWithImage    int             `json:"total_with_image"`
	TotalWithoutImage int             `json:"total_without_image"`
	TotalScheduled    int             `json:"total_scheduled"`
	DailyActivity     []DailyActivity `json:"daily_activity"`
	ProductCoverage   map[string]int  `json:"product_coverage"` // name -> count
	ProductStatus     map[string]int  `json:"product_status"`   // active/pending/inactive
}

type DailyActivity struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}
