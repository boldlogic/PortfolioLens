package httputils

import (
	"encoding/json"
	"net/http"
)

type HTTPErr struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title,omitempty"`
	Status   int    `json:"status,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

func WriteResp(w http.ResponseWriter, status int, data any) {
	body, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if _, werr := w.Write(body); werr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func NotFound(detail string) HTTPErr {
	return HTTPErr{
		Title:  "NOT_FOUND",
		Status: 404,
		Detail: detail,
	}
}

func BadRequest(detail string) HTTPErr {
	return HTTPErr{
		Title:  "VALIDATION_ERROR",
		Status: http.StatusBadRequest,
		Detail: detail,
	}
}

func UnprocessableEntity(detail string) HTTPErr {
	return HTTPErr{
		Title:  "BUSINESS_VALIDATION_ERROR",
		Status: http.StatusUnprocessableEntity,
		Detail: detail,
	}
}
func Internal(detail string) HTTPErr {
	return HTTPErr{
		Title:  "SERVER_ERROR",
		Status: 500,
		Detail: "Что-то пошло не так",
	}
}

func Conflict(detail string) HTTPErr {
	return HTTPErr{
		Title:  "CONFLICT",
		Status: 409,
		Detail: detail,
	}
}
