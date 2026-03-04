package service

import (
	"context"
	"time"

	md "github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/models/quik"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
	"go.uber.org/zap"
)

type Service struct {
	logger       *zap.Logger
	quikRefsRepo QuikRefsRepository
	limitsRepo   LimitsRepository
}

type QuikRefsRepository interface {
	GetInstrumentTypeId(ctx context.Context, title string) (quik.InstrumentType, error)
	InsInstrumentType(ctx context.Context, title string) (quik.InstrumentType, error)
	SyncInstrumentTypesFromQuotes(ctx context.Context) error
	GetInstrumentSubTypeId(ctx context.Context, typeId uint8, title string) (quik.InstrumentSubType, error)
	InsInstrumentSubType(ctx context.Context, typeId uint8, title string) (quik.InstrumentSubType, error)
	SyncInstrumentSubTypesFromQuotes(ctx context.Context) error
	SyncBoardsFromQuotes(ctx context.Context) error
	TagBoardsTradePointId(ctx context.Context) error

	GetTradePoints(ctx context.Context) ([]md.TradePoint, error)
	GetTradePointByID(ctx context.Context, id uint8) (md.TradePoint, error)
	GetBoards(ctx context.Context) ([]quik.Board, error)
	GetBoardByID(ctx context.Context, id uint8) (quik.Board, error)
	GetBoardByIDWithTradePoint(ctx context.Context, id uint8) (quik.Board, error)
	GetBoardsWithTradePoint(ctx context.Context) ([]quik.Board, error)
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
	InsertFirm(ctx context.Context, code string, name string) (quik.Firm, error)
	GetFirmByName(ctx context.Context, name string) (quik.Firm, error)
}

func NewService(ctx context.Context, quikRefsRepo QuikRefsRepository, limitsRepo LimitsRepository, logger *zap.Logger) *Service {
	return &Service{
		logger:       logger,
		quikRefsRepo: quikRefsRepo,
		limitsRepo:   limitsRepo,
	}
}
