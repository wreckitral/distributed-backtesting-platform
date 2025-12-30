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

	fmt.Println("Worker starting...")
	fmt.Printf("Worker Pool Size: %d\n", cfg.Worker.WorkerPoolSize)
	fmt.Printf("Python Path: %s\n", cfg.Worker.PythonPath)
	fmt.Printf("Log Level: %s\n", cfg.LogLevel)
	fmt.Printf("Environment: %s\n", cfg.Environment)
}
