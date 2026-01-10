package marketdata

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/wreckitral/distributed-backtesting-platform/internal/domain"
)

type CSVProvider struct {
	dataDir string
	cache   map[string][]domain.Bar
	mu      sync.RWMutex
}

func NewCSVProvider(dataDir string) (*CSVProvider, error) {
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("data directory does not exist: %s", dataDir)
	}

	return &CSVProvider{
		dataDir: dataDir,
		cache:   make(map[string][]domain.Bar),
	}, nil
}

func (p *CSVProvider) GetBars(ctx context.Context, symbol string, start, end time.Time) ([]domain.Bar, error) {
	// check cache first
	p.mu.RLock()
	bars, exists := p.cache[symbol]
	p.mu.RUnlock()

	if !exists {
		loadedBars, err := p.loadCSV(symbol)
		if err != nil {
			return nil, err
		}

		p.mu.Lock()
		p.cache[symbol] = loadedBars
		p.mu.Unlock()

		bars = loadedBars
	}

	filtered := []domain.Bar{}
	for _, bar := range bars {
		if (bar.Timestamp.Equal(start) || bar.Timestamp.After(start)) && bar.Timestamp.Before(end) {
			filtered = append(filtered, bar)
		}
	}

	return filtered, nil
}

func (p *CSVProvider) GetLatestBar(ctx context.Context, symbol string) (domain.Bar, error) {
	allBars, err := p.GetBars(ctx, symbol, time.Time{}, time.Now().AddDate(100, 0, 0))
	if err != nil {
		return domain.Bar{}, err
	}

	if len(allBars) == 0 {
		return domain.Bar{}, fmt.Errorf("no bars found for symbol %s", symbol)
	}

	return allBars[len(allBars)-1], nil
}

func (p *CSVProvider) ListSymbols(ctx context.Context) ([]string, error) {
	var filenames []string
	err := filepath.WalkDir(p.dataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), "_daily.csv") {
			filenames = append(filenames, d.Name())
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list symbols: %w", err)
	}

	symbols := make([]string, 0, len(filenames))
	for _, name := range filenames {
		symbol := strings.TrimSuffix(name, "_daily.csv")
		symbols = append(symbols, symbol)
	}

	return symbols, nil
}

func (p *CSVProvider) loadCSV(symbol string) ([]domain.Bar, error) {
	filename := filepath.Join(p.dataDir, symbol+"_daily.csv")

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV %s: %w", symbol, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	bars := []domain.Bar{}
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV: %w", err)
		}

		bar, err := p.parseRow(symbol, row)
		if err != nil {
			return nil, fmt.Errorf("error parsing row: %w", err)
		}

		if err := bar.Validate(); err != nil {
			return nil, fmt.Errorf("invalid bar data: %w", err)
		}

		bars = append(bars, bar)
	}

	return bars, nil
}

func (p *CSVProvider) parseRow(symbol string, row []string) (domain.Bar, error) {
	// CSV format: Date,Open,High,Low,Close,Volume
	// row[0] = Date
	// row[1] = Open
	// row[2] = High
	// row[3] = Low
	// row[4] = Close
	// row[5] = Volume

	if len(row) != 6 {
		return domain.Bar{}, fmt.Errorf("expected 6 columns, got %d", len(row))
	}

	timestamp, err := time.Parse("2006-01-02 15:04:05-07:00", row[0])
	if err != nil {
		return domain.Bar{}, fmt.Errorf("invalid date format: %w", err)
	}

	open, err := strconv.ParseFloat(row[1], 64)
	if err != nil {
		return domain.Bar{}, fmt.Errorf("invalid open price: %w", err)
	}

	high, err := strconv.ParseFloat(row[2], 64)
	if err != nil {
		return domain.Bar{}, fmt.Errorf("invalid high price: %w", err)
	}

	low, err := strconv.ParseFloat(row[3], 64)
	if err != nil {
		return domain.Bar{}, fmt.Errorf("invalid low price: %w", err)
	}

	close, err := strconv.ParseFloat(row[4], 64)
	if err != nil {
		return domain.Bar{}, fmt.Errorf("invalid close price: %w", err)
	}

	volume, err := strconv.ParseInt(row[5], 10, 64)
	if err != nil {
		return domain.Bar{}, fmt.Errorf("invalid volume: %w", err)
	}

	return domain.Bar{
		Symbol:    symbol,
		Timestamp: timestamp,
		Open:      open,
		High:      high,
		Low:       low,
		Close:     close,
		Volume:    volume,
	}, nil
}
