package v1

import (
	"errors"
	"fmt"
	"net/http"

	md "github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/utils"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/apperrors"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
	"go.uber.org/zap"
)

func (h *Handler) GetSecurityLimits(r *http.Request) (any, string, error) {
	ctx := r.Context()

	date, err := h.readGetLimitsRequest(r)
	if err != nil {
		return nil, err.Error(), apperrors.ErrValidation
	}
	sls, err := h.service.GetSL(ctx, *date)
	h.logger.Debug("", zap.Error(err), zap.Any("lim", sls))

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, fmt.Sprintf("позиции по бумагам за %s не найдены", date.Format(utils.DateFormat)), err
		}
		return nil, "", err

	}

	return convertSecurityLimit(sls), "", nil
}

func convertSecurityLimit(sls []models.SecurityLimit) []securityLimitDTO {
	var res []securityLimitDTO
	for _, sl := range sls {
		dto := securityLimitDTO{
			LoadDate:       sl.LoadDate.Format(md.DateFormat),
			ClientCode:     sl.ClientCode,
			Ticker:         sl.Ticker,
			FirmName:       sl.FirmName,
			Balance:        sl.Balance,
			AcquisitionCcy: sl.AcquisitionCcy,
		}
		if sl.ISIN != nil {
			dto.ISIN = *sl.ISIN
		}
		res = append(res, dto)
	}
	return res
}

type securityLimitDTO struct {
	LoadDate       string  `json:"loadDate"`
	ClientCode     string  `json:"clientCode"`
	Ticker         string  `json:"ticker"`
	FirmName       string  `json:"firmName"`
	Balance        float64 `json:"balance"`
	AcquisitionCcy string  `json:"acquisitionCcy"`
	ISIN           string  `json:"isin,omitempty"`
}
