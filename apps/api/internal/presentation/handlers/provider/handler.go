package provider

import (
	"encoding/json"
	"net/http"

	"seo-backend/internal/domain/provider"
	"seo-backend/internal/middleware/auth"

	"github.com/go-chi/chi/v5"
)

// ProviderHandler handles HTTP requests for providers
type ProviderHandler struct {
	service provider.ProviderService
}

// NewProviderHandler creates a new provider handler
func NewProviderHandler(service provider.ProviderService) *ProviderHandler {
	return &ProviderHandler{service: service}
}

// userContextAdapter adapts auth middleware to service interface
type userContextAdapter struct {
	userID string
	teamID string
	role   string
}

func (u *userContextAdapter) GetUserID() string   { return u.userID }
func (u *userContextAdapter) GetTeamID() string   { return u.teamID }
func (u *userContextAdapter) GetUserRole() string { return u.role }

func (h *ProviderHandler) getUserContext(r *http.Request) *userContextAdapter {
	ctx := r.Context()
	return &userContextAdapter{
		userID: auth.GetUserID(ctx),
		teamID: auth.GetTeamID(ctx),
		role:   auth.GetUserRole(ctx),
	}
}

// Helper functions
func (h *ProviderHandler) writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *ProviderHandler) writeError(w http.ResponseWriter, message string, status int) {
	h.writeJSON(w, map[string]string{"error": message}, status)
}

// Create handles POST /providers
// @Summary Create a new API provider
// @Tags providers
// @Accept json
// @Produce json
// @Param request body provider.CreateProviderRequest true "Provider data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
func (h *ProviderHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := h.getUserContext(r)

	var req provider.CreateProviderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateProvider(ctx, req, userCtx)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "access denied: admin role required":
			status = http.StatusForbidden
		case "name is required", "displayName is required", "baseUrl is required",
			"authType is required", "authHeader is required", "textEndpoint is required",
			"responseTextPath is required":
			status = http.StatusBadRequest
		}
		h.writeError(w, err.Error(), status)
		return
	}

	h.writeJSON(w, map[string]string{
		"id":      id,
		"message": "Provider created successfully",
	}, http.StatusCreated)
}

// GetByID handles GET /providers/{id}
// @Summary Get provider by ID
// @Tags providers
// @Produce json
// @Param id path string true "Provider ID"
// @Success 200 {object} provider.APIProvider
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
func (h *ProviderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	provider, err := h.service.GetProviderByID(ctx, id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "provider not found" {
			status = http.StatusNotFound
		}
		h.writeError(w, err.Error(), status)
		return
	}

	h.writeJSON(w, provider, http.StatusOK)
}

// GetAll handles GET /providers
// @Summary Get all providers
// @Tags providers
// @Produce json
// @Success 200 {array} provider.APIProvider
// @Failure 500 {object} map[string]string
func (h *ProviderHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := h.getUserContext(r)

	providers, err := h.service.GetAllProviders(ctx, userCtx)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, providers, http.StatusOK)
}

// Update handles PUT /providers/{id}
// @Summary Update a provider
// @Tags providers
// @Accept json
// @Produce json
// @Param id path string true "Provider ID"
// @Param request body provider.UpdateProviderRequest true "Update data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
func (h *ProviderHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	userCtx := h.getUserContext(r)

	var req provider.UpdateProviderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateProvider(ctx, id, req, userCtx); err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "access denied: admin role required":
			status = http.StatusForbidden
		case "provider not found":
			status = http.StatusNotFound
		case "provider id is required":
			status = http.StatusBadRequest
		}
		h.writeError(w, err.Error(), status)
		return
	}

	h.writeJSON(w, map[string]string{
		"message": "Provider updated successfully",
	}, http.StatusOK)
}

// Delete handles DELETE /providers/{id}
// @Summary Delete a provider
// @Tags providers
// @Produce json
// @Param id path string true "Provider ID"
// @Success 204 "No Content"
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
func (h *ProviderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	userCtx := h.getUserContext(r)

	if err := h.service.DeleteProvider(ctx, id, userCtx); err != nil {
		status := http.StatusInternalServerError
		switch {
		case err.Error() == "access denied: admin role required":
			status = http.StatusForbidden
		case err.Error() == "provider not found":
			status = http.StatusNotFound
		case len(err.Error()) > 0 && err.Error()[:7] == "cannot delete provider because it is used":
			status = http.StatusConflict
		}
		h.writeError(w, err.Error(), status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RegisterRoutes registers all provider routes
func (h *ProviderHandler) RegisterRoutes(r chi.Router) {
	r.Route("/providers", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.GetAll)
		r.Get("/{id}", h.GetByID)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}
