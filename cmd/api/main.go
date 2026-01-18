package main

import (
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/wreckitral/distributed-backtesting-platform/internal/api"
	"github.com/wreckitral/distributed-backtesting-platform/internal/config"
	"github.com/wreckitral/distributed-backtesting-platform/internal/repository/postgres"

	// Import swagger docs
	_ "github.com/wreckitral/distributed-backtesting-platform/docs"
)

//	@title			Distributed Backtesting Platform API
//	@version		1.0
//	@description	API for distributed financial backtesting platform
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8080
// @BasePath	/
// @schemes	http https
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database using your existing config structure
	db, err := postgres.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer postgres.Close(db)

	log.Println("Database connected successfully")

	// Get data directory from environment or use default
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./data/sample"
	}

	// Create and start server
	server, err := api.NewServer(db, dataDir)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	defer server.Close()

	// Use port from config
	port := ":" + cfg.OrchestratorPort

	log.Printf("API server starting on port %s", port)
	log.Printf("Swagger UI: http://localhost%s/swagger/index.html", port)
	log.Printf("Health check: http://localhost%s/health", port)
	log.Printf("Environment: %s", cfg.Environment)

	if err := server.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
