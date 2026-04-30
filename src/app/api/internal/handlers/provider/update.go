package provider

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"seo-backend/internal/builder"
	"seo-backend/internal/database"
)

// Update updates an existing API provider (admin only)
func (h *ProviderHandler) Update(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)

	if !h.isAdmin(userCtx.Role) {
		h.writeError(w, "Forbidden", http.StatusForbidden)
		return
	}

	id := chi.URLParam(r, "id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		h.writeError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Field mapping
	fieldMap := map[string]string{
		"name":              "name",
		"displayName":       "display_name",
		"description":       "description",
		"baseUrl":           "base_url",
		"authType":          "auth_type",
		"authHeader":        "auth_header",
		"authPrefix":        "auth_prefix",
		"textEndpoint":      "text_endpoint",
		"imageEndpoint":     "image_endpoint",
		"defaultHeaders":    "default_headers",
		"requestTemplate":   "request_template",
		"responseTextPath":  "response_text_path",
		"responseImagePath": "response_image_path",
		"isActive":          "is_active",
	}

	// Build update query
	qb := builder.NewQueryBuilder("api_providers")
	updateBuilder := qb.Update()

	for key, value := range updates {
		if dbField, ok := fieldMap[key]; ok && value != nil {
			// Handle JSON fields
			if key == "defaultHeaders" || key == "requestTemplate" {
				jsonValue, _ := json.Marshal(value)
				updateBuilder = updateBuilder.Set(dbField, jsonValue)
			} else {
				updateBuilder = updateBuilder.Set(dbField, value)
			}
		}
	}

	updateBuilder = updateBuilder.Set("updated_at", time.Now())
	updateBuilder = updateBuilder.WhereEq("id", id)

	sqlQuery, args, err := updateBuilder.Build()
	if err != nil {
		h.writeError(w, "Failed to update provider", http.StatusInternalServerError)
		return
	}

	// Execute query
	_, err = database.GetDB().Exec(sqlQuery, args...)
	if err != nil {
		h.writeError(w, "Failed to update provider", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, map[string]string{"message": "Provider updated successfully"}, http.StatusOK)
}
