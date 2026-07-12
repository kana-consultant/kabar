// internal/interfaces/api/aimodel/handler.go
package aimodel

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"seo-backend/common"
	aimodel "seo-backend/internal/domain/ai_model"
	"seo-backend/internal/domain/paginate"
	auth "seo-backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// AIModelHandler handles AI model related HTTP requests
type AIModelHandler struct {
	aimodelService aimodel.Service
}

// NewAIModelHandler creates a new AIModelHandler instance
func NewAIModelHandler(aimodelService aimodel.Service) *AIModelHandler {
	return &AIModelHandler{
		aimodelService: aimodelService,
	}
}

// GetAll godoc
// @Summary Get all AI models
// @Description Get all AI models with pagination and search
// @Tags AI Models
// @Accept json
// @Produce json
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Param page query int false "Page number (default: 1)"
// @Param search query string false "Search term for model name"
// @Success 200 {object} paginate.PaginatedResult[aimodel.Response] "Success"
// @Failure 400 {object} common.ErrorResponse400 "Bad request"
// @Failure 401 {object} common.ErrorResponse401 "Unauthorized"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/ai-models [get]
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
		Offset: (page - 1) * limit,
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

// GetAllWithStatus godoc
// @Summary Get all AI models with API key status
// @Description Get all AI models with their API key status information
// @Tags AI Models
// @Accept json
// @Produce json
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Param page query int false "Page number (default: 1)"
// @Param search query string false "Search term for model name"
// @Success 200 {object} paginate.PaginatedResult[aimodel.ModelWithStatus] "Success"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/ai-models/with-status [get]
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
		Offset: (page - 1) * limit,
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

// GetByID godoc
// @Summary Get AI model by ID
// @Description Get a specific AI model by its ID
// @Tags AI Models
// @Accept json
// @Produce json
// @Param id path string true "Model ID"
// @Success 200 {object} aimodel.Response "Success"
// @Failure 404 {object} common.ErrorResponse404 "Model not found"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/ai-models/{id} [get]
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

// GetSchemaByID godoc
// @Summary Get model schema by ID
// @Description Get the schema definition for a specific AI model
// @Tags AI Models
// @Accept json
// @Produce json
// @Param id path string true "Model ID"
// @Success 200 {object} model_family.ModelFamilyWithSchema "Success"
// @Failure 404 {object} common.ErrorResponse404 "Model or schema not found"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/ai-models/{id}/schema [get]
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

// GetByProvider godoc
// @Summary Get AI models by provider
// @Description Get AI models filtered by provider ID with pagination
// @Tags AI Models
// @Accept json
// @Produce json
// @Param providerId path string true "Provider ID"
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Param page query int false "Page number (default: 1)"
// @Param search query string false "Search term for model name"
// @Success 200 {object} paginate.PaginatedResult[aimodel.Response] "Success"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/ai-models/provider/{providerId} [get]
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
		Offset: (page - 1) * limit,
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

// GetByFamily godoc
// @Summary Get AI models by family
// @Description Get AI models filtered by family ID with pagination
// @Tags AI Models
// @Accept json
// @Produce json
// @Param familyId path string true "Family ID"
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Param page query int false "Page number (default: 1)"
// @Param search query string false "Search term for model name"
// @Success 200 {object} paginate.PaginatedResult[aimodel.Response] "Success"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/ai-models/family/{familyId} [get]
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
		Offset: (page - 1) * limit,
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

// GetByTeam godoc
// @Summary Get AI models by team
// @Description Get AI models filtered by team ID with pagination
// @Tags AI Models
// @Accept json
// @Produce json
// @Param teamId path string true "Team ID"
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Param page query int false "Page number (default: 1)"
// @Param search query string false "Search term for model name"
// @Success 200 {object} paginate.PaginatedResult[aimodel.Response] "Success"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/ai-models/team/{teamId} [get]
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
		Offset: (page - 1) * limit,
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

// GetDefault godoc
// @Summary Get default AI models
// @Description Get all default AI models for each provider
// @Tags AI Models
// @Accept json
// @Produce json
// @Success 200 {array} aimodel.Response "Success"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/ai-models/default [get]
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

// Create godoc
// @Summary Create a new AI model
// @Description Create a new AI model with the provided details
// @Tags AI Models
// @Accept json
// @Produce json
// @Param request body aimodel.CreateRequest true "Model creation request"
// @Success 201 {object} aimodel.Response "Created"
// @Failure 400 {object} common.ErrorResponse400 "Invalid request body"
// @Failure 409 {object} common.ErrorResponse409 "Model with this name already exists"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/ai-models [post]
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

// Update godoc
// @Summary Update an AI model
// @Description Update an existing AI model by ID
// @Tags AI Models
// @Accept json
// @Produce json
// @Param id path string true "Model ID"
// @Param request body aimodel.UpdateRequest true "Model update request"
// @Success 200 {object} aimodel.Response "Success"
// @Failure 400 {object} common.ErrorResponse400 "Invalid request body"
// @Failure 404 {object} common.ErrorResponse404 "Model not found"
// @Failure 409 {object} common.ErrorResponse409 "Model with this name already exists"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/ai-models/{id} [put]
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

// Delete godoc
// @Summary Delete an AI model
// @Description Delete an AI model by ID
// @Tags AI Models
// @Accept json
// @Produce json
// @Param id path string true "Model ID"
// @Success 204 "No Content"
// @Failure 404 {object} common.ErrorResponse404 "Model not found"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/ai-models/{id} [delete]
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

// SetAsDefault godoc
// @Summary Set model as default
// @Description Set a specific AI model as the default for its provider
// @Tags AI Models
// @Accept json
// @Produce json
// @Param id path string true "Model ID"
// @Success 200 {object} map[string]string "Success message"
// @Failure 404 {object} common.ErrorResponse404 "Model not found"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/ai-models/{id}/default [put]
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(common.ErrorResponse500{
		Error:  message,
		Status: status,
	})
}
