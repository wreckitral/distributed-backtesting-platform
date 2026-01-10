package strategy

import (
	"fmt"

	"github.com/wreckitral/distributed-backtesting-platform/internal/domain"
)

// SMA calculates the Simple Moving Average over the last N bars
// returns 0 if not enough bars available
func SMA(bars []domain.Bar, period int) (float64, error) {
	if period <= 0 {
		return 0, fmt.Errorf("period must be positive, got %d", period)
	}

	if len(bars) < period {
		return 0, fmt.Errorf("not enough bars: need %d, have %d", period, len(bars))
	}

	// take the last N bars
	relevantBars := bars[len(bars)-period:]

	sum := 0.0
	for _, bar := range relevantBars {
		sum += bar.Close
	}

	return sum / float64(period), nil
}
