package history

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	app "seo-backend/internal/application/history"
	"seo-backend/internal/domain/history"
	"seo-backend/internal/domain/paginate"
	"seo-backend/internal/helper"
	auth "seo-backend/internal/middleware"
	"seo-backend/internal/models"
)

// HistoryHandler handles HTTP requests for history
type HistoryHandler struct {
	service *app.Service
}

// NewHistoryHandler creates a new history handler
func NewHistoryHandler(service *app.Service) *HistoryHandler {
	return &HistoryHandler{service: service}
}

// =======================
// HELPERS
// =======================

// getUserContext extracts user context from request
func (h *HistoryHandler) getUserContext(r *http.Request) *models.SimpleUserContext {
	ctx := r.Context()
	return &models.SimpleUserContext{
		UserID: auth.GetUserID(ctx),
		TeamID: auth.GetTeamID(ctx),
		Role:   auth.GetUserRole(ctx),
	}
}

// writeJSON writes JSON response
func (h *HistoryHandler) writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

// writeError writes error response
func (h *HistoryHandler) writeError(w http.ResponseWriter, message string, status int) {
	h.writeJSON(w, map[string]string{"error": message}, status)
}

// handleServiceError handles service layer errors
func (h *HistoryHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch err.Error() {
	case "title is required":
		h.writeError(w, err.Error(), http.StatusBadRequest)
	case "topic is required":
		h.writeError(w, err.Error(), http.StatusBadRequest)
	case "content is required":
		h.writeError(w, err.Error(), http.StatusBadRequest)
	case "createdBy is required":
		h.writeError(w, err.Error(), http.StatusBadRequest)
	case "teamId is required":
		h.writeError(w, err.Error(), http.StatusBadRequest)
	case "history id is required":
		h.writeError(w, err.Error(), http.StatusBadRequest)
	case "history not found":
		h.writeError(w, "History not found", http.StatusNotFound)
	default:
		log.Printf("Unexpected error: %v", err)
		h.writeError(w, "Internal server error", http.StatusInternalServerError)
	}
}

// parseFilters parses query parameters into filters
func (h *HistoryHandler) parseFilters(r *http.Request) history.HistoryFilter {
	userCtx := h.getUserContext(r)

	filters := history.HistoryFilter{
		TeamID:  userCtx.GetTeamID(),
		Status:  r.URL.Query().Get("status"),
		Topic:   r.URL.Query().Get("topic"),
		Search:  r.URL.Query().Get("search"),
		OrderBy: r.URL.Query().Get("orderBy"),
	}

	// Parse date filters
	if fromDate := r.URL.Query().Get("fromDate"); fromDate != "" {
		// parse fromDate string to time.Time
	}

	if toDate := r.URL.Query().Get("toDate"); toDate != "" {
		// parse toDate string to time.Time
	}

	return filters
}

// =======================
// CREATE HISTORY
// =======================
// @Summary Create a new history record
// @Tags History
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body app.CreateHistoryRequest true "History data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /history [post]

// =======================
// GET HISTORY BY ID
// =======================
// @Summary Get history by ID
// @Tags History
// @Produce json
// @Security BearerAuth
// @Param id path string true "History ID"
// @Success 200 {object} history.History
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /history/{id} [get]
func (h *HistoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.writeError(w, "History ID is required", http.StatusBadRequest)
		return
	}

	data, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("Failed to get history by ID: %v", err)
		h.handleServiceError(w, err)
		return
	}

	if data == nil {
		h.writeError(w, "History not found", http.StatusNotFound)
		return
	}

	h.writeJSON(w, data, http.StatusOK)
}

// =======================
// GET ALL HISTORY
// =======================
// @Summary Get all history records
// @Tags History
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status"
// @Param topic query string false "Filter by topic"
// @Param search query string false "Search by title or content"
// @Param limit query int false "Page size" default(20)
// @Param offset query int false "Page offset" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /history [get]
func (h *HistoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	filters := h.parseFilters(r)
	var paginate paginate.PaginationParams
	paginate = helper.ParsePaginationParams(r)

	filters.Limit = paginate.Limit
	filters.Offset = paginate.Offset

	user := auth.GetUserContext(r)
	// Use filters to get history

	log.Printf("FILTER DEBUG: %+v\n", filters)
	log.Printf("USER DEBUG: %+v\n", user)

	records, err := h.service.GetAll(r.Context(), &user, filters)
	if err != nil {
		log.Printf("Failed to get history: %v", err)
		h.handleServiceError(w, err)
		return
	}

	h.writeJSON(w, records, http.StatusOK)
}

// =======================
// GET HISTORY BY TEAM
// =======================
// @Summary Get history by team ID
// @Tags History
// @Produce json
// @Security BearerAuth
// @Success 200 {array} history.History
// @Failure 500 {object} map[string]string
// @Router /history/team [get]
func (h *HistoryHandler) GetByTeam(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)
	teamID := userCtx.GetTeamID()

	if teamID == "" {
		h.writeError(w, "Team ID not found", http.StatusBadRequest)
		return
	}

	data, err := h.service.GetByTeamID(r.Context(), teamID)
	if err != nil {
		log.Printf("Failed to get history by team: %v", err)
		h.handleServiceError(w, err)
		return
	}

	h.writeJSON(w, data, http.StatusOK)
}

// =======================
// GET HISTORY BY STATUS
// =======================
// @Summary Get history by status
// @Tags History
// @Produce json
// @Security BearerAuth
// @Param status path string true "Status (pending, completed, failed)"
// @Success 200 {array} history.History
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /history/status/{status} [get]
func (h *HistoryHandler) GetByStatus(w http.ResponseWriter, r *http.Request) {
	status := chi.URLParam(r, "status")
	if status == "" {
		h.writeError(w, "Status is required", http.StatusBadRequest)
		return
	}

	data, err := h.service.GetByStatus(r.Context(), status)
	if err != nil {
		log.Printf("Failed to get history by status: %v", err)
		h.handleServiceError(w, err)
		return
	}

	h.writeJSON(w, data, http.StatusOK)
}

// =======================
// UPDATE HISTORY
// =======================
// @Summary Update a history record
// @Tags History
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "History ID"
// @Param request body app.UpdateHistoryRequest true "Update data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /history/{id} [put]
func (h *HistoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.writeError(w, "History ID is required", http.StatusBadRequest)
		return
	}

	var req app.UpdateHistoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.Update(r.Context(), id, req); err != nil {
		log.Printf("Failed to update history: %v", err)
		h.handleServiceError(w, err)
		return
	}

	h.writeJSON(w, map[string]string{
		"id":      id,
		"message": "History updated successfully",
	}, http.StatusOK)
}

// =======================
// UPDATE HISTORY STATUS
// =======================
// @Summary Update history status
// @Tags History
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "History ID"
// @Param request body object true "Status data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /history/{id}/status [patch]
func (h *HistoryHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.writeError(w, "History ID is required", http.StatusBadRequest)
		return
	}

	var req struct {
		Status       string  `json:"status"`
		ErrorMessage *string `json:"errorMessage,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Status == "" {
		h.writeError(w, "Status is required", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateStatus(r.Context(), id, req.Status, req.ErrorMessage); err != nil {
		log.Printf("Failed to update history status: %v", err)
		h.handleServiceError(w, err)
		return
	}

	h.writeJSON(w, map[string]string{
		"id":      id,
		"status":  req.Status,
		"message": "History status updated successfully",
	}, http.StatusOK)
}

// =======================
// DELETE HISTORY
// =======================
// @Summary Delete a history record
// @Tags History
// @Produce json
// @Security BearerAuth
// @Param id path string true "History ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /history/{id} [delete]
func (h *HistoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.writeError(w, "History ID is required", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		log.Printf("Failed to delete history: %v", err)
		h.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// =======================
// GET HISTORY STATISTICS
// =======================
// @Summary Get history statistics
// @Tags History
// @Produce json
// @Security BearerAuth
// @Success 200 {object} history.HistoryStats
// @Failure 500 {object} map[string]string
// @Router /history/stats [get]
func (h *HistoryHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)

	query := &history.HistoryFilter{
		TeamID: userCtx.GetTeamID(),
	}

	stats, err := h.service.GetStats(r.Context(), query)
	if err != nil {
		log.Printf("Failed to get history stats: %v", err)
		h.handleServiceError(w, err)
		return
	}

	h.writeJSON(w, stats, http.StatusOK)
}

// =======================
// GET RECENT ACTIVITY
// =======================
// @Summary Get recent history activity
// @Tags History
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Number of records" default(10)
// @Success 200 {array} history.History
// @Failure 500 {object} map[string]string
// @Router /history/recent [get]
func (h *HistoryHandler) GetRecentActivity(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)
	teamID := userCtx.GetTeamID()

	if teamID == "" {
		h.writeError(w, "Team ID not found", http.StatusBadRequest)
		return
	}

	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	activities, err := h.service.GetRecentActivity(r.Context(), teamID, limit)
	if err != nil {
		log.Printf("Failed to get recent activity: %v", err)
		h.handleServiceError(w, err)
		return
	}

	h.writeJSON(w, activities, http.StatusOK)
}

// =======================
// DELETE BY TEAM
// =======================
// @Summary Delete all history records for a team
// @Tags History
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /history/team [delete]
func (h *HistoryHandler) DeleteByTeam(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)
	teamID := userCtx.GetTeamID()

	if teamID == "" {
		h.writeError(w, "Team ID not found", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteByTeamID(r.Context(), teamID); err != nil {
		log.Printf("Failed to delete history by team: %v", err)
		h.handleServiceError(w, err)
		return
	}

	h.writeJSON(w, map[string]string{
		"message": "All history records for team deleted successfully",
	}, http.StatusOK)
}
