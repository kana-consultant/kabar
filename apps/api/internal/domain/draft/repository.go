// internal/domain/draft/repository.go
package draft

import (
	"context"

	"seo-backend/internal/domain/paginate"
	"seo-backend/internal/models"
	"time"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*DraftData, error)
	GetAll(ctx context.Context, filter models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[Draft], error)
	GetDashboardStats(ctx context.Context, filter models.UserContext) (*DraftStats, error)
	Create(ctx context.Context, req CreateDraftRequest, userID, teamID string) (string, error)
	Update(ctx context.Context, TeamID string, id string, data map[string]interface{}) error
	UpdateStatus(ctx context.Context, id string, status string, scheduledFor *time.Time) error
	Delete(ctx context.Context, TeamID string, id string) error
	InsertScheduledDraft(ctx context.Context, req ScheduleRequest, scheduledFor time.Time, teamID, userID string) (string, error)
	GetAllScheduled(ctx context.Context, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[Draft], error)
}
