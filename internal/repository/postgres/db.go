package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/wreckitral/distributed-backtesting-platform/internal/config"
)

func NewPostgresDB(cfg config.Database) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	maxOpen := cfg.MaxOpenConns
	if maxOpen == 0 {
		maxOpen = 25
	}
	db.SetMaxOpenConns(maxOpen)

	maxIdle := cfg.MaxIdleConns
	if maxIdle == 0 {
		maxIdle = 5
	}
	db.SetMaxIdleConns(maxIdle)

	maxLifetime := cfg.ConnMaxLifetime
	if maxLifetime == 0 {
		maxLifetime = 5 * time.Minute
	}
	db.SetConnMaxLifetime(maxLifetime)

	maxIdleTime := cfg.ConnMaxIdleTime
	if maxIdleTime == 0 {
		maxIdleTime = 2 * time.Minute
	}
	db.SetConnMaxIdleTime(maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return db, nil
}

func Close(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func HealthCheck(ctx context.Context, db *sql.DB) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return db.PingContext(ctx)
}
