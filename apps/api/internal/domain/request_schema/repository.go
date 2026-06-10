package request_schema

import (
	"context"
	"database/sql"
)

type Repository interface {
	// Create methods
	Create(ctx context.Context, rs *RequestSchema) error
	CreateWithTx(ctx context.Context, tx *sql.Tx, rs *RequestSchema) error

	// Read methods
	GetByID(ctx context.Context, id string) (*RequestSchema, error)
	GetByProviderAndName(ctx context.Context, providerID string, name string) (*RequestSchema, error)
	GetAll(ctx context.Context, limit, offset int) ([]RequestSchema, error)
	GetByProvider(ctx context.Context, providerID string) ([]RequestSchema, error)
	GetByProviderSingle(ctx context.Context, providerID string) (*RequestSchema, error)

	// Update methods
	Update(ctx context.Context, rs *RequestSchema) error
	UpdateWithTx(ctx context.Context, tx *sql.Tx, rs *RequestSchema) error

	// Delete methods
	Delete(ctx context.Context, id string) error
	DeleteWithTx(ctx context.Context, tx *sql.Tx, id string) error
	DeleteByTeam(ctx context.Context, teamID string) error

	// Utility methods
	Exists(ctx context.Context, providerID string, name string) (bool, error)
	Count(ctx context.Context) (int64, error)
	CountByProvider(ctx context.Context, providerID string) (int64, error)
}
