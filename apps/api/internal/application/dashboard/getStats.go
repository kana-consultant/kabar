package dashboard

import (
	"context"
	"log"
	"seo-backend/internal/domain/dashboard"
	"strconv"
	"time"
)

type DashboardService struct {
	repo dashboard.DashboardRepository
}

func NewDashboardService(repo dashboard.DashboardRepository) dashboard.DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetStats(ctx context.Context, userCtx dashboard.DashboardFilter) (dashboard.DashboardStats, error) {

	// sementara masih pakai whereClause (nanti bisa di-refactor ke filter struct)
	whereClause, baseArgs := s.buildWhereClause(userCtx)

	currentStart := time.Now().AddDate(0, 0, -30)
	previousStart := time.Now().AddDate(0, 0, -60)
	previousEnd := currentStart

	stats := dashboard.DashboardStats{}

	// === FETCH DATA ===
	totalContent, _ := s.repo.GetTotalContent(ctx, whereClause, baseArgs)
	totalProducts, _ := s.repo.GetTotalProducts(ctx, whereClause, baseArgs)
	totalPublished, _ := s.repo.GetTotalPublished(ctx, whereClause, baseArgs)
	avgSeo, _ := s.repo.GetAverageSeoScore(ctx, whereClause, baseArgs)

	stats.TotalContent = totalContent
	stats.TotalProducts = totalProducts
	stats.TotalPublished = totalPublished
	stats.AverageSeoScore = avgSeo

	// === BUSINESS LOGIC ===
	if totalContent > 0 {
		stats.PublishedPercentage = int(float64(totalPublished) / float64(totalContent) * 100)
	}

	// === CALCULATE CHANGE ===
	stats.ContentChange = s.calculateContentChange(ctx, whereClause, baseArgs, currentStart, previousStart, previousEnd)
	stats.ProductsChange = s.calculateProductsChange(ctx, whereClause, baseArgs, currentStart, previousStart, previousEnd)
	stats.SeoScoreChange = s.calculateSeoScoreChange(ctx, whereClause, baseArgs, currentStart, previousStart, previousEnd)

	return stats, nil
}

func (s *DashboardService) calculateContentChange(
	ctx context.Context,
	where string,
	args []interface{},
	currentStart, previousStart, previousEnd time.Time,
) string {

	currentWhere, currentArgs := s.buildWhereClauseWithDate(where, args, currentStart, ">=")
	current, _ := s.repo.GetContentCountByPeriod(ctx, currentWhere, currentArgs)

	prevWhere, prevArgs := s.buildWhereClauseWithDate(where, args, previousStart, ">=")
	prevWhere, prevArgs = s.buildWhereClauseWithDate(prevWhere, prevArgs, previousEnd, "<")
	previous, _ := s.repo.GetContentCountByPeriod(ctx, prevWhere, prevArgs)

	diff := current - previous

	if diff > 0 {
		return "+" + strconv.Itoa(diff) + " minggu ini"
	} else if diff < 0 {
		return strconv.Itoa(diff) + " minggu ini"
	}
	return "0 minggu ini"
}

func (s *DashboardService) calculateProductsChange(
	ctx context.Context,
	where string,
	args []interface{},
	currentStart, previousStart, previousEnd time.Time,
) string {

	currentWhere, currentArgs := s.buildWhereClauseWithDate(where, args, currentStart, ">=")
	current, _ := s.repo.GetProductsCountByPeriod(ctx, currentWhere, currentArgs)

	prevWhere, prevArgs := s.buildWhereClauseWithDate(where, args, previousStart, ">=")
	prevWhere, prevArgs = s.buildWhereClauseWithDate(prevWhere, prevArgs, previousEnd, "<")
	previous, _ := s.repo.GetProductsCountByPeriod(ctx, prevWhere, prevArgs)

	diff := current - previous

	if diff > 0 {
		return "+" + strconv.Itoa(diff) + " bulan ini"
	} else if diff < 0 {
		return strconv.Itoa(diff) + " bulan ini"
	}
	return "0 bulan ini"
}

func (s *DashboardService) calculateSeoScoreChange(
	ctx context.Context,
	where string,
	args []interface{},
	currentStart, previousStart, previousEnd time.Time,
) string {

	currentWhere, currentArgs := s.buildWhereClauseWithDate(where, args, currentStart, ">=")
	current, _ := s.repo.GetSeoScoreByPeriod(ctx, currentWhere, currentArgs)

	prevWhere, prevArgs := s.buildWhereClauseWithDate(where, args, previousStart, ">=")
	prevWhere, prevArgs = s.buildWhereClauseWithDate(prevWhere, prevArgs, previousEnd, "<")
	previous, _ := s.repo.GetSeoScoreByPeriod(ctx, prevWhere, prevArgs)

	diff := current - previous

	if diff > 0 {
		return "+" + s.formatFloatValue(diff) + "%"
	} else if diff < 0 {
		return s.formatFloatValue(diff) + "%"
	}
	return "0%"
}

func (s *DashboardService) buildWhereClauseWithDate(
	whereClause string,
	args []interface{},
	date time.Time,
	operator string,
) (string, []interface{}) {

	paramIndex := len(args) + 1

	newWhere := whereClause + " AND created_at " + operator + " $" + strconv.Itoa(paramIndex)
	newArgs := append(args, date)

	return newWhere, newArgs
}
func (s *DashboardService) buildWhereClause(ctx dashboard.DashboardFilter) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	switch ctx.Role {
	case "super_admin":
		// Super admin: melihat semua data
		conditions = append(conditions, "1=1")
		log.Println("Super admin - melihat semua stats")

	case "admin":
		// Admin: hanya melihat data dalam team yang sama
		if ctx.TeamID != "" && ctx.TeamID != "00000000-0000-0000-0000-000000000000" {
			conditions = append(conditions, "team_id = $1")
			args = append(args, ctx.TeamID)
			argIndex++
		} else {
			conditions = append(conditions, "created_by = $1")
			args = append(args, ctx.UserID)
			argIndex++
		}
		log.Printf("Admin - melihat stats untuk team: %s", ctx.TeamID)

	default:
		// Role lain: hanya melihat data sendiri
		conditions = append(conditions, "created_by = $1")
		args = append(args, ctx.UserID)
		log.Printf("User %s - melihat stats sendiri", ctx.UserID)
	}

	return conditions[0], args
}

func (*DashboardService) formatFloatValue(value float64) string {
	if value == 0 {
		return "0"
	}

	isNegative := value < 0
	if isNegative {
		value = -value
	}

	intPart := int(value)
	decimalPart := int((value - float64(intPart)) * 10)

	var result string
	if decimalPart > 0 {
		result = string(rune(intPart+48)) + "." + string(rune(decimalPart+48))
	} else {
		result = string(rune(intPart + 48))
	}

	if isNegative {
		result = "-" + result
	}

	return result
}
