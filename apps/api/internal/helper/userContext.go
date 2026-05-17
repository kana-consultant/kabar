package helper

import (
	"net/http"
	"seo-backend/internal/models"
	auth "seo-backend/internal/presentation/middleware"
)

func GetUserContext(r *http.Request) models.UserContext {
	ctx := r.Context()
	return &models.SimpleUserContext{
		UserID: auth.GetUserID(ctx),
		TeamID: auth.GetTeamID(ctx),
		Role:   auth.GetUserRole(ctx),
	}
}
