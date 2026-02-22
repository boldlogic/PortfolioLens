package v1

import (
	"errors"
	"net/http"

	md "github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/apperrors"
	"go.uber.org/zap"
)

func (h *Handler) GetSecurityLimitsOtc(r *http.Request) (any, string, error) {
	ctx := r.Context()
	date, err := h.readGetLimitsRequest(r)
	if err != nil {
		return nil, err.Error(), apperrors.ErrValidation
	}
	sls, err := h.service.GetSLOtc(ctx, *date)
	h.logger.Debug("GetSecurityLimitsOtc", zap.Error(err), zap.Any("lim", sls))
	if err != nil {
		if errors.Is(err, apperrors.ErrSLNotFound) {
			return nil, err.Error(), apperrors.ErrNotFound
		}
		return nil, "", err
	}
	var resp []securityLimitDTO
	for _, sl := range sls {
		dto := securityLimitDTO{
			LoadDate:       sl.LoadDate.Format(md.DateFormat),
			ClientCode:     sl.ClientCode,
			Ticker:         sl.Ticker,
			TradeAccount:   sl.TradeAccount,
			FirmName:       sl.FirmName,
			Balance:        sl.Balance,
			AcquisitionCcy: sl.AcquisitionCcy,
		}
		if sl.ISIN != nil {
			dto.ISIN = *sl.ISIN
		}
		resp = append(resp, dto)
	}
	return resp, "", nil
}
