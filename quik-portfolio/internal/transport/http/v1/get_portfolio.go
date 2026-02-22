package v1

import (
	"net/http"

	md "github.com/boldlogic/PortfolioLens/pkg/models"
)

func (h *Handler) GetPortfolio(r *http.Request) (any, string, error) {
	ctx := r.Context()

	items, err := h.service.GetPortfolio(ctx)
	if err != nil {
		return nil, "", err
	}

	var resp []portfolioItemDTO
	for _, it := range items {
		dto := portfolioItemDTO{
			LoadDate:       it.LoadDate.Format(md.DateFormat),
			ClientCode:     it.ClientCode,
			Ticker:         it.Ticker,
			TradeAccount:   it.TradeAccount,
			FirmCode:       it.FirmCode,
			FirmName:       it.FirmName,
			Balance:        it.Balance,
			AcquisitionCcy: it.AcquisitionCcy,
			MvRub:          it.MvRub,
		}
		if it.ISIN != nil {
			dto.ISIN = *it.ISIN
		}
		if it.MvCurrency != nil {
			dto.MvCurrency = *it.MvCurrency
		}
		if it.ShortName != nil {
			dto.ShortName = *it.ShortName
		}
		resp = append(resp, dto)
	}

	return resp, "", nil
}

type portfolioItemDTO struct {
	LoadDate       string  `json:"loadDate"`
	ClientCode     string  `json:"clientCode"`
	Ticker         string  `json:"ticker"`
	TradeAccount   string  `json:"tradeAccount"`
	FirmCode       string  `json:"firmCode"`
	FirmName       string  `json:"firmName"`
	Balance        float64 `json:"balance"`
	AcquisitionCcy string  `json:"acquisitionCcy"`
	ISIN           string  `json:"isin,omitempty"`
	MvCurrency     string  `json:"mvCurrency,omitempty"`
	MvRub          float64 `json:"mvRub"`
	ShortName      string  `json:"shortName,omitempty"`
}
