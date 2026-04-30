package apikey

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"seo-backend/internal/database"
	"seo-backend/internal/middleware/auth"
)

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	teamID := auth.GetTeamID(ctx)
	userRole := auth.GetUserRole(ctx)
	userID := auth.GetUserID(ctx)

	query := `
		SELECT 
			ak.id, ak.service, ak.provider_id, ak.model_id,
			ak.is_active, ak.system_prompt, ak.created_by,
			ak.created_at, ak.updated_at,
			p.name, p.display_name,
			m.name, m.display_name
		FROM api_keys ak
		LEFT JOIN api_providers p ON ak.provider_id = p.id
		LEFT JOIN ai_models m ON ak.model_id = m.id
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	switch userRole {
	case "admin":
		if teamID != "" {
			query += fmt.Sprintf(" AND ak.team_id = $%d", argIndex)
			args = append(args, teamID)
			argIndex++
		} else {
			query += fmt.Sprintf(" AND ak.created_by = $%d", argIndex)
			args = append(args, userID)
			argIndex++
		}
	default:
		query += fmt.Sprintf(" AND ak.created_by = $%d", argIndex)
		args = append(args, userID)
		argIndex++
	}

	query += " ORDER BY ak.created_at DESC"

	rows, err := database.GetDB().Query(query, args...)
	if err != nil {
		log.Printf("query error: %v", err)
		http.Error(w, "failed fetch api keys", 500)
		return
	}
	defer rows.Close()

	var result []APIKeyDetail

	for rows.Next() {
		var item APIKeyDetail
		var sysPrompt, providerName, providerDisplay, modelName, modelDisplay sql.NullString

		err := rows.Scan(
			&item.ID, &item.Service, &item.ProviderID, &item.ModelID,
			&item.IsActive, &sysPrompt, &item.CreatedBy,
			&item.CreatedAt, &item.UpdatedAt,
			&providerName, &providerDisplay,
			&modelName, &modelDisplay,
		)
		if err != nil {
			continue
		}

		item.SystemPrompt = ScanString(sysPrompt)
		item.ProviderName = ScanString(providerName)
		item.ProviderDisplayName = ScanString(providerDisplay)
		item.ModelName = ScanString(modelName)
		item.ModelDisplayName = ScanString(modelDisplay)

		result = append(result, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
