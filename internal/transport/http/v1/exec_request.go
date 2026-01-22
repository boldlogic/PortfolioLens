package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func (h *Handler) execRequest(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Получен запрос на выполнение операции")
	ctx := r.Context()
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
	err = h.Service.Executor(ctx, cb.Type)
	//h.log.Info(resp)

	h.SendResponse(w, APIResponse{
		StatusCode: http.StatusOK,
	})

}
