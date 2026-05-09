package scheduler

import (
	"context"
	"time"
)

// ScheduledTask - entity untuk task yang dijadwalkan
type ScheduledTask struct {
	ID             string     `json:"id"`
	DraftID        string     `json:"draft_id"`
	Title          string     `json:"title"`
	Topic          string     `json:"topic"`
	Article        string     `json:"article"`
	ImageURL       string     `json:"image_url"`
	ImagePrompt    string     `json:"image_prompt"`
	TargetProducts []string   `json:"target_products"`
	ScheduledFor   time.Time  `json:"scheduled_for"`
	TeamID         string     `json:"team_id"`
	UserID         string     `json:"user_id"`
	Status         string     `json:"status"`
	RetryCount     int        `json:"retry_count"`
	MaxRetries     int        `json:"max_retries"`
	CreatedAt      time.Time  `json:"created_at"`
	ExecutedAt     *time.Time `json:"executed_at,omitempty"`
	Error          string     `json:"error,omitempty"`
}

// ScheduleTaskData - data untuk scheduling
type ScheduleTaskData struct {
	DraftID        string
	Title          string
	Topic          string
	Article        string
	ImageURL       string
	ImagePrompt    string
	TargetProducts []string
	TeamID         string
	UserID         string
}

// Scheduler - interface untuk scheduling
type Scheduler interface {
	// ScheduleTask menjadwalkan task
	ScheduleTask(ctx context.Context, data ScheduleTaskData, scheduledFor time.Time) (string, error)

	// CancelTask membatalkan task yang sudah dijadwalkan
	CancelTask(ctx context.Context, draftID string) error

	// GetTask mengambil task berdasarkan ID
	GetTask(ctx context.Context, taskID string) (*ScheduledTask, error)

	// GetScheduledTasks mengambil semua task yang dijadwalkan
	GetScheduledTasks(ctx context.Context) ([]*ScheduledTask, error)

	// Start memulai scheduler
	Start(ctx context.Context) error

	// Stop menghentikan scheduler
	Stop(ctx context.Context) error

	// RegisterHandler mendaftarkan handler untuk task
	RegisterHandler(taskName string, handler TaskHandler)
}

// TaskHandler - fungsi untuk menangani task ketika dieksekusi
type TaskHandler func(ctx context.Context, task *ScheduledTask) error
