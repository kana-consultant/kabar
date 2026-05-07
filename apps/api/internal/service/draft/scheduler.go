package draft

import (
	"fmt"
	"log"
	"time"

	"seo-backend/internal/helper"
	"seo-backend/internal/models"
	"seo-backend/internal/scheduler"
)

type Scheduler struct {
	repo           *Repository
	redisScheduler *scheduler.RedisScheduler
}

func NewScheduler(repo *Repository, redisScheduler *scheduler.RedisScheduler) *Scheduler {
	return &Scheduler{
		repo:           repo,
		redisScheduler: redisScheduler,
	}
}

func (s *Scheduler) ScheduleDraft(req models.ScheduleRequest, teamID, userID string) (string, error) {
	scheduledFor, err := helper.ParseWIBTime(req.ScheduledFor)
	if err != nil {
		return "", err
	}

	if scheduledFor.Before(time.Now()) {
		return "", fmt.Errorf("scheduled time must be in the future")
	}

	draftID, err := s.repo.InsertScheduledDraft(req, scheduledFor, teamID, userID)
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
		s.repo.Delete(draftID)
		return "", fmt.Errorf("failed to schedule in Redis: %w", err)
	}

	return draftID, nil
}

func (s *Scheduler) CancelSchedule(draftID string) error {
	if err := s.redisScheduler.CancelScheduledTask(draftID); err != nil {
		log.Printf("Failed to cancel Redis task: %v", err)
	}

	return s.repo.UpdateStatus(draftID, "draft", nil)
}
