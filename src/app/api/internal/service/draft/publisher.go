package draft

import (
	"fmt"
	"log"
	"time"

	"seo-backend/internal/helper"
	"seo-backend/internal/models"
	"seo-backend/internal/scheduler"
)

type PublishResult struct {
	Results      interface{}
	SomeFailed   bool
	AllFailed    bool
	Status       string
	ScheduledFor *time.Time
}

func (s *Service) PublishDraft(id string, scheduledForStr string, teamID, userID string) (*PublishResult, error) {
	draft, err := s.getDraftData(id)
	if err != nil {
		return nil, err
	}

	if scheduledForStr != "" {
		return s.scheduleDraft(id, scheduledForStr, draft, teamID, userID)
	}

	return s.processPublish(draft, id, teamID, userID)
}

func (s *Service) PublishContent(req models.DraftDataPost, teamID, userID string) (*PublishResult, error) {
	if err := validatePublishRequest(req); err != nil {
		return nil, err
	}

	result, someFailed, allFailed, err := s.postService.ProcessDraftProducts(req)
	if err != nil {
		return nil, fmt.Errorf("failed to process products: %w", err)
	}

	if err := s.InsertHistory(req, userID, teamID, "published"); err != nil {
		log.Printf("Failed to insert history: %v", err)
	}

	return &PublishResult{
		Results:    result,
		SomeFailed: someFailed,
		AllFailed:  allFailed,
		Status:     "published",
	}, nil
}

func (s *Service) scheduleDraft(id, scheduledForStr string, draft *DraftData, teamID, userID string) (*PublishResult, error) {
	scheduledFor, err := helper.ParseWIBTime(scheduledForStr)
	if err != nil {
		return nil, err
	}

	if err := s.updateDraftStatus(id, "scheduled", &scheduledFor); err != nil {
		return nil, err
	}

	taskData := &scheduler.ScheduledTask{
		ID:             id,
		Title:          draft.Title,
		Topic:          draft.Topic,
		Article:        draft.Article,
		ImageURL:       *draft.ImageURL,
		ImagePrompt:    draft.ImagePrompt,
		TargetProducts: draft.TargetProducts,
		TeamID:         teamID,
		UserID:         userID,
	}

	if err := s.redisScheduler.ScheduleDraftTask(id, scheduledFor, taskData); err != nil {
		return nil, err
	}

	return &PublishResult{
		Status:       "scheduled",
		ScheduledFor: &scheduledFor,
	}, nil
}

func (s *Service) processPublish(draft *DraftData, id, teamID, userID string) (*PublishResult, error) {
	draftData := models.DraftDataPost{
		Title:          draft.Title,
		Topic:          draft.Topic,
		Article:        draft.Article,
		ImageURL:       draft.ImageURL,
		ImagePrompt:    draft.ImagePrompt,
		TargetProducts: draft.TargetProducts,
	}

	result, someFailed, allFailed, err := s.postService.ProcessDraftProducts(draftData)
	if err != nil {
		return nil, err
	}

	if allFailed {
		return &PublishResult{
			Results:    result,
			AllFailed:  true,
			SomeFailed: someFailed,
			Status:     "failed",
		}, nil
	}

	if err := s.deleteDraftAfterPublish(id); err != nil {
		log.Printf("Failed to delete draft: %v", err)
	}

	if err := s.InsertHistory(draftData, userID, teamID, "published"); err != nil {
		log.Printf("Failed to insert history: %v", err)
	}

	return &PublishResult{
		Results:    result,
		SomeFailed: someFailed,
		AllFailed:  allFailed,
		Status:     "published",
	}, nil
}
