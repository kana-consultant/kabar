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
	cancel            context.CancelFunc
	taskHandlers      map[string]TaskHandler
	db                *sql.DB
	productController product.ProductService
	postService       helper.PostService
	isRunning         bool
}

type TaskHandler func(task *ScheduledTask) error

func NewRedisScheduler(redisClient *redis.Client, db *sql.DB, productController product.ProductService, postService helper.PostService) *RedisScheduler {
	ctx, cancel := context.WithCancel(context.Background())

	return &RedisScheduler{
		redisClient:       redisClient,
		productController: productController,
		postService:       postService,
		cron:              cron.New(cron.WithSeconds()),
		ctx:               ctx,
		cancel:            cancel,
		taskHandlers:      make(map[string]TaskHandler),
		db:                db,
		isRunning:         false,
	}
}

func (s *RedisScheduler) RegisterTaskHandler(taskName string, handler TaskHandler) {
	s.taskHandlers[taskName] = handler
}

// ScheduleDraftTask schedules a draft for publishing
func (s *RedisScheduler) ScheduleDraftTask(draftID string, scheduledFor time.Time, taskData *ScheduledTask, userCtx models.UserContext) error {
	loc, _ := time.LoadLocation("Asia/Jakarta")

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

	// Save to Redis with TTL
	taskKey := fmt.Sprintf("schedule:draft:%s", taskID)
	taskDataBytes, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	// 🔥 FIX 1: Set TTL sesuai dengan waktu yang dijadwalkan
	ttl := time.Until(scheduledFor)
	if ttl < 0 {
		ttl = 0 // Sudah lewat, eksekusi segera
	}

	log.Printf("📝 Scheduling task: key=%s, ttl=%v, scheduledFor=%v", taskKey, ttl, scheduledFor)

	err = s.redisClient.Set(s.ctx, taskKey, taskDataBytes, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	// 🔥 FIX 2: Schedule dengan cron yang benar
	// Format: second minute hour day month weekday
	cronExpr := s.buildCronExpression(scheduledFor)

	log.Printf("🕐 Cron expression: %s", cronExpr)

	// Simpan userCtx ke Redis juga (untuk recovery)
	userCtxKey := fmt.Sprintf("%s:userctx", taskKey)
	userCtxBytes, _ := json.Marshal(userCtx)
	s.redisClient.Set(s.ctx, userCtxKey, userCtxBytes, ttl)

	entryID, err := s.cron.AddFunc(cronExpr, func() {
		s.executeDraftTask(taskID)
	})

	if err != nil {
		return fmt.Errorf("failed to schedule task: %w", err)
	}

	log.Printf("✅ Draft %s scheduled at %s | EntryID: %d | TaskID: %s",
		draftID,
		scheduledFor.Format(time.RFC3339),
		entryID,
		taskID)

	// Log semua cron entries untuk debugging
	s.logCronEntries()

	return nil
}

// 🔥 FIX 3: Build cron expression yang benar
func (s *RedisScheduler) buildCronExpression(t time.Time) string {
	// Format: second minute hour day month weekday
	// Untuk satu kali eksekusi di waktu tertentu
	return fmt.Sprintf("%d %d %d %d %d *",
		t.Second(),     // Detik
		t.Minute(),     // Menit
		t.Hour(),       // Jam
		t.Day(),        // Tanggal
		int(t.Month()), // Bulan
	)
}

// 🔥 FIX 4: Method untuk log cron entries
func (s *RedisScheduler) logCronEntries() {
	entries := s.cron.Entries()
	log.Printf("📊 Total cron entries: %d", len(entries))
	for i, entry := range entries {
		log.Printf("  Entry %d: ID=%d, Next=%v, Schedule=%s",
			i,
			entry.ID,
			entry.Next,
			entry.Schedule)
	}
}

// 🔥 FIX 5: Recovery dari Redis jika aplikasi restart
func (s *RedisScheduler) recoverPendingTasks() {
	log.Println("🔄 Recovering pending tasks from Redis...")

	pattern := "schedule:draft:draft_*"
	keys, err := s.redisClient.Keys(s.ctx, pattern).Result()
	if err != nil {
		log.Printf("❌ Failed to recover tasks: %v", err)
		return
	}

	if len(keys) == 0 {
		log.Println("📭 No pending tasks to recover")
		return
	}

	log.Printf("📦 Found %d pending tasks in Redis", len(keys))

	for _, key := range keys {
		// Ambil task dari Redis
		taskData, err := s.redisClient.Get(s.ctx, key).Bytes()
		if err != nil {
			log.Printf("❌ Failed to get task %s: %v", key, err)
			continue
		}

		var task ScheduledTask
		if err := json.Unmarshal(taskData, &task); err != nil {
			log.Printf("❌ Failed to unmarshal task %s: %v", key, err)
			continue
		}

		// Cek apakah task sudah lewat waktunya
		now := time.Now()
		if task.ScheduledFor.After(now) {
			// Belum waktunya, reschedule
			log.Printf("🔄 Rescheduling task %s for %v", task.ID, task.ScheduledFor)

			// Hitung ulang cron
			cronExpr := s.buildCronExpression(task.ScheduledFor)
			_, err := s.cron.AddFunc(cronExpr, func() {
				s.executeDraftTask(task.ID)
			})

			if err != nil {
				log.Printf("❌ Failed to reschedule task %s: %v", task.ID, err)
			} else {
				log.Printf("✅ Task %s rescheduled successfully", task.ID)
			}
		} else {
			// Sudah lewat, eksekusi sekarang
			log.Printf("⏰ Task %s is overdue, executing now...", task.ID)
			go s.executeDraftTask(task.ID)
		}
	}

	s.logCronEntries()
}

// executeDraftTask runs the scheduled draft publishing
func (s *RedisScheduler) executeDraftTask(taskID string) {
	log.Printf("🚀 Executing scheduled draft task: %s at %v", taskID, time.Now())

	// Get task from Redis
	task, err := s.getTask(taskID)
	if err != nil {
		log.Printf("❌ Failed to get task %s: %v", taskID, err)
		return
	}

	// 🔥 FIX 6: Ambil userCtx dari Redis
	userCtxKey := fmt.Sprintf("schedule:draft:%s:userctx", taskID)
	userCtxBytes, err := s.redisClient.Get(s.ctx, userCtxKey).Bytes()
	if err != nil {
		log.Printf("⚠️ Failed to get user context for task %s: %v", taskID, err)
		// Lanjutkan dengan empty context
	}

	var userCtx models.UserContext
	if err == nil {
		json.Unmarshal(userCtxBytes, &userCtx)
	}

	// Mark as processing
	task.Status = "processing"
	s.updateTask(task)

	// Update draft status to publishing
	err = s.updateDraftStatus(task.DraftID, "publishing")
	if err != nil {
		log.Printf("❌ Failed to update draft %s status: %v", task.DraftID, err)
		task.Status = "failed"
		task.Error = err.Error()
		s.updateTask(task)
		return
	}

	// Execute publishing with retry
	err = s.publishDraft(task, userCtx)
	if err != nil {
		log.Printf("❌ Failed to publish draft %s after retries: %v", task.DraftID, err)
		task.Status = "failed"
		task.Error = err.Error()
		s.updateDraftStatus(task.DraftID, "failed")
	} else {
		log.Printf("✅ Draft %s published successfully", task.DraftID)
		task.Status = "completed"
		s.updateDraftStatus(task.DraftID, "published")
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	task.ExecutedAt = &now
	s.updateTask(task)

	// 🔥 FIX 7: Hapus task dari Redis setelah selesai (opsional)
	// s.redisClient.Del(s.ctx, fmt.Sprintf("schedule:draft:%s", taskID))
	// s.redisClient.Del(s.ctx, fmt.Sprintf("schedule:draft:%s:userctx", taskID))
}

// publishDraft with retry mechanism
func (s *RedisScheduler) publishDraft(task *ScheduledTask, userCtx models.UserContext) error {
	var lastErr error

	for i := 0; i <= task.MaxRetries; i++ {
		if i > 0 {
			log.Printf("🔄 Retrying draft %s (attempt %d/%d)", task.DraftID, i+1, task.MaxRetries+1)
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	result, someFailed, allFailed, err := s.postService.ProcessDraftProducts(ctx, draftData, userCtx)
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

	now := helper.ParseWIBTime(time.Now().Format(time.RFC3339))

	_, err = database.GetDB().Exec(`
		UPDATE drafts 
		SET status = $1,
			published_at = $2,
			updated_at = $3
		WHERE id = $4
	`,
		status,
		now,
		now,
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
		now,
		nil,
		task.TeamID,
		task.UserID,
		errorMessage,
	)

	if err != nil {
		log.Printf("Warning: Failed to insert history: %v", err)
	}

	log.Printf(
		"[Scheduler] ✅ SUCCESS published draft_id=%s status=%s",
		task.DraftID,
		status,
	)

	return nil
}

// updateDraftStatus updates draft status in database
func (s *RedisScheduler) updateDraftStatus(draftID string, status string) error {
	now := helper.ParseWIBTime(time.Now().Format(time.RFC3339))
	_, err := database.GetDB().Exec(`
		UPDATE drafts 
		SET status = $1, updated_at = $2
		WHERE id = $3
	`, status, now, draftID)
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

// 🔥 FIX 8: Start dengan recovery
func (s *RedisScheduler) Start() {
	if s.isRunning {
		log.Println("⚠️ Scheduler already running")
		return
	}

	s.cron.Start()
	s.isRunning = true
	log.Println("✅ Redis Scheduler started")

	// Recovery pending tasks
	go s.recoverPendingTasks()
}

// Stop stops the scheduler
func (s *RedisScheduler) Stop() {
	if !s.isRunning {
		return
	}

	s.cron.Stop()
	s.isRunning = false
	s.cancel()
	log.Println("🛑 Redis Scheduler stopped")
}

// 🔥 FIX 9: Method untuk cleanup expired tasks
func (s *RedisScheduler) CleanupExpiredTasks() error {
	pattern := "schedule:draft:*"
	keys, err := s.redisClient.Keys(s.ctx, pattern).Result()
	if err != nil {
		return err
	}

	var deleted int

	for _, key := range keys {
		// Ambil TTL
		ttl, err := s.redisClient.TTL(s.ctx, key).Result()
		if err != nil {
			continue
		}

		// Jika TTL <= 0 dan key masih ada, berarti sudah expired
		if ttl <= 0 {
			// Cek apakah key masih ada
			exists, _ := s.redisClient.Exists(s.ctx, key).Result()
			if exists == 0 {
				continue
			}

			// Hapus jika sudah expired
			s.redisClient.Del(s.ctx, key)
			deleted++
		}
	}

	if deleted > 0 {
		log.Printf("🧹 Cleaned up %d expired tasks", deleted)
	}

	return nil
}
