package storage

import (
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

//func NewClient(cfg config.ClientConfig, log logrus.FieldLogger, storage *storage.Storage) *Client {

func NewStorage(dsn string) (*Storage, error) {
	db, err := initializeDatabase(dsn)
	if err != nil {
		return nil, err
	}
	return &Storage{
		db: db,
	}, nil
}
