package v1

import (
	"errors"
	"fmt"
	"net/http"

	md "github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/utils"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/apperrors"
)

func (h *Handler) GetMoneyLimits(r *http.Request) (any, string, error) {

	ctx := r.Context()
	date, err := h.readGetLimitsRequest(r)
	if err != nil {
		return nil, err.Error(), apperrors.ErrValidation
	}

	mls, err := h.service.GetML(ctx, *date)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, fmt.Sprintf("позиции по деньгам за %s не найдены", date.Format(utils.DateFormat)), err
		}
		return nil, "", err

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
	return resp, "", nil
}

type moneyLimitdto struct {
	LoadDate   string  `json:"loadDate"`
	ClientCode string  `json:"clientCode"`
	Currency   string  `json:"currency"`
	FirmName   string  `json:"firmName"`
	Balance    float64 `json:"balance"`
}
