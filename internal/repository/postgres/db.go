package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/wreckitral/distributed-backtesting-platform/internal/config"
)

func initDB(cfg config.Database) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, fmt.Errorf("error opening db: %w", err)
    }

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	// this is critical for load balancer which silently drop idle connections
	// after a few minutes
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	// verify connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := db.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("error connecting to db: %w", err)
    }

    return db, nil
}

