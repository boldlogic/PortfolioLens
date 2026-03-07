package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/boldlogic/PortfolioLens/market-data-currency/internal/apperrors"
	"github.com/boldlogic/PortfolioLens/market-data-currency/internal/client"
	"github.com/boldlogic/PortfolioLens/pkg/models/scheduler"
	"go.uber.org/zap"
)

func (s *Service) FetchOneNewTask(ctx context.Context) error {

	task, err := s.schedulerRepo.FetchOneNewTask(ctx)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			s.logger.Debug("нет новых задач")

			return nil
		}
		return err
	}
	s.logger.Debug("", zap.Any("task", task))

	action, err := s.schedulerRepo.SelectAction(ctx, task.ActionId)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			msg := fmt.Sprintf("не найдено действие по задаче %d", task.Id)
			s.logger.Debug(msg)
			_ = s.schedulerRepo.UpdateTaskStatus(ctx, task.Id, scheduler.TaskStatusError, msg)
			return err
		}
		s.logger.Error("ошибка при поиске действия по задаче", zap.Error(err))
		_ = s.schedulerRepo.UpdateTaskStatus(ctx, task.Id, scheduler.TaskStatusError, err.Error())

		return err
	}

	plan, err := s.provider.GetPlan(action.Code)
	if err != nil {
		_ = s.schedulerRepo.UpdateTaskStatus(ctx, task.Id, scheduler.TaskStatusError, err.Error())
		return err
	}

	req, err := s.client.PrepareRequest(ctx, plan)
	if err != nil {
		_ = s.schedulerRepo.UpdateTaskStatus(ctx, task.Id, scheduler.TaskStatusError, err.Error())
		return err
	}

	var resp client.Response
	cnt := 0
	for i := 0; i < plan.RetryCount+1; i++ {
		resp, err = s.client.SendRequest(ctx, req)

		if resp.StatusCode == http.StatusOK && err == nil {
			break
		}
		cnt++
	}
	if err != nil {

		err = fmt.Errorf("для задания id: %d, uuid: %s ошибка при получении данных. Кол-во попыток: %d", task.Id, task.Uuid, cnt+1)
		s.logger.Error("ошибка при выполнении HTTP-запроса по задаче", zap.Error(err))
		_ = s.schedulerRepo.UpdateTaskStatus(ctx, task.Id, scheduler.TaskStatusError, err.Error())
		return err
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("для задания id: %d, uuid: %s запрос завершен c ошибкой. StatusCode: %d. Кол-во попыток: %d", task.Id, task.Uuid, resp.StatusCode, cnt+1)
		s.logger.Error("запрос по задаче завершился с неожиданным HTTP-статусом", zap.Error(err))

		_ = s.schedulerRepo.UpdateTaskStatus(ctx, task.Id, scheduler.TaskStatusError, err.Error())
		return err
	}

	if action.Code == "currency.cb.fetch.currency_list" {
		err = s.GetCbrCurrencies(ctx, resp.Body)
		if err != nil {
			s.logger.Error("ошибка при загрузке справочника валют ЦБР по задаче", zap.Error(err))

			_ = s.schedulerRepo.UpdateTaskStatus(ctx, task.Id, scheduler.TaskStatusError, err.Error())

			return err
		}
	}

	err = s.schedulerRepo.UpdateTaskStatus(ctx, task.Id, scheduler.TaskStatusCompleted, "")
	if err != nil {
		s.logger.Error("ошибка при обновлении статуса задачи на completed", zap.Error(err))

		_ = s.schedulerRepo.UpdateTaskStatus(ctx, task.Id, scheduler.TaskStatusError, err.Error())
		return err
	}

	return nil
}
