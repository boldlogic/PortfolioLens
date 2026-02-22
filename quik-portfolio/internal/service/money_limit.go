package service

import (
	"context"
	"time"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
)

func (s *Service) GetML(ctx context.Context, date time.Time) ([]models.MoneyLimit, error) {

	ml, err := s.limitsRepo.GetMoneyLimits(ctx, date)
	if err != nil {
		return nil, err
	}
	return ml, nil
}
