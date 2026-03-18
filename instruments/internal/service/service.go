package service

import (
	"context"
	"errors"

	"github.com/boldlogic/PortfolioLens/instruments/internal/models"
	md "github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/models/quik"
	"go.uber.org/zap"
)

type Service struct {
	logger        *zap.Logger
	instrRepo     InstrumentRepository
	refsSyncRepo  RefsSyncRepository
	refsQueryRepo RefsQueryRepository
}

type InstrumentRepository interface {
	SelectInstrumentFromNewCurrentQuote(ctx context.Context) (models.Instrument, models.InstrumentBoard, string, error)
	InsInstrument(ctx context.Context, i models.Instrument) (int, error)
	SetInstrument(ctx context.Context, id int, ic string) error
	GetInstrumentId(ctx context.Context, ticker string, tradePointId uint8) (int, error)
	MergeInstrumentBoard(ctx context.Context, ib models.InstrumentBoard) error
}

// RefsSyncRepository — запись справочников (воркеры актуализации).
type RefsSyncRepository interface {
	SyncInstrumentTypesFromQuotes(ctx context.Context) error
	SyncInstrumentSubTypesFromQuotes(ctx context.Context) error
	SyncBoardsFromQuotes(ctx context.Context) error
	TagBoardsTradePointId(ctx context.Context) error
}

// RefsQueryRepository — чтение справочников (HTTP API).
type RefsQueryRepository interface {
	GetTradePoints(ctx context.Context) ([]md.TradePoint, error)
	GetTradePointByID(ctx context.Context, id uint8) (md.TradePoint, error)
	GetBoardsWithTradePoint(ctx context.Context) ([]quik.Board, error)
	GetBoardByIDWithTradePoint(ctx context.Context, id uint8) (quik.Board, error)
}

func NewService(instrRepo InstrumentRepository, refsSyncRepo RefsSyncRepository, refsQueryRepo RefsQueryRepository, logger *zap.Logger) *Service {
	return &Service{
		logger:        logger,
		instrRepo:     instrRepo,
		refsSyncRepo:  refsSyncRepo,
		refsQueryRepo: refsQueryRepo,
	}
}

func (s *Service) SaveInstrument(ctx context.Context) error {
	instr, instrumentBoard, instrumentClass, err := s.instrRepo.SelectInstrumentFromNewCurrentQuote(ctx)
	if err != nil {
		return err
	}

	id, err := s.instrRepo.GetInstrumentId(ctx, instr.Ticker, instr.TradePointId)
	s.logger.Debug("получен инструмент для котировки", zap.Int("id", id))

	if err != nil && !errors.Is(err, md.ErrNotFound) {
		s.logger.Error("ошибка получения инструмента", zap.String("Ticker", instr.Ticker), zap.Error(err))
		return err
	} else if id > 0 {
		if err = s.instrRepo.SetInstrument(ctx, id, instrumentClass); err != nil {
			s.logger.Error("ошибка обновления котировки", zap.String("Ticker", instr.Ticker), zap.Error(err))
			return err
		}
		s.logger.Debug("успешно обновлён инструмент для котировки", zap.String("Ticker", instr.Ticker))
	} else {
		id, err = s.instrRepo.InsInstrument(ctx, instr)
		if err != nil {
			s.logger.Error("ошибка создания инструмента", zap.String("Ticker", instr.Ticker), zap.Error(err))
			return err
		}
		if err = s.instrRepo.SetInstrument(ctx, id, instrumentClass); err != nil {
			s.logger.Error("ошибка обновления котировки", zap.String("Ticker", instr.Ticker), zap.Error(err))
			return err
		}
		s.logger.Debug("успешно создан инструмент", zap.Int("id", id))
	}

	instrumentBoard.InstrumentId = id
	if err = s.instrRepo.MergeInstrumentBoard(ctx, instrumentBoard); err != nil {
		return err
	}

	s.logger.Debug("успешно создана связь борда с инструментом", zap.Int("id", id), zap.Uint8("BoardId", instrumentBoard.BoardId))
	return nil
}
