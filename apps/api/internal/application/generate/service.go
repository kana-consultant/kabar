// internal/infrastructure/ai/service/generate_service.go
package generate

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"

	"seo-backend/internal/domain/generate"
	"seo-backend/internal/helper"
	"seo-backend/internal/infrastructure/ai/builder"
	"seo-backend/internal/infrastructure/ai/parser"
	"seo-backend/internal/infrastructure/http/client"
	"seo-backend/internal/infrastructure/http/minio"
)

const (
	imageDownloadTimeout   = 60 * time.Second
	imageUploadTimeout     = 60 * time.Second
	imageProcessTimeout    = 5 * time.Minute // budget total utk semua gambar dalam 1 artikel
	imageDownloadRetries   = 2
	imageGenerationTimeout = 60 * time.Second // timeout per image generation
)

type GenerateServiceImpl struct {
	repo           generate.Repository
	httpClient     *client.HTTPClient
	minioClient    *minio.MinioService
	promptBuilder  *builder.PromptBuilder
	requestBuilder *builder.RequestBuilder
	responseParser *parser.ResponseParser
}

func NewService(
	repo generate.Repository,
	httpClient *client.HTTPClient,
	minioClient *minio.MinioService,
	promptBuilder *builder.PromptBuilder,
	requestBuilder *builder.RequestBuilder,
	responseParser *parser.ResponseParser,
) generate.Service {
	return &GenerateServiceImpl{
		repo:           repo,
		httpClient:     httpClient,
		promptBuilder:  promptBuilder,
		requestBuilder: requestBuilder,
		responseParser: responseParser,
		minioClient:    minioClient,
	}
}

func (s *GenerateServiceImpl) GenerateArticle(ctx context.Context, params generate.ArticleGenerationParams) (*generate.ArticleResult, error) {
	log.Println("========== GENERATE ARTICLE ==========")
	log.Printf("[INFO] Starting article generation with topic: %s", params.Topic)
	log.Printf("[INFO] Model ID: %s, Tone: %s, Length: %s, Language: %s, AutoGenerateImage: %v, ImageModelID: %s",
		params.ModelID, params.Tone, params.Length, params.Language, params.AutoGenerateImage, params.ImageModelID)

	defer func() {
		log.Println("========== END GENERATE ARTICLE ==========")
	}()

	// Validate params
	log.Println("[INFO] Validating request parameters...")
	if err := s.validateArticleParams(params); err != nil {
		log.Printf("[ERROR] Parameter validation failed: %v", err)
		return nil, err
	}
	log.Println("[INFO] Parameters validation passed")

	// Get model configuration
	log.Printf("[INFO] Fetching model configuration for model ID: %s (service: text)", params.ModelID)
	config, err := s.repo.GetModelConfig(ctx, params.ModelID, "text")
	if err != nil {
		log.Printf("[ERROR] Failed to get model config: %v", err)
		return nil, fmt.Errorf("failed to get model config: %w", err)
	}
	log.Printf("[INFO] Model config loaded - Model: %s, BaseURL: %s, Endpoint: %s, MaxTokens: %d, Temperature: %.2f",
		config.ModelName, config.BaseURL, config.Endpoint, config.MaxTokens, config.Temperature)

	// Build prompt
	log.Println("[INFO] Building article prompt...")
	prompt := s.promptBuilder.BuildArticlePrompt(params)

	// Combine system prompt
	fullPrompt := config.SystemPrompt + "\n\n" + prompt
	log.Printf("[INFO] Prompt built successfully (length: %d characters)", len(fullPrompt))

	// Build request body with max_tokens and temperature
	log.Println("[INFO] Building request body...")
	requestBody, err := s.requestBuilder.BuildArticleRequestBody(config, fullPrompt)
	if err != nil {
		log.Printf("[ERROR] Failed to build request: %v", err)
		return nil, fmt.Errorf("failed to build request: %w", err)
	}
	log.Printf("[INFO] Request body built successfully (size: %d bytes)", len(requestBody))

	// Send request
	log.Printf("[INFO] Sending request to AI provider (timeout: 300)...")
	response, err := s.httpClient.SendRequest(ctx, config, requestBody, 300*time.Second)
	if err != nil {
		log.Printf("[ERROR] Failed to send request: %v", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	log.Printf("[INFO] Response received successfully (size: %d bytes)", len(response))

	// Parse response
	log.Printf("[INFO] Parsing response with path: %s", config.ResponsePath)
	result, err := s.responseParser.ParseArticleResponse(response, config.ResponsePath)
	if err != nil {
		log.Printf("[ERROR] Failed to parse response: %v", err)
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	log.Printf("[INFO] Response parsed successfully")

	// 🔥 PROSES GAMBAR - handle both existing images and image placeholders
	if result.Content != "" {
		log.Println("[INFO] Processing images from AI-generated content...")

		imgCtx, cancel := context.WithTimeout(context.Background(), imageProcessTimeout)
		defer cancel()

		var processedContent string
		var stats imageProcessStats
		var processErr error

		// Check if auto-generate images is enabled
		if params.AutoGenerateImage && params.ImageModelID != "" {
			// 🔥 AUTO GENERATE IMAGES FROM PLACEHOLDERS
			log.Println("[INFO] Auto-generate images enabled, processing image placeholders...")
			processedContent, stats, processErr = s.processImagePlaceholders(imgCtx, result.Content, params)
		} else if params.AutoGenerateImage && params.ImageModelID == "" {
			// 🔥 FIX: Auto-generate enabled but no model ID - keep placeholders as-is
			log.Println("[WARNING] Auto-generate images is enabled but no image model ID provided")

			// Check if there are image placeholders
			if s.hasImagePlaceholders(result.Content) {
				log.Println("[INFO] Image placeholders found but no image model configured")
				log.Println("[INFO] Keeping placeholders in content (they will remain as <img prompt='...'>)")

				// Count placeholders
				placeholderCount := strings.Count(result.Content, "<img prompt=")
				stats = imageProcessStats{
					Total:   placeholderCount,
					Skipped: placeholderCount,
				}
				processedContent = result.Content
				processErr = nil
			} else {
				// No placeholders, process existing images with src
				log.Println("[INFO] No image placeholders found, processing existing images...")
				processedContent, stats, processErr = s.processImages(imgCtx, result.Content, params)
			}
		} else {
			// 🔥 AUTO GENERATE DISABLED - process existing images with src
			log.Println("[INFO] Auto-generate images disabled, processing existing images...")
			processedContent, stats, processErr = s.processImages(imgCtx, result.Content, params)
		}

		if processErr != nil {
			log.Printf("[WARNING] Failed to process images: %v", processErr)
			// Jangan gagalkan proses, lanjutkan dengan content asli
		} else {
			result.Content = processedContent
			log.Printf("[INFO] Images processed - total: %d, replaced: %d, failed: %d, skipped: %d",
				stats.Total, stats.Replaced, stats.Failed, stats.Skipped)

			if stats.Total > 0 && stats.Replaced == 0 {
				log.Printf("[WARNING] No images were successfully processed; original content kept")
			} else if stats.Failed > 0 {
				log.Printf("[WARNING] %d/%d images failed to process", stats.Failed, stats.Total)
			}
		}
	}

	// Hitung SEO score otomatis
	result.SeoScore = helper.CalculateSEOScoreSimple(result.Title, result.Content, result.Excerpt, params.Topic)

	log.Printf("SUCCESS title=%s words=%d seo_score=%d slug=%s", result.Title, result.WordCount, result.SeoScore, result.Slug)

	return result, nil
}

// imageProcessStats - ringkasan hasil proses gambar dalam satu artikel
type imageProcessStats struct {
	Total    int
	Replaced int
	Failed   int
	Skipped  int // img tanpa src atau prompt
}

// 🔥 HELPER: Check if content has image placeholders
func (s *GenerateServiceImpl) hasImagePlaceholders(content string) bool {
	return strings.Contains(content, "<img prompt=")
}

// 🔥 FUNGSI BARU: Proses placeholder <img prompt="..."> dan generate gambar
func (s *GenerateServiceImpl) processImagePlaceholders(ctx context.Context, htmlContent string, params generate.ArticleGenerationParams) (string, imageProcessStats, error) {
	log.Println("[INFO] Starting image placeholder processing...")

	stats := imageProcessStats{}

	// Parse HTML untuk mencari semua tag img dengan attribute prompt
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return "", stats, fmt.Errorf("failed to parse HTML: %w", err)
	}

	imgSelections := doc.Find("img[prompt]")
	stats.Total = imgSelections.Length()

	if stats.Total == 0 {
		log.Println("[INFO] No image placeholders found in content")
		return htmlContent, stats, nil
	}

	log.Printf("[INFO] Found %d image placeholders to generate", stats.Total)

	var mu sync.Mutex
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 3) // Limit concurrent image generation to 3

	imgSelections.Each(func(i int, img *goquery.Selection) {
		if ctx.Err() != nil {
			log.Printf("[WARNING] Context done, stopping image generation")
			mu.Lock()
			stats.Failed++
			mu.Unlock()
			return
		}

		// Get prompt attribute
		prompt, exists := img.Attr("prompt")
		if !exists || prompt == "" {
			log.Printf("[WARNING] Image #%d has empty prompt, skipping", i+1)
			mu.Lock()
			stats.Skipped++
			mu.Unlock()
			return
		}

		wg.Add(1)
		go func(index int, promptText string, imgElement *goquery.Selection) {
			defer wg.Done()

			// Acquire semaphore
			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			case <-ctx.Done():
				log.Printf("[WARNING] Context done before processing image #%d", index+1)
				mu.Lock()
				stats.Failed++
				mu.Unlock()
				return
			}

			log.Printf("[INFO] Generating image #%d with prompt: %s", index+1, promptText)

			// Generate image using AI
			imageParams := generate.ImageGenerationParams{
				Prompt:    promptText,
				ModelID:   params.ImageModelID,
				Slug:      params.Slug,
				ArticleID: params.ArticleID,
			}

			// Generate image dengan timeout
			imageCtx, cancel := context.WithTimeout(ctx, imageGenerationTimeout)
			defer cancel()

			imageResult, err := s.GenerateImage(imageCtx, imageParams)
			if err != nil {
				log.Printf("[ERROR] Failed to generate image #%d: %v", index+1, err)
				mu.Lock()
				stats.Failed++
				mu.Unlock()
				return
			}

			log.Printf("[INFO] Image #%d generated successfully: %s", index+1, imageResult.ImageURL)

			// Replace img tag with generated image
			// Remove prompt attribute and set src
			imgElement.RemoveAttr("prompt")
			imgElement.SetAttr("src", imageResult.ImageURL)

			// Add alt attribute for accessibility
			imgElement.SetAttr("alt", promptText)

			mu.Lock()
			stats.Replaced++
			mu.Unlock()
			log.Printf("[INFO] Image #%d placeholder replaced with generated image", index+1)
		}(i, prompt, img)
	})

	wg.Wait()

	log.Printf("[INFO] Image placeholders processed - total: %d, replaced: %d, failed: %d, skipped: %d",
		stats.Total, stats.Replaced, stats.Failed, stats.Skipped)

	// Generate HTML baru
	newHTML, err := doc.Html()
	if err != nil {
		return "", stats, fmt.Errorf("failed to generate new HTML: %w", err)
	}

	return newHTML, stats, nil
}

// 🔥 FUNGSI: Proses semua gambar di HTML (download dan rehost ke MinIO)
func (s *GenerateServiceImpl) processImages(ctx context.Context, htmlContent string, params generate.ArticleGenerationParams) (string, imageProcessStats, error) {
	log.Println("[INFO] Starting image processing...")

	stats := imageProcessStats{}

	// Parse HTML untuk mencari semua tag img
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return "", stats, fmt.Errorf("failed to parse HTML: %w", err)
	}

	imgSelections := doc.Find("img")
	stats.Total = imgSelections.Length()

	var mu sync.Mutex
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 5) // Limit concurrent downloads

	imgSelections.Each(func(i int, img *goquery.Selection) {
		// Kalau budget waktu keseluruhan sudah habis, hentikan lebih awal
		if ctx.Err() != nil {
			log.Printf("[WARNING] Image processing context done (%v), skipping remaining images", ctx.Err())
			mu.Lock()
			stats.Failed++
			mu.Unlock()
			return
		}

		log.Printf("[INFO] Processing image #%d", i+1)

		// Dapatkan src attribute
		src, exists := img.Attr("src")
		if !exists || src == "" {
			log.Printf("[WARNING] Image #%d has no src attribute, skipping", i+1)
			mu.Lock()
			stats.Skipped++
			mu.Unlock()
			return
		}

		// Skip if it's already a MinIO URL (to avoid re-processing)
		if strings.Contains(src, "minio") || strings.Contains(src, "storage.googleapis.com") {
			log.Printf("[INFO] Image #%d already uses MinIO/storage URL, skipping", i+1)
			mu.Lock()
			stats.Skipped++
			mu.Unlock()
			return
		}

		wg.Add(1)
		go func(index int, imageSrc string, imgElement *goquery.Selection) {
			defer wg.Done()

			// Acquire semaphore
			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			case <-ctx.Done():
				log.Printf("[WARNING] Context done before processing image #%d", index+1)
				mu.Lock()
				stats.Failed++
				mu.Unlock()
				return
			}

			log.Printf("[INFO] Image #%d src: %s", index+1, imageSrc)

			// Download gambar (dengan retry, lebih resilient terhadap kegagalan sementara)
			imageData, contentType, err := s.downloadImageWithRetry(ctx, imageSrc, imageDownloadRetries)
			if err != nil {
				log.Printf("[ERROR] Failed to download image #%d: %v", index+1, err)
				mu.Lock()
				stats.Failed++
				mu.Unlock()
				return
			}

			log.Printf("[INFO] Image #%d downloaded (size: %d bytes, type: %s)", index+1, len(imageData), contentType)

			// Upload ke MinIO
			minioURL, err := s.uploadToMinio(ctx, imageData, contentType, params)
			if err != nil {
				log.Printf("[ERROR] Failed to upload image #%d to MinIO: %v", index+1, err)
				mu.Lock()
				stats.Failed++
				mu.Unlock()
				return
			}

			log.Printf("[INFO] Image #%d uploaded to MinIO: %s", index+1, minioURL)

			// Ganti src attribute dengan URL MinIO
			imgElement.SetAttr("src", minioURL)
			mu.Lock()
			stats.Replaced++
			mu.Unlock()
			log.Printf("[INFO] Image #%d src replaced with MinIO URL", index+1)
		}(i, src, img)
	})

	wg.Wait()

	log.Printf("[INFO] Total images found: %d, replaced: %d, failed: %d, skipped: %d",
		stats.Total, stats.Replaced, stats.Failed, stats.Skipped)

	// Generate HTML baru
	newHTML, err := doc.Html()
	if err != nil {
		return "", stats, fmt.Errorf("failed to generate new HTML: %w", err)
	}

	return newHTML, stats, nil
}

// 🔥 FUNGSI: Download gambar dari URL
func (s *GenerateServiceImpl) downloadImage(ctx context.Context, imageURL string) ([]byte, string, error) {
	log.Printf("[INFO] Downloading image from: %s", imageURL)

	// Handle base64 image
	if strings.HasPrefix(imageURL, "data:image") {
		log.Println("[INFO] Detected base64 image, decoding...")
		return s.decodeBase64Image(imageURL)
	}

	// 🔥 Sanitasi URL sebelum request
	sanitizedURL, sanitizeErr := sanitizeImageURL(imageURL)
	if sanitizeErr != nil {
		log.Printf("[WARNING] Failed to sanitize image URL, using as-is: %v", sanitizeErr)
		sanitizedURL = imageURL
	} else if sanitizedURL != imageURL {
		log.Printf("[INFO] URL sanitized for request: %s -> %s", imageURL, sanitizedURL)
	}

	// 🔥 BUAT CONTEXT BARU DENGAN TIMEOUT LEBIH LAMA
	downloadCtx, cancel := context.WithTimeout(ctx, imageDownloadTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(downloadCtx, "GET", sanitizedURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "image/webp,image/apng,image/*,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	client := &http.Client{
		Timeout: imageDownloadTimeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyPreview, _ := io.ReadAll(io.LimitReader(resp.Body, 500))
		return nil, "", fmt.Errorf("download failed with status: %s, body: %s", resp.Status, string(bodyPreview))
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read image data: %w", err)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(imageData)
	}

	log.Printf("[INFO] Image downloaded (size: %d bytes, type: %s)", len(imageData), contentType)

	return imageData, contentType, nil
}

// sanitizeImageURL membersihkan URL gambar/chart yang berasal dari AI
func sanitizeImageURL(rawURL string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "", fmt.Errorf("empty url")
	}

	// data: URL (base64 atau svg mentah) tidak disentuh sama sekali
	if strings.HasPrefix(rawURL, "data:") {
		return rawURL, nil
	}

	idx := strings.Index(rawURL, "?")
	if idx == -1 {
		// Tidak ada query string sama sekali
		if _, err := url.ParseRequestURI(rawURL); err != nil {
			return "", fmt.Errorf("invalid url (no query): %w", err)
		}
		return rawURL, nil
	}

	base := rawURL[:idx]
	rawQuery := rawURL[idx+1:]

	if _, err := url.ParseRequestURI(base); err != nil {
		return "", fmt.Errorf("invalid base url: %w", err)
	}

	// Pisahkan jadi pasangan key=value berdasarkan "&" di level TOP saja
	pairs := splitTopLevelQueryPairs(rawQuery)

	rebuilt := make([]string, 0, len(pairs))
	for _, pair := range pairs {
		eqIdx := strings.Index(pair, "=")
		if eqIdx == -1 {
			rebuilt = append(rebuilt, pair)
			continue
		}

		key := pair[:eqIdx]
		value := pair[eqIdx+1:]

		normalizedValue := normalizeQueryValue(value)
		rebuilt = append(rebuilt, key+"="+url.QueryEscape(normalizedValue))
	}

	sanitized := base + "?" + strings.Join(rebuilt, "&")

	if _, err := url.Parse(sanitized); err != nil {
		return "", fmt.Errorf("sanitized url still invalid: %w", err)
	}

	return sanitized, nil
}

// splitTopLevelQueryPairs memisahkan query string jadi pasangan key=value
func splitTopLevelQueryPairs(rawQuery string) []string {
	var pairs []string
	depth := 0
	start := 0

	for i, r := range rawQuery {
		switch r {
		case '{', '[':
			depth++
		case '}', ']':
			if depth > 0 {
				depth--
			}
		case '&':
			if depth == 0 {
				pairs = append(pairs, rawQuery[start:i])
				start = i + 1
			}
		}
	}
	pairs = append(pairs, rawQuery[start:])

	return pairs
}

// normalizeQueryValue membuat proses idempotent
func normalizeQueryValue(value string) string {
	if decoded, err := url.QueryUnescape(value); err == nil {
		return decoded
	}
	return value
}

// decodeBase64Image - Versi sederhana
func (s *GenerateServiceImpl) decodeBase64Image(dataURL string) ([]byte, string, error) {
	// Cek prefix
	if !strings.HasPrefix(dataURL, "data:image") {
		return nil, "", fmt.Errorf("invalid base64 image format")
	}

	// Pisahkan header dan data
	parts := strings.SplitN(dataURL, ",", 2)
	if len(parts) != 2 {
		return nil, "", fmt.Errorf("invalid base64 format")
	}

	// Ekstrak content type
	header := parts[0]
	contentType := strings.TrimPrefix(header, "data:")
	contentType = strings.Split(contentType, ";")[0]

	// Decode base64
	imageData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Jika contentType kosong, detect dari data
	if contentType == "" {
		contentType = http.DetectContentType(imageData)
	}

	return imageData, contentType, nil
}

// downloadImageWithRetry - dipanggil oleh processImages
func (s *GenerateServiceImpl) downloadImageWithRetry(ctx context.Context, imageURL string, maxRetries int) ([]byte, string, error) {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if ctx.Err() != nil {
			return nil, "", fmt.Errorf("context done before retry %d: %w", attempt+1, ctx.Err())
		}

		if attempt > 0 {
			log.Printf("[INFO] Retry attempt %d for: %s", attempt+1, imageURL)

			backoff := time.Duration(attempt*2) * time.Second
			timer := time.NewTimer(backoff)
			select {
			case <-ctx.Done():
				timer.Stop()
				return nil, "", fmt.Errorf("context done during backoff: %w", ctx.Err())
			case <-timer.C:
			}
		}

		imageData, contentType, err := s.downloadImage(ctx, imageURL)
		if err == nil {
			return imageData, contentType, nil
		}

		lastErr = err
		log.Printf("[WARNING] Download attempt %d failed: %v", attempt+1, err)
	}

	return nil, "", fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}

// uploadToMinio - upload ke MinIO
func (s *GenerateServiceImpl) uploadToMinio(ctx context.Context, imageData []byte, contentType string, params generate.ArticleGenerationParams) (string, error) {
	log.Println("[INFO] Uploading image to MinIO...")

	// Generate object name
	ext := getExtensionFromContentType(contentType)
	timestamp := time.Now().Unix()
	random := generateRandomString(8)

	// Gunakan slug dari params atau generate dari topic
	slug := params.Slug
	if slug == "" {
		slug = helper.GenerateSlug(params.Topic)
	}

	objectName := fmt.Sprintf("articles/%s/%d_%s.%s", slug, timestamp, random, ext)
	log.Printf("[INFO] Object name: %s", objectName)

	uploadCtx, cancel := context.WithTimeout(ctx, imageUploadTimeout)
	defer cancel()

	// Upload ke MinIO
	reader := bytes.NewReader(imageData)
	uploadedName, err := s.minioClient.Upload(uploadCtx, objectName, reader, int64(len(imageData)), contentType)
	if err != nil {
		log.Printf("[ERROR] Failed to upload to MinIO: %v", err)
		return "", fmt.Errorf("failed to upload image to MinIO: %w", err)
	}
	log.Printf("[DEBUG] Uploaded to MinIO: %s", uploadedName)

	// Ambil presigned URL dengan expiry 7 hari
	imageURLResult, err := s.minioClient.GetURL(uploadCtx, uploadedName, 7*24*time.Hour)
	if err != nil {
		log.Printf("[ERROR] Failed to get MinIO URL: %v", err)
		return "", fmt.Errorf("failed to get MinIO URL: %w", err)
	}
	log.Printf("[DEBUG] Presigned URL generated: %s", imageURLResult)

	return imageURLResult, nil
}

// 🔥 Update GenerateImage method
func (s *GenerateServiceImpl) GenerateImage(ctx context.Context, params generate.ImageGenerationParams) (*generate.ImageResult, error) {
	log.Printf("[DEBUG] GenerateImage started, modelID: %s, prompt: %s", params.ModelID, params.Prompt)

	// Validate params
	if err := s.validateImageParams(params); err != nil {
		log.Printf("[ERROR] Validation failed: %v", err)
		return nil, err
	}

	// Get model configuration
	config, err := s.repo.GetModelConfig(ctx, params.ModelID, "image")
	if err != nil {
		log.Printf("[ERROR] Failed to get model config: %v", err)
		return nil, fmt.Errorf("failed to get model config: %w", err)
	}
	log.Printf("[DEBUG] Model config loaded: %s, MaxTokens: %d, Temperature: %.2f",
		config.ModelName, config.MaxTokens, config.Temperature)

	// Build image prompt with system prompt
	imagePrompt := params.Prompt
	if config.SystemPrompt != "" {
		imagePrompt = config.SystemPrompt + " " + params.Prompt
	}
	log.Printf("[DEBUG] Image prompt built: %s", imagePrompt)

	// Build request body with max_tokens and temperature
	requestBody, err := s.requestBuilder.BuildImageRequestBody(config, imagePrompt)
	if err != nil {
		log.Printf("[ERROR] Failed to build request body: %v", err)
		return nil, fmt.Errorf("failed to build request: %w", err)
	}
	log.Printf("[DEBUG] Request body built successfully")

	// Send request with longer timeout for images
	response, err := s.httpClient.SendRequest(ctx, config, requestBody, 120*time.Second)
	if err != nil {
		log.Printf("[ERROR] Failed to send request to AI provider: %v", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	log.Printf("[DEBUG] AI provider response received")

	// Parse response - dapatkan Base64 string dari AI provider
	base64String, contentType, err := s.responseParser.ParseImageResponseBase64(response, config.ResponseImagePath)
	if err != nil {
		log.Printf("[ERROR] Failed to parse image response: %v", err)
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	log.Printf("[DEBUG] Image parsed, contentType: %s, base64 length: %d", contentType, len(base64String))

	// Remove data URL prefix jika ada
	base64String = s.cleanBase64String(base64String)

	// Decode Base64 ke binary
	imageData, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		log.Printf("[ERROR] Failed to decode base64: %v", err)
		return nil, fmt.Errorf("failed to decode base64 image: %w", err)
	}
	log.Printf("[DEBUG] Base64 decoded, image size: %d bytes", len(imageData))

	// Tentukan ekstensi berdasarkan content-type
	ext := s.getImageExtension(contentType)
	log.Printf("[DEBUG] Image extension: %s", ext)

	// Generate object name - gunakan slug dari params jika tersedia
	var objectName string
	if params.Slug != "" {
		timestamp := time.Now().UnixNano()
		random := generateRandomString(8)
		objectName = fmt.Sprintf("articles/%s/images/%d_%s.%s", params.Slug, timestamp, random, ext)
	} else if params.ArticleID != "" {
		objectName = fmt.Sprintf("articles/%s/images/%d.%s", params.ArticleID, time.Now().UnixNano(), ext)
	} else {
		objectName = fmt.Sprintf("images/%d.%s", time.Now().UnixNano(), ext)
	}
	log.Printf("[DEBUG] Object name: %s", objectName)

	// 🔥 Upload ke MinIO dengan timeout terpisah
	uploadCtx, cancel := context.WithTimeout(context.Background(), imageUploadTimeout)
	defer cancel()

	// Upload ke MinIO
	reader := bytes.NewReader(imageData)
	uploadedName, err := s.minioClient.Upload(uploadCtx, objectName, reader, int64(len(imageData)), contentType)
	if err != nil {
		log.Printf("[ERROR] Failed to upload to Minio: %v", err)
		return nil, fmt.Errorf("failed to upload image to Minio: %w", err)
	}
	log.Printf("[DEBUG] Uploaded to Minio: %s", uploadedName)

	// Ambil presigned URL dengan expiry 7 hari
	imageURL, err := s.minioClient.GetURL(uploadCtx, uploadedName, 7*24*time.Hour)
	if err != nil {
		log.Printf("[ERROR] Failed to get Minio URL: %v", err)
		return nil, fmt.Errorf("failed to get Minio URL: %w", err)
	}
	log.Printf("[DEBUG] Presigned URL generated: %s", imageURL)

	result := &generate.ImageResult{
		ImageURL:    imageURL,
		Prompt:      params.Prompt,
		GeneratedAt: s.nowWIB().Format(time.RFC3339),
		Model:       config.ModelName,
	}

	log.Printf("[DEBUG] GenerateImage completed successfully, imageURL: %s", imageURL)
	return result, nil
}

// 🔥 HELPER: Get extension from content type
func getExtensionFromContentType(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return "jpg"
	case "image/png":
		return "png"
	case "image/gif":
		return "gif"
	case "image/webp":
		return "webp"
	case "image/svg+xml":
		return "svg"
	case "image/bmp":
		return "bmp"
	default:
		return "jpg"
	}
}

// 🔥 HELPER: Generate random string
var (
	randMu  sync.Mutex
	randSrc = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)

	randMu.Lock()
	for i := range b {
		b[i] = charset[randSrc.Intn(len(charset))]
	}
	randMu.Unlock()

	return string(b)
}

// Private helper methods
func (s *GenerateServiceImpl) validateArticleParams(params generate.ArticleGenerationParams) error {
	if params.Topic == "" {
		return fmt.Errorf("topic is required")
	}
	if params.ModelID == "" {
		return fmt.Errorf("modelId is required")
	}
	return nil
}

func (s *GenerateServiceImpl) validateImageParams(params generate.ImageGenerationParams) error {
	if params.Prompt == "" {
		return fmt.Errorf("prompt is required")
	}
	if params.ModelID == "" {
		return fmt.Errorf("modelId is required")
	}
	return nil
}

func (s *GenerateServiceImpl) saveHistory(ctx context.Context, params generate.ArticleGenerationParams, result *generate.ArticleResult) error {
	history := &generate.GenerationHistory{
		Type:      "article",
		Topic:     params.Topic,
		Result:    result.Content,
		ModelID:   params.ModelID,
		CreatedAt: s.nowWIB(),
	}
	return s.repo.SaveHistory(ctx, history)
}

func (s *GenerateServiceImpl) nowWIB() time.Time {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	return time.Now().In(loc)
}

func (s *GenerateServiceImpl) cleanBase64String(base64Str string) string {
	if strings.Contains(base64Str, "base64,") {
		parts := strings.Split(base64Str, "base64,")
		base64Str = parts[len(parts)-1]
	}
	base64Str = strings.TrimPrefix(base64Str, "data:image/png;base64,")
	base64Str = strings.TrimPrefix(base64Str, "data:image/jpeg;base64,")
	base64Str = strings.TrimPrefix(base64Str, "data:image/webp;base64,")
	base64Str = strings.TrimPrefix(base64Str, "data:image/gif;base64,")
	return base64Str
}

func (s *GenerateServiceImpl) getImageExtension(contentType string) string {
	switch contentType {
	case "image/jpeg", "image/jpg":
		return "jpg"
	case "image/webp":
		return "webp"
	case "image/gif":
		return "gif"
	default:
		return "png"
	}
}
