package strategy

type Signal int

const (
	SignalHold Signal = iota
	SignalBuy
	SignalSell
)

type Strategy interface {
	Name() string
	Generate(ctx *Context) (Signal, error)
}

func (s Signal) String() string {
	switch s {
	case SignalHold:
		return "HOLD"
	case SignalBuy:
		return "BUY"
	case SignalSell:
		return "SELL"
	default:
		return "UNKNOWN"
	}
}
