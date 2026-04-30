package dashboard

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"seo-backend/internal/middleware/auth"
)

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

// DashboardStats represents the dashboard statistics
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

// UserContext for dashboard
type dashboardUserContext struct {
	UserID string
	TeamID string
	Role   string
}

// Helper: Get user context from request
func (h *DashboardHandler) getUserContext(r *http.Request) dashboardUserContext {
	ctx := r.Context()
	return dashboardUserContext{
		UserID: auth.GetUserID(ctx),
		TeamID: auth.GetTeamID(ctx),
		Role:   auth.GetUserRole(ctx),
	}
}

// Helper: Write JSON response
func (h *DashboardHandler) writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

// Helper: Build filter conditions based on role
func (h *DashboardHandler) buildWhereClause(ctx dashboardUserContext) (string, []interface{}) {
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

// Helper: Build where clause with date filter
func (h *DashboardHandler) buildWhereClauseWithDate(whereClause string, args []interface{}, date time.Time, operator string) (string, []interface{}) {
	newWhereClause := whereClause + " AND created_at " + operator + " $" + string(rune(len(args)+1+48))
	newArgs := append(args, date)
	return newWhereClause, newArgs
}

// Helper: Format float values for display
func formatFloatValue(value float64) string {
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

// Helper: Format change with percentage
func formatChange(change int, isFloat bool, value float64) string {
	if change > 0 {
		if isFloat {
			return "+" + formatFloatValue(value) + "%"
		}
		return "+" + string(rune(change+48))
	} else if change < 0 {
		if isFloat {
			return formatFloatValue(value) + "%"
		}
		return string(rune(change + 48))
	}
	return "0"
}
