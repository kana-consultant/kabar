package routes

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	"seo-backend/internal/config"
	"seo-backend/internal/container"
	"seo-backend/internal/middleware/auth"

	// Swagger docs (akan digenerate)
	_ "seo-backend/docs"
)

// SetupRoutes configures all application routes
func SetupRoutes(cfg *config.Config, container *container.Container) *chi.Mux {
	r := chi.NewRouter()

	// GLOBAL MIDDLEWARE

	// Standard middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Compression middleware
	r.Use(middleware.Compress(5, "gzip"))

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Request-ID"},
		ExposedHeaders:   []string{"Link", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Security headers middleware
	r.Use(securityHeaders)

	// SWAGGER DOCUMENTATION

	// Swagger UI endpoint - mengambil dari YAML
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.yaml"), // Mengambil dari YAML
	))

	// Swagger YAML endpoint
	r.Get("/swagger/doc.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.yaml")
	})

	// HEALTH CHECK ENDPOINTS (No auth)
	r.Get("/health", healthCheckHandler)
	r.Get("/ready", readyCheckHandler)
	r.Get("/live", liveCheckHandler)

	// PUBLIC ROUTES (No auth required)
	r.Group(func(r chi.Router) {
		// Auth endpoints - PUBLIC
		r.Post("/api/auth/login", container.AuthHandler.Login)
		r.Post("/api/auth/register", container.AuthHandler.Register)
		r.Post("/api/auth/logout", container.AuthHandler.Logout)
		r.Post("/api/auth/forgot-password", container.AuthHandler.ForgotPassword)
	})

	// PROTECTED ROUTES (JWT Auth required)
	r.Group(func(r chi.Router) {
		// JWT Authentication middleware
		r.Use(auth.JWTMiddleware(cfg))

		r.Get("/api/dashboard/stats", container.DashboardHandler.GetStats)

		// AUTH ROUTES (Protected)

		r.Get("/api/auth/me", container.AuthHandler.GetMe)
		r.Post("/api/auth/change-password", container.AuthHandler.ChangePassword)

		// USER ROUTES

		r.Route("/api/users", func(r chi.Router) {
			r.Get("/", container.UserHandler.GetAll)
			r.Get("/me", container.UserHandler.GetCurrentUser)
			r.Post("/", container.UserHandler.Create)
			r.Get("/{id}", container.UserHandler.GetByID)
			r.Put("/{id}", container.UserHandler.Update)
			r.Delete("/{id}", container.UserHandler.Delete)
			r.Post("/{id}/password", container.UserHandler.UpdatePassword)
		})

		r.Route("/api/drafts", func(r chi.Router) {

			r.Get("/", container.DraftHandler.GetAll)
			r.Get("/scheduled", container.DraftHandler.GetAllScheduled)
			r.Post("/", container.DraftHandler.Create)
			r.Get("/{id}", container.DraftHandler.GetByID)
			r.Put("/{id}", container.DraftHandler.Update)
			r.Delete("/{id}", container.DraftHandler.Delete)

			r.Post("/{id}/publish", container.DraftHandler.Publish)
			r.Post("/publish", container.DraftHandler.PublishContent)

			r.Post("/schedule", container.DraftHandler.ScheduleDraft)
			r.Post("/schedule/cancel", container.DraftHandler.CancelScheduledDraft)
		})

		r.Route("/api/products", func(r chi.Router) {
			r.Get("/", container.ProductHandler.GetAll)
			r.Post("/", container.ProductHandler.Create)
			r.Get("/{id}", container.ProductHandler.GetByID)
			r.Put("/{id}", container.ProductHandler.Update)
			r.Delete("/{id}", container.ProductHandler.Delete)
			r.Post("/{id}/test", container.ProductHandler.UpdateConnectionStatus)
		})

		r.Route("/api/api-keys", func(r chi.Router) {
			r.Get("/", container.APIKeyHandler.GetAll)
			r.Post("/", container.APIKeyHandler.Create)
			r.Get("/{id}", container.APIKeyHandler.GetByID)
			r.Put("/{id}", container.APIKeyHandler.Update)
			r.Delete("/{id}", container.APIKeyHandler.Delete)
		})

		r.Route("/api/teams", func(r chi.Router) {
			// CRUD
			r.Get("/", container.TeamHandler.GetAll)
			r.Post("/", container.TeamHandler.Create)
			r.Get("/{id}", container.TeamHandler.GetByID)
			r.Put("/{id}", container.TeamHandler.Update)
			r.Delete("/{id}", container.TeamHandler.Delete)

			// Members
			r.Get("/{id}/members", container.TeamHandler.GetTeamMembers)
			r.Post("/{id}/members", container.TeamHandler.AddMember)
			r.Put("/{id}/members/{userId}/role", container.TeamHandler.UpdateMemberRole)
			r.Delete("/{id}/members/{userId}", container.TeamHandler.RemoveMember)

			// User teams
			r.Get("/user/{userId}", container.TeamHandler.GetUserTeams)
		})

		r.Route("/api/providers", func(r chi.Router) {
			r.Get("/", container.ProviderHandler.GetAll)
			r.Post("/", container.ProviderHandler.Create)
			r.Get("/{id}", container.ProviderHandler.GetByID)
			r.Put("/{id}", container.ProviderHandler.Update)
			r.Delete("/{id}", container.ProviderHandler.Delete)
		})

		r.Route("/api/history", func(r chi.Router) {
			r.Get("/", container.HistoryHandler.GetAll)
			r.Post("/", container.HistoryHandler.Create)
			r.Get("/{id}", container.HistoryHandler.GetByID)
			r.Put("/{id}", container.HistoryHandler.Update)
			r.Delete("/{id}", container.HistoryHandler.Delete)
		})

		r.Route("/api/generate", func(r chi.Router) {
			r.Post("/article", container.GenerateHandler.GenerateArticle)
			r.Post("/image", container.GenerateHandler.GenerateImage)
		})

		r.Route("/api/models", func(r chi.Router) {
			r.Get("/", container.ModelHandler.GetAll)
			r.Get("/with-status", container.ModelHandler.GetAllWithStatus)
			r.Get("/{id}", container.ModelHandler.GetByID)
			r.Get("/provider/{providerId}", container.ModelHandler.GetByProvider)
		})
	})

	// 404 Not Found handler
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"Resource not found"}`))
	})

	// 405 Method Not Allowed handler
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error":"Method not allowed"}`))
	})

	return r
}

// MIDDLEWARE FUNCTION

// securityHeaders adds security-related headers
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		next.ServeHTTP(w, r)
	})
}

// getClientIP extracts client IP from request
func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}

// HEALTH CHECK HANDLER

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
}

func readyCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ready","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
}

func liveCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"alive","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
}
