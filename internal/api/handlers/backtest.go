package handlers

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/wreckitral/distributed-backtesting-platform/internal/api/dto"
	"github.com/wreckitral/distributed-backtesting-platform/internal/domain"
	"github.com/wreckitral/distributed-backtesting-platform/internal/marketdata"
	"github.com/wreckitral/distributed-backtesting-platform/internal/metrics"
	"github.com/wreckitral/distributed-backtesting-platform/internal/repository"
	"github.com/wreckitral/distributed-backtesting-platform/internal/strategy"
)

type BacktestHandler struct {
	backtestRepo repository.BacktestRepository
	tradeRepo    repository.TradeRepository
	metricsRepo  repository.MetricsRepository
	provider     marketdata.Provider
	validate     *validator.Validate
}

// NewBacktestHandler creates a new backtest handler
func NewBacktestHandler(
	backtestRepo repository.BacktestRepository,
	tradeRepo repository.TradeRepository,
	metricsRepo repository.MetricsRepository,
	provider marketdata.Provider,
) *BacktestHandler {
	return &BacktestHandler{
		backtestRepo: backtestRepo,
		tradeRepo:    tradeRepo,
		metricsRepo:  metricsRepo,
		provider:     provider,
		validate:     validator.New(),
	}
}

func (h *BacktestHandler) executeBacktest(backtest *domain.Backtest) {
	ctx := context.Background()

	backtest.Status = domain.BacktestStatusRunning
	h.backtestRepo.Update(ctx, backtest)

	var strat strategy.Strategy
	switch backtest.StrategyID {
	case "buy_hold":
		strat = strategy.NewBuyHold()
	case "sma_crossover":
		strat = strategy.NewSMACrossover(10, 30)
	case "sma_crossover_20_50":
		strat = strategy.NewSMACrossover(20, 50)
	default:
		backtest.Status = domain.BacktestStatusFailed
		backtest.ErrorMessage = fmt.Sprintf("Unknown strategy: %s", backtest.StrategyID)
		h.backtestRepo.Update(ctx, backtest)
		return
	}

	// execute the strategy
	executor := strategy.NewExecutor(strat, h.provider, backtest.InitialCapital)
	trades, err := executor.Run(ctx, backtest.Symbol, backtest.StartDate, backtest.EndDate)
	if err != nil {
		backtest.Status = domain.BacktestStatusFailed
		backtest.ErrorMessage = err.Error()
		h.backtestRepo.Update(ctx, backtest)
		return
	}

	// save trades
	for i := range trades {
		trades[i].BacktestID = backtest.ID
		if err := h.tradeRepo.Create(ctx, &trades[i]); err != nil {
			backtest.Status = domain.BacktestStatusFailed
			backtest.ErrorMessage = fmt.Sprintf("Failed to save trade: %v", err)
			h.backtestRepo.Update(ctx, backtest)
			return
		}
	}

	// calculate metrics
	calculator := metrics.NewCalculator(backtest.InitialCapital)
	results, err := calculator.Calculate(trades, backtest.StartDate, backtest.EndDate)
	if err != nil {
		backtest.Status = domain.BacktestStatusFailed
		backtest.ErrorMessage = fmt.Sprintf("Failed to calculate metrics: %v", err)
		h.backtestRepo.Update(ctx, backtest)
		return
	}

	// calculate annualized return: ((1 + return_rate)^(365/days) - 1) * 100
	annualizedReturn := 0.0
	if results.Duration > 0 {
		returnRate := results.ReturnPct / 100.0
		annualizedReturn = (math.Pow(1+returnRate, 365.0/float64(results.Duration)) - 1) * 100
	}

	// save metrics
	metricsEntity := &domain.Metrics{
		ID:               uuid.New(),
		BacktestID:       backtest.ID,
		TotalReturn:      results.TotalReturn,
		AnnualizedReturn: annualizedReturn,
		SharpeRatio:      results.SharpeRatio,
		MaxDrawdown:      results.MaxDrawdown,
		WinRate:          results.WinRate,
		TotalTrades:      results.TotalTrades,
		WinningTrades:    results.WinningTrades,
		LosingTrades:     results.LosingTrades,
		ProfitFactor:     results.ProfitFactor(),
		AvgWin:           results.AverageWin,
		AvgLoss:          results.AverageLoss,
	}

	if err := h.metricsRepo.Create(ctx, metricsEntity); err != nil {
		backtest.Status = domain.BacktestStatusFailed
		backtest.ErrorMessage = fmt.Sprintf("Failed to save metrics: %v", err)
		h.backtestRepo.Update(ctx, backtest)
		return
	}

	// update backtest status
	now := time.Now()
	backtest.Status = domain.BacktestStatusCompleted
	backtest.CompletedAt = &now
	backtest.ErrorMessage = ""
	h.backtestRepo.Update(ctx, backtest)
}

// CreateBacktest godoc
//
//	@Summary		Create a new backtest
//	@Description	Create and execute a backtest with the given parameters
//	@Tags			backtests
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateBacktestRequest	true	"Backtest parameters"
//	@Success		201		{object}	dto.BacktestResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/api/v1/backtests [post]
func (h *BacktestHandler) CreateBacktest(c *gin.Context) {
	var req dto.CreateBacktestRequest

	// parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	// validate request
	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
		return
	}

	// parse dates
	startDate, endDate, err := dto.ParseBacktestDates(req.StartDate, req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid date format",
			Message: "Use YYYY-MM-DD format",
		})
		return
	}

	// create backtest domain object
	backtest := &domain.Backtest{
		ID:             uuid.New(),
		StrategyID:     req.StrategyID,
		Symbol:         req.Symbol,
		StartDate:      startDate,
		EndDate:        endDate,
		InitialCapital: req.InitialCapital,
		Status:         domain.BacktestStatusPending,
	}

	// save to database
	ctx := context.Background()
	if err := h.backtestRepo.Create(ctx, backtest); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to create backtest",
			Message: err.Error(),
		})
		return
	}

	// execute backtest asynchronously
	go h.executeBacktest(backtest)

	c.JSON(http.StatusCreated, dto.FromDomainBacktest(backtest))
}

// GetBacktest godoc
//
//	@Summary		Get backtest by ID
//	@Description	Get details of a specific backtest
//	@Tags			backtests
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Backtest ID"
//	@Success		200	{object}	dto.BacktestResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/api/v1/backtests/{id} [get]
func (h *BacktestHandler) GetBacktest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid backtest ID",
			Message: err.Error(),
		})
		return
	}

	ctx := context.Background()
	backtest, err := h.backtestRepo.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "Backtest not found",
		})
		return
	}

	c.JSON(http.StatusOK, dto.FromDomainBacktest(backtest))
}

// ListBacktests godoc
//
//	@Summary		List all backtests
//	@Description	Get a list of all backtests
//	@Tags			backtests
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.ListResponse
//	@Router			/api/v1/backtests [get]
func (h *BacktestHandler) ListBacktests(c *gin.Context) {
	ctx := context.Background()
	backtests, err := h.backtestRepo.List(ctx, 100, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to fetch backtests",
		})
		return
	}

	responses := make([]dto.BacktestResponse, len(backtests))
	for i, b := range backtests {
		responses[i] = dto.FromDomainBacktest(b)
	}

	c.JSON(http.StatusOK, dto.ListResponse{
		Items: responses,
		Total: len(responses),
		Page:  1,
		Limit: 100,
	})
}

// GetBacktestMetrics godoc
//
//	@Summary		Get backtest metrics
//	@Description	Get performance metrics for a specific backtest
//	@Tags			backtests
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Backtest ID"
//	@Success		200	{object}	dto.MetricsResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/api/v1/backtests/{id}/metrics [get]
func (h *BacktestHandler) GetBacktestMetrics(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid backtest ID",
		})
		return
	}

	ctx := context.Background()
	metricsEntity, err := h.metricsRepo.GetByBacktestID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "Metrics not found",
		})
		return
	}

	// ✅ FIXED: Use the converter function instead of manual mapping
	response := dto.FromDomainMetrics(metricsEntity)

	c.JSON(http.StatusOK, response)
}

// GetBacktestTrades godoc
//
//	@Summary		Get backtest trades
//	@Description	Get all trades for a specific backtest
//	@Tags			backtests
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Backtest ID"
//	@Success		200	{object}	dto.ListResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/api/v1/backtests/{id}/trades [get]
func (h *BacktestHandler) GetBacktestTrades(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid backtest ID",
		})
		return
	}

	ctx := context.Background()
	trades, err := h.tradeRepo.GetByBacktestID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to fetch trades",
		})
		return
	}

	responses := make([]dto.TradeResponse, len(trades))
	for i, t := range trades {
		responses[i] = dto.FromDomainTrade(t)
	}

	c.JSON(http.StatusOK, dto.ListResponse{
		Items: responses,
		Total: len(responses),
		Page:  1,
		Limit: 1000,
	})
}

// DeleteBacktest godoc
//
//	@Summary		Delete a backtest
//	@Description	Delete a backtest and all associated data
//	@Tags			backtests
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Backtest ID"
//	@Success		200	{object}	dto.SuccessResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/api/v1/backtests/{id} [delete]
func (h *BacktestHandler) DeleteBacktest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid backtest ID",
			Message: err.Error(),
		})
		return
	}

	ctx := context.Background()

	_, err = h.backtestRepo.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "Backtest not found",
		})
		return
	}

	if err := h.backtestRepo.Delete(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to delete backtest",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Backtest deleted successfully",
	})
}
