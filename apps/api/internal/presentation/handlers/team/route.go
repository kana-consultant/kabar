package team

import (
	"database/sql"
	teamService "seo-backend/internal/application/team"
	baseRoutes "seo-backend/internal/domain/base"
	"seo-backend/internal/infrastructure/db/repositories"
	services "seo-backend/internal/service"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	baseRoutes  baseRoutes.Route
	TeamHandler TeamHandler
}

func NewRoute(db *sql.DB, chi *chi.Mux, emailService *services.SMTPEmailService) *Route {
	teamRepo := repositories.NewTeamRepository(db)
	memberRepository := repositories.NewMemberRepository(db)
	teamQueryBuilder := teamService.NewQueryBuilder()
	// ========== TEAM AUTHORIZER & VALIDATOR ==========
	teamAuthorizer := teamService.NewAuthorizer()
	teamValidator := teamService.NewValidator(memberRepository)
	invite := repositories.NewTeamInviteRepository(db)
	userRepo := repositories.NewUserRepository(db)

	teamService := teamService.NewService(
		teamRepo,
		memberRepository,
		teamQueryBuilder,
		teamAuthorizer,
		teamValidator,
		invite,
		userRepo,
		*emailService,
	)
	TeamHandler := NewTeamHandler(teamService)

	return &Route{
		baseRoutes: baseRoutes.Route{
			DB:  db,
			CHI: chi,
		},
		TeamHandler: *TeamHandler,
	}
}

func (h *Route) SetupRoutes() *chi.Mux {
	r := h.baseRoutes.CHI

	r.Route("/api/teams", func(r chi.Router) {

		// CRUD Team
		r.Get("/", h.TeamHandler.GetAll)
		r.Post("/", h.TeamHandler.Create)
		r.Get("/{id}", h.TeamHandler.GetByID)
		r.Put("/{id}", h.TeamHandler.Update)
		r.Delete("/{id}", h.TeamHandler.Delete)

		// Invite
		r.Post("/{teamId}/invites", h.TeamHandler.InviteMember)

		// Team Members
		r.Post("/{teamId}/members", h.TeamHandler.AddMember)
		r.Delete("/{teamId}/members/{userId}", h.TeamHandler.RemoveMember)
		r.Put("/{teamId}/members/{userId}/role", h.TeamHandler.UpdateMemberRole)
		r.Get("/{teamId}/members", h.TeamHandler.GetTeamMembers)

		// User Teams
		r.Get("/user/{userId}", h.TeamHandler.GetUserTeams)
	})

	return r
}
