// internal/domain/aimodel/service.go
package aimodel

import "context"

type Service interface {
	GetAllModels(ctx context.Context) ([]AIModel, error)
	GetAllModelsWithStatus(ctx context.Context, userRole string) ([]ModelWithStatus, error)
	GetModelByID(ctx context.Context, id string) (*AIModel, error)
	GetModelsByProvider(ctx context.Context, providerID string) ([]AIModel, error)
}
