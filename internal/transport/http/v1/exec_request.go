package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type CbRequest struct {
	Type string `json:"type"`
	Uuid string `json:"uuid,omitempty"`
}

func (h *Handler) execRequest(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Получен запрос на выполнение операции")
	//ctx := r.Context()
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		h.log.Warnf("Не удалось прочитать тело запроса: %v", err)
		h.SendResponse(w, APIResponse{
			StatusCode: http.StatusBadRequest,
			Body: Body{
				Error: err.Error(),
			},
		})
		return
	}
	var cb CbRequest
	err = json.Unmarshal(buf.Bytes(), &cb)
	if err != nil {
		h.log.Warnf("Не удалось распарсить тело запроса: %v", err)
		h.SendResponse(w, APIResponse{
			StatusCode: http.StatusBadRequest,
			Body: Body{
				Error: err.Error(),
			},
		})
		return
	}
	h.SendResponse(w, APIResponse{
		StatusCode: http.StatusOK,
	})

}
