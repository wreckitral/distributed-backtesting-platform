package common

import "errors"

var (
	ErrNotFound        = errors.New("resource not found")
	ErrInvalidInput    = errors.New("invalid input")
	ErrBacktestRunning = errors.New("backtest already running")
)

const (
	DefaultCommission     = 0.001 // 0.1%
	DefaultInitialCapital = 100000.0
	MaxBacktestDuration   = 10 * 365 // 10 years in days
)

type Timeframe string

const (
	TimeframeDaily  = "1d"
	TimeframeHourly = "1h"
)
