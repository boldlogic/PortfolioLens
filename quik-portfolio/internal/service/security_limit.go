package service

import (
	"context"
	"time"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/apperrors"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
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
