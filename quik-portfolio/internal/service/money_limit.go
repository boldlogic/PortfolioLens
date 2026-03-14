package service

import (
	"context"
	"time"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
)

func (s *Service) GetML(ctx context.Context, date time.Time) ([]models.MoneyLimit, error) {
	return s.repo.GetMoneyLimits(ctx, date)
}
