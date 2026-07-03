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
	"seo-backend/internal/domain/history"
	"seo-backend/internal/domain/paginate"
	"seo-backend/internal/domain/product"
	"seo-backend/internal/helper"
	"seo-backend/internal/scheduler"

	// "seo-backend/internal/helper"
	"seo-backend/internal/models"

	"golang.org/x/net/html"
)

type DraftServiceImpl struct {
	repo              draft.Repository
	redisScheduler    *scheduler.RedisScheduler
	repoHistory       history.HistoryRepository
	productController product.ProductService
	productService    product.ProductRepository
	postService       helper.PostService
}

func NewService(
	repo draft.Repository,
	repoHistory history.HistoryRepository,
	redisScheduler *scheduler.RedisScheduler,
	productController product.ProductService,
	productService product.ProductRepository,
	postService helper.PostService,
) draft.Service {

	return &DraftServiceImpl{
		repo:              repo,
		repoHistory:       repoHistory,
		redisScheduler:    redisScheduler,
		productController: productController,
		productService:    productService,
		postService:       postService,
	}
}

func (s *DraftServiceImpl) GetAll(ctx context.Context, userFilter models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[draft.Draft], error) {
	return s.repo.GetAll(ctx, userFilter, params)
}

func (s *DraftServiceImpl) GetDashboardStats(ctx context.Context, filter models.UserContext) (*draft.DraftStats, error) {
	return s.repo.GetDashboardStats(ctx, filter)
}

func (s *DraftServiceImpl) GetAllScheduled(ctx context.Context, usrCtx models.UserContext, params paginate.PaginationParams) (*paginate.PaginatedResult[draft.Draft], error) {

	return s.repo.GetAllScheduled(ctx, usrCtx, params)
}

// CreateDraft implements draft.Service
func (s *DraftServiceImpl) CreateDraft(ctx context.Context, req draft.CreateDraftRequest, userID, teamID string) (string, error) {
	req.SEOScore = draft.CalculateSEOScore(req.Title, req.Article, req.Topic, req.Topic, req.Keywords).Total
	return s.repo.Create(ctx, req, userID, teamID)
}

// UpdateDraft implements draft.Service
func (s *DraftServiceImpl) UpdateDraft(ctx context.Context, id string, TeamID string, updates map[string]interface{}) error {
	data := prepareUpdateData(updates)

	if len(data) == 0 {
		return fmt.Errorf("no fields to update")
	}
	return s.repo.Update(ctx, id, TeamID, data)
}

// DeleteDraft implements draft.Service
func (s *DraftServiceImpl) DeleteDraft(ctx context.Context, TeamID string, id string) error {
	return s.repo.Delete(ctx, TeamID, id)
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
	userCtx models.UserContext,
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

	seoScore := draftData.SEOScore

	excerpt := draftData.Excerpt
	if len(req.TargetProducts) > 0 {
		excerpt = req.Excerpt
	}

	// TAMBAHAN: Generate meta tags dan inject ke article
	article = s.injectMetaTagsToArticle(article, title, topic, excerpt, *imageURL)

	draftData.Title = title
	draftData.Topic = topic
	draftData.Article = article
	draftData.ImageURL = imageURL
	draftData.TargetProducts = targetProducts
	draftData.SEOScore = seoScore
	draftData.Excerpt = excerpt

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
		SEOScore:       draftData.SEOScore,
		Excerpt:        excerpt,
		Slug:           draftData.Slug,
	}

	// Jika ada schedule
	if req.ScheduledFor != "" {
		log.Printf("[PublishDraft] SCHEDULE MODE id=%s scheduledFor=%s", id, req.ScheduledFor)

		result, err := s.scheduleDraft(ctx, id, req.ScheduledFor, draftData, userCtx)
		if err != nil {
			log.Printf("[PublishDraft] ERROR scheduleDraft id=%s err=%v", id, err)
			return nil, err
		}

		log.Printf("[PublishDraft] SUCCESS scheduleDraft id=%s", id)
		return result, nil
	}

	log.Printf("[PublishDraft] DIRECT PUBLISH MODE id=%s", id)

	result, err := s.processPublish(ctx, draftData, id, userCtx)

	if err != nil {
		log.Printf("[PublishDraft] ERROR processPublish id=%s err=%v", id, err)

		if histErr := s.repoHistory.Create(ctx, historyReq, userCtx.GetUserID(), userCtx.GetTeamID(), "failed"); histErr != nil {
			log.Printf("[PublishDraft] ERROR InsertHistory(failed) id=%s err=%v", id, histErr)
		} else {
			log.Printf("[PublishDraft] SUCCESS InsertHistory(failed) id=%s", id)
		}

		return nil, err
	}

	log.Printf("[PublishDraft] SUCCESS processPublish id=%s result=%+v", id, result)

	if histErr := s.repoHistory.Create(ctx, historyReq, userCtx.GetUserID(), userCtx.GetTeamID(), "published"); histErr != nil {
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
	userCtx models.UserContext,
) (*draft.PublishResult, error) {

	log.Println("========== START PublishContent ==========")

	log.Printf("REQUEST TEAM_ID=%s USER_ID=%s", userCtx.GetTeamID(), userCtx.GetUserID())

	log.Printf(
		"REQUEST DATA => Title=%s Topic=%s ImageURL=%v TargetProducts=%v keywords=%v",
		req.Title,
		req.Topic,
		req.ImageURL,
		req.TargetProducts,
		req.Keywords,
	)

	// Validasi request
	if err := validatePublishRequest(req); err != nil {
		log.Printf("VALIDATION ERROR => %v", err)
		log.Println("========== END PublishContent ==========")
		return nil, err
	}

	log.Println("VALIDATION SUCCESS")

	// TAMBAHAN: Generate meta tags dan inject ke article
	req.Article = s.injectMetaTagsToArticle(
		req.Article,
		req.Title,
		req.Topic,
		req.Excerpt,
		*req.ImageURL,
	)

	// Prepare history request
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

	// Process products
	result, someFailed, allFailed, err := s.postService.ProcessDraftProducts(ctx, req, userCtx)

	log.Printf("PROCESS RESULT => %+v", result)
	log.Printf("PROCESS FLAGS => someFailed=%v allFailed=%v", someFailed, allFailed)

	// Determine status based on result
	status := helper.DeterminePublishStatus(someFailed, allFailed, err)

	// Insert history (always attempt, even on error)
	historyErr := s.insertPublishHistory(ctx, historyReq, userCtx, status)
	if historyErr != nil {
		log.Printf("WARNING: Failed to insert history: %v", historyErr)
		// Don't return error, just log it
	}

	// Handle processing error
	if err != nil {
		log.Printf("PROCESS ERROR => %v", err)
		log.Println("========== END PublishContent ==========")
		return nil, fmt.Errorf("failed to process products: %w", err)
	}

	log.Println("PROCESS SUCCESS")

	// Build final result
	finalResult := &draft.PublishResult{
		Results:    result,
		SomeFailed: someFailed,
		AllFailed:  allFailed,
		Status:     status,
		Message:    helper.GetStatusMessage(someFailed, allFailed),
	}

	log.Printf("FINAL RESPONSE => %+v", finalResult)
	log.Println("========== END PublishContent ==========")

	return finalResult, nil
}

// injectMetaTagsToArticle menambahkan meta tag ke dalam article HTML
func (s *DraftServiceImpl) injectMetaTagsToArticle(
	article string,
	title string,
	topic string,
	excerpt string,
	imageURL string,
) string {
	// Jika article kosong, return apa adanya
	if article == "" {
		return article
	}

	// Generate meta tags
	metaTags := helper.GenerateMetaTags(title, topic, excerpt, imageURL)

	// Inject meta tags ke dalam <head> atau di awal article
	// Cek apakah ada tag <head>
	if strings.Contains(article, "<head>") {
		// Inject setelah <head>
		article = strings.Replace(article, "<head>", "<head>\n"+metaTags, 1)
	} else if strings.Contains(article, "<html") {
		// Inject setelah <html>
		article = strings.Replace(article, "<html", "<html>\n<head>\n"+metaTags+"\n</head>", 1)
	} else {
		// Jika tidak ada struktur HTML, tambahkan di awal
		article = "<!DOCTYPE html>\n<html>\n<head>\n" + metaTags + "\n</head>\n<body>\n" + article + "\n</body>\n</html>"
	}

	return article
}

// Helper function to insert history with retry
func (s *DraftServiceImpl) insertPublishHistory(
	ctx context.Context,
	historyReq draft.PublishHistoryRequest,
	userCtx models.UserContext,
	status string,
) error {
	const maxRetries = 3

	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			time.Sleep(time.Duration(i) * time.Second)
		}

		err := s.repoHistory.Create(
			ctx,
			historyReq,
			userCtx.GetUserID(),
			userCtx.GetTeamID(),
			status,
		)

		if err == nil {
			log.Printf("INSERT HISTORY SUCCESS (status=%s)", status)
			return nil
		}

		lastErr = err
		log.Printf("INSERT HISTORY ATTEMPT %d/%d FAILED: %v", i+1, maxRetries, err)
	}

	return fmt.Errorf("failed to insert history after %d attempts: %w", maxRetries, lastErr)
}

// ScheduleDraft implements draft.Service
func (s *DraftServiceImpl) ScheduleDraft(ctx context.Context, req draft.ScheduleRequest, userCtx models.UserContext) (string, error) {
	scheduledFor := helper.ParseWIBTime(req.ScheduledFor)

	loc, _ := time.LoadLocation("Asia/Jakarta")

	if scheduledFor.Before(time.Now().In(loc)) {
		return "", fmt.Errorf("scheduled time must be in the future")
	}

	draftID, err := s.repo.InsertScheduledDraft(ctx, req, scheduledFor, userCtx.GetTeamID(), userCtx.GetUserID())
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
		TeamID:         userCtx.GetTeamID(),
		UserID:         userCtx.GetUserID(),
	}

	if err := s.redisScheduler.ScheduleDraftTask(draftID, scheduledFor, taskData, userCtx); err != nil {
		s.repo.Delete(ctx, userCtx.GetTeamID(), draftID)
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
func (s *DraftServiceImpl) scheduleDraft(ctx context.Context, id, scheduledForStr string, draftData *draft.DraftData, userCtx models.UserContext) (*draft.PublishResult, error) {
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
		TeamID:         userCtx.GetTeamID(),
		UserID:         userCtx.GetUserID(),
	}

	if err := s.redisScheduler.ScheduleDraftTask(id, scheduledFor, taskData, userCtx); err != nil {
		log.Printf("ERROR on redis: %v", err)
		return nil, err
	}

	return &draft.PublishResult{
		Status:       "scheduled",
		ScheduledFor: &scheduledFor,
	}, nil
}

func (s *DraftServiceImpl) processPublish(ctx context.Context, draftData *draft.DraftData, id string, userCtx models.UserContext) (*draft.PublishResult, error) {
	draftPost := draft.DraftDataPost{
		Id:             id,
		Title:          draftData.Title,
		Topic:          draftData.Topic,
		Article:        draftData.Article,
		ImageURL:       draftData.ImageURL,
		ImagePrompt:    draftData.ImagePrompt,
		TargetProducts: draftData.TargetProducts,
		Excerpt:        draftData.Excerpt,
	}

	result, someFailed, allFailed, err := s.postService.ProcessDraftProducts(ctx, draftPost, userCtx)
	log.Printf("IS ERROR,%v", err)

	if err := s.repo.Delete(ctx, userCtx.GetTeamID(), id); err != nil {
		log.Printf("Failed to delete draft: %v", err)
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

	score := draft.CalculateSEOScore(draftData.Title, draftData.Article, draftData.Topic, draftData.Topic, draftData.Keywords)
	return &score, nil
}

// CheckSimilarity implements draft.Service
func (s *DraftServiceImpl) CheckSimilarity(ctx context.Context, id string, useRole models.UserContext, params paginate.PaginationParams) ([]draft.SimilarityResult, error) {
	target, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("draft not found: %w", err)
	}

	allDrafts, err := s.repo.GetAll(ctx, useRole, params)
	if err != nil || allDrafts == nil {
		return []draft.SimilarityResult{}, nil
	}

	var others []draft.Draft
	for _, d := range allDrafts.Data {
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

	fieldMap := map[string]string{
		"id":              "id",
		"title":           "title",
		"topic":           "topic",
		"article":         "article",
		"image_url":       "image_url",
		"image_prompt":    "image_prompt",
		"slug":            "slug",
		"target_products": "target_products",
		"status":          "status",
		"scheduled_for":   "scheduled_for",
		"has_image":       "has_image",
		"excerpt":         "excerpt",
		"team_id":         "team_id",
		"user_id":         "user_id",
		"created_by":      "created_by",
		"created_at":      "created_at",
		"updated_at":      "updated_at",
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
			log.Printf("[WARNING] Field '%s' not found in fieldMap, skipping", key)
			continue
		}

		log.Printf("[INFO] Mapped DB Field: %s", dbField)

		if key == "target_products" { // lowercase!
			log.Println("[INFO] Marshaling target_products to JSONB")

			jsonValue, err := json.Marshal(value)
			if err != nil {
				log.Printf("[ERROR] Failed to marshal target_products: %v", err)
				continue
			}

			log.Printf("[INFO] JSON Result: %s", string(jsonValue))

			data[dbField] = jsonValue // atau string(jsonValue)

		} else {
			log.Printf("[INFO] Assigning value to '%s'", dbField)
			log.Printf("[INFO] Value Type: %T", value)
			log.Printf("[INFO] Value: %#v", value)
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

// Helper function to mask sensitive values
func maskSensitiveValue(value string) string {
	if len(value) <= 4 {
		return "****"
	}
	if len(value) <= 8 {
		return value[:2] + "****" + value[len(value)-2:]
	}
	return value[:3] + "..." + value[len(value)-3:]
}

// Helper functions
func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func countSuccessfulNodes(nodeResults map[string]interface{}) int {
	count := 0
	for _, result := range nodeResults {
		if resultMap, ok := result.(map[string]interface{}); ok {
			if success, ok := resultMap["success"].(bool); ok && success {
				count++
			}
		}
	}
	return count
}

// ============================================================
// 🔑 HELPER: Enrich draft dengan previous results dari config
// ============================================================
func (s *DraftServiceImpl) enrichDraftWithPreviousResults(
	draftData draft.DraftDataPost,
	cfg *product.ProductConfig,
) (draft.DraftDataPost, map[string]interface{}) {

	if cfg == nil {
		return draftData, nil
	}

	results := cfg.GetAllExecutionResults()
	if len(results) == 0 {
		return draftData, nil
	}

	return draftData, results
}

func structToMap(data interface{}) map[string]interface{} {
	jsonBytes, _ := json.Marshal(data)
	var result map[string]interface{}
	json.Unmarshal(jsonBytes, &result)
	return result
}
