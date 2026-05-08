// internal/application/generate/client.go
package generate

import (
	"context"

	"seo-backend/internal/domain/generate"
)

type ClientImpl struct {
	service generate.Service
}

func NewClient(service generate.Service) generate.Client {
	return &ClientImpl{
		service: service,
	}
}

func (c *ClientImpl) GenerateArticleWithDefaults(ctx context.Context, topic, modelID string) (*generate.ArticleResult, error) {
	return c.service.GenerateArticle(ctx, generate.ArticleGenerationParams{
		Topic:    topic,
		ModelID:  modelID,
		Tone:     "professional",
		Length:   "medium",
		Language: "English",
	})
}

func (c *ClientImpl) GenerateImageWithPrompt(ctx context.Context, prompt, modelID string) (*generate.ImageResult, error) {
	return c.service.GenerateImage(ctx, generate.ImageGenerationParams{
		Prompt:  prompt,
		ModelID: modelID,
	})
}
