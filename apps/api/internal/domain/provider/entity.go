package provider

import (
	"encoding/json"
	"errors"
	model_family "seo-backend/internal/domain/modelfamily"
	"time"

	"github.com/google/uuid"
)

// Errors
var (
	// Basic errors
	ErrNotFound       = errors.New("API provider not found")
	ErrDuplicate      = errors.New("API provider with this name already exists")
	ErrInvalidID      = errors.New("invalid ID format")
	ErrInvalidName    = errors.New("invalid provider name")
	ErrInvalidBaseURL = errors.New("invalid base URL")
	ErrDatabase       = errors.New("database error")

	// Additional errors
	ErrEmptyName          = errors.New("provider name cannot be empty")
	ErrEmptyDisplayName   = errors.New("display name cannot be empty")
	ErrEmptyBaseURL       = errors.New("base URL cannot be empty")
	ErrInvalidAuthType    = errors.New("invalid auth type, must be 'bearer', 'basic', 'api_key', or 'none'")
	ErrInvalidAuthHeader  = errors.New("invalid auth header format")
	ErrInvalidAuthPrefix  = errors.New("invalid auth prefix")
	ErrInvalidHeaders     = errors.New("invalid default headers format")
	ErrProviderInUse      = errors.New("provider is still in use by model families or models")
	ErrProviderInactive   = errors.New("provider is inactive")
	ErrProviderActive     = errors.New("provider is already active")
	ErrProviderInactiveOp = errors.New("cannot perform operation on inactive provider")
	ErrUnauthorized       = errors.New("unauthorized access to provider")
	ErrForbidden          = errors.New("access forbidden: insufficient permissions")
	ErrInvalidUpdateData  = errors.New("invalid update data provided")
	ErrNoUpdatesProvided  = errors.New("no updates provided")
	ErrTransactionFailed  = errors.New("transaction failed")
	ErrOperationTimeout   = errors.New("operation timed out")
)

// ValidationError represents a validation error with field context
type ValidationError struct {
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string, value interface{}) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	}
}

// MultiError represents multiple validation errors
type MultiError struct {
	Errors []error `json:"errors"`
}

func (e *MultiError) Error() string {
	return "multiple validation errors occurred"
}

func (e *MultiError) Add(err error) {
	e.Errors = append(e.Errors, err)
}

func (e *MultiError) HasErrors() bool {
	return len(e.Errors) > 0
}

// APIProvider represents the api_providers table entity
type APIProvider struct {
	ID             uuid.UUID                            `json:"id"`
	Name           string                               `json:"name"`
	DisplayName    string                               `json:"display_name"`
	Description    *string                              `json:"description,omitempty"`
	BaseURL        string                               `json:"base_url"`
	AuthType       *string                              `json:"auth_type"`
	AuthHeader     *string                              `json:"auth_header"`
	AuthPrefix     *string                              `json:"auth_prefix"`
	DefaultHeaders json.RawMessage                      `json:"default_headers"`
	IsActive       *bool                                `json:"is_active"`
	CreatedAt      time.Time                            `json:"created_at"`
	UpdatedAt      time.Time                            `json:"updated_at"`
	Families       []model_family.ModelFamilyWithSchema `json:"families"`
}

// ProviderWithStatus for API response with limited fields
type ProviderWithStatus struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	BaseURL     string `json:"base_url"`
	IsActive    bool   `json:"is_active"`
}
