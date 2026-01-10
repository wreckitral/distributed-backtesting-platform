package strategy

// buyHold is the simplest strategy, buy first day and hold forever
type BuyHold struct {}

func NewBuyHold() *BuyHold {
	return &BuyHold{}
}

func (s *BuyHold) Name() string {
	return "Buy and Hold"
}

func (s *BuyHold) Generate(ctx *Context) (Signal, error) {
	// if we dont have position yet, buy
	if !ctx.HasPosition() {
		return SignalBuy, nil
	}

	// if we already have a position, hold
	return SignalHold, nil
}

