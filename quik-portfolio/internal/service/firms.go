package service

import (
	"context"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
)

func (s *Service) SaveFirm(ctx context.Context, code string, name string) (models.Firm, error) {

	res, err := s.limitsRepo.InsertFirm(ctx, code, name)
	if err != nil {
		return models.Firm{}, err
	}
	return res, nil
}
