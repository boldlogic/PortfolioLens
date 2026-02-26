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
	getBoards = `
		SELECT board_id, code, name, trade_point_id, is_traded
		FROM quik.boards
		ORDER BY code`

	getBoardByID = `
		SELECT board_id, code, name, trade_point_id, is_traded
		FROM quik.boards
		WHERE board_id = @p1`

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
				LTRIM(RTRIM(class_code)) AS class_code,
				LTRIM(RTRIM(class_name)) AS class_name
			FROM quik.current_quotes
		) MERGE INTO quik.boards AS tgt USING brd ON tgt.code = brd.class_code 
		WHEN MATCHED
		AND (tgt.name <> brd.class_name) 
			THEN UPDATE
			SET tgt.name = brd.class_name 
		WHEN NOT MATCHED BY TARGET 
			THEN INSERT (code, NAME)
		VALUES
		(brd.class_code, brd.class_name);`

	tagBoardsTradePointId = `
		UPDATE b
		SET b.trade_point_id = tp.point_id
		FROM quik.boards b
		INNER JOIN quik.trade_points tp ON tp.code = CASE
			WHEN b.name LIKE 'SPB OTC:%' THEN 'SPB_OTC'
			WHEN b.name LIKE 'SPB:%' THEN 'SPB'
			WHEN b.name LIKE 'FORTS:%' OR b.name LIKE N'%: FORTS' THEN 'FORTS'
			WHEN b.name LIKE N'МБ%' AND (b.name LIKE N'%OTC%' OR b.name LIKE N'%ОТС%' OR b.name LIKE N'%OТС%') THEN 'MOEX_OTC'
			WHEN b.name LIKE N'МБ%' THEN 'MOEX'
			WHEN b.name LIKE N'БКС%' THEN 'BCS_OTC'
			WHEN b.name LIKE N'ВТБ Брокер%' THEN 'VTB_OTC'
		END`
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

func (r *Repository) TagBoardsTradePointId(ctx context.Context) error {
	r.logger.Debug("разметка бордов по торговым площадкам")
	_, err := r.db.ExecContext(ctx, tagBoardsTradePointId)
	if err != nil {
		if IsExceeded(err) {
			return err
		}
		r.logger.Error("ошибка разметки бордов", zap.Error(err))
		return apperrors.ErrSavingData
	}
	return nil
}

func (r *Repository) GetBoards(ctx context.Context) ([]models.Board, error) {
	rows, err := r.db.QueryContext(ctx, getBoards)
	if err != nil {
		if IsExceeded(err) {
			return nil, err
		}
		r.logger.Error("ошибка получения бордов", zap.Error(err))
		return nil, apperrors.ErrRetrievingData
	}
	defer rows.Close()

	var result []models.Board
	for rows.Next() {
		var row models.Board
		var tradePointID sql.NullInt32
		err = rows.Scan(&row.Id, &row.Code, &row.Name, &tradePointID, &row.IsTraded)
		if err != nil {
			if IsExceeded(err) {
				return nil, err
			}
			r.logger.Error("ошибка сканирования борда", zap.Error(err))
			return nil, apperrors.ErrRetrievingData
		}
		setBoardTradePointId(&row, tradePointID)
		result = append(result, row)
	}
	if rows.Err() != nil {
		return nil, apperrors.ErrRetrievingData
	}
	return result, nil
}

func (r *Repository) GetBoardByID(ctx context.Context, id uint8) (models.Board, error) {
	var row models.Board
	var tradePointID sql.NullInt32
	err := r.db.QueryRowContext(ctx, getBoardByID, id).Scan(&row.Id, &row.Code, &row.Name, &tradePointID, &row.IsTraded)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Board{}, apperrors.ErrNotFound
		}
		if IsExceeded(err) {
			return models.Board{}, err
		}
		r.logger.Error("ошибка получения борда", zap.Uint8("id", id), zap.Error(err))
		return models.Board{}, apperrors.ErrRetrievingData
	}
	setBoardTradePointId(&row, tradePointID)
	return row, nil
}

func setBoardTradePointId(row *models.Board, n sql.NullInt32) {
	if n.Valid {
		id := uint8(n.Int32)
		row.TradePointId = &id
	}
}
