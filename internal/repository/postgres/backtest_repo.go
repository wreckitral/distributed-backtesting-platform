package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wreckitral/distributed-backtesting-platform/internal/domain"
)

type backtestRepository struct {
	db *sql.DB
}

func NewBacktestRepository(db *sql.DB) *backtestRepository {
	return &backtestRepository{db: db}
}

func (r *backtestRepository) Create(ctx context.Context, b *domain.Backtest) error {
	if b.CreatedAt.IsZero() {
		b.CreatedAt = time.Now()
	}
	b.UpdatedAt = b.CreatedAt

	query := `
		INSERT INTO backtests (
			strategy_id, symbol, status, start_date, end_date,
			initial_capital, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	err := r.db.QueryRowContext(
		ctx,
		query,
		b.StrategyID,
		b.Symbol,
		b.Status.String(),
		b.StartDate,
		b.EndDate,
		b.InitialCapital,
		b.CreatedAt,
		b.UpdatedAt,
	).Scan(&b.ID)

	if err != nil {
		return fmt.Errorf("failed to insert backtest: %w", err)
	}

	return nil
}

func (r *backtestRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Backtest, error) {
	query := `
		SELECT id, strategy_id, symbol, status, start_date, end_date,
		       initial_capital, created_at, updated_at, completed_at, error_message
		FROM backtests
		WHERE id = $1`

	b := &domain.Backtest{}
	var statusStr string
	var completedAt sql.NullTime
	var errorMessage sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&b.ID,
		&b.StrategyID,
		&b.Symbol,
		&statusStr,
		&b.StartDate,
		&b.EndDate,
		&b.InitialCapital,
		&b.CreatedAt,
		&b.UpdatedAt,
		&completedAt,
		&errorMessage,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("backtest not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching backtest: %w", err)
	}

	b.Status = parseStatus(statusStr)

	if completedAt.Valid {
		b.CompletedAt = &completedAt.Time
	}
	if errorMessage.Valid {
		b.ErrorMessage = errorMessage.String
	}

	return b, nil
}

func (r *backtestRepository) Update(ctx context.Context, backtest *domain.Backtest) error {
	backtest.UpdatedAt = time.Now()

	query := `
		UPDATE backtests
		SET strategy_id = $1, symbol = $2, status = $3,
		    start_date = $4, end_date = $5, initial_capital = $6,
		    updated_at = $7, completed_at = $8, error_message = $9
		WHERE id = $10`

	// Handle nullable fields
	var completedAt sql.NullTime
	if backtest.CompletedAt != nil {
		completedAt = sql.NullTime{Time: *backtest.CompletedAt, Valid: true}
	}

	var errorMessage sql.NullString
	if backtest.ErrorMessage != "" {
		errorMessage = sql.NullString{String: backtest.ErrorMessage, Valid: true}
	}

	result, err := r.db.ExecContext(
		ctx,
		query,
		backtest.StrategyID,
		backtest.Symbol,
		backtest.Status.String(),
		backtest.StartDate,
		backtest.EndDate,
		backtest.InitialCapital,
		backtest.UpdatedAt,
		completedAt,
		errorMessage,
		backtest.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update backtest: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("backtest not found")
	}

	return nil
}

func (r *backtestRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.BacktestStatus) error {
	now := time.Now()

	query := `
		UPDATE backtests
		SET status = $1, updated_at = $2
		WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, status.String(), now, id)
	if err != nil {
		return fmt.Errorf("failed to update backtest status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("backtest not found")
	}

	return nil
}

func (r *backtestRepository) MarkAsCompleted(ctx context.Context, id uuid.UUID) error {
	now := time.Now()

	query := `
		UPDATE backtests
		SET status = $1, completed_at = $2, updated_at = $3
		WHERE id = $4`

	result, err := r.db.ExecContext(
		ctx,
		query,
		domain.BacktestStatusCompleted.String(),
		now,
		now,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to mark backtest as completed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("backtest not found")
	}

	return nil
}

func (r *backtestRepository) MarkAsFailed(ctx context.Context, id uuid.UUID, errorMsg string) error {
	now := time.Now()

	query := `
		UPDATE backtests
		SET status = $1, error_message = $2, completed_at = $3, updated_at = $4
		WHERE id = $5`

	result, err := r.db.ExecContext(
		ctx,
		query,
		domain.BacktestStatusFailed.String(),
		errorMsg,
		now,
		now,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to mark backtest as failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("backtest not found")
	}

	return nil
}

func (r *backtestRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM backtests WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete backtest: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("backtest not found")
	}

	return nil
}

func (r *backtestRepository) List(ctx context.Context, limit, offset int) ([]*domain.Backtest, error) {
	query := `
		SELECT id, strategy_id, symbol, status, start_date, end_date,
		       initial_capital, created_at, updated_at, completed_at, error_message
		FROM backtests
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing backtests: %w", err)
	}
	defer rows.Close()

	var backtests []*domain.Backtest
	for rows.Next() {
		b := &domain.Backtest{}
		var statusStr string
		var completedAt sql.NullTime
		var errorMessage sql.NullString

		if err := rows.Scan(
			&b.ID,
			&b.StrategyID,
			&b.Symbol,
			&statusStr,
			&b.StartDate,
			&b.EndDate,
			&b.InitialCapital,
			&b.CreatedAt,
			&b.UpdatedAt,
			&completedAt,
			&errorMessage,
		); err != nil {
			return nil, fmt.Errorf("error scanning backtest: %w", err)
		}

		b.Status = parseStatus(statusStr)

		if completedAt.Valid {
			b.CompletedAt = &completedAt.Time
		}
		if errorMessage.Valid {
			b.ErrorMessage = errorMessage.String
		}

		backtests = append(backtests, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating backtests: %w", err)
	}

	return backtests, nil
}

func (r *backtestRepository) ListByStatus(ctx context.Context, status domain.BacktestStatus) ([]*domain.Backtest, error) {
	query := `
		SELECT id, strategy_id, symbol, status, start_date, end_date,
		       initial_capital, created_at, updated_at, completed_at, error_message
		FROM backtests
		WHERE status = $1
		ORDER BY created_at ASC`

	rows, err := r.db.QueryContext(ctx, query, status.String())
	if err != nil {
		return nil, fmt.Errorf("error listing backtests by status: %w", err)
	}
	defer rows.Close()

	var backtests []*domain.Backtest
	for rows.Next() {
		b := &domain.Backtest{}
		var statusStr string
		var completedAt sql.NullTime
		var errorMessage sql.NullString

		if err := rows.Scan(
			&b.ID,
			&b.StrategyID,
			&b.Symbol,
			&statusStr,
			&b.StartDate,
			&b.EndDate,
			&b.InitialCapital,
			&b.CreatedAt,
			&b.UpdatedAt,
			&completedAt,
			&errorMessage,
		); err != nil {
			return nil, fmt.Errorf("error scanning backtest: %w", err)
		}

		b.Status = parseStatus(statusStr)

		if completedAt.Valid {
			b.CompletedAt = &completedAt.Time
		}
		if errorMessage.Valid {
			b.ErrorMessage = errorMessage.String
		}

		backtests = append(backtests, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating backtests: %w", err)
	}

	return backtests, nil
}

func (r *backtestRepository) ListByStrategy(ctx context.Context, strategyID string) ([]*domain.Backtest, error) {
	query := `
		SELECT id, strategy_id, symbol, status, start_date, end_date,
		       initial_capital, created_at, updated_at, completed_at, error_message
		FROM backtests
		WHERE strategy_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, strategyID)
	if err != nil {
		return nil, fmt.Errorf("error listing backtests by strategy: %w", err)
	}
	defer rows.Close()

	var backtests []*domain.Backtest
	for rows.Next() {
		b := &domain.Backtest{}
		var statusStr string
		var completedAt sql.NullTime
		var errorMessage sql.NullString

		if err := rows.Scan(
			&b.ID,
			&b.StrategyID,
			&b.Symbol,
			&statusStr,
			&b.StartDate,
			&b.EndDate,
			&b.InitialCapital,
			&b.CreatedAt,
			&b.UpdatedAt,
			&completedAt,
			&errorMessage,
		); err != nil {
			return nil, fmt.Errorf("error scanning backtest: %w", err)
		}

		b.Status = parseStatus(statusStr)

		if completedAt.Valid {
			b.CompletedAt = &completedAt.Time
		}
		if errorMessage.Valid {
			b.ErrorMessage = errorMessage.String
		}

		backtests = append(backtests, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating backtests: %w", err)
	}

	return backtests, nil
}

func parseStatus(s string) domain.BacktestStatus {
	switch s {
	case "PENDING":
		return domain.BacktestStatusPending
	case "QUEUED":
		return domain.BacktestStatusQueued
	case "RUNNING":
		return domain.BacktestStatusRunning
	case "COMPLETED":
		return domain.BacktestStatusCompleted
	case "FAILED":
		return domain.BacktestStatusFailed
	default:
		return domain.BacktestStatusPending
	}
}
