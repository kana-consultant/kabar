// internal/domain/draft/repository.go
package draft

import (
	"context"
	"time"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*DraftData, error)
	GetAll(ctx context.Context, teamID string) (*[]Draft, error)
	Create(ctx context.Context, req CreateDraftRequest, userID, teamID string) (string, error)
	Update(ctx context.Context, id string, data map[string]interface{}) error
	UpdateStatus(ctx context.Context, id string, status string, scheduledFor *time.Time) error
	Delete(ctx context.Context, id string) error
	InsertScheduledDraft(ctx context.Context, req ScheduleRequest, scheduledFor time.Time, teamID, userID string) (string, error)
	InsertHistory(ctx context.Context, req PublishHistoryRequest, userID, teamID, action string) error
	GetAllScheduled(ctx context.Context, teamID string) (*[]Draft, error)
}
