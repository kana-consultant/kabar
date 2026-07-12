package draft

import (
	"context"
	"seo-backend/internal/domain/paginate"
	"seo-backend/internal/models"

	"time"
)

type Service interface {
	// Draft management
	CreateDraft(ctx context.Context, req CreateDraftRequest, userID, teamID string) (string, error)
	UpdateDraft(ctx context.Context, id string, userID, TeamID string, updates CreateDraftRequest) error
	DeleteDraft(ctx context.Context, TeamID string, id string) error
	GetDraftByID(ctx context.Context, id string) (*DraftData, error)
	GetAll(ctx context.Context, userFilter models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[Draft], error)
	GetDashboardStats(ctx context.Context, userFilter models.UserContext) (*DraftStats, error)
	GetAllScheduled(ctx context.Context, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[Draft], error)
	RescheduleDraft(ctx context.Context, draftID string, newScheduleTime time.Time, userCtx models.UserContext) (*PublishResult, error)
	// Publishing
	PublishDraft(ctx context.Context, id string, req CreateDraftRequest, userContext models.UserContext) (*PublishResult, error)
	PublishContent(ctx context.Context, req CreateDraftRequest, userCtx models.UserContext) (*PublishResult, error)

	// Scheduling
	ScheduleDraft(ctx context.Context, req ScheduleRequest, userCtx models.UserContext) (string, error)
	CancelSchedule(ctx context.Context, draftID string) error

	// SEO
	GetSEOScore(ctx context.Context, id string) (*SEOScore, error)

	// Similarity
	CheckSimilarity(ctx context.Context, id string, useRole models.UserContext, params paginate.PaginationParams) ([]SimilarityResult, error)
}
