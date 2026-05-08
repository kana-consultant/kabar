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
	ProviderID          string    `json:"providerId"`
	ProviderName        string    `json:"providerName"`
	ProviderDisplayName string    `json:"providerDisplayName"`
	ModelID             string    `json:"modelId"`
	ModelName           string    `json:"modelName"`
	ModelDisplayName    string    `json:"modelDisplayName"`
	IsActive            bool      `json:"isActive"`
	SystemPrompt        string    `json:"systemPrompt"`
	CreatedBy           *string   `json:"createdBy"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
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
