package draft

import (
	"context"
	"seo-backend/internal/domain/paginate"
	"seo-backend/internal/models"
)

type Service interface {
	// Draft management
	CreateDraft(ctx context.Context, req CreateDraftRequest, userID, teamID string) (string, error)
	UpdateDraft(ctx context.Context, id string, TeamID string, updates map[string]interface{}) error
	DeleteDraft(ctx context.Context, TeamID string, id string) error
	GetDraftByID(ctx context.Context, id string) (*DraftData, error)
	GetAll(ctx context.Context, userFilter models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[Draft], error)
	GetDashboardStats(ctx context.Context, userFilter models.UserContext) (*DraftStats, error)
	GetAllScheduled(ctx context.Context, userCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[Draft], error)
	// Publishing
	PublishDraft(ctx context.Context, id string, req CreateDraftRequest, userContext models.UserContext) (*PublishResult, error)
	PublishContent(ctx context.Context, req DraftDataPost, userCtx models.UserContext) (*PublishResult, error)

	// Scheduling
	ScheduleDraft(ctx context.Context, req ScheduleRequest, userCtx models.UserContext) (string, error)
	CancelSchedule(ctx context.Context, draftID string) error

	// SEO
	GetSEOScore(ctx context.Context, id string) (*SEOScore, error)

	// Similarity
	CheckSimilarity(ctx context.Context, id string, useRole models.UserContext, params paginate.PaginationParams) ([]SimilarityResult, error)
}
