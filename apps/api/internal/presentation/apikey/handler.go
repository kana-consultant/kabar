package apikey

import (
	"encoding/json"
	"net/http"

	"seo-backend/internal/domain/apikey"
	auth "seo-backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

type APIKeyHandler struct {
	service apikey.Service
}

func NewAPIKeyHandler(service apikey.Service) *APIKeyHandler {
	return &APIKeyHandler{service: service}
}

// UserContextAdapter - adapts auth middleware to service interface
type userContextAdapter struct {
	userID string
	teamID string
	role   string
}

func (u *userContextAdapter) GetUserID() string   { return u.userID }
func (u *userContextAdapter) GetTeamID() string   { return u.teamID }
func (u *userContextAdapter) GetUserRole() string { return u.role }

// Helper functions
func (h *APIKeyHandler) writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *APIKeyHandler) writeError(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

// Create godoc
// @Summary Create a new API key
// @Description Create a new API key for accessing AI models
// @Tags API Keys
// @Accept json
// @Produce json
// @Param request body apikey.CreateAPIKeyRequest true "API key creation request"
// @Success 201 {object} map[string]string "id: key_id, message: created successfully"
// @Failure 400 {object} map[string]string "Bad request - missing required fields (service, provider_id, model_id, api_key)"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/api-keys [post]
func (h *APIKeyHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := auth.GetUserContext(r)

	var req apikey.CreateAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateAPIKey(ctx, req, userCtx)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "service is required" ||
			err.Error() == "provider id is required" ||
			err.Error() == "model id is required" ||
			err.Error() == "API key is required" {
			status = http.StatusBadRequest
		}
		h.writeError(w, err, status)
		return
	}

	h.writeJSON(w, map[string]string{
		"id":      id,
		"message": "created successfully",
	}, http.StatusCreated)
}

// GetByID godoc
// @Summary Get API key by ID
// @Description Get a specific API key by its ID
// @Tags API Keys
// @Accept json
// @Produce json
// @Param id path string true "API Key ID"
// @Success 200 {object} apikey.APIKey "Success"
// @Failure 404 {object} map[string]string "API key not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/api-keys/{id} [get]
func (h *APIKeyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	key, err := h.service.GetAPIKeyByID(ctx, id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "API key not found" {
			status = http.StatusNotFound
		}
		h.writeError(w, err, status)
		return
	}

	h.writeJSON(w, key, http.StatusOK)
}

// GetAll godoc
// @Summary Get all API keys
// @Description Get all API keys for the current user/team
// @Tags API Keys
// @Accept json
// @Produce json
// @Success 200 {array} apikey.APIKey "Success"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/api-keys [get]
func (h *APIKeyHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := auth.GetUserContext(r)

	keys, err := h.service.GetAllAPIKeys(ctx, userCtx)
	if err != nil {
		h.writeError(w, err, http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, keys, http.StatusOK)
}

// Update godoc
// @Summary Update an API key
// @Description Update an existing API key by ID
// @Tags API Keys
// @Accept json
// @Produce json
// @Param id path string true "API Key ID"
// @Param request body apikey.UpdateAPIKeyRequest true "API key update request"
// @Success 200 {object} map[string]string "id: key_id, message: updated successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 403 {object} map[string]string "Access denied"
// @Failure 404 {object} map[string]string "API key not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/api-keys/{id} [put]
func (h *APIKeyHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	userCtx := auth.GetUserContext(r)

	var req apikey.UpdateAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateAPIKey(ctx, id, req, userCtx); err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "API key not found":
			status = http.StatusNotFound
		case "access denied":
			status = http.StatusForbidden
		}
		h.writeError(w, err, status)
		return
	}

	h.writeJSON(w, map[string]string{
		"id":      id,
		"message": "updated successfully",
	}, http.StatusOK)
}

// Delete godoc
// @Summary Delete an API key
// @Description Delete an API key by ID
// @Tags API Keys
// @Accept json
// @Produce json
// @Param id path string true "API Key ID"
// @Success 200 {object} map[string]string "id: key_id, message: deleted successfully"
// @Failure 403 {object} map[string]string "Access denied"
// @Failure 404 {object} map[string]string "API key not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/api-keys/{id} [delete]
func (h *APIKeyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	userCtx := auth.GetUserContext(r)

	if err := h.service.DeleteAPIKey(ctx, id, userCtx); err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "API key not found":
			status = http.StatusNotFound
		case "access denied":
			status = http.StatusForbidden
		}
		h.writeError(w, err, status)
		return
	}

	h.writeJSON(w, map[string]string{
		"id":      id,
		"message": "deleted successfully",
	}, http.StatusOK)
}
