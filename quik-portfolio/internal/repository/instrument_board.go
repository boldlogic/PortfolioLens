package repository

import (
	"context"
	"database/sql"

	"github.com/boldlogic/PortfolioLens/pkg/shutdown"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/apperrors"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
	"go.uber.org/zap"
)

type instrumentBoard struct {
	InstrumentId int
	BoardId      uint8

	TypeId    uint8
	SubTypeId sql.NullInt16

	CurrencyId        sql.NullInt64
	BaseCurrencyId    sql.NullInt64
	QuoteCurrencyId   sql.NullInt64
	CounterCurrencyId sql.NullInt64
	IsPrimary         bool
}

const (
	selectInstrumentBoard = `
		SELECT instrument_id
			,board_id
			,type_id
			,subtype_id
			,currency_id
			,base_currency_id
			,quote_currency_id
			,counter_currency_id
			,is_primary
		FROM quik.instrument_boards
		WHERE instrument_id=@p1 and board_id=@p2
	`
	mergeInstrumentBoard = `
		WITH
		src AS (
			SELECT
			instrument_id = CAST(@p1 AS BIGINT),
			board_id = CAST(@p2 AS TINYINT),
			type_id = CAST(@p3 AS TINYINT),
			subtype_id = CAST(@p4 AS TINYINT),
			currency_id = CAST(@p5 AS BIGINT),
			base_currency_id = CAST(@p6 AS BIGINT),
			quote_currency_id = CAST(@p7 AS BIGINT),
			counter_currency_id = CAST(@p8 AS BIGINT)
		) MERGE INTO quik.instrument_boards AS tgt USING src ON tgt.instrument_id = src.instrument_id
		AND tgt.board_id = src.board_id
		WHEN MATCHED AND (
			tgt.type_id <> src.type_id
			OR ISNULL(tgt.subtype_id, 0) <> ISNULL(src.subtype_id, 0)
			OR ISNULL(tgt.currency_id, -1) <> ISNULL(src.currency_id, -1)
			OR ISNULL(tgt.base_currency_id, -1) <> ISNULL(src.base_currency_id, -1)
			OR ISNULL(tgt.quote_currency_id, -1) <> ISNULL(src.quote_currency_id, -1)
			OR ISNULL(tgt.counter_currency_id, -1) <> ISNULL(src.counter_currency_id, -1)
		) THEN
		UPDATE
		SET
		tgt.type_id = src.type_id,
		tgt.subtype_id = src.subtype_id,
		tgt.currency_id = src.currency_id,
		tgt.base_currency_id = src.base_currency_id,
		tgt.quote_currency_id = src.quote_currency_id,
		tgt.counter_currency_id = src.counter_currency_id WHEN NOT matched BY target THEN INSERT (
			instrument_id,
			board_id,
			type_id,
			subtype_id,
			currency_id,
			base_currency_id,
			quote_currency_id,
			counter_currency_id
		)
		VALUES
		(CAST(@p1 AS BIGINT), CAST(@p2 AS TINYINT), CAST(@p3 AS TINYINT), CAST(@p4 AS TINYINT), CAST(@p5 AS BIGINT), CAST(@p6 AS BIGINT), CAST(@p7 AS BIGINT), CAST(@p8 AS BIGINT));
	`
)

// To-do func (r *Repository) SelectInstrumentBoard(ctx context.Context, id int, board uint8) (models.InstrumentBoard, error) {

// To-do func (r *Repository) SelectInstrumentBoards(ctx context.Context, id int) ([]models.InstrumentBoard, error) {

func (r *Repository) MergeInstrumentBoard(ctx context.Context, ib models.InstrumentBoard) error {
	_, err := r.db.ExecContext(ctx, mergeInstrumentBoard,
		ib.InstrumentId,
		ib.BoardId,
		ib.TypeId,
		ib.SubTypeId,
		ib.CurrencyId,
		ib.BaseCurrencyId,
		ib.QuoteCurrencyId,
		ib.CounterCurrencyId)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}

		r.logger.Error("ошибка сохранения кода класса для инструмента", zap.Int("instrument_id", ib.InstrumentId), zap.Uint8("board_id", ib.BoardId), zap.Error(err))
		return apperrors.ErrSavingData
	}
	return nil

}
