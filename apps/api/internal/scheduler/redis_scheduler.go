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
	"strings"
	"sync"
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
	mu                sync.RWMutex
	activeTasks       map[string]context.CancelFunc
	stopChan          chan struct{} // Channel untuk graceful shutdown
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
		activeTasks:       make(map[string]context.CancelFunc),
		stopChan:          make(chan struct{}),
	}
}

func (s *RedisScheduler) RegisterTaskHandler(taskName string, handler TaskHandler) {
	s.taskHandlers[taskName] = handler
}

// ScheduleDraftTask schedules a draft for publishing
func (s *RedisScheduler) ScheduleDraftTask(
	ctx context.Context,
	draftID string,
	scheduledFor time.Time,
	taskData *ScheduledTask,
	userCtx models.UserContext,
) error {
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

	taskKey := fmt.Sprintf("schedule:draft:%s", taskID)
	taskDataBytes, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	ttl := time.Until(scheduledFor)
	if ttl < 0 {
		ttl = 0
	}

	log.Printf("📝 Scheduling task: key=%s, ttl=%v, scheduledFor=%v", taskKey, ttl, scheduledFor)

	// Gunakan context dari parameter
	err = s.redisClient.Set(ctx, taskKey, taskDataBytes, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	cronExpr := s.buildCronExpression(scheduledFor)
	log.Printf("🕐 Cron expression: %s", cronExpr)

	userCtxKey := fmt.Sprintf("%s:userctx", taskKey)
	userCtxBytes, _ := json.Marshal(userCtx)
	s.redisClient.Set(ctx, userCtxKey, userCtxBytes, ttl)

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

	s.logCronEntries()
	return nil
}

// buildCronExpression creates cron expression for specific time
func (s *RedisScheduler) buildCronExpression(t time.Time) string {
	return fmt.Sprintf("%d %d %d %d %d *",
		t.Second(),
		t.Minute(),
		t.Hour(),
		t.Day(),
		int(t.Month()),
	)
}

// logCronEntries logs all cron entries
func (s *RedisScheduler) logCronEntries() {
	entries := s.cron.Entries()
	log.Printf("📊 Total cron entries: %d", len(entries))
	for i, entry := range entries {
		log.Printf("  Entry %d: ID=%d, Next=%v",
			i,
			entry.ID,
			entry.Next,
		)
	}
}

// recoverPendingTasks recovers pending tasks from Redis if application restarts
func (s *RedisScheduler) recoverPendingTasks() {
	log.Println("🔄 Recovering pending tasks from Redis...")

	// GUNAKAN CONTEXT BACKGROUND AGAR TIDAK TERGANTUNG s.ctx
	ctx := context.Background()

	// Hanya ambil key yang sesuai format schedule:draft:draft_*
	pattern := "schedule:draft:draft_*"
	keys, err := s.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		log.Printf("❌ Failed to recover tasks: %v", err)
		return
	}

	// Filter key yang bukan userctx
	var taskKeys []string
	for _, key := range keys {
		if !strings.Contains(key, ":userctx") {
			taskKeys = append(taskKeys, key)
		}
	}

	if len(taskKeys) == 0 {
		log.Println("📭 No pending tasks to recover")
		return
	}

	log.Printf("📦 Found %d pending tasks in Redis", len(taskKeys))

	for _, key := range taskKeys {
		// Ambil task dari Redis dengan context background
		taskData, err := s.redisClient.Get(ctx, key).Bytes()
		if err != nil {
			log.Printf("❌ Failed to get task %s: %v", key, err)
			continue
		}

		var task ScheduledTask
		if err := json.Unmarshal(taskData, &task); err != nil {
			log.Printf("❌ Failed to unmarshal task %s: %v", key, err)
			continue
		}

		// Validasi task ID tidak kosong
		if task.ID == "" {
			log.Printf("⚠️ Task with empty ID found, deleting key: %s", key)
			s.redisClient.Del(ctx, key)
			s.redisClient.Del(ctx, key+":userctx")
			continue
		}

		// Cek apakah task sudah lewat waktunya
		now := time.Now()
		if task.ScheduledFor.After(now) {
			// Belum waktunya, reschedule
			log.Printf("🔄 Rescheduling task %s for %v", task.ID, task.ScheduledFor)

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
			// Sudah lewat, eksekusi sekarang dengan goroutine
			log.Printf("⏰ Task %s is overdue, executing now...", task.ID)
			go s.executeDraftTask(task.ID)
		}
	}

	s.logCronEntries()
}

// executeDraftTask runs the scheduled draft publishing
func (s *RedisScheduler) executeDraftTask(taskID string) {
	// CEK APAKAH SCHEDULER MASIH RUNNING
	if !s.isRunning {
		log.Printf("⚠️ Scheduler is not running, skipping task %s", taskID)
		return
	}

	log.Printf("🚀 Executing scheduled draft task: %s at %v", taskID, time.Now())

	// Validasi taskID tidak kosong
	if taskID == "" {
		log.Printf("❌ Empty taskID, skipping execution")
		return
	}

	// Gunakan context dengan timeout untuk mencegah context canceled
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Track active task
	s.mu.Lock()
	s.activeTasks[taskID] = cancel
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.activeTasks, taskID)
		s.mu.Unlock()
	}()

	// Get task from Redis dengan context baru
	task, err := s.getTaskWithContext(ctx, taskID)
	if err != nil {
		log.Printf("❌ Failed to get task %s: %v", taskID, err)
		return
	}

	// Ambil userCtx dari Redis
	userCtxKey := fmt.Sprintf("schedule:draft:%s:userctx", taskID)
	userCtxBytes, err := s.redisClient.Get(ctx, userCtxKey).Bytes()
	if err != nil {
		log.Printf("⚠️ Failed to get user context for task %s: %v", taskID, err)
	}

	var userCtx models.UserContext
	if err == nil {
		json.Unmarshal(userCtxBytes, &userCtx)
	}

	// PENTING: Bersihkan cache Redis di akhir, apapun hasilnya
	defer func() {
		taskKey := fmt.Sprintf("schedule:draft:%s", taskID)
		userCtxKey := fmt.Sprintf("schedule:draft:%s:userctx", taskID)

		// Gunakan context background untuk cleanup
		cleanupCtx := context.Background()
		delCount, err := s.redisClient.Del(cleanupCtx, taskKey, userCtxKey).Result()
		if err != nil {
			log.Printf("⚠️ Failed to clean Redis cache for task %s: %v", taskID, err)
		} else {
			log.Printf("🧹 Cleaned Redis cache for task %s (deleted %d keys)", taskID, delCount)
		}
	}()

	// Mark as processing
	task.Status = "processing"
	s.updateTaskWithContext(ctx, task)

	// Update draft status to publishing
	err = s.updateDraftStatus(task.DraftID, "publishing")
	if err != nil {
		log.Printf("❌ Failed to update draft %s status: %v", task.DraftID, err)
		task.Status = "failed"
		task.Error = err.Error()
		s.updateTaskWithContext(ctx, task)
		return
	}

	// Execute publishing with retry
	err = s.publishDraftWithContext(ctx, task, userCtx)
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
	s.updateTaskWithContext(ctx, task)
}

// publishDraftWithContext with retry mechanism and context
func (s *RedisScheduler) publishDraftWithContext(ctx context.Context, task *ScheduledTask, userCtx models.UserContext) error {
	var lastErr error

	for i := 0; i <= task.MaxRetries; i++ {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
		}

		if i > 0 {
			log.Printf("🔄 Retrying draft %s (attempt %d/%d)", task.DraftID, i+1, task.MaxRetries+1)

			// Sleep with context awareness
			select {
			case <-ctx.Done():
				return fmt.Errorf("context cancelled during retry: %w", ctx.Err())
			case <-time.After(time.Duration(i*5) * time.Second):
			}
		}

		err := s.doPublishDraft(task, userCtx)
		if err == nil {
			return nil
		}

		lastErr = err
		task.RetryCount = i + 1
		s.updateTaskWithContext(ctx, task)
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

// getTask retrieves task from Redis using scheduler context
func (s *RedisScheduler) getTask(taskID string) (*ScheduledTask, error) {
	return s.getTaskWithContext(s.ctx, taskID)
}

// getTaskWithContext retrieves task from Redis with specific context
func (s *RedisScheduler) getTaskWithContext(ctx context.Context, taskID string) (*ScheduledTask, error) {
	taskKey := fmt.Sprintf("schedule:draft:%s", taskID)
	taskData, err := s.redisClient.Get(ctx, taskKey).Bytes()
	if err != nil {
		return nil, err
	}

	var task ScheduledTask
	if err := json.Unmarshal(taskData, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

// updateTask updates task in Redis using scheduler context
func (s *RedisScheduler) updateTask(task *ScheduledTask) error {
	return s.updateTaskWithContext(s.ctx, task)
}

// updateTaskWithContext updates task in Redis with specific context
func (s *RedisScheduler) updateTaskWithContext(ctx context.Context, task *ScheduledTask) error {
	taskKey := fmt.Sprintf("schedule:draft:%s", task.ID)
	taskData, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return s.redisClient.Set(ctx, taskKey, taskData, 0).Err()
}

// CancelScheduledTask cancels a scheduled draft
func (s *RedisScheduler) CancelScheduledTask(ctx context.Context, draftID string) error {
	pattern := fmt.Sprintf("schedule:draft:draft_%s_*", draftID)
	keys, err := s.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		// Hapus task key
		if err := s.redisClient.Del(ctx, key).Err(); err != nil {
			log.Printf("Failed to delete task key %s: %v", key, err)
		}
		// Hapus userctx key
		s.redisClient.Del(ctx, key+":userctx")
	}

	log.Printf("Cancelled scheduled tasks for draft %s", draftID)
	return nil
}

// GetScheduledTasks gets all scheduled drafts
func (s *RedisScheduler) GetScheduledTasks() ([]*ScheduledTask, error) {
	// Gunakan context background untuk operasi read
	ctx := context.Background()

	pattern := "schedule:draft:*"
	keys, err := s.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	var tasks []*ScheduledTask
	for _, key := range keys {
		// Skip userctx keys
		if strings.Contains(key, ":userctx") {
			continue
		}

		taskData, err := s.redisClient.Get(ctx, key).Bytes()
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
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		log.Println("⚠️ Scheduler already running")
		return
	}

	log.Println("✅ Redis Scheduler starting...")

	// Start cron
	s.cron.Start()
	s.isRunning = true
	log.Println("✅ Redis Scheduler started")

	// Recovery pending tasks dengan delay kecil agar server siap
	go func() {
		time.Sleep(2 * time.Second)
		s.recoverPendingTasks()
	}()

	// Start health check routine
	go s.healthCheck()
}

// healthCheck melakukan pengecekan kesehatan scheduler secara periodik
func (s *RedisScheduler) healthCheck() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.mu.RLock()
			if s.isRunning {
				entries := s.cron.Entries()
				log.Printf("💚 Scheduler health check: running, %d active cron entries", len(entries))
			}
			s.mu.RUnlock()
		}
	}
}

// Stop stops the scheduler gracefully
func (s *RedisScheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		log.Println("⚠️ Scheduler already stopped")
		return
	}

	log.Println("🛑 Stopping Redis Scheduler...")

	// Set isRunning ke false dulu agar tidak ada task baru yang dijalankan
	s.isRunning = false

	// Cancel semua active tasks
	s.mu.Unlock()
	s.mu.RLock()
	for taskID, cancel := range s.activeTasks {
		log.Printf("⚠️ Cancelling active task: %s", taskID)
		cancel()
	}
	s.mu.RUnlock()
	s.mu.Lock()

	// Tunggu sebentar untuk menyelesaikan task yang sedang berjalan
	time.Sleep(2 * time.Second)

	// Stop cron
	<-s.cron.Stop().Done()

	// Cancel context
	s.cancel()

	log.Println("🛑 Redis Scheduler stopped")
}

// CleanupExpiredTasks cleans up expired tasks from Redis
func (s *RedisScheduler) CleanupExpiredTasks() error {
	// Gunakan context background
	ctx := context.Background()

	pattern := "schedule:draft:*"
	keys, err := s.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	var deleted int

	for _, key := range keys {
		// Skip userctx keys, akan dihapus bersama task key
		if strings.Contains(key, ":userctx") {
			continue
		}

		// Ambil TTL
		ttl, err := s.redisClient.TTL(ctx, key).Result()
		if err != nil {
			continue
		}

		// Jika TTL <= 0 dan key masih ada, berarti sudah expired
		if ttl <= 0 {
			// Cek apakah key masih ada
			exists, _ := s.redisClient.Exists(ctx, key).Result()
			if exists == 0 {
				continue
			}

			// Hapus task key dan userctx key
			s.redisClient.Del(ctx, key, key+":userctx")
			deleted++
		}
	}

	if deleted > 0 {
		log.Printf("🧹 Cleaned up %d expired tasks", deleted)
	}

	return nil
}

// IsRunning returns whether scheduler is running
func (s *RedisScheduler) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isRunning
}
