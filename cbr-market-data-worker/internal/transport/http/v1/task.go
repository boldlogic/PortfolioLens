package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/boldlogic/PortfolioLens/pkg/models"
)

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Получен запрос на выполнение операции CreateTask")
	ctx := r.Context()
	reqTask, action, err := parseTask(r.Body)

	if err != nil {
		h.log.Errorf("Ошибка: %w", err)
		h.SendResponse(w, APIResponse{
			StatusCode: http.StatusBadRequest,
			Body: Body{
				Error: fmt.Sprintf("%v", err),
			},
		})

		return
	}
	created, err := h.Service.CreateTask(ctx, reqTask, action)
	if err != nil {
		h.log.Errorf("Ошибка: %w", err)
		if errors.Is(err, models.ErrActionNotFound) {

			h.SendResponse(w, APIResponse{
				StatusCode: http.StatusBadRequest,
				Body: Body{
					Error: "Action не существует. Проверьте введенные данные",
				},
			})

			return
		} else if errors.Is(err, models.ErrTaskAlreadyExists) {
			h.SendResponse(w, APIResponse{
				StatusCode: http.StatusConflict,
				Body: Body{
					Error: "Задача уже существует",
				},
			})

			return
		} else {
			h.SendResponse(w, APIResponse{
				StatusCode: http.StatusInternalServerError,
				Body: Body{
					Error: "Ошибка при создании задачи",
				},
			})

			return
		}
	}
	h.SendResponse(w, APIResponse{
		StatusCode: http.StatusAccepted,
		Body: newTaskRespDTO{
			Id:          created.Id,
			Uuid:        created.Uuid,
			CreatedAt:   created.CreatedAt,
			ScheduledAt: created.ScheduledAt,
		},
	})
}

func parseTask(r io.ReadCloser) (models.Task, string, error) {
	reqTask, err := getTaskFromBody(r)
	if err != nil {
		return models.Task{}, "", err
	}
	err = validateNewTask(reqTask)
	if err != nil {
		return models.Task{}, "", err
	}

	action := reqTask.Action

	var dateFrom *time.Time
	var dateTo *time.Time
	if reqTask.Params.DateFrom != "" {
		dateFrom, err = parseDate(reqTask.Params.DateFrom)
		if err != nil {
			return models.Task{}, action, fmt.Errorf("Некорректный формат dateFrom. Ожидается YYYY-MM-DD")
		}
	}
	if reqTask.Params.DateTo != "" {
		dateTo, err = parseDate(reqTask.Params.DateTo)
		if err != nil {
			return models.Task{}, action, fmt.Errorf("Некорректный формат dateTo. Ожидается YYYY-MM-DD")
		}
	}
	task := models.Task{
		Uuid: reqTask.Uuid,
		//ActionId: models.Action(reqTask.Action),
		CharCode: &reqTask.Params.CcyCode,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	return task, action, nil
}

func getTaskFromBody(r io.ReadCloser) (newTaskDTO, error) {
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r)
	if err != nil {
		return newTaskDTO{}, fmt.Errorf("Не удалось прочитать тело запроса: %v", err)
	}
	var tsk newTaskDTO
	err = json.Unmarshal(buf.Bytes(), &tsk)
	if err != nil {
		return newTaskDTO{}, fmt.Errorf("Не удалось распарсить тело запроса: %v", err)
	}
	return tsk, nil
}

func validateNewTask(task newTaskDTO) error {
	if task.Action == "" {
		return fmt.Errorf("Поле action обязательно")
	}
	if (task.Params.DateFrom != "" && task.Params.DateTo == "") || (task.Params.DateFrom == "" && task.Params.DateTo != "") {
		return fmt.Errorf("Некорректный период")
	}

	return nil
}

func parseDate(date string) (*time.Time, error) {
	var dt *time.Time
	if date != "" {

		parsed, err := time.Parse(models.DateFormat, date) //2006-01-11
		if err != nil {
			return nil, err
		}
		dt = &parsed
	}
	return dt, nil
}
