package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/boldlogic/PortfolioLens/market-data-currency/internal/apperrors"
	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/shutdown"
	"go.uber.org/zap"
)

const (
	selectCurrency = `
		select 
			iso_code
			,iso_char_code
			,currency_name
			,lat_name
			,minor_units
			,created_at
			,updated_at
		from dbo.currencies
		where iso_char_code=@p1`

	selectCurrencies = `
		select 
			iso_code
			,iso_char_code
			,currency_name
			,lat_name
			,minor_units
			,created_at
			,updated_at
		from dbo.currencies`

	selectNewCurrenciesFromCurrentQuotes = `
			WITH c AS
		(
			SELECT  distinct (q.currency)
			FROM quik.current_quotes q
			UNION
			SELECT  distinct (q.base_currency)
			FROM quik.current_quotes q
			UNION
			SELECT  distinct (q.counter_currency)
			FROM quik.current_quotes q
			UNION
			SELECT  distinct (q.quote_currency)
			FROM quik.current_quotes q
		), cur as(
		SELECT  distinct currency = CASE WHEN c.currency IN ('SUR','RUR','RUB') THEN 'RUB' else c.currency end
		FROM c )
		SELECT  cur.currency
			,currency_name = coalesce(q.full_name,q.short_name)
		FROM cur
		LEFT JOIN quik.current_quotes q
		ON q.ticker = cur.currency AND q.class_code = 'CROSSRATE'
		WHERE cur.currency is not null
		AND len (cur.currency) <= 3
		AND not exists (
		SELECT  1
		FROM dbo.currencies c
		WHERE c.iso_char_code = CASE WHEN cur.currency IN ('SUR', 'RUR', 'RUB') THEN 'RUB' else cur.currency end )`

	countCurrencies = `
		select count(1) from dbo.currencies`
	mergeCurrencies = `
		WITH src AS (
			SELECT iso_code=@p1
						,iso_char_code=@p2
						,currency_name=@p3
						,lat_name=@p4
						,minor_units=@p5
						,ext_system_id=@p6
		)
		MERGE INTO dbo.currencies AS tgt USING src ON tgt.iso_code=src.iso_code
			AND tgt.iso_char_code=src.iso_char_code 
		WHEN MATCHED
			AND (
				tgt.currency_name<>src.currency_name
				OR tgt.lat_name<>src.lat_name
				OR tgt.minor_units<>src.minor_units
			) THEN
				UPDATE
					SET
						tgt.currency_name=src.currency_name,
						tgt.lat_name=src.lat_name,
						tgt.minor_units=src.minor_units,
						tgt.updated_at=getdate () 
		WHEN NOT MATCHED BY TARGET 
			THEN 
				INSERT (
						iso_code,
						iso_char_code,
						currency_name,
						lat_name,
						minor_units,
						updated_at,
						ext_system_id
						)
				VALUES
					(
						src.iso_code,
						src.iso_char_code,
						src.currency_name,
						src.lat_name,
						src.minor_units,
						getdate () ,
						src.ext_system_id
					);
	`
	setEmptyNamesFromQuik = `
		WITH
		NAMES AS (
			SELECT
			iso_char_code = RTRIM(COALESCE(c.iso_char_code, norm.ticker)),
			currency_name = MAX(COALESCE(q.full_name, q.short_name))
			FROM
			(
				SELECT DISTINCT
				ticker = CASE
					WHEN q.ticker IN ('SUR', 'RUR', 'RUB') THEN 'RUB'
					ELSE q.ticker
				END,
				full_name = q.full_name,
				short_name = q.short_name
				FROM
				quik.current_quotes q
				WHERE
				q.class_code = 'CROSSRATE'
				AND LEN (q.ticker) <= 3
			) q CROSS APPLY (
				SELECT
				ticker = q.ticker
			) norm
			LEFT JOIN dbo.external_codes ec ON ec.ext_system_id = 2
			AND ec.ext_code_type_id = 1
			AND ec.ext_code = norm.ticker
			LEFT JOIN dbo.currencies c ON c.iso_code = ec.internal_id
			GROUP BY
			COALESCE(c.iso_char_code, norm.ticker)
		)
		UPDATE c
		SET
		c.currency_name = n.currency_name,
		updated_at = getdate (),
		ext_system_id = 2
		FROM
		dbo.currencies c
		INNER JOIN NAMES n ON c.iso_char_code = n.iso_char_code
		WHERE
		c.currency_name IS NULL;
			`
)

func (r *Repository) SetEmptyCurrencyNamesFromQuik(ctx context.Context) error {

	_, err := r.Db.ExecContext(ctx, setEmptyNamesFromQuik)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}

		return apperrors.ErrSavingData
	}

	return nil

}

func (r *Repository) SelectCountCurrencies(ctx context.Context) (int, error) {
	var res int

	row := r.Db.QueryRowContext(ctx, countCurrencies)
	err := row.Scan(&res)

	if err != nil {
		if shutdown.IsExceeded(err) {
			return 0, err
		}
		return 0, err
	}
	return res, nil
}
func (r *Repository) MergeCurrencies(ctx context.Context, currencies []models.Currency) error {

	if len(currencies) == 0 {
		return nil
	}

	for _, ccy := range currencies {
		_, err := r.Db.ExecContext(ctx, mergeCurrencies,
			ccy.ISOCode,
			ccy.ISOCharCode,
			ccy.Name,
			ccy.LatName,
			ccy.MinorUnits,
			ccy.ExtSystemId,
		)
		if err != nil {
			if shutdown.IsExceeded(err) {
				return err
			}

			r.Logger.Error("ошибка сохранения валюты", zap.Int16("iso_code", ccy.ISOCode), zap.Error(err))
			return apperrors.ErrSavingData
		}
	}

	return nil

}

func (r *Repository) SelectCurrency(ctx context.Context, charCode string) (models.Currency, error) {
	var res models.Currency

	row := r.Db.QueryRowContext(ctx, selectCurrency, charCode)
	err := row.Scan(&res.ISOCode, &res.ISOCharCode, &res.Name, &res.LatName, &res.MinorUnits, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		if shutdown.IsExceeded(err) {
			return models.Currency{}, err
		}
		if errors.Is(err, sql.ErrNoRows) {
			return models.Currency{}, apperrors.ErrNotFound
		}

		return models.Currency{}, err
	}
	return res, nil
}

func (r *Repository) SelectCurrencies(ctx context.Context) ([]models.Currency, error) {
	var res []models.Currency

	rows, err := r.Db.QueryContext(ctx, selectCurrencies)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}

		return nil, apperrors.ErrRetrievingData
	}
	defer rows.Close()

	for rows.Next() {
		var row models.Currency
		err = rows.Scan(&row.ISOCode,
			&row.ISOCharCode,
			&row.Name,
			&row.LatName,
			&row.MinorUnits,
			&row.CreatedAt,
			&row.UpdatedAt)
		if err != nil {
			if shutdown.IsExceeded(err) {
				return nil, err
			}

			return nil, apperrors.ErrRetrievingData
		}
		res = append(res, row)
	}

	if rows.Err() != nil {
		return nil, apperrors.ErrRetrievingData
	}
	return res, nil
}

type quoteCurrency struct {
	ISOCharCode string
	Name        sql.NullString
}

func (r *Repository) SelectNewCurrenciesFromCurrentQuotes(ctx context.Context) ([]models.Currency, error) {
	var rawQuoteCurrencies []quoteCurrency

	rows, err := r.Db.QueryContext(ctx, selectNewCurrenciesFromCurrentQuotes)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}

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

			return nil, apperrors.ErrRetrievingData
		}
		rawQuoteCurrencies = append(rawQuoteCurrencies, row)
	}

	if rows.Err() != nil {
		return nil, apperrors.ErrRetrievingData
	}

	if len(rawQuoteCurrencies) == 0 {
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
