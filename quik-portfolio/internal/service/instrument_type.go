package service

import "context"

func (s *Service) ActualizeInstrumentTypes(ctx context.Context) error {
	return s.quikRefsRepo.SyncInstrumentTypesFromQuotes(ctx)
}

func (s *Service) SaveInstrumentTypesFromQuotes(ctx context.Context) error {
	return s.quikRefsRepo.SyncInstrumentTypesFromQuotes(ctx)
}

func (s *Service) ActualizeInstrumentSubTypes(ctx context.Context) error {

	err := s.quikRefsRepo.SyncInstrumentSubTypesFromQuotes(ctx)
	return err

}

func (s *Service) ActualizeBoards(ctx context.Context) error {
	if err := s.quikRefsRepo.SyncBoardsFromQuotes(ctx); err != nil {
		return err
	}
	return s.quikRefsRepo.TagBoardsTradePointId(ctx)
}
