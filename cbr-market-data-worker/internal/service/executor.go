package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/boldlogic/PortfolioLens/cbr-market-data-worker/internal/client"
	"github.com/boldlogic/PortfolioLens/cbr-market-data-worker/internal/service/cbr"
	"github.com/boldlogic/PortfolioLens/pkg/models"
)

func (s *Service) StartWorker(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	s.log.Debug("Worker started")
	for {
		select {
		case <-ticker.C:
			err := s.ExecuteCreated(ctx)

			if err != nil {
				if !errors.Is(err, models.ErrNoNewTasks) {
					s.log.Error("Worker error:", err)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
func (c *Service) ExecuteCreated(ctx context.Context) error {

	c.log.Debug("выбор задания")
	task, err := c.schedulerRepo.FetchTask(models.TaskStatusCreated, models.TaskStatusInProgress)
	if err != nil {
		if errors.Is(err, models.ErrNoNewTasks) {
			c.log.Debugf("%v", models.ErrNoNewTasks)
			return nil
		}
		c.log.Errorf("произошла ошибка при выборе задания %v", err)
		return err
	}
	c.log.Infof("выбрано задание для обработки id: %d, uuid: %s", task.Id, task.Uuid)

	action, err := c.schedulerRepo.GetAction(task.ActionId)
	if err != nil {
		_ = c.schedulerRepo.SetTaskStatusError(task.Id, fmt.Sprintf("%v", err))
		return err
	}
	c.log.Debug("определен тип задания %s для id: %d, uuid: %s", action.Code, task.Id, task.Uuid)

	plan, err := c.Provider.GetPlan(action.Code)
	if err != nil {
		_ = c.schedulerRepo.SetTaskStatusError(task.Id, fmt.Sprintf("%v", err))
		return err
	}
	var ccy models.Currency
	var reqParams = make(map[string]string)
	for _, param := range plan.Params {
		if param.Type == "cbcode" {

			ccy, err = c.CurrencyRepo.GetCurrency(*task.CharCode)

			if err != nil {
				_ = c.schedulerRepo.SetTaskStatusError(task.Id, fmt.Sprintf("%v", err))
				return err
			}
			if ccy.CbCode == "" {
				c.log.Errorf("для задания id: %d не удалось определить обязательный параметр: Код валюты ЦБ", task.Id)
				err = fmt.Errorf("не удалось определить обязательный параметр: Код валюты ЦБ")
				_ = c.schedulerRepo.SetTaskStatusError(task.Id, fmt.Sprintf("%v", err))
				return err
			}
			reqParams[param.Name] = ccy.CbCode
			c.log.Debug("для задания id: %d, uuid: %s определен Код валюты ЦБ: %s по Коду валюты ISO: %s", task.Id, task.Uuid, ccy.CbCode, *task.CharCode)
		}
		if param.Type == "datefrom" {
			if task.DateFrom != nil {
				reqParams[param.Name] = task.DateFrom.Format(cbr.DateFormat)
			} else {
				err = fmt.Errorf("параметр dateFrom обязателен для задания id: %d", task.Id)
				_ = c.schedulerRepo.SetTaskStatusError(task.Id, fmt.Sprintf("%v", err))
				return err
			}
		}
		if param.Type == "dateto" {
			if task.DateTo != nil {
				reqParams[param.Name] = task.DateTo.Format(cbr.DateFormat)
			} else {
				err = fmt.Errorf("параметр dateTo обязателен для задания id: %d", task.Id)
				_ = c.schedulerRepo.SetTaskStatusError(task.Id, fmt.Sprintf("%v", err))
				return err
			}
		}
	}
	req, err := c.client.PrepareRequestWithParams(ctx, plan, reqParams)
	if err != nil {
		_ = c.schedulerRepo.SetTaskStatusError(task.Id, fmt.Sprintf("%v", err))
		return err
	}
	c.log.Debug("для задания id: %d, uuid: %s подготовлен URL %s", task.Id, task.Uuid, req.URL)

	var resp client.Response
	cnt := 0
	for i := 0; i < plan.RetryCount+1; i++ {
		resp, err = c.client.SendRequest(ctx, req)

		if resp.StatusCode == http.StatusOK && err == nil {
			break
		}
		cnt++
	}
	if err != nil {

		err = fmt.Errorf("для задания id: %d, uuid: %s ошибка при получении данных. Кол-во попыток: %d", task.Id, task.Uuid, cnt+1)
		c.log.Errorf("%v", err)
		_ = c.schedulerRepo.SetTaskStatusError(task.Id, fmt.Sprintf("%v", err))
		return err
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("для задания id: %d, uuid: %s запрос завершен c ошибкой. StatusCode: %d. Кол-во попыток: %d", task.Id, task.Uuid, resp.StatusCode, cnt+1)
		c.log.Errorf("%v", err)
		_ = c.schedulerRepo.SetTaskStatusError(task.Id, fmt.Sprintf("%v", err))
		return err
	}
	c.log.Debug("для задания id: %d, uuid: %s запрос завершен успешно. StatusCode: %d. Кол-во попыток: %d", task.Id, task.Uuid, resp.StatusCode, cnt+1)

	if action.Code == "currency.cb.fetch.currency_list" {
		err = c.GetCbrCurrencies(ctx, resp.Body)
		if err != nil {
			c.log.Errorf("%v", err)
			_ = c.schedulerRepo.SetTaskStatusError(task.Id, fmt.Sprintf("%v", err))

			return err
		}
	}
	if action.Code == "currency.cb.fetch.rates_today" {
		err = c.GetCurrencyRates(ctx, resp.Body)
		if err != nil {
			c.log.Errorf("%v", err)
			_ = c.schedulerRepo.SetTaskStatusError(task.Id, fmt.Sprintf("%v", err))

			return err
		}
	}
	if action.Code == "currency.cb.fetch.historical_rates" {
		err = c.GetCurrencyRatesDynamic(ctx, resp.Body, ccy)
		if err != nil {
			c.log.Errorf("", err)
			_ = c.schedulerRepo.SetTaskStatusError(task.Id, fmt.Sprintf("%v", err))

			return err
		}
	}

	err = c.schedulerRepo.SetTaskStatusCompleted(task.Id)
	c.log.Debug("для задания id: %d, uuid: %s установлен статус %d", task.Id, task.Uuid, models.TaskStatusCompleted)
	if err != nil { //&& !errors.Is(err, models.ErrNoNewTasks)
		c.log.Errorf("%v", err)
		_ = c.schedulerRepo.SetTaskStatusError(task.Id, fmt.Sprintf("%v", err))
		return err
	}

	c.log.Info("задание id: %d, uuid: %s выполнено успешно", task.Id, task.Uuid)
	return nil
}
