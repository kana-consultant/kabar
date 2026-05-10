package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"

	"seo-backend/internal/domain/scheduler"
	"seo-backend/internal/helper"
)

type RedisSchedulerImpl struct {
	redisClient  *redis.Client
	cron         *cron.Cron
	taskHandlers map[string]scheduler.TaskHandler
}

func NewRedisSchedulerImpl(redisClient *redis.Client) *RedisSchedulerImpl {
	return &RedisSchedulerImpl{
		redisClient:  redisClient,
		cron:         cron.New(cron.WithSeconds()),
		taskHandlers: make(map[string]scheduler.TaskHandler),
	}
}

// ScheduleTask implements scheduler.Scheduler
func (s *RedisSchedulerImpl) ScheduleTask(ctx context.Context, data scheduler.ScheduleTaskData, scheduledFor time.Time) (string, error) {
	taskID := fmt.Sprintf("draft_%s_%d", data.DraftID, scheduledFor.Unix())

	task := &scheduler.ScheduledTask{
		ID:             taskID,
		DraftID:        data.DraftID,
		Title:          data.Title,
		Topic:          data.Topic,
		Article:        data.Article,
		ImageURL:       data.ImageURL,
		ImagePrompt:    data.ImagePrompt,
		TargetProducts: data.TargetProducts,
		ScheduledFor:   scheduledFor,
		TeamID:         data.TeamID,
		UserID:         data.UserID,
		Status:         "pending",
		RetryCount:     0,
		MaxRetries:     3,
		CreatedAt:      helper.ParseWIBTime(time.Now().Format(time.RFC3339)),
	}

	// Save to Redis
	taskKey := fmt.Sprintf("schedule:draft:%s", taskID)
	taskDataBytes, err := json.Marshal(task)
	if err != nil {
		return "", fmt.Errorf("failed to marshal task: %w", err)
	}

	err = s.redisClient.Set(ctx, taskKey, taskDataBytes, 0).Err()
	if err != nil {
		return "", fmt.Errorf("failed to save task to redis: %w", err)
	}

	// Schedule with cron
	cronExpr := fmt.Sprintf("0 %d %d %d %d *",
		scheduledFor.Minute(),
		scheduledFor.Hour(),
		scheduledFor.Day(),
		int(scheduledFor.Month()))

	_, err = s.cron.AddFunc(cronExpr, func() {
		s.executeDraftTask(ctx, taskID)
	})

	if err != nil {
		return "", fmt.Errorf("failed to schedule cron: %w", err)
	}

	log.Printf("[REDIS] Draft %s scheduled at %s", data.DraftID, scheduledFor.Format(time.RFC3339))
	return taskID, nil
}

// CancelTask implements scheduler.Scheduler
func (s *RedisSchedulerImpl) CancelTask(ctx context.Context, draftID string) error {
	pattern := fmt.Sprintf("schedule:draft:draft_%s_*", draftID)
	keys, err := s.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to find tasks: %w", err)
	}

	for _, key := range keys {
		if err := s.redisClient.Del(ctx, key).Err(); err != nil {
			log.Printf("[REDIS] Failed to delete key %s: %v", key, err)
		}
	}

	log.Printf("[REDIS] Cancelled scheduled tasks for draft %s", draftID)
	return nil
}

// GetTask implements scheduler.Scheduler
func (s *RedisSchedulerImpl) GetTask(ctx context.Context, taskID string) (*scheduler.ScheduledTask, error) {
	taskKey := fmt.Sprintf("schedule:draft:%s", taskID)
	taskData, err := s.redisClient.Get(ctx, taskKey).Bytes()
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	var task scheduler.ScheduledTask
	if err := json.Unmarshal(taskData, &task); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task: %w", err)
	}

	return &task, nil
}

// GetScheduledTasks implements scheduler.Scheduler
func (s *RedisSchedulerImpl) GetScheduledTasks(ctx context.Context) ([]*scheduler.ScheduledTask, error) {
	pattern := "schedule:draft:*"
	keys, err := s.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys: %w", err)
	}

	var tasks []*scheduler.ScheduledTask
	for _, key := range keys {
		taskData, err := s.redisClient.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}

		var task scheduler.ScheduledTask
		if err := json.Unmarshal(taskData, &task); err != nil {
			continue
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}

// Start implements scheduler.Scheduler
func (s *RedisSchedulerImpl) Start(ctx context.Context) error {
	s.cron.Start()
	log.Println("[REDIS] Scheduler started")
	return nil
}

// Stop implements scheduler.Scheduler
func (s *RedisSchedulerImpl) Stop(ctx context.Context) error {
	s.cron.Stop()
	log.Println("[REDIS] Scheduler stopped")
	return nil
}

// RegisterHandler implements scheduler.Scheduler
func (s *RedisSchedulerImpl) RegisterHandler(taskName string, handler scheduler.TaskHandler) {
	s.taskHandlers[taskName] = handler
}

// executeDraftTask - internal function to execute scheduled task
func (s *RedisSchedulerImpl) executeDraftTask(ctx context.Context, taskID string) {
	log.Printf("[REDIS] Executing scheduled draft task: %s", taskID)

	// Get task from Redis
	task, err := s.GetTask(ctx, taskID)
	if err != nil {
		log.Printf("[REDIS] Failed to get task %s: %v", taskID, err)
		return
	}

	// Mark as processing
	task.Status = "processing"
	s.updateTask(ctx, task)

	// Find and execute handler
	handler, exists := s.taskHandlers["draft_publish"]
	if !exists {
		log.Printf("[REDIS] No handler registered for draft_publish")
		task.Status = "failed"
		task.Error = "no handler registered"
		s.updateTask(ctx, task)
		return
	}

	// Execute handler
	err = handler(ctx, task)
	if err != nil {
		log.Printf("[REDIS] Failed to execute task %s: %v", taskID, err)
		task.Status = "failed"
		task.Error = err.Error()
	} else {
		task.Status = "completed"
	}

	now := helper.ParseWIBTime(time.Now().Format(time.RFC3339))
	task.ExecutedAt = &now
	s.updateTask(ctx, task)
}

// updateTask - update task in Redis
func (s *RedisSchedulerImpl) updateTask(ctx context.Context, task *scheduler.ScheduledTask) error {
	taskKey := fmt.Sprintf("schedule:draft:%s", task.ID)
	taskData, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return s.redisClient.Set(ctx, taskKey, taskData, 0).Err()
}

func (s *RedisSchedulerImpl) ScheduleDraftTask(ctx context.Context, draftID string, scheduledFor time.Time, taskData *scheduler.ScheduledTask) error {
	taskID := fmt.Sprintf(
		"draft_%s_%d",
		draftID,
		scheduledFor.Unix(),
	)

	task := &scheduler.ScheduledTask{
		ID:             taskID,
		DraftID:        draftID,
		Title:          taskData.Title,
		Topic:          taskData.Topic,
		ImageURL:       taskData.ImageURL,
		Article:        taskData.Article,
		ImagePrompt:    taskData.ImagePrompt,
		TargetProducts: taskData.TargetProducts,
		ScheduledFor:   scheduledFor,
		TeamID:         taskData.TeamID,
		UserID:         taskData.UserID,
		Status:         "pending",
		RetryCount:     0,
		MaxRetries:     3,
		CreatedAt:      helper.ParseWIBTime(time.Now().Format(time.RFC3339)),
	}

	// Save to Redis
	taskKey := fmt.Sprintf("schedule:draft:%s", taskID)
	taskDataBytes, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	log.Println("taskKey:", taskKey)

	log.Println("scheduler ptr:", s)
	log.Println("redisClient nil:", s.redisClient == nil)

	// log.Println("scheduler nil:", s == nil)
	// log.Println("handler scheduler nil?", s.schedule)
	// log.Println("cron nil:", s.cron == nil)
	// log.Println("ctx nil:", s.ctx == nil)

	err = s.redisClient.Set(ctx, taskKey, taskDataBytes, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	// Schedule with cron
	cronExpr := fmt.Sprintf("0 %d %d %d %d *",
		scheduledFor.Minute(),
		scheduledFor.Hour(),
		scheduledFor.Day(),
		int(scheduledFor.Month()))

	_, err = s.cron.AddFunc(cronExpr, func() {
		s.executeDraftTask(ctx, taskID)
	})

	if err != nil {
		return fmt.Errorf("failed to schedule task: %w", err)
	}

	log.Printf("Draft %d scheduled at %s with task ID: %s", draftID, scheduledFor.Format(time.RFC3339), taskID)
	return nil
}
