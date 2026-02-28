package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/apperrors"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
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

func (r *Repository) GetFirmByName(ctx context.Context, name string) (models.Firm, error) {
	var res models.Firm
	r.logger.Debug("получение фирмы по имени", zap.String("name", name))
	row := r.db.QueryRowContext(ctx, selectFirmByName, name)
	err := row.Scan(&res.Id, &res.Code, &res.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug("фирма не найдена", zap.String("name", name))
			return models.Firm{}, apperrors.ErrNotFound
		}
		r.logger.Error("ошибка получения фирмы по имени", zap.String("name", name), zap.Error(err))
		return models.Firm{}, apperrors.ErrRetrievingData
	}
	return res, nil
}

func (r *Repository) InsertFirm(ctx context.Context, code string, name string) (models.Firm, error) {
	res := models.Firm{}
	r.logger.Debug("сохранение фирмы брокера", zap.String("code", code), zap.String("name", name))
	row := r.db.QueryRowContext(ctx, insertFirms, code, name)
	err := row.Scan(&res.Id, &res.Code, &res.Name)

	if err != nil {
		var mssqlErr mssql.Error
		if errors.As(err, &mssqlErr) && (mssqlErr.Number == 2627 || mssqlErr.Number == 2601) {
			r.logger.Warn("фирма с таким кодом уже существует", zap.String("code", code))
			return models.Firm{}, apperrors.ErrConflict
		}
		r.logger.Error("ошибка сохранения фирмы брокера", zap.String("code", code), zap.String("name", name), zap.Error(err))
		return models.Firm{}, apperrors.ErrSavingData
	}
	r.logger.Debug("фирма успешно сохранена", zap.Uint8("id", res.Id), zap.String("code", res.Code), zap.String("name", res.Name))

	return res, nil
}
