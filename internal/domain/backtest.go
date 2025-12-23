package domain

import (
	"time"

	"github.com/google/uuid"
)

type BacktestStatus int

const (
	StatusPending BacktestStatus = iota
	StatusQueued
	StatusRunning
	StatusCompleted
	StatusFailed
)

func (b BacktestStatus) String() string {
	switch b {
	case StatusPending:
		return "PENDING"
	case StatusQueued:
		return "QUEUED"
	case StatusRunning:
		return "RUNNING"
	case StatusCompleted:
		return "COMPLETED"
	case StatusFailed:
		return "FAILED"
	default:
		return "UNKNOWN"
	}
}

type Backtest struct {
	ID 				uuid.UUID
	StrategyID 		string
	Status 			BacktestStatus
	Symbol 			string
	StartDate 		time.Time
	EndDate 		time.Time
	InitialCapital 	float64
	CreatedAt 		time.Time
	UpdatedAt 		time.Time
	CompletedAt 	*time.Time
	ErrorMessage 	string
}

