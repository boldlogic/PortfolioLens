package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func (h *Handler) AddFirm(r *http.Request) (any, error) {

	//ctx := r.Context()
	_, err := h.readFirm(r)
	if err != nil {
		return nil, err
	}
	return nil, nil

}

type firmReqDto struct {
	Code string `json:"firmCode"`
	Name string `json:"firmName"`
}

type firmRespDto struct {
	Id   int8   `json:"id"`
	Code string `json:"firmCode"`
	Name string `json:"firmName"`
}

func (h *Handler) readFirm(r *http.Request) (firmReqDto, error) {

	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		h.logger.Warn("не удалось прочитать тело запроса", zap.Error(err))

		return firmReqDto{}, fmt.Errorf("Некорректный формат запроса")
	}
	var req firmReqDto
	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		h.logger.Warn("не удалось декодировать тело запроса", zap.Error(err))
		return firmReqDto{}, fmt.Errorf("Некорректный формат запроса")
	}
	return req, nil
}
