package strategy

import "github.com/wreckitral/distributed-backtesting-platform/internal/domain"

type Context struct {
	Symbol string

	CurrentBar domain.Bar

	HistoricalBars []domain.Bar

	CurrentPosition *Position

	Cash float64
}

func (c *Context) BarCount() int {
	return len(c.HistoricalBars) + 1
}

func (c *Context) HasPosition() bool {
	return c.CurrentPosition != nil && c.CurrentPosition.Shares > 0
}

func (c *Context) AllBars() []domain.Bar {
	all := make([]domain.Bar, len(c.HistoricalBars)+1)
	copy(all, c.HistoricalBars)
	all[len(all)-1] = c.CurrentBar
	return all
}
