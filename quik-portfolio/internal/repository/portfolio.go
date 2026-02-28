package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/apperrors"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
	"go.uber.org/zap"
)

const getPortfolio = `
WITH
    cte AS (
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
            settle_max = MAX(li.settle_code) OVER (
                PARTITION BY li.load_date, li.client_code, li.ticker, li.trade_account, li.firm_code
            )
        FROM quik.security_limits li
        WHERE li.load_date = cast(getdate() as date)
    ),
    cte_filtered AS (
        SELECT load_date, client_code, ticker, trade_account, firm_code, firm_name, balance, acquisition_ccy, isin
        FROM cte
        WHERE settle_code = settle_max AND balance <> 0
    )
SELECT
    c.load_date,
    c.client_code,
    c.ticker,
    c.trade_account,
    c.firm_code,
    c.firm_name,
    c.balance,
    c.acquisition_ccy,
    c.isin,
    a.mv_currency,
    (isnull(a.price_in_ccy, 0) * c.balance) * coalesce(f.quote_per_unit, 1)
        + (isnull(a.accrued_int, 0) * c.balance) * coalesce(f_accr.quote_per_unit, 1),
    a.short_name
FROM cte_filtered c
OUTER APPLY (
    SELECT TOP 1
        price_in_ccy = case when q.instrument_type = 'Облигации'
            then (isnull(q.face_value, 0) / 100.0) * (case when isnull(q.last_price, 0) <> 0 then q.last_price else q.close_price end)
            else (case when isnull(q.last_price, 0) <> 0 then q.last_price else q.close_price end)
        end,
        q.accrued_int,
        q.base_currency,
        q.counter_currency,
        q.short_name,
        mv_currency = case when q.instrument_type = 'Облигации' then q.base_currency else isnull(q.counter_currency, q.base_currency) end,
        accrued_currency = case when q.instrument_type = 'Облигации' then q.counter_currency else null end
    FROM quik.current_quotes q
    WHERE q.ticker = c.ticker
    ORDER BY case when c.acquisition_ccy = q.base_currency and c.acquisition_ccy = q.counter_currency then 0
        when c.acquisition_ccy = q.base_currency then 1
        when c.acquisition_ccy = q.counter_currency then 2
        else 3 end
) a
LEFT JOIN currencies cv
    ON cv.iso_char_code = case when upper(ltrim(rtrim(isnull(a.mv_currency, '')))) in ('SUR', 'RUR', 'RUB') then 'RUB' else upper(ltrim(rtrim(a.mv_currency))) end
LEFT JOIN currencies ca
    ON ca.iso_char_code = case when upper(ltrim(rtrim(isnull(a.accrued_currency, '')))) in ('SUR', 'RUR', 'RUB') then 'RUB' else upper(ltrim(rtrim(a.accrued_currency))) end
LEFT JOIN fx_rates f ON f.[date] = c.load_date AND f.quote_iso_code = 643 AND f.base_iso_code = cv.iso_code
LEFT JOIN fx_rates f_accr ON f_accr.[date] = c.load_date AND f_accr.quote_iso_code = 643 AND f_accr.base_iso_code = ca.iso_code
ORDER BY c.load_date, c.client_code, c.ticker, c.trade_account, c.firm_code
`

func (r *Repository) GetPortfolio(ctx context.Context) ([]models.PortfolioItem, error) {
	var result []models.PortfolioItem
	r.logger.Debug("получение портфеля (позиции + mv_rub)")

	rows, err := r.db.QueryContext(ctx, getPortfolio)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug("портфель не найден")
			return nil, apperrors.ErrNotFound
		}
		r.logger.Error("ошибка запроса портфеля", zap.Error(err))
		return nil, apperrors.ErrRetrievingData
	}
	defer rows.Close()

	for rows.Next() {
		var row models.PortfolioItem
		var mvCurrency, shortName sql.NullString
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
			&mvCurrency,
			&row.MvRub,
			&shortName,
		)
		if err != nil {
			r.logger.Error("ошибка при сканировании строки портфеля", zap.Error(err))
			return nil, apperrors.ErrRetrievingData
		}
		if mvCurrency.Valid {
			row.MvCurrency = &mvCurrency.String
		}
		if shortName.Valid {
			row.ShortName = &shortName.String
		}
		result = append(result, row)
	}
	if rows.Err() != nil {
		return nil, apperrors.ErrRetrievingData
	}
	if len(result) == 0 {
		r.logger.Debug("позиции портфеля не найдены")
		return nil, apperrors.ErrRetrievingData

	}
	return result, nil
}
