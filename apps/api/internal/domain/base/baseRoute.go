package baseRoutes

import (
	"database/sql"
	"seo-backend/internal/infrastructure/db/repositories/rbac"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	DB        *sql.DB
	CHI       chi.Router
	PermCache *rbac.PermissionCache
}
