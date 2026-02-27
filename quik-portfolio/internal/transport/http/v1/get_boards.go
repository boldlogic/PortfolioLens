package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/apperrors"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
	"github.com/go-chi/chi"
)

func (h *Handler) GetBoards(r *http.Request) (any, string, error) {
	res, err := h.service.GetBoards(r.Context())
	if err != nil {
		return nil, "", err
	}
	return boardsToDTO(res), "", nil
}

func (h *Handler) GetBoard(r *http.Request) (any, string, error) {
	id64, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 8)
	if err != nil {
		return nil, "некорректный id борда", apperrors.ErrValidation
	}
	res, err := h.service.GetBoardByID(r.Context(), uint8(id64))
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, "борд не найден", err
		}
		return nil, "", err
	}
	return boardToDTO(res), "", nil
}

func boardToDTO(b models.Board) BoardDTO {
	out := BoardDTO{
		Id:       b.Id,
		Code:     b.Code,
		Name:     b.Name,
		IsTraded: b.IsTraded,
	}
	if b.TradePoint != nil {
		tr := tradePointToDTO(*b.TradePoint)
		out.TradePoint = &tr
	}

	return out
}

func boardsToDTO(boards []models.Board) []BoardDTO {
	out := make([]BoardDTO, 0, len(boards))
	for _, b := range boards {
		out = append(out, boardToDTO(b))
	}
	return out
}

type BoardDTO struct {
	Id         uint8          `json:"id"`
	Code       string         `json:"code"`
	Name       string         `json:"name"`
	IsTraded   bool           `json:"isTraded"`
	TradePoint *TradePointDTO `json:"tradePoint,omitempty"`
}
