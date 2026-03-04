package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/boldlogic/PortfolioLens/pkg/models/quik"
	"github.com/boldlogic/PortfolioLens/pkg/shutdown"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/apperrors"
	"go.uber.org/zap"
)

const (
	insInstrumentType = `
		insert into quik.instrument_types
		(title)
		output inserted.*
		values (@p1)
	`
	selectInstrumentTypeId = `
		SELECT 
			type_id
			,title
		FROM quik.instrument_types
		WHERE title=@p1
	`

	mergeInstrumentTypesFromQuotes = `
		WITH
		src AS (
			select distinct(instrument_type) from quik.current_quotes
		) MERGE INTO quik.instrument_types AS tgt USING src ON tgt.title = src.instrument_type 

		WHEN NOT MATCHED BY TARGET THEN INSERT (title)
		VALUES
		(src.instrument_type);
	`
)

func (r *Repository) SyncInstrumentTypesFromQuotes(ctx context.Context) error {
	r.Logger.Debug("сохранение типов инструментов")

	_, err := r.Db.ExecContext(ctx, mergeInstrumentTypesFromQuotes)

	if err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}

		r.Logger.Error("ошибка сохранения типов инструментов", zap.Error(err))
		return apperrors.ErrSavingData
	}

	return nil
}

func (r *Repository) InsInstrumentType(ctx context.Context, title string) (quik.InstrumentType, error) {
	res := quik.InstrumentType{}
	r.Logger.Debug("сохранение типа инструмента", zap.String("title", title))
	row := r.Db.QueryRowContext(ctx, insInstrumentType, title)
	err := row.Scan(&res.Id, &res.Title)

	if err != nil {
		if shutdown.IsExceeded(err) {
			return quik.InstrumentType{}, err
		}
		r.Logger.Error("ошибка сохранения типа инструмента", zap.String("title", title), zap.Error(err))
		return quik.InstrumentType{}, apperrors.ErrSavingData
	}

	return res, nil
}

func (r *Repository) GetInstrumentTypeId(ctx context.Context, title string) (quik.InstrumentType, error) {
	//to-do добавить подсчет времени
	res := quik.InstrumentType{}
	r.Logger.Debug("сохранение типа инструмента", zap.String("title", title))
	row := r.Db.QueryRowContext(ctx, selectInstrumentTypeId, title)

	err := row.Scan(&res.Id, &res.Title)

	if err != nil {
		if shutdown.IsExceeded(err) {
			return quik.InstrumentType{}, err
		}
		if errors.Is(err, sql.ErrNoRows) {
			r.Logger.Error("типа инструмента не найден", zap.String("title", title))
			return quik.InstrumentType{}, apperrors.ErrNotFound
		}

		r.Logger.Error("ошибка получения типа инструмента", zap.String("title", title), zap.Error(err))
		return quik.InstrumentType{}, err
	}

	return res, nil
}
