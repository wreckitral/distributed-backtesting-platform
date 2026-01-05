package marketdata

import (
	"context"
	"time"

	"github.com/wreckitral/distributed-backtesting-platform/internal/domain"
)

type Provider interface {
	GetBars(ctx context.Context, symbol string, start, end time.Time) ([]domain.Bar, error)
	GetLatestBar(ctx context.Context, symbol string) (domain.Bar, error)
	ListSymbols(ctx context.Context) ([]string, error)
}
