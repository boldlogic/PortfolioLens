package v1

import (
	"net/http"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

func (h *Handler) GetCurrencies(w http.ResponseWriter, r *http.Request) {
	currencies, err := h.Service.CurrencyRepo.GetCurrencies()
	if err != nil {

		if err == gorm.ErrRecordNotFound {
			h.log.Info("Справочник валют пуст")
			h.SendResponse(w, APIResponse{
				StatusCode: http.StatusOK,
				Body:       []currencyDTO{},
			})

			return
		}

		h.SendResponse(w, APIResponse{
			StatusCode: http.StatusInternalServerError,
			Body: Body{
				Error: "Неожиданная ошибка",
			},
		})

		return
	}
	res := getCurrenciesDTO(currencies)
	h.SendResponse(w, APIResponse{
		StatusCode: http.StatusOK,
		Body:       res,
	})

}

func getCurrenciesDTO(curr []models.Currency) []currencyDTO {
	var res []currencyDTO
	for _, i := range curr {
		res = append(res, currencyDTO{
			ISOCode:  i.ISOCode,
			CharCode: i.ISOCharCode,
			NameRu:   i.Name,
			NameEn:   i.LatName,
		})
	}
	return res

}

func (h *Handler) GetCurrency(w http.ResponseWriter, r *http.Request) {

	code := chi.URLParam(r, "code")
	if code == "" {
		h.SendResponse(w, APIResponse{
			StatusCode: http.StatusBadRequest,
			Body: Body{
				Error: "Некорректный код валюты",
			},
		})

		return
	}
	currency, err := h.Service.CurrencyRepo.GetCurrency(code)
	if err != nil {

		if err == gorm.ErrRecordNotFound {
			h.log.Info("Справочник валют пуст")
			h.SendResponse(w, APIResponse{
				StatusCode: http.StatusNotFound,
				Body: Body{
					Error: "Валюта не найдена",
				},
			})

			return
		}

		h.SendResponse(w, APIResponse{
			StatusCode: http.StatusInternalServerError,
			Body: Body{
				Error: "Неожиданная ошибка",
			},
		})

		return
	}
	res := getcurrencyDTO(currency)
	h.SendResponse(w, APIResponse{
		StatusCode: http.StatusOK,
		Body:       res,
	})

}

func getcurrencyDTO(curr models.Currency) currencyDTO {

	return currencyDTO{
		ISOCode:  curr.ISOCode,
		CharCode: curr.ISOCharCode,
		NameRu:   curr.Name,
		NameEn:   curr.LatName,
	}

}
