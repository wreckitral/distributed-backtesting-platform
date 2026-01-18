package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/wreckitral/distributed-backtesting-platform/internal/domain"
)

type tradeRepository struct {
	db *sql.DB
}

func NewTradeRepository(db *sql.DB) *tradeRepository {
	return &tradeRepository{db: db}
}

func (r *tradeRepository) Create(ctx context.Context, trade *domain.Trade) error {
	query := `
		INSERT INTO trades (
			backtest_id, symbol, direction, quantity, price,
			commission, timestamp, pnl, cumulative_pnl
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	err := r.db.QueryRowContext(
		ctx,
		query,
		trade.BacktestID,
		trade.Symbol,
		trade.Direction.String(),
		trade.Quantity,
		trade.Price,
		trade.Commission,
		trade.Timestamp,
		trade.PnL,
		trade.CumulativePnL,
	).Scan(&trade.ID)

	if err != nil {
		return fmt.Errorf("failed to insert trade: %w", err)
	}

	return nil
}

func (r *tradeRepository) CreateBatch(ctx context.Context, trades []*domain.Trade) error {
	if len(trades) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	valueStrings := make([]string, 0, len(trades))
	valueArgs := make([]interface{}, 0, len(trades)*9)

	for i, trade := range trades {
		valueStrings = append(valueStrings, fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*9+1, i*9+2, i*9+3, i*9+4, i*9+5, i*9+6, i*9+7, i*9+8, i*9+9,
		))

		valueArgs = append(valueArgs,
			trade.BacktestID,
			trade.Symbol,
			trade.Direction.String(),
			trade.Quantity,
			trade.Price,
			trade.Commission,
			trade.Timestamp,
			trade.PnL,
			trade.CumulativePnL,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO trades (
			backtest_id, symbol, direction, quantity, price,
			commission, timestamp, pnl, cumulative_pnl
		)
		VALUES %s
		RETURNING id`,
		strings.Join(valueStrings, ","),
	)

	rows, err := tx.QueryContext(ctx, query, valueArgs...)
	if err != nil {
		return fmt.Errorf("failed to bulk insert trades: %w", err)
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		if i >= len(trades) {
			return fmt.Errorf("more rows returned than trades inserted")
		}
		if err := rows.Scan(&trades[i].ID); err != nil {
			return fmt.Errorf("failed to scan trade ID: %w", err)
		}
		i++
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating trade IDs: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *tradeRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Trade, error) {
	query := `
		SELECT id, backtest_id, symbol, direction, quantity, price,
		       commission, timestamp, pnl, cumulative_pnl
		FROM trades
		WHERE id = $1`

	trade := &domain.Trade{}
	var directionStr string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&trade.ID,
		&trade.BacktestID,
		&trade.Symbol,
		&directionStr,
		&trade.Quantity,
		&trade.Price,
		&trade.Commission,
		&trade.Timestamp,
		&trade.PnL,
		&trade.CumulativePnL,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("trade not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching trade: %w", err)
	}

	trade.Direction = parseDirection(directionStr)
	return trade, nil
}

func (r *tradeRepository) GetByBacktestID(ctx context.Context, backtestID uuid.UUID) ([]*domain.Trade, error) {
	query := `
        SELECT id, backtest_id, symbol, direction, quantity, price,
               commission, timestamp, pnl, cumulative_pnl
        FROM trades
        WHERE backtest_id = $1
        ORDER BY timestamp ASC
    `

	rows, err := r.db.QueryContext(ctx, query, backtestID)
	if err != nil {
		return nil, fmt.Errorf("failed to query trades: %w", err)
	}
	defer rows.Close()

	var trades []*domain.Trade
	for rows.Next() {
		trade := &domain.Trade{}
		err := rows.Scan(
			&trade.ID,
			&trade.BacktestID,
			&trade.Symbol,
			&trade.Direction,
			&trade.Quantity,
			&trade.Price,
			&trade.Commission,
			&trade.Timestamp,
			&trade.PnL,
			&trade.CumulativePnL,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trade: %w", err)
		}
		trades = append(trades, trade)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating trades: %w", err)
	}

	return trades, nil
}

func (r *tradeRepository) ListByBacktest(ctx context.Context, backtestID uuid.UUID) ([]*domain.Trade, error) {
	query := `
		SELECT id, backtest_id, symbol, direction, quantity, price,
		       commission, timestamp, pnl, cumulative_pnl
		FROM trades
		WHERE backtest_id = $1
		ORDER BY timestamp ASC`

	rows, err := r.db.QueryContext(ctx, query, backtestID)
	if err != nil {
		return nil, fmt.Errorf("error listing trades: %w", err)
	}
	defer rows.Close()

	var trades []*domain.Trade
	for rows.Next() {
		trade := &domain.Trade{}
		var directionStr string

		if err := rows.Scan(
			&trade.ID,
			&trade.BacktestID,
			&trade.Symbol,
			&directionStr,
			&trade.Quantity,
			&trade.Price,
			&trade.Commission,
			&trade.Timestamp,
			&trade.PnL,
			&trade.CumulativePnL,
		); err != nil {
			return nil, fmt.Errorf("error scanning trade: %w", err)
		}

		trade.Direction = parseDirection(directionStr)
		trades = append(trades, trade)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating trades: %w", err)
	}

	return trades, nil
}

func (r *tradeRepository) DeleteByBacktest(ctx context.Context, backtestID uuid.UUID) error {
	query := `DELETE FROM trades WHERE backtest_id = $1`

	result, err := r.db.ExecContext(ctx, query, backtestID)
	if err != nil {
		return fmt.Errorf("failed to delete trades: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	_ = rowsAffected

	return nil
}

func (r *tradeRepository) CountByBacktest(ctx context.Context, backtestID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM trades WHERE backtest_id = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, backtestID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting trades: %w", err)
	}

	return count, nil
}

func (r *tradeRepository) GetTradeStats(ctx context.Context, backtestID uuid.UUID) (winning, losing int, err error) {
	query := `
		SELECT
			COUNT(CASE WHEN pnl > 0 THEN 1 END) as winning_trades,
			COUNT(CASE WHEN pnl <= 0 THEN 1 END) as losing_trades
		FROM trades
		WHERE backtest_id = $1`

	err = r.db.QueryRowContext(ctx, query, backtestID).Scan(&winning, &losing)
	if err != nil {
		return 0, 0, fmt.Errorf("error fetching trade stats: %w", err)
	}

	return winning, losing, nil
}

func parseDirection(d string) domain.TradeDirection {
	switch d {
	case "BUY":
		return domain.TradeDirectionBuy
	case "SELL":
		return domain.TradeDirectionSell
	default:
		return domain.TradeDirectionBuy
	}
}
