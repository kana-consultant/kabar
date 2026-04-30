package apikey

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"seo-backend/internal/database"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var key APIKey
	var sysPrompt sql.NullString

	err := database.GetDB().QueryRow(`
		SELECT id, service, provider_id, model_id,
		       is_active, system_prompt, created_at, updated_at
		FROM api_keys WHERE id = $1
	`, id).Scan(
		&key.ID,
		&key.Service,
		&key.ProviderID,
		&key.ModelID,
		&key.IsActive,
		&sysPrompt,
		&key.CreatedAt,
		&key.UpdatedAt,
	)

	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	key.SystemPrompt = ScanString(sysPrompt)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key)
}
