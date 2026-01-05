package marketdata

import (
	"context"
	"testing"
	"time"
)

func TestNewCSVProvider(t *testing.T) {
	provider, err := NewCSVProvider("../../data/sample")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if provider == nil {
		t.Fatal("Expected provider to be created")
	}

	_, err = NewCSVProvider("/invalid/path/does/not/exist")
	if err == nil {
		t.Fatal("Expected error for invalid directory, got nil")
	}
}

func TestGetBars(t *testing.T) {
	provider, err := NewCSVProvider("../../data/sample")
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	ctx := context.Background()

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	bars, err := provider.GetBars(ctx, "AAPL", start, end)
	if err != nil {
		t.Fatalf("GetBars failed: %v", err)
	}

	if len(bars) == 0 {
		t.Fatal("Expected some bars, got 0")
	}

	t.Logf("Got %d bars for AAPL in Jan 2024", len(bars))

	for _, bar := range bars {
		if bar.Timestamp.Before(start) {
			t.Errorf("Bar timestamp %v is before start %v", bar.Timestamp, start)
		}
		if bar.Timestamp.After(end) || bar.Timestamp.Equal(end) {
			t.Errorf("Bar timestamp %v is not before end %v", bar.Timestamp, end)
		}
	}

	firstBar := bars[0]
	if firstBar.Symbol != "AAPL" {
		t.Errorf("Expected symbol AAPL, got %s", firstBar.Symbol)
	}
	if firstBar.Open <= 0 {
		t.Errorf("Expected positive open price, got %f", firstBar.Open)
	}
	if firstBar.Volume <= 0 {
		t.Errorf("Expected positive volume, got %d", firstBar.Volume)
	}

	t.Logf("First bar: %+v", firstBar)
}

func TestGetBarsCache(t *testing.T) {
	provider, err := NewCSVProvider("../../data/sample")
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	ctx := context.Background()
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	bars1, err := provider.GetBars(ctx, "AAPL", start, end)
	if err != nil {
		t.Fatalf("First GetBars failed: %v", err)
	}

	bars2, err := provider.GetBars(ctx, "AAPL", start, end)
	if err != nil {
		t.Fatalf("Second GetBars failed: %v", err)
	}

	if len(bars1) != len(bars2) {
		t.Errorf("Expected same number of bars, got %d and %d", len(bars1), len(bars2))
	}

	provider.mu.RLock()
	_, exists := provider.cache["AAPL"]
	provider.mu.RUnlock()

	if !exists {
		t.Error("Expected AAPL to be in cache")
	}

	t.Logf("Cache working correctly - both calls returned %d bars", len(bars1))
}

func TestGetLatestBar(t *testing.T) {
	provider, err := NewCSVProvider("../../data/sample")
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	ctx := context.Background()

	bar, err := provider.GetLatestBar(ctx, "AAPL")
	if err != nil {
		t.Fatalf("GetLatestBar failed: %v", err)
	}

	if bar.Symbol != "AAPL" {
		t.Errorf("Expected symbol AAPL, got %s", bar.Symbol)
	}
	if bar.Timestamp.IsZero() {
		t.Error("Expected non-zero timestamp")
	}
	if bar.Close <= 0 {
		t.Errorf("Expected positive close price, got %f", bar.Close)
	}

	t.Logf("Latest bar: Date=%v, Close=%f", bar.Timestamp, bar.Close)
}

func TestListSymbols(t *testing.T) {
	provider, err := NewCSVProvider("../../data/sample")
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	ctx := context.Background()

	symbols, err := provider.ListSymbols(ctx)
	if err != nil {
		t.Fatalf("ListSymbols failed: %v", err)
	}

	if len(symbols) == 0 {
		t.Fatal("Expected some symbols, got 0")
	}

	expectedSymbols := map[string]bool{
		"AAPL": false,
		"MSFT": false,
		"SPY":  false,
	}

	for _, symbol := range symbols {
		if _, exists := expectedSymbols[symbol]; exists {
			expectedSymbols[symbol] = true
		}
	}

	for symbol, found := range expectedSymbols {
		if !found {
			t.Errorf("Expected to find symbol %s", symbol)
		}
	}

	t.Logf("Found symbols: %v", symbols)
}

func TestGetBarsInvalidSymbol(t *testing.T) {
	provider, err := NewCSVProvider("../../data/sample")
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	ctx := context.Background()
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	_, err = provider.GetBars(ctx, "INVALID_SYMBOL", start, end)
	if err == nil {
		t.Fatal("Expected error for invalid symbol, got nil")
	}

	t.Logf("Correctly got error for invalid symbol: %v", err)
}

func TestParseRow(t *testing.T) {
	provider := &CSVProvider{}

	row := []string{
		"2024-01-03 00:00:00-05:00",
		"182.496",
		"184.141",
		"181.714",
		"182.526",
		"58414500",
	}

	bar, err := provider.parseRow("AAPL", row)
	if err != nil {
		t.Fatalf("parseRow failed: %v", err)
	}

	if bar.Symbol != "AAPL" {
		t.Errorf("Expected symbol AAPL, got %s", bar.Symbol)
	}
	if bar.Open != 182.496 {
		t.Errorf("Expected open 182.496, got %f", bar.Open)
	}
	if bar.High != 184.141 {
		t.Errorf("Expected high 184.141, got %f", bar.High)
	}
	if bar.Low != 181.714 {
		t.Errorf("Expected low 181.714, got %f", bar.Low)
	}
	if bar.Close != 182.526 {
		t.Errorf("Expected close 182.526, got %f", bar.Close)
	}
	if bar.Volume != 58414500 {
		t.Errorf("Expected volume 58414500, got %d", bar.Volume)
	}

	invalidRow := []string{"2024-01-03", "182.496"}
	_, err = provider.parseRow("AAPL", invalidRow)
	if err == nil {
		t.Error("Expected error for invalid row, got nil")
	}
}
