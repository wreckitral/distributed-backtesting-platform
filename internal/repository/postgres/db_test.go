package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/wreckitral/distributed-backtesting-platform/internal/config"
)

func TestNewPostgresDB(t *testing.T) {
	cfg := config.Database{
		DBHost:          "localhost",
		DBPort:          "5432",
		DBUser:          "faliux",
		DBPassword:      "faliux123",
		DBName:          "finux",
		MaxOpenConns:    10,
		MaxIdleConns:    2,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 2 * time.Minute,
	}

	db, err := NewPostgresDB(cfg)
	if err != nil {
		t.Fatalf("Failed to create database connection: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	if err := HealthCheck(ctx, db); err != nil {
		t.Fatalf("Health check failed: %v", err)
	}

	stats := db.Stats()
	t.Logf("Connection pool stats: %+v", stats)

	if stats.MaxOpenConnections != 10 {
		t.Errorf("Expected MaxOpenConnections=10, got %d", stats.MaxOpenConnections)
	}
}
