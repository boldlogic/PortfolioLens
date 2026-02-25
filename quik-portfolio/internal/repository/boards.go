package repository

import (
	"context"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/apperrors"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
	"go.uber.org/zap"
)

const (
	insBoard = `
		INSERT INTO quik.boards
			(code
			,name)
		output inserted.*
		VALUES (@p1, @p2)`

	mergeBoardsFromQuotes = `
		WITH
		brd AS (
			SELECT DISTINCT
			(class_code),
			class_name
			FROM
			quik.current_quotes
		) MERGE INTO quik.boards AS tgt USING brd ON tgt.code = brd.class_code 
		WHEN MATCHED
		AND (tgt.name <> brd.class_name) 
			THEN UPDATE
			SET tgt.name = brd.class_name 
		WHEN NOT MATCHED BY TARGET 
			THEN INSERT (code, NAME)
		VALUES
		(brd.class_code, brd.class_name);`
)

func (r *Repository) InsBoard(ctx context.Context, code string, name string) (models.Board, error) {
	res := models.Board{}
	r.logger.Debug("сохранение кода класса", zap.String("code", code))
	row := r.db.QueryRowContext(ctx, insBoard, code, name)
	err := row.Scan(&res.Id, &res.Code, &res.Name)

	if err != nil {
		r.logger.Error("ошибка сохранения кода класса", zap.String("code", code), zap.Error(err))
		return models.Board{}, apperrors.ErrSavingData
	}

	return res, nil
}

func (r *Repository) SyncBoardsFromQuotes(ctx context.Context) error {
	r.logger.Debug("сохранение кода класса")
	_, err := r.db.ExecContext(ctx, mergeBoardsFromQuotes)

	if err != nil {
		if IsExceeded(err) {
			return err
		}
		r.logger.Error("ошибка сохранения кода класса", zap.Error(err))
		return apperrors.ErrSavingData
	}

	return nil
}
