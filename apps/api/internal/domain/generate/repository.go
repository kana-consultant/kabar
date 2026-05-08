// internal/domain/generate/repository.go
package generate

import "context"

type Repository interface {
	GetModelConfig(ctx context.Context, modelID, serviceType string) (*ModelConfig, error)
	SaveHistory(ctx context.Context, history *GenerationHistory) error
}
