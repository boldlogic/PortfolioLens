package storage

import (
	"fmt"
	"strconv"

	"github.com/bojanz/currency"
	"github.com/boldlogic/PortfolioLens/pkg/models"
)

func (st *Storage) Migrate() error {
	err := st.db.AutoMigrate(&models.Currency{})
	if err != nil {
		return err
	}
	err = st.db.AutoMigrate(&models.FxRate{})
	if err != nil {
		return err
	}
	err = st.db.AutoMigrate(&models.Action{})
	if err != nil {
		return err
	}
	err = st.db.AutoMigrate(&models.Task{})
	if err != nil {
		return err
	}
	err = st.FillRussianRuble()
	if err != nil {
		return err
	}
	err = st.FillActions()
	if err != nil {
		return err
	}
	return nil

}

func (st *Storage) FillRussianRuble() error {
	rubCode := "RUB"
	numCodeStr, ok := currency.GetNumericCode(rubCode)
	if !ok {
		return fmt.Errorf("Не удалось получить код валюты для RUB из библиотеки")
	}

	numCode, err := strconv.Atoi(numCodeStr)
	if err != nil {
		return fmt.Errorf("не удалось сконвертировать код валюты: %w", err)
	}

	rub := models.Currency{
		ISOCharCode: rubCode,
		CbCode:      "",
		Name:        "Российский рубль",
		LatName:     "Russian Ruble",
		Nominal:     1,
		ParentCode:  "",
		ISOCode:     numCode,
	}
	result := st.db.Where(models.Currency{ISOCode: numCode}).Assign(rub).FirstOrCreate(&rub)
	if result.Error != nil {
		return fmt.Errorf("%w", result.Error)
	}

	return nil
}

func (st *Storage) FillActions() error {
	act1 := models.Action{
		Code: "currency.cb.fetch.currency_list",
		Name: "Получение справочника валют из ЦБ по www.cbr.ru/scripts/XML_valFull.asp",
	}
	result := st.db.Where(models.Action{Code: act1.Code}).Assign(act1).FirstOrCreate(&act1)
	if result.Error != nil {
		return fmt.Errorf("%w", result.Error)
	}
	act2 := models.Action{
		Code: "currency.cb.fetch.rates_today",
		Name: "Получение справочника валют из ЦБ по www.cbr.ru/scripts/XML_daily.asp",
	}
	result = st.db.Where(models.Action{Code: act2.Code}).Assign(act2).FirstOrCreate(&act2)
	if result.Error != nil {
		return fmt.Errorf("%w", result.Error)
	}
	act3 := models.Action{
		Code: "currency.cb.fetch.historical_rates",
		Name: "Получение справочника валют из ЦБ по www.cbr.ru/scripts/XML_dynamic.asp",
	}
	result = st.db.Where(models.Action{Code: act3.Code}).Assign(act3).FirstOrCreate(&act3)
	if result.Error != nil {
		return fmt.Errorf("%w", result.Error)
	}
	return nil
}
