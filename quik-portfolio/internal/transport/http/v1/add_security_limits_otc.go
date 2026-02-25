package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/boldlogic/PortfolioLens/pkg/utils"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/apperrors"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
	"go.uber.org/zap"
)

func (h *Handler) AddSecurityLimitOtc(r *http.Request) (any, string, error) {
	ctx := r.Context()
	lim, err := h.readSecurityLimitOtcReq(r)
	if err != nil {
		return nil, err.Error(), apperrors.ErrValidation
	}
	err = h.service.SaveSLOtc(ctx, lim)
	if err != nil {
		if errors.Is(err, apperrors.ErrSettleCode) {
			return nil, err.Error(), apperrors.ErrBusinessValidation
		}
		return nil, "", err
	}
	return nil, "", nil
}

func (h *Handler) readSecurityLimitOtcReq(r *http.Request) (models.SecurityLimit, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		h.logger.Warn("не удалось прочитать тело запроса", zap.Error(err))
		return models.SecurityLimit{}, fmt.Errorf("Некорректный формат запроса")
	}
	var req securityLimitOtcReqDTO
	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		h.logger.Warn("не удалось декодировать тело запроса", zap.Error(err))
		return models.SecurityLimit{}, fmt.Errorf("Некорректный формат запроса")
	}
	var date time.Time
	if req.LoadDate != "" {
		date, err = utils.ParseDateDefault(req.LoadDate)
		if err != nil {
			return models.SecurityLimit{}, fmt.Errorf("Некорректный формат loadDate. Ожидается YYYY-MM-DD")
		}
	} else {
		date = time.Now()
	}
	if req.ClientCode == "" {
		return models.SecurityLimit{}, fmt.Errorf("clientCode должен быть заполнен")
	}
	if req.Ticker == "" {
		return models.SecurityLimit{}, fmt.Errorf("ticker должен быть заполнен")
	}
	if req.FirmName == "" {
		return models.SecurityLimit{}, fmt.Errorf("firmName должен быть заполнен")
	}
	isin := (*string)(nil)
	if req.ISIN != "" {
		isin = &req.ISIN
	}
	return models.SecurityLimit{
		LoadDate:       date,
		ClientCode:     req.ClientCode,
		Ticker:         req.Ticker,
		FirmName:       req.FirmName,
		Balance:        req.Balance,
		AcquisitionCcy: req.AcquisitionCcy,
		ISIN:           isin,
	}, nil
}

type securityLimitOtcReqDTO struct {
	LoadDate       string  `json:"loadDate,omitempty"`
	ClientCode     string  `json:"clientCode"`
	Ticker         string  `json:"ticker"`
	FirmName       string  `json:"firmName"`
	Balance        float64 `json:"balance"`
	AcquisitionCcy string  `json:"acquisitionCcy,omitempty"`
	ISIN           string  `json:"isin,omitempty"`
}
