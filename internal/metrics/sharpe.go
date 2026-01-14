package metrics

import (
	"math"

	"github.com/wreckitral/distributed-backtesting-platform/internal/domain"
)

// calculateSharpe computes the Sharpe ratio
// sharpe = (Average Return - Risk Free Rate) / Standard Deviation of Returns
// this is using 0% risk-free rate for simplicity
func (c *Calculator) calculateSharpe(trades []domain.Trade, m *Metrics) {
	if len(trades) < 2 {
		m.SharpeRatio = 0
		return
	}

	// collect returns from each closed trade (sell trades only)
	var returns []float64
	for _, trade := range trades {
		if trade.Direction == domain.TradeDirectionSell {
			// Return as percentage of capital
			returnPct := (trade.PnL / c.initialCapital) * 100
			returns = append(returns, returnPct)
		}
	}

	if len(returns) < 2 {
		m.SharpeRatio = 0
		return
	}

	// calculate average return
	avgReturn := 0.0
	for _, r := range returns {
		avgReturn += r
	}
	avgReturn /= float64(len(returns))

	// calculate standard deviation
	variance := 0.0
	for _, r := range returns {
		variance += math.Pow(r-avgReturn, 2)
	}
	variance /= float64(len(returns))
	stdDev := math.Sqrt(variance)

	// sharpe ratio (annualized)
	// assuming ~252 trading days per year
	if stdDev > 0 {
		m.SharpeRatio = (avgReturn / stdDev) * math.Sqrt(252)
	}
}
