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
