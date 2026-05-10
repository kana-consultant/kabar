// internal/domain/draft/service.go
package draft

import (
	"context"
)

type Service interface {
	// Draft management
	CreateDraft(ctx context.Context, req CreateDraftRequest, userID, teamID string) (string, error)
	UpdateDraft(ctx context.Context, id string, updates map[string]interface{}) error
	DeleteDraft(ctx context.Context, id string) error
	GetDraftByID(ctx context.Context, id string) (*DraftData, error)
	GetAll(ctx context.Context, TeamID string) (*[]Draft, error)
	GetAllScheduled(ctx context.Context, TeamID string) (*[]Draft, error)

	// Publishing
	PublishDraft(ctx context.Context, id string, req CreateDraftRequest, teamID, userID string) (*PublishResult, error)
	PublishContent(ctx context.Context, req DraftDataPost, teamID, userID string) (*PublishResult, error)

	// Scheduling
	ScheduleDraft(ctx context.Context, req ScheduleRequest, teamID, userID string) (string, error)
	CancelSchedule(ctx context.Context, draftID string) error
}
