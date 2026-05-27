package container

import (
	"database/sql"
	"net/http"
	"time"

	// Infrastructure
	"seo-backend/internal/config"
	"seo-backend/internal/helper"
	"seo-backend/internal/pkg/jwt"
	auth "seo-backend/internal/presentation/middleware"
	"seo-backend/internal/scheduler"
	services "seo-backend/internal/service"

	// Handlers (Routers)
	aimodelHandler "seo-backend/internal/presentation/handlers/aimodel"
	apikeyHandler "seo-backend/internal/presentation/handlers/apikey"
	authHandler "seo-backend/internal/presentation/handlers/auth"
	dashboardHandler "seo-backend/internal/presentation/handlers/dashboard"
	draftHandler "seo-backend/internal/presentation/handlers/draft"
	generateHandler "seo-backend/internal/presentation/handlers/generate"
	historyHandler "seo-backend/internal/presentation/handlers/history"
	productHandler "seo-backend/internal/presentation/handlers/product"
	providerHandler "seo-backend/internal/presentation/handlers/provider"
	teamHandler "seo-backend/internal/presentation/handlers/team"
	userHandler "seo-backend/internal/presentation/handlers/user"

	// Utilities

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-redis/redis/v8"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Container struct {
	Router *chi.Mux
}

func NewContainer(
	cfg *config.Config,
	db *sql.DB,
	redisScheduler *scheduler.RedisScheduler,
	redisClient *redis.Client,
	emailService *services.SMTPEmailService,
) *Container {

	r := chi.NewRouter()

	// middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Compress(5, "gzip"))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Request-ID"},
		ExposedHeaders:   []string{"Link", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(securityHeaders)

	// Swagger
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.yaml"),
	))

	r.Get("/swagger/doc.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.yaml")
	})

	// Health check
	r.Get("/health", healthCheckHandler)
	r.Get("/ready", readyCheckHandler)
	r.Get("/live", liveCheckHandler)

	jwtSecret := cfg.JWTSecret
	jwtExpiry := cfg.JWTExpiry
	jwtGenerator, _ := jwt.NewGenerator(jwtSecret, jwtExpiry)

	// Public routes
	// Public team routes untuk invite acceptance

	authHandler.NewRoute(db, r, jwtGenerator).SetupRoute()
	teamHandler.NewRoute(db, r, emailService).SetupPublicRoutes()

	// Protected routes
	r.Route("/api/", func(protected chi.Router) {
		protected.Use(auth.JWTMiddleware(cfg))
		authHandler.NewRoute(db, protected, jwtGenerator).AuthSettingRoute()
		dashboardHandler.NewRoute(db, protected).SetupRoutes()
		draftHandler.NewRoute(db, protected, redisClient, redisScheduler).SetupRoutes()
		generateHandler.NewRoute(db, protected).SetupRoutes()
		historyHandler.NewHistoryRoute(db, protected).SetupRoute()
		productHandler.NewRoute(db, protected).SetupRoutes()
		providerHandler.NewRoute(db, protected).SetupRoutes()
		teamHandler.NewRoute(db, protected, emailService).SetupRoutes()
		userHandler.NewRoute(db, protected).SetupRoutes()
		apikeyHandler.NewRoute(db, protected).SetupRoutes()
		aimodelHandler.NewRoute(db, protected).SetupRoute()
	})

	return &Container{
		Router: r,
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","timestamp":"` + helper.ParseWIBTime(time.Now().Format(time.RFC3339)).Format(time.RFC3339) + `"}`))
}

func readyCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ready","timestamp":"` + helper.ParseWIBTime(time.Now().Format(time.RFC3339)).Format(time.RFC3339) + `"}`))
}

func liveCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"alive","timestamp":"` + helper.ParseWIBTime(time.Now().Format(time.RFC3339)).Format(time.RFC3339) + `"}`))
}

// securityHeaders adds security-related headers
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		next.ServeHTTP(w, r)
	})
}
