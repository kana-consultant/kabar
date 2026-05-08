// internal/application/aimodel/service.go
package aimodel

import (
	"context"

	aimodel "seo-backend/internal/domain/model"
)

type ServiceImpl struct {
	repo aimodel.Repository
}

func NewService(repo aimodel.Repository) aimodel.Service {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) GetAllModels(ctx context.Context) ([]aimodel.AIModel, error) {
	return s.repo.GetAll(ctx)
}

func (s *ServiceImpl) GetAllModelsWithStatus(ctx context.Context, userRole string) ([]aimodel.ModelWithStatus, error) {
	return s.repo.GetAllWithStatus(ctx, userRole)
}

func (s *ServiceImpl) GetModelByID(ctx context.Context, id string) (*aimodel.AIModel, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ServiceImpl) GetModelsByProvider(ctx context.Context, providerID string) ([]aimodel.AIModel, error) {
	return s.repo.GetByProvider(ctx, providerID)
}
