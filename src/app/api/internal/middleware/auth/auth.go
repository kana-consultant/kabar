// internal/middleware/auth.go
package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"seo-backend/internal/config"
)

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	UserRoleKey  contextKey = "user_role"
	UserEmailKey contextKey = "user_email"
	TeamIDKey    contextKey = "team_id"
)

// JWTMiddleware untuk authentication
func JWTMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				http.Error(w, "Invalid token format", http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, http.ErrAbortHandler
				}
				return []byte(cfg.JWTSecret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid claims", http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			if userID, ok := claims["user_id"].(string); ok {
				ctx = context.WithValue(ctx, UserIDKey, userID)
			}
			if role, ok := claims["role"].(string); ok {
				ctx = context.WithValue(ctx, UserRoleKey, role)
			}
			if email, ok := claims["email"].(string); ok {
				ctx = context.WithValue(ctx, UserEmailKey, email)
			}
			if teamID, ok := claims["team_id"].(string); ok {
				ctx = context.WithValue(ctx, TeamIDKey, teamID)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID mengambil user_id dari context
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

// GetUserRole mengambil user_role dari context
func GetUserRole(ctx context.Context) string {
	if role, ok := ctx.Value(UserRoleKey).(string); ok {
		return role
	}
	return ""
}

// GetUserEmail mengambil user_email dari context
func GetUserEmail(ctx context.Context) string {
	if email, ok := ctx.Value(UserEmailKey).(string); ok {
		return email
	}
	return ""
}

// GetTeamID mengambil team_id dari context
func GetTeamID(ctx context.Context) string {
	if teamID, ok := ctx.Value(TeamIDKey).(string); ok {
		return teamID
	}
	return ""
}
