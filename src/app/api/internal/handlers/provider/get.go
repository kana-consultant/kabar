package provider

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"seo-backend/internal/builder"
	"seo-backend/internal/database"
)

// GetByID returns a specific API provider by ID
func (h *ProviderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Build query
	qb := builder.NewQueryBuilder("api_providers")
	sqlQuery, args, err := qb.Select(
		"id", "name", "display_name", "description",
		"base_url", "auth_type", "auth_header", "auth_prefix",
		"text_endpoint", "image_endpoint",
		"default_headers", "request_template",
		"response_text_path", "response_image_path",
		"is_active", "created_at", "updated_at",
	).WhereEq("id", id).Build()

	if err != nil {
		h.writeError(w, "Failed to fetch provider", http.StatusInternalServerError)
		return
	}

	// Execute query
	row := database.GetDB().QueryRow(sqlQuery, args...)
	provider, err := h.scanProvider(row)
	if err != nil {
		h.writeError(w, "Provider not found", http.StatusNotFound)
		return
	}

	h.writeJSON(w, provider, http.StatusOK)
}
