package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"

	"seo-backend/internal/config"
	"seo-backend/internal/container"
	"seo-backend/internal/database"
	"seo-backend/internal/scheduler"
	services "seo-backend/internal/service"

	_ "seo-backend/docs"
)

// @title SEO Backend API
// @version 1.0
// @description API Documentation for SEO Backend Service
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@seo-backend.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

var (
	redisScheduler *scheduler.RedisScheduler
	appContainer   *container.Container
)

func main() {
	// 1. LOAD ENVIRONMENT
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// 2. LOAD CONFIGURATION
	cfg := config.Load()

	// 3. INITIALIZE DATABASE
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// 4. INITIALIZE REDIS
	if err := database.InitRedis(cfg); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer database.CloseRedis()

	smtpConfig := config.LoadSMTPConfig()
	emailService := services.NewSMTPEmailService(smtpConfig)

	// 6. INITIALIZE CONTAINER
	appContainer = container.NewContainer(cfg, database.GetDB(), database.RedisClient, emailService)

	// 7. SETUP SWAGGER ROUTE (jika belum di container)
	setupSwagger(appContainer)

	// 8. CREATE HTTP SERVER
	server := &http.Server{
		Addr:         ":8080",
		Handler:      appContainer.Router,
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
		IdleTimeout:  120 * time.Second,
	}

	// 9. START SERVER
	startServerWithGracefulShutdown(server, cfg)
}

// setupSwagger adds Swagger UI route
func setupSwagger(container *container.Container) {
	// Jika container belum punya route Swagger
	container.Router.Get("/swagger/*", httpSwagger.WrapHandler)

	// Atau jika ingin custom URL
	// container.Router.Get("/swagger/*", httpSwagger.Handler(
	//     httpSwagger.URL("/api/docs/swagger.json"),
	// ))
}

func startServerWithGracefulShutdown(server *http.Server, cfg *config.Config) {
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		log.Printf("Swagger UI: http://localhost:%s/swagger/index.html", cfg.ServerPort)
		log.Printf("Health check: http://localhost:%s/health", cfg.ServerPort)
		log.Printf("API: http://localhost:%s/api", cfg.ServerPort)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatal("Server failed to start:", err)

	case sig := <-shutdown:
		log.Printf("Received signal: %v, starting graceful shutdown...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			if err := server.Close(); err != nil {
				log.Printf("Forced shutdown failed: %v", err)
			}
		}

		log.Println("Server stopped gracefully")
	}
}
