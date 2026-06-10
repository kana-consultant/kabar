package request_schema

import (
	"errors"
	"time"
)

// Errors
var (
	ErrNotFound            = errors.New("request schema not found")
	ErrDuplicate           = errors.New("request schema with this provider and name already exists")
	ErrInvalidProviderID   = errors.New("invalid provider ID")
	ErrInvalidName         = errors.New("invalid schema name")
	ErrInvalidEndpointPath = errors.New("endpoint path is required")
	ErrInvalidID           = errors.New("invalid ID format")
	ErrDatabase            = errors.New("database error")
)

// CreateRequest DTO
type CreateRequest struct {
	ProviderID          string  `json:"provider_id"`
	Name                string  `json:"name"`
	EndpointPath        string  `json:"endpoint_path"`
	MaxTokensKey        *string `json:"max_tokens_key"`
	SystemRoleKey       *string `json:"system_role_key"`
	ResponseTextPath    *string `json:"response_text_path"`
	ResponseImagePath   *string `json:"response_image_path"`
	RequestTemplate     *string `json:"request_template"`
	SupportsTemperature *bool   `json:"supports_temperature"`
	SupportsStreaming   *bool   `json:"supports_streaming"`
}

// UpdateRequest DTO
type UpdateRequest struct {
	Name                *string `json:"name"`
	EndpointPath        *string `json:"endpoint_path"`
	MaxTokensKey        *string `json:"max_tokens_key"`
	SystemRoleKey       *string `json:"system_role_key"`
	ResponseTextPath    *string `json:"response_text_path"`
	ResponseImagePath   *string `json:"response_image_path"`
	RequestTemplate     *string `json:"request_template"`
	SupportsTemperature *bool   `json:"supports_temperature"`
	SupportsStreaming   *bool   `json:"supports_streaming"`
}

// Response DTO
type Response struct {
	ID                  string    `json:"id"`
	ProviderID          string    `json:"provider_id"`
	Name                string    `json:"name"`
	EndpointPath        string    `json:"endpoint_path"`
	MaxTokensKey        *string   `json:"max_tokens_key"`
	SystemRoleKey       *string   `json:"system_role_key"`
	ResponseTextPath    *string   `json:"response_text_path"`
	ResponseImagePath   *string   `json:"response_image_path"`
	RequestTemplate     *string   `json:"request_template"`
	SupportsTemperature bool      `json:"supports_temperature"`
	SupportsStreaming   bool      `json:"supports_streaming"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
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
func ToResponse(rs *RequestSchema) *Response {
	supportsTemp := true
	if rs.SupportsTemperature != nil {
		supportsTemp = *rs.SupportsTemperature
	}

	supportsStream := true
	if rs.SupportsStreaming != nil {
		supportsStream = *rs.SupportsStreaming
	}

	return &Response{
		ID:                  rs.ID,
		ProviderID:          rs.ProviderID,
		Name:                rs.Name,
		EndpointPath:        rs.EndpointPath,
		MaxTokensKey:        rs.MaxTokensKey,
		SystemRoleKey:       rs.SystemRoleKey,
		ResponseTextPath:    rs.ResponseTextPath,
		ResponseImagePath:   rs.ResponseImagePath,
		RequestTemplate:     rs.RequestTemplate,
		SupportsTemperature: supportsTemp,
		SupportsStreaming:   supportsStream,
		CreatedAt:           rs.CreatedAt,
		UpdatedAt:           rs.UpdatedAt,
	}
}

// ToResponseList converts entity list to response list
func ToResponseList(schemas []RequestSchema) []Response {
	responses := make([]Response, len(schemas))
	for i, schema := range schemas {
		responses[i] = *ToResponse(&schema)
	}
	return responses
}
