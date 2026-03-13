package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/models/scheduler"
	"go.uber.org/zap"
)

type newTaskDTO struct {
	Action string    `json:"action"`
	Uuid   string    `json:"uuid,omitempty"`
	Params paramsDTO `json:"params,omitempty"`
}

func taskToDTO(t scheduler.Task) newTaskRespDTO {
	return newTaskRespDTO{
		Id:          t.Id,
		Uuid:        t.Uuid.String(),
		CreatedAt:   t.CreatedAt,
		ScheduledAt: t.ScheduledAt,
	}
}

type paramsDTO struct {
	CcyCode  string `json:"ccyCode,omitempty"`
	DateFrom string `json:"dateFrom,omitempty"`
	DateTo   string `json:"dateTo,omitempty"`
}

func buildTaskParams(p paramsDTO) (map[string]string, error) {
	params := make(map[string]string)

	if p.CcyCode != "" {
		params["char_code"] = p.CcyCode
	}

	if (p.DateFrom != "" && p.DateTo == "") ||
		(p.DateFrom == "" && p.DateTo != "") {
		return nil, fmt.Errorf("dateFrom и dateTo должны быть указаны вместе")
	}
	if p.DateFrom != "" {
		if _, err := time.Parse(models.ISODateFormat, p.DateFrom); err != nil {
			return nil, fmt.Errorf("некорректный формат dateFrom, ожидается YYYY-MM-DD")
		}
		params["date_from"] = p.DateFrom
	}
	if p.DateTo != "" {
		if _, err := time.Parse(models.ISODateFormat, p.DateTo); err != nil {
			return nil, fmt.Errorf("некорректный формат dateTo, ожидается YYYY-MM-DD")
		}
		params["date_to"] = p.DateTo
	}
	if params["date_from"] > params["date_to"] {
		return nil, fmt.Errorf("dateFrom не может быть больше dateTo")
	}

	return params, nil
}

type newTaskRespDTO struct {
	Id          int64     `json:"id"`
	Uuid        string    `json:"uuid"`
	CreatedAt   time.Time `json:"createdAt"`
	ScheduledAt time.Time `json:"scheduledAt"`
}

func (h *Handler) CreateTask(r *http.Request) (any, string, error) {
	ctx := r.Context()

	var dto newTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.logger.Warn("некорректный формат запроса", zap.Error(err))
		return nil, "некорректный формат запроса", models.ErrValidation
	}
	if dto.Action == "" {
		h.logger.Warn("некорректный формат запроса: поле 'action' обязательно")

		return nil, "поле 'action' обязательно", models.ErrValidation
	}

	params, err := buildTaskParams(dto.Params)
	if err != nil {
		h.logger.Warn("некорректный запрос", zap.Error(err))

		return nil, err.Error(), models.ErrValidation
	}

	created, err := h.taskSvc.CreateTask(ctx, dto.Action, dto.Uuid, params)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrBusinessValidation):
			h.logger.Warn("некорректный запрос: код action не существует или параметр не найден", zap.Error(err))
			return nil, "некорректный код action или параметр", err
		case errors.Is(err, models.ErrConflict):
			h.logger.Warn("конфликт: задача с таким uuid уже существует", zap.Error(err))
			return nil, "задача с таким uuid уже существует", err
		default:
			h.logger.Error("ошибка при сохранении задачи", zap.Error(err))
			return nil, "внутренняя ошибка", err
		}
	}

	return taskToDTO(created), "", nil
}
