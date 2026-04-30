package main

import (
	"log"
	"seo-backend/internal/config"
	"seo-backend/internal/database"
)

func main() {
	// Load config
	cfg := config.Load()

	// Connect database
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect database:", err)
	}
	defer database.Close()

	// Run seed
	if err := database.RunSeed(); err != nil {
		log.Fatal("Failed to seed database:", err)
	}

	log.Println("Seeding completed successfully")
}