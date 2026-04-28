// internal/models/api_key.go
package models

import (
	"time"
)

// internal/models/api_key.go
type APIKey struct {
	ID           string    `json:"id"`
	Service      string    `json:"service"`    // 'text', 'image'
	ProviderID   string    `json:"providerId"` // FK to api_providers
	Model        string    `json:"model"`      // 'gemini-2.5-flash', 'gpt-4'
	KeyEncrypted string    `json:"-"`
	SystemPrompt string    `json:"systemPrompt,omitempty"`
	TeamID       *string   `json:"teamId,omitempty"`
	Config       string    `json:"config,omitempty"`
	IsActive     bool      `json:"isActive"`
	CreatedBy    *string   `json:"createdBy,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// internal/models/api_key.go
type CreateAPIKeyRequest struct {
	Service      string `json:"service" binding:"required"`
	ProviderID   string `json:"providerId" binding:"required"`
	Model        string `json:"model" binding:"required"`
	Key          string `json:"key" binding:"required"`
	SystemPrompt string `json:"systemPrompt,omitempty"`
	TeamID       string `json:"teamId,omitempty"`
}

type UpdateAPIKeyRequest struct {
	Service      *string `json:"service,omitempty"`
	ProviderID   *string `json:"providerId,omitempty"`
	Model        *string `json:"model,omitempty"`
	Key          *string `json:"key,omitempty"`
	SystemPrompt *string `json:"systemPrompt,omitempty"`
	IsActive     *bool   `json:"isActive,omitempty"`
}
