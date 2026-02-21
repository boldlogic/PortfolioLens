package service

import (
	"context"
	"errors"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
)

func (s *Service) GetLimits(ctx context.Context) ([]models.Limit, error) {
	var res []models.Limit
	ml, err := s.limitsRepo.GetMoneyLimits(ctx)
	if err != nil && !errors.Is(err, models.ErrMLNotFound) {
		return nil, err
	}

	for _, m := range ml {
		res = append(res, models.Limit{
			LoadDate:   m.LoadDate,
			ClientCode: m.ClientCode,
			Ticker:     m.Currency,
			FirmCode:   m.FirmCode,
			FirmName:   m.FirmName,
			Balance:    m.Balance,
		})
	}

	sl, err := s.limitsRepo.GetSecurityLimits(ctx)
	if err != nil && !errors.Is(err, models.ErrSLNotFound) {
		return nil, err
	}

	for _, s := range sl {
		res = append(res, models.Limit{
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
		return nil, models.ErrLimitsNotFound
	}

	return res, nil
}

func (s *Service) GetPortfolio(ctx context.Context) ([]models.PortfolioItem, error) {
	return s.limitsRepo.GetPortfolio(ctx)
}
