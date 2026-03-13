package repository

import (
	"context"

	"github.com/boldlogic/PortfolioLens/pkg/dbzap"
	"go.uber.org/zap"
)

type Repository struct {
	*dbzap.Pool
}

func NewRepository(ctx context.Context, dsn string, logger *zap.Logger) (*Repository, error) {
	pool, err := dbzap.New(ctx, dsn, logger)
	if err != nil {
		logger.Error("ошибка подключения к БД", zap.Error(err))
		return nil, err
	}
	return &Repository{Pool: pool}, nil
}
