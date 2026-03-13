package dbzap

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/microsoft/go-mssqldb"
	"go.uber.org/zap"
)

type Pool struct {
	Db     *sql.DB
	Logger *zap.Logger
}

func New(ctx context.Context, dsn string, logger *zap.Logger) (*Pool, error) {
	db, err := openDB(ctx, dsn)
	if err != nil {
		logger.Error("не удалось подключиться к БД", zap.Error(err))
		return nil, err
	}
	return &Pool{Db: db, Logger: logger}, nil
}

func openDB(ctx context.Context, dsn string) (*sql.DB, error) {
	conn, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, err
	}
	if err := conn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("потеряно соединение: %w", err)
	}
	return conn, nil
}

// Close закрывает подключение к БД.
func (p *Pool) Close() {
	p.Db.Close()
}
