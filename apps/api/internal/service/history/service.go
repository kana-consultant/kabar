package history

import (
	"fmt"

	"seo-backend/internal/database"
	"seo-backend/internal/models"
)

type Service struct {
	repo         *Repository
	auth         *Authorizer
	queryBuilder *QueryBuilder
	statistics   *Statistics
}

func NewService() *Service {
	db := database.GetDB()
	repo := NewRepository(db)
	qb := NewQueryBuilder()

	return &Service{
		repo:         repo,
		auth:         NewAuthorizer(),
		queryBuilder: qb,
		statistics:   NewStatistics(repo, qb),
	}
}

// Public methods
func (s *Service) GetAll(ctx models.UserContext, filter HistoryFilter) ([]models.History, error) {
	query, args := s.queryBuilder.BuildListQuery(ctx, filter)
	return s.repo.GetAll(query, args)
}

func (s *Service) GetByID(id string, ctx models.UserContext) (*models.History, error) {
	history, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if history == nil {
		return nil, nil
	}

	if !s.auth.CanAccess(history, ctx) {
		return nil, fmt.Errorf("access denied")
	}

	return history, nil
}

func (s *Service) Create(req CreateHistoryRequest) (string, error) {
	return s.repo.Create(req)
}

func (s *Service) Delete(id string, ctx models.UserContext) error {
	history, err := s.GetByID(id, ctx)
	if err != nil {
		return err
	}

	if history == nil {
		return fmt.Errorf("history not found")
	}

	return s.repo.Delete(id)
}

func (s *Service) ClearAll(ctx models.UserContext) error {
	if err := s.auth.ValidateAdminAccess(ctx); err != nil {
		return err
	}
	return s.repo.ClearAll()
}

func (s *Service) GetStats(ctx models.UserContext) (*HistoryStats, error) {
	return s.statistics.GetStats(ctx)
}
