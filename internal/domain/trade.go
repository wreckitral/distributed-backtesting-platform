package domain

import (
	"time"

	"github.com/google/uuid"
)

type Trade struct {
	ID            uuid.UUID
	BacktestID    uuid.UUID
	Symbol        string
	Direction     TradeDirection
	Quantity      float64
	Price         float64
	Commission    float64
	Timestamp     time.Time
	PnL           float64
	CumulativePnL float64
}

type TradeDirection int

const (
	TradeDirectionBuy = iota
	TradeDirectionSell
)

func (td TradeDirection) String() string {
	switch td {
	case TradeDirectionBuy:
		return "BUY"
	case TradeDirectionSell:
		return "SELL"
	default:
		return "UNKNOWN"
	}
}

func (t Trade) Value() float64 {
	return t.Quantity * t.Price
}

func (t Trade) TotalCost() float64 {
	return t.Value() + t.Commission
}
