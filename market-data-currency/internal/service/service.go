package service

import (
	"context"
	"net/http"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/models/requestplan"
	"github.com/boldlogic/PortfolioLens/pkg/models/scheduler"
	"go.uber.org/zap"
)

type CBRParser interface {
	ParseCurrenciesXML(bdy []byte) ([]models.Currency, []models.ExternalCode, error)
	ParseFxRatesXML(bdy []byte) ([]models.FxRate, error)
	ParseFxRateDynamicXML(bdy []byte, base int) ([]models.FxRate, error)
}

type Client interface {
	SendRequest(ctx context.Context, req *http.Request) (int, []byte, error)
	SendWithRetry(ctx context.Context, req *http.Request, retryCount int) (int, []byte, int, error)
}

type responseHandlerFunc func(ctx context.Context, body []byte, taskId int64, taskParams map[string]string) error

type Service struct {
	client Client

	logger           *zap.Logger
	cbrParser        CBRParser
	currencyRepo     CurrencyRepository
	schedulerRepo    SchedulerRepository
	responseHandlers map[string]responseHandlerFunc
}

func NewService(
	client Client,
	cbrParser CBRParser,
	currencyRepo CurrencyRepository,
	schedulerRepo SchedulerRepository,
	logger *zap.Logger) *Service {

	s := &Service{
		client:        client,
		logger:        logger,
		cbrParser:     cbrParser,
		currencyRepo:  currencyRepo,
		schedulerRepo: schedulerRepo,
	}
	s.responseHandlers = s.buildResponseHandlers()
	return s
}

func (s *Service) buildResponseHandlers() map[string]responseHandlerFunc {
	return map[string]responseHandlerFunc{
		scheduler.ActionCodeCbrCurrencyList:    s.handleCbrCurrencyList,
		scheduler.ActionCodeCbrRatesToday:     s.handleCbrRatesToday,
		scheduler.ActionCodeCbrHistoricalRates: s.handleCbrHistoricalRates,
	}
}

type CurrencyRepository interface {
	SelectCurrencies(ctx context.Context) ([]models.Currency, error)
	SelectCurrency(ctx context.Context, charCode string) (models.Currency, error)
	MergeCurrencies(ctx context.Context, currencies []models.Currency) error
	MergeExternalCodes(ctx context.Context, codes []models.ExternalCode) error
	SelectCountCurrencies(ctx context.Context) (int, error)
	SetEmptyCurrencyNamesFromQuik(ctx context.Context) error
	MergeFxCBRRates(ctx context.Context, rates []models.FxRate) error
	MergeFxCBRRatesQuik(ctx context.Context) error
}

type SchedulerRepository interface {
	FetchOneNewTask(ctx context.Context) (scheduler.Task, error)
	SelectAction(ctx context.Context, id uint8) (scheduler.Action, error)
	UpdateTaskStatus(ctx context.Context, id int64, newStatus scheduler.TaskStatusID, errMsg string) error
	SelectTaskParams(ctx context.Context, taskId int64) ([]scheduler.TaskParam, error)
	SelectExternalCodeByCurrency(ctx context.Context, isoCharCode string, extCodeTypeId uint8, actionId uint8) (string, error)
	SelectRequestPlan(ctx context.Context, actionId uint8) (requestplan.RequestPlan, error)
}
