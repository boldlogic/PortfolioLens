package service

import (
	"github.com/boldlogic/PortfolioLens/cbr-market-data-worker/internal/client"
	"github.com/boldlogic/PortfolioLens/cbr-market-data-worker/internal/service/request_catalog"
	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/sirupsen/logrus"
)

type Service struct {
	client   *client.Client
	Provider *request_catalog.Provider
	//Storage      *storage.Storage
	log           logrus.FieldLogger
	CurrencyRepo  CurrencyRepository
	fxRateRepo    FxRateRepository
	schedulerRepo JobRepository
}

func NewService(cl *client.Client, registry *request_catalog.Provider, currencyRepo CurrencyRepository, fxRateRepo FxRateRepository, schedulerRepo JobRepository, log logrus.FieldLogger) *Service {

	return &Service{
		client:   cl,
		Provider: registry,
		//Storage:      storage,
		CurrencyRepo:  currencyRepo,
		fxRateRepo:    fxRateRepo,
		schedulerRepo: schedulerRepo,
		log:           log,
	}
}

type CurrencyRepository interface {
	SaveCurrencies([]models.Currency) []error
	GetCurrencies() ([]models.Currency, error)
	GetCurrency(charCode string) (models.Currency, error)
}

type FxRateRepository interface {
	SaveFxRates([]models.FxRate) []error
}

type JobRepository interface {
	GetActionId(code string) (models.Action, error)
	GetAction(id int) (models.Action, error)
	SaveAction(row *models.Action) error

	GetTask(uuid string) (models.Task, error)

	CreateTask(row *models.Task) (models.Task, error)
	FetchTask(status int, newStatus int) (models.Task, error)
	SetTaskStatusCompleted(id int) error
	SetTaskStatusError(id int, err string) error
}
