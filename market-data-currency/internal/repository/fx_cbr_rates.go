package repository

import (
	"context"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"go.uber.org/zap"
)

const mergeFxCBRRate = `
	WITH src AS (
		SELECT
			date             = @p1,
			quote_iso_code   = @p2,
			base_iso_code    = @p3,
			rate_quote_per_base   = @p4,
			rate_base_per_quote = @p5,
			ext_system_id = @p6
	)
	MERGE INTO dbo.fx_cbr_rates AS tgt
	USING src ON tgt.date = src.date
		AND tgt.quote_iso_code = src.quote_iso_code
		AND tgt.base_iso_code  = src.base_iso_code
	WHEN MATCHED THEN UPDATE SET
		tgt.rate_quote_per_base      = src.rate_quote_per_base,
		tgt.rate_base_per_quote = src.rate_base_per_quote,
		tgt.updated_at          = SYSDATETIMEOFFSET(),
		tgt.ext_system_id=src.ext_system_id
	WHEN NOT MATCHED BY TARGET THEN INSERT (
		date, quote_iso_code, base_iso_code,
		rate_quote_per_base, rate_base_per_quote,
		created_at, updated_at,ext_system_id
	) VALUES (
		src.date, src.quote_iso_code, src.base_iso_code,
		src.rate_quote_per_base, src.rate_base_per_quote,
		SYSDATETIMEOFFSET(), SYSDATETIMEOFFSET(), src.ext_system_id
	);`

func (r *Repository) MergeFxCBRRates(ctx context.Context, rates []models.FxRate) error {
	if len(rates) == 0 {
		return nil
	}
	for _, rate := range rates {
		_, err := r.Db.ExecContext(ctx, mergeFxCBRRate,
			rate.Date,
			rate.QuoteISOCode,
			rate.BaseISOCode,
			rate.RateQuotePerBase,
			rate.RateBasePerQuote,
			rate.ExtSystemId,
		)
		if err != nil {
			if r.isShutdown(err) {
				return err
			}
			r.Logger.Error("ошибка сохранения курса валюты",
				zap.Time("date", rate.Date),
				zap.Int("base_iso", rate.BaseISOCode),
				zap.Error(err))
			return models.ErrSavingData
		}
	}
	return nil
}
