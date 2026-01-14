package metrics

import "time"

type Metrics struct {
	// basic returns
	InitialCapital float64 // starting cash
	FinalCapital   float64 // ending cash
	TotalReturn    float64 // Total return in dollars
	ReturnPct      float64 // Total return in percentage

	// trade statistic
	TotalTrades   int     // number of trades executed
	WinningTrades int     // number of profitable trades
	LosingTrades  int     // number of losing trades
	WinRate       float64 // total return as percentage

	// profit/Loss
	GrossProfit  float64 // total profit from winning trades
	GrossLoss    float64 // total loss from losing trades
	NetProfit    float64 // gross profit, gross loss
	AverageTrade float64 // average P&L per trade
	AverageWin   float64 // average profit on winning trades
	AverageLoss  float64 // average loss on losing trades

	// risk metrics
	MaxDrawdown    float64 // largest peak-to-trough decline (%)
	MaxDrawdownAmt float64 // largest decline in dollars
	SharpeRatio    float64 // risk-adjusted return

	// time
	StartDate time.Time // backtest start date
	EndDate   time.Time // backtest end date
	Duration  int       // number of days
}

// ratio of gross profit and gross loss
func (m *Metrics) ProfitFactor() float64 {
	if m.GrossLoss == 0 {
		return 0
	}

	return m.GrossProfit / m.GrossLoss
}

// expected profit per trade
func (m *Metrics) Expectancy() float64 {
	if m.TotalTrades == 0 {
		return 0
	}

	winProb := float64(m.WinningTrades) / float64(m.TotalTrades)
	lossProb := float64(m.LosingTrades) / float64(m.TotalTrades)
	return (winProb * m.AverageWin) - (lossProb * m.AverageLoss)
}
