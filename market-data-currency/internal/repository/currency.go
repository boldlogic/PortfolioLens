package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/boldlogic/PortfolioLens/market-data-currency/internal/apperrors"
	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/shutdown"
)

// iso_code       SMALLINT          NOT NULL,
//         iso_char_code  NVARCHAR(3)       NOT NULL,
//         currency_name  NVARCHAR(100)     NULL,
//         lat_name       NVARCHAR(100)     NULL,
//         minor_units        INT               NULL,
//         created_at     DATETIMEOFFSET(7) NULL,
//         updated_at     DATETIMEOFFSET(7) NULL,

// func (r *Repository) saveCurrency(row *models.Currency) error {

// 	if row.ISOCode <= 0 {
// 		return fmt.Errorf("ISOCode должен быть больше 0. ISOCharCode: %v", row.ISOCharCode)
// 	}

// 	result := st.db.Where(models.Currency{ISOCode: row.ISOCode}).Assign(models.Currency{
// 		ISOCharCode: row.ISOCharCode,
// 		//CbCode:      row.CbCode,
// 		Name:    row.Name,
// 		LatName: row.LatName,
// 		MinorUnits: row.MinorUnits,
// 		//ParentCode:  row.ParentCode,
// 		ISOCode: row.ISOCode,
// 	}).FirstOrCreate(row)
// 	if result.Error != nil {
// 		return fmt.Errorf("Currency. %v", result.Error)
// 	}

// 	return nil
// }

// func (r *Repository) SaveCurrencies(rows []models.Currency) []error {
// 	var errs []error
// 	for i := range rows {
// 		err := st.saveCurrency(&rows[i])
// 		if err != nil {
// 			errs = append(errs, err)
// 		}

// 	}
// 	return errs

// }

// func (r *Repository) GetCurrencies() ([]models.Currency, error) {
// 	var result []models.Currency
// 	err := st.db.Table("currencies c").Select("c.iso_code, c.iso_char_code,c.name, c.lat_name").Scan(&result).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

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
)

type quoteCurrency struct {
	ISOCharCode string
	Name        sql.NullString
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
		res.Name = name
	}

	return res

}
