package service

import (
	"context"
)

func (s *Service) SaveFxCBRRatesFromQuik(ctx context.Context) error {
	return s.currencyRepo.MergeFxCBRRatesQuik(ctx)
}
