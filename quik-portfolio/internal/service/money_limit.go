package service

import (
	"context"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
)

func (s *Service) GetML(ctx context.Context) ([]models.MoneyLimit, error) {

	ml, err := s.limitsRepo.GetMoneyLimits(ctx)
	if err != nil {
		return nil, err
	}
	return ml, nil
}
