// internal/domain/generate/client.go (Optional)
package generate

import "context"

type Client interface {
	GenerateArticleWithDefaults(ctx context.Context, topic, modelID string) (*ArticleResult, error)
	GenerateImageWithPrompt(ctx context.Context, prompt, modelID string) (*ImageResult, error)
}
