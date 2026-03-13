package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"go.uber.org/zap"
)

const (
	selectCurrency = `
		SELECT
			iso_code,
			iso_char_code,
			currency_name,
			lat_name,
			minor_units,
			created_at,
			updated_at
		FROM dbo.currencies
		WHERE iso_char_code = @p1`

	selectCurrencies = `
		SELECT
			iso_code,
			iso_char_code,
			currency_name,
			lat_name,
			minor_units,
			created_at,
			updated_at
		FROM dbo.currencies`

	countCurrencies = `
		SELECT COUNT(1) FROM dbo.currencies`

	mergeCurrencies = `
		WITH src AS (
			SELECT
				iso_code = @p1,
				iso_char_code = @p2,
				currency_name = @p3,
				lat_name = @p4,
				minor_units = @p5,
				ext_system_id = @p6
		)
		MERGE INTO dbo.currencies AS tgt
		USING src ON tgt.iso_code = src.iso_code
		WHEN MATCHED
			AND (
				tgt.currency_name <> src.currency_name
				OR tgt.lat_name <> src.lat_name
				OR tgt.minor_units <> src.minor_units
			)
		THEN UPDATE SET
			tgt.currency_name = src.currency_name,
			tgt.lat_name = src.lat_name,
			tgt.minor_units = src.minor_units,
			tgt.updated_at = SYSDATETIMEOFFSET(),
			tgt.ext_system_id=src.ext_system_id
		WHEN NOT MATCHED BY TARGET
		THEN INSERT (
			iso_code,
			iso_char_code,
			currency_name,
			lat_name,
			minor_units,
			updated_at,
			ext_system_id
		)
		VALUES (
			src.iso_code,
			src.iso_char_code,
			src.currency_name,
			src.lat_name,
			src.minor_units,
			SYSDATETIMEOFFSET(),
			src.ext_system_id
		);`
)

func (r *Repository) SelectCountCurrencies(ctx context.Context) (int, error) {
	var res int

	row := r.Db.QueryRowContext(ctx, countCurrencies)
	err := row.Scan(&res)

	if err != nil {
		if r.isShutdown(err) {
			return 0, err
		}
		r.Logger.Error("ошибка при получении количества валют в справочнике currencies", zap.Error(err))

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
			if r.isShutdown(err) {
				return err
			}
			r.Logger.Error("ошибка сохранения валюты", zap.Int16("iso_code", ccy.ISOCode), zap.Error(err))
			return models.ErrSavingData
		}
	}

	return nil

}

func (r *Repository) SelectCurrency(ctx context.Context, charCode string) (models.Currency, error) {
	var res models.Currency

	row := r.Db.QueryRowContext(ctx, selectCurrency, charCode)
	err := row.Scan(&res.ISOCode, &res.ISOCharCode, &res.Name, &res.LatName, &res.MinorUnits, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		if r.isShutdown(err) {
			return models.Currency{}, err
		}
		if errors.Is(err, sql.ErrNoRows) {
			return models.Currency{}, models.ErrNotFound
		}
		r.Logger.Error("ошибка при получении валюты из справочника currencies", zap.String("iso_char_code", charCode), zap.Error(err))

		return models.Currency{}, models.ErrRetrievingData
	}
	return res, nil
}

func (r *Repository) SelectCurrencies(ctx context.Context) ([]models.Currency, error) {
	var res []models.Currency

	rows, err := r.Db.QueryContext(ctx, selectCurrencies)
	if err != nil {
		if r.isShutdown(err) {
			return nil, err
		}
		r.Logger.Error("ошибка при получении валют из справочника currencies", zap.Error(err))
		return nil, models.ErrRetrievingData
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
			if r.isShutdown(err) {
				return nil, err
			}
			r.Logger.Error("ошибка при чтении валюты из справочника currencies", zap.Error(err))
			return nil, models.ErrRetrievingData
		}
		res = append(res, row)
	}

	if rows.Err() != nil {
		r.Logger.Error("ошибка при получении валют из справочника", zap.Error(rows.Err()))

		return nil, models.ErrRetrievingData
	}
	return res, nil
}
