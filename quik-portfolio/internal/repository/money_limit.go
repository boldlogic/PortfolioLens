package repository

import (
	"context"
	"time"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/shutdown"
	qmodels "github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
	"go.uber.org/zap"
)

const (
	getMoneyLimits = `
;WITH cte AS (
    SELECT
        li.load_date,
        li.client_code,
        li.ccy,
        li.position_code,
        li.firm_code,
        li.settle_code,
        li.firm_name,
        li.balance,
        settle_max = MAX(li.settle_code) OVER (
            PARTITION BY li.load_date, li.client_code, li.ccy, li.position_code, li.firm_code
        )
    FROM quik.money_limits li
	WHERE li.load_date=cast(@p1 as date) 
)
SELECT
    load_date,
    client_code,
    ccy,
    position_code,
    firm_code,
    firm_name,
    balance
FROM cte
WHERE settle_code = settle_max and balance<>0
ORDER BY load_date, client_code, ccy, position_code, firm_code;
`
)

func (r *Repository) GetMoneyLimits(ctx context.Context, date time.Time) ([]qmodels.MoneyLimit, error) {
	var result []qmodels.MoneyLimit

	r.Logger.Debug("получение текущих позиций по деньгам")

	rows, err := r.Db.QueryContext(ctx, getMoneyLimits, date)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}
		r.Logger.Error("текущие позиции по деньгам не найдены", zap.Error(err))
		return nil, models.ErrRetrievingData
	}
	defer rows.Close()

	for rows.Next() {
		row := qmodels.MoneyLimit{}
		err = rows.Scan(&row.LoadDate, &row.ClientCode, &row.Currency, &row.PositionCode, &row.FirmCode, &row.FirmName, &row.Balance)
		if err != nil {
			if shutdown.IsExceeded(err) {
				return nil, err
			}
			r.Logger.Error("ошибка при получении текущих позиций по деньгам", zap.Error(err))
			return nil, models.ErrRetrievingData
		}
		result = append(result, row)
	}
	if rows.Err() != nil {
		return nil, models.ErrRetrievingData
	}
	r.Logger.Debug("результаты получения позиций по деньгам", zap.Int("count", len(result)))

	if len(result) == 0 {
		r.Logger.Warn("позиции по деньгам не найдены")
		return nil, models.ErrNotFound
	}
	return result, nil
}
