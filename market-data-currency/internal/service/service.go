package service

import (
	"context"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"go.uber.org/zap"
)

type Service struct {
	logger       *zap.Logger
	currencyRepo CurrencyRepository
}

func NewService(ctx context.Context,
	currencyRepo CurrencyRepository,
	logger *zap.Logger) *Service {

	return &Service{
		logger:       logger,
		currencyRepo: currencyRepo,
	}
}

type CurrencyRepository interface {
	SelectCurrencies(ctx context.Context) ([]models.Currency, error)
	SelectCurrency(ctx context.Context, charCode string) (models.Currency, error)
	SelectNewCurrenciesFromCurrentQuotes(ctx context.Context) ([]models.Currency, error)
	MergeCurrencies(ctx context.Context, currencies []models.Currency) error
	SelectCountCurrencies(ctx context.Context) (int, error)
	SetEmptyCurrencyNamesFromQuik(ctx context.Context) error
}
