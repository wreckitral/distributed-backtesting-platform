package domain

import (
	"time"

	"github.com/google/uuid"
)

type Metrics struct {
	BacktestID          uuid.UUID
	TotalReturn         float64
	AnnualizedReturn    float64
	SharpeRatio         float64
	MaxDrawdown         float64
	MaxDrawdownDuration int
	WinRate             float64
	TotalTrades         int
	WinningTrades       int
	LosingTrades        int
	ProfitFactor        float64
	AvgWin              float64
	AvgLoss             float64
	LargestWin          float64
	LargestLoss         float64
}

type EquityCurve struct {
	Timestamp time.Time
	Equity    float64
}
