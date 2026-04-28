package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"seo-backend/internal/builder"
	"seo-backend/internal/database"
	"seo-backend/internal/middleware/auth"
	"seo-backend/internal/models"
)

type HistoryHandler struct{}

func NewHistoryHandler() *HistoryHandler {
	return &HistoryHandler{}
}

// GetAll history with filters
func (h *HistoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	teamID := auth.GetTeamID(ctx)
	userRole := auth.GetUserRole(ctx)
	userID := auth.GetUserID(ctx)

	status := r.URL.Query().Get("status")
	action := r.URL.Query().Get("action")
	search := r.URL.Query().Get("search")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	// Build query dengan generic builder
	qb := builder.NewQueryBuilder("histories")

	query := qb.Select(
		"id", "title", "topic", "content", "image_url", "target_products",
		"status", "action", "error_message", "published_at", "scheduled_for",
		"created_by", "team_id", "created_at",
	).OrderBy("published_at DESC")

	// Logic filter berdasarkan role
	switch userRole {
	case "super_admin":
		// Super admin: melihat SEMUA history
		log.Printf("Super admin - melihat semua history")

	case "admin":
		// Admin: hanya melihat history dalam team yang sama
		if teamID != "" {
			query = query.WhereEq("team_id", teamID)
			log.Printf("Admin - melihat history dalam team: %s", teamID)
		} else {
			// Admin tanpa team hanya melihat history miliknya sendiri
			query = query.WhereEq("created_by", userID)

			log.Printf("Admin tanpa team - hanya melihat history sendiri")
		}

	default:
		// Role lain (manager, editor, viewer): hanya melihat history milik sendiri
		query = query.WhereEq("created_by", userID)
		log.Printf("User role %s - hanya melihat history sendiri", userRole)
	}

	// Filter tambahan
	if status != "" {
		query = query.WhereEq("status", status)
	}

	if action != "" {
		query = query.WhereEq("action", action)
	}

	if search != "" {
		// Menggunakan Where untuk OR condition
		query = query.Where("(title ILIKE ? OR topic ILIKE ?)", "%"+search+"%", "%"+search+"%")
	}

	// Pagination
	if limit != "" {
		var limitUint uint64
		fmt.Sscanf(limit, "%d", &limitUint)
		query = query.Limit(limitUint)
	}

	if offset != "" {
		var offsetUint uint64
		fmt.Sscanf(offset, "%d", &offsetUint)
		query = query.Offset(offsetUint)
	}

	sqlQuery, args, err := query.Build()
	if err != nil {
		log.Printf("Failed to build query: %v", err)
		http.Error(w, "Failed to fetch history", http.StatusInternalServerError)
		return
	}

	log.Printf("Query: %s | Args: %v", sqlQuery, args)

	var rows *sql.Rows
	rows, err = database.GetDB().Query(sqlQuery, args...)
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Failed to fetch history", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var histories []models.History
	for rows.Next() {
		var h models.History
		var targetProductsJSON []byte

		err := rows.Scan(
			&h.ID, &h.Title, &h.Topic, &h.Content, &h.ImageURL,
			&targetProductsJSON, &h.Status, &h.Action, &h.ErrorMessage,
			&h.PublishedAt, &h.ScheduledFor, &h.CreatedBy, &h.TeamID, &h.CreatedAt,
		)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}

		json.Unmarshal(targetProductsJSON, &h.TargetProducts)
		histories = append(histories, h)
	}

	log.Printf("Successfully fetched %d history records", len(histories))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(histories)
}

// GetByID retrieves a single history entry
func (h *HistoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	teamID := auth.GetTeamID(ctx)
	userRole := auth.GetUserRole(ctx)
	userID := auth.GetUserID(ctx)

	id := chi.URLParam(r, "id")

	qb := builder.NewQueryBuilder("histories")

	sqlQuery, args, err := qb.Select(
		"id", "title", "topic", "content", "image_url", "target_products",
		"status", "action", "error_message", "published_at", "scheduled_for",
		"created_by", "team_id", "created_at",
	).WhereEq("id", id).Build()

	if err != nil {
		log.Printf("Failed to build query: %v", err)
		http.Error(w, "Failed to fetch history", http.StatusInternalServerError)
		return
	}

	var history models.History
	var targetProductsJSON []byte

	err = database.GetDB().QueryRow(sqlQuery, args...).Scan(
		&history.ID, &history.Title, &history.Topic, &history.Content, &history.ImageURL,
		&targetProductsJSON, &history.Status, &history.Action, &history.ErrorMessage,
		&history.PublishedAt, &history.ScheduledFor, &history.CreatedBy, &history.TeamID, &history.CreatedAt,
	)
	if err != nil {
		http.Error(w, "History not found", http.StatusNotFound)
		return
	}

	// Authorization check
	if userRole != "admin" && userRole != "super_admin" {
		if history.TeamID != nil && teamID != "" && *history.TeamID != teamID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if history.CreatedBy != nil && *history.CreatedBy != userID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}

	json.Unmarshal(targetProductsJSON, &history.TargetProducts)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

// Delete history entry
func (h *HistoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userRole := auth.GetUserRole(ctx)
	teamID := auth.GetTeamID(ctx)
	userID := auth.GetUserID(ctx)

	id := chi.URLParam(r, "id")

	// Check permission berdasarkan role
	var createdBy *string
	var historyTeamID *string
	err := database.GetDB().QueryRow(`
		SELECT created_by, team_id FROM histories WHERE id = $1
	`, id).Scan(&createdBy, &historyTeamID)

	if err != nil {
		http.Error(w, "History not found", http.StatusNotFound)
		return
	}

	// Authorization check berdasarkan role
	switch userRole {
	case "super_admin":
		// Super admin: bisa hapus semua history
		log.Printf("Super admin - can delete history %s", id)

	case "admin":
		// Admin: hanya bisa hapus history dalam team yang sama
		if historyTeamID != nil && teamID != "" && *historyTeamID != teamID {
			log.Printf("Admin access denied - history team %v != user team %s", historyTeamID, teamID)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		log.Printf("Admin - can delete history %s", id)

	default:
		// Role lain: hanya bisa hapus history milik sendiri
		if createdBy == nil || *createdBy != userID {
			log.Printf("User %s access denied - history owner %v", userID, createdBy)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		log.Printf("User %s - can delete own history %s", userID, id)
	}

	qb := builder.NewQueryBuilder("histories")

	sqlQuery, args, err := qb.Delete().WhereEq("id", id).Build()
	if err != nil {
		log.Printf("Failed to build delete query: %v", err)
		http.Error(w, "Failed to delete history", http.StatusInternalServerError)
		return
	}

	log.Printf("Delete query: %s | Args: %v", sqlQuery, args)

	_, err = database.GetDB().Exec(sqlQuery, args...)
	if err != nil {
		log.Printf("Failed to delete history: %v", err)
		http.Error(w, "Failed to delete history", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ClearAll clears all history (admin only)
func (h *HistoryHandler) ClearAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userRole := auth.GetUserRole(ctx)

	if userRole != "admin" && userRole != "super_admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	qb := builder.NewQueryBuilder("histories")

	sqlQuery, args, err := qb.Delete().Build()
	if err != nil {
		log.Printf("Failed to build delete query: %v", err)
		http.Error(w, "Failed to clear history", http.StatusInternalServerError)
		return
	}

	log.Printf("Clear all query: %s", sqlQuery)

	_, err = database.GetDB().Exec(sqlQuery, args...)
	if err != nil {
		log.Printf("Failed to clear history: %v", err)
		http.Error(w, "Failed to clear history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "All history cleared successfully"})
}

// AddToHistory adds a new history entry
func (h *HistoryHandler) AddToHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := auth.GetUserID(ctx)
	teamID := auth.GetTeamID(ctx)

	var req struct {
		Title          string   `json:"title"`
		Topic          string   `json:"topic"`
		Content        string   `json:"content"`
		ImageURL       *string  `json:"imageUrl,omitempty"`
		TargetProducts []string `json:"targetProducts"`
		Status         string   `json:"status"`
		Action         string   `json:"action"`
		ErrorMessage   *string  `json:"errorMessage,omitempty"`
		ScheduledFor   *string  `json:"scheduledFor,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	targetProductsJSON, _ := json.Marshal(req.TargetProducts)

	var createdBy interface{}
	if userID != "" {
		createdBy = userID
	}

	var teamIDPtr interface{}
	if teamID != "" {
		teamIDPtr = teamID
	}

	qb := builder.NewQueryBuilder("histories")

	sqlQuery, args, err := qb.Insert().
		Columns(
			"title", "topic", "content", "image_url", "target_products",
			"status", "action", "error_message", "scheduled_for", "published_at",
			"created_by", "team_id",
		).
		Values(
			req.Title, req.Topic, req.Content, req.ImageURL, targetProductsJSON,
			req.Status, req.Action, req.ErrorMessage, req.ScheduledFor,
			"NOW()", createdBy, teamIDPtr,
		).
		Returning("id").
		Build()

	if err != nil {
		log.Printf("Failed to build insert query: %v", err)
		http.Error(w, "Failed to add to history", http.StatusInternalServerError)
		return
	}

	log.Printf("Insert query: %s | Args: %v", sqlQuery, args)

	var historyID string
	err = database.GetDB().QueryRow(sqlQuery, args...).Scan(&historyID)
	if err != nil {
		log.Printf("Failed to add to history: %v", err)
		http.Error(w, "Failed to add to history", http.StatusInternalServerError)
		return
	}

	log.Printf("History added with ID: %s", historyID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": historyID, "message": "Added to history successfully"})
}

// GetStats returns history statistics
func (h *HistoryHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	teamID := auth.GetTeamID(ctx)
	userRole := auth.GetUserRole(ctx)
	userID := auth.GetUserID(ctx)

	// Build query dengan builder untuk statistik
	qb := builder.NewQueryBuilder("histories")

	query := qb.Select(
		"COUNT(*) as total",
		"COUNT(CASE WHEN status = 'success' THEN 1 END) as success_count",
		"COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_count",
		"COUNT(CASE WHEN action = 'published' THEN 1 END) as published_count",
		"COUNT(CASE WHEN action = 'scheduled' THEN 1 END) as scheduled_count",
	)

	// Filter berdasarkan role
	if userRole != "admin" && userRole != "super_admin" {
		if teamID != "" {
			query = query.WhereEq("team_id", teamID)
		} else {
			query = query.WhereEq("created_by", userID)
		}
	}

	sqlQuery, args, err := query.Build()
	if err != nil {
		log.Printf("Failed to build stats query: %v", err)
		http.Error(w, "Failed to fetch stats", http.StatusInternalServerError)
		return
	}

	log.Printf("Stats query: %s | Args: %v", sqlQuery, args)

	var total, successCount, failedCount, publishedCount, scheduledCount int

	err = database.GetDB().QueryRow(sqlQuery, args...).Scan(
		&total, &successCount, &failedCount, &publishedCount, &scheduledCount,
	)
	if err != nil {
		log.Printf("Failed to fetch stats: %v", err)
		http.Error(w, "Failed to fetch stats", http.StatusInternalServerError)
		return
	}

	var successRate float64
	if total > 0 {
		successRate = float64(successCount) / float64(total) * 100
	}

	stats := map[string]interface{}{
		"total":       total,
		"success":     successCount,
		"failed":      failedCount,
		"published":   publishedCount,
		"scheduled":   scheduledCount,
		"successRate": successRate,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
