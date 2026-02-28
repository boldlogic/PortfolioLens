package service

import (
	"context"
	"time"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/apperrors"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
	"go.uber.org/zap"
)

type Service struct {
	logger       *zap.Logger
	instrRepo    InstrumentRepository
	quikRefsRepo QuikRefsRepository
	limitsRepo   LimitsRepository
}

type InstrumentRepository interface {
	SelectInstrumentFromNewCurrentQuote(ctx context.Context) (models.Instrument, string, error)
	InsInstrument(ctx context.Context, i models.Instrument) (int, error)
	SetInstrument(ctx context.Context, id int, ic string) error
	GetInstrumentId(ctx context.Context, ticker string) (int, error)
}

type QuikRefsRepository interface {
	GetInstrumentTypeId(ctx context.Context, title string) (models.InstrumentType, error)
	InsInstrumentType(ctx context.Context, title string) (models.InstrumentType, error)
	SyncInstrumentTypesFromQuotes(ctx context.Context) error
	GetInstrumentSubTypeId(ctx context.Context, typeId uint8, title string) (models.InstrumentSubType, error)
	InsInstrumentSubType(ctx context.Context, typeId uint8, title string) (models.InstrumentSubType, error)
	SyncInstrumentSubTypesFromQuotes(ctx context.Context) error
	SyncBoardsFromQuotes(ctx context.Context) error
	TagBoardsTradePointId(ctx context.Context) error

	GetTradePoints(ctx context.Context) ([]models.TradePoint, error)
	GetTradePointByID(ctx context.Context, id uint8) (models.TradePoint, error)
	GetBoards(ctx context.Context) ([]models.Board, error)
	GetBoardByID(ctx context.Context, id uint8) (models.Board, error)
	GetBoardByIDWithTradePoint(ctx context.Context, id uint8) (models.Board, error)
	GetBoardsWithTradePoint(ctx context.Context) ([]models.Board, error)
}

type LimitsRepository interface {
	GetMoneyLimits(ctx context.Context, date time.Time) ([]models.MoneyLimit, error)

	GetSecurityLimits(ctx context.Context, date time.Time) ([]models.SecurityLimit, error)
	SaveSecurityLimit(ctx context.Context, s models.SecurityLimit) error

	GetSecurityLimitsOtc(ctx context.Context, date time.Time) ([]models.SecurityLimit, error)
	SaveSecurityLimitOtc(ctx context.Context, s models.SecurityLimit) error
	GetSecurityLimitsOtcMaxDate(ctx context.Context) (*time.Time, error)
	RollSecurityLimitsOtcFromDateToDate(ctx context.Context, dateFrom time.Time, dateTo time.Time) error
	DeleteSecurityLimitsOtcBeforeDate(ctx context.Context, date time.Time) error

	GetPortfolio(ctx context.Context) ([]models.PortfolioItem, error)
	InsertFirm(ctx context.Context, code string, name string) (models.Firm, error)
	GetFirmByName(ctx context.Context, name string) (models.Firm, error)
}

func NewService(ctx context.Context, intrRepo InstrumentRepository, quikRefsRepo QuikRefsRepository, limitsRepo LimitsRepository, logger *zap.Logger) *Service {
	return &Service{
		logger:       logger,
		instrRepo:    intrRepo,
		quikRefsRepo: quikRefsRepo,
		limitsRepo:   limitsRepo,
	}
}

func (s *Service) SaveInstrument(ctx context.Context) error {
	instr, instrumentClass, err := s.instrRepo.SelectInstrumentFromNewCurrentQuote(ctx)
	if err != nil {
		return err
	}

	id, err := s.instrRepo.GetInstrumentId(ctx, instr.Ticker)
	s.logger.Debug("успешно получен инстуремент для котировки", zap.Int("id", id))

	if err != nil && err != apperrors.ErrNotFound {
		s.logger.Error("ошибка получения inst", zap.String("Ticker", instr.Ticker), zap.Error(err))
		return err
	} else if id > 0 {
		err = s.instrRepo.SetInstrument(ctx, id, instrumentClass)
		s.logger.Debug("успешно обновлен инстуремент для котировки", zap.String("Ticker", instr.Ticker))
		return nil
	} else if err == apperrors.ErrNotFound || id == 0 {

		id, err = s.instrRepo.InsInstrument(ctx, instr)
		if err != nil && id != 0 {
			s.logger.Error("ошибка создания", zap.String("Ticker", *&instr.Ticker), zap.Error(err))

			return err
		}
		err = s.instrRepo.SetInstrument(ctx, id, instrumentClass)
		if err != nil {
			s.logger.Error("ошибка обновления", zap.String("Ticker", *&instr.Ticker), zap.Error(err))
			return err
		}
		s.logger.Debug("успешно создан инструмент", zap.Int("id", id))

		return nil
	}
	return nil
}
