package draft

import (
	"context"
	"database/sql"
	"fmt"

	"seo-backend/internal/helper"
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
		postService:    helper.NewPostService(db),
	}
}

func (s *Service) CreateDraft(ctx context.Context, req CreateDraftRequest, userID, teamID string) (string, error) {
	return s.CreateDraftRecord(req, userID, teamID)
}

func (s *Service) UpdateDraft(ctx context.Context, id string, updates map[string]interface{}) error {
	data := prepareUpdateData(updates)
	if len(data) == 0 {
		return fmt.Errorf("no fields to update")
	}
	return s.UpdateDraftRecord(id, data)
}

func (s *Service) DeleteDraft(ctx context.Context, id string) error {
	return s.DeleteDraftRecord(id)
}

func (s *Service) GetSEOScore(ctx context.Context, id string) (*SEOScore, error) {
	draft, err := s.GetDraftByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("draft not found: %w", err)
	}

	score := CalculateSEOScore(draft.Title, draft.Article, draft.Excerpt, draft.Topic)
	return &score, nil
}

// GetSEOScore implements draft.Service
func (s *DraftServiceImpl) GetSEOScore(ctx context.Context, id string) (*draft.SEOScore, error) {
	draftData, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("draft not found: %w", err)
	}

	score := calculateSEOScore(draftData.Title, draftData.Article, draftData.Topic, draftData.Topic)
	return &score, nil
}