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
	"seo-backend/internal/middleware/auth"
)

type DraftHandler struct {
	draftService draft.Service
}

func NewDraftHandler(draftService draft.Service) *DraftHandler {
	return &DraftHandler{
		draftService: draftService,
	}
}

// GetAll - get all drafts with filters
func (h *DraftHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	usrCtx := helper.GetUserContext(r)

	draftData, err := h.draftService.GetAll(ctx, usrCtx.GetTeamID())

	if err != nil {
		http.Error(w, "Draft not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(draftData)
}

func (h *DraftHandler) GetAllScheduled(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	usrCtx := helper.GetUserContext(r)

	draftData, err := h.draftService.GetAllScheduled(ctx, usrCtx.GetTeamID())

	if err != nil {
		http.Error(w, "Draft not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(draftData)
}

// GetByID - get draft by ID
func (h *DraftHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	draftData, err := h.draftService.GetDraftByID(ctx, id)
	if err != nil {
		http.Error(w, "Draft not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(draftData)
}

// Create - create new draft
func (h *DraftHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := auth.GetUserID(ctx)
	teamID := auth.GetTeamID(ctx)

	var req draft.CreateDraftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := validateCreateRequest(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	draftID, err := h.draftService.CreateDraft(ctx, req, userID, teamID)
	if err != nil {
		log.Printf("Failed to create draft: %v", err)
		http.Error(w, "Failed to create draft", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id":      draftID,
		"message": "Draft created successfully",
	})
}

// Update - update existing draft
func (h *DraftHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.draftService.UpdateDraft(ctx, id, updates); err != nil {
		log.Printf("Failed to update draft: %v", err)
		http.Error(w, "Failed to update draft", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Draft updated successfully"})
}

// Delete - delete draft
func (h *DraftHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if err := h.draftService.DeleteDraft(ctx, id); err != nil {
		http.Error(w, "Draft not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Publish - publish draft by ID
func (h *DraftHandler) Publish(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	teamID := auth.GetTeamID(ctx)
	userID := auth.GetUserID(ctx)

	var req draft.CreateDraftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	result, err := h.draftService.PublishDraft(ctx, id, req, teamID, userID)
	if err != nil {
		log.Printf("Failed to publish draft: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writePublishResponse(w, result)
}

// PublishContent - publish content directly (without saving to drafts)
func (h *DraftHandler) PublishContent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Println("========== PUBLISH CONTENT ==========")
	log.Printf("Method: %s | URL: %s\n", r.Method, r.URL.Path)

	teamID := auth.GetTeamID(ctx)
	userID := auth.GetUserID(ctx)

	log.Printf("TeamID: %s\n", teamID)
	log.Printf("UserID: %s\n", userID)

	// READ RAW BODY
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {

		log.Printf("Failed read request body: %v\n", err)

		http.Error(
			w,
			"Failed to read request body",
			http.StatusBadRequest,
		)

		return
	}

	log.Printf("RAW BODY: %s\n", string(bodyBytes))

	// restore body
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var req draft.DraftDataPost

	log.Println("Decoding request body...")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		log.Printf("Failed decode request body: %v\n", err)

		http.Error(
			w,
			"Invalid request body",
			http.StatusBadRequest,
		)

		return
	}

	log.Printf("Parsed Request: %+v\n", req)

	log.Println("Calling PublishContent service...")

	log.Printf("==================DATA : %v", req)

	result, err := h.draftService.PublishContent(
		ctx,
		req,
		teamID,
		userID,
	)

	if err != nil {

		log.Printf("Failed to publish content: %v\n", err)

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	log.Printf("Publish Result: %+v\n", result)

	log.Println("Writing response...")

	h.writePublishResponse(w, result)

	log.Println("========== END PUBLISH CONTENT ==========")
}

// ScheduleDraft - schedule a draft for future publishing
func (h *DraftHandler) ScheduleDraft(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := auth.GetUserID(ctx)
	teamID := auth.GetTeamID(ctx)

	var req draft.ScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	draftID, err := h.draftService.ScheduleDraft(ctx, req, teamID, userID)
	if err != nil {
		log.Printf("Failed to schedule draft: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Draft scheduled successfully",
		"draft_id": draftID,
		"status":   "scheduled",
	})
}

// CancelScheduledDraft - cancel scheduled draft
func (h *DraftHandler) CancelScheduledDraft(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		DraftID string `json:"draft_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.draftService.CancelSchedule(ctx, req.DraftID); err != nil {
		log.Printf("Failed to cancel schedule: %v", err)
		http.Error(w, "Failed to cancel schedule", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Schedule cancelled",
		"draft_id": req.DraftID,
		"status":   "draft",
	})
}

// GetSEOScore - get SEO score for a draft
func (h *DraftHandler) GetSEOScore(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	score, err := h.draftService.GetSEOScore(r.Context(), id)
	if err != nil {
		http.Error(w, "Draft not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(score)
}

// Helper methods
func (h *DraftHandler) writePublishResponse(w http.ResponseWriter, result *draft.PublishResult) {
	response := map[string]interface{}{
		"message": "Draft processed",
		"status":  result.Status,
		"results": result.Results,
	}

	if result.ScheduledFor != nil {
		response["scheduled_for"] = result.ScheduledFor.Format(time.RFC3339)
	}

	statusCode := http.StatusOK
	if result.AllFailed {
		statusCode = http.StatusBadGateway
	} else if result.SomeFailed {
		statusCode = http.StatusMultiStatus
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func validateCreateRequest(req draft.CreateDraftRequest) error {
	if req.Title == "" || req.Topic == "" || req.Article == "" {
		return fmt.Errorf("title, topic, and article are required")
	}
	return nil
}
