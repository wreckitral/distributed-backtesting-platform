package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	goose "github.com/pressly/goose/v3"
	"github.com/wreckitral/distributed-backtesting-platform/internal/config"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env file not found")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.DBHost, cfg.Database.DBPort, cfg.Database.DBUser,
		cfg.Database.DBPassword, cfg.Database.DBName,
	)

	if err := goose.SetDialect("postgres"); err != nil {
        log.Fatalf("Failed to set dialect: %v", err)
    }

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	migrationsDir := "scripts/migrations"

	if err := goose.Up(db, migrationsDir); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully")

}
