package service

import (
	"github.com/boldlogic/cbr-market-data-worker/internal/client"
	"github.com/boldlogic/cbr-market-data-worker/internal/service/request_catalog"
	"github.com/boldlogic/cbr-market-data-worker/internal/storage"
	"github.com/sirupsen/logrus"
)

type Service struct {
	client   *client.Client
	Provider *request_catalog.Provider
	Storage  *storage.Storage
	log      logrus.FieldLogger
}

func NewService(cl *client.Client, registry *request_catalog.Provider, storage *storage.Storage, log logrus.FieldLogger) *Service {

	return &Service{
		client:   cl,
		Provider: registry,
		Storage:  storage,
		log:      log,
	}
}
