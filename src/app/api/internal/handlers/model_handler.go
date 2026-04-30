package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"seo-backend/internal/database"
	"seo-backend/internal/middleware/auth"
)

type AIModelHandler struct{}

func NewAIModelHandler() *AIModelHandler {
	return &AIModelHandler{}
}

func (h *AIModelHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	rows, _ := database.GetDB().Query(`
        SELECT id, name, display_name, provider_id 
        FROM ai_models 
        WHERE is_active = true
    `)
	defer rows.Close()

	var models []map[string]interface{}
	for rows.Next() {
		var id, name, displayName, providerID string
		rows.Scan(&id, &name, &displayName, &providerID)
		models = append(models, map[string]interface{}{
			"id": id, "name": name, "displayName": displayName, "providerId": providerID,
		})
	}
	json.NewEncoder(w).Encode(models)
}

// GetAllWithStatus - mendapatkan semua model dengan status apakah sudah ada API key
func (h *AIModelHandler) GetAllWithStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userRole := auth.GetUserRole(ctx)

	query := `
		SELECT DISTINCT
			m.id, 
			m.name, 
			m.provider_id, 
			m.display_name
		FROM ai_models m
		INNER JOIN api_keys ak ON ak.model_id = m.id
		WHERE ak.is_active = true AND ak.service = 'text'
	`

	// Filter berdasarkan role - PAKE AND, BUKAN WHERE LAGI
	if userRole != "super_admin" {
		query += " AND m.is_active = true"
	}

	rows, err := database.GetDB().Query(query)
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Failed to fetch models", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type ModelWithStatus struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		ProviderID  string `json:"providerId"`
		DisplayName string `json:"displayName"`
	}

	var models []ModelWithStatus
	for rows.Next() {
		var m ModelWithStatus
		err := rows.Scan(&m.ID, &m.Name, &m.ProviderID, &m.DisplayName)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		models = append(models, m)
	}

	log.Printf("Found %d models", len(models))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models)
}
