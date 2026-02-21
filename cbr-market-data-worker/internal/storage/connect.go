package storage

import (
	"fmt"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func initializeDatabase(dsn string) (*gorm.DB, error) {

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к БД: %w", err)
	}
	if err := checkDatabaseConnection(db); err != nil {
		return nil, err
	}
	return db, nil
}

func checkDatabaseConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("не удалось подключиться к БД: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("потеряно соединение с БД: %w", err)
	}
	return nil
}
