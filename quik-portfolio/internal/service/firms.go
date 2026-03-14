package service

import (
	"context"

	"github.com/boldlogic/PortfolioLens/pkg/models/quik"
)

func (s *Service) SaveFirm(ctx context.Context, code string, name string) (quik.Firm, error) {
	return s.repo.InsertFirm(ctx, code, name)
}
