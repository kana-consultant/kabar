package container

import (
	"database/sql"
	"net/http"
	"time"

	// Infrastructure
	"seo-backend/internal/config"
	"seo-backend/internal/helper"
	"seo-backend/internal/infrastructure/db/repositories"
	"seo-backend/internal/infrastructure/db/repositories/rbac"
	auth "seo-backend/internal/middleware"
	"seo-backend/internal/pkg/jwt"
	"seo-backend/internal/scheduler"
	services "seo-backend/internal/service"

	// Handlers (Routers)
	aimodelHandler "seo-backend/internal/presentation/aimodel"
	apikeyHandler "seo-backend/internal/presentation/apikey"
	authHandler "seo-backend/internal/presentation/auth"
	dashboardHandler "seo-backend/internal/presentation/dashboard"
	draftHandler "seo-backend/internal/presentation/draft"
	generateHandler "seo-backend/internal/presentation/generate"
	historyHandler "seo-backend/internal/presentation/history"
	familiesHandler "seo-backend/internal/presentation/model_families"
	productHandler "seo-backend/internal/presentation/product"
	providerHandler "seo-backend/internal/presentation/provider"
	schemasHandler "seo-backend/internal/presentation/request_schema"
	teamHandler "seo-backend/internal/presentation/team"
	userHandler "seo-backend/internal/presentation/user"
	workflowNodeHandler "seo-backend/internal/presentation/workflow_node"

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

	// Public team routes untuk invite acceptance

	// ← buat SATU instance permCache di sini
	rbacRepo := repositories.NewRBACRepository(db)
	permCache := rbac.NewPermissionCache(redisClient, rbacRepo, 10*time.Minute)

	authHandler.NewRoute(db, r, jwtGenerator, permCache).SetupRoute()
	teamHandler.NewRoute(db, r, emailService).SetupPublicRoutes()

	// Protected routes
	r.Route("/api/", func(protected chi.Router) {
		protected.Use(auth.JWTMiddleware(cfg))

		authHandler.NewRoute(db, protected, jwtGenerator, permCache).AuthSettingRoute()
		dashboardHandler.NewRoute(db, protected).SetupRoutes()
		draftHandler.NewRoute(db, protected, redisClient, redisScheduler, permCache).SetupRoutes()
		generateHandler.NewRoute(db, protected, cfg).SetupRoutes()
		historyHandler.NewHistoryRoute(db, protected, permCache).SetupRoute()
		productHandler.NewRoute(db, protected, permCache).SetupRoutes()
		providerHandler.NewRoute(db, protected, redisClient).SetupRoutes()
		teamHandler.NewRoute(db, protected, emailService).SetupRoutes()
		userHandler.NewRoute(db, protected, permCache).SetupRoutes()
		apikeyHandler.NewRoute(db, protected).SetupRoutes()
		aimodelHandler.NewRoute(db, protected, redisClient).SetupRoute()
		schemasHandler.NewRoute(db, protected, redisClient).SetupRoutes()
		familiesHandler.NewRoute(db, protected, redisClient).SetupRoute()
		workflowNodeHandler.NewRoute(db, protected).SetupRoutes()
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
