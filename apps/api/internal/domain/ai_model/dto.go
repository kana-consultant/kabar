package ai_model

import (
	"errors"
	"time"
)

// Errors
var (
	ErrInvalidID          = errors.New("invalid ID format")
	ErrInvalidName        = errors.New("invalid model name")
	ErrInvalidDisplayName = errors.New("invalid display name")
	ErrInvalidProvider    = errors.New("provider ID is required when no family ID is provided")
	ErrInvalidSchema      = errors.New("schema ID is required when no family ID is provided")
	ErrNoFamilyOrProvider = errors.New("either family_id or provider_id must be provided")
)

// CreateRequest DTO
type CreateRequest struct {
	FamilyID      *string  `json:"family_id"`
	ProviderID    *string  `json:"provider_id"`
	SchemaID      *string  `json:"schema_id"`
	TeamID        *string  `json:"team_id"`
	Name          string   `json:"name"`
	DisplayName   string   `json:"display_name"`
	Description   *string  `json:"description"`
	SystemPrompt  *string  `json:"system_prompt"`
	MaxTokens     *int     `json:"max_tokens"`
	Temperature   *float64 `json:"temperature"`
	ContextWindow *int     `json:"context_window"`
	IsActive      *bool    `json:"is_active"`
	IsDefault     *bool    `json:"is_default"`
	CreatedBy     *string  `json:"created_by"`
}

// UpdateRequest DTO
type UpdateRequest struct {
	FamilyID      *string  `json:"family_id"`
	ProviderID    *string  `json:"provider_id"`
	SchemaID      *string  `json:"schema_id"`
	TeamID        *string  `json:"team_id"`
	Name          *string  `json:"name"`
	DisplayName   *string  `json:"display_name"`
	Description   *string  `json:"description"`
	SystemPrompt  *string  `json:"system_prompt"`
	MaxTokens     *int     `json:"max_tokens"`
	Temperature   *float64 `json:"temperature"`
	ContextWindow *int     `json:"context_window"`
	IsActive      *bool    `json:"is_active"`
	IsDefault     *bool    `json:"is_default"`
	CreatedBy     *string  `json:"created_by"`
}

// Response DTO
type Response struct {
	ID            string    `json:"id"`
	FamilyID      *string   `json:"family_id,omitempty"`
	ProviderID    *string   `json:"provider_id,omitempty"`
	SchemaID      *string   `json:"schema_id,omitempty"`
	TeamID        *string   `json:"team_id,omitempty"`
	Name          string    `json:"name"`
	DisplayName   string    `json:"display_name"`
	Description   *string   `json:"description,omitempty"`
	SystemPrompt  *string   `json:"system_prompt,omitempty"`
	MaxTokens     *int      `json:"max_tokens,omitempty"`
	Temperature   *float64  `json:"temperature,omitempty"`
	ContextWindow *int      `json:"context_window,omitempty"`
	IsActive      bool      `json:"is_active"`
	IsDefault     bool      `json:"is_default"`
	CreatedBy     *string   `json:"created_by,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ListResponse DTO
type ListResponse struct {
	Data       []Response `json:"data"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
	TotalPages int        `json:"total_pages"`
}

// ToResponse converts entity to response DTO
func ToResponse(model *AIModel) *Response {
	isActive := true

	return &Response{
		ID:            model.ID,
		FamilyID:      model.FamilyID,
		ProviderID:    model.ProviderID,
		TeamID:        model.TeamID,
		Name:          model.Name,
		DisplayName:   *model.DisplayName,
		Description:   model.Description,
		SystemPrompt:  model.SystemPrompt,
		MaxTokens:     model.MaxTokens,
		Temperature:   model.Temperature,
		ContextWindow: model.ContextWindow,
		IsActive:      isActive,
		CreatedBy:     model.CreatedBy,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}
}

// RequestSchemaResponse (schema saja)
type RequestSchemaResponse struct {
	ID                  string     `json:"id"`
	ProviderID          string     `json:"provider_id"`
	Name                string     `json:"name"`
	EndpointPath        string     `json:"endpoint_path"`
	RequestTemplate     *string    `json:"request_template,omitempty"`
	ResponseTextPath    string     `json:"response_text_path"`
	ResponseImagePath   *string    `json:"response_image_path,omitempty"`
	SupportsTemperature bool       `json:"supports_temperature"`
	SupportsStreaming   bool       `json:"supports_streaming"`
	CreatedAt           *time.Time `json:"created_at,omitempty"`
	UpdatedAt           *time.Time `json:"updated_at,omitempty"`
}

type AIModelResponse struct {
	ID            string                `json:"id"`
	FamilyID      *string               `json:"family_id,omitempty"`
	ProviderID    string                `json:"provider_id"`
	SchemaID      *string               `json:"schema_id,omitempty"`
	TeamID        *string               `json:"team_id,omitempty"`
	Name          string                `json:"name"`
	DisplayName   string                `json:"display_name"`
	Description   *string               `json:"description,omitempty"`
	SystemPrompt  *string               `json:"system_prompt,omitempty"`
	MaxTokens     *int                  `json:"max_tokens,omitempty"`
	Temperature   *float64              `json:"temperature,omitempty"`
	ContextWindow *int                  `json:"context_window,omitempty"`
	IsActive      bool                  `json:"is_active"`
	IsDefault     bool                  `json:"is_default"`
	Schema        RequestSchemaResponse `json:"schema"`
	CreatedBy     *string               `json:"created_by,omitempty"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
}

// ToResponseSchema maps AIModel entity to RequestSchemaResponse (mengambil schema dari dalam AIModel)
func ToResponseSchema(entity *AiModelSchema) *RequestSchemaResponse {
	if entity == nil {
		return nil
	}

	// Cek apakah schema ada (ID tidak nil)
	if entity.Schema.ID == nil {
		return nil
	}

	return &RequestSchemaResponse{
		ID:                  *entity.Schema.ID,
		ProviderID:          *entity.Schema.ProviderID,
		EndpointPath:        *entity.Schema.EndpointPath,
		RequestTemplate:     entity.Schema.RequestTemplate, // sudah pointer, bisa nil
		ResponseTextPath:    *entity.Schema.ResponseTextPath,
		ResponseImagePath:   entity.Schema.ResponseImagePath, // sudah pointer, bisa nil
		SupportsTemperature: *entity.Schema.SupportsTemperature,
		SupportsStreaming:   *entity.Schema.SupportsStreaming,
		CreatedAt:           entity.Schema.CreatedAt, // sudah pointer, bisa nil
		UpdatedAt:           entity.Schema.UpdatedAt, // sudah pointer, bisa nil
	}
}

// ToResponseSlice maps slice of entities to slice of responses
func ToResponseSlice(entities []*AIModel) []*Response {
	if entities == nil {
		return nil
	}

	responses := make([]*Response, 0, len(entities))
	for _, entity := range entities {
		responses = append(responses, ToResponse(entity))
	}
	return responses
}

// ToResponseList converts entity list to response list
func ToResponseList(models []AIModel) []Response {
	responses := make([]Response, len(models))
	for i, model := range models {
		responses[i] = *ToResponse(&model)
	}
	return responses
}
