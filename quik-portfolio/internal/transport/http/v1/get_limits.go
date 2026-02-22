package v1

import (
	"fmt"
	"net/http"
	"time"

	md "github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/utils"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/apperrors"
)

func (h *Handler) readGetLimitsRequest(r *http.Request) (*time.Time, error) {
	var date time.Time
	dateReq := r.URL.Query().Get("date")
	if dateReq != "" {
		date, err := utils.ParseDate(dateReq)
		if err != nil {
			return nil, fmt.Errorf("Некорректный формат date. Ожидается YYYY-MM-DD")
		}
		return date, nil
	}

	date = time.Now()

	return &date, nil
}

func (h *Handler) GetLimits(r *http.Request) (any, string, error) {

	ctx := r.Context()
	date, err := h.readGetLimitsRequest(r)
	if err != nil {
		return nil, err.Error(), apperrors.ErrValidation
	}
	lim, err := h.service.GetLimits(ctx, *date)

	if err != nil {
		return nil, "", err
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

	return resp, "", nil

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
