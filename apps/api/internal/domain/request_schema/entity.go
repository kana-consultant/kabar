package request_schema

import (
	"time"
)

// RequestSchema represents the request_schemas table entity
type RequestSchema struct {
	ID                  string    `json:"id"`
	ProviderID          string    `json:"provider_id"`
	Name                string    `json:"name"`
	EndpointPath        string    `json:"endpoint_path"`
	MaxTokensKey        *string   `json:"max_tokens_key"`
	SystemRoleKey       *string   `json:"system_role_key"`
	ResponseTextPath    *string   `json:"response_text_path"`
	ResponseImagePath   *string   `json:"response_image_path"`
	RequestTemplate     *string   `json:"request_template"`
	SupportsTemperature *bool     `json:"supports_temperature"`
	SupportsStreaming   *bool     `json:"supports_streaming"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// NewRequestSchema creates a new RequestSchema entity
func NewRequestSchema(
	providerID string,
	name string,
	endpointPath string,
	maxTokensKey *string,
	systemRoleKey *string,
	responseTextPath *string,
	responseImagePath *string,
	requestTemplate *string,
	supportsTemperature *bool,
	supportsStreaming *bool,
) *RequestSchema {
	now := time.Now()
	return &RequestSchema{
		ProviderID:          providerID,
		Name:                name,
		EndpointPath:        endpointPath,
		MaxTokensKey:        maxTokensKey,
		SystemRoleKey:       systemRoleKey,
		ResponseTextPath:    responseTextPath,
		ResponseImagePath:   responseImagePath,
		RequestTemplate:     requestTemplate,
		SupportsTemperature: supportsTemperature,
		SupportsStreaming:   supportsStreaming,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
}
