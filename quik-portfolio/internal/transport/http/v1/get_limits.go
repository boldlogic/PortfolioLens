package v1

import (
	"net/http"

	md "github.com/boldlogic/PortfolioLens/pkg/models"
)

func (h *Handler) GetLimits(r *http.Request) (any, error) {

	ctx := r.Context()
	lim, err := h.service.GetLimits(ctx)

	if err != nil {
		return nil, err
	}

	var resp []limitDTO

	for _, l := range lim {
		resp = append(resp, limitDTO{
			LoadDate:       l.LoadDate.Format(md.DateFormat),
			ClientCode:     l.ClientCode,
			Ticker:         l.Ticker,
			ISIN:           l.ISIN,
			FirmCode:       l.FirmCode,
			FirmName:       l.FirmName,
			Balance:        l.Balance,
			AcquisitionCcy: l.AcquisitionCcy,
		})
	}

	return resp, nil

}

type limitDTO struct {
	LoadDate       string  `json:"loadDate,omitempty"`
	ClientCode     string  `json:"clientCode,omitempty"`
	Ticker         string  `json:"ticker,omitempty"`
	ISIN           *string `json:"isin,omitempty"`
	FirmCode       string  `json:"firmCode,omitempty"`
	FirmName       string  `json:"firmName,omitempty"`
	Balance        float64 `json:"balance,omitempty"`
	AcquisitionCcy string  `json:"acquisitionCcy,omitempty"`
}
