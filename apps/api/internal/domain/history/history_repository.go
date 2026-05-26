package history

import (
	"context"
	"time"

	"seo-backend/internal/domain/paginate"
	"seo-backend/internal/models"
)

// HistoryQuery for filtering
type HistoryQuery struct {
	TeamID   string
	Status   string
	Topic    string
	Search   string
	FromDate *time.Time
	ToDate   *time.Time
	Limit    int
	Offset   int
	OrderBy  string
}

// Repository interface
type HistoryRepository interface {
	// Create
	Create(ctx context.Context, data *History) (string, error)

	// Read
	GetByID(ctx context.Context, id string) (*History, error)
	GetAll(ctx context.Context, userCtx models.UserContext, filter HistoryFilter) (*paginate.PaginatedResult[History], error)
	GetAllWithQuery(ctx context.Context, query string, args []interface{}) ([]History, error)
	GetByTeamID(ctx context.Context, teamID string) ([]History, error)
	GetByCreatedBy(ctx context.Context, createdBy string) ([]History, error)
	GetByStatus(ctx context.Context, status string) ([]History, error)
	GetRecentActivity(ctx context.Context, teamID string, limit int) ([]History, error)

	// Count
	Count(ctx context.Context, query HistoryFilter) (int, error)
	GetCountByStatus(ctx context.Context, teamID string) (map[string]int, error)

	// Update
	Update(ctx context.Context, id string, updates map[string]interface{}) error

	// Delete
	Delete(ctx context.Context, id string) error
	DeleteByTeamID(ctx context.Context, teamID string) error
	DeleteByStatus(ctx context.Context, status string) error

	// Statistics
	GetStats(ctx context.Context, query *HistoryFilter) (*HistoryStats, error)
}

// QueryBuilder interface
type QueryBuilder interface {
	BuildListQuery(filters HistoryQuery) (string, []interface{})
	BuildCountQuery(filters HistoryQuery) (string, []interface{})
}
