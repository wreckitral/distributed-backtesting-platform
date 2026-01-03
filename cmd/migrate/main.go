package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/wreckitral/distributed-backtesting-platform/internal/config"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.DBHost,
		cfg.Database.DBPort,
		cfg.Database.DBUser,
		cfg.Database.DBPassword,
		cfg.Database.DBName,
	)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Failed to set dialect: %v", err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to database successfully")

	migrationsDir := filepath.Join("scripts", "migrations")

	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		migrationsDir = filepath.Join("..", "..", "scripts", "migrations")
	}

	log.Printf("Using migrations directory: %s", migrationsDir)

	if err := goose.Up(db, migrationsDir); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully")
}
