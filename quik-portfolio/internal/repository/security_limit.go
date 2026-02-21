package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
	"go.uber.org/zap"
)

const (
	getSecurityLimits = `
;WITH cte AS (
    SELECT
        li.load_date,
        li.client_code,
        li.ticker,
        li.trade_account,
        li.firm_code,
        li.settle_code,
        li.firm_name,
        li.balance,
        li.acquisition_ccy,
        li.isin,
        li.ts,
        settle_max = MAX(li.settle_code) OVER (
            PARTITION BY li.load_date, li.client_code, li.ticker, li.trade_account, li.firm_code
        )
    FROM quik.security_limits li
	where li.load_date=cast(getdate()as date) 
)
SELECT
    load_date,
    client_code,
    ticker,
    trade_account,
    firm_code,
    --settle_code,
    firm_name,
    balance,
    acquisition_ccy,
    isin
   -- ts
FROM cte
WHERE settle_code = settle_max and balance<>0
ORDER BY load_date, client_code, ticker, trade_account, firm_code;

`
)

func (r *Repository) GetSecurityLimits(ctx context.Context) ([]models.SecurityLimit, error) {
	var result []models.SecurityLimit

	r.logger.Debug("получение позиций по бумагам")

	rows, err := r.db.QueryContext(ctx, getSecurityLimits)
	r.logger.Debug("", zap.Error(err))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug("позиции по бумагам не найдены")
			return nil, models.ErrSLNotFound
		}
		r.logger.Error("ошибка запроса позиций по бумагам", zap.Error(err))
		return nil, models.ErrSLRetrieving
	}
	defer rows.Close()

	for rows.Next() {
		row := models.SecurityLimit{}
		err = rows.Scan(
			&row.LoadDate,
			&row.ClientCode,
			&row.Ticker,
			&row.TradeAccount,
			&row.FirmCode,
			&row.FirmName,
			&row.Balance,
			&row.AcquisitionCcy,
			&row.ISIN,
		)
		if err != nil {
			r.logger.Error("ошибка при сканировании позиции по бумагам", zap.Error(err))
			return nil, models.ErrSLRetrieving
		}
		result = append(result, row)
	}
	if rows.Err() != nil {
		return nil, models.ErrSLRetrieving
	}
	r.logger.Debug("результаты получения позиций по бумагам", zap.Int("", len(result)))
	if len(result) == 0 {
		r.logger.Debug("позиции по бумагам не найдены")
		return nil, models.ErrSLNotFound

	}

	return result, nil
}
