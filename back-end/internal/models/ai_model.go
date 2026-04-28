package models

import "time"

type AIModel struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`        // gpt-4, claude-3, gemini-2.0-flash
	Provider    string    `json:"provider"`    // openai, anthropic, google, local
	DisplayName string    `json:"displayName"` // "GPT-4 Turbo", "Claude 3 Opus"
	Description string    `json:"description"`
	IsActive    bool      `json:"isActive"`
	IsDefault   bool      `json:"isDefault"`
	MaxTokens   int       `json:"maxTokens"`
	Temperature float64   `json:"temperature"`
	TeamID      *string   `json:"teamId,omitempty"`
	CreatedBy   *string   `json:"createdBy,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreateModelRequest struct {
	Name        string  `json:"name" binding:"required"`
	Provider    string  `json:"provider" binding:"required"`
	DisplayName string  `json:"displayName" binding:"required"`
	Description string  `json:"description"`
	IsDefault   bool    `json:"isDefault"`
	MaxTokens   int     `json:"maxTokens"`
	Temperature float64 `json:"temperature"`
}
