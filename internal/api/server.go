package api

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/wreckitral/distributed-backtesting-platform/internal/api/handlers"
	"github.com/wreckitral/distributed-backtesting-platform/internal/marketdata"
	"github.com/wreckitral/distributed-backtesting-platform/internal/repository/postgres"
)

type Server struct {
	router *gin.Engine
	db     *sql.DB
}

func NewServer(db *sql.DB, dataDir string) (*Server, error) {
	// Set Gin mode (release/debug)
	gin.SetMode(gin.ReleaseMode)

	// Create router
	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Initialize repositories
	backtestRepo := postgres.NewBacktestRepository(db)
	tradeRepo := postgres.NewTradeRepository(db)
	metricsRepo := postgres.NewMetricsRepository(db)
	// strategyRepo := postgres.NewStrategyRepository(db) // TODO: Will be used in Day 7 for strategy listing

	// Initialize market data provider
	provider, err := marketdata.NewCSVProvider(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create market data provider: %w", err)
	}

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	backtestHandler := handlers.NewBacktestHandler(
		backtestRepo,
		tradeRepo,
		metricsRepo,
		provider,
	)

	// Register routes
	registerRoutes(router, healthHandler, backtestHandler)

	return &Server{
		router: router,
		db:     db,
	}, nil
}

func registerRoutes(
	router *gin.Engine,
	healthHandler *handlers.HealthHandler,
	backtestHandler *handlers.BacktestHandler,
) {
	// Health check
	router.GET("/health", healthHandler.GetHealth)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Backtest routes
		backtests := v1.Group("/backtests")
		{
			backtests.POST("", backtestHandler.CreateBacktest)
			backtests.GET("", backtestHandler.ListBacktests)
			backtests.GET("/:id", backtestHandler.GetBacktest)
			backtests.GET("/:id/metrics", backtestHandler.GetBacktestMetrics)
			backtests.GET("/:id/trades", backtestHandler.GetBacktestTrades)
			backtests.DELETE("/:id", backtestHandler.DeleteBacktest)
		}

		// TODO: Day 7 - Strategy routes
		// strategies := v1.Group("/strategies")
		// {
		//     strategies.GET("", strategyHandler.ListStrategies)
		// }

		// TODO: Day 7 - Symbol routes
		// symbols := v1.Group("/symbols")
		// {
		//     symbols.GET("", symbolHandler.ListSymbols)
		// }
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (s *Server) Run(port string) error {
	log.Printf("Starting API server on port %s", port)
	log.Printf("Swagger docs available at http://localhost%s/swagger/index.html", port)
	return s.router.Run(port)
}

func (s *Server) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
