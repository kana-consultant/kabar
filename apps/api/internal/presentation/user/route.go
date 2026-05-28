package user

import (
	"database/sql"
	userApp "seo-backend/internal/application/user"
	baseRoutes "seo-backend/internal/domain/base"
	"seo-backend/internal/infrastructure/db/repositories"
	rbacCache "seo-backend/internal/infrastructure/db/repositories/rbac"
	authmw "seo-backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	baseRoutes  baseRoutes.Route
	UserHandler UserHandler
	permCache   *rbacCache.PermissionCache
}

func NewRoute(db *sql.DB, chi chi.Router, permCache *rbacCache.PermissionCache) *Route {
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
		permCache:   permCache,
	}
}

func (h *Route) SetupRoutes() chi.Router {
	r := h.baseRoutes.CHI
	c := h.permCache
	r.Route("/users", func(r chi.Router) {
		r.Get("/", h.UserHandler.GetAll)
		r.With(authmw.UserManage(c)).Post("/", h.UserHandler.Create)
		r.With(authmw.UserManage(c)).Get("/{id}", h.UserHandler.GetByID)
		r.With(authmw.UserManage(c)).Put("/{id}", h.UserHandler.Update)
		r.With(authmw.UserManage(c)).Delete("/{id}", h.UserHandler.Delete)
		r.With(authmw.UserManage(c)).Post("/{id}/password", h.UserHandler.UpdatePassword)

		// endpoint ini tidak perlu permission khusus, cukup JWT valid
		r.Get("/me", h.UserHandler.GetCurrentUser)
	})
	return r
}
