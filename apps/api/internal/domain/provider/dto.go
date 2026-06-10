package provider

import (
	"encoding/json"
	"time"

	model_family "seo-backend/internal/domain/modelfamily"

	"github.com/google/uuid"
)

// CreateRequest DTO
type CreateRequest struct {
	Name           string                               `json:"name" binding:"required"`
	DisplayName    string                               `json:"display_name" binding:"required"`
	Description    *string                              `json:"description"`
	BaseURL        string                               `json:"base_url" binding:"required,url"`
	AuthType       *string                              `json:"auth_type"`
	AuthHeader     *string                              `json:"auth_header"`
	AuthPrefix     *string                              `json:"auth_prefix"`
	DefaultHeaders json.RawMessage                      `json:"default_headers"`
	IsActive       *bool                                `json:"is_active"`
	Families       []model_family.ModelFamilyWithSchema `json:"families"`
}

// UpdateRequest DTO
type UpdateRequest struct {
	Name           *string                              `json:"name"`
	DisplayName    *string                              `json:"display_name"`
	Description    *string                              `json:"description"`
	BaseURL        *string                              `json:"base_url"`
	AuthType       *string                              `json:"auth_type"`
	AuthHeader     *string                              `json:"auth_header"`
	AuthPrefix     *string                              `json:"auth_prefix"`
	DefaultHeaders json.RawMessage                      `json:"default_headers"`
	IsActive       *bool                                `json:"is_active"`
	Families       []model_family.ModelFamilyWithSchema `json:"families,omitempty"`
}

// Response DTO
type Response struct {
	ID             uuid.UUID                            `json:"id"`
	Name           string                               `json:"name"`
	DisplayName    string                               `json:"display_name"`
	Description    *string                              `json:"description,omitempty"`
	BaseURL        string                               `json:"base_url"`
	AuthType       *string                              `json:"auth_type"`
	AuthHeader     *string                              `json:"auth_header"`
	AuthPrefix     *string                              `json:"auth_prefix"`
	DefaultHeaders json.RawMessage                      `json:"default_headers"`
	IsActive       bool                                 `json:"is_active"`
	CreatedAt      time.Time                            `json:"created_at"`
	UpdatedAt      time.Time                            `json:"updated_at"`
	Families       []model_family.ModelFamilyWithSchema `json:"families,omitempty"`
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
func ToResponse(provider *APIProvider) *Response {
	return &Response{
		ID:             provider.ID,
		Name:           provider.Name,
		DisplayName:    provider.DisplayName,
		Description:    provider.Description,
		BaseURL:        provider.BaseURL,
		AuthType:       provider.AuthType,
		AuthHeader:     provider.AuthHeader,
		AuthPrefix:     provider.AuthPrefix,
		DefaultHeaders: provider.DefaultHeaders,
		IsActive:       *provider.IsActive,
		CreatedAt:      provider.CreatedAt,
		UpdatedAt:      provider.UpdatedAt,
		Families:       provider.Families,
	}
}

// ToResponseList converts entity list to response list
func ToResponseList(providers []APIProvider) []Response {
	responses := make([]Response, len(providers))
	for i, provider := range providers {
		responses[i] = *ToResponse(&provider)
	}
	return responses
}
