// internal/domain/provider/service.go
package provider

import "context"

// ProviderService defines the provider business logic interface
type ProviderService interface {
	// CreateProvider creates a new provider (admin only)
	CreateProvider(ctx context.Context, req CreateProviderRequest, userCtx UserContext) (string, error)

	// GetProviderByID retrieves a provider by ID
	GetProviderByID(ctx context.Context, id string) (*APIProvider, error)

	// GetAllProviders retrieves all providers based on user role
	GetAllProviders(ctx context.Context, userCtx UserContext) ([]APIProvider, error)

	// UpdateProvider updates a provider (admin only)
	UpdateProvider(ctx context.Context, id string, req UpdateProviderRequest, userCtx UserContext) error

	// DeleteProvider deletes a provider (admin only)
	DeleteProvider(ctx context.Context, id string, userCtx UserContext) error
}
