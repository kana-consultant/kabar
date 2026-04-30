package draft

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"seo-backend/internal/database"
	"seo-backend/internal/middleware/auth"
	"seo-backend/internal/models"
	"seo-backend/internal/scheduler"
	services "seo-backend/internal/service/draft"
)

type DraftHandler struct {
	draftService *services.Service
	s            *services.Scheduler
}

func NewDraftHandler(redisScheduler *scheduler.RedisScheduler) *DraftHandler {
	repo := services.NewRepository(database.GetDB())
	return &DraftHandler{
		draftService: services.NewService(database.GetDB(), redisScheduler),
		s:            services.NewScheduler(repo, redisScheduler),
	}
}

func (h *DraftHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	teamID := auth.GetTeamID(ctx)
	userRole := auth.GetUserRole(ctx)
	userID := auth.GetUserID(ctx)

	query := NewDraftQueryBuilder()
	query.SetUserContext(userRole, teamID, userID)

	// Apply filters
	if status := r.URL.Query().Get("status"); status != "" {
		query.WithStatus(status)
	}
	if search := r.URL.Query().Get("search"); search != "" {
		query.WithSearch(search)
	}

	drafts, err := query.Execute()
	if err != nil {
		log.Printf("Failed to fetch drafts: %v", err)
		http.Error(w, "Failed to fetch drafts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drafts)
}

func (h *DraftHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	draft, err := h.draftService.GetDraftByID(id)
	if err != nil {
		http.Error(w, "Draft not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(draft)
}

func (h *DraftHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := auth.GetUserID(ctx)
	teamID := auth.GetTeamID(ctx)

	var req models.CreateDraftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := validateCreateRequest(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	draftID, err := h.draftService.CreateDraft(req, userID, teamID)
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

func (h *DraftHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.draftService.UpdateDraft(id, updates); err != nil {
		log.Printf("Failed to update draft: %v", err)
		http.Error(w, "Failed to update draft", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Draft updated successfully"})
}

func (h *DraftHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.draftService.DeleteDraft(id); err != nil {
		http.Error(w, "Draft not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *DraftHandler) Publish(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	teamID := auth.GetTeamID(ctx)
	userID := auth.GetUserID(ctx)

	var req struct {
		ScheduledFor string `json:"scheduledFor,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.draftService.PublishDraft(id, req.ScheduledFor, teamID, userID)
	if err != nil {
		log.Printf("Failed to publish draft: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writePublishResponse(w, result)
}

func (h *DraftHandler) PublishContent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	teamID := auth.GetTeamID(ctx)
	userID := auth.GetUserID(ctx)

	var req models.DraftDataPost
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.draftService.PublishContent(req, teamID, userID)
	if err != nil {
		log.Printf("Failed to publish content: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writePublishResponse(w, result)
}

func (h *DraftHandler) ScheduleDraft(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := auth.GetUserID(ctx)
	teamID := auth.GetTeamID(ctx)

	var req models.ScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	draftID, err := h.s.ScheduleDraft(req, teamID, userID)
	if err != nil {
		log.Printf("Failed to schedule draft: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Draft scheduled successfully",
		"draft_id": draftID,
		"status":   "scheduled",
	})
}

func (h *DraftHandler) CancelScheduledDraft(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DraftID string `json:"draft_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.s.CancelSchedule(req.DraftID); err != nil {
		log.Printf("Failed to cancel schedule: %v", err)
		http.Error(w, "Failed to cancel schedule", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Schedule cancelled",
		"draft_id": req.DraftID,
		"status":   "draft",
	})
}

// Helper methods
func (h *DraftHandler) writePublishResponse(w http.ResponseWriter, result *services.PublishResult) {
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

func validateCreateRequest(req models.CreateDraftRequest) error {
	if req.Title == "" || req.Topic == "" || req.Article == "" {
		return fmt.Errorf("title, topic, and article are required")
	}
	return nil
}
