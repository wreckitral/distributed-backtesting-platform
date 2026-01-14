package metrics

import (
	"testing"
	"time"

	"github.com/wreckitral/distributed-backtesting-platform/internal/domain"
	"github.com/google/uuid"
)

// TestCalculateEmptyTrades tests metrics with no trades
func TestCalculateEmptyTrades(t *testing.T) {
	calculator := NewCalculator(10000.0)

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	metrics, err := calculator.Calculate([]domain.Trade{}, start, end)
	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}

	// Verify no trading activity
	if metrics.TotalTrades != 0 {
		t.Errorf("Expected 0 trades, got %d", metrics.TotalTrades)
	}

	if metrics.TotalReturn != 0 {
		t.Errorf("Expected 0 return, got %.2f", metrics.TotalReturn)
	}

	if metrics.FinalCapital != 10000.0 {
		t.Errorf("Expected final capital 10000, got %.2f", metrics.FinalCapital)
	}

	t.Logf("Empty trades: Initial=%.2f, Final=%.2f",
		metrics.InitialCapital, metrics.FinalCapital)
}

// TestCalculateProfitableTrade tests a simple profitable round trip
func TestCalculateProfitableTrade(t *testing.T) {
	calculator := NewCalculator(10000.0)

	// Buy 100 shares at $100
	buyTrade := domain.Trade{
		ID:            uuid.New(),
		BacktestID:    uuid.New(),
		Symbol:        "AAPL",
		Direction:     domain.TradeDirectionBuy,
		Quantity:      100,
		Price:         100.0,
		Commission:    0,
		Timestamp:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PnL:           0,
		CumulativePnL: 0,
	}

	// Sell 100 shares at $110 (made $1000 profit)
	sellTrade := domain.Trade{
		ID:            uuid.New(),
		BacktestID:    buyTrade.BacktestID,
		Symbol:        "AAPL",
		Direction:     domain.TradeDirectionSell,
		Quantity:      100,
		Price:         110.0,
		Commission:    0,
		Timestamp:     time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		PnL:           1000.0, // $10 profit per share * 100 shares
		CumulativePnL: 1000.0,
	}

	trades := []domain.Trade{buyTrade, sellTrade}

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	metrics, err := calculator.Calculate(trades, start, end)
	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}

	// Verify basic metrics
	if metrics.TotalTrades != 2 {
		t.Errorf("Expected 2 trades, got %d", metrics.TotalTrades)
	}

	if metrics.WinningTrades != 1 {
		t.Errorf("Expected 1 winning trade, got %d", metrics.WinningTrades)
	}

	if metrics.WinRate != 100.0 {
		t.Errorf("Expected 100%% win rate, got %.2f%%", metrics.WinRate)
	}

	if metrics.NetProfit != 1000.0 {
		t.Errorf("Expected net profit 1000, got %.2f", metrics.NetProfit)
	}

	if metrics.FinalCapital != 11000.0 {
		t.Errorf("Expected final capital 11000, got %.2f", metrics.FinalCapital)
	}

	if metrics.ReturnPct != 10.0 {
		t.Errorf("Expected 10%% return, got %.2f%%", metrics.ReturnPct)
	}

	t.Logf("Profitable trade:")
	t.Logf("Initial Capital: $%.2f", metrics.InitialCapital)
	t.Logf("Final Capital: $%.2f", metrics.FinalCapital)
	t.Logf("Total Return: $%.2f (%.2f%%)", metrics.TotalReturn, metrics.ReturnPct)
	t.Logf("Win Rate: %.2f%%", metrics.WinRate)
	t.Logf("Net Profit: $%.2f", metrics.NetProfit)
}

// TestCalculateMixedTrades tests multiple trades with wins and losses
func TestCalculateMixedTrades(t *testing.T) {
	calculator := NewCalculator(10000.0)

	trades := []domain.Trade{
		// Trade 1: Win +$500
		{ID: uuid.New(), Direction: domain.TradeDirectionBuy, Quantity: 50, Price: 100, PnL: 0, Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{ID: uuid.New(), Direction: domain.TradeDirectionSell, Quantity: 50, Price: 110, PnL: 500, Timestamp: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)},

		// Trade 2: Loss -$300
		{ID: uuid.New(), Direction: domain.TradeDirectionBuy, Quantity: 30, Price: 100, PnL: 0, Timestamp: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
		{ID: uuid.New(), Direction: domain.TradeDirectionSell, Quantity: 30, Price: 90, PnL: -300, Timestamp: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC)},

		// Trade 3: Win +$800
		{ID: uuid.New(), Direction: domain.TradeDirectionBuy, Quantity: 40, Price: 100, PnL: 0, Timestamp: time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC)},
		{ID: uuid.New(), Direction: domain.TradeDirectionSell, Quantity: 40, Price: 120, PnL: 800, Timestamp: time.Date(2024, 1, 30, 0, 0, 0, 0, time.UTC)},
	}

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	metrics, err := calculator.Calculate(trades, start, end)
	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}

	// Verify trade counts
	if metrics.TotalTrades != 6 {
		t.Errorf("Expected 6 trades, got %d", metrics.TotalTrades)
	}

	if metrics.WinningTrades != 2 {
		t.Errorf("Expected 2 winning trades, got %d", metrics.WinningTrades)
	}

	if metrics.LosingTrades != 1 {
		t.Errorf("Expected 1 losing trade, got %d", metrics.LosingTrades)
	}

	// Verify P&L
	expectedNetProfit := 500.0 - 300.0 + 800.0 // = 1000
	if metrics.NetProfit != expectedNetProfit {
		t.Errorf("Expected net profit %.2f, got %.2f", expectedNetProfit, metrics.NetProfit)
	}

	// Verify win rate (2 wins out of 3 closed positions = 66.67%)
	if metrics.WinRate < 66.0 || metrics.WinRate > 67.0 {
		t.Errorf("Expected win rate ~66.67%%, got %.2f%%", metrics.WinRate)
	}

	// Verify profit factor (gross profit / gross loss)
	profitFactor := metrics.ProfitFactor()
	if profitFactor < 4.0 || profitFactor > 4.5 {
		t.Errorf("Expected profit factor ~4.33, got %.2f", profitFactor)
	}

	t.Logf("Mixed trades:")
	t.Logf("   Trades: %d (Wins: %d, Losses: %d)",
		metrics.TotalTrades, metrics.WinningTrades, metrics.LosingTrades)
	t.Logf("   Win Rate: %.2f%%", metrics.WinRate)
	t.Logf("   Gross Profit: $%.2f", metrics.GrossProfit)
	t.Logf("   Gross Loss: $%.2f", metrics.GrossLoss)
	t.Logf("   Net Profit: $%.2f", metrics.NetProfit)
	t.Logf("   Profit Factor: %.2f", profitFactor)
	t.Logf("   Average Trade: $%.2f", metrics.AverageTrade)
}

// TestCalculateDrawdown tests maximum drawdown calculation
func TestCalculateDrawdown(t *testing.T) {
	calculator := NewCalculator(10000.0)

	trades := []domain.Trade{
		// Start: $10,000
		{ID: uuid.New(), Direction: domain.TradeDirectionBuy, Quantity: 100, Price: 100, PnL: 0},
		{ID: uuid.New(), Direction: domain.TradeDirectionSell, Quantity: 100, Price: 120, PnL: 2000}, // $12,000 (peak)

		{ID: uuid.New(), Direction: domain.TradeDirectionBuy, Quantity: 100, Price: 100, PnL: 0},
		{ID: uuid.New(), Direction: domain.TradeDirectionSell, Quantity: 100, Price: 90, PnL: -1000}, // $11,000

		{ID: uuid.New(), Direction: domain.TradeDirectionBuy, Quantity: 100, Price: 100, PnL: 0},
		{ID: uuid.New(), Direction: domain.TradeDirectionSell, Quantity: 100, Price: 70, PnL: -3000}, // $8,000 (drawdown!)

		{ID: uuid.New(), Direction: domain.TradeDirectionBuy, Quantity: 100, Price: 100, PnL: 0},
		{ID: uuid.New(), Direction: domain.TradeDirectionSell, Quantity: 100, Price: 130, PnL: 3000}, // $11,000 (recovery)
	}

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	metrics, err := calculator.Calculate(trades, start, end)
	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}

	if metrics.MaxDrawdownAmt != 4000.0 {
		t.Errorf("Expected max drawdown amount $4000, got $%.2f", metrics.MaxDrawdownAmt)
	}

	if metrics.MaxDrawdown < 33.0 || metrics.MaxDrawdown > 34.0 {
		t.Errorf("Expected max drawdown ~33.33%%, got %.2f%%", metrics.MaxDrawdown)
	}

	t.Logf("Drawdown test:")
	t.Logf("Max Drawdown: %.2f%% ($%.2f)", metrics.MaxDrawdown, metrics.MaxDrawdownAmt)
	t.Logf("Final Capital: $%.2f", metrics.FinalCapital)
}
