// internal/application/draft/service.go
package draft

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/url"
	"regexp"
	"strings"
	"time"

	"seo-backend/internal/domain/draft"
	"seo-backend/internal/domain/history"
	"seo-backend/internal/domain/paginate"
	"seo-backend/internal/domain/product"
	"seo-backend/internal/helper"
	"seo-backend/internal/infrastructure/http/minio"
	"seo-backend/internal/scheduler"

	// "seo-backend/internal/helper"
	"seo-backend/internal/models"

	"github.com/google/uuid"
	"golang.org/x/net/html"
)

type DraftServiceImpl struct {
	repo              draft.Repository
	redisScheduler    *scheduler.RedisScheduler
	repoHistory       history.HistoryRepository
	productController product.ProductService
	productService    product.ProductRepository
	postService       helper.PostService
	minioClient       *minio.MinioService
}

func NewService(
	repo draft.Repository,
	repoHistory history.HistoryRepository,
	redisScheduler *scheduler.RedisScheduler,
	productController product.ProductService,
	productService product.ProductRepository,
	postService helper.PostService,
	minioClient *minio.MinioService,
) draft.Service {

	return &DraftServiceImpl{
		repo:              repo,
		repoHistory:       repoHistory,
		redisScheduler:    redisScheduler,
		productController: productController,
		productService:    productService,
		postService:       postService,
		minioClient:       minioClient,
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
	data := prepareUpdateData(req, s.minioClient)
	data.SEOScore = draft.CalculateSEOScore(req.Title, req.Article, req.Topic, req.Topic, req.Keywords).Total
	return s.repo.Create(ctx, data, userID, teamID)
}

// UpdateDraft implements draft.Service
func (s *DraftServiceImpl) UpdateDraft(ctx context.Context, id string, userID, TeamID string, updates draft.CreateDraftRequest) error {
	data := prepareUpdateData(updates, s.minioClient)
	data.TeamID = TeamID
	data.UserID = userID
	return s.repo.Update(ctx, id, TeamID, data)
}

// DeleteDraft implements draft.Service
func (s *DraftServiceImpl) DeleteDraft(ctx context.Context, TeamID string, id string) error {
	return s.repo.Delete(ctx, TeamID, id)
}

// GetDraftByID implements draft.Service
func (s *DraftServiceImpl) GetDraftByID(ctx context.Context, id string) (*draft.DraftData, error) {
	log.Printf("[INFO] GetDraftByID: %s", id)

	draftData, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("[ERROR] GetDraftByID failed: %v", err)
		return nil, err
	}

	// Refresh image_url
	if draftData.ImageURL != nil && *draftData.ImageURL != "" {
		log.Printf("[INFO] Refreshing image_url: %s", *draftData.ImageURL)
		newURL, err := s.minioClient.GetURL(ctx, *draftData.ImageURL, 7*24*time.Hour)
		if err != nil {
			log.Printf("[ERROR] Failed to refresh image_url: %v", err)
		} else {
			log.Printf("[SUCCESS] image_url refreshed: %s", newURL)
			*draftData.ImageURL = newURL
		}
	} else {
		log.Println("[INFO] No image_url to refresh")
	}

	// Refresh semua <img> di article
	if draftData.Article != "" {
		log.Println("[INFO] Refreshing images in article...")

		re := regexp.MustCompile(`<img[^>]+src="([^"]+)"`)
		count := 0

		draftData.Article = re.ReplaceAllStringFunc(draftData.Article, func(tag string) string {
			reSrc := regexp.MustCompile(`src="([^"]+)"`)
			match := reSrc.FindStringSubmatch(tag)
			if len(match) < 2 {
				log.Printf("[WARNING] Cannot extract src from tag: %s", tag[:min(50, len(tag))])
				return tag
			}

			oldURL := match[1]
			log.Printf("[INFO] Refreshing image #%d: %s", count+1, oldURL[:min(80, len(oldURL))])

			newURL, err := s.minioClient.GetURL(ctx, oldURL, 7*24*time.Hour)
			if err != nil {
				log.Printf("[ERROR] Failed to refresh image #%d: %v", count+1, err)
				return tag
			}

			count++
			log.Printf("[SUCCESS] Image #%d refreshed", count)
			return strings.Replace(tag, oldURL, newURL, 1)
		})

		log.Printf("[INFO] Refreshed %d images in article", count)
	} else {
		log.Println("[INFO] No article to refresh")
	}

	log.Printf("[SUCCESS] GetDraftByID completed: %s", id)
	return draftData, nil
}

// PublishDraft implements draft.Service
func (s *DraftServiceImpl) PublishDraft(
	ctx context.Context,
	id string,
	req draft.CreateDraftRequest,
	userCtx models.UserContext,
) (*draft.PublishResult, error) {

	log.Printf("========== START PublishDraft ==========")
	log.Printf("DRAFT_ID=%s TEAM_ID=%s USER_ID=%s", id, userCtx.GetTeamID(), userCtx.GetUserID())

	// Get draft data from repository
	log.Printf("Fetching draft data for ID: %s", id)
	draftData, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("ERROR GetByID: %v", err)
		return nil, err
	}
	log.Printf("SUCCESS GetByID for ID: %s", id)

	// Log draft data with safe dereference
	imageURLStr := "nil"
	if draftData.ImageURL != nil {
		imageURLStr = *draftData.ImageURL
	}

	log.Printf("DRAFT_DATA => ID=%s Title=%s Topic=%s ImageURL=%s TargetProducts=%v keywords=%v slug=%s excerpt=%s",
		id,
		draftData.Title,
		draftData.Topic,
		imageURLStr,
		draftData.TargetProducts,
		draftData.Keywords,
		draftData.Slug,
		draftData.Excerpt,
	)

	// Fallback ke draftData jika req kosong
	title := draftData.Title
	if req.Title != "" {
		title = req.Title
		log.Printf("Using request title: %s", title)
	}

	topic := draftData.Topic
	if req.Topic != "" {
		topic = req.Topic
		log.Printf("Using request topic: %s", topic)
	}

	article := draftData.Article
	if req.Article != "" {
		article = req.Article
		log.Printf("Using request article (length: %d)", len(article))
	}

	// FIX: Safe handling imageURL
	imageURL := draftData.ImageURL
	if req.ImageURL != nil {
		imageURL = req.ImageURL
		log.Printf("Using request imageURL: %s", *imageURL)
	}

	// Safe handling imageURL for logging
	var finalImageURLStr string
	if imageURL != nil {
		finalImageURLStr = *imageURL
		log.Printf("IMAGE_URL from request: %s", finalImageURLStr)
	} else {
		finalImageURLStr = ""
		log.Printf("WARNING: ImageURL is nil, using empty string")
	}

	// Safe handling targetProducts
	targetProducts := draftData.TargetProducts
	if len(req.TargetProducts) > 0 {
		targetProducts = req.TargetProducts
		log.Printf("Using request targetProducts: %v", targetProducts)
	}
	if len(targetProducts) == 0 {
		log.Printf("WARNING: TargetProducts is empty")
	}

	// Safe handling excerpt
	excerpt := draftData.Excerpt
	if req.Excerpt != "" {
		excerpt = req.Excerpt
		log.Printf("Using request excerpt: %s", excerpt)
	}
	if excerpt == "" {
		log.Printf("WARNING: Excerpt is empty")
	}

	// Safe handling slug
	slug := draftData.Slug
	if slug == "" {
		log.Printf("WARNING: Slug is empty")
	}

	// Safe handling keywords
	keywords := draftData.Keywords
	if len(keywords) == 0 {
		log.Printf("WARNING: Keywords is empty")
	}

	// Cek jika article juga bisa kosong
	if article == "" {
		log.Printf("WARNING: Article is empty")
	}

	// Inject meta tags dengan safe imageURL
	log.Printf("Injecting meta tags to article...")
	article = s.injectMetaTagsToArticle(
		article,
		title,
		topic,
		excerpt,
		finalImageURLStr, // Gunakan string, bukan pointer
	)
	log.Printf("Article injected successfully (length: %d)", len(article))

	// Update draftData with all values
	draftData.Title = title
	draftData.Topic = topic
	draftData.Article = article
	draftData.ImageURL = imageURL
	draftData.TargetProducts = targetProducts
	draftData.Excerpt = excerpt
	// Keep existing SEOScore and Slug

	log.Printf("DRAFT_DATA updated: slug=%s keywords_count=%d products_count=%d",
		slug,
		len(keywords),
		len(targetProducts),
	)

	// Prepare history request
	historyReq := draft.PublishHistoryRequest{
		Title:          title,
		Topic:          topic,
		Article:        article,
		ImageURL:       imageURL,
		TargetProducts: targetProducts,
		Keywords:       keywords,
		SEOScore:       draftData.SEOScore,
		Excerpt:        excerpt,
		Slug:           slug,
	}
	log.Printf("HISTORY_REQ prepared for ID: %s", id)

	// Jika ada schedule
	if req.ScheduledFor != "" {
		log.Printf("SCHEDULE_MODE: scheduledFor=%s", req.ScheduledFor)

		log.Printf("Calling scheduleDraft...")
		result, err := s.scheduleDraft(ctx, id, req.ScheduledFor, draftData, userCtx)
		if err != nil {
			log.Printf("ERROR scheduleDraft: %v", err)

			// Log history failed
			if histErr := s.repoHistory.Create(ctx, historyReq, userCtx.GetUserID(), userCtx.GetTeamID(), "failed"); histErr != nil {
				log.Printf("ERROR InsertHistory(failed): %v", histErr)
			}

			return nil, err
		}

		// Log history scheduled
		if histErr := s.repoHistory.Create(ctx, historyReq, userCtx.GetUserID(), userCtx.GetTeamID(), "scheduled"); histErr != nil {
			log.Printf("ERROR InsertHistory(scheduled): %v", histErr)
		} else {
			log.Printf("SUCCESS InsertHistory(scheduled)")
		}

		log.Printf("SUCCESS scheduleDraft for ID: %s", id)
		log.Printf("========== END PublishDraft ==========")
		return result, nil
	}

	// DIRECT PUBLISH MODE
	log.Printf("DIRECT_PUBLISH_MODE for ID: %s", id)

	// Process publish
	log.Printf("Calling processPublish...")
	result, err := s.processPublish(ctx, draftData, id, userCtx)
	if result.AllFailed {
		log.Printf("ERROR processPublish: %v", err)

		// Log history failed
		historyReq := draft.PublishHistoryRequest{
			Title:          req.Title,
			Topic:          req.Topic,
			Article:        article,
			ImageURL:       req.ImageURL,
			TargetProducts: req.TargetProducts,
			Keywords:       keywords,
			SEOScore:       0,
			Excerpt:        excerpt,
			Slug:           slug,
		}

		if histErr := s.repoHistory.Create(ctx, historyReq, userCtx.GetUserID(), userCtx.GetTeamID(), "failed"); histErr != nil {
			log.Printf("ERROR InsertHistory(failed): %v", histErr)
		}

		return result, err
	}

	log.Printf("SUCCESS processPublish")

	// Log history success
	historyReq = draft.PublishHistoryRequest{
		Title:          req.Title,
		Topic:          req.Topic,
		Article:        article,
		ImageURL:       req.ImageURL,
		TargetProducts: req.TargetProducts,
		Keywords:       keywords,
		SEOScore:       draftData.SEOScore,
		Excerpt:        excerpt,
		Slug:           slug,
	}

	if histErr := s.repoHistory.Create(ctx, historyReq, userCtx.GetUserID(), userCtx.GetTeamID(), "published"); histErr != nil {
		log.Printf("ERROR InsertHistory(published): %v", histErr)
	} else {
		log.Printf("SUCCESS InsertHistory(published)")
	}

	log.Printf("========== END PublishContent ==========")
	return result, nil
}

// Helper function untuk insert history
func (s *DraftServiceImpl) insertHistory(
	ctx context.Context,
	historyReq draft.PublishHistoryRequest,
	userCtx models.UserContext,
	status string,
) {
	if err := s.repoHistory.Create(ctx, historyReq, userCtx.GetUserID(), userCtx.GetTeamID(), status); err != nil {
		log.Printf("[PublishDraft] ERROR InsertHistory(%s) slug=%s err=%v", status, historyReq.Slug, err)
	} else {
		log.Printf("[PublishDraft] SUCCESS InsertHistory(%s) slug=%s", status, historyReq.Slug)
	}
}

// PublishContent implements draft.Service
func (s *DraftServiceImpl) PublishContent(
	ctx context.Context,
	req draft.CreateDraftRequest,
	userCtx models.UserContext,
) (*draft.PublishResult, error) {

	log.Printf("========== START PublishContent ==========")
	log.Printf("REQUEST TEAM_ID=%s USER_ID=%s", userCtx.GetTeamID(), userCtx.GetUserID())

	// Log request data dengan safe dereference
	imageURLStr := "nil"
	if req.ImageURL != nil {
		imageURLStr = *req.ImageURL
	}

	log.Printf("REQUEST DATA => Title=%s Topic=%s ImageURL=%s TargetProducts=%v keywords=%v slug=%s excerpt=%s",
		req.Title,
		req.Topic,
		imageURLStr,
		req.TargetProducts,
		req.Keywords,
		req.Slug,
		req.Excerpt,
	)

	// Validasi basic
	if req.Title == "" {
		log.Printf("VALIDATION ERROR: Title is empty")
		return nil, fmt.Errorf("title is required")
	}
	if req.Topic == "" {
		log.Printf("VALIDATION ERROR: Topic is empty")
		return nil, fmt.Errorf("topic is required")
	}

	log.Printf("VALIDATION SUCCESS")

	// FIX: Safe handling imageURL
	var articleImageURL string
	if req.ImageURL != nil {
		articleImageURL = *req.ImageURL
		log.Printf("IMAGE_URL from request: %s", articleImageURL)
	} else {
		articleImageURL = ""
		log.Printf("WARNING: ImageURL is nil, using empty string")
	}

	// Safe handling excerpt
	excerpt := req.Excerpt
	if excerpt == "" {
		log.Printf("WARNING: Excerpt is empty")
	}

	// Safe handling slug
	slug := req.Slug
	if slug == "" {
		log.Printf("WARNING: Slug is empty, generating from title")
		slug = s.generateSlug(req.Title)
	}

	// Safe handling keywords
	keywords := req.Keywords
	if len(keywords) == 0 {
		log.Printf("WARNING: Keywords is empty")
	}

	// BARIS 354: Kemungkinan error disini - injectMetaTagsToArticle dengan imageURL
	// FIX: Gunakan articleImageURL yang sudah di-safe
	log.Printf("Generating article with meta tags...")
	article := req.Article

	// Cek jika article juga bisa kosong
	if article == "" {
		log.Printf("WARNING: Article is empty")
	}

	// Inject meta tags dengan safe imageURL
	article = s.injectMetaTagsToArticle(
		article,
		req.Title,
		req.Topic,
		excerpt,
		articleImageURL, // Gunakan string, bukan pointer
	)

	log.Printf("Article generated successfully (length: %d)", len(article))

	// Create draft data untuk dipublish
	draftData := &draft.DraftData{
		Title:          req.Title,
		Topic:          req.Topic,
		Article:        article,
		ImageURL:       req.ImageURL,
		ImagePrompt:    req.ImagePrompt,
		TargetProducts: req.TargetProducts,
		Keywords:       keywords,
		Slug:           slug,
		Excerpt:        excerpt,
		SEOScore:       0, // Akan dihitung ulang
	}

	log.Printf("DRAFT_DATA created: slug=%s keywords_count=%d products_count=%d",
		slug,
		len(keywords),
		len(req.TargetProducts),
	)

	// Process publish
	log.Printf("Calling processPublish...")
	result, err := s.processPublish(ctx, draftData, "", userCtx)
	if result.AllFailed {
		log.Printf("ERROR processPublish: %v", err)

		// Log history failed
		historyReq := draft.PublishHistoryRequest{
			Title:          req.Title,
			Topic:          req.Topic,
			Article:        article,
			ImageURL:       req.ImageURL,
			TargetProducts: req.TargetProducts,
			Keywords:       keywords,
			SEOScore:       0,
			Excerpt:        excerpt,
			Slug:           slug,
		}

		if histErr := s.repoHistory.Create(ctx, historyReq, userCtx.GetUserID(), userCtx.GetTeamID(), "failed"); histErr != nil {
			log.Printf("ERROR InsertHistory(failed): %v", histErr)
		}

		return result, err
	}

	log.Printf("SUCCESS processPublish")

	// Log history success
	historyReq := draft.PublishHistoryRequest{
		Title:          req.Title,
		Topic:          req.Topic,
		Article:        article,
		ImageURL:       req.ImageURL,
		TargetProducts: req.TargetProducts,
		Keywords:       keywords,
		SEOScore:       draftData.SEOScore,
		Excerpt:        excerpt,
		Slug:           slug,
	}

	if histErr := s.repoHistory.Create(ctx, historyReq, userCtx.GetUserID(), userCtx.GetTeamID(), "published"); histErr != nil {
		log.Printf("ERROR InsertHistory(published): %v", histErr)
	} else {
		log.Printf("SUCCESS InsertHistory(published)")
	}

	log.Printf("========== END PublishContent ==========")
	return result, nil
}

// Helper function untuk generate slug jika kosong
func (s *DraftServiceImpl) generateSlug(title string) string {
	// Simple slug generation
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove special characters
	reg := regexp.MustCompile("[^a-z0-9-]")
	slug = reg.ReplaceAllString(slug, "")
	return slug
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

	if err := s.redisScheduler.ScheduleDraftTask(ctx, draftID, scheduledFor, taskData, userCtx); err != nil {
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

	if err := s.redisScheduler.ScheduleDraftTask(ctx, id, scheduledFor, taskData, userCtx); err != nil {
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
	if !allFailed {
		if err := s.repo.Delete(ctx, userCtx.GetTeamID(), id); err != nil {
			log.Printf("Failed to delete draft: %v", err)
		}
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
// Tambahkan parameter minioService
func prepareUpdateData(updates draft.CreateDraftRequest, minioService *minio.MinioService) draft.CreateDraftRequest {
	log.Println("========== PREPARE UPDATE DATA ==========")

	// ========== PROCESS IMAGES ==========
	ctx := context.Background()

	// Create a copy of updates to modify
	processedUpdates := updates

	// Process image_url
	if processedUpdates.ImageURL != nil {
		if isBase64Image(processedUpdates.ImageURL) {
			log.Println("[INFO] Detected base64 in image_url, uploading...")
			uploadedURL, err := uploadBase64ToMinio(ctx, minioService, *processedUpdates.ImageURL, "image_url")
			if err != nil {
				log.Printf("[ERROR] Failed to upload base64 image: %v", err)
				emptyStr := ""
				processedUpdates.ImageURL = &emptyStr
			} else {
				log.Printf("[SUCCESS] Uploaded image_url: %s", uploadedURL)
				processedUpdates.ImageURL = &uploadedURL
			}
		} else if strings.Contains(*processedUpdates.ImageURL, minioService.Bucket) {
			objectName := extractObjectName(*processedUpdates.ImageURL)
			log.Printf("[INFO] Extracted object name from URL: %s", objectName)
			processedUpdates.ImageURL = &objectName
		} else {
			log.Printf("[INFO] Keeping original image_url: %s", *processedUpdates.ImageURL)
		}
	} else {
		log.Println("[INFO] No image_url provided")
	}

	// Process article images
	if processedUpdates.Article != "" {
		log.Println("[INFO] Processing article images...")
		processedArticle := processArticleImages(ctx, minioService, processedUpdates.Article)
		processedUpdates.Article = processedArticle
		log.Printf("[INFO] Article processed, length: %d", len(processedArticle))
	}

	log.Printf("[SUCCESS] prepareUpdateData completed")
	return processedUpdates
}

func processArticleImages(ctx context.Context, minioService *minio.MinioService, article string) string {
	re := regexp.MustCompile(`<img[^>]+src="([^"]+)"`)

	return re.ReplaceAllStringFunc(article, func(tag string) string {
		reSrc := regexp.MustCompile(`src="([^"]+)"`)
		match := reSrc.FindStringSubmatch(tag)
		if len(match) < 2 {
			return tag
		}

		oldURL := match[1]

		// Base64 → Upload
		if isBase64Image(&oldURL) {
			newURL, err := uploadBase64ToMinio(ctx, minioService, oldURL, "article_img")
			if err != nil {
				log.Printf("[ERROR] Failed to upload: %v", err)
				return tag
			}
			log.Printf("[SUCCESS] Base64 uploaded: %s", newURL)
			return strings.Replace(tag, oldURL, newURL, 1)
		}

		// URL Minio → Ambil object name
		if strings.Contains(oldURL, minioService.Bucket) {
			objectName := extractObjectName(oldURL)
			log.Printf("[INFO] Extracted object name: %s", objectName)
			return strings.Replace(tag, oldURL, objectName, 1)
		}

		// URL external → Biarkan
		return tag
	})
}

func extractObjectName(imageURL string) string {
	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return imageURL
	}

	// /kabar-storage/blog-images/2026/07/05/uuid.jpg?X-Amz-...
	// → blog-images/2026/07/05/uuid.jpg
	path := strings.TrimPrefix(parsedURL.Path, "/")
	parts := strings.SplitN(path, "/", 2) // ["kabar-storage", "blog-images/..."]

	if len(parts) == 2 {
		return parts[1]
	}

	return path
}

// Helper function: Proses semua base64 images di article
func processArticleBase64Images(ctx context.Context, minioService *minio.MinioService, article string) (string, error) {
	// Regex untuk mencari tag <img> dengan src base64
	imgTagRegex := regexp.MustCompile(`<img[^>]+src="(data:image\/[^"]+)"[^>]*>`)
	matches := imgTagRegex.FindAllStringSubmatch(article, -1)

	log.Printf("[INFO] Found %d base64 images in article", len(matches))

	processedArticle := article
	for i, match := range matches {
		fullTag := match[0]   // full <img ...> tag
		base64Src := match[1] // hanya data:image... base64

		log.Printf("[INFO] Processing article image %d/%d", i+1, len(matches))

		// Upload ke Minio
		uploadedURL, err := uploadBase64ToMinio(ctx, minioService, base64Src, fmt.Sprintf("article_img_%d", i))
		if err != nil {
			log.Printf("[ERROR] Failed to upload article image %d: %v", i, err)
			continue
		}

		// Replace src di tag img
		newTag := strings.Replace(fullTag, base64Src, uploadedURL, 1)
		processedArticle = strings.Replace(processedArticle, fullTag, newTag, 1)
		log.Printf("[SUCCESS] Replaced image %d with: %s", i+1, uploadedURL)
	}

	return processedArticle, nil
}

// Helper function: Cek apakah string adalah base64 image
func isBase64Image(data *string) bool {
	base64Pattern := regexp.MustCompile(`^data:image\/(png|jpeg|jpg|gif|webp|svg\+xml);base64,`)
	return base64Pattern.MatchString(*data)
}

// Helper function: Upload base64 image ke Minio
func uploadBase64ToMinio(ctx context.Context, minioService *minio.MinioService, base64Data string, source string) (string, error) {
	log.Printf("[INFO] Uploading base64 from %s to Minio", source)

	re := regexp.MustCompile(`^data:(image\/[^;]+);base64,(.*)`)
	matches := re.FindStringSubmatch(base64Data)
	if len(matches) != 3 {
		return "", fmt.Errorf("invalid base64 image format")
	}

	mimeType := matches[1]
	base64String := matches[2]

	imageData, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	ext := getExtensionFromMimeType(mimeType)
	objectName := fmt.Sprintf("blog-images/%s/%s%s",
		time.Now().Format("2006/01/02"),
		uuid.New().String(),
		ext,
	)
	log.Printf("[INFO] Object name: %s", objectName)

	reader := bytes.NewReader(imageData)
	uploadedResult, err := minioService.Upload(
		ctx,
		objectName,
		reader,
		int64(len(imageData)),
		mimeType,
	)
	if err != nil {
		log.Printf("[ERROR] Failed to upload to Minio: %v", err)
		return "", fmt.Errorf("failed to upload image to Minio: %w", err)
	}
	log.Printf("[DEBUG] Minio upload result: %s", uploadedResult)

	// 🔥 LANGSUNG RETURN PRESIGNED URL, JANGAN DI PARSE!
	return uploadedResult, nil
}

// Helper function: Dapatkan extension dari MIME type
func getExtensionFromMimeType(mimeType string) string {
	switch mimeType {
	case "image/png":
		return ".png"
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "image/svg+xml":
		return ".svg"
	default:
		return ".jpg"
	}
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
