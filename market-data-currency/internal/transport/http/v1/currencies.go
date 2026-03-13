package v1

import (
	"net/http"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func (h *Handler) GetCurrencies(r *http.Request) (any, string, error) {
	currencies, err := h.currencySvc.GetCurrencies(r.Context())
	if err != nil {
		h.logger.Error("ошибка при получении валют", zap.Error(err))
		return nil, err.Error(), err
	}

	result := make([]currencyDTO, 0, len(currencies))
	for _, c := range currencies {
		result = append(result, currencyToDTO(c))
	}
	return result, "", nil
}

func (h *Handler) GetCurrency(r *http.Request) (any, string, error) {
	code := chi.URLParam(r, "code")

	if code == "" {
		return nil, "код валюты не может быть пустым", models.ErrValidation
	}

	ccy, detail, err := h.currencySvc.GetCurrency(r.Context(), code)
	if err != nil {
		h.logger.Error("ошибка при получении валюты", zap.Error(err), zap.String("detail", detail), zap.String("code", code))

		return nil, detail, err
	}

	return currencyToDTO(ccy), "", nil
}
