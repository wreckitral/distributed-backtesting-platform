package domain

import "time"

type Bar struct {
	Symbol    string
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    int64
}

func (b Bar) Validate() error {
	if b.High < b.Low {
		return ErrInvalidMarketData{Reason: "high must be higher than low"}
	}
	if b.High < b.Open || b.High < b.Close {
		return ErrInvalidMarketData{Reason: "high must be higher than open/close"}
	}
	if b.Open < 0 || b.High < 0 || b.Low < 0 || b.Close < 0 {
		return ErrInvalidMarketData{Reason: "prices cannot be negative"}
	}
	if b.Volume < 0 {
		return ErrInvalidMarketData{Reason: "volume cannot be negative"}
	}
	if b.Timestamp.IsZero() {
		return ErrInvalidMarketData{Reason: "timestamp is required"}
	}

	return nil
}

func (b Bar) TypicalPrice() float64 {
	return (b.High + b.Low + b.Close) / 3
}

type ErrInvalidMarketData struct {
	Reason string
}

func (e ErrInvalidMarketData) Error() string {
	return "invalid market data: " + e.Reason
}
