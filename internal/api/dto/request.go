package dto

import (
	"fmt"
	"time"
)

type CreateBacktestRequest struct {
	StrategyID     string  `json:"strategy_id" binding:"required" example:"buy_hold"`
	Symbol         string  `json:"symbol" binding:"required" example:"AAPL"`
	StartDate      string  `json:"start_date" binding:"required" example:"2024-01-01"`
	EndDate        string  `json:"end_date" binding:"required" example:"2024-12-31"`
	InitialCapital float64 `json:"initial_capital" binding:"required,gt=0" example:"10000"`
}

func ParseBacktestDates(startDateStr, endDateStr string) (time.Time, time.Time, error) {
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start_date format: %w", err)
	}

	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end_date format: %w", err)
	}

	if endDate.Before(startDate) {
		return time.Time{}, time.Time{}, fmt.Errorf("end_date must be after start_date")
	}

	return startDate, endDate, nil
}
