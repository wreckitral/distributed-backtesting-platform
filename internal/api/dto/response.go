package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/wreckitral/distributed-backtesting-platform/internal/domain"
)

type BacktestResponse struct {
	ID             uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	StrategyID     string    `json:"strategy_id" example:"SMA Crossover"`
	Symbol         string    `json:"symbol" example:"AAPL"`
	StartDate      string    `json:"start_date" example:"2024-01-01"`
	EndDate        string    `json:"end_date" example:"2024-12-31"`
	InitialCapital float64   `json:"initial_capital" example:"10000"`
	Status         string    `json:"status" example:"completed"`
	CreatedAt      time.Time `json:"created_at" example:"2025-01-15T10:30:00Z"`
	UpdatedAt      time.Time `json:"updated_at" example:"2025-01-15T10:35:00Z"`
}

type MetricsResponse struct {
	BacktestID    uuid.UUID `json:"backtest_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	TotalReturn   float64   `json:"total_return" example:"2500.00"`
	ReturnPct     float64   `json:"return_pct" example:"25.00"`
	TotalTrades   int       `json:"total_trades" example:"8"`
	WinningTrades int       `json:"winning_trades" example:"6"`
	LosingTrades  int       `json:"losing_trades" example:"2"`
	WinRate       float64   `json:"win_rate" example:"75.00"`
	SharpeRatio   float64   `json:"sharpe_ratio" example:"1.82"`
	MaxDrawdown   float64   `json:"max_drawdown" example:"15.50"`
	ProfitFactor  float64   `json:"profit_factor" example:"4.33"`
}

type TradeResponse struct {
	ID        uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Symbol    string    `json:"symbol" example:"AAPL"`
	Direction string    `json:"direction" example:"BUY"`
	Quantity  float64   `json:"quantity" example:"54.79"`
	Price     float64   `json:"price" example:"182.50"`
	PnL       float64   `json:"pnl" example:"1000.00"`
	Timestamp time.Time `json:"timestamp" example:"2024-01-15T09:30:00Z"`
}

type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid request"`
	Message string `json:"message,omitempty" example:"strategy field is required"`
}

type SuccessResponse struct {
	Message string `json:"message" example:"Backtest created successfully"`
	Data    any    `json:"data,omitempty"`
}

type ListResponse struct {
	Items any `json:"items"`
	Total int `json:"total" example:"42"`
	Page  int `json:"page" example:"1"`
	Limit int `json:"limit" example:"20"`
}

func FromDomainBacktest(b *domain.Backtest) BacktestResponse {
	return BacktestResponse{
		ID:             b.ID,
		StrategyID:     b.StrategyID,
		Symbol:         b.Symbol,
		StartDate:      b.StartDate.Format("2006-01-02"),
		EndDate:        b.EndDate.Format("2006-01-02"),
		InitialCapital: b.InitialCapital,
		Status:         b.Status.String(),
		CreatedAt:      b.CreatedAt,
		UpdatedAt:      b.UpdatedAt,
	}
}

func FromDomainTrade(t *domain.Trade) TradeResponse {
	return TradeResponse{
		ID:        t.ID,
		Symbol:    t.Symbol,
		Direction: t.Direction.String(),
		Quantity:  t.Quantity,
		Price:     t.Price,
		PnL:       t.PnL,
		Timestamp: t.Timestamp,
	}
}

func FromDomainMetrics(m *domain.Metrics) MetricsResponse {
	return MetricsResponse{
		BacktestID:    m.BacktestID,
		TotalReturn:   m.TotalReturn,
		ReturnPct:     m.AnnualizedReturn * 100, // Convert to percentage
		TotalTrades:   m.TotalTrades,
		WinningTrades: m.WinningTrades,
		LosingTrades:  m.LosingTrades,
		WinRate:       m.WinRate * 100, // Convert to percentage if stored as decimal
		SharpeRatio:   m.SharpeRatio,
		MaxDrawdown:   m.MaxDrawdown * 100, // Convert to percentage if stored as decimal
		ProfitFactor:  m.ProfitFactor,
	}
}
