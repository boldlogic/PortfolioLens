package service

import (
	"context"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
)

func (s *Service) GetBoards(ctx context.Context) ([]models.Board, error) {
	return s.quikRefsRepo.GetBoards(ctx)
}

func (s *Service) GetBoardByID(ctx context.Context, id uint8) (models.Board, error) {
	return s.quikRefsRepo.GetBoardByID(ctx, id)
}
