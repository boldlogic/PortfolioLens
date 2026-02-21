package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
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

	fillInstrumentTypes = `
		WITH
		src AS (
			select distinct(instrument_type) from quik.current_quotes
		) MERGE INTO quik.instrument_types AS tgt USING src ON tgt.title = src.instrument_type 

		WHEN NOT MATCHED BY TARGET THEN INSERT (title)
		VALUES
		(src.instrument_type);
	`
)

func (r *Repository) ActualizeInstrumentTypes(ctx context.Context) error {
	r.logger.Debug("сохранение типов инструментов")
	_, err := r.db.ExecContext(ctx, fillInstrumentTypes)

	if err != nil {
		r.logger.Error("ошибка сохранения типов инструментов", zap.Error(err))
		return models.ErrInstrumentTypesMerging
	}

	return nil
}

func (r *Repository) InsInstrumentType(ctx context.Context, title string) (models.InstrumentType, error) {
	res := models.InstrumentType{}
	r.logger.Debug("сохранение типа инструмента", zap.String("title", title))
	row := r.db.QueryRowContext(ctx, insInstrumentType, title)
	err := row.Scan(&res.Id, &res.Title)

	if err != nil {
		r.logger.Error("ошибка сохранения типа инструмента", zap.String("title", title), zap.Error(err))
		return models.InstrumentType{}, models.ErrInstrumentTypeCreating
	}

	return res, nil
}

func (r *Repository) GetInstrumentTypeId(ctx context.Context, title string) (models.InstrumentType, error) {
	//to-do добавить подсчет времени
	res := models.InstrumentType{}
	r.logger.Debug("сохранение типа инструмента", zap.String("title", title))
	row := r.db.QueryRowContext(ctx, selectInstrumentTypeId, title)

	err := row.Scan(&res.Id, &res.Title)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Error("типа инструмента не найден", zap.String("title", title))
			return models.InstrumentType{}, models.ErrInstrumentTypeNotFound
		}

		r.logger.Error("ошибка получения типа инструмента", zap.String("title", title), zap.Error(err))
		return models.InstrumentType{}, err
	}

	return res, nil
}
