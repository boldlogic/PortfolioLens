package service

import (
	"context"
	"errors"
	"time"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	qmodels "github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
)

func (s *Service) GetLimits(ctx context.Context, date time.Time) ([]qmodels.Limit, error) {
	var res []qmodels.Limit

	ml, err := s.repo.GetMoneyLimits(ctx, date)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return nil, err
	}
	for _, m := range ml {
		res = append(res, qmodels.Limit{
			LoadDate:   m.LoadDate,
			ClientCode: m.ClientCode,
			Ticker:     m.Currency,
			FirmCode:   m.FirmCode,
			FirmName:   m.FirmName,
			Balance:    m.Balance,
		})
	}

	sl, err := s.repo.GetSecurityLimits(ctx, date)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return nil, err
	}
	for _, s := range sl {
		res = append(res, qmodels.Limit{
			LoadDate:       s.LoadDate,
			ClientCode:     s.ClientCode,
			Ticker:         s.Ticker,
			FirmCode:       s.FirmCode,
			FirmName:       s.FirmName,
			Balance:        s.Balance,
			ISIN:           s.ISIN,
			AcquisitionCcy: s.AcquisitionCcy,
		})
	}

	if len(res) == 0 {
		return nil, models.ErrNotFound
	}
	return res, nil
}

func (s *Service) GetPortfolio(ctx context.Context) ([]qmodels.PortfolioItem, error) {
	return s.repo.GetPortfolio(ctx)
}
