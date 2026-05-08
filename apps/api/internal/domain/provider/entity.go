package provider

import (
	"database/sql"
	"time"
)

// APIProvider represents an API provider configuration
type APIProvider struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	DisplayName       string                 `json:"displayName"`
	Description       string                 `json:"description"`
	BaseURL           string                 `json:"baseUrl"`
	AuthType          string                 `json:"authType"`
	AuthHeader        string                 `json:"authHeader"`
	AuthPrefix        sql.NullString         `json:"authPrefix,omitempty"`
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
	Name              string                 `json:"name"`
	DisplayName       string                 `json:"displayName"`
	Description       string                 `json:"description"`
	BaseURL           string                 `json:"baseUrl"`
	AuthType          string                 `json:"authType"`
	AuthHeader        string                 `json:"authHeader"`
	AuthPrefix        string                 `json:"authPrefix"`
	TextEndpoint      string                 `json:"textEndpoint"`
	ImageEndpoint     string                 `json:"imageEndpoint"`
	DefaultHeaders    map[string]string      `json:"defaultHeaders"`
	RequestTemplate   map[string]interface{} `json:"requestTemplate"`
	ResponseTextPath  string                 `json:"responseTextPath"`
	ResponseImagePath string                 `json:"responseImagePath"`
}

// UpdateProviderRequest represents request to update a provider
type UpdateProviderRequest struct {
	Name              *string                `json:"name,omitempty"`
	DisplayName       *string                `json:"displayName,omitempty"`
	Description       *string                `json:"description,omitempty"`
	BaseURL           *string                `json:"baseUrl,omitempty"`
	AuthType          *string                `json:"authType,omitempty"`
	AuthHeader        *string                `json:"authHeader,omitempty"`
	AuthPrefix        *string                `json:"authPrefix,omitempty"`
	TextEndpoint      *string                `json:"textEndpoint,omitempty"`
	ImageEndpoint     *string                `json:"imageEndpoint,omitempty"`
	DefaultHeaders    map[string]string      `json:"defaultHeaders,omitempty"`
	RequestTemplate   map[string]interface{} `json:"requestTemplate,omitempty"`
	ResponseTextPath  *string                `json:"responseTextPath,omitempty"`
	ResponseImagePath *string                `json:"responseImagePath,omitempty"`
	IsActive          *bool                  `json:"isActive,omitempty"`
}

// ProviderFilters for filtering providers
type ProviderFilters struct {
	IsActive *bool
	Search   string
}

// UserContext for permission checking
type UserContext interface {
	GetUserID() string
	GetTeamID() string
	GetUserRole() string
}
