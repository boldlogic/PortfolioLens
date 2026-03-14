package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/utils"
	qmodels "github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
	"go.uber.org/zap"
)

func (h *Handler) AddSecurityLimit(r *http.Request) (any, string, error) {
	ctx := r.Context()
	lim, err := h.readSecurityLimit(r)
	if err != nil {
		return nil, err.Error(), models.ErrValidation
	}
	err = h.service.SaveSL(ctx, lim)
	if err != nil {
		if errors.Is(err, models.ErrBusinessValidation) {
			return nil, "settleCode должен быть T0, T1, T2 или Tx", err
		}
		return nil, "", err
	}
	return nil, "", nil
}

func (h *Handler) readSecurityLimit(r *http.Request) (qmodels.SecurityLimit, error) {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		h.logger.Warn("не удалось прочитать тело запроса", zap.Error(err))
		return qmodels.SecurityLimit{}, fmt.Errorf("некорректный формат запроса")
	}
	var req securityLimitReqDTO
	if err := json.Unmarshal(buf.Bytes(), &req); err != nil {
		h.logger.Warn("не удалось декодировать тело запроса", zap.Error(err))
		return qmodels.SecurityLimit{}, fmt.Errorf("некорректный формат запроса")
	}

	var date time.Time
	var err error
	if req.LoadDate != "" {
		date, err = utils.ParseDateDefault(req.LoadDate)
		if err != nil {
			return qmodels.SecurityLimit{}, fmt.Errorf("некорректный формат loadDate. Ожидается YYYY-MM-DD")
		}
	} else {
		date = time.Now()
	}

	if req.ClientCode == "" {
		return qmodels.SecurityLimit{}, fmt.Errorf("clientCode должен быть заполнен")
	}
	if req.Ticker == "" {
		return qmodels.SecurityLimit{}, fmt.Errorf("ticker должен быть заполнен")
	}
	if req.TradeAccount == "" {
		return qmodels.SecurityLimit{}, fmt.Errorf("tradeAccount должен быть заполнен")
	}
	if req.FirmName == "" {
		return qmodels.SecurityLimit{}, fmt.Errorf("firmName должен быть заполнен")
	}

	return qmodels.SecurityLimit{
		LoadDate:       date,
		ClientCode:     req.ClientCode,
		Ticker:         req.Ticker,
		TradeAccount:   req.TradeAccount,
		SettleCode:     req.SettleCode,
		FirmName:       req.FirmName,
		Balance:        req.Balance,
		AcquisitionCcy: req.AcquisitionCcy,
		ISIN:           &req.ISIN,
	}, nil
}

type securityLimitReqDTO struct {
	LoadDate       string  `json:"loadDate,omitempty"`
	ClientCode     string  `json:"clientCode"`
	Ticker         string  `json:"ticker"`
	TradeAccount   string  `json:"tradeAccount"`
	SettleCode     string  `json:"settleCode,omitempty"`
	FirmName       string  `json:"firmName"`
	Balance        float64 `json:"balance"`
	AcquisitionCcy string  `json:"acquisitionCcy,omitempty"`
	ISIN           string  `json:"isin,omitempty"`
}
