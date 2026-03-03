package service

import (
	"context"

	"github.com/boldlogic/PortfolioLens/pkg/models/quik"
)

func (s *Service) SaveFirm(ctx context.Context, code string, name string) (quik.Firm, error) {

	res, err := s.limitsRepo.InsertFirm(ctx, code, name)
	if err != nil {
		return quik.Firm{}, err
	}
	return res, nil
}
