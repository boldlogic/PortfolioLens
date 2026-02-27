package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/apperrors"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
	"github.com/go-chi/chi"
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

	return tradePointsToDTO(res), "", nil
}

func (h *Handler) GetTradePoint(r *http.Request) (any, string, error) {
	id64, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 8)
	if err != nil {
		return nil, "некорректный id торговой площадки", apperrors.ErrValidation
	}
	res, err := h.service.GetTradePointByID(r.Context(), uint8(id64))
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, "торговая площадка не найдена", err
		}
		return nil, "", err
	}
	return tradePointToDTO(res), "", nil
}

func tradePointsToDTO(tradepoints []models.TradePoint) []TradePointDTO {
	var res = make([]TradePointDTO, 0, len(tradepoints))

	for _, tr := range tradepoints {
		res = append(res, tradePointToDTO(tr))
	}
	return res

}

func tradePointToDTO(tradepoint models.TradePoint) TradePointDTO {
	return TradePointDTO{
		Id:   tradepoint.Id,
		Code: tradepoint.Code,
		Name: tradepoint.Name,
	}
}

type TradePointDTO struct {
	Id   uint8  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}
