-- +goose Up
-- +goose StatementBegin
-- Indexes for backtests
CREATE INDEX idx_backtests_strategy_id ON backtests(strategy_id);
CREATE INDEX idx_backtests_status ON backtests(status);
CREATE INDEX idx_backtests_symbol ON backtests(symbol);
CREATE INDEX idx_backtests_created_at ON backtests(created_at);

-- Indexes for trades
CREATE INDEX idx_trades_backtest_id ON trades(backtest_id);
CREATE INDEX idx_trades_timestamp ON trades(timestamp);
CREATE INDEX idx_trades_symbol ON trades(symbol);

-- Indexes for strategies
CREATE INDEX idx_strategies_name ON strategies(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_strategies_name;
DROP INDEX IF EXISTS idx_trades_symbol;
DROP INDEX IF EXISTS idx_trades_timestamp;
DROP INDEX IF EXISTS idx_trades_backtest_id;
DROP INDEX IF EXISTS idx_backtests_created_at;
DROP INDEX IF EXISTS idx_backtests_symbol;
DROP INDEX IF EXISTS idx_backtests_status;
DROP INDEX IF EXISTS idx_backtests_strategy_id;
-- +goose StatementEnd
