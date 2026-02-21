package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
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
	WHERE li.load_date=cast(getdate()as date) 
)
SELECT
    load_date,
    client_code,
    ccy,
    position_code,
    firm_code,
   -- settle_code,
    firm_name,
    balance
FROM cte
WHERE settle_code = settle_max and balance<>0
ORDER BY load_date, client_code, ccy, position_code, firm_code;
`
)

func (r *Repository) GetMoneyLimits(ctx context.Context) ([]models.MoneyLimit, error) {
	var result []models.MoneyLimit

	r.logger.Debug("получение текущих позиций по деньгам")

	rows, err := r.db.QueryContext(ctx, getMoneyLimits)

	if errors.Is(err, sql.ErrNoRows) {
		r.logger.Error("текущие позиции по деньгам не найдены")
		return nil, models.ErrMLNotFound
	} else if err != nil {
		r.logger.Error("текущие позиции по деньгам не найдены", zap.Error(err))

		return nil, models.ErrMLRetrieving
	}
	defer rows.Close()

	for rows.Next() {
		row := models.MoneyLimit{}
		err = rows.Scan(&row.LoadDate, &row.ClientCode, &row.Currency, &row.PositionCode, &row.FirmCode, &row.FirmName, &row.Balance)
		if err != nil {
			r.logger.Error("ошибка при получении текущих позиций по деньгам", zap.Error(err))

			return nil, models.ErrMLRetrieving
		}
		result = append(result, row)
	}
	if rows.Err() != nil {
		return nil, models.ErrMLRetrieving
	}
	r.logger.Debug("результаты получения позиций по деньгам", zap.Int("", len(result)))

	if len(result) == 0 {
		r.logger.Debug("позиции по деньгам не найдены")
		return nil, models.ErrMLNotFound

	}
	return result, nil
}
