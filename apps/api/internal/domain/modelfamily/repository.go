package model_family

import (
	"context"
	"database/sql"
	"seo-backend/internal/domain/paginate"
	"seo-backend/internal/models"
)

// Repository defines the data access interface for ModelFamily
type Repository interface {
	Create(ctx context.Context, modelFamily *ModelFamily) error
	CreateBatchWithTx(ctx context.Context, tx *sql.Tx, families []ModelFamilyWithSchema) error
	CreateWithTx(ctx context.Context, tx *sql.Tx, modelFamily *ModelFamily) error
	GetByID(ctx context.Context, id string) (*ModelFamily, error)
	GetByProviderAndName(ctx context.Context, providerID string, name string) (*ModelFamily, error)
	GetByProviderID(ctx context.Context, providerID string) ([]ModelFamilyWithSchema, error)
	GetAll(ctx context.Context, userCtx models.UserContext, params paginate.PaginationParams) ([]ModelFamilyWithProvider, error)
	GetBySchema(ctx context.Context, schemaID string) ([]ModelFamily, error)
	Update(ctx context.Context, modelFamily *ModelFamily) error
	UpdateWithTx(ctx context.Context, tx *sql.Tx, modelFamily *ModelFamilyWithSchema) error
	Delete(ctx context.Context, id string) error
	DeleteWithTx(ctx context.Context, tx *sql.Tx, id string) error
	Exists(ctx context.Context, providerID string, name string) (bool, error)
	Count(ctx context.Context) (int64, error)
	CountByProvider(ctx context.Context, providerID string) (int64, error)
	GetWithSchemaByID(ctx context.Context, id string) (*ModelFamilyWithSchema, error)
}
