package apikey

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"seo-backend/internal/builder"
	"seo-backend/internal/database"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req struct {
		Service      *string `json:"service"`
		ProviderID   *string `json:"providerId"`
		ModelID      *string `json:"modelId"`
		Key          *string `json:"key"`
		SystemPrompt *string `json:"systemPrompt"`
		IsActive     *bool   `json:"isActive"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", 400)
		return
	}

	qb := builder.NewQueryBuilder("api_keys")
	update := qb.Update()

	if req.Service != nil {
		update = update.Set("service", *req.Service)
	}
	if req.ProviderID != nil {
		update = update.Set("provider_id", *req.ProviderID)
	}
	if req.ModelID != nil {
		update = update.Set("model_id", *req.ModelID)
	}
	if req.Key != nil && *req.Key != "" {
		encrypted, err := h.encrypt(*req.Key)
		if err != nil {
			http.Error(w, "encrypt failed", 500)
			return
		}
		update = update.Set("key_encrypted", encrypted)
	}
	if req.SystemPrompt != nil {
		update = update.Set("system_prompt", *req.SystemPrompt)
	}
	if req.IsActive != nil {
		update = update.Set("is_active", *req.IsActive)
	}

	update = update.Set("updated_at", time.Now())
	update = update.WhereEq("id", id)

	sqlQuery, args, err := update.Build()
	if err != nil {
		http.Error(w, "build update failed", 500)
		return
	}

	_, err = database.GetDB().Exec(sqlQuery, args...)
	if err != nil {
		log.Printf("update error: %v", err)
		http.Error(w, "update failed", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "updated successfully",
	})
}
