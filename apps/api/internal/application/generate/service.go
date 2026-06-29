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
	imageDownloadTimeout = 60 * time.Second
	imageUploadTimeout   = 60 * time.Second
	imageProcessTimeout  = 5 * time.Minute // budget total utk semua gambar dalam 1 artikel
	imageDownloadRetries = 2
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
	log.Printf("[INFO] Model ID: %s, Tone: %s, Length: %s, Language: %s",
		params.ModelID, params.Tone, params.Length, params.Language)

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

	// 🔥 PROSES GAMBAR LANGSUNG (tanpa cek auto generate)
	// Selalu proses gambar yang ada di content
	if result.Content != "" {
		log.Println("[INFO] Processing images from AI-generated content...")

		// 🔥 PENTING: jangan pakai ctx request (HTTP) yang mungkin
		// sudah hampir/sudah expired setelah AI generation (bisa makan puluhan detik).
		// Image pipeline diberi budget waktu sendiri yang independen.
		imgCtx, cancel := context.WithTimeout(context.Background(), imageProcessTimeout)
		processedContent, stats, err := s.processImages(imgCtx, result.Content, params)
		cancel()

		if err != nil {
			log.Printf("[WARNING] Failed to process images: %v", err)
			// Jangan gagalkan proses, lanjutkan dengan content asli
		} else {
			result.Content = processedContent
			log.Printf("[INFO] Images processed - total: %d, replaced: %d, failed: %d, skipped: %d",
				stats.Total, stats.Replaced, stats.Failed, stats.Skipped)

			if stats.Total > 0 && stats.Replaced == 0 {
				log.Printf("[WARNING] No images were successfully rehosted to MinIO; original source URLs are kept in content")
			} else if stats.Failed > 0 {
				log.Printf("[WARNING] %d/%d images failed to rehost; their original source URLs are kept in content", stats.Failed, stats.Total)
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
	Skipped  int // img tanpa src
}

// 🔥 FUNGSI BARU: Proses semua gambar di HTML
// Setiap outcome (sukses, gagal, skip) ditrack via `stats` supaya caller
// tahu hasil REAL, bukan asumsi "processed successfully" padahal 0 yang
// benar-benar ke-replace.
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

	imgSelections.Each(func(i int, img *goquery.Selection) {
		// Kalau budget waktu keseluruhan sudah habis, hentikan lebih awal
		// daripada mencoba dan pasti gagal dengan deadline exceeded berulang.
		if ctx.Err() != nil {
			log.Printf("[WARNING] Image processing context done (%v), skipping remaining images", ctx.Err())
			stats.Failed++
			return
		}

		log.Printf("[INFO] Processing image #%d", i+1)

		// Dapatkan src attribute
		src, exists := img.Attr("src")
		if !exists || src == "" {
			log.Printf("[WARNING] Image #%d has no src attribute, skipping", i+1)
			stats.Skipped++
			return
		}

		log.Printf("[INFO] Image #%d src: %s", i+1, src)

		// Download gambar (dengan retry, lebih resilient terhadap kegagalan sementara)
		imageData, contentType, err := s.downloadImageWithRetry(ctx, src, imageDownloadRetries)
		if err != nil {
			log.Printf("[ERROR] Failed to download image #%d: %v", i+1, err)
			stats.Failed++
			return
		}

		log.Printf("[INFO] Image #%d downloaded (size: %d bytes, type: %s)", i+1, len(imageData), contentType)

		// Upload ke MinIO
		minioURL, err := s.uploadToMinio(ctx, imageData, contentType, params)
		if err != nil {
			log.Printf("[ERROR] Failed to upload image #%d to MinIO: %v", i+1, err)
			stats.Failed++
			return
		}

		log.Printf("[INFO] Image #%d uploaded to MinIO: %s", i+1, minioURL)

		// Ganti src attribute dengan URL MinIO
		img.SetAttr("src", minioURL)
		stats.Replaced++
		log.Printf("[INFO] Image #%d src replaced with MinIO URL", i+1)
	})

	log.Printf("[INFO] Total images found: %d, replaced: %d, failed: %d, skipped: %d",
		stats.Total, stats.Replaced, stats.Failed, stats.Skipped)

	// Generate HTML baru
	newHTML, err := doc.Html()
	if err != nil {
		return "", stats, fmt.Errorf("failed to generate new HTML: %w", err)
	}

	return newHTML, stats, nil
}

// 🔥 FUNGSI BARU: Download gambar dari URL
func (s *GenerateServiceImpl) downloadImage(ctx context.Context, imageURL string) ([]byte, string, error) {
	log.Printf("[INFO] Downloading image from: %s", imageURL)

	// Handle base64 image
	if strings.HasPrefix(imageURL, "data:image") {
		log.Println("[INFO] Detected base64 image, decoding...")
		return s.decodeBase64Image(imageURL)
	}

	// 🔥 Sanitasi URL sebelum request. AI (terutama saat menulis URL
	// chart-as-a-service seperti QuickChart.io untuk visualisasi) sering
	// menaruh query string sebagai JS object literal MENTAH tanpa
	// percent-encoding, contoh:
	//   ?c={type:'sankey',data:{datasets:[{from:'Raw Data',to:'X'}]}}
	// Karakter seperti {, }, ', dan spasi mentah membuat server menolak
	// request dengan 400 Bad Request. sanitizeImageURL bersifat idempotent:
	// aman dipanggil baik untuk URL yang sudah ter-encode benar maupun
	// yang belum di-encode sama sekali.
	sanitizedURL, sanitizeErr := sanitizeImageURL(imageURL)
	if sanitizeErr != nil {
		log.Printf("[WARNING] Failed to sanitize image URL, using as-is: %v", sanitizeErr)
		sanitizedURL = imageURL
	} else if sanitizedURL != imageURL {
		log.Printf("[INFO] URL sanitized for request: %s -> %s", imageURL, sanitizedURL)
	}

	// 🔥 BUAT CONTEXT BARU DENGAN TIMEOUT LEBIH LAMA, tetap "anak" dari ctx
	// pemanggil (imgCtx di processImages) supaya tetap menghormati pembatalan
	// budget total, bukan context.Background() yang sepenuhnya lepas.
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
		Timeout: imageDownloadTimeout, // 🔥 60 detik timeout
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Sertakan preview body untuk debugging - kalau provider chart
		// mengembalikan pesan error berguna (misal "invalid JSON config"),
		// ini akan kelihatan di log tanpa perlu reproduce manual.
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

// sanitizeImageURL membersihkan URL gambar/chart yang berasal dari AI, yang
// kerap menulis query string (terutama untuk layanan chart-as-a-service
// seperti QuickChart.io) sebagai JS object literal MENTAH tanpa
// percent-encoding, contoh:
//
//	?c={type:'bar',data:{labels:['Raw Data','Pretraining']}}
//
// Strategi: decode-lalu-encode-ulang membuat proses ini idempotent -
// tidak masalah apakah input sudah ter-encode sebagian, penuh, atau
// tidak sama sekali; hasil akhirnya selalu valid dan konsisten, dan tidak
// akan men-double-encode URL yang sebenarnya sudah benar.
func sanitizeImageURL(rawURL string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "", fmt.Errorf("empty url")
	}

	// data: URL (base64 atau svg mentah) tidak disentuh sama sekali,
	// itu ditangani jalur lain (decodeBase64Image).
	if strings.HasPrefix(rawURL, "data:") {
		return rawURL, nil
	}

	idx := strings.Index(rawURL, "?")
	if idx == -1 {
		// Tidak ada query string sama sekali, tidak ada yang perlu di-encode.
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

	// Pisahkan jadi pasangan key=value berdasarkan "&" di level TOP saja,
	// mengabaikan "&" yang berada di dalam {} atau [] (umum muncul di
	// JSON/JS object literal milik QuickChart), supaya tidak salah memecah
	// value JSON yang kebetulan mengandung "&" literal di dalamnya.
	pairs := splitTopLevelQueryPairs(rawQuery)

	rebuilt := make([]string, 0, len(pairs))
	for _, pair := range pairs {
		eqIdx := strings.Index(pair, "=")
		if eqIdx == -1 {
			// Param tanpa value, biarkan apa adanya.
			rebuilt = append(rebuilt, pair)
			continue
		}

		key := pair[:eqIdx]
		value := pair[eqIdx+1:]

		normalizedValue := normalizeQueryValue(value)
		rebuilt = append(rebuilt, key+"="+url.QueryEscape(normalizedValue))
	}

	sanitized := base + "?" + strings.Join(rebuilt, "&")

	// Validasi akhir: pastikan hasilnya benar-benar parseable.
	if _, err := url.Parse(sanitized); err != nil {
		return "", fmt.Errorf("sanitized url still invalid: %w", err)
	}

	return sanitized, nil
}

// splitTopLevelQueryPairs memisahkan query string jadi pasangan key=value
// berdasarkan "&", TAPI mengabaikan "&" yang berada di dalam tanda kurung
// {} atau [] (umum muncul di JSON/JS object literal milik QuickChart),
// supaya tidak salah memecah value JSON yang mengandung "&" literal.
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

// normalizeQueryValue membuat proses idempotent: kalau value SUDAH
// ter-percent-encode (mengandung urutan %XX valid), decode dulu supaya
// tidak terjadi double-encoding saat di-escape ulang nanti. Kalau value
// memang belum di-encode sama sekali (karakter mentah seperti {, ', spasi),
// url.QueryUnescape akan gagal secara aman dan value mentah dipakai
// langsung sebagai basis untuk encoding berikutnya.
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

// downloadImageWithRetry - dipanggil oleh processImages. Menghormati ctx:
// kalau ctx sudah dibatalkan/timeout, hentikan retry lebih awal daripada
// menunggu backoff yang tidak akan berguna.
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

// uploadToMinio - pakai ctx (anak dari imgCtx di processImages) dengan
// timeout sendiri per-upload, jadi tidak lagi tersandung deadline parent
// yang pendek (penyebab "context deadline exceeded" sebelumnya), tapi
// tetap bisa dibatalkan kalau budget total imageProcessTimeout habis.
func (s *GenerateServiceImpl) uploadToMinio(ctx context.Context, imageData []byte, contentType string, params generate.ArticleGenerationParams) (string, error) {
	log.Println("[INFO] Uploading image to MinIO...")

	// Generate object name
	// Format: articles/{slug}/{timestamp}_{random}.{ext}
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
		return "jpg" // default
	}
}

// 🔥 HELPER: Generate random string
// Pakai *rand.Rand package-level dengan mutex agar thread-safe kalau
// suatu saat processImages diparalelkan.
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

	objectName := fmt.Sprintf("images/%d.%s", time.Now().UnixNano(), ext)
	log.Printf("[DEBUG] Object name: %s", objectName)

	// 🔥 Sama seperti uploadToMinio di processImages: jangan biarkan upload
	// MinIO terikat pada ctx request yang bisa saja sudah hampir expired
	// setelah SendRequest ke AI provider untuk image generation (sampai 120s).
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
