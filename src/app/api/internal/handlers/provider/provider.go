package provider

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"seo-backend/internal/middleware/auth"
)

type ProviderHandler struct{}

func NewProviderHandler() *ProviderHandler {
	return &ProviderHandler{}
}

// APIProvider represents an API provider configuration
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

// CreateProviderRequest represents request to create a provider
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

// ProviderFilters for filtering providers
type ProviderFilters struct {
	IsActive *bool
	Search   string
}

// Helper: Get user context from request
func (h *ProviderHandler) getUserContext(r *http.Request) struct {
	Role   string
	TeamID string
} {
	ctx := r.Context()
	return struct {
		Role   string
		TeamID string
	}{
		Role:   auth.GetUserRole(ctx),
		TeamID: auth.GetTeamID(ctx),
	}
}

// Helper: Write JSON response
func (h *ProviderHandler) writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

// Helper: Write error response
func (h *ProviderHandler) writeError(w http.ResponseWriter, message string, status int) {
	h.writeJSON(w, map[string]string{"error": message}, status)
}

// Helper: Check if user is admin
func (h *ProviderHandler) isAdmin(role string) bool {
	return role == "admin" || role == "super_admin"
}

// Helper: Build where clause for role-based access
func (h *ProviderHandler) buildWhereClause(role string) string {
	switch role {
	case "super_admin", "admin":
		return ""
	case "manager", "editor", "viewer":
		return "WHERE is_active = true"
	default:
		return "WHERE is_active = true"
	}
}

// Helper: Scan provider from row
func (h *ProviderHandler) scanProvider(scanner interface {
	Scan(dest ...interface{}) error
}) (*APIProvider, error) {
	var p APIProvider
	var defaultHeadersJSON []byte
	var requestTemplateJSON []byte
	var imageEndpoint *string
	var responseImagePath *string

	err := scanner.Scan(
		&p.ID, &p.Name, &p.DisplayName, &p.Description,
		&p.BaseURL, &p.AuthType, &p.AuthHeader, &p.AuthPrefix,
		&p.TextEndpoint, &imageEndpoint,
		&defaultHeadersJSON, &requestTemplateJSON,
		&p.ResponseTextPath, &responseImagePath,
		&p.IsActive, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON fields
	if len(defaultHeadersJSON) > 0 {
		if err := json.Unmarshal(defaultHeadersJSON, &p.DefaultHeaders); err != nil {
			log.Printf("Failed to unmarshal default headers: %v", err)
			p.DefaultHeaders = make(map[string]string)
		}
	} else {
		p.DefaultHeaders = make(map[string]string)
	}

	if len(requestTemplateJSON) > 0 {
		if err := json.Unmarshal(requestTemplateJSON, &p.RequestTemplate); err != nil {
			log.Printf("Failed to unmarshal request template: %v", err)
			p.RequestTemplate = make(map[string]interface{})
		}
	} else {
		p.RequestTemplate = make(map[string]interface{})
	}

	p.ImageEndpoint = imageEndpoint
	p.ResponseImagePath = responseImagePath

	return &p, nil
}

// Helper: Null check for empty string
func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
