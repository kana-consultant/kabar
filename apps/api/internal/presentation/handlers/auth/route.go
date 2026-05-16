package auth

import (
	"database/sql"
	AuthService "seo-backend/internal/application/auth"
	"seo-backend/internal/domain/auth"
	baseRoutes "seo-backend/internal/domain/base"
	"seo-backend/internal/infrastructure/db/repositories"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	baseRoute   baseRoutes.Route
	AuthHandler AuthHandler
}

func NewRoute(db *sql.DB, chi chi.Router, tokenGen auth.TokenGenerator) *Route {
	AuthRepo := repositories.NewAuthRepository(db)
	AuthService := AuthService.NewService(db, AuthRepo, tokenGen)
	authHandler := NewAuthHandler(AuthService)
	return &Route{
		baseRoute: baseRoutes.Route{
			DB:  db,
			CHI: chi,
		},
		AuthHandler: *authHandler,
	}
}

func (h *Route) SetupRoute() chi.Router {
	r := h.baseRoute.CHI

	r.Post("/api/auth/login", h.AuthHandler.Login)
	r.Post("/api/auth/register", h.AuthHandler.Register)
	r.Post("/api/auth/logout", h.AuthHandler.Logout)
	r.Post("/api/auth/forgot-password", h.AuthHandler.ForgotPassword)

	return r
}

func (h *Route) AuthSettingRoute() chi.Router {
	r := h.baseRoute.CHI

	r.Post("/auth/logout", h.AuthHandler.Logout)
	r.Post("/auth/change-password", h.AuthHandler.ChangePassword)

	return r
}
