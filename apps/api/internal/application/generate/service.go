package generate

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"seo-backend/internal/domain/generate"
	"seo-backend/internal/helper"
	"seo-backend/internal/infrastructure/ai/builder"
	"seo-backend/internal/infrastructure/ai/parser"
	"seo-backend/internal/infrastructure/http/client"
	"seo-backend/internal/infrastructure/http/minio"
	"time"
)

type GenereteServiceImpl struct {
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
	return &GenereteServiceImpl{
		repo:           repo,
		httpClient:     httpClient,
		promptBuilder:  promptBuilder,
		requestBuilder: requestBuilder,
		responseParser: responseParser,
		minioClient:    minioClient,
	}
}

func (s *GenereteServiceImpl) GenerateArticle(ctx context.Context, params generate.ArticleGenerationParams) (*generate.ArticleResult, error) {
	log.Println("========== GENERATE ARTICLE ==========")
	log.Printf("[INFO] Starting article generation with topic: %s", params.Topic)
	log.Printf("[INFO] Model ID: %s, Tone: %s, Length: %s, Language: %s",
		params.ModelID, params.Tone, params.Length, params.Language)
	log.Printf("[INFO] Auto generate image: %v", params.AutoGenerateImage)

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
	log.Printf("[INFO] Model config loaded - Model: %s, BaseURL: %s, Endpoint: %s",
		config.ModelName, config.BaseURL, config.Endpoint)

	// Build prompt
	log.Println("[INFO] Building article prompt...")

	log.Printf("PROMPT SYSTEM ====== %v", config.SystemPrompt)

	prompt := config.SystemPrompt
	prompt += s.promptBuilder.BuildArticlePrompt(params)

	log.Printf("[INFO] Prompt built successfully (length: %d characters)", len(prompt))

	// Build request body
	log.Println("[INFO] Building request body...")
	requestBody, err := s.requestBuilder.BuildArticleRequestBody(config, prompt)
	if err != nil {
		log.Printf("[ERROR] Failed to build request: %v", err)
		return nil, fmt.Errorf("failed to build request: %w", err)
	}
	log.Printf("[INFO] Request body built successfully (size: %d bytes)", len(requestBody))

	// Send request
	log.Printf("[INFO] Sending request to AI provider (timeout: 90s)...")
	response, err := s.httpClient.SendRequest(ctx, config, requestBody, 90*time.Second)
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

	// Hitung SEO score otomatis
	result.SeoScore = helper.CalculateSEOScoreSimple(result.Title, result.Content, result.Excerpt, params.Topic)

	log.Printf("SUCCESS title=%s words=%d seo_score=%d slug=%s", result.Title, result.WordCount, result.SeoScore, result.Slug)
	return result, nil
}

func (s *GenereteServiceImpl) GenerateImage(ctx context.Context, params generate.ImageGenerationParams) (*generate.ImageResult, error) {
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
	log.Printf("[DEBUG] Model config loaded: %s", config.ModelName)

	SystemPrompt := config.SystemPrompt + "Tentang" + params.Prompt

	// Build request body
	requestBody, err := s.requestBuilder.BuildImageRequestBody(config, SystemPrompt)
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
	log.Printf("[DEBUG] AI provider raw response: %s", string(response))

	// Parse response - dapatkan Base64 string dari AI provider
	base64String, contentType, err := s.responseParser.ParseImageResponseBase64(response, config.ResponseImagePath)
	if err != nil {
		log.Printf("[ERROR] Failed to parse image response: %v", err)
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	log.Printf("[DEBUG] Image parsed, contentType: %s, base64 length: %d", contentType, len(base64String))

	// Remove data URL prefix jika ada
	if strings.Contains(base64String, "base64,") {
		parts := strings.Split(base64String, "base64,")
		base64String = parts[len(parts)-1]
	}
	base64String = strings.TrimPrefix(base64String, "data:image/png;base64,")
	base64String = strings.TrimPrefix(base64String, "data:image/jpeg;base64,")
	base64String = strings.TrimPrefix(base64String, "data:image/webp;base64,")

	// Decode Base64 ke binary
	imageData, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		log.Printf("[ERROR] Failed to decode base64: %v", err)
		return nil, fmt.Errorf("failed to decode base64 image: %w", err)
	}
	log.Printf("[DEBUG] Base64 decoded, image size: %d bytes", len(imageData))

	// Tentukan ekstensi berdasarkan content-type
	ext := "png"
	if contentType == "image/jpeg" || contentType == "image/jpg" {
		ext = "jpg"
	} else if contentType == "image/webp" {
		ext = "webp"
	} else if contentType == "image/gif" {
		ext = "gif"
	}
	log.Printf("[DEBUG] Image extension: %s", ext)

	objectName := fmt.Sprintf("images/%d.%s", time.Now().UnixNano(), ext)
	log.Printf("[DEBUG] Object name: %s", objectName)

	// Upload ke MinIO
	reader := bytes.NewReader(imageData)
	uploadedName, err := s.minioClient.Upload(ctx, objectName, reader, int64(len(imageData)), contentType)
	if err != nil {
		log.Printf("[ERROR] Failed to upload to Minio: %v", err)
		return nil, fmt.Errorf("failed to upload image to Minio: %w", err)
	}
	log.Printf("[DEBUG] Uploaded to Minio: %s", uploadedName)

	// Ambil presigned URL dengan expiry 7 hari
	imageURL, err := s.minioClient.GetURL(ctx, uploadedName, 7*24*time.Hour)
	if err != nil {
		log.Printf("[ERROR] Failed to get Minio URL: %v", err)
		return nil, fmt.Errorf("failed to get Minio URL: %w", err)
	}
	log.Printf("[DEBUG] Presigned URL generated: %s", imageURL)

	result := &generate.ImageResult{
		ImageURL:    imageURL,
		Prompt:      params.Prompt,
		GeneratedAt: helper.ParseWIBTime(time.Now().Format(time.RFC3339)).Format(time.RFC3339),
		Model:       config.ModelName,
	}

	log.Printf("[DEBUG] GenerateImage completed successfully, imageURL: %s", imageURL)
	return result, nil
}

// Private helper methods
func (s *GenereteServiceImpl) validateArticleParams(params generate.ArticleGenerationParams) error {
	if params.Topic == "" {
		return fmt.Errorf("topic is required")
	}
	if params.ModelID == "" {
		return fmt.Errorf("modelId is required")
	}
	return nil
}

func (s *GenereteServiceImpl) validateImageParams(params generate.ImageGenerationParams) error {
	if params.Prompt == "" {
		return fmt.Errorf("prompt is required")
	}
	if params.ModelID == "" {
		return fmt.Errorf("modelId is required")
	}
	return nil
}

func (s *GenereteServiceImpl) saveHistory(ctx context.Context, params generate.ArticleGenerationParams, result *generate.ArticleResult) error {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	history := &generate.GenerationHistory{
		Type:      "article",
		Topic:     params.Topic,
		Result:    result.Content,
		ModelID:   params.ModelID,
		CreatedAt: time.Now().In(loc),
	}
	return s.repo.SaveHistory(ctx, history)
}
