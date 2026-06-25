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

// GetAll - get all drafts with filters
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

// GetByID - get draft by ID
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

// Create - create new draft
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

// Update - update existing draft
func (h *DraftHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	teamID := auth.GetUserContext(r).GetTeamID()
	id := chi.URLParam(r, "id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		log.Printf("[ERROR] Invalid request body: %v", err)
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.draftService.UpdateDraft(ctx, id, teamID, updates); err != nil {
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

// Delete - delete draft
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

// Publish - publish draft by ID
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

// PublishContent - publish content directly (without saving to drafts)
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

	var req draft.DraftDataPost
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

// ScheduleDraft - schedule a draft for future publishing
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

// CancelScheduledDraft - cancel scheduled draft
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

// GetSEOScore - get SEO score for a draft
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

// CheckSimilarity - check similarity of a draft with others
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

func validateCreateRequest(req draft.CreateDraftRequest) error {
	if req.Title == "" || req.Topic == "" || req.Article == "" {
		return fmt.Errorf("title, topic, and article are required")
	}
	return nil
}
