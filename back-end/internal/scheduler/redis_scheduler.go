// internal/scheduler/redis_scheduler.go
package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"seo-backend/internal/database"
	"seo-backend/internal/helper"

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
	Status         string     `json:"status"`
	RetryCount     int        `json:"retry_count"`
	MaxRetries     int        `json:"max_retries"`
	CreatedAt      time.Time  `json:"created_at"`
	ExecutedAt     *time.Time `json:"executed_at,omitempty"`
	Error          string     `json:"error,omitempty"`
}

type RedisScheduler struct {
	redisClient  *redis.Client
	cron         *cron.Cron
	ctx          context.Context
	taskHandlers map[string]TaskHandler
}

type TaskHandler func(task *ScheduledTask) error

func NewRedisScheduler(redisClient *redis.Client) *RedisScheduler {
	return &RedisScheduler{
		redisClient:  redisClient,
		cron:         cron.New(cron.WithSeconds()),
		ctx:          context.Background(),
		taskHandlers: make(map[string]TaskHandler),
	}
}

func (s *RedisScheduler) RegisterTaskHandler(taskName string, handler TaskHandler) {
	s.taskHandlers[taskName] = handler
}

// =====================================================
// SCHEDULE DRAFT TASK
// =====================================================
func (s *RedisScheduler) ScheduleDraftTask(
	draftID string,
	scheduledFor time.Time,
	taskData *ScheduledTask,
) error {

	taskID := fmt.Sprintf(
		"draft_%s_%d",
		draftID,
		scheduledFor.Unix(),
	)

	task := &ScheduledTask{
		ID:             taskID,
		DraftID:        draftID,
		Title:          taskData.Title,
		Topic:          taskData.Topic,
		Article:        taskData.Article,
		ImageURL:       taskData.ImageURL,
		ImagePrompt:    taskData.ImagePrompt,
		TargetProducts: taskData.TargetProducts,
		ScheduledFor:   scheduledFor,
		TeamID:         taskData.TeamID,
		UserID:         taskData.UserID,
		Status:         "pending",
		RetryCount:     0,
		MaxRetries:     3,
		CreatedAt:      time.Now(),
	}

	taskKey := fmt.Sprintf("schedule:draft:%s", taskID)

	taskDataBytes, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("marshal task failed: %w", err)
	}

	// ================= LOGGING =================
	log.Println("===================================")
	log.Println("[SCHED] NEW TASK")
	log.Printf("[SCHED] taskID        : %s", taskID)
	log.Printf("[SCHED] draftID       : %s", draftID)
	log.Printf("[SCHED] scheduledFor  : %s", scheduledFor.Format(time.RFC3339))
	log.Printf("[SCHED] redisKey      : %s", taskKey)
	log.Printf("[SCHED] redis nil?    : %v", s.redisClient == nil)
	log.Println("===================================")

	err = s.redisClient.Set(
		s.ctx,
		taskKey,
		taskDataBytes,
		0,
	).Err()

	if err != nil {
		return fmt.Errorf("failed save task redis: %w", err)
	}

	log.Println("[SCHED] Task persisted to redis")

	// second minute hour day month weekday
	cronExpr := fmt.Sprintf(
		"0 %d %d %d %d *",
		scheduledFor.Minute(),
		scheduledFor.Hour(),
		scheduledFor.Day(),
		int(scheduledFor.Month()),
	)

	log.Printf("[CRON] Registering cron => %s", cronExpr)

	_, err = s.cron.AddFunc(
		cronExpr,
		func() {

			log.Println("===================================")
			log.Println("[CRON] CRON FIRED")
			log.Printf("[CRON] taskID : %s", taskID)
			log.Printf("[CRON] now    : %s", time.Now().Format(time.RFC3339))
			log.Println("===================================")

			s.executeDraftTask(taskID)
		},
	)

	if err != nil {
		return fmt.Errorf("cron registration failed: %w", err)
	}

	log.Printf("[SCHED] Task registered successfully %s", taskID)

	return nil
}

// =====================================================
// EXECUTE TASK
// =====================================================
func (s *RedisScheduler) executeDraftTask(taskID string) {

	log.Println("===================================")
	log.Println("[TASK] EXECUTION START")
	log.Printf("[TASK] taskID : %s", taskID)
	log.Println("===================================")

	task, err := s.getTask(taskID)
	if err != nil {
		log.Printf("[TASK] Redis fetch failed: %v", err)
		return
	}

	log.Printf("[TASK] Loaded DraftID=%s", task.DraftID)
	log.Printf("[TASK] Current status=%s", task.Status)

	task.Status = "processing"

	if err := s.updateTask(task); err != nil {
		log.Printf("[TASK] updateTask error: %v", err)
	}

	log.Println("[DB] Updating draft -> publishing")

	err = s.updateDraftStatus(
		task.DraftID,
		"publishing",
	)

	if err != nil {
		log.Printf(
			"[DB] Failed update draft status %s : %v",
			task.DraftID,
			err,
		)

		task.Status = "failed"
		task.Error = err.Error()
		s.updateTask(task)

		return
	}

	err = s.publishDraft(task)

	if err != nil {

		log.Println("===================================")
		log.Printf(
			"[TASK] PUBLISH FAILED draft=%s err=%v",
			task.DraftID,
			err,
		)
		log.Println("===================================")

		task.Status = "failed"
		task.Error = err.Error()

		s.updateDraftStatus(
			task.DraftID,
			"failed",
		)

	} else {

		log.Println("===================================")
		log.Printf(
			"[TASK] PUBLISH SUCCESS draft=%s",
			task.DraftID,
		)
		log.Println("===================================")

		task.Status = "completed"

		s.updateDraftStatus(
			task.DraftID,
			"published",
		)
	}

	now := time.Now()
	task.ExecutedAt = &now

	s.updateTask(task)

	log.Printf(
		"[TASK] Execution finished at %s",
		now.Format(time.RFC3339),
	)
}

// =====================================================
// RETRY PUBLISH
// =====================================================
func (s *RedisScheduler) publishDraft(
	task *ScheduledTask,
) error {

	var lastErr error

	log.Printf(
		"[PUBLISH] Starting publish retries max=%d",
		task.MaxRetries,
	)

	for i := 0; i <= task.MaxRetries; i++ {

		log.Printf(
			"[PUBLISH] Attempt %d/%d Draft=%s",
			i+1,
			task.MaxRetries+1,
			task.DraftID,
		)

		if i > 0 {
			time.Sleep(
				time.Duration(i*5) * time.Second,
			)
		}

		err := s.doPublishDraft(task)

		if err == nil {

			log.Printf(
				"[PUBLISH] Success on attempt=%d",
				i+1,
			)

			return nil
		}

		log.Printf(
			"[PUBLISH] Attempt failed: %v",
			err,
		)

		lastErr = err

		task.RetryCount = i + 1

		s.updateTask(task)
	}

	return fmt.Errorf(
		"max retries exceeded: %w",
		lastErr,
	)
}

// =====================================================
// ACTUAL PUBLISH
// =====================================================
func (s *RedisScheduler) doPublishDraft(
	task *ScheduledTask,
) error {

	log.Println("===================================")
	log.Printf(
		"[PUBLISH] Loading draft from DB %s",
		task.DraftID,
	)
	log.Println("===================================")

	var draft struct {
		Title          string
		Topic          string
		Article        string
		ImageURL       *string
		TargetProducts []string
		ImagePrompt    string
	}

	var targetProductsJSON []byte

	err := database.GetDB().QueryRow(`
	SELECT title, topic, article, image_url, target_products, image_prompt
	FROM drafts
	WHERE id = $1
`, task.DraftID).Scan(
		&draft.Title,
		&draft.Topic,
		&draft.Article,
		&draft.ImageURL,
		&targetProductsJSON,
		&draft.ImagePrompt,
	)

	if err != nil {
		return fmt.Errorf("failed get draft: %w", err)
	}

	// decode JSONB -> []string
	if len(targetProductsJSON) > 0 {
		err = json.Unmarshal(targetProductsJSON, &draft.TargetProducts)
		if err != nil {
			return fmt.Errorf("failed decode target_products: %w", err)
		}
	}

	log.Printf("[PUBLISH] Draft title=%s", draft.Title)
	log.Printf("[PUBLISH] Products=%v", draft.TargetProducts)

	postService := helper.NewPostService()

	log.Println("[PUBLISH] Sending to ProcessDraftProducts...")

	result, _, _, err :=
		postService.ProcessDraftProducts(draft)

	if err != nil {
		return fmt.Errorf(
			"process draft failed: %w",
			err,
		)
	}

	log.Printf(
		"[PUBLISH] Process result => %+v",
		result,
	)

	log.Println("[DB] Updating status -> published")

	_, err = database.GetDB().Exec(`
		UPDATE drafts
		SET status='published',
			published_at=$1,
			updated_at=$2
		WHERE id=$3
	`,
		time.Now(),
		time.Now(),
		task.DraftID,
	)

	if err != nil {
		return fmt.Errorf(
			"database update failed: %w",
			err,
		)
	}

	log.Println("[DB] Moving draft to histories")

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
			user_id
		)
		SELECT
			title,
			topic,
			article,
			image_url,
			target_products,
			'published',
			'auto_publish',
			$1,
			created_by,
			team_id,
			user_id
		FROM drafts
		WHERE id=$2
	`,
		time.Now(),
		task.DraftID,
	)

	if err != nil {
		log.Printf(
			"[DB] Warning history insert failed: %v",
			err,
		)
	}

	log.Println("[PUBLISH] Draft fully published")

	return nil
}

// =====================================================
// DB HELPERS
// =====================================================
func (s *RedisScheduler) updateDraftStatus(
	draftID string,
	status string,
) error {

	log.Printf(
		"[DB] updateDraftStatus draft=%s status=%s",
		draftID,
		status,
	)

	_, err := database.GetDB().Exec(`
		UPDATE drafts
		SET status=$1,
			updated_at=$2
		WHERE id=$3
	`,
		status,
		time.Now(),
		draftID,
	)

	return err
}

func (s *RedisScheduler) getTask(
	taskID string,
) (*ScheduledTask, error) {

	taskKey := fmt.Sprintf(
		"schedule:draft:%s",
		taskID,
	)

	log.Printf("[REDIS] GET %s", taskKey)

	taskData, err :=
		s.redisClient.Get(
			s.ctx,
			taskKey,
		).Bytes()

	if err != nil {
		return nil, err
	}

	var task ScheduledTask

	if err := json.Unmarshal(
		taskData,
		&task,
	); err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *RedisScheduler) updateTask(
	task *ScheduledTask,
) error {

	taskKey := fmt.Sprintf(
		"schedule:draft:%s",
		task.ID,
	)

	log.Printf(
		"[REDIS] UPDATE %s status=%s retry=%d",
		taskKey,
		task.Status,
		task.RetryCount,
	)

	taskData, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return s.redisClient.Set(
		s.ctx,
		taskKey,
		taskData,
		0,
	).Err()
}

// =====================================================
// START / STOP
// =====================================================
func (s *RedisScheduler) Start() {
	log.Println("===================================")
	log.Println("[BOOT] Starting Redis Scheduler...")
	log.Printf("[BOOT] scheduler=%p", s)
	log.Printf("[BOOT] redis nil? %v", s.redisClient == nil)
	log.Println("===================================")

	s.cron.Start()

	log.Println("[BOOT] Redis Scheduler started")
}

func (s *RedisScheduler) Stop() {
	s.cron.Stop()
	log.Println("[BOOT] Redis Scheduler stopped")
}
