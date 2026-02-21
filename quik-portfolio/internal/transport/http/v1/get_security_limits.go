package v1

import (
	"net/http"

	md "github.com/boldlogic/PortfolioLens/pkg/models"
	"go.uber.org/zap"
)

func (h *Handler) GetSecurityLimits(r *http.Request) (any, error) {
	ctx := r.Context()
	sls, err := h.service.GetSL(ctx)
	h.logger.Debug("", zap.Error(err), zap.Any("lim", sls))

	if err != nil {
		return nil, err
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

	return resp, nil
}

type securityLimitDTO struct {
	LoadDate       string  `json:"loadDate"`
	ClientCode     string  `json:"clientCode"`
	Ticker         string  `json:"ticker"`
	TradeAccount   string  `json:"tradeAccount"`
	FirmName       string  `json:"firmName"`
	Balance        float64 `json:"balance"`
	AcquisitionCcy string  `json:"acquisitionCcy"`
	ISIN           string  `json:"isin,omitempty"`
}
