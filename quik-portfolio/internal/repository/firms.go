package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/models/quik"
	"github.com/boldlogic/PortfolioLens/pkg/shutdown"
	mssql "github.com/microsoft/go-mssqldb"
	"go.uber.org/zap"
)

const (
	insertFirms = `
		insert into quik.firms
		(code, name)
		output inserted.*
		values (@p1, @p2)
	`
	selectFirms = `
		SELECT  firm_id
			,code
			,name
		FROM quik.firms
`
	selectFirmByName = `
		SELECT  firm_id
			,code
			,name
		FROM quik.firms
		where name=@p1
`
)

func (r *Repository) GetFirmByName(ctx context.Context, name string) (quik.Firm, error) {
	var res quik.Firm
	r.Logger.Debug("получение фирмы по имени", zap.String("name", name))
	row := r.Db.QueryRowContext(ctx, selectFirmByName, name)
	err := row.Scan(&res.Id, &res.Code, &res.Name)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return quik.Firm{}, err
		}
		if errors.Is(err, sql.ErrNoRows) {
			r.Logger.Debug("фирма не найдена", zap.String("name", name))
			return quik.Firm{}, models.ErrNotFound
		}
		r.Logger.Error("ошибка получения фирмы по имени", zap.String("name", name), zap.Error(err))
		return quik.Firm{}, models.ErrRetrievingData
	}
	return res, nil
}

func (r *Repository) InsertFirm(ctx context.Context, code string, name string) (quik.Firm, error) {
	res := quik.Firm{}
	r.Logger.Debug("сохранение фирмы брокера", zap.String("code", code), zap.String("name", name))
	row := r.Db.QueryRowContext(ctx, insertFirms, code, name)
	err := row.Scan(&res.Id, &res.Code, &res.Name)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return quik.Firm{}, err
		}
		var mssqlErr mssql.Error
		if errors.As(err, &mssqlErr) && (mssqlErr.Number == 2627 || mssqlErr.Number == 2601) {
			r.Logger.Warn("фирма с таким кодом уже существует", zap.String("code", code))
			return quik.Firm{}, models.ErrConflict
		}
		r.Logger.Error("ошибка сохранения фирмы брокера", zap.String("code", code), zap.String("name", name), zap.Error(err))
		return quik.Firm{}, models.ErrSavingData
	}
	r.Logger.Debug("фирма успешно сохранена", zap.Uint8("id", res.Id), zap.String("code", res.Code), zap.String("name", res.Name))
	return res, nil
}
