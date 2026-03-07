package service

import (
	"context"

	"github.com/boldlogic/PortfolioLens/market-data-currency/internal/client"
	"github.com/boldlogic/PortfolioLens/market-data-currency/internal/service/request_catalog"
	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/models/scheduler"
	"go.uber.org/zap"
)

type Service struct {
	client *client.Client

	logger        *zap.Logger
	currencyRepo  CurrencyRepository
	schedulerRepo SchedulerRepository
	provider      *request_catalog.Provider
}

func NewService(ctx context.Context,
	client *client.Client,
	currencyRepo CurrencyRepository,
	schedulerRepo SchedulerRepository,
	provider *request_catalog.Provider,
	logger *zap.Logger) *Service {

	return &Service{
		client:        client,
		logger:        logger,
		currencyRepo:  currencyRepo,
		schedulerRepo: schedulerRepo,
		provider:      provider,
	}
}

type CurrencyRepository interface {
	SelectCurrencies(ctx context.Context) ([]models.Currency, error)
	SelectCurrency(ctx context.Context, charCode string) (models.Currency, error)
	SelectNewCurrenciesFromCurrentQuotes(ctx context.Context) ([]models.Currency, error)
	MergeCurrencies(ctx context.Context, currencies []models.Currency) error
	MergeExternalCodes(ctx context.Context, codes []models.ExternalCode) error
	SelectCountCurrencies(ctx context.Context) (int, error)
	SetEmptyCurrencyNamesFromQuik(ctx context.Context) error
}

type SchedulerRepository interface {
	FetchOneNewTask(ctx context.Context) (scheduler.Task, error)
	SelectAction(ctx context.Context, id uint8) (scheduler.Action, error)
	UpdateTaskStatus(ctx context.Context, id int64, newStatus scheduler.TaskStatusID, errMsg string) error
}
