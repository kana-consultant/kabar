// internal/handlers/provider_handler.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"seo-backend/internal/builder"
	"seo-backend/internal/database"
	"seo-backend/internal/middleware/auth"
)

type ProviderHandler struct{}

func NewProviderHandler() *ProviderHandler {
	return &ProviderHandler{}
}

type APIProvider struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	DisplayName       string                 `json:"displayName"`
	Description       string                 `json:"description"`
	BaseURL           string                 `json:"baseUrl"`
	AuthType          string                 `json:"authType"`
	AuthHeader        string                 `json:"authHeader"`
	AuthPrefix        sql.NullString         `json:"authPrefix"`
	TextEndpoint      string                 `json:"textEndpoint"`
	ImageEndpoint     *string                `json:"imageEndpoint,omitempty"`
	DefaultHeaders    map[string]string      `json:"defaultHeaders"`
	RequestTemplate   map[string]interface{} `json:"requestTemplate"`
	ResponseTextPath  string                 `json:"responseTextPath"`
	ResponseImagePath *string                `json:"responseImagePath,omitempty"`
	IsActive          bool                   `json:"isActive"`
	CreatedAt         time.Time              `json:"createdAt"`
	UpdatedAt         time.Time              `json:"updatedAt"`
}

type CreateProviderRequest struct {
	Name              string                 `json:"name" binding:"required"`
	DisplayName       string                 `json:"displayName" binding:"required"`
	Description       string                 `json:"description"`
	BaseURL           string                 `json:"baseUrl" binding:"required"`
	AuthType          string                 `json:"authType" binding:"required"`
	AuthHeader        string                 `json:"authHeader" binding:"required"`
	AuthPrefix        string                 `json:"authPrefix"`
	TextEndpoint      string                 `json:"textEndpoint" binding:"required"`
	ImageEndpoint     string                 `json:"imageEndpoint"`
	DefaultHeaders    map[string]string      `json:"defaultHeaders"`
	RequestTemplate   map[string]interface{} `json:"requestTemplate" binding:"required"`
	ResponseTextPath  string                 `json:"responseTextPath" binding:"required"`
	ResponseImagePath string                 `json:"responseImagePath"`
}

// GetAll providers
func (h *ProviderHandler) GetAll(
	w http.ResponseWriter,
	r *http.Request,
) {

	ctx := r.Context()

	userRole := auth.GetUserRole(ctx)
	teamID := auth.GetTeamID(ctx)

	log.Println("==== GET ALL PROVIDERS DEBUG ====")
	log.Println("userRole:", userRole)
	log.Println("teamID:", teamID)

	qb := builder.NewQueryBuilder("api_providers")

	query := qb.Select(
		"id",
		"name",
		"display_name",
		"description",
		"base_url",
		"auth_type",
		"auth_header",
		"auth_prefix",
		"text_endpoint",
		"image_endpoint",
		"default_headers",
		"request_template",
		"response_text_path",
		"response_image_path",
		"is_active",
		"created_at",
		"updated_at",
	).OrderBy("name ASC")

	switch userRole {

	case "super_admin", "admin":
		log.Println("role can access all providers")

	case "manager", "editor", "viewer":
		log.Println("filtering active providers only")
		query = query.Where("is_active = true")

	default:
		log.Println("unknown role fallback active=true")
		query = query.Where("is_active = true")
	}

	if teamID != "" &&
		(userRole == "admin" || userRole == "super_admin") {

		log.Println("teamID exists:", teamID)
	}

	sqlQuery, args, err := query.Build()

	if err != nil {
		log.Println("query build error:", err)
		http.Error(
			w,
			"Failed to fetch providers",
			http.StatusInternalServerError,
		)
		return
	}

	log.Println("SQL:", sqlQuery)
	log.Printf("ARGS: %#v\n", args)

	rows, err := database.GetDB().Query(
		sqlQuery,
		args...,
	)

	if err != nil {
		log.Println("query execute error:", err)
		http.Error(
			w,
			"Failed to fetch providers",
			http.StatusInternalServerError,
		)
		return
	}

	defer rows.Close()

	var providers []APIProvider
	rowCount := 0

	for rows.Next() {

		rowCount++

		var p APIProvider
		var defaultHeadersJSON []byte
		var requestTemplateJSON []byte
		var imageEndpoint *string
		var responseImagePath *string

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.DisplayName,
			&p.Description,
			&p.BaseURL,
			&p.AuthType,
			&p.AuthHeader,
			&p.AuthPrefix,
			&p.TextEndpoint,
			&imageEndpoint,
			&defaultHeadersJSON,
			&requestTemplateJSON,
			&p.ResponseTextPath,
			&responseImagePath,
			&p.IsActive,
			&p.CreatedAt,
			&p.UpdatedAt,
		)

		if err != nil {
			log.Println(
				"scan error row",
				rowCount,
				":",
				err,
			)
			continue
		}

		log.Println("provider loaded:", p.ID, p.Name)

		if len(defaultHeadersJSON) > 0 {

			log.Println(
				"default_headers raw:",
				string(defaultHeadersJSON),
			)

			if err := json.Unmarshal(
				defaultHeadersJSON,
				&p.DefaultHeaders,
			); err != nil {

				log.Println(
					"default headers unmarshal error:",
					err,
				)
			}
		}

		if len(requestTemplateJSON) > 0 {

			log.Println(
				"request_template raw:",
				string(requestTemplateJSON),
			)

			if err := json.Unmarshal(
				requestTemplateJSON,
				&p.RequestTemplate,
			); err != nil {

				log.Println(
					"request template unmarshal error:",
					err,
				)
			}
		}

		p.ImageEndpoint = imageEndpoint
		p.ResponseImagePath = responseImagePath

		providers = append(
			providers,
			p,
		)
	}

	if err := rows.Err(); err != nil {
		log.Println("rows iteration error:", err)
	}

	log.Println("total rows:", rowCount)
	log.Println("providers returned:", len(providers))
	log.Println("==== END DEBUG ====")

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	if err := json.NewEncoder(w).Encode(providers); err != nil {
		log.Println("response encode error:", err)
	}
}

// GetByID provider
func (h *ProviderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

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
		http.Error(w, "Failed to fetch provider", http.StatusInternalServerError)
		return
	}

	var p APIProvider
	var defaultHeadersJSON, requestTemplateJSON []byte
	var imageEndpoint, responseImagePath *string

	err = database.GetDB().QueryRow(sqlQuery, args...).Scan(
		&p.ID, &p.Name, &p.DisplayName, &p.Description,
		&p.BaseURL, &p.AuthType, &p.AuthHeader, &p.AuthPrefix,
		&p.TextEndpoint, &imageEndpoint,
		&defaultHeadersJSON, &requestTemplateJSON,
		&p.ResponseTextPath, &responseImagePath,
		&p.IsActive, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		http.Error(w, "Provider not found", http.StatusNotFound)
		return
	}

	if len(defaultHeadersJSON) > 0 {
		json.Unmarshal(defaultHeadersJSON, &p.DefaultHeaders)
	}
	if len(requestTemplateJSON) > 0 {
		json.Unmarshal(requestTemplateJSON, &p.RequestTemplate)
	}

	p.ImageEndpoint = imageEndpoint
	p.ResponseImagePath = responseImagePath

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// Create provider
func (h *ProviderHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userRole := auth.GetUserRole(ctx)

	// Only admin and super_admin can create providers
	if userRole != "admin" && userRole != "super_admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req CreateProviderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" || req.DisplayName == "" || req.BaseURL == "" ||
		req.AuthType == "" || req.AuthHeader == "" || req.TextEndpoint == "" ||
		req.ResponseTextPath == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Marshal JSON fields
	defaultHeadersJSON, _ := json.Marshal(req.DefaultHeaders)
	requestTemplateJSON, _ := json.Marshal(req.RequestTemplate)

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
		http.Error(w, "Failed to create provider", http.StatusInternalServerError)
		return
	}

	var providerID string
	err = database.GetDB().QueryRow(sqlQuery, args...).Scan(&providerID)
	if err != nil {
		http.Error(w, "Failed to create provider", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id":      providerID,
		"message": "Provider created successfully",
	})
}

// Update provider
func (h *ProviderHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userRole := auth.GetUserRole(ctx)

	if userRole != "admin" && userRole != "super_admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	id := chi.URLParam(r, "id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

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
		http.Error(w, "Failed to update provider", http.StatusInternalServerError)
		return
	}

	_, err = database.GetDB().Exec(sqlQuery, args...)
	if err != nil {
		http.Error(w, "Failed to update provider", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Provider updated successfully"})
}

// Delete provider
func (h *ProviderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userRole := auth.GetUserRole(ctx)

	if userRole != "admin" && userRole != "super_admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	id := chi.URLParam(r, "id")

	// Check if provider is used by any models
	var modelCount int
	err := database.GetDB().QueryRow(`
		SELECT COUNT(*) FROM ai_models WHERE provider_id = $1
	`, id).Scan(&modelCount)
	if err != nil {
		http.Error(w, "Failed to check provider usage", http.StatusInternalServerError)
		return
	}

	if modelCount > 0 {
		http.Error(w, "Cannot delete provider because it is used by existing models", http.StatusConflict)
		return
	}

	qb := builder.NewQueryBuilder("api_providers")
	sqlQuery, args, err := qb.Delete().WhereEq("id", id).Build()
	if err != nil {
		http.Error(w, "Failed to delete provider", http.StatusInternalServerError)
		return
	}

	_, err = database.GetDB().Exec(sqlQuery, args...)
	if err != nil {
		http.Error(w, "Failed to delete provider", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper function
func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
