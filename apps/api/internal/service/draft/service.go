package draft

import (
	"database/sql"
	"fmt"

	"seo-backend/internal/helper"
	"seo-backend/internal/models"
	"seo-backend/internal/scheduler"
)

type Service struct {
	db             *sql.DB
	redisScheduler *scheduler.RedisScheduler
	postService    *helper.PostService
}

type DraftData struct {
	ID             string
	Title          string
	Topic          string
	Article        string
	ImageURL       *string
	ImagePrompt    string
	TargetProducts []string
	TeamID         string
	UserID         string
}

func NewService(db *sql.DB, redisScheduler *scheduler.RedisScheduler) *Service {
	return &Service{
		db:             db,
		redisScheduler: redisScheduler,
		postService:    helper.NewPostService(),
	}
}

// Public methods
func (s *Service) CreateDraft(req models.CreateDraftRequest, userID, teamID string) (string, error) {
	return s.CreateDraftRecord(req, userID, teamID)
}

func (s *Service) UpdateDraft(id string, updates map[string]interface{}) error {
	data := prepareUpdateData(updates)
	if len(data) == 0 {
		return fmt.Errorf("no fields to update")
	}
	return s.UpdateDraftRecord(id, data)
}

func (s *Service) DeleteDraft(id string) error {
	return s.DeleteDraftRecord(id)
}
