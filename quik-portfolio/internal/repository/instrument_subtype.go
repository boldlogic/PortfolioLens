package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/apperrors"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
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
	r.logger.Debug("сохранение подтипов инструментов")
	_, err := r.db.ExecContext(ctx, mergeInstrumentSubTypesFromQuotes)

	if err != nil {
		if IsExceeded(err) {
			return err
		}
		r.logger.Error("ошибка сохранения подтипов инструментов", zap.Error(err))
		return apperrors.ErrSavingData
	}

	return nil
}

func (r *Repository) GetInstrumentSubTypeId(ctx context.Context, typeId int16, title string) (models.InstrumentSubType, error) {
	//to-do добавить подсчет времени
	res := models.InstrumentSubType{}
	r.logger.Debug("сохранение подтипа инструмента", zap.String("title", title))
	row := r.db.QueryRowContext(ctx, selectInstrumentSubTypeId, title, typeId)

	err := row.Scan(&res.SubTypeId, &res.TypeId, &res.Title)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Error("подтип инструмента не найден", zap.String("title", title))
			return models.InstrumentSubType{}, apperrors.ErrNotFound
		}

		r.logger.Error("ошибка получения подтипа инструмента", zap.String("title", title), zap.Error(err))
		return models.InstrumentSubType{}, err
	}

	return res, nil
}

func (r *Repository) InsInstrumentSubType(ctx context.Context, typeId int16, title string) (models.InstrumentSubType, error) {
	res := models.InstrumentSubType{}
	r.logger.Debug("сохранение подтипа инструмента", zap.String("title", title))
	row := r.db.QueryRowContext(ctx, insInstrumentSubType, typeId, title)
	err := row.Scan(&res.SubTypeId, &res.TypeId, &res.Title)

	if err != nil {
		r.logger.Error("ошибка сохранения подтипа инструмента", zap.String("title", title), zap.Error(err))
		return models.InstrumentSubType{}, apperrors.ErrSavingData
	}

	return res, nil
}
