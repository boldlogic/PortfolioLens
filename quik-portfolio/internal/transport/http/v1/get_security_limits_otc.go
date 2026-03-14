package v1

import (
	"errors"
	"fmt"
	"net/http"

	md "github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/utils"
)

func (h *Handler) GetSecurityLimitsOtc(r *http.Request) (any, string, error) {
	ctx := r.Context()
	date, err := h.readGetLimitsRequest(r)
	if err != nil {
		return nil, err.Error(), md.ErrValidation
	}

	sls, err := h.service.GetSLOtc(ctx, *date)
	if err != nil {
		if errors.Is(err, md.ErrNotFound) {
			return nil, fmt.Sprintf("позиции по бумагам за %s не найдены", date.Format(utils.DateFormat)), err
		}
		return nil, "", err
	}
	return convertSecurityLimit(sls), "", nil
}
