package strategy

import (
	"testing"
	"time"

	"github.com/wreckitral/distributed-backtesting-platform/internal/domain"
)

func TestSMA(t *testing.T) {
	bars := []domain.Bar{
		{Close: 100, Timestamp: time.Now()},
		{Close: 110, Timestamp: time.Now()},
		{Close: 120, Timestamp: time.Now()},
		{Close: 130, Timestamp: time.Now()},
		{Close: 140, Timestamp: time.Now()},
	}

	sma, err := SMA(bars, 3)
	if err != nil {
		t.Fatalf("SMA failed: %v", err)
	}

	expected := 130.0
	if sma != expected {
		t.Errorf("Expected SMA of %.2f, got %.2f", expected, sma)
	}

	sma, err = SMA(bars, 5)
	if err != nil {
		t.Fatalf("SMA failed: %v", err)
	}

	expected = 120.0
	if sma != expected {
		t.Errorf("Expected SMA of %.2f, got %.2f", expected, sma)
	}

	_, err = SMA(bars, 10)
	if err == nil {
		t.Error("Expected error for insufficient bars, got nil")
	}
}
