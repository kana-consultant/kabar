package model_family

import (
	"errors"
	"time"
)

// Errors
var (
	ErrNotFound           = errors.New("model family not found")
	ErrDuplicate          = errors.New("model family with this provider and name already exists")
	ErrInvalidProviderID  = errors.New("invalid provider ID")
	ErrInvalidSchemaID    = errors.New("invalid schema ID")
	ErrInvalidName        = errors.New("invalid family name")
	ErrInvalidDisplayName = errors.New("invalid display name")
	ErrInvalidID          = errors.New("invalid ID format")
	ErrDatabase           = errors.New("database error")
	ErrInvalidMaxTokens   = errors.New("max_tokens must be greater than 0")
	ErrInvalidTemperature = errors.New("temperature must be between 0 and 2")
)

// CreateRequest DTO
type CreateRequest struct {
	ProviderID   string  `json:"provider_id" validate:"required"`
	SchemaID     string  `json:"schema_id" validate:"required"`
	Name         string  `json:"name" validate:"required"`
	DisplayName  string  `json:"display_name" validate:"required"`
	Description  *string `json:"description"`
	MaxTokens    int     `json:"max_tokens" validate:"min=1"`
	Temperature  float64 `json:"temperature" validate:"min=0,max=2"`
	SystemPrompt string  `json:"system_prompt"`
}

// UpdateRequest DTO
type UpdateRequest struct {
	SchemaID     *string  `json:"schema_id"`
	Name         *string  `json:"name"`
	DisplayName  *string  `json:"display_name"`
	Description  *string  `json:"description"`
	MaxTokens    *int     `json:"max_tokens" validate:"omitempty,min=1"`
	Temperature  *float64 `json:"temperature" validate:"omitempty,min=0,max=2"`
	SystemPrompt *string  `json:"system_prompt"`
}

// Response DTO
type Response struct {
	ID           string    `json:"id"`
	ProviderID   string    `json:"provider_id"`
	SchemaID     string    `json:"schema_id"`
	Name         string    `json:"name"`
	DisplayName  string    `json:"display_name"`
	Description  *string   `json:"description"`
	MaxTokens    int       `json:"max_tokens"`
	Temperature  float64   `json:"temperature"`
	SystemPrompt string    `json:"system_prompt"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ListResponse DTO
type ListResponse struct {
	Data       []Response `json:"data"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
	TotalPages int        `json:"total_pages"`
}

// ModelFamilyWithStatus for API response with limited fields
type ModelFamilyWithStatus struct {
	ID           string  `json:"id"`
	ProviderID   string  `json:"provider_id"`
	Name         string  `json:"name"`
	DisplayName  string  `json:"display_name"`
	SchemaID     string  `json:"schema_id"`
	MaxTokens    int     `json:"max_tokens"`
	Temperature  float64 `json:"temperature"`
	SystemPrompt string  `json:"system_prompt"`
}

// ToResponse converts entity to response DTO
func ToResponse(mf *ModelFamily) *Response {
	if mf == nil {
		return nil
	}
	return &Response{
		ID:           mf.ID,
		ProviderID:   mf.ProviderID,
		SchemaID:     mf.SchemaID,
		Name:         mf.Name,
		DisplayName:  mf.DisplayName,
		Description:  mf.Description,
		MaxTokens:    mf.MaxTokens,
		Temperature:  mf.Temperature,
		SystemPrompt: mf.SystemPrompt,
		CreatedAt:    mf.CreatedAt,
		UpdatedAt:    mf.UpdatedAt,
	}
}

// ToResponseFromSchema converts ModelFamilyWithSchema to response DTO
func ToResponseFromSchema(mf *ModelFamilyWithSchema) *Response {
	if mf == nil {
		return nil
	}
	return &Response{
		ID:           mf.ID,
		ProviderID:   mf.ProviderID,
		SchemaID:     mf.SchemaID,
		Name:         mf.Name,
		DisplayName:  mf.DisplayName,
		Description:  mf.Description,
		MaxTokens:    mf.MaxTokens,
		Temperature:  mf.Temperature,
		SystemPrompt: mf.SystemPrompt,
		CreatedAt:    mf.CreatedAt,
		UpdatedAt:    mf.UpdatedAt,
	}
}

// ToResponseList converts entity list to response list
func ToResponseList(families []ModelFamily) []Response {
	responses := make([]Response, len(families))
	for i, family := range families {
		resp := ToResponse(&family)
		if resp != nil {
			responses[i] = *resp
		}
	}
	return responses
}

// ToResponseListFromSchema converts ModelFamilyWithSchema list to response list
func ToResponseListFromSchema(families []ModelFamilyWithSchema) []Response {
	responses := make([]Response, len(families))
	for i, family := range families {
		resp := ToResponseFromSchema(&family)
		if resp != nil {
			responses[i] = *resp
		}
	}
	return responses
}

// Deprecated: Use ToResponseListFromSchema instead
func ToResponseListSchem(families []ModelFamilyWithSchema) []Response {
	return ToResponseListFromSchema(families)
}

// ToResponseListWithProvider converts ModelFamilyWithProvider list to response list
func ToResponseListWithProvider(families []ModelFamilyWithProvider) []Response {
	responses := make([]Response, len(families))
	for i, family := range families {
		responses[i] = Response{
			ID:           family.ID,
			ProviderID:   family.ProviderID,
			SchemaID:     family.SchemaID,
			Name:         family.Name,
			DisplayName:  family.DisplayName,
			Description:  family.Description,
			MaxTokens:    family.MaxTokens,
			Temperature:  family.Temperature,
			SystemPrompt: family.SystemPrompt,
			CreatedAt:    family.CreatedAt,
			UpdatedAt:    family.UpdatedAt,
		}
	}
	return responses
}
