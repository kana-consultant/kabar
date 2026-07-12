package modelfamily

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	model_family "seo-backend/internal/domain/modelfamily"
	"seo-backend/internal/domain/paginate"
	auth "seo-backend/internal/middleware"
	"seo-backend/internal/models"

	"github.com/go-chi/chi/v5"
)

type ModelFamilyHandler struct {
	service model_family.Service
}

func NewModelFamilyHandler(service model_family.Service) *ModelFamilyHandler {
	return &ModelFamilyHandler{
		service: service,
	}
}

// Create godoc
// @Summary Create a new model family
// @Description Create a new AI model family with provider and schema
// @Tags Model Families
// @Accept json
// @Produce json
// @Param request body model_family.CreateRequest true "Model family creation request"
// @Success 201 {object} model_family.ModelFamily "Model family created"
// @Failure 400 {object} map[string]string "Bad request - missing required fields"
// @Failure 409 {object} map[string]string "Duplicate entry"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /families [post]
func (h *ModelFamilyHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model_family.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.ProviderID == "" {
		http.Error(w, model_family.ErrInvalidProviderID.Error(), http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, model_family.ErrInvalidName.Error(), http.StatusBadRequest)
		return
	}
	if req.SchemaID == "" {
		http.Error(w, model_family.ErrInvalidSchemaID.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.service.Create(ctx, &req)
	if err != nil {
		if errors.Is(err, model_family.ErrDuplicate) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetByID godoc
// @Summary Get model family by ID
// @Description Get a specific model family by its ID
// @Tags Model Families
// @Accept json
// @Produce json
// @Param id path string true "Model Family ID"
// @Success 200 {object} model_family.ModelFamily "Model family details"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Model family not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /families/{id} [get]
func (h *ModelFamilyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, model_family.ErrInvalidID.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.service.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, model_family.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetByProviderAndName godoc
// @Summary Get model family by provider and name
// @Description Get a model family by provider ID and name
// @Tags Model Families
// @Accept json
// @Produce json
// @Param provider_id query string true "Provider ID"
// @Param name query string true "Family name"
// @Success 200 {object} model_family.ModelFamily "Model family details"
// @Failure 400 {object} map[string]string "Missing required parameters"
// @Failure 404 {object} map[string]string "Model family not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /families/by-provider-name [get]
func (h *ModelFamilyHandler) GetByProviderAndName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	providerID := r.URL.Query().Get("provider_id")
	name := r.URL.Query().Get("name")

	if providerID == "" {
		http.Error(w, model_family.ErrInvalidProviderID.Error(), http.StatusBadRequest)
		return
	}
	if name == "" {
		http.Error(w, model_family.ErrInvalidName.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.service.GetByProviderAndName(ctx, providerID, name)
	if err != nil {
		if errors.Is(err, model_family.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetAll godoc
// @Summary Get all model families
// @Description Get all model families with pagination and search
// @Tags Model Families
// @Accept json
// @Produce json
// @Param limit query int false "Items per page (default: 10)"
// @Param page query int false "Page number (default: 1)"
// @Param search query string false "Search term"
// @Success 200 {object} paginate.PaginatedResult[model_family.Response] "List of model families"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /families [get]
func (h *ModelFamilyHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	search := r.URL.Query().Get("search")

	userCtx, ok := ctx.Value("user_context").(models.UserContext)
	if !ok {
		userCtx = auth.GetUserContext(r)
	}

	params := paginate.PaginationParams{
		Limit:  limit,
		Offset: page,
		Search: search,
	}

	result, err := h.service.GetAll(ctx, userCtx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetByProvider godoc
// @Summary Get model families by provider
// @Description Get all model families for a specific provider
// @Tags Model Families
// @Accept json
// @Produce json
// @Param provider_id path string true "Provider ID"
// @Success 200 {array} model_family.ModelFamily "List of model families"
// @Failure 400 {object} map[string]string "Invalid provider ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /families/providers/{provider_id}/families [get]
func (h *ModelFamilyHandler) GetByProvider(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	providerID := chi.URLParam(r, "provider_id")

	if providerID == "" {
		http.Error(w, model_family.ErrInvalidProviderID.Error(), http.StatusBadRequest)
		return
	}

	families, err := h.service.GetByProvider(ctx, providerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(families)
}

// GetBySchema godoc
// @Summary Get model families by schema
// @Description Get all model families for a specific schema
// @Tags Model Families
// @Accept json
// @Produce json
// @Param schema_id path string true "Schema ID"
// @Success 200 {array} model_family.ModelFamily "List of model families"
// @Failure 400 {object} map[string]string "Invalid schema ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /families/schemas/{schema_id}/families [get]
func (h *ModelFamilyHandler) GetBySchema(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	schemaID := chi.URLParam(r, "schema_id")

	if schemaID == "" {
		http.Error(w, model_family.ErrInvalidSchemaID.Error(), http.StatusBadRequest)
		return
	}

	families, err := h.service.GetBySchema(ctx, schemaID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(families)
}

// Update godoc
// @Summary Update a model family
// @Description Update an existing model family by ID
// @Tags Model Families
// @Accept json
// @Produce json
// @Param id path string true "Model Family ID"
// @Param request body model_family.UpdateRequest true "Update request"
// @Success 200 {object} model_family.ModelFamily "Updated model family"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "Model family not found"
// @Failure 409 {object} map[string]string "Duplicate entry"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /families/{id} [put]
func (h *ModelFamilyHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, model_family.ErrInvalidID.Error(), http.StatusBadRequest)
		return
	}

	var req model_family.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.service.Update(ctx, id, &req)
	if err != nil {
		if errors.Is(err, model_family.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, model_family.ErrDuplicate) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Delete godoc
// @Summary Delete a model family
// @Description Delete a model family by ID
// @Tags Model Families
// @Accept json
// @Produce json
// @Param id path string true "Model Family ID"
// @Success 200 {object} map[string]string "message: deleted successfully"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Model family not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /families/{id} [delete]
func (h *ModelFamilyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, model_family.ErrInvalidID.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		if errors.Is(err, model_family.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "deleted successfully"})
}

// CheckExists godoc
// @Summary Check if model family exists
// @Description Check if a model family exists for a specific provider and name
// @Tags Model Families
// @Accept json
// @Produce json
// @Param provider_id query string true "Provider ID"
// @Param name query string true "Family name"
// @Success 200 {object} map[string]bool "exists: true/false"
// @Failure 400 {object} map[string]string "Missing required parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /families/exists [get]
func (h *ModelFamilyHandler) CheckExists(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	providerID := r.URL.Query().Get("provider_id")
	name := r.URL.Query().Get("name")

	if providerID == "" {
		http.Error(w, model_family.ErrInvalidProviderID.Error(), http.StatusBadRequest)
		return
	}
	if name == "" {
		http.Error(w, model_family.ErrInvalidName.Error(), http.StatusBadRequest)
		return
	}

	// Validate method from service
	err := h.service.Validate(ctx, &model_family.ModelFamily{
		ProviderID: providerID,
		Name:       name,
	})

	exists := err == nil
	if err != nil && !errors.Is(err, model_family.ErrDuplicate) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"exists": exists})
}
