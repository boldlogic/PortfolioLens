package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/boldlogic/PortfolioLens/cbr-market-data-worker/internal/service/cbr"
)

func (c *Service) GetCbrCurrencies(ctx context.Context, bdy []byte) error {
	// ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	// defer cancel()
	currencies, err := cbr.ParseCurrenciesXML(bdy)
	if err != nil {
		return err
	}
	errs := c.CurrencyRepo.SaveCurrencies(currencies)
	if len(errs) > 0 {
		if err := errors.Join(errs...); err != nil {
			return fmt.Errorf("%w", err)
		}
	}
	return nil
}
