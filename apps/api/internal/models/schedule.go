package models

type ScheduleStatus string

const (
	ScheduleStatusPending   ScheduleStatus = "pending"
	ScheduleStatusProcessed ScheduleStatus = "processed"
	ScheduleStatusFailed    ScheduleStatus = "failed"
)

type ScheduleRequest struct {
	Title          string   `json:"title,omitempty"`
	Topic          string   `json:"topic,omitempty"`
	Article        string   `json:"article,omitempty"`
	ImageURL       string   `json:"imageUrl,omitempty"`
	ImagePrompt    string   `json:"imagePrompt,omitempty"`
	Status         string   `json:"status,omitempty"`
	ScheduledFor   string   `json:"scheduledFor,omitempty"`
	TargetProducts []string `json:"targetProducts,omitempty"`
	HasImage       bool     `json:"hasImage,omitempty"`
}
