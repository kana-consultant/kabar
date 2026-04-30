package routes

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	"seo-backend/internal/config"
	"seo-backend/internal/database"
	"seo-backend/internal/handlers"
	apikey "seo-backend/internal/handlers/api_key"
	auth_handler "seo-backend/internal/handlers/auth"
	"seo-backend/internal/handlers/dashboard"
	"seo-backend/internal/handlers/draft"
	"seo-backend/internal/handlers/history"
	"seo-backend/internal/handlers/product"
	"seo-backend/internal/handlers/provider"
	"seo-backend/internal/handlers/team"
	"seo-backend/internal/handlers/user"
	"seo-backend/internal/middleware/auth"
	"seo-backend/internal/scheduler"
)

func SetupRoutes(cfg *config.Config, sched *scheduler.RedisScheduler) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000", "http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Initialize handlers

	// Initialize handlers
	authHandler := auth_handler.NewAuthHandler(cfg)
	userHandler := user.NewUserHandler()
	productHandler := product.NewProductHandler()
	draftHandler := draft.NewDraftHandler(sched)
	historyHandler := history.NewHistoryHandler()
	generateHandler := handlers.NewGenerateHandler()
	teamHandler := team.NewTeamHandler()
	dashboardHandler := dashboard.NewDashboardHandler()
	apiKeyHandler := apikey.NewAPIKeyHandler()
	modelHandler := handlers.NewAIModelHandler()
	providerHandler := provider.NewProviderHandler()

	// =====================================================
	// SWAGGER DOCUMENTATION
	// =====================================================
	r.Get("/swagger/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "docs/swagger.json")
	})

	r.Get("/swagger/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/yaml")
		http.ServeFile(w, r, "docs/swagger.yaml")
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/swagger.yaml"),
	))

	// =====================================================
	// ROOT ROUTE
	// =====================================================
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"name":        "SEO Backend API",
			"version":     "1.0.0",
			"description": "API for SEO content management system",
			"status":      "running",
			"timestamp":   time.Now().Format(time.RFC3339),
			"endpoints": map[string]string{
				"health":  "/health",
				"api":     "/api",
				"docs":    "/swagger/index.html",
				"metrics": "/metrics",
			},
		})
	})

	// =====================================================
	// API INFO ROUTE
	// =====================================================
	r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"name":    "SEO Backend API",
			"version": "1.0.0",
			"routes": map[string]interface{}{
				"auth": map[string]string{
					"login":    "POST /api/auth/login",
					"register": "POST /api/auth/register",
				},
				"dashboard": map[string]string{
					"stats": "GET /api/dashboard/stats",
				},
				"users": map[string]string{
					"list":   "GET /api/users",
					"detail": "GET /api/users/{id}",
					"create": "POST /api/users",
					"update": "PUT /api/users/{id}",
					"delete": "DELETE /api/users/{id}",
				},
				"teams": map[string]string{
					"list":   "GET /api/teams",
					"detail": "GET /api/teams/{id}",
					"create": "POST /api/teams",
					"update": "PUT /api/teams/{id}",
					"delete": "DELETE /api/teams/{id}",
				},
				"products": map[string]string{
					"list":   "GET /api/products",
					"detail": "GET /api/products/{id}",
					"create": "POST /api/products",
					"update": "PUT /api/products/{id}",
					"delete": "DELETE /api/products/{id}",
					"test":   "POST /api/products/{id}/test",
				},
				"drafts": map[string]string{
					"list":    "GET /api/drafts",
					"detail":  "GET /api/drafts/{id}",
					"create":  "POST /api/drafts",
					"update":  "PUT /api/drafts/{id}",
					"delete":  "DELETE /api/drafts/{id}",
					"publish": "POST /api/drafts/{id}/publish",
				},
				"history": map[string]string{
					"list":   "GET /api/history",
					"add":    "POST /api/history",
					"delete": "DELETE /api/history/{id}",
					"clear":  "DELETE /api/history",
				},
				"generate": map[string]string{
					"article": "POST /api/generate/article",
					"image":   "POST /api/generate/image",
				},
			},
		})
	})

	// =====================================================
	// HEALTH CHECK
	// =====================================================
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		dbStatus := "connected"
		if err := database.GetDB().Ping(); err != nil {
			dbStatus = "disconnected"
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
			"services": map[string]string{
				"database": dbStatus,
				"api":      "running",
			},
			"system": map[string]interface{}{
				"go_version": runtime.Version(),
				"os":         runtime.GOOS,
				"arch":       runtime.GOARCH,
				"goroutines": runtime.NumGoroutine(),
			},
		})
	})

	// =====================================================
	// METRICS
	// =====================================================
	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"goroutines": runtime.NumGoroutine(),
			"cpu_cores":  runtime.NumCPU(),
			"go_version": runtime.Version(),
			"compiler":   runtime.Compiler,
			"os":         runtime.GOOS,
			"arch":       runtime.GOARCH,
			"memory": map[string]interface{}{
				"alloc":          mem.Alloc,
				"total_alloc":    mem.TotalAlloc,
				"sys":            mem.Sys,
				"heap_alloc":     mem.HeapAlloc,
				"heap_sys":       mem.HeapSys,
				"heap_objects":   mem.HeapObjects,
				"stack_inuse":    mem.StackInuse,
				"stack_sys":      mem.StackSys,
				"gc_cycles":      mem.NumGC,
				"last_gc":        mem.LastGC,
				"next_gc":        mem.NextGC,
				"pause_total_ns": mem.PauseTotalNs,
			},
		})
	})

	// =====================================================
	// PUBLIC AUTH ROUTES (No auth required)
	// =====================================================
	r.Post("/api/auth/login", authHandler.Login)
	r.Post("/api/auth/register", authHandler.Register)
	r.Post("/api/auth/forgot-password", authHandler.ForgotPassword)

	// =====================================================
	// PROTECTED AUTH ROUTES (Auth required)
	// =====================================================
	r.Group(func(r chi.Router) {
		r.Use(auth.JWTMiddleware(cfg))

		r.Get("/api/auth/me", authHandler.GetMe)
		r.Post("/api/auth/change-password", authHandler.ChangePassword)
	})

	// =====================================================
	// PROTECTED API ROUTES (Auth required)
	// =====================================================
	r.Group(func(r chi.Router) {
		r.Use(auth.JWTMiddleware(cfg))

		// Dashboard
		r.Get("/api/dashboard/stats", dashboardHandler.GetStats)

		// Users routes
		r.Get("/api/users", userHandler.GetAll)
		r.Get("/api/users/{id}", userHandler.GetByID)
		r.Post("/api/users", userHandler.Create)
		r.Put("/api/users/{id}", userHandler.Update)
		r.Delete("/api/users/{id}", userHandler.Delete)

		// Teams routes
		r.Get("/api/teams", teamHandler.GetAll)
		r.Get("/api/teams/{id}", teamHandler.GetByID)
		r.Post("/api/teams", teamHandler.Create)
		r.Put("/api/teams/{id}", teamHandler.Update)
		r.Delete("/api/teams/{id}", teamHandler.Delete)

		// Team members routes
		r.Get("/api/teams/{id}/members", teamHandler.GetTeamMembers)
		r.Post("/api/teams/{id}/members", teamHandler.AddMember)
		r.Put("/api/teams/{id}/members/{userId}", teamHandler.UpdateMemberRole)
		r.Delete("/api/teams/{id}/members/{userId}", teamHandler.RemoveMember)

		// User teams
		r.Get("/api/users/{userId}/teams", teamHandler.GetUserTeams)

		// Products routes
		r.Get("/api/products", productHandler.GetAll)
		r.Get("/api/products/{id}", productHandler.GetByID)
		r.Post("/api/products", productHandler.Create)
		r.Put("/api/products/{id}", productHandler.Update)
		r.Delete("/api/products/{id}", productHandler.Delete)
		r.Post("/api/products/{id}/test", productHandler.TestConnection)

		// Drafts routes
		r.Get("/api/drafts", draftHandler.GetAll)
		r.Get("/api/drafts/{id}", draftHandler.GetByID)
		r.Post("/api/drafts", draftHandler.Create)
		r.Put("/api/drafts/{id}", draftHandler.Update)
		r.Delete("/api/drafts/{id}", draftHandler.Delete)
		r.Post("/api/drafts/{id}/publish", draftHandler.Publish)
		r.Post("/api/drafts/publish", draftHandler.PublishContent)
		r.Post("/api/drafts/schedule", draftHandler.ScheduleDraft)

		// History routes
		r.Get("/api/history", historyHandler.GetAll)
		r.Get("/api/history/{id}", historyHandler.GetByID)
		r.Post("/api/history", historyHandler.AddToHistory)
		r.Delete("/api/history/{id}", historyHandler.Delete)
		r.Delete("/api/history", historyHandler.ClearAll)
		r.Get("/api/history/stats", historyHandler.GetStats)

		// Generate routes
		r.Post("/api/generate/article", generateHandler.GenerateArticle)
		r.Post("/api/generate/image", generateHandler.GenerateImage)

		// API Keys routes
		r.Get("/api/api-keys", apiKeyHandler.GetAll)
		r.Get("/api/api-keys/{id}", apiKeyHandler.GetByID)
		r.Post("/api/api-keys", apiKeyHandler.Create)
		r.Put("/api/api-keys/{id}", apiKeyHandler.Update)
		r.Get("/api/models", modelHandler.GetAll)
		r.Get("/api/models/with-status", modelHandler.GetAllWithStatus)

		r.Get("/api/models/with-status", modelHandler.GetAllWithStatus)

		r.Get("/api/providers", providerHandler.GetAll)
		r.Get("/api/providers/{id}", providerHandler.GetByID)
		r.Post("/api/providers", providerHandler.Create)
		r.Put("/api/providers/{id}", providerHandler.Update)
		r.Delete("/api/providers/{id}", providerHandler.Delete)
	})

	return r
}
