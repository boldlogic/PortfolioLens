package cbr

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/JohannesJHN/iso4217"
	"github.com/boldlogic/PortfolioLens/pkg/models"
	"go.uber.org/zap"
	"golang.org/x/net/html/charset"
)

type valItem struct {
	Id          string `xml:"ID,attr"`
	Name        string `xml:"Name"`
	EngName     string `xml:"EngName"`
	Nominal     int    `xml:"Nominal"`
	ParentCode  string `xml:"ParentCode"`
	ISONumCode  int    `xml:"ISO_Num_Code"`
	ISOCharCode string `xml:"ISO_Char_Code"`
}

func cbrCurrencyItemToCurrency(it valItem) (*models.Currency, *models.ExternalCode) {
	var cur models.Currency
	var ext models.ExternalCode

	if it.ISONumCode == 0 {
		return nil, nil
	}
	charCode := strings.TrimSpace(it.ISOCharCode)
	if charCode == "" {
		return nil, nil
	}

	cur.ISOCode = int16(it.ISONumCode)
	cur.ISOCharCode = charCode
	cur.LatName = strings.TrimSpace(it.EngName)

	lib, ok := iso4217.LookupByAlpha3(charCode)
	if ok {
		cur.MinorUnits = int32(lib.MinorUnits)
	}

	name := strings.TrimSpace(it.Name)
	cur.Name = &name

	extSys := models.CBRSystem
	cur.ExtSystemId = &extSys

	ext.IntId = int64(it.ISONumCode)
	ext.Code = strings.TrimSpace(it.Id)
	ext.Type = models.ExCodeTypeCurrency
	ext.ExternalSystemId = extSys

	return &cur, &ext
}

type valuta struct {
	Name string    `xml:"name,attr"`
	Item []valItem `xml:"Item"`
}

func (p *Parser) ParseCurrenciesXML(bdy []byte) ([]models.Currency, []models.ExternalCode, error) {
	p.logger.Debug("разбор XML справочника валют", zap.Int("body_len", len(bdy)))
	decoder := xml.NewDecoder(bytes.NewReader(bdy))
	decoder.CharsetReader = charset.NewReaderLabel
	var v valuta

	if err := decoder.Decode(&v); err != nil {
		p.logger.Error("ошибка разбора XML справочника валют", zap.Error(err))
		return nil, nil, fmt.Errorf("Не удалось получить справочник валют: %w", err)
	}
	p.logger.Debug("в ответе справочника валют записей", zap.Int("Item", len(v.Item)))
	currencies := make([]models.Currency, 0, len(v.Item))
	extCodes := make([]models.ExternalCode, 0, len(v.Item))
	var skipped int

	for _, item := range v.Item {
		ccy, ext := cbrCurrencyItemToCurrency(item)
		if ccy == nil || ext == nil {
			skipped++
			p.logger.Debug("пропуск записи справочника: невалидный ISO_Num_Code или пустой ISO_Char_Code",
				zap.String("ID", item.Id),
				zap.String("Name", item.Name),
				zap.String("EngName", item.EngName),
				zap.Int("ISO_Num_Code", item.ISONumCode),
				zap.String("ISO_Char_Code", item.ISOCharCode),
				zap.Int("Nominal", item.Nominal),
			)
			continue
		}
		currencies = append(currencies, *ccy)
		extCodes = append(extCodes, *ext)
	}
	p.logger.Info("справочник валют: распознано валют",
		zap.Int("count", len(currencies)),
		zap.Int("skipped", skipped))
	if skipped > 0 {
		p.logger.Warn("часть записей справочника пропущена", zap.Int("skipped", skipped), zap.Int("Item", len(v.Item)))
	}
	if len(currencies) == 0 {
		p.logger.Warn("справочник валют пуст")
	}
	return currencies, extCodes, nil
}
