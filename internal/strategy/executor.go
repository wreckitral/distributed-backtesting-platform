package strategy

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wreckitral/distributed-backtesting-platform/internal/domain"
	"github.com/wreckitral/distributed-backtesting-platform/internal/marketdata"
)

type Executor struct {
	strategy    Strategy
	provider    marketdata.Provider
	initialCash float64
}

func NewExecutor(strategy Strategy, provider marketdata.Provider, initialCash float64) *Executor {
	return &Executor{
		strategy:    strategy,
		provider:    provider,
		initialCash: initialCash,
	}
}

func (e *Executor) Run(ctx context.Context, symbol string, start, end time.Time) ([]domain.Trade, error) {
	bars, err := e.provider.GetBars(ctx, symbol, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get bars: %w", err)
	}

	if len(bars) == 0 {
		return nil, fmt.Errorf("no bars found for %s between %s and %s", symbol, start, end)
	}

	// initialize tracking variables
	var position *Position = nil
	cash := e.initialCash
	trades := []domain.Trade{}

	for i, bar := range bars {
		strategyCtx := &Context{
			Symbol:          symbol,
			CurrentBar:      bar,
			HistoricalBars:  bars[0:i], // all bars before today
			CurrentPosition: position,
			Cash:            cash,
		}

		signal, err := e.strategy.Generate(strategyCtx)
		if err != nil {
			return nil, fmt.Errorf("strategy error on %s: %w", bar.Timestamp, err)
		}

		switch signal {
		case SignalBuy:
			if position == nil || !position.IsOpen() {
				if cash >= bar.Close {
					shares := cash / bar.Close

					if shares > 0 {
						trade := domain.Trade{
							ID:            uuid.New(),
							BacktestID:    uuid.Nil,
							Symbol:        symbol,
							Direction:     domain.TradeDirectionBuy,
							Quantity:      shares,
							Price:         bar.Close,
							Commission:    0,
							Timestamp:     bar.Timestamp,
							PnL:           0,
							CumulativePnL: 0,
						}
						trades = append(trades, trade)

						position = &Position{
							Symbol:     symbol,
							Shares:     shares,
							EntryPrice: bar.Close,
							EntryTime:  bar.Timestamp,
						}

						cash -= shares * bar.Close
					}
				}
			}

		case SignalSell:
			if position != nil && position.IsOpen() {
				sellValue := position.Shares * bar.Close
				buyValue := position.CostBasis()
				pnl := sellValue - buyValue

				trade := domain.Trade{
					ID:            uuid.New(),
					BacktestID:    uuid.Nil,
					Symbol:        symbol,
					Direction:     domain.TradeDirectionSell,
					Quantity:      position.Shares,
					Price:         bar.Close,
					Commission:    0,
					Timestamp:     bar.Timestamp,
					PnL:           pnl,
					CumulativePnL: 0,
				}
				trades = append(trades, trade)

				cash += position.Shares * bar.Close
				position = nil
			}

		case SignalHold:
			// do nothing
		}
	}

	return trades, nil
}
