package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"seo-backend/internal/builder"
	"seo-backend/internal/database"
	"seo-backend/internal/helper"
	"seo-backend/internal/middleware/auth"
	"seo-backend/internal/models"
	"seo-backend/internal/scheduler"
)

type DraftHandler struct{}

func NewDraftHandler() *DraftHandler {
	return &DraftHandler{}
}

var redisScheduler *scheduler.RedisScheduler

// InitScheduler initializes the scheduler when app starts
func InitScheduler() {
	redisScheduler = scheduler.NewRedisScheduler(database.RedisClient)
	redisScheduler.Start()
}

// GetAll drafts - menggunakan query builder
func (h *DraftHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	teamID := auth.GetTeamID(ctx)
	userRole := auth.GetUserRole(ctx)
	userID := auth.GetUserID(ctx)

	status := r.URL.Query().Get("status")
	search := r.URL.Query().Get("search")

	// Build query dengan generic builder
	qb := builder.NewQueryBuilder("drafts")

	query := qb.Select(
		"id", "title", "topic", "article", "image_url", "image_prompt",
		"status", "scheduled_for", "target_products", "has_image",
		"team_id", "user_id", "created_at", "updated_at",
	).OrderBy("created_at DESC")

	// Logic filter berdasarkan role
	switch userRole {
	case "super_admin":
		// Super admin: melihat SEMUA draft
		log.Printf("Super admin - melihat semua draft")

	case "admin":
		// Admin: hanya melihat draft dalam team yang sama
		if teamID != "" {
			query = query.WhereEq("team_id", teamID)
			log.Printf("Admin - melihat draft dalam team: %s", teamID)
		} else {
			// Admin tanpa team hanya melihat draft miliknya sendiri
			query = query.WhereEq("user_id", userID)
			log.Printf("Admin tanpa team - hanya melihat draft sendiri")
		}

	default:
		// Role lain (manager, editor, viewer): hanya melihat draft milik sendiri
		query = query.WhereEq("user_id", userID)
		log.Printf("User role %s - hanya melihat draft sendiri", userRole)
	}

	// Filter tambahan
	if status != "" {
		query = query.WhereEq("status", status)
	}

	if search != "" {
		query = query.WhereLike("title", search)
	}

	sqlQuery, args, err := query.Build()
	if err != nil {
		log.Printf("Failed to build query: %v", err)
		http.Error(w, "Failed to fetch drafts", http.StatusInternalServerError)
		return
	}

	log.Printf("Query: %s | Args: %v", sqlQuery, args)

	rows, err := database.GetDB().Query(sqlQuery, args...)
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Failed to fetch drafts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var drafts []models.Draft
	for rows.Next() {
		var d models.Draft
		var targetProductsJSON []byte

		err := rows.Scan(
			&d.ID, &d.Title, &d.Topic, &d.Article, &d.ImageURL, &d.ImagePrompt,
			&d.Status, &d.ScheduledFor, &targetProductsJSON, &d.HasImage,
			&d.TeamID, &d.UserID, &d.CreatedAt, &d.UpdatedAt,
		)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}

		json.Unmarshal(targetProductsJSON, &d.TargetProducts)
		drafts = append(drafts, d)
	}

	log.Printf("Successfully fetched %d drafts", len(drafts))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drafts)
}

// GetByID draft
func (h *DraftHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	qb := builder.NewQueryBuilder("drafts")

	sqlQuery, args, err := qb.Select(
		"id", "title", "topic", "article", "image_url", "image_prompt",
		"status", "scheduled_for", "target_products", "has_image",
		"team_id", "user_id", "created_at", "updated_at",
	).WhereEq("id", id).Build()

	if err != nil {
		log.Printf("Failed to build query: %v", err)
		http.Error(w, "Failed to fetch draft", http.StatusInternalServerError)
		return
	}

	var draft models.Draft
	var targetProductsJSON []byte

	err = database.GetDB().QueryRow(sqlQuery, args...).Scan(
		&draft.ID, &draft.Title, &draft.Topic, &draft.Article,
		&draft.ImageURL, &draft.ImagePrompt, &draft.Status,
		&draft.ScheduledFor, &targetProductsJSON, &draft.HasImage,
		&draft.TeamID, &draft.UserID, &draft.CreatedAt, &draft.UpdatedAt,
	)
	if err != nil {
		http.Error(w, "Draft not found", http.StatusNotFound)
		return
	}

	json.Unmarshal(targetProductsJSON, &draft.TargetProducts)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(draft)
}

// Create draft
// internal/handlers/draft_handler.go
func (h *DraftHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := auth.GetUserID(ctx)
	teamID := auth.GetTeamID(ctx)

	var req models.CreateDraftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Title == "" || req.Topic == "" || req.Article == "" {
		http.Error(w, "Title, topic, and article are required", http.StatusBadRequest)
		return
	}

	targetProductsJSON, _ := json.Marshal(req.TargetProducts)

	// Handle UUID null values
	var createdBy interface{}
	if userID != "" && userID != "00000000-0000-0000-0000-000000000000" {
		createdBy = userID
	} else {
		createdBy = nil
	}

	var teamIDPtr interface{}
	if teamID != "" && teamID != "00000000-0000-0000-0000-000000000000" {
		teamIDPtr = teamID
	} else {
		teamIDPtr = nil
	}

	qb := builder.NewQueryBuilder("drafts")

	sqlQuery, args, err := qb.Insert().
		Columns(
			"title", "topic", "article", "image_url", "image_prompt",
			"status", "target_products", "has_image", "created_by", "team_id", "user_id",
		).
		Values(
			req.Title, req.Topic, req.Article, req.ImageURL, req.ImagePrompt,
			"draft", targetProductsJSON, req.HasImage, createdBy, teamIDPtr, createdBy,
		).
		Returning("id").
		Build()

	if err != nil {
		log.Printf("Failed to build insert query: %v", err)
		http.Error(w, "Failed to create draft", http.StatusInternalServerError)
		return
	}

	log.Printf("Insert query: %s", sqlQuery)
	log.Printf("Args: %v", args)

	var draftID string
	err = database.GetDB().QueryRow(sqlQuery, args...).Scan(&draftID)
	if err != nil {
		log.Printf("Failed to create draft: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create draft: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Draft created with ID: %s", draftID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id":      draftID,
		"message": "Draft created successfully",
	})
}

// Update draft
func (h *DraftHandler) Update(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Map field names dari JSON ke database
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
			} else if key == "scheduledFor" && value != nil {
				// Parse scheduledFor string to time.Time
				if scheduledStr, ok := value.(string); ok && scheduledStr != "" {
					// Handle daily schedule format
					if strings.HasPrefix(scheduledStr, "daily:") {
						// Store as string or handle differently
						data[dbField] = scheduledStr
					} else {
						// Parse ISO timestamp
						parsedTime, err := time.Parse(time.RFC3339, scheduledStr)
						if err == nil {
							data[dbField] = parsedTime
						} else {
							data[dbField] = value
						}
					}
				} else {
					data[dbField] = nil
				}
			} else {
				data[dbField] = value
			}
		}
	}
	data["updated_at"] = time.Now()

	if len(data) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	sqlQuery, args, err := buildUpdateByIDQuery("drafts", id, data)

	if err != nil {
		http.Error(w, "Draft not found", http.StatusNotFound)
		return
	}

	log.Printf("Update query: %s | Args: %v", sqlQuery, args)

	_, err = database.GetDB().Exec(sqlQuery, args...)
	if err != nil {
		log.Printf("Failed to update draft: %v", err)
		http.Error(w, "Failed to update draft", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Draft updated successfully"})
}

// Delete draft
func (h *DraftHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Check permission berdasarkan role
	var draftTeamID *string
	var draftUserID *string
	err := database.GetDB().QueryRow(`
		SELECT team_id, user_id FROM drafts WHERE id = $1
	`, id).Scan(&draftTeamID, &draftUserID)

	if err != nil {
		http.Error(w, "Draft not found", http.StatusNotFound)
		return
	}

	qb := builder.NewQueryBuilder("drafts")

	sqlQuery, args, err := qb.Delete().
		WhereEq("id", id).
		Build()

	if err != nil {
		log.Printf("Failed to build delete query: %v", err)
		http.Error(w, "Failed to delete draft", http.StatusInternalServerError)
		return
	}

	log.Printf("Delete query: %s | Args: %v", sqlQuery, args)

	_, err = database.GetDB().Exec(sqlQuery, args...)
	if err != nil {
		log.Printf("Failed to delete draft: %v", err)
		http.Error(w, "Failed to delete draft", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Publish draft
// internal/handlers/draft_handler.go

func (h *DraftHandler) Publish(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()

	teamID := auth.GetTeamID(ctx)
	userID := auth.GetUserID(ctx)

	var req struct {
		ScheduledFor string `json:"scheduledFor,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Publishing draft %s with scheduledFor: %s", id, req.ScheduledFor)

	// =========================
	// GET DRAFT DATA
	// =========================
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
		SELECT title, topic, article, image_url, target_products, COALESCE(image_prompt, '')
		FROM drafts WHERE id = $1
	`, id).Scan(
		&draft.Title,
		&draft.Topic,
		&draft.Article,
		&draft.ImageURL,
		&targetProductsJSON,
		&draft.ImagePrompt,
	)

	if err != nil {
		log.Printf("Failed to get draft data: %v", err)
		http.Error(w, "Draft not found", http.StatusNotFound)
		return
	}

	_ = json.Unmarshal(targetProductsJSON, &draft.TargetProducts)

	log.Printf("Target products: %v", draft.TargetProducts)

	// =========================
	// DETERMINE STATUS
	// =========================
	status := "published"
	if req.ScheduledFor != "" {

		ok := helper.ScheduleDraft(id, req.ScheduledFor)

		if !ok {
			http.Error(w, "Failed to schedule draft", http.StatusInternalServerError)
			return
		}

		// Schedule in Redis
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

		ScheduledFor, err := parseTimeFlexible(req.ScheduledFor)
		if err != nil {
			panic(err)
		}

		err = redisScheduler.ScheduleDraftTask(id, ScheduledFor, taskData)
		if err != nil {
			log.Printf("Failed to schedule in Redis: %v", err)
			// Rollback the insert
			database.GetDB().Exec("DELETE FROM drafts WHERE id = $1", id)
			http.Error(w, "Failed to schedule draft", http.StatusInternalServerError)
			return
		}

		// Response
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":       "Draft scheduled successfully",
			"draft_id":      id,
			"status":        "scheduled",
			"scheduled_for": ScheduledFor.Format(time.RFC3339),
		})

		return
	}

	// =========================
	// PROCESS PRODUCTS
	// =========================
	postService := helper.NewPostService()

	result, someFailed, allFailed, err := postService.ProcessDraftProducts(draft)
	if err != nil {
		log.Printf("ProcessDraftProducts error: %v", err)
		http.Error(w, "Failed processing products", http.StatusInternalServerError)
		return
	}

	log.Printf("[PROCESS RESULT] %+v", result)

	// =========================
	// ALL FAILED HANDLING
	// =========================
	if allFailed {
		log.Printf("All products failed for draft %s", id)
		writeAllProductsFailed(w, result)
		return
	}

	// =========================
	// DELETE DRAFT AFTER SUCCESS
	// =========================
	_, err = database.GetDB().Exec(`
		DELETE FROM drafts WHERE id = $1
	`, id)

	if err != nil {
		log.Printf("Failed to delete draft: %v", err)
		http.Error(w, "Failed to delete draft", http.StatusInternalServerError)
		return
	}

	log.Printf("Draft %s successfully published and deleted", id)

	// =========================
	// INSERT HISTORY (IMPORTANT FIX)
	// =========================
	sqlHistory, argsHistory, err := buildHistoryInsert(draft, userID, teamID)
	if err != nil {
		log.Printf("System error build history: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = database.GetDB().Exec(sqlHistory, argsHistory...)
	if err != nil {
		log.Printf("Failed to insert history: %v", err)
	}

	// =========================
	// RESPONSE
	// =========================
	response := map[string]interface{}{
		"message": "Draft processed",
		"status":  status,
		"results": result,
	}

	if someFailed {
		w.WriteHeader(http.StatusMultiStatus)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	log.Printf("Draft %s processed. Results: %+v", id, result)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *DraftHandler) PublishContent(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title          string   `json:"Title"`
		Topic          string   `json:"Topic"`
		Article        string   `json:"Article"`
		ImageURL       *string  `json:"ImageUrl"`
		TargetProducts []string `json:"TargetProducts"`
		ImagePrompt    string   `json:"ImagePrompt"`
	}

	ctx := r.Context()
	teamID := auth.GetTeamID(ctx)
	userID := auth.GetUserID(ctx)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// =====================
	// VALIDATION
	// =====================
	if req.Title == "" || req.Article == "" || len(req.TargetProducts) == 0 {
		log.Printf("Validation failed title=%q article_len=%d products=%v",
			req.Title,
			len(req.Article),
			req.TargetProducts,
		)

		http.Error(w, "title, article, and target_products are required", http.StatusBadRequest)
		return
	}

	draftData := models.DraftDataPost{
		Title:          req.Title,
		Topic:          req.Topic,
		Article:        req.Article,
		ImageURL:       req.ImageURL,
		TargetProducts: req.TargetProducts,
		ImagePrompt:    req.ImagePrompt,
	}

	// =====================
	// PROCESS PRODUCTS
	// =====================
	postService := helper.NewPostService()

	result, someFailed, allFailed, err := postService.ProcessDraftProducts(draftData)
	if err != nil {
		log.Printf("System error while processing products: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("[PROCESS RESULT] %+v", result)

	// =====================
	// HISTORY INSERT (FIXED)
	// =====================
	sqlHistory, argsHistory, err := buildHistoryInsert(draftData, userID, teamID)
	if err != nil {
		log.Printf("Failed to build history: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = database.GetDB().Exec(sqlHistory, argsHistory...)
	if err != nil {
		log.Printf("Failed to insert history: %v", err)
	}

	// =====================
	// ALL FAILED HANDLING
	// =====================
	if allFailed {
		log.Printf("All products failed for title: %s", req.Title)

		writeAllProductsFailed(w, result)
		return
	}

	// =====================
	// RESPONSE
	// =====================
	response := map[string]interface{}{
		"message": "Post published successfully",
		"title":   req.Title,
		"results": result,
	}

	if someFailed {
		w.WriteHeader(http.StatusMultiStatus) // 207
	} else {
		w.WriteHeader(http.StatusOK) // 200
	}

	log.Printf("Direct publish completed for title: %s | Results: %+v", req.Title, result)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper function to parse time flexibly
func parseTimeFlexible(timeStr string) (time.Time, error) {
	// List of possible formats
	formats := []string{
		time.RFC3339,          // "2006-01-02T15:04:05Z07:00"
		"2006-01-02T15:04:05", // Without timezone
		"2006-01-02 15:04:05",
		"2006-01-02",
		"2006-01-02T15:04",
		"2006-01-02 15:04",
	}

	for _, format := range formats {
		if parsed, err := time.Parse(format, timeStr); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time: %s", timeStr)
}

// ScheduleDraft - Endpoint khusus untuk menjadwalkan draft
func (h *DraftHandler) ScheduleDraft(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := auth.GetUserID(ctx)
	teamID := auth.GetTeamID(ctx)

	var req models.ScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.ScheduledFor == "" {
		http.Error(w, "scheduledFor is required", http.StatusBadRequest)
		return
	}

	now := time.Now()

	// Parse time
	var scheduledFor time.Time
	var err error

	scheduledFor, err = time.Parse(time.RFC3339, req.ScheduledFor)
	if err != nil {
		scheduledFor, err = time.Parse("2006-01-02T15:04:05", req.ScheduledFor)
		if err != nil {
			http.Error(w, "Invalid scheduledFor format", http.StatusBadRequest)
			return
		}
	}

	// Validate scheduled time is in the future
	if scheduledFor.Before(now) {
		http.Error(w, "scheduledFor must be in the future", http.StatusBadRequest)
		return
	}

	targetProductsJSON, _ := json.Marshal(req.TargetProducts)

	// Insert draft with scheduled status
	var draftID string
	err = database.GetDB().QueryRow(`
		INSERT INTO drafts (
			title, topic, article, image_url, image_prompt,
			target_products, has_image, status, scheduled_for,
			created_at, created_by, team_id, user_id
		) VALUES ($1,$2,$3,$4,$5,$6::jsonb,$7,$8,$9,$10,$11,$12,$13)
		RETURNING id
	`,
		req.Title, req.Topic, req.Article, req.ImageURL, req.ImagePrompt,
		string(targetProductsJSON), req.HasImage, "scheduled",
		scheduledFor, now, userID, teamID, userID,
	).Scan(&draftID)

	if err != nil {
		log.Printf("Failed to insert scheduled draft: %v", err)
		http.Error(w, "Failed to schedule draft", http.StatusInternalServerError)
		return
	}

	// Schedule in Redis
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

	err = redisScheduler.ScheduleDraftTask(draftID, scheduledFor, taskData)
	if err != nil {
		log.Printf("Failed to schedule in Redis: %v", err)
		// Rollback the insert
		database.GetDB().Exec("DELETE FROM drafts WHERE id = $1", draftID)
		http.Error(w, "Failed to schedule draft", http.StatusInternalServerError)
		return
	}

	// Response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "Draft scheduled successfully",
		"draft_id":      draftID,
		"status":        "scheduled",
		"scheduled_for": scheduledFor.Format(time.RFC3339),
	})
}

// CancelScheduledDraft cancels a scheduled draft
func (h *DraftHandler) CancelScheduledDraft(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DraftID int64 `json:"draft_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Cancel in Redis
	err := redisScheduler.CancelScheduledTask(req.DraftID)
	if err != nil {
		log.Printf("Failed to cancel scheduled task: %v", err)
	}

	// Update draft status
	_, err = database.GetDB().Exec(`
		UPDATE drafts 
		SET status = 'draft', updated_at = $1
		WHERE id = $2 AND status = 'scheduled'
	`, time.Now(), req.DraftID)

	if err != nil {
		log.Printf("Failed to update draft status: %v", err)
		http.Error(w, "Failed to cancel schedule", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Schedule cancelled",
		"draft_id": req.DraftID,
		"status":   "draft",
	})
}

// GetScheduledDrafts gets all scheduled drafts
func (h *DraftHandler) GetScheduledDrafts(w http.ResponseWriter, r *http.Request) {
	tasks, err := redisScheduler.GetScheduledTasks()
	if err != nil {
		log.Printf("Failed to get scheduled tasks: %v", err)
		http.Error(w, "Failed to get schedules", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func writeAllProductsFailed(
	w http.ResponseWriter,
	result interface{},
) {

	response := map[string]interface{}{
		"message": "Failed to publish draft: all products failed",
		"status":  "failed",
		"results": result,
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusBadGateway)

	json.NewEncoder(w).Encode(response)
}

func buildUpdateByIDQuery(
	table string,
	id interface{},
	data map[string]interface{},
) (string, []interface{}, error) {
	qb := builder.NewQueryBuilder(table)
	return qb.
		Update().
		SetMap(data).
		WhereEq("id", id).
		Build()
}

func buildHistoryInsert(req models.DraftDataPost, userID string, teamID string) (
	string,
	[]interface{},
	error,
) {

	log.Printf(
		"[HISTORY INSERT] req=%+v userID=%s teamID=%s",
		req,
		userID,
		teamID,
	)

	createdAt := time.Now()
	publishedAt := time.Now()

	qb := builder.NewQueryBuilder("histories")

	targetProductsJSON, _ := json.Marshal(req.TargetProducts)

	sqlQuery, args, err := qb.
		Insert().
		Columns(
			"title",
			"topic",
			"content",
			"image_url",
			"target_products",
			"status",
			"action",
			"error_message",
			"published_at",
			"scheduled_for",
			"created_by",
			"team_id",
			"created_at",
		).
		Values(
			req.Title,
			req.Topic,
			req.Article,
			req.ImageURL,
			targetProductsJSON,
			"published",
			"publish",
			nil,
			publishedAt,
			nil,
			userID,
			teamID,
			createdAt,
		).
		Build()

	if err != nil {
		log.Printf(
			"[HISTORY INSERT BUILD ERROR] %v",
			err,
		)
		return "", nil, err
	}

	log.Printf("[HISTORY SQL] %s", sqlQuery)
	log.Printf("[HISTORY ARGS] %#v", args)

	return sqlQuery, args, nil
}
