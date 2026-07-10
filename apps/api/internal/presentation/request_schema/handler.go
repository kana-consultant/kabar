package request_schema

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"seo-backend/internal/domain/request_schema"

	"github.com/go-chi/chi/v5"
)

type RequestSchemaHandler struct {
	service request_schema.Service
}

func NewRequestSchemaHandler(service request_schema.Service) *RequestSchemaHandler {
	return &RequestSchemaHandler{
		service: service,
	}
}

// Create godoc
// @Summary Create a new request schema
// @Description Create a new request schema for API validation
// @Tags Request Schemas
// @Accept json
// @Produce json
// @Param request body request_schema.CreateRequest true "Request schema creation request"
// @Success 201 {object} request_schema.RequestSchema "Request schema created"
// @Failure 400 {object} map[string]string "Bad request - invalid provider_id, name, or endpoint_path"
// @Failure 409 {object} map[string]string "Duplicate entry"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/request-schemas [post]
func (h *RequestSchemaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req request_schema.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.service.Create(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, request_schema.ErrDuplicate):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, request_schema.ErrInvalidProviderID):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, request_schema.ErrInvalidName):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, request_schema.ErrInvalidEndpointPath):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetByID godoc
// @Summary Get request schema by ID
// @Description Get a specific request schema by its ID
// @Tags Request Schemas
// @Accept json
// @Produce json
// @Param id path string true "Request Schema ID"
// @Success 200 {object} request_schema.RequestSchema "Request schema details"
// @Failure 404 {object} map[string]string "Request schema not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/request-schemas/{id} [get]
func (h *RequestSchemaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	resp, err := h.service.GetByID(r.Context(), idStr)
	if err != nil {
		switch {
		case errors.Is(err, request_schema.ErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetAll godoc
// @Summary Get all request schemas
// @Description Get all request schemas with pagination
// @Tags Request Schemas
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Success 200 {object} map[string]interface{} "List of request schemas with pagination"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/request-schemas [get]
func (h *RequestSchemaHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	resp, err := h.service.GetAll(r.Context(), page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetByProvider godoc
// @Summary Get request schemas by provider
// @Description Get all request schemas for a specific provider
// @Tags Request Schemas
// @Accept json
// @Produce json
// @Param provider_id path string true "Provider ID"
// @Success 200 {array} request_schema.RequestSchema "List of request schemas"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/request-schemas/provider/{provider_id} [get]
func (h *RequestSchemaHandler) GetByProvider(w http.ResponseWriter, r *http.Request) {
	providerIDStr := chi.URLParam(r, "provider_id")

	resp, err := h.service.GetByProvider(r.Context(), providerIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Update godoc
// @Summary Update a request schema
// @Description Update an existing request schema by ID
// @Tags Request Schemas
// @Accept json
// @Produce json
// @Param id path string true "Request Schema ID"
// @Param request body request_schema.UpdateRequest true "Request schema update request"
// @Success 200 {object} request_schema.RequestSchema "Updated request schema"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 404 {object} map[string]string "Request schema not found"
// @Failure 409 {object} map[string]string "Duplicate entry"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/request-schemas/{id} [put]
func (h *RequestSchemaHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	var req request_schema.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.service.Update(r.Context(), idStr, &req)
	if err != nil {
		switch {
		case errors.Is(err, request_schema.ErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, request_schema.ErrDuplicate):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Delete godoc
// @Summary Delete a request schema
// @Description Delete a request schema by ID
// @Tags Request Schemas
// @Accept json
// @Produce json
// @Param id path string true "Request Schema ID"
// @Success 200 {object} map[string]string "message: deleted successfully"
// @Failure 404 {object} map[string]string "Request schema not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/request-schemas/{id} [delete]
func (h *RequestSchemaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	if err := h.service.Delete(r.Context(), idStr); err != nil {
		switch {
		case errors.Is(err, request_schema.ErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "deleted successfully"})
}
