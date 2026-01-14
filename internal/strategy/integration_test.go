package strategy

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/wreckitral/distributed-backtesting-platform/internal/marketdata"
	"github.com/wreckitral/distributed-backtesting-platform/internal/metrics"
)

// TestFullBacktestFlow tests the complete backtest pipeline
func TestFullBacktestFlow(t *testing.T) {
	// 1. Setup market data
	provider, err := marketdata.NewCSVProvider("../../data/sample")
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// 2. Create strategy (SMA Crossover)
	strategy := NewSMACrossover(10, 30)

	// 3. Create executor
	initialCash := 10000.0
	executor := NewExecutor(strategy, provider, initialCash)

	// 4. Run backtest on AAPL for 2024
	ctx := context.Background()
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Logf("Running backtest: %s on AAPL from %s to %s",
		strategy.Name(),
		start.Format("2006-01-02"),
		end.Format("2006-01-02"))

	trades, err := executor.Run(ctx, "AAPL", start, end)
	if err != nil {
		t.Fatalf("Executor failed: %v", err)
	}

	t.Logf("Backtest complete: %d trades executed", len(trades))

	// 5. Calculate metrics
	calculator := metrics.NewCalculator(initialCash)
	results, err := calculator.Calculate(trades, start, end)
	if err != nil {
		t.Fatalf("Metrics calculation failed: %v", err)
	}

	// 6. Print comprehensive results
	t.Logf("\n%s", strings.Repeat("=", 60))
	t.Logf("BACKTEST RESULTS")
	t.Logf("%s", strings.Repeat("=", 60))

	t.Logf("\nCAPITAL")
	t.Logf("Initial Capital: $%.2f", results.InitialCapital)
	t.Logf("Final Capital: $%.2f", results.FinalCapital)
	t.Logf("Total Return: $%.2f (%.2f%%)", results.TotalReturn, results.ReturnPct)

	t.Logf("\nTRADES")
	t.Logf("Total Trades: %d", results.TotalTrades)
	t.Logf("Winning Trades: %d", results.WinningTrades)
	t.Logf("Losing Trades: %d", results.LosingTrades)
	t.Logf("Win Rate: %.2f%%", results.WinRate)

	t.Logf("\nPROFIT/LOSS")
	t.Logf("Gross Profit: $%.2f", results.GrossProfit)
	t.Logf("Gross Loss: $%.2f", results.GrossLoss)
	t.Logf("Net Profit: $%.2f", results.NetProfit)
	t.Logf("Average Trade: $%.2f", results.AverageTrade)
	t.Logf("Average Win: $%.2f", results.AverageWin)
	t.Logf("Average Loss: $%.2f", results.AverageLoss)
	t.Logf("Profit Factor: %.2f", results.ProfitFactor())
	t.Logf("Expectancy: $%.2f", results.Expectancy())

	t.Logf("\nRISK METRICS")
	t.Logf("Max Drawdown: %.2f%% ($%.2f)", results.MaxDrawdown, results.MaxDrawdownAmt)
	t.Logf("Sharpe Ratio: %.2f", results.SharpeRatio)

	t.Logf("\nDURATION")
	t.Logf("Start Date: %s", results.StartDate.Format("2006-01-02"))
	t.Logf("End Date: %s", results.EndDate.Format("2006-01-02"))
	t.Logf("Days: %d", results.Duration)

	t.Logf("\n%s", strings.Repeat("=", 60))

	// 7. Print sample trades
	if len(trades) > 0 {
		t.Logf("\nSAMPLE TRADES (first 5):")
		for i, trade := range trades {
			if i >= 5 {
				break
			}
			t.Logf("%d. %s %.2f shares @ $%.2f on %s (P&L: $%.2f)",
				i+1,
				trade.Direction,
				trade.Quantity,
				trade.Price,
				trade.Timestamp.Format("2006-01-02"),
				trade.PnL)
		}
		if len(trades) > 5 {
			t.Logf("... and %d more trades", len(trades)-5)
		}
	}

	// 8. Validate reasonable results
	if results.TotalTrades == 0 {
		t.Error("Expected some trades to be executed")
	}

	if results.FinalCapital <= 0 {
		t.Error("Final capital should be positive")
	}

	// Strategy evaluation message
	t.Logf("\nSTRATEGY EVALUATION:")
	if results.ReturnPct > 0 {
		t.Logf("Strategy was PROFITABLE (+%.2f%%)", results.ReturnPct)
	} else {
		t.Logf("Strategy LOST MONEY (%.2f%%)", results.ReturnPct)
	}

	if results.SharpeRatio > 1.0 {
		t.Logf("Good risk-adjusted returns (Sharpe: %.2f)", results.SharpeRatio)
	} else if results.SharpeRatio > 0 {
		t.Logf("Moderate risk-adjusted returns (Sharpe: %.2f)", results.SharpeRatio)
	} else {
		t.Logf("Poor risk-adjusted returns (Sharpe: %.2f)", results.SharpeRatio)
	}

	if results.WinRate > 50 {
		t.Logf("Win rate above 50%% (%.2f%%)", results.WinRate)
	} else {
		t.Logf("Win rate below 50%% (%.2f%%)", results.WinRate)
	}

	if results.MaxDrawdown < 20 {
		t.Logf("Low drawdown risk (<20%%: %.2f%%)", results.MaxDrawdown)
	} else {
		t.Logf("High drawdown risk (%.2f%%)", results.MaxDrawdown)
	}
}

// TestCompareStrategies compares Buy and Hold vs SMA Crossover
func TestCompareStrategies(t *testing.T) {
	provider, err := marketdata.NewCSVProvider("../../data/sample")
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	ctx := context.Background()
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	initialCash := 10000.0

	strategies := []Strategy{
		NewBuyHold(),
		NewSMACrossover(10, 30),
		NewSMACrossover(20, 50),
	}

	t.Logf("\n%s", strings.Repeat("=", 80))
	t.Logf("STRATEGY COMPARISON ON AAPL (2024)")
	t.Logf("%s", strings.Repeat("=", 80))

	for _, strategy := range strategies {
		executor := NewExecutor(strategy, provider, initialCash)
		trades, err := executor.Run(ctx, "AAPL", start, end)
		if err != nil {
			t.Logf("%s failed: %v", strategy.Name(), err)
			continue
		}

		calculator := metrics.NewCalculator(initialCash)
		results, err := calculator.Calculate(trades, start, end)
		if err != nil {
			t.Logf("Metrics for %s failed: %v", strategy.Name(), err)
			continue
		}

		t.Logf("\n%s", strategy.Name())
		t.Logf("Return: %+.2f%% ($%.2f)", results.ReturnPct, results.TotalReturn)
		t.Logf("Trades: %d", results.TotalTrades)
		t.Logf("Win Rate: %.2f%%", results.WinRate)
		t.Logf("Sharpe Ratio: %.2f", results.SharpeRatio)
		t.Logf("Max Drawdown: %.2f%%", results.MaxDrawdown)
		t.Logf("Profit Factor: %.2f", results.ProfitFactor())
	}

	t.Logf("\n%s", strings.Repeat("=", 80))
}
