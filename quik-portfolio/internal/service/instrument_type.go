package service

import "context"

func (s *Service) ActualizeInstrumentTypes(ctx context.Context) error {

	err := s.instrTypeRepo.SyncInstrumentTypesFromQuotes(ctx)
	return err

}

func (s *Service) ActualizeInstrumentSubTypes(ctx context.Context) error {

	err := s.instrTypeRepo.SyncInstrumentSubTypesFromQuotes(ctx)
	return err

}

func (s *Service) ActualizeBoards(ctx context.Context) error {
	err := s.instrTypeRepo.SyncBoardsFromQuotes(ctx)
	return err
}
