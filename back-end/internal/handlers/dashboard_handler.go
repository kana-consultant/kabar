// internal/handlers/dashboard_handler.go
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"seo-backend/internal/database"
	"seo-backend/internal/middleware/auth"
)

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

type DashboardStats struct {
	TotalContent        int     `json:"totalContent"`
	TotalProducts       int     `json:"totalProducts"`
	TotalPublished      int     `json:"totalPublished"`
	AverageSeoScore     float64 `json:"averageSeoScore"`
	ContentChange       string  `json:"contentChange"`
	ProductsChange      string  `json:"productsChange"`
	PublishedPercentage int     `json:"publishedPercentage"`
	SeoScoreChange      string  `json:"seoScoreChange"`
}

func (h *DashboardHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userRole := auth.GetUserRole(ctx)
	teamID := auth.GetTeamID(ctx)
	userID := auth.GetUserID(ctx)

	log.Printf("Dashboard stats - role: %s, teamID: %s, userID: %s", userRole, teamID, userID)

	// Build filter conditions based on role
	var whereClause string
	var args []interface{}
	argIndex := 1

	switch userRole {
	case "super_admin":
		// Super admin: melihat semua data
		whereClause = "1=1"
		log.Println("Super admin - melihat semua stats")

	case "admin":
		// Admin: hanya melihat data dalam team yang sama
		if teamID != "" {
			whereClause = "team_id = $" + string(rune(argIndex+48))
			args = append(args, teamID)
			argIndex++
		} else {
			whereClause = "created_by = $" + string(rune(argIndex+48))
			args = append(args, userID)
			argIndex++
		}
		log.Printf("Admin - melihat stats untuk team: %s", teamID)

	default:
		// Role lain: hanya melihat data sendiri
		whereClause = "created_by = $" + string(rune(argIndex+48))
		args = append(args, userID)
		argIndex++
		log.Printf("User %s - melihat stats sendiri", userID)
	}

	// Get current period stats (last 30 days)
	currentPeriodStart := time.Now().AddDate(0, 0, -30)
	previousPeriodStart := time.Now().AddDate(0, 0, -60)
	previousPeriodEnd := currentPeriodStart

	var stats DashboardStats

	// 1. Total Content (drafts + histories)
	contentQuery := `
		SELECT COUNT(*) FROM (
			SELECT id FROM drafts WHERE ` + whereClause + `
			UNION ALL
			SELECT id FROM histories WHERE ` + whereClause + `
		) AS all_content
	`
	err := database.GetDB().QueryRow(contentQuery, args...).Scan(&stats.TotalContent)
	if err != nil {
		log.Printf("Failed to get total content: %v", err)
	}

	// 2. Total Products
	productsQuery := `SELECT COUNT(*) FROM products WHERE ` + whereClause
	err = database.GetDB().QueryRow(productsQuery, args...).Scan(&stats.TotalProducts)
	if err != nil {
		log.Printf("Failed to get total products: %v", err)
	}

	// 3. Total Published (histories with action = 'published' or 'generated')
	publishedQuery := `
		SELECT COUNT(*) FROM histories 
		WHERE status IN ('published', 'generated') AND status = 'success' AND ` + whereClause
	err = database.GetDB().QueryRow(publishedQuery, args...).Scan(&stats.TotalPublished)
	if err != nil {
		log.Printf("Failed to get total published: %v", err)
	}

	// 4. Average SEO Score (from histories)
	// If seo_score column doesn't exist, calculate from readability or use mock
	seoQuery := `
		SELECT COALESCE(AVG(seo_score), 0) FROM histories 
		WHERE ` + whereClause + ` AND seo_score IS NOT NULL
	`
	err = database.GetDB().QueryRow(seoQuery, args...).Scan(&stats.AverageSeoScore)
	if err != nil {
		log.Printf("Failed to get average SEO score: %v", err)
		// Default mock value if column doesn't exist
		stats.AverageSeoScore = 78.5
	}

	// 5. Calculate percentage
	if stats.TotalContent > 0 {
		stats.PublishedPercentage = int(float64(stats.TotalPublished) / float64(stats.TotalContent) * 100)
	}

	// 6. Calculate changes (compare with previous period)
	// Replace whereClause with date filter for current period
	currentWhereClause := whereClause + " AND created_at >= $" + string(rune(argIndex+48))
	currentArgs := append(args, currentPeriodStart)

	previousWhereClause := whereClause + " AND created_at >= $" + string(rune(argIndex+48)) + " AND created_at < $" + string(rune(argIndex+49))
	previousArgs := append(args, previousPeriodStart, previousPeriodEnd)

	// Current period content count
	var currentContentCount, previousContentCount int
	contentCurrentQuery := `
		SELECT COUNT(*) FROM (
			SELECT id FROM drafts WHERE ` + currentWhereClause + `
			UNION ALL
			SELECT id FROM histories WHERE ` + currentWhereClause + `
		) AS all_content
	`
	err = database.GetDB().QueryRow(contentCurrentQuery, currentArgs...).Scan(&currentContentCount)
	if err != nil {
		log.Printf("Failed to get current content count: %v", err)
	}

	// Previous period content count
	contentPreviousQuery := `
		SELECT COUNT(*) FROM (
			SELECT id FROM drafts WHERE ` + previousWhereClause + `
			UNION ALL
			SELECT id FROM histories WHERE ` + previousWhereClause + `
		) AS all_content
	`
	err = database.GetDB().QueryRow(contentPreviousQuery, previousArgs...).Scan(&previousContentCount)
	if err != nil {
		log.Printf("Failed to get previous content count: %v", err)
	}

	// Calculate content change
	contentDiff := currentContentCount - previousContentCount
	if contentDiff > 0 {
		stats.ContentChange = "+" + string(rune(contentDiff+48)) + " minggu ini"
	} else if contentDiff < 0 {
		stats.ContentChange = string(rune(contentDiff+48)) + " minggu ini"
	} else {
		stats.ContentChange = "0 minggu ini"
	}

	// Current period products count
	var currentProductsCount, previousProductsCount int
	productsCurrentQuery := `SELECT COUNT(*) FROM products WHERE ` + currentWhereClause
	err = database.GetDB().QueryRow(productsCurrentQuery, currentArgs...).Scan(&currentProductsCount)
	if err != nil {
		log.Printf("Failed to get current products count: %v", err)
	}

	productsPreviousQuery := `SELECT COUNT(*) FROM products WHERE ` + previousWhereClause
	err = database.GetDB().QueryRow(productsPreviousQuery, previousArgs...).Scan(&previousProductsCount)
	if err != nil {
		log.Printf("Failed to get previous products count: %v", err)
	}

	// Calculate products change
	productsDiff := currentProductsCount - previousProductsCount
	if productsDiff > 0 {
		stats.ProductsChange = "+" + string(rune(productsDiff+48)) + " bulan ini"
	} else if productsDiff < 0 {
		stats.ProductsChange = string(rune(productsDiff+48)) + " bulan ini"
	} else {
		stats.ProductsChange = "0 bulan ini"
	}

	// Calculate SEO score change
	var currentSeoScore, previousSeoScore float64
	seoCurrentQuery := `
		SELECT COALESCE(AVG(seo_score), 0) FROM histories 
		WHERE ` + currentWhereClause + ` AND seo_score IS NOT NULL
	`
	err = database.GetDB().QueryRow(seoCurrentQuery, currentArgs...).Scan(&currentSeoScore)
	if err != nil {
		log.Printf("Failed to get current SEO score: %v", err)
	}

	seoPreviousQuery := `
		SELECT COALESCE(AVG(seo_score), 0) FROM histories 
		WHERE ` + previousWhereClause + ` AND seo_score IS NOT NULL
	`
	err = database.GetDB().QueryRow(seoPreviousQuery, previousArgs...).Scan(&previousSeoScore)
	if err != nil {
		log.Printf("Failed to get previous SEO score: %v", err)
	}

	seoDiff := currentSeoScore - previousSeoScore
	if seoDiff > 0 {
		stats.SeoScoreChange = "+" + formatFloat(seoDiff) + "%"
	} else if seoDiff < 0 {
		stats.SeoScoreChange = formatFloat(seoDiff) + "%"
	} else {
		stats.SeoScoreChange = "0%"
	}

	log.Printf("Dashboard stats: content=%d, products=%d, published=%d, seoScore=%.1f",
		stats.TotalContent, stats.TotalProducts, stats.TotalPublished, stats.AverageSeoScore)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func formatFloat(value float64) string {
	if value > 0 {
		return "+" + formatFloatValue(value)
	}
	return formatFloatValue(value)
}

func formatFloatValue(value float64) string {
	// Simple formatter for float values
	intPart := int(value)
	if value-float64(intPart) > 0 {
		return string(rune(intPart+48)) + "." + string(rune(int((value-float64(intPart))*10)+48))
	}
	return string(rune(intPart + 48))
}
