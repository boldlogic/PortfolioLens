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
	insertSecurityLimit = `
INSERT INTO quik.security_limits
           (load_date
           ,client_code
           ,ticker
           ,trade_account
           ,settle_code
           ,firm_code
           ,firm_name
           ,balance
           ,acquisition_ccy
           ,isin)
		  -- output inserted.*
     VALUES (@p1, @p2, @p3, @p4, @p5, @p6,@p7,@p8,@p9, @p10)	
	`
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
	where li.load_date=cast(@p1 as date) 
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

func (r *Repository) SaveSecurityLimit(ctx context.Context, s models.SecurityLimit) error {

	_, err := r.db.ExecContext(ctx, insertSecurityLimit,
		s.LoadDate, s.ClientCode, s.Ticker, s.TradeAccount, s.SettleCode,
		s.FirmCode, s.FirmName, s.Balance, s.AcquisitionCcy, s.ISIN)
	if err != nil {
		var msErr mssql.Error
		if errors.As(err, &msErr) && (msErr.Number == 2627 || msErr.Number == 2601) {
			r.logger.Warn("лимит по бумаге уже существует", zap.Time("LoadDate", s.LoadDate), zap.String("ClientCode", s.ClientCode), zap.String("Ticker", s.Ticker),
				zap.String("TradeAccount", s.TradeAccount), zap.String("SettleCode", s.SettleCode), zap.String("SettleCode", s.SettleCode), zap.String("FirmCode", s.FirmCode))
			return apperrors.ErrConflict
		}
		r.logger.Error("ошибка при создании лимита по бумаге", zap.Time("LoadDate", s.LoadDate), zap.String("ClientCode", s.ClientCode), zap.String("Ticker", s.Ticker),
			zap.String("TradeAccount", s.TradeAccount), zap.String("SettleCode", s.SettleCode), zap.String("SettleCode", s.SettleCode), zap.String("FirmCode", s.FirmCode), zap.Error(err))
		return apperrors.ErrSavingData
	}
	r.logger.Debug("лимит по бумаге успешно сохранен", zap.Time("LoadDate", s.LoadDate), zap.String("ClientCode", s.ClientCode), zap.String("Ticker", s.Ticker),
		zap.String("TradeAccount", s.TradeAccount), zap.String("SettleCode", s.SettleCode), zap.String("SettleCode", s.SettleCode), zap.String("FirmCode", s.FirmCode))
	return nil
}

func (r *Repository) GetSecurityLimits(ctx context.Context, date time.Time) ([]models.SecurityLimit, error) {
	var result []models.SecurityLimit

	rows, err := r.db.QueryContext(ctx, getSecurityLimits, date)
	r.logger.Debug("", zap.Error(err))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug("позиции по бумагам не найдены")
			return nil, apperrors.ErrSLNotFound
		}
		r.logger.Error("ошибка запроса позиций по бумагам", zap.Error(err))
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
			r.logger.Error("ошибка при сканировании позиции по бумагам", zap.Error(err))
			return nil, apperrors.ErrRetrievingData
		}
		result = append(result, row)
	}
	if rows.Err() != nil {
		return nil, apperrors.ErrRetrievingData
	}
	r.logger.Debug("результаты получения позиций по бумагам", zap.Int("", len(result)))
	if len(result) == 0 {
		r.logger.Debug("позиции по бумагам не найдены")
		return nil, apperrors.ErrSLNotFound

	}

	return result, nil
}
