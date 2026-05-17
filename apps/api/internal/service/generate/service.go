package generate

import (
	"log"
	"seo-backend/internal/helper"
	"time"
)

type Service struct {
	httpClient *HTTPClient
	repo       *Repository
}

func NewService() *Service {
	return &Service{
		httpClient: NewHTTPClient(90 * time.Second),
		repo:       NewRepository(),
	}
}

func (s *Service) GenerateArticle(params ArticleGenerationParams) (*ArticleResult, error) {
	log.Println("========== GENERATE ARTICLE ==========")
	defer log.Println("========== END GENERATE ARTICLE ==========")

	if err := validateArticleParams(params); err != nil {
		return nil, err
	}

	config, err := s.repo.GetModelConfig(params.ModelID, "text")
	if err != nil {
		return nil, err
	}

	prompt := buildArticlePrompt(params)

	requestBody, err := buildArticleRequestBody(config, prompt)
	if err != nil {
		return nil, err
	}

	response, err := s.httpClient.SendRequest(config, requestBody, 90*time.Second)
	if err != nil {
		return nil, err
	}

	result, err := parseArticleResponse(response, config.ResponsePath)
	if err != nil {
		return nil, err
	}

	// Hitung SEO score otomatis
	result.SeoScore = calculateSEOScoreSimple(result.Title, result.Content, result.Excerpt, params.Topic)

	log.Printf("SUCCESS title=%s words=%d seo_score=%d", result.Title, result.WordCount, result.SeoScore)
	return result, nil
}

func (s *Service) GenerateImage(params ImageGenerationParams) (*ImageResult, error) {
	log.Println("========== GENERATE IMAGE ==========")
	defer log.Println("========== END GENERATE IMAGE ==========")

	if err := validateImageParams(params); err != nil {
		return nil, err
	}

	config, err := s.repo.GetModelConfig(params.ModelID, "image")
	if err != nil {
		return nil, err
	}

	requestBody, err := buildImageRequestBody(config, params.Prompt)
	if err != nil {
		return nil, err
	}

	response, err := s.httpClient.SendRequest(config, requestBody, 120*time.Second)
	if err != nil {
		return nil, err
	}

	imageURL, err := parseImageResponse(response, config.ResponseImagePath)
	if err != nil {
		return nil, err
	}

	result := &ImageResult{
		ImageURL:    imageURL,
		Prompt:      params.Prompt,
		GeneratedAt: helper.ParseWIBTime(time.Now().Format(time.RFC3339)).Format(time.RFC3339),
		Model:       config.ModelName,
	}

	log.Printf("SUCCESS: Image generated for prompt: %s", params.Prompt)
	return result, nil
}