package service

import (
	"context"

	md "github.com/boldlogic/PortfolioLens/pkg/models"
)

func (s *Service) GetTradePoints(ctx context.Context) ([]md.TradePoint, error) {
	res, err := s.quikRefsRepo.GetTradePoints(ctx)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) GetTradePointByID(ctx context.Context, id uint8) (md.TradePoint, error) {
	return s.quikRefsRepo.GetTradePointByID(ctx, id)
}
