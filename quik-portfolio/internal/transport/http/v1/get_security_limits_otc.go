package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/boldlogic/PortfolioLens/pkg/utils"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/apperrors"
)

func (h *Handler) GetSecurityLimitsOtc(r *http.Request) (any, string, error) {
	ctx := r.Context()
	date, err := h.readGetLimitsRequest(r)
	if err != nil {
		return nil, err.Error(), apperrors.ErrValidation
	}

	sls, err := h.service.GetSLOtc(ctx, *date)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, fmt.Sprintf("позиции по бумагам за %s не найдены", date.Format(utils.DateFormat)), err
		}
		return nil, "", err
	}

	return convertSecurityLimit(sls), "", nil
}
