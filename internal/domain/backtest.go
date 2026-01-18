package domain

import (
	"time"

	"github.com/google/uuid"
)

type Backtest struct {
	ID             uuid.UUID
	StrategyID     string
	Status         BacktestStatus
	Symbol         string
	StartDate      time.Time
	EndDate        time.Time
	InitialCapital float64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CompletedAt    *time.Time
	ErrorMessage   string
}

type BacktestStatus int

const (
	BacktestStatusPending BacktestStatus = iota
	BacktestStatusQueued
	BacktestStatusRunning
	BacktestStatusCompleted
	BacktestStatusFailed
)

func (b BacktestStatus) String() string {
	switch b {
	case BacktestStatusPending:
		return "PENDING"
	case BacktestStatusQueued:
		return "QUEUED"
	case BacktestStatusRunning:
		return "RUNNING"
	case BacktestStatusCompleted:
		return "COMPLETED"
	case BacktestStatusFailed:
		return "FAILED"
	default:
		return "UNKNOWN"
	}
}
