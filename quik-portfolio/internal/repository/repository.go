package repository

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/microsoft/go-mssqldb" // MS SQL Server driver
	"go.uber.org/zap"
)

type Repository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewRepository(ctx context.Context, dsn string, logger *zap.Logger) (*Repository, error) {
	db, err := initializeDatabase(ctx, dsn)
	if err != nil {
		logger.Error("не удалось подключиться к БД", zap.Error(err))

		return nil, err
	}
	return &Repository{
		db:     db,
		logger: logger,
	}, nil
}

func initializeDatabase(ctx context.Context, dsn string) (*sql.DB, error) {
	conn, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, err
	}
	if err := conn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("потеряно соединение: %w", err)
	}

	return conn, nil
}

func (r *Repository) Close() {
	r.db.Close()
}
