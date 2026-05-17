package baseRoutes

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	DB  *sql.DB
	CHI chi.Router
}
