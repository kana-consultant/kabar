package generate

import (
	"log"
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

// GenerateArticle generates an article based on parameters
func (s *Service) GenerateArticle(params ArticleGenerationParams) (*ArticleResult, error) {
	log.Println("========== GENERATE ARTICLE ==========")
	defer log.Println("========== END GENERATE ARTICLE ==========")

	// Validate params
	if err := validateArticleParams(params); err != nil {
		return nil, err
	}

	// Get model configuration
	config, err := s.repo.GetModelConfig(params.ModelID, "text")
	if err != nil {
		return nil, err
	}

	// Build prompt
	prompt := buildArticlePrompt(params)

	// Build request body
	requestBody, err := buildArticleRequestBody(config, prompt)
	if err != nil {
		return nil, err
	}

	// Send request
	response, err := s.httpClient.SendRequest(config, requestBody, 90*time.Second)
	if err != nil {
		return nil, err
	}

	// Parse response
	result, err := parseArticleResponse(response, config.ResponsePath)
	if err != nil {
		return nil, err
	}

	log.Printf("SUCCESS title=%s words=%d", result.Title, result.WordCount)
	return result, nil
}

// GenerateImage generates an image based on prompt
func (s *Service) GenerateImage(params ImageGenerationParams) (*ImageResult, error) {
	log.Println("========== GENERATE IMAGE ==========")
	defer log.Println("========== END GENERATE IMAGE ==========")

	// Validate params
	if err := validateImageParams(params); err != nil {
		return nil, err
	}

	// Get model configuration
	config, err := s.repo.GetModelConfig(params.ModelID, "image")
	if err != nil {
		return nil, err
	}

	// Build request body
	requestBody, err := buildImageRequestBody(config, params.Prompt)
	if err != nil {
		return nil, err
	}

	// Send request with longer timeout for images
	response, err := s.httpClient.SendRequest(config, requestBody, 120*time.Second)
	if err != nil {
		return nil, err
	}

	// Parse response
	imageURL, err := parseImageResponse(response, config.ResponsePath)
	if err != nil {
		return nil, err
	}

	result := &ImageResult{
		ImageURL:    imageURL,
		Prompt:      params.Prompt,
		GeneratedAt: time.Now().Format(time.RFC3339),
		Model:       config.ModelName,
	}

	log.Printf("SUCCESS: Image generated for prompt: %s", params.Prompt)
	return result, nil
}
