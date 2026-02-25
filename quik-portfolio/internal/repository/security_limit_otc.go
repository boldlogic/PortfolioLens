package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/apperrors"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
	mssql "github.com/microsoft/go-mssqldb"
	"go.uber.org/zap"
)

const (
	insertSecurityLimitOtc = `
INSERT INTO quik.security_limits_otc
           (load_date, client_code, ticker, trade_account, settle_code, firm_code, firm_name, balance, acquisition_ccy, isin)
     VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10)
	`
	getSecurityLimitsOtc = `
SELECT load_date, client_code, ticker, trade_account, firm_code, firm_name, balance, acquisition_ccy, isin
FROM quik.security_limits_otc
WHERE load_date = CAST(@p1 AS date) AND balance <> 0
ORDER BY load_date, client_code, ticker, trade_account, firm_code
	`
	getSecurityLimitsOtcMaxDate = `
	select max(load_date) FROM quik.security_limits_otc 
	`
	deleteSecurityLimitsOtcBeforeDate = `
DELETE FROM
    quik.security_limits_otc
WHERE
    load_date<CAST(@p1 AS date)	
	`
	rollSecurityLimitsOtcFromDateToDate = `
INSERT INTO
    quik.security_limits_otc (
        load_date,
        client_code,
        ticker,
        trade_account,
        settle_code,
        firm_code,
        firm_name,
        balance,
        acquisition_ccy,
        isin
    )
SELECT
    CAST(@p1 AS date),
    client_code,
    ticker,
    trade_account,
    settle_code,
    firm_code,
    firm_name,
    balance,
    acquisition_ccy,
    isin
FROM
    quik.security_limits_otc
WHERE
    load_date=CAST(@p2 AS date)
    AND balance<>0
ORDER BY
    load_date,
    client_code,
    ticker,
    trade_account,
    firm_code	
	`
)

func (r *Repository) DeleteSecurityLimitsOtcBeforeDate(ctx context.Context, date time.Time) error {
	_, err := r.db.ExecContext(ctx, deleteSecurityLimitsOtcBeforeDate, date)
	if err != nil {

		if IsExceeded(err) {
			return err
		}

		r.logger.Error("ошибка при удалении лимитов OTC по бумагам",
			zap.Time("LoadDate", date),
			zap.Error(err))
		return apperrors.ErrSavingData
	}

	r.logger.Debug("лимиты OTC по бумаге успешно удалены",
		zap.Time("LoadDate", date))
	return nil
}

func (r *Repository) RollSecurityLimitsOtcFromDateToDate(ctx context.Context, dateFrom time.Time, dateTo time.Time) error {
	_, err := r.db.ExecContext(ctx, rollSecurityLimitsOtcFromDateToDate,
		dateTo, dateFrom)
	if err != nil {

		if IsExceeded(err) {
			return err
		}

		var msErr mssql.Error
		if errors.As(err, &msErr) && (msErr.Number == 2627 || msErr.Number == 2601) {
			r.logger.Warn("лимит OTC по бумаге уже существует",
				zap.Time("LoadDate", dateTo))
			return apperrors.ErrConflict
		}

		r.logger.Error("ошибка при создании лимита OTC по бумаге",
			zap.Time("LoadDate", dateTo),
			zap.Error(err))
		return apperrors.ErrSavingData
	}
	r.logger.Debug("лимит OTC по бумаге успешно сохранен",
		zap.Time("LoadDate", dateTo))
	return nil
}

func (r *Repository) SaveSecurityLimitOtc(ctx context.Context, s models.SecurityLimit) error {
	_, err := r.db.ExecContext(ctx, insertSecurityLimitOtc,
		s.LoadDate,
		s.ClientCode,
		s.Ticker,
		s.TradeAccount,
		s.SettleCode,
		s.FirmCode,
		s.FirmName,
		s.Balance,
		s.AcquisitionCcy,
		s.ISIN)

	if err != nil {
		if IsExceeded(err) {
			return err
		}

		var msErr mssql.Error
		if errors.As(err, &msErr) && (msErr.Number == 2627 || msErr.Number == 2601) {
			r.logger.Warn("лимит OTC по бумаге уже существует",
				zap.Time("LoadDate", s.LoadDate),
				zap.String("ClientCode", s.ClientCode),
				zap.String("Ticker", s.Ticker),
				zap.String("TradeAccount", s.TradeAccount),
				zap.String("FirmCode", s.FirmCode))
			return apperrors.ErrConflict
		}

		r.logger.Error("ошибка при создании лимита OTC по бумаге",
			zap.Time("LoadDate", s.LoadDate),
			zap.String("ClientCode", s.ClientCode),
			zap.String("Ticker", s.Ticker),
			zap.Error(err))
		return apperrors.ErrSavingData
	}
	r.logger.Debug("лимит OTC по бумаге успешно сохранен",
		zap.Time("LoadDate", s.LoadDate),
		zap.String("ClientCode", s.ClientCode),
		zap.String("Ticker", s.Ticker))
	return nil
}

func (r *Repository) GetSecurityLimitsOtcMaxDate(ctx context.Context) (*time.Time, error) {

	var date *time.Time
	r.logger.Debug("получение максимальной даты из quik.security_limits_otc")
	row := r.db.QueryRowContext(ctx, getSecurityLimitsOtcMaxDate)

	err := row.Scan(&date)
	if err != nil {
		if IsExceeded(err) {
			return nil, err
		}

		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug("лимиты не найдены")
			return nil, apperrors.ErrNotFound
		}
		r.logger.Error("ошибка получения максимальной даты из quik.security_limits_otc", zap.Error(err))
		return nil, apperrors.ErrRetrievingData
	}
	return date, nil

}

func (r *Repository) GetSecurityLimitsOtc(ctx context.Context, date time.Time) ([]models.SecurityLimit, error) {
	var result []models.SecurityLimit
	rows, err := r.db.QueryContext(ctx, getSecurityLimitsOtc, date)
	if err != nil {
		if IsExceeded(err) {
			return nil, err
		}

		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug("позиции OTC по бумагам не найдены")
			return nil, apperrors.ErrNotFound
		}
		r.logger.Error("ошибка запроса позиций OTC по бумагам", zap.Error(err))
		return nil, apperrors.ErrRetrievingData
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
			if IsExceeded(err) {
				return nil, err
			}

			r.logger.Error("ошибка при сканировании позиции OTC по бумагам", zap.Error(err))
			return nil, apperrors.ErrRetrievingData
		}
		result = append(result, row)
	}
	if rows.Err() != nil {
		return nil, apperrors.ErrRetrievingData
	}
	if len(result) == 0 {
		r.logger.Debug("позиции OTC по бумагам не найдены")
		return nil, apperrors.ErrNotFound
	}
	return result, nil
}
