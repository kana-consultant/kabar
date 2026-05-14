package user

import (
	"database/sql"
	userApp "seo-backend/internal/application/user"
	baseRoutes "seo-backend/internal/domain/base"
	"seo-backend/internal/infrastructure/db/repositories"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	baseRoutes  baseRoutes.Route
	UserHandler UserHandler
}

func NewRoute(db *sql.DB, chi *chi.Mux) *Route {

	userRepo := repositories.NewUserRepository(db)
	userQueryBuilder := userApp.NewQueryBuilder()
	userAuthorizer := userApp.NewAuthorizer(db)
	userPasswordService := userApp.NewPasswordService()
	userValidator := userApp.NewValidator(userRepo)
	userService := userApp.NewService(db, userRepo, userQueryBuilder, userAuthorizer, userPasswordService, userValidator)
	UserHandler := NewUserHandler(userService)
	return &Route{
		baseRoutes: baseRoutes.Route{
			DB:  db,
			CHI: chi,
		},
		UserHandler: *UserHandler,
	}
}

func (h *Route) SetupRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/api/users", func(r chi.Router) {
		r.Get("/", h.UserHandler.GetAll)
		r.Get("/me", h.UserHandler.GetCurrentUser)
		r.Post("/", h.UserHandler.Create)
		r.Get("/{id}", h.UserHandler.GetByID)
		r.Put("/{id}", h.UserHandler.Update)
		r.Delete("/{id}", h.UserHandler.Delete)
		r.Post("/{id}/password", h.UserHandler.UpdatePassword)
	})
	return r
}
