package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/wreckitral/distributed-backtesting-platform/internal/domain"
)

type StrategyRepository interface {
	Create(ctx context.Context, strategy *domain.Strategy) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Strategy, error)
	Update(ctx context.Context, strategy *domain.Strategy) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*domain.Strategy, error)
}

type BacktestRepository interface {
	Create(ctx context.Context, backtest *domain.Backtest) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Backtest, error)
	Update(ctx context.Context, backtest *domain.Backtest) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.BacktestStatus) error
	MarkAsCompleted(ctx context.Context, id uuid.UUID) error
	MarkAsFailed(ctx context.Context, id uuid.UUID, errorMsg string) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*domain.Backtest, error)
	ListByStatus(ctx context.Context, status domain.BacktestStatus) ([]*domain.Backtest, error)
	ListByStrategy(ctx context.Context, strategyID string) ([]*domain.Backtest, error)
}

type TradeRepository interface {
	Create(ctx context.Context, trade *domain.Trade) error
	CreateBatch(ctx context.Context, trades []*domain.Trade) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Trade, error)
	GetByBacktestID(ctx context.Context, backtestID uuid.UUID) ([]*domain.Trade, error)
	ListByBacktest(ctx context.Context, backtestID uuid.UUID) ([]*domain.Trade, error)
	DeleteByBacktest(ctx context.Context, backtestID uuid.UUID) error
	CountByBacktest(ctx context.Context, backtestID uuid.UUID) (int, error)
	GetTradeStats(ctx context.Context, backtestID uuid.UUID) (winning, losing int, err error)
}

type MetricsRepository interface {
	Create(ctx context.Context, metrics *domain.Metrics) error
	GetByBacktestID(ctx context.Context, backtestID uuid.UUID) (*domain.Metrics, error)
	Update(ctx context.Context, metrics *domain.Metrics) error
	Delete(ctx context.Context, backtestID uuid.UUID) error
	Exists(ctx context.Context, backtestID uuid.UUID) (bool, error)
	ListTopPerformers(ctx context.Context, limit int) ([]*domain.Metrics, error)
}
