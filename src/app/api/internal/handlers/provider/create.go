package provider

import (
	"encoding/json"
	"fmt"
	"net/http"

	"seo-backend/internal/builder"
	"seo-backend/internal/database"
)

// Create creates a new API provider (admin only)
func (h *ProviderHandler) Create(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)

	// Only admin and super_admin can create providers
	if !h.isAdmin(userCtx.Role) {
		h.writeError(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req CreateProviderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if err := h.validateCreateRequest(req); err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Marshal JSON fields
	defaultHeadersJSON, _ := json.Marshal(req.DefaultHeaders)
	requestTemplateJSON, _ := json.Marshal(req.RequestTemplate)

	// Build insert query
	qb := builder.NewQueryBuilder("api_providers")
	sqlQuery, args, err := qb.Insert().
		Columns(
			"name", "display_name", "description", "base_url",
			"auth_type", "auth_header", "auth_prefix",
			"text_endpoint", "image_endpoint",
			"default_headers", "request_template",
			"response_text_path", "response_image_path",
			"is_active",
		).
		Values(
			req.Name, req.DisplayName, req.Description, req.BaseURL,
			req.AuthType, req.AuthHeader, req.AuthPrefix,
			req.TextEndpoint, nullIfEmpty(req.ImageEndpoint),
			defaultHeadersJSON, requestTemplateJSON,
			req.ResponseTextPath, nullIfEmpty(req.ResponseImagePath),
			true,
		).
		Returning("id").
		Build()

	if err != nil {
		h.writeError(w, "Failed to create provider", http.StatusInternalServerError)
		return
	}

	// Execute query
	var providerID string
	err = database.GetDB().QueryRow(sqlQuery, args...).Scan(&providerID)
	if err != nil {
		h.writeError(w, "Failed to create provider", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, map[string]string{
		"id":      providerID,
		"message": "Provider created successfully",
	}, http.StatusCreated)
}

// Helper: Validate create request
func (h *ProviderHandler) validateCreateRequest(req CreateProviderRequest) error {
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if req.DisplayName == "" {
		return fmt.Errorf("displayName is required")
	}
	if req.BaseURL == "" {
		return fmt.Errorf("baseUrl is required")
	}
	if req.AuthType == "" {
		return fmt.Errorf("authType is required")
	}
	if req.AuthHeader == "" {
		return fmt.Errorf("authHeader is required")
	}
	if req.TextEndpoint == "" {
		return fmt.Errorf("textEndpoint is required")
	}
	if req.ResponseTextPath == "" {
		return fmt.Errorf("responseTextPath is required")
	}
	return nil
}
