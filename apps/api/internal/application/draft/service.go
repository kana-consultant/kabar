// internal/application/draft/service.go
package draft

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"seo-backend/internal/domain/draft"
	"seo-backend/internal/domain/product"
	"seo-backend/internal/helper"
	"seo-backend/internal/scheduler"
)

type DraftServiceImpl struct {
	repo           draft.Repository
	redisScheduler *scheduler.RedisScheduler
	postService    *helper.PostService
	productService product.ProductRepository
}

func NewService(
	repo draft.Repository,
	redisScheduler *scheduler.RedisScheduler,
	postService *helper.PostService,
	productService product.ProductRepository,
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
func (s *DraftServiceImpl) PublishDraft(
	ctx context.Context,
	id string,
	req draft.CreateDraftRequest,
	teamID, userID string,
) (*draft.PublishResult, error) {

	log.Printf("[PublishDraft] START id=%s", id)

	draftData, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("[PublishDraft] ERROR GetByID id=%s err=%v", id, err)
		return nil, err
	}

	log.Printf("[PublishDraft] SUCCESS GetByID id=%s", id)

	// fallback ke draftData jika req kosong
	title := draftData.Title
	if req.Title != "" {
		log.Printf("[PublishDraft] get by req -> title")
		title = req.Title
	}

	topic := draftData.Topic
	if req.Topic != "" {
		log.Printf("[PublishDraft] get by req -> topic")
		topic = req.Topic
	}

	article := draftData.Article
	if req.Article != "" {
		log.Printf("[PublishDraft] get by req -> article")
		article = req.Article
	}

	imageURL := draftData.ImageURL
	if req.ImageURL != nil {
		log.Printf("[PublishDraft] get by req -> image_url")
		imageURL = req.ImageURL
	}

	targetProducts := draftData.TargetProducts
	if len(req.TargetProducts) > 0 {
		log.Printf("[PublishDraft] get by req -> target_products")
		targetProducts = req.TargetProducts
	}

	// replace draftData juga
	draftData.Title = title
	draftData.Topic = topic
	draftData.Article = article
	draftData.ImageURL = imageURL
	draftData.TargetProducts = targetProducts

	log.Printf(
		"[PublishDraft] PAYLOAD title=%s topic=%s target_products=%v",
		title,
		topic,
		targetProducts,
	)

	historyReq := draft.PublishHistoryRequest{
		Title:          title,
		Topic:          topic,
		Article:        article,
		ImageURL:       imageURL,
		TargetProducts: targetProducts,
	}

	log.Printf("[PublishDraft] Draft fetched successfully id=%s", id)

	// jika ada schedule
	if req.ScheduledFor != "" {

		log.Printf(
			"[PublishDraft] SCHEDULE MODE id=%s scheduledFor=%s",
			id,
			req.ScheduledFor,
		)

		result, err := s.scheduleDraft(ctx, id, req.ScheduledFor, draftData, teamID, userID)
		if err != nil {
			log.Printf("[PublishDraft] ERROR scheduleDraft id=%s err=%v", id, err)
			return nil, err
		}

		log.Printf("[PublishDraft] SUCCESS scheduleDraft id=%s", id)

		return result, nil
	}

	log.Printf("[PublishDraft] DIRECT PUBLISH MODE id=%s", id)

	result, err := s.processPublish(ctx, draftData, id, teamID, userID)

	log.Printf("[PublishDraft] processPublish result=%+v", result)
	log.Printf("[PublishDraft] processPublish err=%v", err)

	if err != nil {
		log.Printf("[PublishDraft] INSERT HISTORY FAILED STATUS id=%s", id)

		s.repo.InsertHistory(ctx, historyReq, userID, teamID, "failed")

		log.Printf("[PublishDraft] ERROR processPublish id=%s err=%v", id, err)

		return nil, err
	}

	log.Printf("[PublishDraft] SUCCESS processPublish id=%s", id)

	return result, nil
}

// PublishContent implements draft.Service
func (s *DraftServiceImpl) PublishContent(ctx context.Context, req draft.DraftDataPost, teamID, userID string) (*draft.PublishResult, error) {
	if err := validatePublishRequest(req); err != nil {
		return nil, err
	}

	historyReq := draft.PublishHistoryRequest{
		Title:          req.Title,
		Topic:          req.Topic,
		Article:        req.Article,
		ImageURL:       req.ImageURL,
		TargetProducts: req.TargetProducts,
	}

	result, someFailed, allFailed, err := s.postService.ProcessDraftProducts(req)
	log.Printf("========== ERROR %v", err)
	if err != nil {
		s.repo.InsertHistory(ctx, historyReq, userID, teamID, "failed")
		return nil, fmt.Errorf("failed to process products: %w", err)
	} else {
		err = s.repo.InsertHistory(ctx, historyReq, userID, teamID, "published")

	}

	if err != nil {
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
	scheduledFor := helper.ParseWIBTime(req.ScheduledFor)

	loc, _ := time.LoadLocation("Asia/Jakarta")

	if scheduledFor.Before(time.Now().In(loc)) {
		return "", fmt.Errorf("scheduled time must be in the future")
	}

	draftID, err := s.repo.InsertScheduledDraft(ctx, req, scheduledFor, teamID, userID)
	if err != nil {
		log.Printf("ERROR InsertScheduledDraft : %v", err)
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
		log.Printf("ERROR ScheduleDraftTask : %v", err)
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
	scheduledFor := helper.ParseWIBTime(scheduledForStr)

	if err := s.repo.UpdateStatus(ctx, id, "scheduled", &scheduledFor); err != nil {
		log.Printf("ERROR : %v", err)
		return nil, err
	}

	var imageURL string

	if draftData.ImageURL != nil {
		imageURL = *draftData.ImageURL
	}

	taskData := &scheduler.ScheduledTask{
		DraftID:        id,
		Title:          draftData.Title,
		Topic:          draftData.Topic,
		Article:        draftData.Article,
		ImageURL:       imageURL,
		ImagePrompt:    draftData.ImagePrompt,
		TargetProducts: draftData.TargetProducts,
		TeamID:         teamID,
		UserID:         userID,
	}

	if err := s.redisScheduler.ScheduleDraftTask(id, scheduledFor, taskData); err != nil {
		log.Printf("ERROR on redis: %v", err)
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

	historyReq := draft.PublishHistoryRequest{
		Title:          draftPost.Title,
		Topic:          draftPost.Topic,
		Article:        draftPost.Article,
		ImageURL:       draftPost.ImageURL,
		TargetProducts: draftPost.TargetProducts,
	}

	result, someFailed, allFailed, err := s.postService.ProcessDraftProducts(draftPost)
	log.Printf("IS ERROR,%v", err)
	if err != nil {
		s.repo.InsertHistory(ctx, historyReq, userID, teamID, "failed")
		return nil, err
	}

	if allFailed {
		s.repo.InsertHistory(ctx, historyReq, userID, teamID, "failed")
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

	log.Println("========== PREPARE UPDATE DATA ==========")

	// mapping berdasarkan request yang masuk
	fieldMap := map[string]string{
		"id":             "id",
		"title":          "title",
		"topic":          "topic",
		"article":        "article",
		"ImageUrl":       "image_url",
		"ImagePrompt":    "image_prompt",
		"targetProducts": "target_products",
		"status":         "status",
		"scheduledFor":   "scheduled_for",
		"hasImage":       "has_image",
		"teamId":         "team_id",
		"userId":         "user_id",
		"createdBy":      "created_by",
		"createdAt":      "created_at",
		"updatedAt":      "updated_at",
	}

	log.Println("[INFO] Incoming Updates:")
	for k, v := range updates {
		log.Printf("  %s => %#v\n", k, v)
	}

	data := make(map[string]interface{})

	for key, value := range updates {

		log.Println("--------------------------------------------------")
		log.Println("[INFO] Processing Field:", key)

		dbField, ok := fieldMap[key]

		if !ok {
			log.Println("[WARNING] Field not found in fieldMap:", key)
			continue
		}

		log.Println("[INFO] Mapped DB Field:", dbField)

		// khusus target products
		if key == "TargetProducts" {

			log.Println("[INFO] Marshaling TargetProducts to JSON")

			jsonValue, err := json.Marshal(value)
			if err != nil {

				log.Println("[ERROR] Failed to marshal TargetProducts:", err)
				continue
			}

			log.Println("[INFO] JSON Result:", string(jsonValue))

			data[dbField] = jsonValue

		} else {

			log.Printf("[INFO] Assigning value to '%s'\n", dbField)
			log.Printf("[INFO] Value Type: %T\n", value)
			log.Printf("[INFO] Value: %#v\n", value)

			data[dbField] = value
		}
	}

	log.Println("==================================================")
	log.Println("[INFO] Final Prepared Data:")

	for k, v := range data {
		log.Printf("  %s => %#v\n", k, v)
	}

	log.Println("[SUCCESS] prepareUpdateData completed")

	return data
}

func validatePublishRequest(req draft.DraftDataPost) error {
	if req.Title == "" || req.Article == "" || len(req.TargetProducts) == 0 {
		return fmt.Errorf("title, article, and target_products are required")
	}
	return nil
}
