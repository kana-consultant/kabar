package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"seo-backend/internal/config"
	"seo-backend/internal/database"
	"seo-backend/internal/routes"
	"seo-backend/internal/scheduler"
)

var redisScheduler *scheduler.RedisScheduler

func InitScheduler() {
	redisScheduler = scheduler.NewRedisScheduler(database.RedisClient)
	redisScheduler.Start()
}
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.Load()

	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Initialize Redis
	if err := database.InitRedis(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer database.CloseRedis()
	r := routes.SetupRoutes(cfg)

	InitScheduler()

	log.Printf("Server starting on port %s", cfg.ServerPort)
	log.Printf("Dashboard stats: http://localhost:%s/api/dashboard/stats", cfg.ServerPort)

	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
