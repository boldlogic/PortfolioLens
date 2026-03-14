package service

import (
	"context"
	"time"

	"github.com/boldlogic/PortfolioLens/pkg/models/quik"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
	"go.uber.org/zap"
)

type MoneyLimitsRepo interface {
	GetMoneyLimits(ctx context.Context, date time.Time) ([]models.MoneyLimit, error)
}

type SecurityLimitsRepo interface {
	GetSecurityLimits(ctx context.Context, date time.Time) ([]models.SecurityLimit, error)
	SaveSecurityLimit(ctx context.Context, s models.SecurityLimit) error
}

type SecurityLimitsOtcRepo interface {
	GetSecurityLimitsOtc(ctx context.Context, date time.Time) ([]models.SecurityLimit, error)
	SaveSecurityLimitOtc(ctx context.Context, s models.SecurityLimit) error
	GetSecurityLimitsOtcMaxDate(ctx context.Context) (*time.Time, error)
	RollSecurityLimitsOtcFromDateToDate(ctx context.Context, dateFrom time.Time, dateTo time.Time) error
	DeleteSecurityLimitsOtcBeforeDate(ctx context.Context, date time.Time) error
}

type PortfolioRepo interface {
	GetPortfolio(ctx context.Context) ([]models.PortfolioItem, error)
}

type FirmsRepo interface {
	InsertFirm(ctx context.Context, code string, name string) (quik.Firm, error)
	GetFirmByName(ctx context.Context, name string) (quik.Firm, error)
}

type Repository interface {
	MoneyLimitsRepo
	SecurityLimitsRepo
	SecurityLimitsOtcRepo
	PortfolioRepo
	FirmsRepo
}

type Service struct {
	logger *zap.Logger
	repo   Repository
}

func NewService(repo Repository, logger *zap.Logger) *Service {
	return &Service{
		logger: logger,
		repo:   repo,
	}
}
