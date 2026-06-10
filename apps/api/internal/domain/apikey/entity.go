package apikey

import (
	"time"
)

// APIKey entity - domain model
type APIKey struct {
	ID           string    `json:"id"`
	Service      string    `json:"service"`
	ProviderID   string    `json:"providerId"`
	ModelID      string    `json:"modelId"`
	KeyEncrypted string    `json:"-"` // tidak diekspos ke JSON
	IsActive     bool      `json:"isActive"`
	SystemPrompt string    `json:"systemPrompt"`
	TeamID       *string   `json:"teamId,omitempty"`
	CreatedBy    string    `json:"createdBy"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// APIKeyDetail - enriched entity with provider & model names
type APIKeyDetail struct {
	ID                  string    `json:"id"`
	Service             string    `json:"service"`
	ProviderID          string    `json:"provider_id"`
	ProviderName        string    `json:"provider_name"`
	ProviderDisplayName string    `json:"provider_display_name"`
	ModelID             string    `json:"model_id"`
	ModelName           string    `json:"model_name"`
	ModelDisplayName    string    `json:"model_display_name"`
	IsActive            bool      `json:"is_active"`
	SystemPrompt        string    `json:"system_prompt"`
	CreatedBy           *string   `json:"created_by"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// CreateAPIKeyRequest - DTO for creation
type CreateAPIKeyRequest struct {
	Service      string `json:"service"`
	ProviderID   string `json:"providerId"`
	ModelID      string `json:"modelId"`
	Key          string `json:"key"`
	SystemPrompt string `json:"systemPrompt"`
}

// UpdateAPIKeyRequest - DTO for update
type UpdateAPIKeyRequest struct {
	Service      *string `json:"service,omitempty"`
	ProviderID   *string `json:"providerId,omitempty"`
	ModelID      *string `json:"modelId,omitempty"`
	Key          *string `json:"key,omitempty"`
	IsActive     *bool   `json:"isActive,omitempty"`
	SystemPrompt *string `json:"systemPrompt,omitempty"`
}
