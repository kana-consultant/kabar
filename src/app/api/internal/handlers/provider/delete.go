package provider

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"seo-backend/internal/builder"
	"seo-backend/internal/database"
)

// Delete removes an API provider (admin only)
func (h *ProviderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)

	if !h.isAdmin(userCtx.Role) {
		h.writeError(w, "Forbidden", http.StatusForbidden)
		return
	}

	id := chi.URLParam(r, "id")

	// Check if provider is used by any models
	var modelCount int
	err := database.GetDB().QueryRow(`
		SELECT COUNT(*) FROM ai_models WHERE provider_id = $1
	`, id).Scan(&modelCount)
	if err != nil {
		h.writeError(w, "Failed to check provider usage", http.StatusInternalServerError)
		return
	}

	if modelCount > 0 {
		h.writeError(w, "Cannot delete provider because it is used by existing models", http.StatusConflict)
		return
	}

	// Build delete query
	qb := builder.NewQueryBuilder("api_providers")
	sqlQuery, args, err := qb.Delete().WhereEq("id", id).Build()
	if err != nil {
		h.writeError(w, "Failed to delete provider", http.StatusInternalServerError)
		return
	}

	// Execute query
	_, err = database.GetDB().Exec(sqlQuery, args...)
	if err != nil {
		h.writeError(w, "Failed to delete provider", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
