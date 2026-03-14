package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"go.uber.org/zap"
)

func (h *Handler) AddFirm(r *http.Request) (any, string, error) {
	ctx := r.Context()
	req, err := h.readFirm(r)
	if err != nil {
		return nil, err.Error(), models.ErrValidation
	}

	firm, err := h.service.SaveFirm(ctx, req.Code, req.Name)
	if err != nil {
		return nil, "", err
	}
	return firmRespDto{Id: firm.Id, Code: firm.Code, Name: firm.Name}, "", nil
}

type firmReqDto struct {
	Code string `json:"firmCode"`
	Name string `json:"firmName"`
}

type firmRespDto struct {
	Id   uint8  `json:"id"`
	Code string `json:"firmCode"`
	Name string `json:"firmName"`
}

func (h *Handler) readFirm(r *http.Request) (firmReqDto, error) {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		h.logger.Warn("не удалось прочитать тело запроса", zap.Error(err))
		return firmReqDto{}, fmt.Errorf("некорректный формат запроса")
	}
	var req firmReqDto
	if err := json.Unmarshal(buf.Bytes(), &req); err != nil {
		h.logger.Warn("не удалось декодировать тело запроса", zap.Error(err))
		return firmReqDto{}, fmt.Errorf("некорректный формат запроса")
	}
	if req.Code == "" {
		return firmReqDto{}, fmt.Errorf("поле Code должно быть заполнено")
	}
	if req.Name == "" {
		return firmReqDto{}, fmt.Errorf("поле Name должно быть заполнено")
	}
	return req, nil
}
