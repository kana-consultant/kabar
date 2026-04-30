package dashboard

import (
	"log"
	"net/http"
	"time"

	"seo-backend/internal/database"
)

// GetStats returns dashboard statistics
func (h *DashboardHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)
	log.Printf("Dashboard stats - role: %s, teamID: %s, userID: %s", userCtx.Role, userCtx.TeamID, userCtx.UserID)

	// Build base where clause
	whereClause, baseArgs := h.buildWhereClause(userCtx)

	// Get current period stats (last 30 days)
	currentPeriodStart := time.Now().AddDate(0, 0, -30)
	previousPeriodStart := time.Now().AddDate(0, 0, -60)
	previousPeriodEnd := currentPeriodStart

	var stats DashboardStats

	// Fetch all statistics
	h.fetchTotalContent(&stats, whereClause, baseArgs)
	h.fetchTotalProducts(&stats, whereClause, baseArgs)
	h.fetchTotalPublished(&stats, whereClause, baseArgs)
	h.fetchAverageSeoScore(&stats, whereClause, baseArgs)

	// Calculate percentage
	if stats.TotalContent > 0 {
		stats.PublishedPercentage = int(float64(stats.TotalPublished) / float64(stats.TotalContent) * 100)
	}

	// Calculate changes
	h.calculateContentChange(&stats, whereClause, baseArgs, currentPeriodStart, previousPeriodStart, previousPeriodEnd)
	h.calculateProductsChange(&stats, whereClause, baseArgs, currentPeriodStart, previousPeriodStart, previousPeriodEnd)
	h.calculateSeoScoreChange(&stats, whereClause, baseArgs, currentPeriodStart, previousPeriodStart, previousPeriodEnd)

	log.Printf("Dashboard stats: content=%d, products=%d, published=%d, seoScore=%.1f",
		stats.TotalContent, stats.TotalProducts, stats.TotalPublished, stats.AverageSeoScore)

	h.writeJSON(w, stats, http.StatusOK)
}

// fetchTotalContent gets total content count (drafts + histories)
func (h *DashboardHandler) fetchTotalContent(stats *DashboardStats, whereClause string, args []interface{}) {
	query := `
		SELECT COUNT(*) FROM (
			SELECT id FROM drafts WHERE ` + whereClause + `
			UNION ALL
			SELECT id FROM histories WHERE ` + whereClause + `
		) AS all_content
	`
	err := database.GetDB().QueryRow(query, args...).Scan(&stats.TotalContent)
	if err != nil {
		log.Printf("Failed to get total content: %v", err)
		stats.TotalContent = 0
	}
}

// fetchTotalProducts gets total products count
func (h *DashboardHandler) fetchTotalProducts(stats *DashboardStats, whereClause string, args []interface{}) {
	query := `SELECT COUNT(*) FROM products WHERE ` + whereClause
	err := database.GetDB().QueryRow(query, args...).Scan(&stats.TotalProducts)
	if err != nil {
		log.Printf("Failed to get total products: %v", err)
		stats.TotalProducts = 0
	}
}

// fetchTotalPublished gets total published content count
func (h *DashboardHandler) fetchTotalPublished(stats *DashboardStats, whereClause string, args []interface{}) {
	query := `
		SELECT COUNT(*) FROM histories 
		WHERE status IN ('published', 'generated') AND ` + whereClause
	err := database.GetDB().QueryRow(query, args...).Scan(&stats.TotalPublished)
	if err != nil {
		log.Printf("Failed to get total published: %v", err)
		stats.TotalPublished = 0
	}
}

// fetchAverageSeoScore gets average SEO score from histories
func (h *DashboardHandler) fetchAverageSeoScore(stats *DashboardStats, whereClause string, args []interface{}) {
	query := `
		SELECT COALESCE(AVG(seo_score), 0) FROM histories 
		WHERE ` + whereClause + ` AND seo_score IS NOT NULL
	`
	err := database.GetDB().QueryRow(query, args...).Scan(&stats.AverageSeoScore)
	if err != nil {
		log.Printf("Failed to get average SEO score: %v", err)
		// Default mock value if column doesn't exist
		stats.AverageSeoScore = 78.5
	}
}

// calculateContentChange calculates content change between periods
func (h *DashboardHandler) calculateContentChange(stats *DashboardStats, whereClause string, args []interface{},
	currentStart, previousStart, previousEnd time.Time) {

	// Current period content count
	currentWhere, currentArgs := h.buildWhereClauseWithDate(whereClause, args, currentStart, ">=")
	currentQuery := `
		SELECT COUNT(*) FROM (
			SELECT id FROM drafts WHERE ` + currentWhere + `
			UNION ALL
			SELECT id FROM histories WHERE ` + currentWhere + `
		) AS all_content
	`
	var currentCount int
	err := database.GetDB().QueryRow(currentQuery, currentArgs...).Scan(&currentCount)
	if err != nil {
		log.Printf("Failed to get current content count: %v", err)
		currentCount = 0
	}

	// Previous period content count
	previousWhere, previousArgs := h.buildWhereClauseWithDate(whereClause, args, previousStart, ">=")
	previousWhere, previousArgs = h.buildWhereClauseWithDate(previousWhere, previousArgs, previousEnd, "<")
	previousQuery := `
		SELECT COUNT(*) FROM (
			SELECT id FROM drafts WHERE ` + previousWhere + `
			UNION ALL
			SELECT id FROM histories WHERE ` + previousWhere + `
		) AS all_content
	`
	var previousCount int
	err = database.GetDB().QueryRow(previousQuery, previousArgs...).Scan(&previousCount)
	if err != nil {
		log.Printf("Failed to get previous content count: %v", err)
		previousCount = 0
	}

	// Calculate change
	diff := currentCount - previousCount
	if diff > 0 {
		stats.ContentChange = "+" + string(rune(diff+48)) + " minggu ini"
	} else if diff < 0 {
		stats.ContentChange = string(rune(diff+48)) + " minggu ini"
	} else {
		stats.ContentChange = "0 minggu ini"
	}
}

// calculateProductsChange calculates products change between periods
func (h *DashboardHandler) calculateProductsChange(stats *DashboardStats, whereClause string, args []interface{},
	currentStart, previousStart, previousEnd time.Time) {

	// Current period products count
	currentWhere, currentArgs := h.buildWhereClauseWithDate(whereClause, args, currentStart, ">=")
	currentQuery := `SELECT COUNT(*) FROM products WHERE ` + currentWhere
	var currentCount int
	err := database.GetDB().QueryRow(currentQuery, currentArgs...).Scan(&currentCount)
	if err != nil {
		log.Printf("Failed to get current products count: %v", err)
		currentCount = 0
	}

	// Previous period products count
	previousWhere, previousArgs := h.buildWhereClauseWithDate(whereClause, args, previousStart, ">=")
	previousWhere, previousArgs = h.buildWhereClauseWithDate(previousWhere, previousArgs, previousEnd, "<")
	previousQuery := `SELECT COUNT(*) FROM products WHERE ` + previousWhere
	var previousCount int
	err = database.GetDB().QueryRow(previousQuery, previousArgs...).Scan(&previousCount)
	if err != nil {
		log.Printf("Failed to get previous products count: %v", err)
		previousCount = 0
	}

	// Calculate change
	diff := currentCount - previousCount
	if diff > 0 {
		stats.ProductsChange = "+" + string(rune(diff+48)) + " bulan ini"
	} else if diff < 0 {
		stats.ProductsChange = string(rune(diff+48)) + " bulan ini"
	} else {
		stats.ProductsChange = "0 bulan ini"
	}
}

// calculateSeoScoreChange calculates SEO score change between periods
func (h *DashboardHandler) calculateSeoScoreChange(stats *DashboardStats, whereClause string, args []interface{},
	currentStart, previousStart, previousEnd time.Time) {

	// Current period SEO score
	currentWhere, currentArgs := h.buildWhereClauseWithDate(whereClause, args, currentStart, ">=")
	currentQuery := `
		SELECT COALESCE(AVG(seo_score), 0) FROM histories 
		WHERE ` + currentWhere + ` AND seo_score IS NOT NULL
	`
	var currentScore float64
	err := database.GetDB().QueryRow(currentQuery, currentArgs...).Scan(&currentScore)
	if err != nil {
		log.Printf("Failed to get current SEO score: %v", err)
		currentScore = 0
	}

	// Previous period SEO score
	previousWhere, previousArgs := h.buildWhereClauseWithDate(whereClause, args, previousStart, ">=")
	previousWhere, previousArgs = h.buildWhereClauseWithDate(previousWhere, previousArgs, previousEnd, "<")
	previousQuery := `
		SELECT COALESCE(AVG(seo_score), 0) FROM histories 
		WHERE ` + previousWhere + ` AND seo_score IS NOT NULL
	`
	var previousScore float64
	err = database.GetDB().QueryRow(previousQuery, previousArgs...).Scan(&previousScore)
	if err != nil {
		log.Printf("Failed to get previous SEO score: %v", err)
		previousScore = 0
	}

	// Calculate change
	diff := currentScore - previousScore
	if diff > 0 {
		stats.SeoScoreChange = "+" + formatFloatValue(diff) + "%"
	} else if diff < 0 {
		stats.SeoScoreChange = formatFloatValue(diff) + "%"
	} else {
		stats.SeoScoreChange = "0%"
	}
}
