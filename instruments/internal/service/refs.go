package service

import (
	"context"

	md "github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/models/quik"
)

// ActualizeRefs последовательно синхронизирует все справочники.
// Порядок важен: subtypes зависят от types, boards теггируются после sync.
func (s *Service) ActualizeRefs(ctx context.Context) error {
	if err := s.refsSyncRepo.SyncInstrumentTypesFromQuotes(ctx); err != nil {
		return err
	}
	if err := s.refsSyncRepo.SyncInstrumentSubTypesFromQuotes(ctx); err != nil {
		return err
	}
	if err := s.refsSyncRepo.SyncBoardsFromQuotes(ctx); err != nil {
		return err
	}
	return s.refsSyncRepo.TagBoardsTradePointId(ctx)
}

func (s *Service) GetBoards(ctx context.Context) ([]quik.Board, error) {
	return s.refsQueryRepo.GetBoardsWithTradePoint(ctx)
}

func (s *Service) GetBoardByID(ctx context.Context, id uint8) (quik.Board, error) {
	return s.refsQueryRepo.GetBoardByIDWithTradePoint(ctx, id)
}

func (s *Service) GetTradePoints(ctx context.Context) ([]md.TradePoint, error) {
	return s.refsQueryRepo.GetTradePoints(ctx)
}

func (s *Service) GetTradePointByID(ctx context.Context, id uint8) (md.TradePoint, error) {
	return s.refsQueryRepo.GetTradePointByID(ctx, id)
}
