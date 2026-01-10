package strategy

import (
	"context"
	"testing"
	"time"

	"github.com/wreckitral/distributed-backtesting-platform/internal/domain"
	"github.com/wreckitral/distributed-backtesting-platform/internal/marketdata"
)

func TestExecutorBuyAndHold(t *testing.T) {
	provider, err := marketdata.NewCSVProvider("../../data/sample")
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	strategy := NewBuyHold()

	executor := NewExecutor(strategy, provider, 10000.0)

	ctx := context.Background()
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	trades, err := executor.Run(ctx, "AAPL", start, end)
	if err != nil {
		t.Fatalf("Executor failed: %v", err)
	}

	if len(trades) != 1 {
		t.Errorf("Expected 1 trade (buy), got %d", len(trades))
	}

	if len(trades) > 0 {
		firstTrade := trades[0]
		if firstTrade.Direction != domain.TradeDirectionBuy {
			t.Errorf("Expected BUY trade, got %s", firstTrade.Direction)
		}
		if firstTrade.Symbol != "AAPL" {
			t.Errorf("Expected AAPL, got %s", firstTrade.Symbol)
		}
		if firstTrade.Quantity <= 0 {
			t.Errorf("Expected positive quantity, got %f", firstTrade.Quantity)
		}

		t.Logf("Buy and Hold executed successfully!")
		t.Logf("Trade: %s %.2f shares @ $%.2f on %s",
			firstTrade.Direction,
			firstTrade.Quantity,
			firstTrade.Price,
			firstTrade.Timestamp.Format("2006-01-02"))
		t.Logf("Total cost: $%.2f", firstTrade.TotalCost())
	}
}

func TestExecutorSMACrossover(t *testing.T) {
	provider, err := marketdata.NewCSVProvider("../../data/sample")
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	strategy := NewSMACrossover(10, 30)

	executor := NewExecutor(strategy, provider, 10000.0)

	ctx := context.Background()
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	trades, err := executor.Run(ctx, "AAPL", start, end)
	if err != nil {
		t.Fatalf("Executor failed: %v", err)
	}

	t.Logf("SMA Crossover executed successfully!")
	t.Logf("Total trades: %d", len(trades))

	for i, trade := range trades {
		t.Logf("Trade %d: %s %.2f shares @ $%.2f on %s (P&L: $%.2f)",
			i+1,
			trade.Direction,
			trade.Quantity,
			trade.Price,
			trade.Timestamp.Format("2006-01-02"),
			trade.PnL)
	}

	totalPnL := 0.0
	for _, trade := range trades {
		totalPnL += trade.PnL
	}
	t.Logf("Total P&L: $%.2f", totalPnL)
}
