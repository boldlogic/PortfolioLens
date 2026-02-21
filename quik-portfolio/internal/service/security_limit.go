package service

import (
	"context"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
)

func (s *Service) GetSL(ctx context.Context) ([]models.SecurityLimit, error) {
	sl, err := s.limitsRepo.GetSecurityLimits(ctx)
	if err != nil {
		return nil, err
	}
	return sl, nil
}
