package repository

import (
	"context"
	"database/sql"
	"strings"

	"github.com/boldlogic/PortfolioLens/market-data-currency/internal/apperrors"
	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/shutdown"
	"go.uber.org/zap"
)

const (
	selectNewCurrenciesFromCurrentQuotes = `
		WITH c AS (
			SELECT DISTINCT q.currency
			FROM quik.current_quotes q
			UNION
			SELECT DISTINCT q.base_currency
			FROM quik.current_quotes q
			UNION
			SELECT DISTINCT q.counter_currency
			FROM quik.current_quotes q
			UNION
			SELECT DISTINCT q.quote_currency
			FROM quik.current_quotes q
		),
		cur AS (
			SELECT DISTINCT
				currency = CASE WHEN c.currency IN ('SUR', 'RUR', 'RUB') THEN 'RUB' ELSE c.currency END
			FROM c
		)
		SELECT
			cur.currency,
			currency_name = COALESCE(q.full_name, q.short_name)
		FROM cur
		LEFT JOIN quik.current_quotes q ON q.ticker = cur.currency AND q.class_code = 'CROSSRATE'
		WHERE cur.currency IS NOT NULL
			AND LEN(cur.currency) <= 3
			AND NOT EXISTS (
				SELECT 1
				FROM dbo.currencies c
				WHERE c.iso_char_code = CASE WHEN cur.currency IN ('SUR', 'RUR', 'RUB') THEN 'RUB' ELSE cur.currency END
			)`

	setEmptyNamesFromQuik = `
		WITH NAMES AS (
			SELECT
				iso_char_code = RTRIM(COALESCE(c.iso_char_code, norm.ticker)),
				currency_name = MAX(COALESCE(q.full_name, q.short_name))
			FROM (
				SELECT DISTINCT
					ticker = CASE WHEN q.ticker IN ('SUR', 'RUR', 'RUB') THEN 'RUB' ELSE q.ticker END,
					full_name = q.full_name,
					short_name = q.short_name
				FROM quik.current_quotes q
				WHERE q.class_code = 'CROSSRATE'
					AND LEN(q.ticker) <= 3
			) q
			CROSS APPLY (SELECT ticker = q.ticker) norm
			LEFT JOIN dbo.external_codes ec ON ec.ext_system_id = 2 AND ec.ext_code_type_id = 1 AND ec.ext_code = norm.ticker
			LEFT JOIN dbo.currencies c ON c.iso_code = ec.internal_id
			GROUP BY COALESCE(c.iso_char_code, norm.ticker)
		)
		UPDATE c
		SET
			c.currency_name = n.currency_name,
			updated_at = GETDATE(),
			ext_system_id = 2
		FROM dbo.currencies c
		INNER JOIN NAMES n ON c.iso_char_code = n.iso_char_code
		WHERE c.currency_name IS NULL;`
)

type quoteCurrency struct {
	ISOCharCode string
	Name        sql.NullString
}

func (r *Repository) SetEmptyCurrencyNamesFromQuik(ctx context.Context) error {

	_, err := r.Db.ExecContext(ctx, setEmptyNamesFromQuik)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}
		r.Logger.Error("ошибка при обновлении currency_name в currencies", zap.Error(err))
		return apperrors.ErrSavingData
	}

	return nil

}

func (r *Repository) SelectNewCurrenciesFromCurrentQuotes(ctx context.Context) ([]models.Currency, error) {
	var rawQuoteCurrencies []quoteCurrency

	rows, err := r.Db.QueryContext(ctx, selectNewCurrenciesFromCurrentQuotes)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}
		r.Logger.Error("ошибка при получении новых валют из current_quotes", zap.Error(err))

		return nil, apperrors.ErrRetrievingData
	}
	defer rows.Close()

	for rows.Next() {
		var row quoteCurrency
		err = rows.Scan(&row.ISOCharCode,
			&row.Name)
		if err != nil {
			if shutdown.IsExceeded(err) {
				return nil, err
			}
			r.Logger.Error("ошибка при чтении новых валют из current_quotes", zap.Error(err))

			return nil, apperrors.ErrRetrievingData
		}
		rawQuoteCurrencies = append(rawQuoteCurrencies, row)
	}

	if rows.Err() != nil {
		r.Logger.Error("ошибка при чтении новых валют из current_quotes", zap.Error(err))

		return nil, apperrors.ErrRetrievingData
	}

	r.Logger.Info("в current_quotes найдено новых валют", zap.Int("new_currencies_count", len(rawQuoteCurrencies)))

	if len(rawQuoteCurrencies) == 0 {
		r.Logger.Info("новые валюты в current_quotes не найдены")

		return nil, apperrors.ErrNotFound
	}
	res := quotesToCurrencies(rawQuoteCurrencies)

	return res, nil
}

func quotesToCurrencies(qs []quoteCurrency) []models.Currency {
	var out []models.Currency

	for _, q := range qs {
		out = append(out, quoteToCurrency(q))
	}
	return out
}

func quoteToCurrency(q quoteCurrency) models.Currency {
	var res models.Currency

	res.ISOCharCode = strings.TrimSpace(q.ISOCharCode)
	if q.Name.Valid {
		name := strings.TrimSpace(q.Name.String)
		res.Name = &name
	}

	return res

}
