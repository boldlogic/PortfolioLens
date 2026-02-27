package service

import (
	"context"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
)

func (s *Service) GetTradePoints(ctx context.Context) ([]models.TradePoint, error) {
	res, err := s.quikRefsRepo.GetTradePoints(ctx)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) GetTradePointByID(ctx context.Context, id uint8) (models.TradePoint, error) {
	return s.quikRefsRepo.GetTradePointByID(ctx, id)
}
