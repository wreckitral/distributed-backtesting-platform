-- +goose Up
-- +goose StatementBegin
-- Create strategies table first (backtests references it)
CREATE TABLE IF NOT EXISTS strategies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    language VARCHAR(20) NOT NULL,
    code TEXT NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create backtests table
CREATE TABLE IF NOT EXISTS backtests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    strategy_id UUID NOT NULL REFERENCES strategies(id),
    symbol VARCHAR(10) NOT NULL,
    status VARCHAR(20) NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    initial_capital DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    error_message TEXT
);

-- Create trades table
CREATE TABLE IF NOT EXISTS trades (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    backtest_id UUID NOT NULL REFERENCES backtests(id) ON DELETE CASCADE,
    symbol VARCHAR(10) NOT NULL,
    direction VARCHAR(10) NOT NULL,
    quantity DOUBLE PRECISION NOT NULL,
    price DOUBLE PRECISION NOT NULL,
    commission DOUBLE PRECISION NOT NULL DEFAULT 0,
    timestamp TIMESTAMPTZ NOT NULL,
    pnl DOUBLE PRECISION NOT NULL,
    cumulative_pnl DOUBLE PRECISION NOT NULL
);

-- Create metrics table
CREATE TABLE IF NOT EXISTS metrics (
    backtest_id UUID PRIMARY KEY REFERENCES backtests(id) ON DELETE CASCADE,
    total_return DOUBLE PRECISION,
    annualized_return DOUBLE PRECISION,
    sharpe_ratio DOUBLE PRECISION,
    max_drawdown DOUBLE PRECISION,
    max_drawdown_duration INTEGER,
    win_rate DOUBLE PRECISION,
    total_trades INTEGER,
    winning_trades INTEGER,
    losing_trades INTEGER,
    profit_factor DOUBLE PRECISION,
    avg_win DOUBLE PRECISION,
    avg_loss DOUBLE PRECISION,
    largest_win DOUBLE PRECISION,
    largest_loss DOUBLE PRECISION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS metrics;
DROP TABLE IF EXISTS trades;
DROP TABLE IF EXISTS backtests;
DROP TABLE IF EXISTS strategies;
-- +goose StatementEnd
