package model_family

import (
	"seo-backend/internal/domain/request_schema"
	"time"
)

// ModelFamily represents the model_families table entity
type ModelFamily struct {
	ID           string    `json:"id"`
	ProviderID   string    `json:"provider_id"`
	SchemaID     string    `json:"schema_id"`
	Name         string    `json:"name"`
	DisplayName  string    `json:"display_name"`
	Description  *string   `json:"description"`
	MaxTokens    int       `json:"max_token"`
	Temperature  float64   `json:"temperature"`
	SystemPrompt string    `json:"system_prompt"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ModelFamilyWithSchema struct {
	ID           string                       `json:"id"`
	ProviderID   string                       `json:"provider_id"`
	SchemaID     string                       `json:"schema_id"`
	Name         string                       `json:"name"`
	DisplayName  string                       `json:"display_name"`
	Description  *string                      `json:"description"`
	MaxTokens    int                          `json:"max_token"`
	Temperature  float64                      `json:"temperature"`
	SystemPrompt string                       `json:"system_prompt"`
	CreatedAt    time.Time                    `json:"created_at"`
	UpdatedAt    time.Time                    `json:"updated_at"`
	Schema       request_schema.RequestSchema `json:"schema"`
}

type ModelFamilyWithProvider struct {
	ModelFamily
	ProviderName        string `json:"provider_name"`
	ProviderDisplayName string `json:"provider_display_name"`
}

// NewModelFamily creates a new ModelFamily entity
func NewModelFamily(
	providerID string,
	schemaID string,
	name string,
	displayName string,
	description *string,
	maxTokens int,
	temperature float64,
	systemPrompt string,
) *ModelFamily {
	now := time.Now()
	return &ModelFamily{
		ProviderID:   providerID,
		SchemaID:     schemaID,
		Name:         name,
		DisplayName:  displayName,
		Description:  description,
		MaxTokens:    maxTokens,
		Temperature:  temperature,
		SystemPrompt: systemPrompt,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// NewDefaultModelFamily creates a new ModelFamily with default values
func NewDefaultModelFamily(
	providerID string,
	schemaID string,
	name string,
	displayName string,
	description *string,
) *ModelFamily {
	now := time.Now()
	return &ModelFamily{
		ProviderID:   providerID,
		SchemaID:     schemaID,
		Name:         name,
		DisplayName:  displayName,
		Description:  description,
		MaxTokens:    1024,
		Temperature:  1.0,
		SystemPrompt: "You are a helpful assistant.",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// Update updates the model family fields
func (m *ModelFamily) Update(
	name string,
	displayName string,
	description *string,
	maxTokens int,
	temperature float64,
	systemPrompt string,
) {
	m.Name = name
	m.DisplayName = displayName
	m.Description = description
	m.MaxTokens = maxTokens
	m.Temperature = temperature
	m.SystemPrompt = systemPrompt
	m.UpdatedAt = time.Now()
}

// Validate validates the model family entity
func (m *ModelFamily) Validate() error {
	if m.ProviderID == "" {
		return ErrInvalidProviderID
	}
	if m.SchemaID == "" {
		return ErrInvalidSchemaID
	}
	return nil
}
