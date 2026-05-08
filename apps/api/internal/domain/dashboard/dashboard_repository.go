package dashboard

import "context"

type DashboardRepository interface {
	GetTotalContent(ctx context.Context, where string, args []interface{}) (int, error)
	GetTotalProducts(ctx context.Context, where string, args []interface{}) (int, error)
	GetTotalPublished(ctx context.Context, where string, args []interface{}) (int, error)
	GetAverageSeoScore(ctx context.Context, where string, args []interface{}) (float64, error)

	GetContentCountByPeriod(ctx context.Context, where string, args []interface{}) (int, error)
	GetProductsCountByPeriod(ctx context.Context, where string, args []interface{}) (int, error)
	GetSeoScoreByPeriod(ctx context.Context, where string, args []interface{}) (float64, error)
}

type DashboardFilter struct {
	UserID string
	TeamID string
	Role   string
}
