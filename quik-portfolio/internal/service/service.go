package service

import (
	"context"
	"errors"
	"time"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/apperrors"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
	"go.uber.org/zap"
)

type Service struct {
	logger        *zap.Logger
	instrRepo     InstrumentRepository
	instrTypeRepo InstrumentTypeRepository
	limitsRepo    LimitsRepository
}

type InstrumentRepository interface {
	SelectNewCurrentQuote(ctx context.Context) (models.CurrentQuote, error)
	InsInstrument(ctx context.Context, i models.Instrument) (int, error)
	SetInstrument(ctx context.Context, id int, ic string) error
	GetInstrumentId(ctx context.Context, ticker string) (int, error)
}

type InstrumentTypeRepository interface {
	GetInstrumentTypeId(ctx context.Context, title string) (models.InstrumentType, error)
	InsInstrumentType(ctx context.Context, title string) (models.InstrumentType, error)
	SyncInstrumentTypesFromQuotes(ctx context.Context) error
	GetInstrumentSubTypeId(ctx context.Context, typeId int16, title string) (models.InstrumentSubType, error)
	InsInstrumentSubType(ctx context.Context, typeId int16, title string) (models.InstrumentSubType, error)
	SyncInstrumentSubTypesFromQuotes(ctx context.Context) error
	SyncBoardsFromQuotes(ctx context.Context) error
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

func NewService(ctx context.Context, intrRepo InstrumentRepository, instrTypeRepo InstrumentTypeRepository, limitsRepo LimitsRepository, logger *zap.Logger) *Service {
	return &Service{
		logger:        logger,
		instrRepo:     intrRepo,
		instrTypeRepo: instrTypeRepo,
		limitsRepo:    limitsRepo,
	}
}

func (s *Service) SaveInstrument(ctx context.Context) error {
	quote, err := s.instrRepo.SelectNewCurrentQuote(ctx)
	if err != nil {
		return err
	}
	quote.Clear()

	id, err := s.instrRepo.GetInstrumentId(ctx, quote.Ticker)
	s.logger.Debug("успешно получен инстуремент для котировки", zap.Int("id", id))
	if err != nil && err != models.ErrInstrumentNotFound {
		s.logger.Error("ошибка получения inst", zap.String("Ticker", quote.Ticker), zap.Error(err))
		return err
	} else if id > 0 {
		err = s.instrRepo.SetInstrument(ctx, id, quote.InstrumentClass)
		s.logger.Debug("успешно обновлен инстуремент для котировки", zap.String("Ticker", quote.Ticker))
		return nil
	} else if err == models.ErrInstrumentNotFound || id == 0 {

		intrType, err := s.instrTypeRepo.GetInstrumentTypeId(ctx, quote.InstrumentType)
		if err != nil && errors.Is(err, apperrors.ErrNotFound) {
			s.logger.Warn("тип инструмента не найден, создаем", zap.String("тип", quote.InstrumentType), zap.Error(err))

			intrType, err = s.instrTypeRepo.InsInstrumentType(ctx, quote.InstrumentType)

			if err != nil {
				s.logger.Error("ошибка создания", zap.String("тип", quote.InstrumentType), zap.Error(err))

				return err
			}
			s.logger.Debug("создание типа инструмента успешно", zap.String("тип", quote.InstrumentType))

		}
		var instrSubType models.InstrumentSubType
		if quote.InstrumentSubtype != nil {
			instrSubType, err = s.instrTypeRepo.GetInstrumentSubTypeId(ctx, intrType.Id, *quote.InstrumentSubtype)
			if err != nil && errors.Is(err, models.ErrInstrumentSubTypeNotFound) {
				s.logger.Warn("подтип инструмента не найден, создаем", zap.String("подтип", *quote.InstrumentSubtype), zap.Error(err))
				instrSubType, err = s.instrTypeRepo.InsInstrumentSubType(ctx, intrType.Id, *quote.InstrumentSubtype)

				if err != nil {
					s.logger.Error("ошибка создания", zap.String("подтип", *quote.InstrumentSubtype), zap.Error(err))

					return err
				}
				s.logger.Debug("создание подтипа инструмента успешно", zap.String("тип", *quote.InstrumentSubtype), zap.Int16("id", instrSubType.SubTypeId))
			}
		}
		inst := models.Instrument{
			Ticker:             quote.Ticker,
			RegistrationNumber: quote.RegistrationNumber,
			FullName:           quote.FullName,
			ShortName:          quote.ShortName,
			ClassCode:          quote.ClassCode,
			ClassName:          quote.ClassCode,
			TypeId:             intrType.Id,
			SubTypeId:          &instrSubType.SubTypeId,
			ISIN:               quote.ISIN,
			FaceValue:          quote.FaceValue,
			BaseCurrency:       quote.BaseCurrency,
			QuoteCurrency:      quote.QuoteCurrency,
			CounterCurrency:    quote.CounterCurrency,
			MaturityDate:       quote.MaturityDate,
			CouponDuration:     quote.CouponDuration,
		}
		id, err = s.instrRepo.InsInstrument(ctx, inst)
		if err != nil && id != 0 {
			s.logger.Error("ошибка создания", zap.String("Ticker", *&quote.Ticker), zap.Error(err))

			return err
		}
		err = s.instrRepo.SetInstrument(ctx, id, quote.InstrumentClass)
		if err != nil {
			s.logger.Error("ошибка обновления", zap.String("Ticker", *&quote.Ticker), zap.Error(err))
			return err
		}
		s.logger.Debug("успешно создан инструмент", zap.Int("id", id))

		return nil
	}
	return nil
}
