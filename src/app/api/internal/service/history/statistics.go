package history

import "seo-backend/internal/models"

type Statistics struct {
	repo         *Repository
	queryBuilder *QueryBuilder
}

func NewStatistics(repo *Repository, qb *QueryBuilder) *Statistics {
	return &Statistics{
		repo:         repo,
		queryBuilder: qb,
	}
}

func (s *Statistics) GetStats(ctx models.UserContext) (*HistoryStats, error) {
	query, args := s.queryBuilder.BuildStatsQuery(ctx)

	stats, err := s.repo.GetStats(query, args)
	if err != nil {
		return nil, err
	}

	// Calculate success rate
	if stats.Total > 0 {
		stats.SuccessRate = float64(stats.SuccessCount) / float64(stats.Total) * 100
	}

	return stats, nil
}
