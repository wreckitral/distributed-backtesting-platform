package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Database         Database
	OrchestratorPort string
	Worker           Worker
	LogLevel         string
	Environment      string
}

type Database struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Connection pool settings
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type Worker struct {
	WorkerPoolSize int
	PythonPath     string
}

func Load() (*Config, error) {
	cfg := &Config{
		Database: Database{
			DBHost:     os.Getenv("DB_HOST"),
			DBPort:     os.Getenv("DB_PORT"),
			DBUser:     os.Getenv("DB_USER"),
			DBPassword: os.Getenv("DB_PASSWORD"),
			DBName:     os.Getenv("DB_NAME"),

			// Optional pool settings (with defaults)
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			ConnMaxIdleTime: getEnvAsDuration("DB_CONN_MAX_IDLE_TIME", 2*time.Minute),
		},
		OrchestratorPort: os.Getenv("ORCHESTRATOR_PORT"),
		Worker: Worker{
			WorkerPoolSize: getEnvAsInt("WORKER_POOL_SIZE", 4),
			PythonPath:     os.Getenv("PYTHON_PATH"),
		},
		LogLevel:    os.Getenv("LOG_LEVEL"),
		Environment: os.Getenv("ENVIRONMENT"),
	}

	setDefaults(cfg)

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func setDefaults(c *Config) {
	if c.Database.DBHost == "" {
		c.Database.DBHost = "localhost"
	}
	if c.Database.DBPort == "" {
		c.Database.DBPort = "5432"
	}
	if c.OrchestratorPort == "" {
		c.OrchestratorPort = "8080"
	}
	if c.Worker.PythonPath == "" {
		c.Worker.PythonPath = "/usr/bin/python3"
	}
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}
	if c.Environment == "" {
		c.Environment = "development"
	}
}

func (c *Config) Validate() error {
	if c.Database.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.Database.DBUser == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.Worker.WorkerPoolSize < 1 {
		return fmt.Errorf("WORKER_POOL_SIZE must be at least 1")
	}

	validLogLevels := []string{"debug", "info", "warn", "error"}
	if !contains(validLogLevels, c.LogLevel) {
		return fmt.Errorf("invalid LOG_LEVEL: %s", c.LogLevel)
	}

	return nil
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}

func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := time.ParseDuration(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
