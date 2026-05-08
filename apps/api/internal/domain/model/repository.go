// internal/domain/aimodel/repository.go
package aimodel

import "context"

type Repository interface {
	GetAll(ctx context.Context) ([]AIModel, error)
	GetAllWithStatus(ctx context.Context, userRole string) ([]ModelWithStatus, error)
	GetByID(ctx context.Context, id string) (*AIModel, error)
	GetByProvider(ctx context.Context, providerID string) ([]AIModel, error)
}
