package service

import (
	"context"

	"github.com/boldlogic/PortfolioLens/pkg/models/quik"
)

func (s *Service) GetBoards(ctx context.Context) ([]quik.Board, error) {
	return s.quikRefsRepo.GetBoardsWithTradePoint(ctx)
}

func (s *Service) GetBoardByID(ctx context.Context, id uint8) (quik.Board, error) {
	return s.quikRefsRepo.GetBoardByIDWithTradePoint(ctx, id)
}
