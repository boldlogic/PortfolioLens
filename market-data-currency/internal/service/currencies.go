package service

import (
	"context"

	"github.com/JohannesJHN/iso4217"
	"go.uber.org/zap"
)

// func (c *Service) GetCbrCurrencies(ctx context.Context, bdy []byte) error {
// 	// ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
// 	// defer cancel()
// 	currencies, err := cbr.ParseCurrenciesXML(bdy)
// 	if err != nil {
// 		return err
// 	}
// 	errs := c.CurrencyRepo.SaveCurrencies(currencies)
// 	if len(errs) > 0 {
// 		if err := errors.Join(errs...); err != nil {
// 			return fmt.Errorf("%w", err)
// 		}
// 	}
// 	return nil
// }

func (s *Service) GetNewCurrenciesFromQuik(ctx context.Context) error {

	currencies, err := s.currencyRepo.SelectNewCurrenciesFromCurrentQuotes(ctx)
	if err != nil {
		return err
	}

	var na []string
	var tickerToISO = map[string]string{
		"GLD": "XAU",
		"SLV": "XAG",
		"PLT": "XPT",
		"PLD": "XPD",
	}

	for i, cur := range currencies {
		charCode := cur.ISOCharCode
		if iso, ok := tickerToISO[charCode]; ok {
			charCode = iso
		}
		iso, ok := iso4217.LookupByAlpha3(charCode)
		if !ok {
			s.logger.Warn("не удалось получить код валюты для", zap.String("ISOCharCode", cur.ISOCharCode))
			na = append(na, cur.ISOCharCode)
		}
		currencies[i].ISOCode = int16(iso.Numeric)
		currencies[i].LatName = iso.Name
		currencies[i].MinorUnits = int32(iso.MinorUnits)

	}
	s.logger.Info("получены валюты", zap.Any("", currencies))
	return nil
}
