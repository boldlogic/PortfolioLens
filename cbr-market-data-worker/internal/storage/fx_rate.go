package storage

import (
	"fmt"
	"time"

	"github.com/boldlogic/PortfolioLens/pkg/models"
)

func (st *Storage) SaveFxRates(rows []models.FxRate) []error {
	var errs []error
	for _, row := range rows {
		err := st.saveFxRate(&row)
		if err != nil {
			errs = append(errs, err)
		}

	}
	return errs

}

func truncateToDateUTC(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func (st *Storage) saveFxRate(row *models.FxRate) error {
	if row.BaseISOCode <= 0 || row.QuoteISOCode <= 0 {
		return fmt.Errorf("Код валюты должен быть больше 0. BaseISOCode: %v, QuoteISOCode: %v", row.BaseISOCode, row.QuoteISOCode)
	}
	if row.Date.IsZero() {
		return fmt.Errorf("Дата не может быть пустой: %v", row.Date)
	}
	row.Date = truncateToDateUTC(row.Date)
	result := st.db.Where(models.FxRate{Date: row.Date, QuoteISOCode: row.QuoteISOCode, BaseISOCode: row.BaseISOCode}).Assign(models.FxRate{
		Date:             row.Date,
		QuoteISOCode:     row.QuoteISOCode,
		BaseISOCode:      row.BaseISOCode,
		Nominal:          row.Nominal,
		QuoteForNominal:  row.QuoteForNominal,
		QuotePerUnit:     row.QuotePerUnit,
		BasePerQuoteUnit: row.BasePerQuoteUnit,
	}).FirstOrCreate(&row)
	if result.Error != nil {
		return fmt.Errorf("Currency. %v", result.Error)
	}

	return nil
}
