package service

import (
	"fmt"
	"time"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/utils"
)

func checkLimitDate(date time.Time, minAvailable time.Time) error {
	loc := time.Now().Location()
	today := utils.Today()
	dateTrunc := utils.TruncateToDateIn(date, loc)
	minTrunc := utils.TruncateToDateIn(minAvailable, loc)

	if dateTrunc.Before(minTrunc) || dateTrunc.After(today) {
		return fmt.Errorf("%w: дата должна быть в диапазоне от %s до %s",
			models.ErrBusinessValidation,
			minTrunc.Format(models.ISODateFormat),
			today.Format(models.ISODateFormat),
		)
	}
	return nil
}

func minRollForwardDate(maxDate *time.Time) time.Time {
	if maxDate == nil {
		return utils.Today()
	}
	return *maxDate
}
