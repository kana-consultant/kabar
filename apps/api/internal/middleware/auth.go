package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"seo-backend/internal/config"
	domainAuth "seo-backend/internal/domain/auth"
	"seo-backend/internal/models"
)

type contextKey string

const ClaimsKey contextKey = "claims"

func writeError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func JWTMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeError(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				writeError(w, "Invalid token format", http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(cfg.JWTSecret), nil
			})

			if err != nil {
				if errors.Is(err, jwt.ErrTokenExpired) {
					writeError(w, "Token expired", http.StatusUnauthorized)
					return
				}
				writeError(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				writeError(w, "Token is not valid", http.StatusUnauthorized)
				return
			}

			mapClaims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				writeError(w, "Invalid claims", http.StatusUnauthorized)
				return
			}

			userID, ok := mapClaims["user_id"].(string)
			if !ok || userID == "" {
				writeError(w, "Invalid token claims: missing user_id", http.StatusUnauthorized)
				return
			}

			// TeamID optional, ambil jika ada
			teamID, _ := mapClaims["team_id"].(string)

			claims := &domainAuth.Claims{
				UserID: userID,
				Role:   mapClaims["role"].(string),
				Email:  mapClaims["email"].(string),
				TeamID: teamID,
			}

			ctx := context.WithValue(r.Context(), ClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetClaims(ctx context.Context) *domainAuth.Claims {
	claims, _ := ctx.Value(ClaimsKey).(*domainAuth.Claims)
	return claims
}

func GetUserID(ctx context.Context) string {
	if c := GetClaims(ctx); c != nil {
		return c.UserID
	}
	return ""
}

func GetUserRole(ctx context.Context) string {
	if c := GetClaims(ctx); c != nil {
		return c.Role
	}
	return ""
}

func GetUserEmail(ctx context.Context) string {
	if c := GetClaims(ctx); c != nil {
		return c.Email
	}
	return ""
}

func GetTeamID(ctx context.Context) string {
	if c := GetClaims(ctx); c != nil {
		return c.TeamID
	}
	return ""
}

func GetUserContext(r *http.Request) models.UserContext {
	ctx := r.Context()
	return &models.SimpleUserContext{
		UserID: GetUserID(ctx),
		TeamID: GetTeamID(ctx),
		Role:   GetUserRole(ctx),
	}
}
