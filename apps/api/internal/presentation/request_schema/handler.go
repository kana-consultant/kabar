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

// Create
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

// GetByID
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

// GetAll
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

// GetByProvider
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

// Update
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

// Delete
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
