package draft

import (
	"context"
	"time"

	"seo-backend/internal/models"
)

type Repository interface {
	// Basic CRUD
	GetAll(ctx context.Context, filters GetAllFilters) ([]models.Draft, error)
	GetByID(ctx context.Context, id string) (*models.Draft, error)
	Create(ctx context.Context, draft *models.Draft) (string, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error

	// Publishing related
	GetDraftForPublish(ctx context.Context, id string) (*PublishData, error)
	UpdatePublishStatus(ctx context.Context, id string, status string, scheduledFor *time.Time, pattern string) error
	UpdateProductSyncStatus(ctx context.Context, productID string) error
	SaveToHistory(ctx context.Context, history *PublishHistory) error
}

type GetAllFilters struct {
	TeamID *string
	UserID *string
	Role   string
	Status string
	Search string
}

type PublishData struct {
	ID             string
	Title          string
	Topic          string
	Article        string
	ImageURL       *string
	ImagePrompt    string
	TargetProducts []string
	TeamID         *string
}

type PublishHistory struct {
	Title       string
	Topic       string
	Content     string
	ImageURL    *string
	Products    []string
	Status      string
	Action      string
	Error       string
	PublishedAt time.Time
	CreatedBy   string
	TeamID      *string
}
