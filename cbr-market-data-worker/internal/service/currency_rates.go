package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/boldlogic/PortfolioLens/cbr-market-data-worker/internal/service/cbr"
	"github.com/boldlogic/PortfolioLens/pkg/models"
)

func (c *Service) GetCurrencyRates(ctx context.Context, bdy []byte) error {

	rates, err := cbr.ParseFxRatesXML(bdy)
	if err != nil {
		return err
	}
	errs := c.fxRateRepo.SaveFxRates(rates)
	if len(errs) > 0 {
		if err := errors.Join(errs...); err != nil {
			return fmt.Errorf("%w", err)
		}
	}
	return nil
}

func (c *Service) GetCurrencyRatesDynamic(ctx context.Context, bdy []byte, ccy models.Currency) error {

	rates, err := cbr.ParseFxRateDynamicXML(bdy, ccy.ISOCode)
	if err != nil {
		return err
	}
	errs := c.fxRateRepo.SaveFxRates(rates)
	if len(errs) > 0 {
		if err := errors.Join(errs...); err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	return nil
}
