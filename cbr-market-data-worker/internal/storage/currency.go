package storage

import (
	"fmt"

	"github.com/boldlogic/PortfolioLens/pkg/models"
)

func (st *Storage) saveCurrency(row *models.Currency) error {
	if row.ISOCode <= 0 {
		return fmt.Errorf("ISOCode должен быть больше 0. CbCode: %v, ISOCharCode: %v", row.CbCode, row.ISOCharCode)
	}

	result := st.db.Where(models.Currency{ISOCode: row.ISOCode}).Assign(models.Currency{
		ISOCharCode: row.ISOCharCode,
		CbCode:      row.CbCode,
		Name:        row.Name,
		LatName:     row.LatName,
		Nominal:     row.Nominal,
		ParentCode:  row.ParentCode,
		ISOCode:     row.ISOCode,
	}).FirstOrCreate(row)
	if result.Error != nil {
		return fmt.Errorf("Currency. %v", result.Error)
	}

	return nil
}

func (st *Storage) SaveCurrencies(rows []models.Currency) []error {
	var errs []error
	for i := range rows {
		err := st.saveCurrency(&rows[i])
		if err != nil {
			errs = append(errs, err)
		}

	}
	return errs

}

func (st *Storage) GetCurrencies() ([]models.Currency, error) {
	var result []models.Currency
	err := st.db.Table("currencies c").Select("c.iso_code, c.iso_char_code,c.name, c.lat_name").Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (st *Storage) GetCurrency(charCode string) (models.Currency, error) {
	var result models.Currency

	err := st.db.Where("iso_char_code=?", charCode).First(&result).Error

	if err != nil {

		return models.Currency{}, err
	}
	return result, nil
}
