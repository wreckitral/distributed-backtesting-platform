package metrics

import (
	"math"
	"time"

	"github.com/wreckitral/distributed-backtesting-platform/internal/domain"
)

// computes performance metrics from trades
type Calculator struct {
	initialCapital float64
}

// NewCalculator creates a new metrics calculator
func NewCalculator(initialCapital float64) *Calculator {
	return &Calculator{
		initialCapital: initialCapital,
	}
}

// computes all metrics from a list of trades
func (c *Calculator) Calculate(trades []domain.Trade, startDate, endDate time.Time) (*Metrics, error) {
	if len(trades) == 0 {
		return c.emptyMetrics(startDate, endDate), nil
	}

	m := &Metrics{
		InitialCapital: c.initialCapital,
		StartDate:      startDate,
		EndDate:        endDate,
		Duration:       int(endDate.Sub(startDate).Hours() / 24),
	}

	// basic trade statistics
	c.calculateTradeStats(trades, m)

	// returns
	c.calculateReturns(trades, m)

	// drawdown
	c.calculateDrawdown(trades, m)

	// sharpe ratio
	c.calculateSharpe(trades, m)

	return m, nil
}

// emptyMetrics returns metrics for a backtest with no trades
func (c *Calculator) emptyMetrics(startDate, endDate time.Time) *Metrics {
	return &Metrics{
		InitialCapital: c.initialCapital,
		FinalCapital:   c.initialCapital,
		TotalReturn:    0,
		ReturnPct:      0,
		TotalTrades:    0,
		StartDate:      startDate,
		EndDate:        endDate,
		Duration:       int(endDate.Sub(startDate).Hours() / 24),
	}
}

// calculateTradeStats computes trade-level statistics
func (c *Calculator) calculateTradeStats(trades []domain.Trade, m *Metrics) {
	m.TotalTrades = len(trades)

	for _, trade := range trades {
		// only count P&L from sell trades
		if trade.Direction == domain.TradeDirectionSell {
			if trade.PnL > 0 {
				m.WinningTrades++
				m.GrossProfit += trade.PnL
			} else if trade.PnL < 0 {
				m.LosingTrades++
				m.GrossLoss += math.Abs(trade.PnL)
			}
		}
	}

	// calculate averages
	if m.TotalTrades > 0 {
		totalPnL := 0.0
		for _, t := range trades {
			if t.Direction == domain.TradeDirectionSell {
				totalPnL += t.PnL
			}
		}
		m.NetProfit = totalPnL

		// average per closed trade (buy+sell = 1 round trip)
		roundTrips := len(trades) / 2
		if roundTrips > 0 {
			m.AverageTrade = totalPnL / float64(roundTrips)
		}
	}

	if m.WinningTrades > 0 {
		m.AverageWin = m.GrossProfit / float64(m.WinningTrades)
	}

	if m.LosingTrades > 0 {
		m.AverageLoss = m.GrossLoss / float64(m.LosingTrades)
	}

	if m.WinningTrades+m.LosingTrades > 0 {
		m.WinRate = float64(m.WinningTrades) / float64(m.WinningTrades+m.LosingTrades) * 100
	}
}

// calculateReturns computes return metrics
func (c *Calculator) calculateReturns(trades []domain.Trade, m *Metrics) {
	// final capital = initial + all realized P&L
	m.FinalCapital = c.initialCapital

	for _, trade := range trades {
		if trade.Direction == domain.TradeDirectionSell {
			m.FinalCapital += trade.PnL
		}
	}

	m.TotalReturn = m.FinalCapital - c.initialCapital

	if c.initialCapital > 0 {
		m.ReturnPct = (m.TotalReturn / c.initialCapital) * 100
	}
}

// computes maximum drawdown
func (c *Calculator) calculateDrawdown(trades []domain.Trade, m *Metrics) {
	if len(trades) == 0 {
		return
	}

	// track equity curve (cash balance over time)
	equity := c.initialCapital
	peak := equity
	maxDrawdown := 0.0
	maxDrawdownAmt := 0.0

	for _, trade := range trades {
		// update equity after each trade
		if trade.Direction == domain.TradeDirectionSell {
			equity += trade.PnL
		}

		// update peak if we've reached a new high
		if equity > peak {
			peak = equity
		}

		// calculate current drawdown from peak
		drawdown := (peak - equity) / peak * 100
		drawdownAmt := peak - equity

		// track maximum drawdown
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
			maxDrawdownAmt = drawdownAmt
		}
	}

	m.MaxDrawdown = maxDrawdown
	m.MaxDrawdownAmt = maxDrawdownAmt
}
