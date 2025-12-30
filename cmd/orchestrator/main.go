package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
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

	fmt.Println("Orchestrator starting...")
	fmt.Printf("Environment: %s\n", cfg.Environment)
	fmt.Printf("Log Level: %s\n", cfg.LogLevel)
	fmt.Printf("Port: %s\n", cfg.OrchestratorPort)
	fmt.Printf("Database: %s@%s:%s/%s\n",
		cfg.Database.DBUser,
		cfg.Database.DBHost,
		cfg.Database.DBPort,
		cfg.Database.DBName,
	)
	fmt.Printf("Worker Pool Size: %d\n", cfg.Worker.WorkerPoolSize)
}
