package service

import (
	"context"
	"time"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/apperrors"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
	"go.uber.org/zap"
)

func (s *Service) GetSL(ctx context.Context, date time.Time) ([]models.SecurityLimit, error) {
	sl, err := s.limitsRepo.GetSecurityLimits(ctx, date)
	if err != nil {
		return nil, err
	}
	return sl, nil
}

func (s *Service) GetSLOtc(ctx context.Context, date time.Time) ([]models.SecurityLimit, error) {
	sl, err := s.limitsRepo.GetSecurityLimitsOtc(ctx, date)
	if err != nil {
		return nil, err
	}
	return sl, nil
}

func (s *Service) SaveSL(ctx context.Context, sec models.SecurityLimit) error {

	firm, err := s.limitsRepo.GetFirmByName(ctx, sec.FirmName)
	if err != nil {
		return apperrors.ErrSavingData
	}
	sec.FirmCode = firm.Code

	if err = checkSettleCode(sec.SettleCode); err != nil {
		return err
	}

	err = s.limitsRepo.SaveSecurityLimit(ctx, sec)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) SaveSLOtc(ctx context.Context, sec models.SecurityLimit) error {
	firm, err := s.limitsRepo.GetFirmByName(ctx, sec.FirmName)
	if err != nil {
		return apperrors.ErrSavingData
	}
	sec.FirmCode = firm.Code
	sec.TradeAccount = "OTC"
	sec.SettleCode = "Tx"
	err = s.limitsRepo.SaveSecurityLimitOtc(ctx, sec)
	if err != nil {
		return err
	}
	return nil
}

func checkSettleCode(code string) error {
	allowedSettle := map[string]bool{"T0": true, "T1": true, "T2": true, "Tx": true}
	if !allowedSettle[code] {
		return apperrors.ErrSettleCode
	}
	return nil
}

func (s *Service) RollForwardSecurityLimitsOtc(ctx context.Context) error {
	for {
		timer := time.NewTimer(60 * time.Second)
		select {
		case <-timer.C:
			date, err := s.limitsRepo.GetSecurityLimitsOtcMaxDate(ctx)
			if err != nil {
				s.logger.Error("ошибка получения макс. даты", zap.Error(err))
				continue
			}
			if date == nil {
				continue
			}
			now := time.Now()
			loc := now.Location()
			maxDateOnly := time.Date(date.In(loc).Year(), date.In(loc).Month(), date.In(loc).Day(), 0, 0, 0, 0, loc)
			todayOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
			if !todayOnly.After(maxDateOnly) {
				continue
			}

			if err := s.limitsRepo.RollSecurityLimitsOtcFromDateToDate(ctx, maxDateOnly, todayOnly); err != nil {
				s.logger.Error("ошибка переноса лимитов", zap.Error(err))
				select {
				case <-time.After(5 * time.Second):
				case <-ctx.Done():
					s.logger.Debug("получена команда завершать")
					return nil
				}
				continue
			}
			err = s.limitsRepo.DeleteSecurityLimitsOtcBeforeDate(ctx, maxDateOnly)
			if err != nil {
				s.logger.Error("ошибка удаления лимитов", zap.Error(err))
			}
		case <-ctx.Done():
			timer.Stop()
			s.logger.Debug("получена команда завершать")
			return nil
		}
	}
}
