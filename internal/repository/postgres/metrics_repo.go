package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/wreckitral/distributed-backtesting-platform/internal/domain"
)

type metricsRepository struct {
	db *sql.DB
}

func NewMetricsRepository(db *sql.DB) *metricsRepository {
	return &metricsRepository{db: db}
}

func (r *metricsRepository) Create(ctx context.Context, metrics *domain.Metrics) error {
	query := `
		INSERT INTO metrics (
			backtest_id, total_return, annualized_return, sharpe_ratio,
			max_drawdown, max_drawdown_duration, win_rate, total_trades,
			winning_trades, losing_trades, profit_factor, avg_win, avg_loss,
			largest_win, largest_loss
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

	_, err := r.db.ExecContext(
		ctx,
		query,
		metrics.BacktestID,
		metrics.TotalReturn,
		metrics.AnnualizedReturn,
		metrics.SharpeRatio,
		metrics.MaxDrawdown,
		metrics.MaxDrawdownDuration,
		metrics.WinRate,
		metrics.TotalTrades,
		metrics.WinningTrades,
		metrics.LosingTrades,
		metrics.ProfitFactor,
		metrics.AvgWin,
		metrics.AvgLoss,
		metrics.LargestWin,
		metrics.LargestLoss,
	)

	if err != nil {
		return fmt.Errorf("failed to insert metrics: %w", err)
	}

	return nil
}

func (r *metricsRepository) GetByBacktestID(ctx context.Context, backtestID uuid.UUID) (*domain.Metrics, error) {
	query := `
		SELECT backtest_id, total_return, annualized_return, sharpe_ratio,
		       max_drawdown, max_drawdown_duration, win_rate, total_trades,
		       winning_trades, losing_trades, profit_factor, avg_win, avg_loss,
		       largest_win, largest_loss
		FROM metrics
		WHERE backtest_id = $1`

	metrics := &domain.Metrics{}

	err := r.db.QueryRowContext(ctx, query, backtestID).Scan(
		&metrics.BacktestID,
		&metrics.TotalReturn,
		&metrics.AnnualizedReturn,
		&metrics.SharpeRatio,
		&metrics.MaxDrawdown,
		&metrics.MaxDrawdownDuration,
		&metrics.WinRate,
		&metrics.TotalTrades,
		&metrics.WinningTrades,
		&metrics.LosingTrades,
		&metrics.ProfitFactor,
		&metrics.AvgWin,
		&metrics.AvgLoss,
		&metrics.LargestWin,
		&metrics.LargestLoss,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("metrics not found for backtest")
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching metrics: %w", err)
	}

	return metrics, nil
}

func (r *metricsRepository) Update(ctx context.Context, metrics *domain.Metrics) error {
	query := `
		UPDATE metrics
		SET total_return = $1, annualized_return = $2, sharpe_ratio = $3,
		    max_drawdown = $4, max_drawdown_duration = $5, win_rate = $6,
		    total_trades = $7, winning_trades = $8, losing_trades = $9,
		    profit_factor = $10, avg_win = $11, avg_loss = $12,
		    largest_win = $13, largest_loss = $14
		WHERE backtest_id = $15`

	result, err := r.db.ExecContext(
		ctx,
		query,
		metrics.TotalReturn,
		metrics.AnnualizedReturn,
		metrics.SharpeRatio,
		metrics.MaxDrawdown,
		metrics.MaxDrawdownDuration,
		metrics.WinRate,
		metrics.TotalTrades,
		metrics.WinningTrades,
		metrics.LosingTrades,
		metrics.ProfitFactor,
		metrics.AvgWin,
		metrics.AvgLoss,
		metrics.LargestWin,
		metrics.LargestLoss,
		metrics.BacktestID,
	)

	if err != nil {
		return fmt.Errorf("failed to update metrics: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("metrics not found for backtest")
	}

	return nil
}

func (r *metricsRepository) Delete(ctx context.Context, backtestID uuid.UUID) error {
	query := `DELETE FROM metrics WHERE backtest_id = $1`

	result, err := r.db.ExecContext(ctx, query, backtestID)
	if err != nil {
		return fmt.Errorf("failed to delete metrics: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("metrics not found for backtest")
	}

	return nil
}

func (r *metricsRepository) Exists(ctx context.Context, backtestID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM metrics WHERE backtest_id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, backtestID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking metrics existence: %w", err)
	}

	return exists, nil
}

func (r *metricsRepository) ListTopPerformers(ctx context.Context, limit int) ([]*domain.Metrics, error) {
	query := `
		SELECT backtest_id, total_return, annualized_return, sharpe_ratio,
		       max_drawdown, max_drawdown_duration, win_rate, total_trades,
		       winning_trades, losing_trades, profit_factor, avg_win, avg_loss,
		       largest_win, largest_loss
		FROM metrics
		ORDER BY sharpe_ratio DESC
		LIMIT $1`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("error listing top performers: %w", err)
	}
	defer rows.Close()

	var metricsList []*domain.Metrics
	for rows.Next() {
		metrics := &domain.Metrics{}

		if err := rows.Scan(
			&metrics.BacktestID,
			&metrics.TotalReturn,
			&metrics.AnnualizedReturn,
			&metrics.SharpeRatio,
			&metrics.MaxDrawdown,
			&metrics.MaxDrawdownDuration,
			&metrics.WinRate,
			&metrics.TotalTrades,
			&metrics.WinningTrades,
			&metrics.LosingTrades,
			&metrics.ProfitFactor,
			&metrics.AvgWin,
			&metrics.AvgLoss,
			&metrics.LargestWin,
			&metrics.LargestLoss,
		); err != nil {
			return nil, fmt.Errorf("error scanning metrics: %w", err)
		}

		metricsList = append(metricsList, metrics)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating metrics: %w", err)
	}

	return metricsList, nil
}
