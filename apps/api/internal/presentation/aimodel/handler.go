// internal/interfaces/api/aimodel/handler.go
package aimodel

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	aimodel "seo-backend/internal/domain/ai_model"
	"seo-backend/internal/domain/paginate"
	auth "seo-backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

type AIModelHandler struct {
	aimodelService aimodel.Service
}

func NewAIModelHandler(aimodelService aimodel.Service) *AIModelHandler {
	return &AIModelHandler{
		aimodelService: aimodelService,
	}
}

// GetAll returns all AI models with pagination
func (h *AIModelHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := auth.GetUserContext(r)

	// Parse pagination params
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	search := r.URL.Query().Get("search")

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if page <= 0 {
		page = 1
	}

	params := paginate.PaginationParams{
		Limit:  limit,
		Offset: page,
		Search: search,
	}

	result, err := h.aimodelService.GetAll(ctx, userCtx, params)
	if err != nil {
		log.Printf("Failed to fetch models: %v", err)
		h.writeError(w, "Failed to fetch models", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, result, http.StatusOK)
}

// GetAllWithStatus returns all models with API key status
func (h *AIModelHandler) GetAllWithStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := auth.GetUserContext(r)

	// Parse pagination params
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	search := r.URL.Query().Get("search")

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if page <= 0 {
		page = 1
	}

	params := paginate.PaginationParams{
		Limit:  limit,
		Offset: page,
		Search: search,
	}

	result, err := h.aimodelService.GetAllWithStatus(ctx, userCtx, params)
	if err != nil {
		log.Printf("Failed to fetch models with status: %v", err)
		h.writeError(w, "Failed to fetch models", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, result, http.StatusOK)
}

// GetByID returns a specific model by ID
func (h *AIModelHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	userCtx := auth.GetUserContext(r)

	model, err := h.aimodelService.GetByID(ctx, id, userCtx)
	if err != nil {
		log.Printf("Failed to fetch model: %v", err)
		if err == aimodel.ErrNotFound {
			h.writeError(w, "Model not found", http.StatusNotFound)
			return
		}
		h.writeError(w, "Failed to fetch model", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, model, http.StatusOK)
}

// GetSchemaByID returns schema for a specific model
func (h *AIModelHandler) GetSchemaByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	userCtx := auth.GetUserContext(r)

	result, err := h.aimodelService.GetSchemaByID(ctx, id, userCtx)
	if err != nil {
		log.Printf("Failed to fetch schema for model %s: %v", id, err)
		if err == aimodel.ErrNotFound {
			h.writeError(w, "Model or schema not found", http.StatusNotFound)
			return
		}
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, result, http.StatusOK)
}

// GetByProvider returns models by provider ID with pagination
func (h *AIModelHandler) GetByProvider(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	providerID := chi.URLParam(r, "providerId")
	userCtx := auth.GetUserContext(r)

	// Parse pagination params
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	search := r.URL.Query().Get("search")

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if page <= 0 {
		page = 1
	}

	params := paginate.PaginationParams{
		Limit:  limit,
		Offset: page,
		Search: search,
	}

	result, err := h.aimodelService.GetByProvider(ctx, providerID, userCtx, params)
	if err != nil {
		log.Printf("Failed to fetch models by provider: %v", err)
		h.writeError(w, "Failed to fetch models", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, result, http.StatusOK)
}

// GetByFamily returns models by family ID with pagination
func (h *AIModelHandler) GetByFamily(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	familyID := chi.URLParam(r, "familyId")
	userCtx := auth.GetUserContext(r)

	// Parse pagination params
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	search := r.URL.Query().Get("search")

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if page <= 0 {
		page = 1
	}

	params := paginate.PaginationParams{
		Limit:  limit,
		Offset: page,
		Search: search,
	}

	result, err := h.aimodelService.GetByFamily(ctx, familyID, userCtx, params)
	if err != nil {
		log.Printf("Failed to fetch models by family: %v", err)
		h.writeError(w, "Failed to fetch models", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, result, http.StatusOK)
}

// GetByTeam returns models by team ID with pagination
func (h *AIModelHandler) GetByTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	teamID := chi.URLParam(r, "teamId")
	userCtx := auth.GetUserContext(r)

	// Parse pagination params
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	search := r.URL.Query().Get("search")

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if page <= 0 {
		page = 1
	}

	params := paginate.PaginationParams{
		Limit:  limit,
		Offset: page,
		Search: search,
	}

	result, err := h.aimodelService.GetByTeam(ctx, teamID, userCtx, params)
	if err != nil {
		log.Printf("Failed to fetch models by team: %v", err)
		h.writeError(w, "Failed to fetch models", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, result, http.StatusOK)
}

// GetDefault returns default models
func (h *AIModelHandler) GetDefault(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := auth.GetUserContext(r)

	models, err := h.aimodelService.GetDefault(ctx, userCtx)
	if err != nil {
		log.Printf("Failed to fetch default models: %v", err)
		h.writeError(w, "Failed to fetch default models", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, models, http.StatusOK)
}

// Create creates a new AI model
func (h *AIModelHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req aimodel.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userCtx := auth.GetUserContext(r)

	resp, err := h.aimodelService.Create(ctx, &req, userCtx)
	if err != nil {
		log.Printf("Failed to create model: %v", err)
		if err == aimodel.ErrDuplicate {
			h.writeError(w, "Model with this name already exists", http.StatusConflict)
			return
		}
		h.writeError(w, "Failed to create model", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, resp, http.StatusCreated)
}

// Update updates an existing AI model
func (h *AIModelHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	var req aimodel.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userCtx := auth.GetUserContext(r)

	resp, err := h.aimodelService.Update(ctx, id, &req, userCtx)
	if err != nil {
		log.Printf("Failed to update model: %v", err)
		if err == aimodel.ErrNotFound {
			h.writeError(w, "Model not found", http.StatusNotFound)
			return
		}
		if err == aimodel.ErrDuplicate {
			h.writeError(w, "Model with this name already exists", http.StatusConflict)
			return
		}
		h.writeError(w, "Failed to update model", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, resp, http.StatusOK)
}

// Delete deletes an AI model
func (h *AIModelHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	userCtx := auth.GetUserContext(r)

	err := h.aimodelService.Delete(ctx, id, userCtx)
	if err != nil {
		log.Printf("Failed to delete model: %v", err)
		if err == aimodel.ErrNotFound {
			h.writeError(w, "Model not found", http.StatusNotFound)
			return
		}
		h.writeError(w, "Failed to delete model", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SetAsDefault sets a model as default for its provider
func (h *AIModelHandler) SetAsDefault(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	userCtx := auth.GetUserContext(r)

	err := h.aimodelService.SetAsDefault(ctx, id, userCtx)
	if err != nil {
		log.Printf("Failed to set model as default: %v", err)
		if err == aimodel.ErrNotFound {
			h.writeError(w, "Model not found", http.StatusNotFound)
			return
		}
		h.writeError(w, "Failed to set model as default", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, map[string]string{"message": "Model set as default successfully"}, http.StatusOK)
}

// Helper methods
func (h *AIModelHandler) writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

func (h *AIModelHandler) writeError(w http.ResponseWriter, message string, status int) {
	h.writeJSON(w, map[string]string{"error": message}, status)
}
