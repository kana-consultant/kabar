// internal/domain/dashboard/service.go
package dashboard

import (
	"context"
)

type DashboardService interface {
	GetStats(ctx context.Context, userCtx DashboardFilter) (DashboardStats, error)
}
