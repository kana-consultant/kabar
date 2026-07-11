package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"seo-backend/common"
	"seo-backend/internal/domain/paginate"
	"seo-backend/internal/domain/provider"

	auth "seo-backend/internal/middleware"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// ProviderHandler handles HTTP requests for providers
type ProviderHandler struct {
	service provider.Service
}

// NewProviderHandler creates a new provider handler
func NewProviderHandler(service provider.Service) *ProviderHandler {
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	switch status {
	case http.StatusBadRequest:
		json.NewEncoder(w).Encode(common.ErrorResponse400{
			Error:  message,
			Status: status,
		})
	case http.StatusForbidden:
		json.NewEncoder(w).Encode(common.ErrorResponse403{
			Error:  message,
			Status: status,
		})
	case http.StatusNotFound:
		json.NewEncoder(w).Encode(common.ErrorResponse404{
			Error:  message,
			Status: status,
		})
	case http.StatusConflict:
		json.NewEncoder(w).Encode(common.ErrorResponse409{
			Error:  message,
			Status: status,
		})
	default:
		json.NewEncoder(w).Encode(common.ErrorResponse500{
			Error:  message,
			Status: status,
		})
	}
}

// Create godoc
// @Summary Create a new provider
// @Description Create a new AI provider with configuration
// @Tags Providers
// @Accept json
// @Produce json
// @Param request body provider.CreateRequest true "Provider creation request"
// @Success 201 {object} provider.Response "Provider created"
// @Failure 400 {object} common.ErrorResponse400 "Bad request - missing required fields"
// @Failure 403 {object} common.ErrorResponse403 "Access denied - admin role required"
// @Failure 409 {object} common.ErrorResponse409 "Provider with this name already exists"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/providers [post]
func (h *ProviderHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := auth.GetUserContext(r)

	// Log request received
	log.Printf("[ProviderHandler.Create] Received create provider request from user: %v", userCtx.GetUserID())

	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ProviderHandler.Create] Failed to read request body: %v", err)
		h.writeError(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Log raw request body
	log.Printf("[ProviderHandler.Create] Raw request body: %s", string(body))

	// Restore body for decoding
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	var req provider.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[ProviderHandler.Create] JSON decode error: %v", err)
		log.Printf("[ProviderHandler.Create] Raw body that caused error: %s", string(body))
		h.writeError(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Log decoded request
	log.Printf("[ProviderHandler.Create] Decoded request: %+v", req)

	// Validate required fields
	if req.Name == "" {
		log.Printf("[ProviderHandler.Create] Validation failed: name is required")
		h.writeError(w, "name is required", http.StatusBadRequest)
		return
	}
	if req.DisplayName == "" {
		log.Printf("[ProviderHandler.Create] Validation failed: display_name is required")
		h.writeError(w, "display_name is required", http.StatusBadRequest)
		return
	}
	if req.BaseURL == "" {
		log.Printf("[ProviderHandler.Create] Validation failed: base_url is required")
		h.writeError(w, "base_url is required", http.StatusBadRequest)
		return
	}

	resp, err := h.service.Create(ctx, &req, userCtx)
	if err != nil {
		log.Printf("[ProviderHandler.Create] Service error: %v", err)

		status := http.StatusInternalServerError
		switch err.Error() {
		case "access denied: admin role required":
			status = http.StatusForbidden
		case "name is required", "display_name is required", "base_url is required":
			status = http.StatusBadRequest
		case "provider with this name already exists":
			status = http.StatusConflict
		}
		h.writeError(w, err.Error(), status)
		return
	}

	log.Printf("[ProviderHandler.Create] Successfully created provider with ID: %v", resp.ID)
	h.writeJSON(w, resp, http.StatusCreated)
}

// GetByID godoc
// @Summary Get provider by ID
// @Description Get a specific provider by its ID
// @Tags Providers
// @Accept json
// @Produce json
// @Param id path string true "Provider ID"
// @Success 200 {object} provider.Response "Provider details"
// @Failure 404 {object} common.ErrorResponse404 "Provider not found"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/providers/{id} [get]
func (h *ProviderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	userCtx := auth.GetUserContext(r)

	provider, err := h.service.GetByID(ctx, id, userCtx)
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

// GetByName godoc
// @Summary Get provider by name
// @Description Get a specific provider by its name
// @Tags Providers
// @Accept json
// @Produce json
// @Param name path string true "Provider name"
// @Success 200 {object} provider.Response "Provider details"
// @Failure 404 {object} common.ErrorResponse404 "Provider not found"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/providers/name/{name} [get]
func (h *ProviderHandler) GetByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	name := chi.URLParam(r, "name")
	userCtx := auth.GetUserContext(r)

	provider, err := h.service.GetByName(ctx, name, userCtx)
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

// GetAll godoc
// @Summary Get all providers
// @Description Get all providers with pagination and search
// @Tags Providers
// @Accept json
// @Produce json
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Param search query string false "Search term"
// @Success 200 {object} paginate.PaginatedResult[provider.Response] "List of providers"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/providers [get]
func (h *ProviderHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := auth.GetUserContext(r)

	// Parse pagination params
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	search := r.URL.Query().Get("search")

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	params := paginate.PaginationParams{
		Limit:  limit,
		Offset: offset,
		Search: search,
	}

	result, err := h.service.GetAll(ctx, userCtx, params)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, result, http.StatusOK)
}

// GetActive godoc
// @Summary Get active providers
// @Description Get all active providers with pagination
// @Tags Providers
// @Accept json
// @Produce json
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Param search query string false "Search term"
// @Success 200 {object} paginate.PaginatedResult[provider.Response] "List of active providers"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/providers/active [get]
func (h *ProviderHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := auth.GetUserContext(r)

	// Parse pagination params
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	search := r.URL.Query().Get("search")

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	params := paginate.PaginationParams{
		Limit:  limit,
		Offset: offset,
		Search: search,
	}

	result, err := h.service.GetActive(ctx, userCtx, params)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, result, http.StatusOK)
}

// Update godoc
// @Summary Update a provider
// @Description Update an existing provider by ID (admin only)
// @Tags Providers
// @Accept json
// @Produce json
// @Param id path string true "Provider ID"
// @Param request body provider.UpdateRequest true "Provider update data"
// @Success 200 {object} provider.Response "Updated provider"
// @Failure 400 {object} common.ErrorResponse400 "Bad request - no fields to update"
// @Failure 403 {object} common.ErrorResponse403 "Access denied - admin role required"
// @Failure 404 {object} common.ErrorResponse404 "Provider not found"
// @Failure 409 {object} common.ErrorResponse409 "Provider with this name already exists"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/providers/{id} [put]
func (h *ProviderHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	userCtx := auth.GetUserContext(r)

	// Log request received
	log.Printf("[ProviderHandler.Update] Received update request for provider ID: %s from user: %v", id, userCtx.GetUserID())

	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ProviderHandler.Update] Failed to read request body: %v", err)
		h.writeError(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Log raw request body
	log.Printf("[ProviderHandler.Update] Raw request body: %s", string(body))

	// Restore body for decoding
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	var req provider.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[ProviderHandler.Update] JSON decode error: %v", err)
		log.Printf("[ProviderHandler.Update] Raw body that caused error: %s", string(body))
		h.writeError(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Log decoded request
	log.Printf("[ProviderHandler.Update] Decoded request: %+v", req)

	// Validate at least one field to update
	if req.Name == nil && req.DisplayName == nil && req.Description == nil &&
		req.BaseURL == nil && req.AuthType == nil && req.AuthHeader == nil &&
		req.AuthPrefix == nil && req.DefaultHeaders == nil && req.IsActive == nil {
		log.Printf("[ProviderHandler.Update] Validation failed: no fields to update")
		h.writeError(w, "No fields to update", http.StatusBadRequest)
		return
	}

	resp, err := h.service.Update(ctx, id, &req, userCtx)
	if err != nil {
		log.Printf("[ProviderHandler.Update] Service error: %v", err)

		status := http.StatusInternalServerError
		switch err.Error() {
		case "access denied: admin role required":
			status = http.StatusForbidden
		case "provider not found":
			status = http.StatusNotFound
		case "provider with this name already exists":
			status = http.StatusConflict
		}
		h.writeError(w, err.Error(), status)
		return
	}

	log.Printf("[ProviderHandler.Update] Successfully updated provider ID: %s", id)
	h.writeJSON(w, resp, http.StatusOK)
}

// Delete godoc
// @Summary Soft delete a provider
// @Description Soft delete a provider by ID (admin only)
// @Tags Providers
// @Accept json
// @Produce json
// @Param id path string true "Provider ID"
// @Success 204 "No Content"
// @Failure 403 {object} common.ErrorResponse403 "Access denied - admin role required"
// @Failure 404 {object} common.ErrorResponse404 "Provider not found"
// @Failure 409 {object} common.ErrorResponse409 "Cannot delete provider with existing references"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/providers/{id} [delete]
func (h *ProviderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	userCtx := auth.GetUserContext(r)

	if err := h.service.Delete(ctx, id, userCtx); err != nil {
		status := http.StatusInternalServerError
		switch {
		case err.Error() == "access denied: admin role required":
			status = http.StatusForbidden
		case err.Error() == "provider not found":
			status = http.StatusNotFound
		case len(err.Error()) > 0 && err.Error()[:7] == "cannot delete provider":
			status = http.StatusConflict
		}
		h.writeError(w, err.Error(), status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HardDelete godoc
// @Summary Hard delete a provider
// @Description Permanently delete a provider by ID (admin only)
// @Tags Providers
// @Accept json
// @Produce json
// @Param id path string true "Provider ID"
// @Success 204 "No Content"
// @Failure 403 {object} common.ErrorResponse403 "Access denied - admin role required"
// @Failure 404 {object} common.ErrorResponse404 "Provider not found"
// @Failure 409 {object} common.ErrorResponse409 "Cannot delete provider with existing references"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/providers/{id}/hard [delete]
func (h *ProviderHandler) HardDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	userCtx := auth.GetUserContext(r)

	if err := h.service.HardDelete(ctx, id, userCtx); err != nil {
		status := http.StatusInternalServerError
		switch {
		case err.Error() == "access denied: admin role required":
			status = http.StatusForbidden
		case err.Error() == "provider not found":
			status = http.StatusNotFound
		case len(err.Error()) > 0 && err.Error()[:7] == "cannot delete provider":
			status = http.StatusConflict
		}
		h.writeError(w, err.Error(), status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ToggleActive godoc
// @Summary Toggle provider active status
// @Description Toggle the active status of a provider
// @Tags Providers
// @Accept json
// @Produce json
// @Param id path string true "Provider ID"
// @Success 200 {object} common.SuccessMessage "Provider status toggled successfully"
// @Failure 404 {object} common.ErrorResponse404 "Provider not found"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/providers/{id}/toggle-active [patch]
func (h *ProviderHandler) ToggleActive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	userCtx := auth.GetUserContext(r)

	if err := h.service.ToggleActive(ctx, id, userCtx); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "provider not found" {
			status = http.StatusNotFound
		}
		h.writeError(w, err.Error(), status)
		return
	}

	h.writeJSON(w, common.SuccessMessage{Message: "Provider status toggled successfully"}, http.StatusOK)
}

// UpdateHeaders godoc
// @Summary Update provider headers
// @Description Update default headers for a provider
// @Tags Providers
// @Accept json
// @Produce json
// @Param id path string true "Provider ID"
// @Param request body object true "Headers to update" example({"X-API-Key": "new-key", "X-API-Version": "2024-01"})
// @Success 200 {object} common.SuccessMessage "Headers updated successfully"
// @Failure 400 {object} common.ErrorResponse400 "Invalid request body"
// @Failure 404 {object} common.ErrorResponse404 "Provider not found"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/providers/{id}/headers [patch]
func (h *ProviderHandler) UpdateHeaders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	userCtx := auth.GetUserContext(r)

	var headers map[string]string
	if err := json.NewDecoder(r.Body).Decode(&headers); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateHeaders(ctx, id, headers, userCtx); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "provider not found" {
			status = http.StatusNotFound
		}
		h.writeError(w, err.Error(), status)
		return
	}

	h.writeJSON(w, common.SuccessMessage{Message: "Headers updated successfully"}, http.StatusOK)
}
