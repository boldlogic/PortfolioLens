package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/boldlogic/cbr-market-data-worker/internal/client"
)

func (c *Service) Executor(ctx context.Context, reqType string) error {

	plan, err := c.Provider.GetPlan(reqType)

	if err != nil {
		return err
	}
	req, err := c.client.PrepareRequest(ctx, plan)
	if err != nil {
		return err
	}

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
		c.log.Errorf("Ошибка при получении данных. Кол-во попыток: %v", cnt)
		return fmt.Errorf("Ошибка при получении данных")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Плохой статус ответа")
	}

	if reqType == "CBR_CURRENCIES" {
		err = c.GetCbrCurrencies(ctx, resp.Body)
		if err != nil {
			return err
		}
	} else if reqType == "CBR_CURRENCY_RATES" {
		err = c.GetCurrencyRates(ctx, resp.Body)
		if err != nil {
			return err
		}
	}

	return nil

}
