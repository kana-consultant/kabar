package ai_model

import (
	"time"
)

// AIModel represents the ai_models table entity
type AIModel struct {
	ID         string  `json:"id"`
	FamilyID   *string `json:"family_id"`
	ProviderID *string `json:"provider_id"`
	// SchemaID      *string   `json:"schema_id"` // <-- DIHAPUS
	TeamID        *string   `json:"team_id"`
	Name          string    `json:"name"`
	DisplayName   *string   `json:"display_name"`
	Description   *string   `json:"description"`
	SystemPrompt  *string   `json:"system_prompt"`
	MaxTokens     *int      `json:"max_tokens"`
	Temperature   *float64  `json:"temperature"`
	ContextWindow *int      `json:"context_window"`
	IsActive      *bool     `json:"is_active"`
	IsDefault     *bool     `json:"is_default"`
	CreatedBy     *string   `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// NewAIModel creates a new AIModel entity
func NewAIModel(
	familyID *string,
	providerID *string,
	// schemaID *string, // <-- DIHAPUS
	teamID *string,
	name *string,
	displayName *string,
	description *string,
	systemPrompt *string,
	maxTokens *int,
	temperature *float64,
	contextWindow *int,
	isActive *bool,
	isDefault *bool,
	createdBy *string,
) *AIModel {
	now := time.Now()
	return &AIModel{
		FamilyID:   familyID,
		ProviderID: providerID,
		// SchemaID:      schemaID, // <-- DIHAPUS
		TeamID:        teamID,
		Name:          *name,
		DisplayName:   displayName,
		Description:   description,
		SystemPrompt:  systemPrompt,
		MaxTokens:     maxTokens,
		Temperature:   temperature,
		ContextWindow: contextWindow,
		IsActive:      isActive,
		IsDefault:     isDefault,
		CreatedBy:     createdBy,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

type AiModelSchema struct {
	ID         string  `db:"id"`
	FamilyID   *string `db:"family_id"`
	ProviderID string  `db:"provider_id"`
	// SchemaID      *string       `db:"schema_id"` // <-- DIHAPUS
	TeamID *string `db:"team_id"`
	// Name          string        `db:"name"` // <-- DIHAPUS
	DisplayName   string        `db:"display_name"`
	Description   *string       `db:"description"`
	SystemPrompt  *string       `db:"system_prompt"`
	MaxTokens     *int          `db:"max_tokens"`
	Temperature   *float64      `db:"temperature"`
	ContextWindow *int          `db:"context_window"`
	Schema        RequestSchema `db:"schema"`
	IsActive      bool          `db:"is_active"`
	IsDefault     bool          `db:"is_default"`
	CreatedBy     *string       `db:"created_by"`
	CreatedAt     time.Time     `db:"created_at"`
	UpdatedAt     time.Time     `db:"updated_at"`
}

type RequestSchema struct {
	ID         *string `db:"id"`
	ProviderID *string `db:"provider_id"`
	// Name                *string    `db:"name"` // <-- DIHAPUS
	EndpointPath        *string    `db:"endpoint_path"`
	RequestTemplate     *string    `db:"request_template"`
	ResponseTextPath    *string    `db:"response_text_path"`
	ResponseImagePath   *string    `db:"response_image_path"`
	SupportsTemperature *bool      `db:"supports_temperature"`
	SupportsStreaming   *bool      `db:"supports_streaming"`
	CreatedAt           *time.Time `db:"created_at"`
	UpdatedAt           *time.Time `db:"updated_at"`
}

type ResponseSchema struct {
	ID         string  `json:"id"`
	FamilyID   *string `json:"family_id,omitempty"`
	ProviderID string  `json:"provider_id"`
	// SchemaID      *string       `json:"schema_id,omitempty"` // <-- DIHAPUS
	TeamID        *string       `json:"team_id,omitempty"`
	Name          string        `json:"name"`
	DisplayName   string        `json:"display_name"`
	Description   *string       `json:"description,omitempty"`
	SystemPrompt  *string       `json:"system_prompt,omitempty"`
	MaxTokens     *int          `json:"max_tokens,omitempty"`
	Temperature   *float64      `json:"temperature,omitempty"`
	ContextWindow *int          `json:"context_window,omitempty"`
	Schema        RequestSchema `json:"schema"`
	IsActive      bool          `json:"is_active"`
	IsDefault     bool          `json:"is_default"`
	CreatedBy     *string       `json:"created_by,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}
