package service

import (
	"context"
	"strings"

	"github.com/boldlogic/PortfolioLens/pkg/currencies"
	"github.com/boldlogic/PortfolioLens/pkg/models"
)

func (s *Service) GetCurrencies(ctx context.Context) ([]models.Currency, error) {
	return s.currencyRepo.SelectCurrencies(ctx)
}

func (s *Service) GetCurrency(ctx context.Context, charCode string) (models.Currency, string, error) {
	code := strings.ToUpper(strings.TrimSpace(charCode))

	err := currencies.CheckCurrencyCode(code)
	if err != nil {
		return models.Currency{}, err.Error(), models.ErrBusinessValidation
	}
	cur, err := s.currencyRepo.SelectCurrency(ctx, code)
	if err != nil {
		return models.Currency{}, "", err
	}

	return cur, "", nil
}
