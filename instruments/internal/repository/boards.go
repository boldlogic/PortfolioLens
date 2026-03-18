package repository

import (
	"context"
	"database/sql"
	"errors"

	md "github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/models/quik"
	"github.com/boldlogic/PortfolioLens/pkg/shutdown"

	"go.uber.org/zap"
)

const (
	getBoards = `
		SELECT board_id, code, name, trade_point_id, is_traded
		FROM quik.boards
		ORDER BY code`

	getBoardsWithTradePoint = `
		SELECT
			b.board_id,
			b.code,
			b.name,
			b.trade_point_id,
			b.is_traded,
			p.code,
			p.name
		FROM
			quik.boards b
			left join quik.trade_points p on p.point_id=b.trade_point_id
		ORDER BY
			b.board_id`

	getBoardByID = `
		SELECT board_id, code, name, trade_point_id, is_traded
		FROM quik.boards
		WHERE board_id = @p1`

	getBoardByIDWithTradePoint = `
		SELECT
			b.board_id,
			b.code,
			b.name,
			b.trade_point_id,
			b.is_traded,
			p.code,
			p.name
		FROM
			quik.boards b
			LEFT JOIN quik.trade_points p ON p.point_id = b.trade_point_id
		WHERE
			b.board_id = @p1`

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
			WHEN b.name LIKE 'FORTS:%' OR b.name LIKE N'%: FORTS' THEN 'MOEX'
			WHEN b.name LIKE N'МБ Деривативы:%' THEN 'MOEX'
			WHEN b.name LIKE N'МБ%' AND (b.name LIKE N'%OTC%' OR b.name LIKE N'%ОТС%' OR b.name LIKE N'%OТС%') THEN 'MOEX_OTC'
			WHEN b.name LIKE N'МБ%' THEN 'MOEX'
			WHEN b.name LIKE N'БКС%' THEN 'BCS_OTC'
			WHEN b.name LIKE N'ВТБ Брокер%' THEN 'VTB_OTC'
		END`
)

func (r *Repository) InsBoard(ctx context.Context, code string, name string) (quik.Board, error) {
	res := quik.Board{}
	row := r.Db.QueryRowContext(ctx, insBoard, code, name)
	err := row.Scan(&res.Id, &res.Code, &res.Name)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return quik.Board{}, err
		}

		r.Logger.Error("ошибка сохранения кода класса", zap.String("code", code), zap.Error(err))
		return quik.Board{}, md.ErrSavingData
	}

	r.Logger.Debug("кода класса успешно сохранен", zap.String("code", code), zap.Uint8("board_id", res.Id))
	return res, nil
}

func (r *Repository) SyncBoardsFromQuotes(ctx context.Context) error {
	_, err := r.Db.ExecContext(ctx, mergeBoardsFromQuotes)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}

		r.Logger.Error("ошибка сохранения кодов классов из котировок", zap.Error(err))
		return md.ErrSavingData
	}

	r.Logger.Debug("коды классов успешно сохранены из котировок")
	return nil
}

func (r *Repository) TagBoardsTradePointId(ctx context.Context) error {
	_, err := r.Db.ExecContext(ctx, tagBoardsTradePointId)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}

		r.Logger.Error("ошибка разметки кодов классов", zap.Error(err))
		return md.ErrSavingData
	}

	r.Logger.Debug("разметка кодов классов по торговым площадкам завершена успешно")
	return nil
}

func (r *Repository) GetBoards(ctx context.Context) ([]quik.Board, error) {
	rows, err := r.Db.QueryContext(ctx, getBoards)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}

		r.Logger.Error("ошибка получения кодов классов", zap.Error(err))
		return nil, md.ErrRetrievingData
	}
	defer rows.Close()

	var result []quik.Board
	for rows.Next() {
		var row quik.Board
		var tradePointID sql.NullInt32
		err = rows.Scan(&row.Id, &row.Code, &row.Name, &tradePointID, &row.IsTraded)
		if err != nil {
			if shutdown.IsExceeded(err) {
				return nil, err
			}

			r.Logger.Error("ошибка чтения кода класса", zap.Error(err))
			return nil, md.ErrRetrievingData
		}

		setBoardTradePointId(&row, tradePointID)
		result = append(result, row)
	}
	if rows.Err() != nil {
		return nil, md.ErrRetrievingData
	}
	return result, nil
}

func (r *Repository) GetBoardsWithTradePoint(ctx context.Context) ([]quik.Board, error) {
	rows, err := r.Db.QueryContext(ctx, getBoardsWithTradePoint)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}
		r.Logger.Error("ошибка получения бордов", zap.Error(err))
		return nil, md.ErrRetrievingData
	}
	defer rows.Close()

	r.Logger.Debug("получение борда")
	var result []quik.Board
	for rows.Next() {
		var row quik.Board
		var tradePointID sql.NullInt32
		var tradePointCode sql.NullString
		var tradePointName sql.NullString

		err = rows.Scan(&row.Id, &row.Code, &row.Name, &tradePointID, &row.IsTraded, &tradePointCode, &tradePointName)

		if err != nil {
			if shutdown.IsExceeded(err) {
				return nil, err
			}
			r.Logger.Error("ошибка сканирования борда", zap.Error(err))
			return nil, md.ErrRetrievingData
		}
		r.Logger.Debug("получение борда", zap.Any("tradePointID", tradePointID), zap.Any("tradePointCode", tradePointCode), zap.Any("tradePointName", tradePointName))
		setBoardTradePoint(&row, tradePointID, tradePointCode, tradePointName)
		result = append(result, row)
	}
	if rows.Err() != nil {
		return nil, md.ErrRetrievingData
	}
	return result, nil
}

func (r *Repository) GetBoardByID(ctx context.Context, id uint8) (quik.Board, error) {
	var row quik.Board
	var tradePointID sql.NullInt32
	err := r.Db.QueryRowContext(ctx, getBoardByID, id).Scan(&row.Id, &row.Code, &row.Name, &tradePointID, &row.IsTraded)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return quik.Board{}, md.ErrNotFound
		}
		if shutdown.IsExceeded(err) {
			return quik.Board{}, err
		}
		r.Logger.Error("ошибка получения борда", zap.Uint8("id", id), zap.Error(err))
		return quik.Board{}, md.ErrRetrievingData
	}
	setBoardTradePointId(&row, tradePointID)
	return row, nil
}

func (r *Repository) GetBoardByIDWithTradePoint(ctx context.Context, id uint8) (quik.Board, error) {
	var row quik.Board
	var tradePointID sql.NullInt32
	var tradePointCode sql.NullString
	var tradePointName sql.NullString
	err := r.Db.QueryRowContext(ctx, getBoardByIDWithTradePoint, id).Scan(
		&row.Id, &row.Code, &row.Name, &tradePointID, &row.IsTraded,
		&tradePointCode, &tradePointName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return quik.Board{}, md.ErrNotFound
		}
		if shutdown.IsExceeded(err) {
			return quik.Board{}, err
		}
		r.Logger.Error("ошибка получения борда", zap.Uint8("id", id), zap.Error(err))
		return quik.Board{}, md.ErrRetrievingData
	}
	setBoardTradePoint(&row, tradePointID, tradePointCode, tradePointName)
	return row, nil
}

func setBoardTradePointId(row *quik.Board, n sql.NullInt32) {
	if n.Valid {
		id := uint8(n.Int32)
		row.TradePointId = &id
	}
}

func setBoardTradePoint(row *quik.Board, i sql.NullInt32, c sql.NullString, n sql.NullString) {
	if !i.Valid {
		return
	}

	id := uint8(i.Int32)
	row.TradePointId = &id

	row.TradePoint = &md.TradePoint{}
	row.TradePoint.Id = id

	if c.Valid {
		row.TradePoint.Code = c.String
	}
	if n.Valid {
		row.TradePoint.Name = n.String
	}
}
