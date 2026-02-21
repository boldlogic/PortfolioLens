package v1

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/sirupsen/logrus"
)

func (h *Handler) SendResponse(w http.ResponseWriter, response APIResponse) {

	h.setHeaders(w, response.Headers)

	if response.Body == nil {
		w.WriteHeader(response.StatusCode)
		return
	}

	body, err := json.Marshal(response.Body)
	if err != nil {
		h.log.WithError(err).Error("Не удалось сериализовать тело ответа")
		w.WriteHeader(http.StatusInternalServerError)
		_, writeErr := w.Write([]byte(`{"error":"Произошла непредвиденная ошибка"}`))
		if writeErr != nil {
			h.log.WithError(writeErr).Warn("Не удалось записать fallback-тело ответа")
		}
		return
	}

	headersCount, headersBytes := headerMetrics(w.Header())
	bodyType := typeName(response.Body)
	bodyBytes := len(body)

	if response.StatusCode >= 500 {
		h.log.WithFields(logrus.Fields{
			"StatusCode":   response.StatusCode,
			"HeadersCount": headersCount,
			"HeadersBytes": headersBytes,
			"BodyType":     bodyType,
			"BodyBytes":    bodyBytes,
		}).Error("Ответ с ошибкой сервера")

		h.log.WithFields(logrus.Fields{
			"StatusCode": response.StatusCode,
			"Headers":    w.Header(),
			"Body":       response.Body,
		}).Debug("Ответ с ошибкой сервера")
	} else if response.StatusCode >= 400 {
		h.log.WithFields(logrus.Fields{
			"StatusCode":   response.StatusCode,
			"HeadersCount": headersCount,
			"HeadersBytes": headersBytes,
			"BodyType":     bodyType,
			"BodyBytes":    bodyBytes,
		}).Warn("Ответ с ошибкой клиента")

		h.log.WithFields(logrus.Fields{
			"StatusCode": response.StatusCode,
			"Headers":    w.Header(),
			"Body":       response.Body,
		}).Debug("Ответ с ошибкой клиента")
	}

	w.WriteHeader(response.StatusCode)
	_, writeErr := w.Write(body)
	if writeErr != nil {
		h.log.WithError(writeErr).Warn("Не удалось записать тело ответа")
	}
	if writeErr == nil && response.StatusCode >= 200 && response.StatusCode < 400 {
		h.log.WithFields(logrus.Fields{
			"StatusCode":   response.StatusCode,
			"HeadersCount": headersCount,
			"HeadersBytes": headersBytes,
			"BodyType":     bodyType,
			"BodyBytes":    bodyBytes,
		}).Info("Успешный ответ")
	}
}

func (h *Handler) setHeaders(w http.ResponseWriter, headers map[string]string) {
	for k, v := range headers {
		w.Header().Set(k, v)
	}

	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	}
}

func headerMetrics(headers http.Header) (count int, bytes int) {
	if len(headers) == 0 {
		return 0, 0
	}
	for k, values := range headers {
		count++
		bytes += len(k)
		for _, v := range values {
			bytes += len(v)
		}
	}
	return count, bytes
}

func typeName(v any) string {
	if v == nil {
		return "<nil>"
	}
	return reflect.TypeOf(v).String()
}
