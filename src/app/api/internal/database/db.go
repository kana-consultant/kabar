package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"seo-backend/internal/config"
)

var db *sql.DB

func Connect(cfg *config.Config) error {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	var err error
	db, err = sql.Open("postgres", connStr)
	log.Printf(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connected successfully")

	// Jalankan migrasi
	if err := RunMigrations(); err != nil {
		log.Printf("Migration warning: %v", err)
	}

	

	return nil
}

func GetDB() *sql.DB {
	return db
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}