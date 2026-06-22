// internal/scheduler/redis_scheduler.go
package scheduler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"seo-backend/internal/database"
	"seo-backend/internal/domain/draft"
	"seo-backend/internal/domain/product"
	"seo-backend/internal/helper"
	"seo-backend/internal/models"

	"time"

	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
)

type ScheduledTask struct {
	ID             string     `json:"id"`
	DraftID        string     `json:"draft_id"`
	Title          string     `json:"title"`
	Topic          string     `json:"topic"`
	Article        string     `json:"article"`
	ImageURL       string     `json:"image_url"`
	ImagePrompt    string     `json:"image_prompt"`
	TargetProducts []string   `json:"target_products"`
	ScheduledFor   time.Time  `json:"scheduled_for"`
	TeamID         string     `json:"team_id"`
	UserID         string     `json:"user_id"`
	Status         string     `json:"status"` // pending, processing, completed, failed
	RetryCount     int        `json:"retry_count"`
	MaxRetries     int        `json:"max_retries"`
	CreatedAt      time.Time  `json:"created_at"`
	ExecutedAt     *time.Time `json:"executed_at,omitempty"`
	Error          string     `json:"error,omitempty"`
}

type RedisScheduler struct {
	redisClient       *redis.Client
	cron              *cron.Cron
	ctx               context.Context
	taskHandlers      map[string]TaskHandler
	db                *sql.DB
	productController product.ProductService
}

type TaskHandler func(task *ScheduledTask) error

func NewRedisScheduler(redisClient *redis.Client, db *sql.DB, productController product.ProductService) *RedisScheduler {
	return &RedisScheduler{
		redisClient:       redisClient,
		productController: productController,
		cron:              cron.New(cron.WithSeconds()),
		ctx:               context.Background(),
		taskHandlers:      make(map[string]TaskHandler),
		db:                db,
	}
}

func (s *RedisScheduler) RegisterTaskHandler(taskName string, handler TaskHandler) {
	s.taskHandlers[taskName] = handler
}

// ScheduleDraftTask schedules a draft for publishing
func (s *RedisScheduler) ScheduleDraftTask(draftID string, scheduledFor time.Time, taskData *ScheduledTask, userCtx models.UserContext) error {
	taskID := fmt.Sprintf(
		"draft_%s_%d",
		draftID,
		scheduledFor.Unix(),
	)
	loc, _ := time.LoadLocation("Asia/Jakarta")

	task := &ScheduledTask{
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
		CreatedAt:      time.Now().In(loc),
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

	err = s.redisClient.Set(s.ctx, taskKey, taskDataBytes, 0).Err()
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
		s.executeDraftTask(taskID, userCtx)
	})

	if err != nil {
		return fmt.Errorf("failed to schedule task: %w", err)
	}

	log.Printf("Draft %s scheduled at %s with task ID: %s", draftID, scheduledFor.Format(time.RFC3339), taskID)
	return nil
}

// executeDraftTask runs the scheduled draft publishing
func (s *RedisScheduler) executeDraftTask(taskID string, userCtx models.UserContext) {
	log.Printf("Executing scheduled draft task: %s", taskID)

	// Get task from Redis
	task, err := s.getTask(taskID)
	if err != nil {
		log.Printf("Failed to get task %s: %v", taskID, err)
		return
	}

	// Mark as processing
	task.Status = "processing"
	s.updateTask(task)

	// Update draft status to publishing
	err = s.updateDraftStatus(task.DraftID, "publishing")
	if err != nil {
		log.Printf("Failed to update draft %s status: %v", task.DraftID, err)
		task.Status = "failed"
		task.Error = err.Error()
		s.updateTask(task)
		return
	}

	// Execute publishing with retry
	err = s.publishDraft(task, userCtx)
	if err != nil {
		log.Printf("Failed to publish draft %s after retries: %v", task.DraftID, err)
		task.Status = "failed"
		task.Error = err.Error()
		s.updateDraftStatus(task.DraftID, "failed")
	} else {
		log.Printf("Draft %s published successfully", task.DraftID)
		task.Status = "completed"
		s.updateDraftStatus(task.DraftID, "published")
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")

	now := time.Now().In(loc)
	task.ExecutedAt = &now
	s.updateTask(task)
}

// publishDraft with retry mechanism
func (s *RedisScheduler) publishDraft(task *ScheduledTask, userCtx models.UserContext) error {
	var lastErr error

	for i := 0; i <= task.MaxRetries; i++ {
		if i > 0 {
			log.Printf("Retrying draft %s (attempt %d/%d)", task.DraftID, i+1, task.MaxRetries+1)
			time.Sleep(time.Duration(i*5) * time.Second)
		}

		err := s.doPublishDraft(task, userCtx)
		if err == nil {
			return nil
		}

		lastErr = err
		task.RetryCount = i + 1
		s.updateTask(task)
	}

	return fmt.Errorf("max retries exceeded: %w", lastErr)
}

// doPublishDraft actually publishes the draft
func (s *RedisScheduler) doPublishDraft(task *ScheduledTask, userCtx models.UserContext) error {

	log.Printf(
		"[Scheduler] START publish draft_id=%s title=%s products=%v",
		task.DraftID,
		task.Title,
		task.TargetProducts,
	)

	// langsung gunakan data dari task redis
	draftData := draft.DraftDataPost{
		Title:          task.Title,
		Topic:          task.Topic,
		Article:        task.Article,
		ImageURL:       &task.ImageURL,
		ImagePrompt:    task.ImagePrompt,
		TargetProducts: task.TargetProducts,
	}

	postService := helper.NewPostService(s.db, s.productController)

	result, someFailed, allFailed, err := postService.ProcessDraftProducts(s.ctx, draftData, userCtx)
	if err != nil {
		return fmt.Errorf("failed to process products: %w", err)
	}

	log.Printf(
		"[PROCESS RESULT] someFailed=%v allFailed=%v result=%+v",
		someFailed,
		allFailed,
		result,
	)

	status := "published"

	if allFailed {
		status = "failed"
	}

	_, err = database.GetDB().Exec(`
        UPDATE drafts 
        SET status = $1,
            published_at = $2,
            updated_at = $3
        WHERE id = $4
    `,
		status,
		helper.ParseWIBTime(time.Now().Format(time.RFC3339)),
		helper.ParseWIBTime(time.Now().Format(time.RFC3339)),
		task.DraftID,
	)

	if err != nil {
		return fmt.Errorf("database update failed: %w", err)
	}

	targetProductsJSON, _ := json.Marshal(task.TargetProducts)

	errorMessage, _ := json.Marshal(map[string]interface{}{
		"all_failed":  allFailed,
		"some_failed": someFailed,
	})

	publishedAt := helper.ParseWIBTime(time.Now().Format(time.RFC3339))

	_, err = database.GetDB().Exec(`
    INSERT INTO histories (
        title,
        topic,
        content,
        image_url,
        target_products,
        status,
        action,
        published_at,
        created_by,
        team_id,
        user_id,
        error_message
    )
    VALUES (
        $1, $2, $3, $4, $5,
        $6, 'auto_publish', $7, $8, $9, $10, $11
    )
`,
		task.Title,
		task.Topic,
		task.Article,
		task.ImageURL,
		targetProductsJSON,
		status,
		publishedAt,
		nil,
		task.TeamID,
		task.UserID,
		errorMessage,
	)

	if err != nil {
		log.Printf(
			"Warning: Failed to insert history: %v",
			err,
		)
	}

	log.Printf(
		"[Scheduler] SUCCESS published draft_id=%s status=%s",
		task.DraftID,
		status,
	)

	return nil
}

// updateDraftStatus updates draft status in database
func (s *RedisScheduler) updateDraftStatus(draftID string, status string) error {

	_, err := database.GetDB().Exec(`
		UPDATE drafts 
		SET status = $1, updated_at = $2
		WHERE id = $3
	`, status, helper.ParseWIBTime(time.Now().Format(time.RFC3339)), draftID)

	return err
}

func (s *RedisScheduler) getTask(taskID string) (*ScheduledTask, error) {
	taskKey := fmt.Sprintf("schedule:draft:%s", taskID)
	taskData, err := s.redisClient.Get(s.ctx, taskKey).Bytes()
	if err != nil {
		return nil, err
	}

	var task ScheduledTask
	if err := json.Unmarshal(taskData, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *RedisScheduler) updateTask(task *ScheduledTask) error {
	taskKey := fmt.Sprintf("schedule:draft:%s", task.ID)
	taskData, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return s.redisClient.Set(s.ctx, taskKey, taskData, 0).Err()
}

// CancelScheduledTask cancels a scheduled draft
func (s *RedisScheduler) CancelScheduledTask(draftID string) error {
	pattern := fmt.Sprintf("schedule:draft:draft_%s_*", draftID)
	keys, err := s.redisClient.Keys(s.ctx, pattern).Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		if err := s.redisClient.Del(s.ctx, key).Err(); err != nil {
			log.Printf("Failed to delete task key %s: %v", key, err)
		}
	}

	log.Printf("Cancelled scheduled tasks for draft %s", draftID)
	return nil
}

// GetScheduledTasks gets all scheduled drafts
func (s *RedisScheduler) GetScheduledTasks() ([]*ScheduledTask, error) {
	pattern := "schedule:draft:*"
	keys, err := s.redisClient.Keys(s.ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	var tasks []*ScheduledTask
	for _, key := range keys {
		taskData, err := s.redisClient.Get(s.ctx, key).Bytes()
		if err != nil {
			continue
		}

		var task ScheduledTask
		if err := json.Unmarshal(taskData, &task); err != nil {
			continue
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}

// Start starts the scheduler
func (s *RedisScheduler) Start() {
	s.cron.Start()
	log.Println("Redis Scheduler started")
}

// Stop stops the scheduler
func (s *RedisScheduler) Stop() {
	s.cron.Stop()
	log.Println("Redis Scheduler stopped")
}
