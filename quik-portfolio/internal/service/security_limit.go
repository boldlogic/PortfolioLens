package service

import (
	"context"
	"time"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	qmodels "github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
)

func (s *Service) GetSL(ctx context.Context, date time.Time) ([]qmodels.SecurityLimit, error) {
	return s.repo.GetSecurityLimits(ctx, date)
}

func (s *Service) GetSLOtc(ctx context.Context, date time.Time) ([]qmodels.SecurityLimit, error) {
	return s.repo.GetSecurityLimitsOtc(ctx, date)
}

func (s *Service) SaveSL(ctx context.Context, sec qmodels.SecurityLimit) error {
	firm, err := s.repo.GetFirmByName(ctx, sec.FirmName)
	if err != nil {
		return err
	}
	sec.FirmCode = firm.Code

	if err = checkSettleCode(sec.SettleCode); err != nil {
		return err
	}

	return s.repo.SaveSecurityLimit(ctx, sec)
}

func (s *Service) SaveSLOtc(ctx context.Context, sec qmodels.SecurityLimit) error {
	firm, err := s.repo.GetFirmByName(ctx, sec.FirmName)
	if err != nil {
		return err
	}
	sec.FirmCode = firm.Code
	sec.TradeAccount = "OTC"
	sec.SettleCode = "Tx"
	return s.repo.SaveSecurityLimitOtc(ctx, sec)
}

func checkSettleCode(code string) error {
	allowedSettle := map[string]bool{"T0": true, "T1": true, "T2": true, "Tx": true}
	if !allowedSettle[code] {
		return models.ErrBusinessValidation
	}
	return nil
}

// DoRollForwardOtc выполняет одну итерацию переноса OTC-лимитов: с макс. даты по сегодня.
func (s *Service) DoRollForwardOtc(ctx context.Context) error {
	date, err := s.repo.GetSecurityLimitsOtcMaxDate(ctx)
	if err != nil {
		return err
	}
	if date == nil {
		return nil
	}
	now := time.Now()
	loc := now.Location()
	maxDateOnly := time.Date(date.In(loc).Year(), date.In(loc).Month(), date.In(loc).Day(), 0, 0, 0, 0, loc)
	todayOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	if !todayOnly.After(maxDateOnly) {
		return nil
	}
	if err := s.repo.RollSecurityLimitsOtcFromDateToDate(ctx, maxDateOnly, todayOnly); err != nil {
		return err
	}
	return s.repo.DeleteSecurityLimitsOtcBeforeDate(ctx, maxDateOnly)
}
