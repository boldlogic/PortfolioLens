package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/apperrors"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
)

func (h *Handler) GetTradePoints(r *http.Request) (any, string, error) {
	ctx := r.Context()
	res, err := h.service.GetTradePoints(ctx)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, fmt.Sprintf("торговые площадки не найдены"), err
		}
		return nil, "", err

	}

	return assertTradePointToDTO(res), "", nil
}

func assertTradePointToDTO(tradepoints []models.TradePoint) []TradePointDTO {
	var res []TradePointDTO

	for _, tr := range tradepoints {
		res = append(res, TradePointDTO{
			Id:   tr.Id,
			Code: tr.Code,
			Name: tr.Name,
		})
	}
	return res

}

type TradePointDTO struct {
	Id   uint8  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}
