package service

import (
	"context"

	"github.com/boldlogic/PortfolioLens/instruments/internal/apperrors"
	"github.com/boldlogic/PortfolioLens/instruments/internal/models"
	"go.uber.org/zap"
)

type Service struct {
	logger    *zap.Logger
	instrRepo InstrumentRepository
}

type InstrumentRepository interface {
	SelectInstrumentFromNewCurrentQuote(ctx context.Context) (models.Instrument, models.InstrumentBoard, string, error)
	InsInstrument(ctx context.Context, i models.Instrument) (int, error)
	SetInstrument(ctx context.Context, id int, ic string) error
	GetInstrumentId(ctx context.Context, ticker string, tradePointId uint8) (int, error)
	MergeInstrumentBoard(ctx context.Context, ib models.InstrumentBoard) error
}

func NewService(ctx context.Context, instrRepo InstrumentRepository, logger *zap.Logger) *Service {
	return &Service{
		logger:    logger,
		instrRepo: instrRepo,
	}
}

func (s *Service) SaveInstrument(ctx context.Context) error {
	instr, instrumentBoard, instrumentClass, err := s.instrRepo.SelectInstrumentFromNewCurrentQuote(ctx)
	if err != nil {
		return err
	}

	//1. сначала инструмент
	id, err := s.instrRepo.GetInstrumentId(ctx, instr.Ticker, instr.TradePointId)
	s.logger.Debug("получен инстуремент для котировки", zap.Int("id", id))

	if err != nil && err != apperrors.ErrNotFound {
		s.logger.Error("ошибка получения inst", zap.String("Ticker", instr.Ticker), zap.Error(err))
		return err
	} else if id > 0 {
		err = s.instrRepo.SetInstrument(ctx, id, instrumentClass)
		s.logger.Debug("успешно обновлен инстуремент для котировки", zap.String("Ticker", instr.Ticker))
		//return nil
	} else if err == apperrors.ErrNotFound || id == 0 {

		id, err = s.instrRepo.InsInstrument(ctx, instr)
		if err != nil && id != 0 {
			s.logger.Error("ошибка создания", zap.String("Ticker", instr.Ticker), zap.Error(err))

			return err
		}
		err = s.instrRepo.SetInstrument(ctx, id, instrumentClass)
		if err != nil {
			s.logger.Error("ошибка обновления", zap.String("Ticker", instr.Ticker), zap.Error(err))
			return err
		}
		s.logger.Debug("успешно создан инструмент", zap.Int("id", id))

		//return nil
	}

	// 2. потом борд
	instrumentBoard.InstrumentId = id
	err = s.instrRepo.MergeInstrumentBoard(ctx, instrumentBoard)
	if err != nil {
		return err
	}

	s.logger.Debug("успешно создана связь борда с инструментом", zap.Int("id", id), zap.Uint8("BoardId", instrumentBoard.BoardId))

	return nil
}
