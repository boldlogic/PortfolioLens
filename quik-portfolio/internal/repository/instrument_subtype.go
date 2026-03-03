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
	insInstrumentSubType = `
	insert into quik.instrument_subtypes
	(type_id, title)
	output inserted.*
	values (@p1, @p2)
	`
	selectInstrumentSubTypeId = `
	SELECT subtype_id
		,type_id
      	,title
  	FROM quik.instrument_subtypes
	WHERE title=@p1 and type_id=@p2
	`

	mergeInstrumentSubTypesFromQuotes = `
	  WITH
		src AS (
						select distinct(q.instrument_subtype),t.type_id from quik.current_quotes q
			join quik.instrument_types t on t.title=q.instrument_type
			where q.instrument_subtype is not null
		) MERGE INTO quik.instrument_subtypes AS tgt USING src ON tgt.title = src.instrument_subtype and src.type_id=tgt.type_id 

		WHEN NOT MATCHED BY TARGET THEN INSERT (title, type_id)
	values (src.instrument_subtype, src.type_id);`
)

func (r *Repository) SyncInstrumentSubTypesFromQuotes(ctx context.Context) error {
	r.Logger.Debug("сохранение подтипов инструментов")
	_, err := r.Db.ExecContext(ctx, mergeInstrumentSubTypesFromQuotes)

	if err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}
		r.Logger.Error("ошибка сохранения подтипов инструментов", zap.Error(err))
		return apperrors.ErrSavingData
	}

	return nil
}

func (r *Repository) GetInstrumentSubTypeId(ctx context.Context, typeId uint8, title string) (quik.InstrumentSubType, error) {
	//to-do добавить подсчет времени
	res := quik.InstrumentSubType{}
	r.Logger.Debug("сохранение подтипа инструмента", zap.String("title", title))
	row := r.Db.QueryRowContext(ctx, selectInstrumentSubTypeId, title, typeId)

	err := row.Scan(&res.SubTypeId, &res.TypeId, &res.Title)

	if err != nil {
		if shutdown.IsExceeded(err) {
			return quik.InstrumentSubType{}, err
		}
		if errors.Is(err, sql.ErrNoRows) {
			r.Logger.Error("подтип инструмента не найден", zap.String("title", title))
			return quik.InstrumentSubType{}, apperrors.ErrNotFound
		}

		r.Logger.Error("ошибка получения подтипа инструмента", zap.String("title", title), zap.Error(err))
		return quik.InstrumentSubType{}, err
	}

	return res, nil
}

func (r *Repository) InsInstrumentSubType(ctx context.Context, typeId uint8, title string) (quik.InstrumentSubType, error) {
	res := quik.InstrumentSubType{}
	r.Logger.Debug("сохранение подтипа инструмента", zap.String("title", title))
	row := r.Db.QueryRowContext(ctx, insInstrumentSubType, typeId, title)
	err := row.Scan(&res.SubTypeId, &res.TypeId, &res.Title)

	if err != nil {
		if shutdown.IsExceeded(err) {
			return quik.InstrumentSubType{}, err
		}
		r.Logger.Error("ошибка сохранения подтипа инструмента", zap.String("title", title), zap.Error(err))
		return quik.InstrumentSubType{}, apperrors.ErrSavingData
	}

	return res, nil
}
