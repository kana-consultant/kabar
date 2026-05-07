// internal/application/draft/service.go
package draft

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"seo-backend/internal/domain/draft"
	"seo-backend/internal/helper"
	"seo-backend/internal/scheduler"
)

type DraftServiceImpl struct {
	repo           draft.Repository
	redisScheduler *scheduler.RedisScheduler
	postService    *helper.PostService
}

func NewService(
	repo draft.Repository,
	redisScheduler *scheduler.RedisScheduler,
	postService *helper.PostService,
) draft.Service {
	return &DraftServiceImpl{
		repo:           repo,
		redisScheduler: redisScheduler,
		postService:    postService,
	}
}

func (s *DraftServiceImpl) GetAll(ctx context.Context, TeamID string) (*[]draft.Draft, error) {

	return s.repo.GetAll(ctx, TeamID)
}

func (s *DraftServiceImpl) GetAllScheduled(ctx context.Context, TeamID string) (*[]draft.Draft, error) {

	return s.repo.GetAllScheduled(ctx, TeamID)
}

// CreateDraft implements draft.Service
func (s *DraftServiceImpl) CreateDraft(ctx context.Context, req draft.CreateDraftRequest, userID, teamID string) (string, error) {
	return s.repo.Create(ctx, req, userID, teamID)
}

// UpdateDraft implements draft.Service
func (s *DraftServiceImpl) UpdateDraft(ctx context.Context, id string, updates map[string]interface{}) error {
	data := prepareUpdateData(updates)
	if len(data) == 0 {
		return fmt.Errorf("no fields to update")
	}
	return s.repo.Update(ctx, id, data)
}

// DeleteDraft implements draft.Service
func (s *DraftServiceImpl) DeleteDraft(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// GetDraftByID implements draft.Service
func (s *DraftServiceImpl) GetDraftByID(ctx context.Context, id string) (*draft.DraftData, error) {
	return s.repo.GetByID(ctx, id)
}

// PublishDraft implements draft.Service
func (s *DraftServiceImpl) PublishDraft(ctx context.Context, id string, scheduledForStr string, teamID, userID string) (*draft.PublishResult, error) {
	draftData, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if scheduledForStr != "" {
		return s.scheduleDraft(ctx, id, scheduledForStr, draftData, teamID, userID)
	}

	return s.processPublish(ctx, draftData, id, teamID, userID)
}

// PublishContent implements draft.Service
func (s *DraftServiceImpl) PublishContent(ctx context.Context, req draft.DraftDataPost, teamID, userID string) (*draft.PublishResult, error) {
	if err := validatePublishRequest(req); err != nil {
		return nil, err
	}

	result, someFailed, allFailed, err := s.postService.ProcessDraftProducts(req)
	if err != nil {
		return nil, fmt.Errorf("failed to process products: %w", err)
	}

	historyReq := draft.PublishHistoryRequest{
		Title:          req.Title,
		Topic:          req.Topic,
		Article:        req.Article,
		ImageURL:       req.ImageURL,
		TargetProducts: req.TargetProducts,
	}

	if err := s.repo.InsertHistory(ctx, historyReq, userID, teamID, "published"); err != nil {
		log.Printf("Failed to insert history: %v", err)
	}

	return &draft.PublishResult{
		Results:    result,
		SomeFailed: someFailed,
		AllFailed:  allFailed,
		Status:     "published",
	}, nil
}

// ScheduleDraft implements draft.Service
func (s *DraftServiceImpl) ScheduleDraft(ctx context.Context, req draft.ScheduleRequest, teamID, userID string) (string, error) {
	scheduledFor, err := helper.ParseWIBTime(req.ScheduledFor)
	if err != nil {
		return "", err
	}

	if scheduledFor.Before(time.Now()) {
		return "", fmt.Errorf("scheduled time must be in the future")
	}

	draftID, err := s.repo.InsertScheduledDraft(ctx, req, scheduledFor, teamID, userID)
	if err != nil {
		return "", err
	}

	taskData := &scheduler.ScheduledTask{
		DraftID:        draftID,
		Title:          req.Title,
		Topic:          req.Topic,
		Article:        req.Article,
		ImageURL:       req.ImageURL,
		ImagePrompt:    req.ImagePrompt,
		TargetProducts: req.TargetProducts,
		TeamID:         teamID,
		UserID:         userID,
	}

	if err := s.redisScheduler.ScheduleDraftTask(draftID, scheduledFor, taskData); err != nil {
		s.repo.Delete(ctx, draftID)
		return "", fmt.Errorf("failed to schedule in Redis: %w", err)
	}

	return draftID, nil
}

// CancelSchedule implements draft.Service
func (s *DraftServiceImpl) CancelSchedule(ctx context.Context, draftID string) error {
	if err := s.redisScheduler.CancelScheduledTask(draftID); err != nil {
		log.Printf("Failed to cancel Redis task: %v", err)
	}
	return s.repo.UpdateStatus(ctx, draftID, "draft", nil)
}

// Private methods
func (s *DraftServiceImpl) scheduleDraft(ctx context.Context, id, scheduledForStr string, draftData *draft.DraftData, teamID, userID string) (*draft.PublishResult, error) {
	scheduledFor, err := helper.ParseWIBTime(scheduledForStr)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateStatus(ctx, id, "scheduled", &scheduledFor); err != nil {
		return nil, err
	}

	taskData := &scheduler.ScheduledTask{
		DraftID:        id,
		Title:          draftData.Title,
		Topic:          draftData.Topic,
		Article:        draftData.Article,
		ImageURL:       *draftData.ImageURL,
		ImagePrompt:    draftData.ImagePrompt,
		TargetProducts: draftData.TargetProducts,
		TeamID:         teamID,
		UserID:         userID,
	}

	if err := s.redisScheduler.ScheduleDraftTask(id, scheduledFor, taskData); err != nil {
		return nil, err
	}

	return &draft.PublishResult{
		Status:       "scheduled",
		ScheduledFor: &scheduledFor,
	}, nil
}

func (s *DraftServiceImpl) processPublish(ctx context.Context, draftData *draft.DraftData, id, teamID, userID string) (*draft.PublishResult, error) {
	draftPost := draft.DraftDataPost{
		Title:          draftData.Title,
		Topic:          draftData.Topic,
		Article:        draftData.Article,
		ImageURL:       draftData.ImageURL,
		ImagePrompt:    draftData.ImagePrompt,
		TargetProducts: draftData.TargetProducts,
	}

	result, someFailed, allFailed, err := s.postService.ProcessDraftProducts(draftPost)
	if err != nil {
		return nil, err
	}

	if allFailed {
		return &draft.PublishResult{
			Results:    result,
			AllFailed:  true,
			SomeFailed: someFailed,
			Status:     "failed",
		}, nil
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		log.Printf("Failed to delete draft: %v", err)
	}

	historyReq := draft.PublishHistoryRequest{
		Title:          draftPost.Title,
		Topic:          draftPost.Topic,
		Article:        draftPost.Article,
		ImageURL:       draftPost.ImageURL,
		TargetProducts: draftPost.TargetProducts,
	}

	if err := s.repo.InsertHistory(ctx, historyReq, userID, teamID, "published"); err != nil {
		log.Printf("Failed to insert history: %v", err)
	}

	return &draft.PublishResult{
		Results:    result,
		SomeFailed: someFailed,
		AllFailed:  allFailed,
		Status:     "published",
	}, nil
}

// Helper functions
func prepareUpdateData(updates map[string]interface{}) map[string]interface{} {
	fieldMap := map[string]string{
		"title":          "title",
		"topic":          "topic",
		"article":        "article",
		"imageUrl":       "image_url",
		"imagePrompt":    "image_prompt",
		"status":         "status",
		"scheduledFor":   "scheduled_for",
		"targetProducts": "target_products",
		"hasImage":       "has_image",
	}

	data := make(map[string]interface{})
	for key, value := range updates {
		if dbField, ok := fieldMap[key]; ok {
			if key == "targetProducts" {
				jsonValue, _ := json.Marshal(value)
				data[dbField] = jsonValue
			} else {
				data[dbField] = value
			}
		}
	}
	return data
}

func validatePublishRequest(req draft.DraftDataPost) error {
	if req.Title == "" || req.Article == "" || len(req.TargetProducts) == 0 {
		return fmt.Errorf("title, article, and target_products are required")
	}
	return nil
}
