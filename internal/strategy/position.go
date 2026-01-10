package strategy

import "time"

type Position struct {
	Symbol string
	Shares float64
	EntryPrice float64
	EntryTime time.Time
}

func (p *Position) IsOpen() bool {
	return p.Shares > 0
}

func (p *Position) CostBasis() float64 {
	return float64(p.Shares) * p.EntryPrice
}

func (p *Position) Value(currentPrice float64) float64 {
	return float64(p.Shares) * currentPrice
}

func (p *Position) ProfitLoss(currentPrice float64) float64 {
	return p.Value(currentPrice) - p.CostBasis()
}

func (p *Position) ProfitLossPercent(currentPrice float64) float64 {
	if p.CostBasis() == 0 {
		return 0
	}
	return (p.ProfitLoss(currentPrice) / p.CostBasis()) * 100
}
