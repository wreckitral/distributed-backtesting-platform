package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/wreckitral/distributed-backtesting-platform/internal/domain"
)

type StrategyRepository interface {
	Create(ctx context.Context, strategy *domain.Strategy) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Strategy, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Strategy, error)
}

type BacktestRepository interface {
	Create(ctx context.Context, backtest *domain.Backtest) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Backtest, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, resultJSON string) error
}

type TradeRepository interface {
	CreateBatch(ctx context.Context, trades []*domain.Trade) error
	ListByBacktest(ctx context.Context, backtestID uuid.UUID) ([]*domain.Trade, error)
}


