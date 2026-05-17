package generate

import (
	"context"
	"fmt"
	"log"

	"seo-backend/internal/domain/generate"
	"seo-backend/internal/helper"
	"seo-backend/internal/infrastructure/ai/builder"
	"seo-backend/internal/infrastructure/ai/parser"
	"seo-backend/internal/infrastructure/http/client"
	"time"
)

type GenereteServiceImpl struct {
	repo           generate.Repository
	httpClient     *client.HTTPClient
	promptBuilder  *builder.PromptBuilder
	requestBuilder *builder.RequestBuilder
	responseParser *parser.ResponseParser
}

func NewService(
	repo generate.Repository,
	httpClient *client.HTTPClient,
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
	prompt := s.promptBuilder.BuildArticlePrompt(params)
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

	log.Printf("SUCCESS title=%s words=%d seo_score=%d", result.Title, result.WordCount, result.SeoScore)
	return result, nil

	log.Printf("[SUCCESS] Article generated - Title: %s, Word count: %d, SEO Score: %d",
		result.Title, result.WordCount, result.SeoScore)
	return result, nil
}

func (s *GenereteServiceImpl) GenerateImage(ctx context.Context, params generate.ImageGenerationParams) (*generate.ImageResult, error) {
	log.Println("========== GENERATE IMAGE ==========")
	defer log.Println("========== END GENERATE IMAGE ==========")

	// Validate params
	if err := s.validateImageParams(params); err != nil {
		return nil, err
	}

	// Get model configuration
	config, err := s.repo.GetModelConfig(ctx, params.ModelID, "image")
	if err != nil {
		return nil, fmt.Errorf("failed to get model config: %w", err)
	}

	// Build request body
	requestBody, err := s.requestBuilder.BuildImageRequestBody(config, params.Prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	// Send request with longer timeout for images
	response, err := s.httpClient.SendRequest(ctx, config, requestBody, 120*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Parse response
	imageURL, err := s.responseParser.ParseImageResponse(response, config.ResponseImagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	result := &generate.ImageResult{
		ImageURL:    imageURL,
		Prompt:      params.Prompt,
		GeneratedAt: helper.ParseWIBTime(time.Now().Format(time.RFC3339)).Format(time.RFC3339),
		Model:       config.ModelName,
	}

	log.Printf("SUCCESS: Image generated for prompt: %s", params.Prompt)
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
