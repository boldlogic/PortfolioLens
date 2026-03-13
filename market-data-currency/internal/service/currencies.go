package service

import (
	"context"

	"github.com/JohannesJHN/iso4217"
	"github.com/boldlogic/PortfolioLens/pkg/models"
	"go.uber.org/zap"
)

func (s *Service) SaveCbrCurrencies(ctx context.Context, bdy []byte) error {
	currencies, extCodes, err := s.cbrParser.ParseCurrenciesXML(bdy)
	if err != nil {
		return err
	}
	err = s.currencyRepo.MergeCurrencies(ctx, currencies)
	if err != nil {
		return err
	}
	err = s.currencyRepo.MergeExternalCodes(ctx, extCodes)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) setEmptyCurrencyNamesFromQuik(ctx context.Context) error {
	return s.currencyRepo.SetEmptyCurrencyNamesFromQuik(ctx)
}

func (s *Service) InitCurrencyDictionary(ctx context.Context) error {
	err := s.getNewCurrenciesFromLib(ctx)
	if err != nil {
		return err
	}
	return s.setEmptyCurrencyNamesFromQuik(ctx)
}

func (s *Service) getNewCurrenciesFromLib(ctx context.Context) error {

	count, err := s.currencyRepo.SelectCountCurrencies(ctx)
	if err != nil {
		return err
	}
	if count != 0 {
		s.logger.Debug("в справочнике уже есть валюты, библиотеку не используем", zap.Int("количество записей", count))
		return nil
	}

	lib := iso4217.AllActive()

	currencies := make([]models.Currency, 0, len(lib))

	for k, v := range lib {
		currencies = append(currencies, models.Currency{
			ISOCode:     int16(v.Numeric),
			ISOCharCode: k,
			LatName:     v.Name,
			MinorUnits:  int32(v.MinorUnits),
		})
	}

	err = s.currencyRepo.MergeCurrencies(ctx, currencies)
	if err != nil {
		s.logger.Error("произошла ошибка при добавлении валют из библиотеки", zap.Error(err))

		return err
	}
	s.logger.Info("справочник валют был пуст. добавлены валюты из библиотеки", zap.Int("количество записей", len(currencies)))

	return nil
}
