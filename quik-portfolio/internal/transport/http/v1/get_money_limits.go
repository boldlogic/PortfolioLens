package v1

import (
	"net/http"

	md "github.com/boldlogic/PortfolioLens/pkg/models"
)

func (h *Handler) GetMoneyLimits(r *http.Request) (any, error) {

	ctx := r.Context()
	mls, err := h.service.GetML(ctx)
	if err != nil {
		return nil, err
	}
	var resp []moneyLimitdto
	for _, ml := range mls {

		resp = append(resp, moneyLimitdto{
			LoadDate:   ml.LoadDate.Format(md.DateFormat),
			ClientCode: ml.ClientCode,
			Currency:   ml.Currency,
			FirmName:   ml.FirmName,
			Balance:    ml.Balance,
		})
	}
	return resp, nil
}

type moneyLimitdto struct {
	LoadDate   string  `json:"loadDate"`
	ClientCode string  `json:"clientCode"`
	Currency   string  `json:"currency"`
	FirmName   string  `json:"firmName"`
	Balance    float64 `json:"balance"`
}
