package v1

import (
	"errors"
	"net/http"
	"strconv"

	md "github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) GetTradePoints(r *http.Request) (any, string, error) {
	res, err := h.refsSvc.GetTradePoints(r.Context())
	if err != nil {
		if errors.Is(err, md.ErrNotFound) {
			return nil, "торговые площадки не найдены", err
		}
		return nil, "", err
	}
	return tradePointsToDTO(res), "", nil
}

func (h *Handler) GetTradePoint(r *http.Request) (any, string, error) {
	id64, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 8)
	if err != nil {
		return nil, "некорректный id торговой площадки", md.ErrValidation
	}
	res, err := h.refsSvc.GetTradePointByID(r.Context(), uint8(id64))
	if err != nil {
		if errors.Is(err, md.ErrNotFound) {
			return nil, "торговая площадка не найдена", err
		}
		return nil, "", err
	}
	return tradePointToDTO(res), "", nil
}

func tradePointsToDTO(tradepoints []md.TradePoint) []TradePointDTO {
	res := make([]TradePointDTO, 0, len(tradepoints))
	for _, tr := range tradepoints {
		res = append(res, tradePointToDTO(tr))
	}
	return res
}

func tradePointToDTO(tradepoint md.TradePoint) TradePointDTO {
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
