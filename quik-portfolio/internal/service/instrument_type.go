package service

import "context"

func (s *Service) ActualizeInstrumentTypes(ctx context.Context) error {
	return s.quikRefsRepo.SyncInstrumentTypesFromQuotes(ctx)
}

func (s *Service) ActualizeInstrumentSubTypes(ctx context.Context) error {
	return s.quikRefsRepo.SyncInstrumentSubTypesFromQuotes(ctx)

}

func (s *Service) ActualizeBoards(ctx context.Context) error {
	if err := s.quikRefsRepo.SyncBoardsFromQuotes(ctx); err != nil {
		return err
	}
	return s.quikRefsRepo.TagBoardsTradePointId(ctx)
}
