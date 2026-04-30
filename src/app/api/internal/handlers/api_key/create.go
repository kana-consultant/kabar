package apikey

import (
	"encoding/json"
	"log"
	"net/http"

	"seo-backend/internal/builder"
	"seo-backend/internal/database"
	"seo-backend/internal/middleware/auth"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := auth.GetUserID(ctx)
	teamID := auth.GetTeamID(ctx)

	var req struct {
		Service      string `json:"service"`
		ProviderID   string `json:"providerId"`
		ModelID      string `json:"modelId"`
		Key          string `json:"key"`
		SystemPrompt string `json:"systemPrompt"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", 400)
		return
	}

	if req.Service == "" || req.ProviderID == "" || req.ModelID == "" || req.Key == "" {
		http.Error(w, "missing required fields", 400)
		return
	}

	encryptedKey, err := h.encrypt(req.Key)
	if err != nil {
		http.Error(w, "failed encrypt key", 500)
		return
	}

	var teamIDPtr any
	if teamID != "" {
		teamIDPtr = teamID
	}

	qb := builder.NewQueryBuilder("api_keys")

	sqlQuery, args, err := qb.Insert().
		Columns(
			"service",
			"provider_id",
			"model_id",
			"key_encrypted",
			"system_prompt",
			"team_id",
			"is_active",
			"created_by",
		).
		Values(
			req.Service,
			req.ProviderID,
			req.ModelID,
			encryptedKey,
			req.SystemPrompt,
			teamIDPtr,
			true,
			userID,
		).
		Returning("id").
		Build()

	if err != nil {
		log.Printf("build insert error: %v", err)
		http.Error(w, "failed create api key", 500)
		return
	}

	var id string
	err = database.GetDB().QueryRow(sqlQuery, args...).Scan(&id)
	if err != nil {
		log.Printf("insert error: %v", err)
		http.Error(w, "failed create api key", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"id":      id,
		"message": "created successfully",
	})
}
