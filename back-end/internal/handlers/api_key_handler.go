package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"

	"seo-backend/internal/builder"
	"seo-backend/internal/database"
	"seo-backend/internal/middleware/auth"
)

type APIKeyHandler struct{}

func NewAPIKeyHandler() *APIKeyHandler {
	return &APIKeyHandler{}
}

// GetAll API Keys
func (h *APIKeyHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	teamID := auth.GetTeamID(ctx)
	userRole := auth.GetUserRole(ctx)
	userID := auth.GetUserID(ctx)

	query := `
		SELECT 
			ak.id, 
			ak.service, 
			ak.provider_id, 
			ak.model_id, 
			ak.is_active, 
			ak.system_prompt, 
			ak.created_by, 
			ak.created_at, 
			ak.updated_at,
			p.name as provider_name,
			p.display_name as provider_display_name,
			m.name as model_name,
			m.display_name as model_display_name
		FROM api_keys ak
		LEFT JOIN api_providers p ON ak.provider_id = p.id
		LEFT JOIN ai_models m ON ak.model_id = m.id
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	switch userRole {
	case "super_admin":
		// no filter
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
		log.Printf("Query error: %v", err)
		http.Error(w, "Failed to fetch API keys", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type APIKeyDetail struct {
		ID                  string    `json:"id"`
		Service             string    `json:"service"`
		ProviderID          string    `json:"providerId"`
		ModelID             string    `json:"modelId"`
		IsActive            bool      `json:"isActive"`
		SystemPrompt        string    `json:"systemPrompt"`
		CreatedBy           *string   `json:"createdBy"`
		CreatedAt           time.Time `json:"createdAt"`
		UpdatedAt           time.Time `json:"updatedAt"`
		ProviderName        string    `json:"providerName"`
		ProviderDisplayName string    `json:"providerDisplayName"`
		ModelName           string    `json:"modelName"`
		ModelDisplayName    string    `json:"modelDisplayName"`
	}

	var result []APIKeyDetail
	for rows.Next() {
		var item APIKeyDetail
		var sysPrompt sql.NullString
		var providerName, providerDisplayName, modelName, modelDisplayName sql.NullString

		err := rows.Scan(
			&item.ID, &item.Service, &item.ProviderID, &item.ModelID,
			&item.IsActive, &sysPrompt,
			&item.CreatedBy, &item.CreatedAt, &item.UpdatedAt,
			&providerName, &providerDisplayName,
			&modelName, &modelDisplayName,
		)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}

		if sysPrompt.Valid {
			item.SystemPrompt = sysPrompt.String
		}
		if providerName.Valid {
			item.ProviderName = providerName.String
			item.ProviderDisplayName = providerDisplayName.String
		}
		if modelName.Valid {
			item.ModelName = modelName.String
			item.ModelDisplayName = modelDisplayName.String
		}

		result = append(result, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetByID API Key
func (h *APIKeyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var key struct {
		ID           string    `json:"id"`
		Service      string    `json:"service"`
		ProviderID   string    `json:"providerId"`
		ModelID      string    `json:"modelId"`
		IsActive     bool      `json:"isActive"`
		SystemPrompt string    `json:"systemPrompt"`
		CreatedAt    time.Time `json:"createdAt"`
		UpdatedAt    time.Time `json:"updatedAt"`
	}
	var sysPrompt sql.NullString

	err := database.GetDB().QueryRow(`
		SELECT id, service, provider_id, model_id, is_active, system_prompt, created_at, updated_at
		FROM api_keys WHERE id = $1
	`, id).Scan(&key.ID, &key.Service, &key.ProviderID, &key.ModelID, &key.IsActive, &sysPrompt, &key.CreatedAt, &key.UpdatedAt)
	if err != nil {
		http.Error(w, "API key not found", http.StatusNotFound)
		return
	}

	key.SystemPrompt = sysPrompt.String

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key)
}

func (h *APIKeyHandler) encrypt(plainText string) (string, error) {
	key := os.Getenv("ENCRYPTION_KEY")
	if key == "" {
		key = "01234567890123456789012345678901"
	}

	block, _ := aes.NewCipher([]byte(key))
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Create API Key
func (h *APIKeyHandler) Create(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Service == "" || req.ProviderID == "" || req.ModelID == "" || req.Key == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	var teamIDPtr interface{}
	if teamID != "" {
		teamIDPtr = teamID
	} else {
		teamIDPtr = nil
	}

	qb := builder.NewQueryBuilder("api_keys")

	// Di Create, encrypt dulu sebelum insert
	encryptedKey, err := h.encrypt(req.Key)
	if err != nil {
		http.Error(w, "Failed to encrypt key", http.StatusInternalServerError)
		return
	}
	sqlQuery, args, err := qb.Insert().
		Columns("service", "provider_id", "model_id", "key_encrypted", "system_prompt", "team_id", "is_active", "created_by").
		Values(req.Service, req.ProviderID, req.ModelID, encryptedKey, req.SystemPrompt, teamIDPtr, true, userID).
		Returning("id").
		Build()

	if err != nil {
		log.Printf("Failed to build insert query: %v", err)
		http.Error(w, "Failed to create API key", http.StatusInternalServerError)
		return
	}

	var keyID string
	err = database.GetDB().QueryRow(sqlQuery, args...).Scan(&keyID)
	if err != nil {
		log.Printf("Failed to create API key: %v", err)
		http.Error(w, "Failed to create API key", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id":      keyID,
		"message": "API key created successfully",
	})
}

// Update API Key
func (h *APIKeyHandler) Update(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	qb := builder.NewQueryBuilder("api_keys")
	updateBuilder := qb.Update()

	if req.Service != nil {
		updateBuilder = updateBuilder.Set("service", *req.Service)
	}
	if req.ProviderID != nil {
		updateBuilder = updateBuilder.Set("provider_id", *req.ProviderID)
	}
	if req.ModelID != nil {
		updateBuilder = updateBuilder.Set("model_id", *req.ModelID)
	}
	if req.Key != nil && *req.Key != "" {
		encryptedKey, err := h.encrypt(*req.Key)
		if err != nil {
			http.Error(w, "Failed to encrypt key", http.StatusInternalServerError)
			return
		}
		updateBuilder = updateBuilder.Set("key_encrypted", encryptedKey)
	}
	if req.SystemPrompt != nil {
		updateBuilder = updateBuilder.Set("system_prompt", *req.SystemPrompt)
	}
	if req.IsActive != nil {
		updateBuilder = updateBuilder.Set("is_active", *req.IsActive)
	}

	updateBuilder = updateBuilder.Set("updated_at", time.Now())
	updateBuilder = updateBuilder.WhereEq("id", id)

	sqlQuery, args, err := updateBuilder.Build()
	if err != nil {
		http.Error(w, "Failed to update API key", http.StatusInternalServerError)
		return
	}

	_, err = database.GetDB().Exec(sqlQuery, args...)
	if err != nil {
		log.Printf("Failed to update API key: %v", err)
		http.Error(w, "Failed to update API key", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "API key updated successfully"})
}

// Delete API Key
func (h *APIKeyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	qb := builder.NewQueryBuilder("api_keys")
	sqlQuery, args, err := qb.Delete().WhereEq("id", id).Build()
	if err != nil {
		http.Error(w, "Failed to delete API key", http.StatusInternalServerError)
		return
	}

	_, err = database.GetDB().Exec(sqlQuery, args...)
	if err != nil {
		log.Printf("Failed to delete API key: %v", err)
		http.Error(w, "Failed to delete API key", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
