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

// Create
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

// GetByID
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

// GetByProviderAndName
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

// GetAll
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

// GetByProvider
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

// GetBySchema
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

// Update
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

// Delete
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

// CheckExists
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

// // GetStatistics
// func (h *ModelFamilyHandler) GetStatistics(w http.ResponseWriter, r *http.Context) {
// 	ctx := r.Context()

// 	// This method is not in Service interface, need to implement or remove
// 	// For now, return not implemented
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"message": "Statistics endpoint not implemented",
// 	})
// }
