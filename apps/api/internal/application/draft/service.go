// internal/application/draft/service.go
package draft

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"seo-backend/internal/domain/draft"
	"seo-backend/internal/domain/product"
	"seo-backend/internal/helper"
	"seo-backend/internal/scheduler"

	"golang.org/x/net/html"
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

	// Fallback ke draftData jika req kosong
	title := draftData.Title
	if req.Title != "" {
		title = req.Title
	}

	topic := draftData.Topic
	if req.Topic != "" {
		topic = req.Topic
	}

	article := draftData.Article
	if req.Article != "" {
		article = req.Article
	}

	imageURL := draftData.ImageURL
	if req.ImageURL != nil {
		imageURL = req.ImageURL
	}

	targetProducts := draftData.TargetProducts
	if len(req.TargetProducts) > 0 {
		targetProducts = req.TargetProducts
	}

	log.Printf("KEYWORDS== %v", req.Keywords)

	keywords := draftData.Keywords
	if req.Keywords != nil {
		keywords = req.Keywords
	}

	draftData.Title = title
	draftData.Topic = topic
	draftData.Article = article
	draftData.ImageURL = imageURL
	draftData.TargetProducts = targetProducts
	draftData.Keywords = keywords

	log.Printf(
		"[PublishDraft] PAYLOAD title=%s topic=%s target_products=%v",
		title, topic, targetProducts,
	)

	historyReq := draft.PublishHistoryRequest{
		Title:          title,
		Topic:          topic,
		Article:        article,
		ImageURL:       imageURL,
		TargetProducts: targetProducts,
		Keywords:       draftData.Keywords,
	}

	// Jika ada schedule
	if req.ScheduledFor != "" {
		log.Printf("[PublishDraft] SCHEDULE MODE id=%s scheduledFor=%s", id, req.ScheduledFor)

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

	if err != nil {
		log.Printf("[PublishDraft] ERROR processPublish id=%s err=%v", id, err)

		if histErr := s.repo.InsertHistory(ctx, historyReq, userID, teamID, "failed"); histErr != nil {
			log.Printf("[PublishDraft] ERROR InsertHistory(failed) id=%s err=%v", id, histErr)
		} else {
			log.Printf("[PublishDraft] SUCCESS InsertHistory(failed) id=%s", id)
		}

		return nil, err
	}

	log.Printf("[PublishDraft] SUCCESS processPublish id=%s result=%+v", id, result)

	if histErr := s.repo.InsertHistory(ctx, historyReq, userID, teamID, "published"); histErr != nil {
		log.Printf("[PublishDraft] ERROR InsertHistory(published) id=%s err=%v", id, histErr)
	} else {
		log.Printf("[PublishDraft] SUCCESS InsertHistory(published) id=%s", id)
	}

	return result, nil
}

// PublishContent implements draft.Service
func (s *DraftServiceImpl) PublishContent(
	ctx context.Context,
	req draft.DraftDataPost,
	teamID,
	userID string,
) (*draft.PublishResult, error) {

	log.Println("========== START PublishContent ==========")

	log.Printf("REQUEST TEAM_ID=%s USER_ID=%s", teamID, userID)

	log.Printf(
		"REQUEST DATA => Title=%s Topic=%s ImageURL=%s TargetProducts=%v keywords=%v",
		req.Title,
		req.Topic,
		req.ImageURL,
		req.TargetProducts,
		req.Keywords,
	)

	if err := validatePublishRequest(req); err != nil {
		log.Printf("VALIDATION ERROR => %v", err)
		return nil, err
	}

	log.Println("VALIDATION SUCCESS")

	historyReq := draft.PublishHistoryRequest{
		Title:          req.Title,
		Topic:          req.Topic,
		Article:        req.Article,
		ImageURL:       req.ImageURL,
		TargetProducts: req.TargetProducts,
		Keywords:       req.Keywords,
	}

	log.Printf("HISTORY REQUEST => %+v", historyReq)

	log.Println("CALLING ProcessDraftProducts...")

	result, someFailed, allFailed, err := s.postService.ProcessDraftProducts(req)

	log.Printf("PROCESS RESULT => %+v", result)
	log.Printf("PROCESS FLAGS => someFailed=%v allFailed=%v", someFailed, allFailed)

	if err != nil {

		log.Printf("PROCESS ERROR => %v", err)

		log.Println("INSERT HISTORY STATUS=failed")

		historyErr := s.repo.InsertHistory(
			ctx,
			historyReq,
			userID,
			teamID,
			"failed",
		)

		if historyErr != nil {
			log.Printf("FAILED INSERT HISTORY => %v", historyErr)
		}

		return nil, fmt.Errorf("failed to process products: %w", err)
	}

	log.Println("PROCESS SUCCESS")

	log.Println("INSERT HISTORY STATUS=published")

	err = s.repo.InsertHistory(
		ctx,
		historyReq,
		userID,
		teamID,
		"published",
	)

	if err != nil {
		log.Printf("FAILED INSERT HISTORY => %v", err)
	} else {
		log.Println("INSERT HISTORY SUCCESS")
	}

	finalResult := &draft.PublishResult{
		Results:    result,
		SomeFailed: someFailed,
		AllFailed:  allFailed,
		Status:     "published",
	}

	log.Printf("FINAL RESPONSE => %+v", finalResult)

	log.Println("========== END PublishContent ==========")

	return finalResult, nil
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
		Id:             id,
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

// GetSEOScore implements draft.Service
func (s *DraftServiceImpl) GetSEOScore(ctx context.Context, id string) (*draft.SEOScore, error) {
	draftData, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("draft not found: %w", err)
	}

	score := draft.CalculateSEOScore(draftData.Title, draftData.Article, draftData.Topic, draftData.Topic)
	return &score, nil
}

// CheckSimilarity implements draft.Service
func (s *DraftServiceImpl) CheckSimilarity(ctx context.Context, id string, teamID string) ([]draft.SimilarityResult, error) {
	target, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("draft not found: %w", err)
	}

	if teamID == "" {
		return []draft.SimilarityResult{}, nil
	}

	allDrafts, err := s.repo.GetAll(ctx, teamID)
	if err != nil || allDrafts == nil || len(*allDrafts) == 0 {
		return []draft.SimilarityResult{}, nil
	}

	var others []draft.Draft
	for _, d := range *allDrafts {
		if d.ID != id {
			others = append(others, d)
		}
	}

	if len(others) == 0 {
		return []draft.SimilarityResult{}, nil
	}

	docs := []string{target.Title + " " + target.Article}
	for _, d := range others {
		docs = append(docs, d.Title+" "+d.Article)
	}

	vectors := draft.ComputeTFIDF(docs)
	targetVector := vectors[0]

	var results []draft.SimilarityResult
	for i, d := range others {
		sim := draft.CosineSimilarity(targetVector, vectors[i+1])
		if sim > 0.3 {
			results = append(results, draft.SimilarityResult{
				DraftID:    d.ID,
				Title:      d.Title,
				Similarity: math.Round(sim*100) / 100,
			})
		}
	}

	return results, nil
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

func CalculateSEOScore(title, content, excerpt, topic string) draft.SEOScore {
	details := map[string]int{}
	suggestions := []string{}
	total := 0

	keyword := strings.ToLower(topic)
	contentLower := strings.ToLower(content)
	titleLower := strings.ToLower(title)

	// 1. Keyword in title (20 pts)
	if strings.Contains(titleLower, keyword) {
		details["keyword_in_title"] = 20
		total += 20
	} else {
		suggestions = append(suggestions, "Add the main keyword to the title")
	}

	// 2. H1 exists (15 pts)
	if strings.Contains(content, "<h1") {
		details["has_h1"] = 15
		total += 15
	} else {
		suggestions = append(suggestions, "Add an H1 heading to the content")
	}

	// 3. H2 exists (10 pts)
	if strings.Contains(content, "<h2") {
		details["has_h2"] = 10
		total += 10
	} else {
		suggestions = append(suggestions, "Add H2 subheadings to the content")
	}

	// 4. Keyword in first 100 words (15 pts)
	words := strings.Fields(stripHTML(contentLower))
	first100 := strings.Join(words[:min(100, len(words))], " ")
	if strings.Contains(first100, keyword) {
		details["keyword_in_intro"] = 15
		total += 15
	} else {
		suggestions = append(suggestions, "Use the keyword in the first 100 words")
	}

	// 5. Meta description length 120-160 chars (15 pts)
	excerptLen := len(excerpt)
	if excerptLen >= 120 && excerptLen <= 160 {
		details["meta_description"] = 15
		total += 15
	} else {
		suggestions = append(suggestions, "Meta description should be 120-160 characters")
	}

	// 6. Content length > 600 words (15 pts)
	wordCount := len(strings.Fields(stripHTML(contentLower)))
	if wordCount >= 600 {
		details["content_length"] = 15
		total += 15
	} else {
		suggestions = append(suggestions, "Content should be at least 600 words")
	}

	// 7. Keyword in content (10 pts)
	if strings.Count(contentLower, keyword) >= 2 {
		details["keyword_density"] = 10
		total += 10
	} else {
		suggestions = append(suggestions, "Use the keyword at least 2 times in the content")
	}

	return draft.SEOScore{
		Total:       total,
		Details:     details,
		Suggestions: suggestions,
	}
}

func stripHTML(content string) string {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return content
	}
	var buf strings.Builder
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return buf.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
