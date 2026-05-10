// internal/domain/draft/request.go
package draft

type CreateDraftRequest struct {
	Title          string   `json:"title"`
	Topic          string   `json:"topic"`
	Article        string   `json:"article"`
	ImageURL       *string  `json:"image_url"`
	ImagePrompt    string   `json:"image_prompt"`
	TargetProducts []string `json:"target_products"`
	HasImage       bool     `json:"has_image"`
	ScheduledFor   string   `json:"scheduled_for"`
}

type DraftDataPost struct {
	Title          string   `json:"title"`
	Topic          string   `json:"topic"`
	Article        string   `json:"article"`
	ImageURL       *string  `json:"ImageURL"`
	ImagePrompt    string   `json:"ImagePrompt"`
	TargetProducts []string `json:"scheduled_for"`
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
	Title          string
	Topic          string
	Article        string
	ImageURL       *string
	TargetProducts []string
}
