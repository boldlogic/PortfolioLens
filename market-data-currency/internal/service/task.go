package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/models/scheduler"
	"github.com/boldlogic/PortfolioLens/pkg/transport/requestplanner"
	"go.uber.org/zap"
)

func (s *Service) FetchOneNewTask(ctx context.Context) error {
	task, err := s.schedulerRepo.FetchOneNewTask(ctx)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil
		}
		return err
	}
	s.logger.Debug("задача получена", zap.Int64("task_id", task.Id))

	action, err := s.schedulerRepo.SelectAction(ctx, task.ActionId)
	if err != nil {
		msg := fmt.Sprintf("не найдено действие по задаче %d", task.Id)
		s.logger.Error(msg, zap.Error(err))
		_ = s.schedulerRepo.UpdateTaskStatus(ctx, task.Id, scheduler.TaskStatusError, msg)
		return err
	}

	rawPlan, err := s.schedulerRepo.SelectRequestPlan(ctx, task.ActionId)
	if err != nil {
		msg := fmt.Sprintf("не найден план запроса для действия %s (task_id=%d)", action.Code, task.Id)
		s.logger.Error(msg, zap.Error(err))
		_ = s.schedulerRepo.UpdateTaskStatus(ctx, task.Id, scheduler.TaskStatusError, msg)
		return err
	}

	taskParams, err := s.loadTaskParams(ctx, task.Id)
	if err != nil {
		_ = s.schedulerRepo.UpdateTaskStatus(ctx, task.Id, scheduler.TaskStatusError, err.Error())
		return err
	}

	resolvedPlan, err := s.fillPlanForTask(ctx, rawPlan, taskParams, task.ActionId)
	if err != nil {
		s.logger.Error("ошибка резолва параметров задачи", zap.Int64("task_id", task.Id), zap.Error(err))
		_ = s.schedulerRepo.UpdateTaskStatus(ctx, task.Id, scheduler.TaskStatusError, err.Error())
		return err
	}

	req, err := requestplanner.PrepareRequest(ctx, resolvedPlan)
	if err != nil {
		s.logger.Error("ошибка построения запроса", zap.Int64("task_id", task.Id), zap.Error(err))
		_ = s.schedulerRepo.UpdateTaskStatus(ctx, task.Id, scheduler.TaskStatusError, err.Error())
		return err
	}

	statusCode, body, attempts, err := s.client.SendWithRetry(ctx, req, rawPlan.RetryCount)
	if err != nil {
		msg := fmt.Sprintf("ошибка HTTP-запроса для task_id=%d после %d попыток", task.Id, attempts)
		s.logger.Error(msg, zap.Error(err))
		_ = s.schedulerRepo.UpdateTaskStatus(ctx, task.Id, scheduler.TaskStatusError, msg)
		return err
	}
	if statusCode != http.StatusOK {
		msg := fmt.Sprintf("HTTP %d для task_id=%d после %d попыток", statusCode, task.Id, attempts)
		s.logger.Error(msg)
		_ = s.schedulerRepo.UpdateTaskStatus(ctx, task.Id, scheduler.TaskStatusError, msg)
		return errors.New(msg)
	}

	if err = s.handleResponse(ctx, action.Code, body, task.Id, taskParams); err != nil {
		_ = s.schedulerRepo.UpdateTaskStatus(ctx, task.Id, scheduler.TaskStatusError, err.Error())
		return err
	}

	if err = s.schedulerRepo.UpdateTaskStatus(ctx, task.Id, scheduler.TaskStatusCompleted, ""); err != nil {
		s.logger.Error("ошибка обновления статуса задачи на выполненную", zap.Error(err))
		_ = s.schedulerRepo.UpdateTaskStatus(ctx, task.Id, scheduler.TaskStatusError, err.Error())
		return err
	}
	s.logger.Info("задача выполнена", zap.Int64("task_id", task.Id), zap.String("action", action.Code))
	return nil
}

func (s *Service) loadTaskParams(ctx context.Context, taskId int64) (map[string]string, error) {
	params, err := s.schedulerRepo.SelectTaskParams(ctx, taskId)
	if err != nil {
		return nil, fmt.Errorf("параметры задачи %d: %w", taskId, err)
	}
	out := make(map[string]string, len(params))
	for _, p := range params {
		out[p.Code] = p.Value
	}
	return out, nil
}

func (s *Service) handleResponse(ctx context.Context, actionCode string, body []byte, taskId int64, taskParams map[string]string) error {
	handler, ok := s.responseHandlers[actionCode]
	if !ok {
		return fmt.Errorf("неизвестный код действия: %s", actionCode)
	}
	return handler(ctx, body, taskId, taskParams)
}

func (s *Service) handleCbrCurrencyList(ctx context.Context, body []byte, _ int64, _ map[string]string) error {
	return s.SaveCbrCurrencies(ctx, body)
}

func (s *Service) handleCbrRatesToday(ctx context.Context, body []byte, _ int64, _ map[string]string) error {
	rates, err := s.cbrParser.ParseFxRatesXML(body)
	if err != nil {
		return err
	}
	return s.currencyRepo.MergeFxCBRRates(ctx, rates)
}

func (s *Service) handleCbrHistoricalRates(ctx context.Context, body []byte, taskId int64, taskParams map[string]string) error {
	charCode, ok := taskParams["char_code"]
	if !ok || charCode == "" {
		return fmt.Errorf("не найден параметр char_code в задаче %d", taskId)
	}
	ccy, err := s.currencyRepo.SelectCurrency(ctx, charCode)
	if err != nil {
		return fmt.Errorf("не найдена валюта '%s' для задачи %d: %w", charCode, taskId, err)
	}
	rates, err := s.cbrParser.ParseFxRateDynamicXML(body, int(ccy.ISOCode))
	if err != nil {
		return err
	}
	return s.currencyRepo.MergeFxCBRRates(ctx, rates)
}
