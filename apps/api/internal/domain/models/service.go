// internal/domain/generate/service.go
package generate

import "context"

type Service interface {
	GenerateArticle(ctx context.Context, params ArticleGenerationParams) (*ArticleGenerationResult, error)
	GenerateImage(ctx context.Context, params ImageGenerationParams) (*ImageGenerationResult, error)
}
