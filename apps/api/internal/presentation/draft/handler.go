// internal/interfaces/api/draft/handler.go
package draft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"seo-backend/internal/domain/draft"
	"seo-backend/internal/helper"
	auth "seo-backend/internal/middleware"
)

type DraftHandler struct {
	draftService draft.Service
}

func NewDraftHandler(draftService draft.Service) *DraftHandler {
	return &DraftHandler{
		draftService: draftService,
	}
}

// Standard response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// Helper function to write JSON response
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

// Helper function to write error response
func writeErrorResponse(w http.ResponseWriter, statusCode int, message string, err error) {
	response := APIResponse{
		Success: false,
		Error:   message,
	}

	if err != nil {
		response.Message = err.Error()
	}

	writeJSONResponse(w, statusCode, response)
}

// Helper function to write success response
func writeSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}, message string) {
	response := APIResponse{
		Success: true,
		Data:    data,
		Message: message,
	}

	writeJSONResponse(w, statusCode, response)
}

// GetAll godoc
// @Summary Get all drafts
// @Description Get all drafts with pagination and filters
// @Tags Drafts
// @Accept json
// @Produce json
// @Param limit query int false "Items per page (default: 10)"
// @Param page query int false "Page number (default: 1)"
// @Param search query string false "Search term"
// @Param status query string false "Filter by status (draft, scheduled, published)" Enums(draft, scheduled, published)
// @Success 200 {object} APIResponse{data=object{drafts=object,stats=object}} "Drafts and stats"
// @Failure 404 {object} APIResponse "Drafts not found"
// @Failure 500 {object} APIResponse "Internal server error"
// @Security BearerAuth
// @Router /api/drafts [get]
func (h *DraftHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	usrCtx := auth.GetUserContext(r)

	draftData, err := h.draftService.GetAll(ctx, usrCtx, helper.ParsePaginationParams(r))
	if err != nil {
		log.Printf("[ERROR] Failed to get drafts: %v", err)
		writeErrorResponse(w, http.StatusNotFound, "Drafts not found", err)
		return
	}

	stats, err := h.draftService.GetDashboardStats(ctx, usrCtx)
	if err != nil {
		log.Printf("[ERROR] Failed to get dashboard stats: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to get dashboard stats", err)
		return
	}

	response := APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"drafts": draftData,
			"stats":  stats,
		},
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// GetAllScheduled godoc
// @Summary Get all scheduled drafts
// @Description Get all drafts that are scheduled for future publishing
// @Tags Drafts
// @Accept json
// @Produce json
// @Param limit query int false "Items per page (default: 10)"
// @Param page query int false "Page number (default: 1)"
// @Param search query string false "Search term"
// @Success 200 {object} APIResponse{data=object} "Scheduled drafts"
// @Failure 404 {object} APIResponse "Scheduled drafts not found"
// @Failure 500 {object} APIResponse "Internal server error"
// @Security BearerAuth
// @Router /api/drafts/scheduled [get]
func (h *DraftHandler) GetAllScheduled(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	usrCtx := auth.GetUserContext(r)

	draftData, err := h.draftService.GetAllScheduled(ctx, usrCtx, helper.ParsePaginationParams(r))
	if err != nil {
		log.Printf("[ERROR] Failed to get scheduled drafts: %v", err)
		writeErrorResponse(w, http.StatusNotFound, "Scheduled drafts not found", err)
		return
	}

	response := APIResponse{
		Success: true,
		Data:    draftData,
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// GetByID godoc
// @Summary Get draft by ID
// @Description Get a specific draft by its ID
// @Tags Drafts
// @Accept json
// @Produce json
// @Param id path string true "Draft ID"
// @Success 200 {object} APIResponse{data=draft.Draft} "Draft details"
// @Failure 404 {object} APIResponse "Draft not found"
// @Security BearerAuth
// @Router /api/drafts/{id} [get]
func (h *DraftHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	draftData, err := h.draftService.GetDraftByID(ctx, id)
	if err != nil {
		log.Printf("[ERROR] Draft not found: %v", err)
		writeErrorResponse(w, http.StatusNotFound, "Draft not found", err)
		return
	}

	response := APIResponse{
		Success: true,
		Data:    draftData,
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// Create godoc
// @Summary Create a new draft
// @Description Create a new draft article
// @Tags Drafts
// @Accept json
// @Produce json
// @Param request body draft.CreateDraftRequest true "Draft creation request"
// @Success 201 {object} APIResponse{data=object{id=string}} "Draft created"
// @Failure 400 {object} APIResponse "Invalid request or validation failed"
// @Failure 500 {object} APIResponse "Internal server error"
// @Security BearerAuth
// @Router /api/drafts [post]
func (h *DraftHandler) Create(w http.ResponseWriter, r *http.Request) {
	log.Println("========== START CREATE DRAFT HANDLER ==========")

	ctx := r.Context()
	userID := auth.GetUserID(ctx)
	teamID := auth.GetTeamID(ctx)

	log.Printf("REQUEST METHOD => %s", r.Method)
	log.Printf("REQUEST URL => %s", r.URL.String())
	log.Printf("USER_ID => %s", userID)
	log.Printf("TEAM_ID => %s", teamID)

	// Read raw body for debugging
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("FAILED READ BODY => %v", err)
		writeErrorResponse(w, http.StatusBadRequest, "Failed to read request body", err)
		return
	}

	log.Printf("RAW BODY => %s", string(bodyBytes))

	// Restore body after reading
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var req draft.CreateDraftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("JSON DECODE ERROR => %v", err)
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	log.Printf("REQUEST STRUCT => %+v", req)

	if err := validateCreateRequest(req); err != nil {
		log.Printf("VALIDATION ERROR => %v", err)
		writeErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	draftID, err := h.draftService.CreateDraft(ctx, req, userID, teamID)
	if err != nil {
		log.Printf("SERVICE ERROR => %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to create draft", err)
		return
	}

	log.Printf("CREATE DRAFT SUCCESS => draftID=%s", draftID)

	response := APIResponse{
		Success: true,
		Data: map[string]string{
			"id": draftID,
		},
		Message: "Draft created successfully",
	}

	writeJSONResponse(w, http.StatusCreated, response)
	log.Println("========== END CREATE DRAFT HANDLER ==========")
}

// Update godoc
// @Summary Update a draft
// @Description Update an existing draft by ID
// @Tags Drafts
// @Accept json
// @Produce json
// @Param id path string true "Draft ID"
// @Param request body draft.CreateDraftRequest true "Draft update request"
// @Success 200 {object} APIResponse "Draft updated"
// @Failure 400 {object} APIResponse "Invalid request"
// @Failure 404 {object} APIResponse "Draft not found"
// @Failure 500 {object} APIResponse "Internal server error"
// @Security BearerAuth
// @Router /api/drafts/{id} [put]
func (h *DraftHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	teamID := auth.GetUserContext(r).GetTeamID()
	userID := auth.GetUserContext(r).GetUserID()
	id := chi.URLParam(r, "id")

	// Decode request body ke CreateDraftRequest
	var updates draft.CreateDraftRequest
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		log.Printf("[ERROR] Invalid request body: %v", err)
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Set team_id dan user_id dari context
	updates.TeamID = teamID
	updates.UserID = userID

	// Panggil service dengan CreateDraftRequest
	if err := h.draftService.UpdateDraft(ctx, id, userID, teamID, updates); err != nil {
		log.Printf("[ERROR] Failed to update draft: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to update draft", err)
		return
	}

	response := APIResponse{
		Success: true,
		Message: "Draft updated successfully",
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// Delete godoc
// @Summary Delete a draft
// @Description Delete a draft by ID
// @Tags Drafts
// @Accept json
// @Produce json
// @Param id path string true "Draft ID"
// @Success 200 {object} APIResponse "Draft deleted"
// @Failure 404 {object} APIResponse "Draft not found"
// @Failure 500 {object} APIResponse "Internal server error"
// @Security BearerAuth
// @Router /api/drafts/{id} [delete]
func (h *DraftHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := auth.GetUserContext(r)
	id := chi.URLParam(r, "id")

	if err := h.draftService.DeleteDraft(ctx, userCtx.GetTeamID(), id); err != nil {
		log.Printf("[ERROR] Draft not found: %v", err)
		writeErrorResponse(w, http.StatusNotFound, "Draft not found", err)
		return
	}

	response := APIResponse{
		Success: true,
		Message: "Draft deleted successfully",
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// Publish godoc
// @Summary Publish a draft
// @Description Publish a draft by ID
// @Tags Drafts
// @Accept json
// @Produce json
// @Param id path string true "Draft ID"
// @Param request body draft.CreateDraftRequest true "Publish request"
// @Success 200 {object} APIResponse{data=object{results=array,some_failed=boolean,all_failed=boolean,status=string,published_at=string,total_products=int,success_count=int,failed_count=int}} "Draft published"
// @Failure 400 {object} APIResponse "Invalid request"
// @Failure 500 {object} APIResponse "Internal server error"
// @Security BearerAuth
// @Router /api/drafts/{id}/publish [post]
func (h *DraftHandler) Publish(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	userContext := auth.GetUserContext(r)

	var req draft.CreateDraftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[ERROR] Invalid request body: %v", err)
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	result, err := h.draftService.PublishDraft(ctx, id, req, userContext)
	if err != nil {
		log.Printf("[ERROR] Failed to publish draft: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to publish draft", err)
		return
	}

	h.writePublishResponse(w, result, "Draft published successfully")
}

// PublishContent godoc
// @Summary Publish content directly
// @Description Publish content directly without saving to drafts
// @Tags Drafts
// @Accept json
// @Produce json
// @Param request body draft.CreateDraftRequest true "Content to publish"
// @Success 200 {object} APIResponse{data=object{results=array,some_failed=boolean,all_failed=boolean,status=string,published_at=string}} "Content published"
// @Failure 400 {object} APIResponse "Invalid request"
// @Failure 500 {object} APIResponse "Internal server error"
// @Security BearerAuth
// @Router /api/drafts/publish [post]
func (h *DraftHandler) PublishContent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("========== PUBLISH CONTENT ==========")
	log.Printf("Method: %s | URL: %s\n", r.Method, r.URL.Path)

	userCtx := auth.GetUserContext(r)

	// Read raw body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed read request body: %v\n", err)
		writeErrorResponse(w, http.StatusBadRequest, "Failed to read request body", err)
		return
	}

	log.Printf("RAW BODY: %s\n", string(bodyBytes))
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var req draft.CreateDraftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed decode request body: %v\n", err)
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	log.Printf("Parsed Request: %+v\n", req)

	result, err := h.draftService.PublishContent(ctx, req, userCtx)
	if err != nil {
		log.Printf("Failed to publish content: %v\n", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to publish content", err)
		return
	}

	log.Printf("Publish Result: %+v\n", result)
	h.writePublishResponse(w, result, "Content published successfully")
	log.Println("========== END PUBLISH CONTENT ==========")
}

// ScheduleDraft godoc
// @Summary Schedule a draft
// @Description Schedule a draft for future publishing
// @Tags Drafts
// @Accept json
// @Produce json
// @Param request body draft.ScheduleRequest true "Schedule request"
// @Success 201 {object} APIResponse{data=object{draft_id=string,status=string}} "Draft scheduled"
// @Failure 400 {object} APIResponse "Invalid request"
// @Failure 500 {object} APIResponse "Internal server error"
// @Security BearerAuth
// @Router /api/drafts/schedule [post]
func (h *DraftHandler) ScheduleDraft(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := auth.GetUserContext(r)

	var req draft.ScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[ERROR] Invalid request body: %v", err)
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	draftID, err := h.draftService.ScheduleDraft(ctx, req, userCtx)
	if err != nil {
		log.Printf("[ERROR] Failed to schedule draft: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to schedule draft", err)
		return
	}

	response := APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"draft_id": draftID,
			"status":   "scheduled",
		},
		Message: "Draft scheduled successfully",
	}

	writeJSONResponse(w, http.StatusCreated, response)
}

// CancelScheduledDraft godoc
// @Summary Cancel scheduled draft
// @Description Cancel a scheduled draft
// @Tags Drafts
// @Accept json
// @Produce json
// @Param request body object{draft_id=string} true "Cancel request"
// @Success 200 {object} APIResponse{data=object{draft_id=string,status=string}} "Schedule cancelled"
// @Failure 400 {object} APIResponse "Invalid request"
// @Failure 500 {object} APIResponse "Internal server error"
// @Security BearerAuth
// @Router /api/drafts/schedule/cancel [post]
func (h *DraftHandler) CancelScheduledDraft(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		DraftID string `json:"draft_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[ERROR] Invalid request body: %v", err)
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.draftService.CancelSchedule(ctx, req.DraftID); err != nil {
		log.Printf("[ERROR] Failed to cancel schedule: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to cancel schedule", err)
		return
	}

	response := APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"draft_id": req.DraftID,
			"status":   "draft",
		},
		Message: "Schedule cancelled successfully",
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// GetSEOScore godoc
// @Summary Get SEO score
// @Description Get SEO score for a draft
// @Tags Drafts
// @Accept json
// @Produce json
// @Param id path string true "Draft ID"
// @Success 200 {object} APIResponse{data=object} "SEO score"
// @Failure 404 {object} APIResponse "Draft not found"
// @Security BearerAuth
// @Router /api/drafts/{id}/seo-score [get]
func (h *DraftHandler) GetSEOScore(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	score, err := h.draftService.GetSEOScore(r.Context(), id)
	if err != nil {
		log.Printf("[ERROR] Draft not found: %v", err)
		writeErrorResponse(w, http.StatusNotFound, "Draft not found", err)
		return
	}

	response := APIResponse{
		Success: true,
		Data:    score,
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// CheckSimilarity godoc
// @Summary Check similarity
// @Description Check similarity of a draft with other drafts
// @Tags Drafts
// @Accept json
// @Produce json
// @Param id path string true "Draft ID"
// @Param limit query int false "Items per page (default: 10)"
// @Param page query int false "Page number (default: 1)"
// @Success 200 {object} APIResponse{data=object{draft_id=string,similar_drafts=array,total=int}} "Similarity results"
// @Failure 404 {object} APIResponse "Draft not found"
// @Security BearerAuth
// @Router /api/drafts/{id}/similarity [get]
func (h *DraftHandler) CheckSimilarity(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userCtx := auth.GetUserContext(r)

	results, err := h.draftService.CheckSimilarity(r.Context(), id, userCtx, helper.ParsePaginationParams(r))
	if err != nil {
		log.Printf("[ERROR] Draft not found: %v", err)
		writeErrorResponse(w, http.StatusNotFound, "Draft not found", err)
		return
	}

	response := APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"draft_id":       id,
			"similar_drafts": results,
			"total":          len(results),
		},
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// RescheduleDraft godoc
// @Summary Reschedule a draft
// @Description Reschedule a scheduled draft to a new date/time
// @Tags Drafts
// @Accept json
// @Produce json
// @Param id path string true "Draft ID"
// @Param request body object{scheduled_for=string} true "Reschedule request (ISO 8601 format: 2026-07-07T15:30:00+07:00)"
// @Success 200 {object} APIResponse{data=object{draft_id=string,scheduled_for=string,status=string,updated_at=string}} "Draft rescheduled"
// @Failure 400 {object} APIResponse "Invalid request or scheduled time must be in future"
// @Failure 404 {object} APIResponse "Draft not found"
// @Failure 500 {object} APIResponse "Internal server error"
// @Security BearerAuth
// @Router /api/drafts/{id}/reschedule [put]
func (h *DraftHandler) RescheduleDraft(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := auth.GetUserContext(r)
	id := chi.URLParam(r, "id")

	log.Printf("========== START RESCHEDULE DRAFT ==========")
	log.Printf("DRAFT_ID=%s TEAM_ID=%s USER_ID=%s",
		id,
		userCtx.GetTeamID(),
		userCtx.GetUserID(),
	)

	// Read raw body for debugging
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read request body: %v", err)
		writeErrorResponse(w, http.StatusBadRequest, "Failed to read request body", err)
		return
	}

	log.Printf("RAW BODY => %s", string(bodyBytes))
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Decode request body
	var req struct {
		ScheduledFor string `json:"scheduled_for"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[ERROR] Invalid request body: %v", err)
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	log.Printf("REQUEST => scheduled_for=%s", req.ScheduledFor)

	// Validate input
	if req.ScheduledFor == "" {
		log.Printf("[ERROR] scheduled_for is required")
		writeErrorResponse(w, http.StatusBadRequest, "scheduled_for is required", nil)
		return
	}

	// Parse and validate datetime format
	scheduledTime, err := time.Parse(time.RFC3339, req.ScheduledFor)
	if err != nil {
		// Coba parse format lain jika RFC3339 gagal
		scheduledTime, err = time.Parse("2006-01-02T15:04:05", req.ScheduledFor)
		if err != nil {
			log.Printf("[ERROR] Invalid date format: %v", err)
			writeErrorResponse(w, http.StatusBadRequest,
				"Invalid date format. Use ISO 8601 format (e.g., 2026-07-07T15:30:00+07:00)",
				err,
			)
			return
		}
	}

	// Validasi waktu harus di masa depan
	if scheduledTime.Before(time.Now()) {
		log.Printf("[ERROR] Scheduled time must be in the future: %v", scheduledTime)
		writeErrorResponse(w, http.StatusBadRequest,
			"Scheduled time must be in the future",
			fmt.Errorf("scheduled time %v is in the past", scheduledTime),
		)
		return
	}

	log.Printf("PARSED SCHEDULED TIME => %v", scheduledTime)

	// Panggil service untuk reschedule
	result, err := h.draftService.RescheduleDraft(ctx, id, scheduledTime, userCtx)
	if err != nil {
		log.Printf("[ERROR] Failed to reschedule draft: %v", err)

		// Handle specific errors
		switch err.Error() {
		case "draft not found":
			writeErrorResponse(w, http.StatusNotFound, "Draft not found", err)
		case "draft is not in scheduled status":
			writeErrorResponse(w, http.StatusBadRequest, "Draft is not in scheduled status", err)
		default:
			writeErrorResponse(w, http.StatusInternalServerError, "Failed to reschedule draft", err)
		}
		return
	}

	log.Printf("✅ RESCHEDULE SUCCESS => draftID=%s, newTime=%v", id, scheduledTime)

	// Build response
	response := APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"draft_id":      id,
			"scheduled_for": scheduledTime.Format(time.RFC3339),
			"status":        "scheduled",
			"updated_at":    time.Now().Format(time.RFC3339),
		},
		Message: "Draft rescheduled successfully",
	}

	// Add publish result info jika ada
	if result != nil {
		if result.ScheduledFor != nil {
			response.Data.(map[string]interface{})["previous_scheduled_for"] = result.ScheduledFor.Format(time.RFC3339)
		}
	}

	writeJSONResponse(w, http.StatusOK, response)
	log.Printf("========== END RESCHEDULE DRAFT ==========")
}

func validateCreateRequest(req draft.CreateDraftRequest) error {
	if req.Title == "" || req.Topic == "" || req.Article == "" {
		return fmt.Errorf("title, topic, and article are required")
	}
	return nil
}

// Helper methods
func (h *DraftHandler) writePublishResponse(w http.ResponseWriter, result *draft.PublishResult, message string) {
	data := map[string]interface{}{
		"results":      result.Results,
		"some_failed":  result.SomeFailed,
		"all_failed":   result.AllFailed,
		"status":       result.Status,
		"published_at": time.Now().Format(time.RFC3339),
	}

	if result.ScheduledFor != nil {
		data["scheduled_for"] = result.ScheduledFor.Format(time.RFC3339)
	}

	// Add statistics if available
	if result.TotalProducts > 0 {
		data["total_products"] = result.TotalProducts
		data["success_count"] = result.SuccessCount
		data["failed_count"] = result.FailedCount
	}

	if len(result.Errors) > 0 {
		data["errors"] = result.Errors
	}

	response := APIResponse{
		Success: !result.AllFailed,
		Data:    data,
		Message: message,
	}

	// Determine status code
	statusCode := http.StatusOK
	if result.AllFailed {
		statusCode = http.StatusBadGateway
		response.Error = "All products failed to publish"
	} else if result.SomeFailed {
		statusCode = http.StatusMultiStatus
		response.Message = "Some products failed to publish"
	}

	writeJSONResponse(w, statusCode, response)
}
