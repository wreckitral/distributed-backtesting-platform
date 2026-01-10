package strategy

type SMACrossover struct {
	ShortPeriod int
	LongPeriod  int
}

func NewSMACrossover(shortPeriod, longPeriod int) *SMACrossover {
	return &SMACrossover{
		ShortPeriod: shortPeriod,
		LongPeriod:  longPeriod,
	}
}

func (s *SMACrossover) Name() string {
	return "SMA Crossover"
}

func (s *SMACrossover) Generate(ctx *Context) (Signal, error) {
	allBars := ctx.AllBars()

	// need enough data for the long period
	if len(allBars) < s.LongPeriod {
		return SignalHold, nil
	}

	// get the simple moving average for short period
	shortSMA, err := SMA(allBars, s.ShortPeriod)
	if err != nil {
		return SignalHold, err
	}

	// get the simple moving average for long period
	longSMA, err := SMA(allBars, s.LongPeriod)
	if err != nil {
		return SignalHold, err
	}

	// need previous SMAs to detect crossover
	if len(allBars) < s.LongPeriod+1 {
		return SignalHold, nil
	}

	// calculate previous days SMA except the current one
	previousBars := allBars[:len(allBars)-1]
	prevShortSMA, err := SMA(previousBars, s.ShortPeriod)
	if err != nil {
		return SignalHold, err
	}

	prevLongSMA, err := SMA(previousBars, s.LongPeriod)
	if err != nil {
		return SignalHold, err
	}

	// detect crossover
	// golden cross, short crosses above long (bullish, buy signal)
	if prevShortSMA <= prevLongSMA && shortSMA > longSMA {
		if !ctx.HasPosition() {
			return SignalBuy, nil
		}
	}

	// death cross, short crosses below long (bearish, sell signal)
	if prevShortSMA >= prevLongSMA && shortSMA < longSMA {
		if !ctx.HasPosition() {
			return SignalSell, nil
		}
	}

	return SignalHold, nil
}
