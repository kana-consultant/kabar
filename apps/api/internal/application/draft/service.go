// internal/application/draft/service.go
package draft

import (
	"bytes"
	"context"
	"encoding/base64"
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
	"seo-backend/internal/models"
	"seo-backend/internal/scheduler"

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

	// Inject meta tags saat create draft
	data.Article = s.injectMetaTagsToArticle(
		data.Article,
		data.Title,
		data.Topic,
		data.Excerpt,
		getImageURLString(data.ImageURL),
	)

	data.SEOScore = draft.CalculateSEOScore(req.Title, req.Article, req.Topic, req.Excerpt, req.Keywords).Total
	return s.repo.Create(ctx, data, userID, teamID)
}

// UpdateDraft implements draft.Service
func (s *DraftServiceImpl) UpdateDraft(ctx context.Context, id string, userID, TeamID string, updates draft.CreateDraftRequest) error {
	data := prepareUpdateData(updates, s.minioClient)

	// Inject meta tags saat update draft
	data.Article = s.injectMetaTagsToArticle(
		data.Article,
		data.Title,
		data.Topic,
		data.Excerpt,
		getImageURLString(data.ImageURL),
	)

	data.TeamID = TeamID
	data.UserID = userID
	return s.repo.Update(ctx, id, TeamID, data)
}

// DeleteDraft implements draft.Service
func (s *DraftServiceImpl) DeleteDraft(ctx context.Context, TeamID string, id string) error {
	// Get draft dulu untuk cek status
	draft, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get draft: %w", err)
	}

	// Jika draft scheduled, hapus schedule di Redis
	if draft.Status == "scheduled" {
		// Hapus scheduled task di Redis
		if s.redisScheduler != nil {
			err := s.redisScheduler.CancelScheduledTask(ctx, id)
			if err != nil {
				log.Printf("⚠️ Failed to cancel scheduled task for draft %s: %v", id, err)
				// Tetap lanjutkan delete meskipun gagal hapus schedule
			} else {
				log.Printf("✅ Cancelled scheduled task for draft %s", id)
			}
		}
	}

	// Hapus draft dari database
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
	}

	// Refresh semua <img> di article
	if draftData.Article != "" {
		log.Println("[INFO] Refreshing images in article...")
		draftData.Article = s.minioClient.RefreshArticleImages(ctx, draftData.Article)
	}

	log.Printf("[SUCCESS] GetDraftByID completed: %s", id)
	return draftData, nil
}

// refreshArticleImages refresh semua URL gambar di article
func (s *DraftServiceImpl) refreshArticleImages(ctx context.Context, article string) string {
	re := regexp.MustCompile(`<img[^>]+src="([^"]+)"`)
	count := 0

	return re.ReplaceAllStringFunc(article, func(tag string) string {
		reSrc := regexp.MustCompile(`src="([^"]+)"`)
		match := reSrc.FindStringSubmatch(tag)
		if len(match) < 2 {
			return tag
		}

		oldURL := match[1]
		newURL, err := s.minioClient.GetURL(ctx, oldURL, 7*24*time.Hour)
		if err != nil {
			log.Printf("[ERROR] Failed to refresh image #%d: %v", count+1, err)
			return tag
		}

		count++
		log.Printf("[SUCCESS] Image #%d refreshed", count)
		return strings.Replace(tag, oldURL, newURL, 1)
	})
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

	// Get draft data
	draftData, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("ERROR GetByID: %v", err)
		return nil, err
	}

	// Merge request data dengan draft data (fallback)
	title := coalesceString(req.Title, draftData.Title)
	topic := coalesceString(req.Topic, draftData.Topic)
	article := coalesceString(req.Article, draftData.Article)
	imageURL := coalescePointer(req.ImageURL, draftData.ImageURL)
	targetProducts := coalesceSlice(req.TargetProducts, draftData.TargetProducts)
	excerpt := coalesceString(req.Excerpt, draftData.Excerpt)
	slug := draftData.Slug
	keywords := draftData.Keywords

	// TIDAK inject meta tags saat publish, gunakan article apa adanya
	// karena meta tags sudah di-inject saat create/update draft

	// Update draftData dengan merged values
	draftData.Title = title
	draftData.Topic = topic
	draftData.Article = article
	draftData.ImageURL = imageURL
	draftData.TargetProducts = targetProducts
	draftData.Excerpt = excerpt

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

	// Schedule mode
	if req.ScheduledFor != "" {
		result, err := s.scheduleDraft(ctx, id, req.ScheduledFor, draftData, userCtx)
		if err != nil {
			s.insertHistory(ctx, historyReq, userCtx, "failed")
			return nil, err
		}
		return result, nil
	}

	// Direct publish mode
	result, err := s.processPublish(ctx, draftData, id, userCtx)
	if result.AllFailed {
		s.insertHistory(ctx, historyReq, userCtx, "failed")
		return result, err
	}

	s.insertHistory(ctx, historyReq, userCtx, "published")
	return result, nil
}

// PublishContent implements draft.Service
func (s *DraftServiceImpl) PublishContent(
	ctx context.Context,
	req draft.CreateDraftRequest,
	userCtx models.UserContext,
) (*draft.PublishResult, error) {
	log.Printf("========== START PublishContent ==========")

	// Validasi
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if req.Topic == "" {
		return nil, fmt.Errorf("topic is required")
	}

	slug := req.Slug
	if slug == "" {
		slug = s.generateSlug(req.Title)
	}

	// TIDAK inject meta tags saat publish content
	// Meta tags sudah ada di article yang disimpan

	draftData := &draft.DraftData{
		Title:          req.Title,
		Topic:          req.Topic,
		Article:        req.Article, // Gunakan article apa adanya
		ImageURL:       req.ImageURL,
		ImagePrompt:    req.ImagePrompt,
		TargetProducts: req.TargetProducts,
		Keywords:       req.Keywords,
		Slug:           slug,
		Excerpt:        req.Excerpt,
		SEOScore:       0,
	}

	// Process publish
	result, err := s.processPublish(ctx, draftData, "", userCtx)
	if result.AllFailed {
		historyReq := draft.PublishHistoryRequest{
			Title:          req.Title,
			Topic:          req.Topic,
			Article:        req.Article,
			ImageURL:       req.ImageURL,
			TargetProducts: req.TargetProducts,
			Keywords:       req.Keywords,
			SEOScore:       0,
			Excerpt:        req.Excerpt,
			Slug:           slug,
		}
		s.insertHistory(ctx, historyReq, userCtx, "failed")
		return result, err
	}

	historyReq := draft.PublishHistoryRequest{
		Title:          req.Title,
		Topic:          req.Topic,
		Article:        req.Article,
		ImageURL:       req.ImageURL,
		TargetProducts: req.TargetProducts,
		Keywords:       req.Keywords,
		SEOScore:       draftData.SEOScore,
		Excerpt:        req.Excerpt,
		Slug:           slug,
	}
	s.insertHistory(ctx, historyReq, userCtx, "published")

	return result, nil
}

// Helper functions untuk coalesce
func coalesceString(newVal, fallback string) string {
	if newVal != "" {
		return newVal
	}
	return fallback
}

func coalescePointer(newVal, fallback *string) *string {
	if newVal != nil {
		return newVal
	}
	return fallback
}

func coalesceSlice(newVal, fallback []string) []string {
	if len(newVal) > 0 {
		return newVal
	}
	return fallback
}

func getImageURLString(imageURL *string) string {
	if imageURL != nil {
		return *imageURL
	}
	return ""
}

// generateSlug generate slug dari title
func (s *DraftServiceImpl) generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	reg := regexp.MustCompile("[^a-z0-9-]")
	slug = reg.ReplaceAllString(slug, "")
	return slug
}

// injectMetaTagsToArticle menambahkan meta tag ke dalam article HTML
// HANYA dipanggil saat create/update draft, TIDAK saat publish
func (s *DraftServiceImpl) injectMetaTagsToArticle(
	article string,
	title string,
	topic string,
	excerpt string,
	imageURL string,
) string {
	if article == "" {
		return article
	}

	// Generate meta tags
	metaTags := helper.GenerateMetaTags(title, topic, excerpt, imageURL)

	// Inject meta tags
	if strings.Contains(article, "<head>") {
		article = strings.Replace(article, "<head>", "<head>\n"+metaTags, 1)
	} else if strings.Contains(article, "<html") {
		article = strings.Replace(article, "<html", "<html>\n<head>\n"+metaTags+"\n</head>", 1)
	} else {
		article = "<!DOCTYPE html>\n<html>\n<head>\n" + metaTags + "\n</head>\n<body>\n" + article + "\n</body>\n</html>"
	}

	return article
}

// insertHistory helper untuk insert history
func (s *DraftServiceImpl) insertHistory(
	ctx context.Context,
	historyReq draft.PublishHistoryRequest,
	userCtx models.UserContext,
	status string,
) {
	if err := s.repoHistory.Create(ctx, historyReq, userCtx.GetUserID(), userCtx.GetTeamID(), status); err != nil {
		log.Printf("[ERROR] InsertHistory(%s) slug=%s err=%v", status, historyReq.Slug, err)
	} else {
		log.Printf("[SUCCESS] InsertHistory(%s) slug=%s", status, historyReq.Slug)
	}
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
	if err := s.redisScheduler.CancelScheduledTask(ctx, draftID); err != nil {
		log.Printf("Failed to cancel Redis task: %v", err)
	}
	return s.repo.UpdateStatus(ctx, draftID, "draft", nil)
}

// scheduleDraft handle schedule logic
func (s *DraftServiceImpl) scheduleDraft(ctx context.Context, id, scheduledForStr string, draftData *draft.DraftData, userCtx models.UserContext) (*draft.PublishResult, error) {
	scheduledFor := helper.ParseWIBTime(scheduledForStr)

	if err := s.repo.UpdateStatus(ctx, id, "scheduled", &scheduledFor); err != nil {
		log.Printf("ERROR : %v", err)
		return nil, err
	}

	imageURL := ""
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

// processPublish handle publish logic
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

	score := draft.CalculateSEOScore(draftData.Title, draftData.Article, draftData.Topic, draftData.Excerpt, draftData.Keywords)
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

// ==================== HELPER FUNCTIONS ====================

func prepareUpdateData(updates draft.CreateDraftRequest, minioService *minio.MinioService) draft.CreateDraftRequest {
	log.Println("========== PREPARE UPDATE DATA ==========")
	ctx := context.Background()
	processedUpdates := updates
	hasImage := false

	// Process image_url
	if processedUpdates.ImageURL != nil && *processedUpdates.ImageURL != "" {
		if isBase64Image(processedUpdates.ImageURL) {
			uploadedURL, err := uploadBase64ToMinio(ctx, minioService, *processedUpdates.ImageURL)
			if err != nil {
				log.Printf("[ERROR] Failed to upload base64 image: %v", err)
				emptyStr := ""
				processedUpdates.ImageURL = &emptyStr
			} else {
				processedUpdates.ImageURL = &uploadedURL
				hasImage = true
			}
		} else if strings.Contains(*processedUpdates.ImageURL, minioService.Bucket) {
			objectName := extractObjectName(*processedUpdates.ImageURL)
			processedUpdates.ImageURL = &objectName
			hasImage = true
		} else {
			hasImage = true
		}
	}

	// Process article images
	if processedUpdates.Article != "" {
		processedArticle := processArticleImages(ctx, minioService, processedUpdates.Article)
		processedUpdates.Article = processedArticle
		if containsImageInArticle(processedArticle) {
			hasImage = true
		}
	}

	processedUpdates.HasImage = hasImage
	return processedUpdates
}

func containsImageInArticle(article string) bool {
	if strings.Contains(article, "<img ") || strings.Contains(article, "<img>") {
		return true
	}
	matched, _ := regexp.MatchString(`!\[.*\]\(.*\)`, article)
	if matched {
		return true
	}
	if strings.Contains(article, "data:image/") {
		return true
	}
	return false
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

		if isBase64Image(&oldURL) {
			newURL, err := uploadBase64ToMinio(ctx, minioService, oldURL)
			if err != nil {
				log.Printf("[ERROR] Failed to upload: %v", err)
				return tag
			}
			return strings.Replace(tag, oldURL, newURL, 1)
		}

		if strings.Contains(oldURL, minioService.Bucket) {
			objectName := extractObjectName(oldURL)
			return strings.Replace(tag, oldURL, objectName, 1)
		}

		return tag
	})
}

func extractObjectName(imageURL string) string {
	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return imageURL
	}
	path := strings.TrimPrefix(parsedURL.Path, "/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return path
}

func isBase64Image(data *string) bool {
	base64Pattern := regexp.MustCompile(`^data:image\/(png|jpeg|jpg|gif|webp|svg\+xml);base64,`)
	return base64Pattern.MatchString(*data)
}

func uploadBase64ToMinio(ctx context.Context, minioService *minio.MinioService, base64Data string) (string, error) {
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

	reader := bytes.NewReader(imageData)
	uploadedResult, err := minioService.Upload(ctx, objectName, reader, int64(len(imageData)), mimeType)
	if err != nil {
		return "", fmt.Errorf("failed to upload image to Minio: %w", err)
	}

	return uploadedResult, nil
}

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

func (s *DraftServiceImpl) RescheduleDraft(
	ctx context.Context,
	draftID string,
	newScheduleTime time.Time,
	userCtx models.UserContext,
) (*draft.PublishResult, error) {
	log.Printf("========== RESCHEDULE DRAFT SERVICE ==========")

	existingDraft, err := s.repo.GetByID(ctx, draftID)
	if err != nil {
		return nil, fmt.Errorf("draft not found")
	}

	if existingDraft.Status != "scheduled" {
		return nil, fmt.Errorf("draft is not in scheduled status")
	}

	if err := s.redisScheduler.CancelScheduledTask(ctx, draftID); err != nil {
		log.Printf("[WARNING] Failed to cancel existing schedule: %v", err)
	}

	oldScheduledTime := existingDraft.ScheduledFor

	updateRequest := draft.CreateDraftRequest{
		Title:        existingDraft.Title,
		Topic:        existingDraft.Topic,
		Article:      existingDraft.Article,
		ImageURL:     existingDraft.ImageURL,
		ImagePrompt:  existingDraft.ImagePrompt,
		ScheduledFor: newScheduleTime.Format(time.RFC3339),
		Slug:         existingDraft.Slug,
		SEOScore:     existingDraft.SEOScore,
		Excerpt:      existingDraft.Excerpt,
		UpdateAt:     time.Now(),
		Status:       "scheduled",
		TeamID:       userCtx.GetTeamID(),
		UserID:       userCtx.GetUserID(),
	}

	if err := s.repo.Update(ctx, userCtx.GetTeamID(), draftID, updateRequest); err != nil {
		return nil, fmt.Errorf("failed to update draft schedule: %w", err)
	}

	taskData := &scheduler.ScheduledTask{
		ID:             draftID,
		Title:          existingDraft.Title,
		Topic:          existingDraft.Topic,
		Article:        existingDraft.Article,
		ImagePrompt:    existingDraft.ImagePrompt,
		TargetProducts: existingDraft.TargetProducts,
		TeamID:         userCtx.GetTeamID(),
		UserID:         userCtx.GetUserID(),
	}

	if err := s.redisScheduler.ScheduleDraftTask(ctx, draftID, newScheduleTime, taskData, userCtx); err != nil {
		rollbackRequest := updateRequest
		rollbackRequest.Status = "draft"
		rollbackRequest.ScheduledFor = ""
		if rollbackErr := s.repo.Update(ctx, userCtx.GetTeamID(), draftID, rollbackRequest); rollbackErr != nil {
			log.Printf("[CRITICAL] Failed to rollback draft status: %v", rollbackErr)
		}
		return nil, fmt.Errorf("failed to create new schedule: %w", err)
	}

	return &draft.PublishResult{
		Status:       "scheduled",
		ScheduledFor: &newScheduleTime,
		Results: []draft.PublishResult{
			{
				Status: "scheduled",
				Message: fmt.Sprintf("Rescheduled from %s to %s",
					oldScheduledTime,
					newScheduleTime,
				),
			},
		},
		TotalProducts: 1,
		SuccessCount:  1,
		FailedCount:   0,
	}, nil
}

// enrichDraftWithPreviousResults
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
