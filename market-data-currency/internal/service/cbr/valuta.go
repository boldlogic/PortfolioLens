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

func cbrCurrencyItemToCurrency(cbr valItem) (*models.Currency, *models.ExternalCode) {
	var cur models.Currency
	var ext models.ExternalCode

	if cbr.ISONumCode == 0 {
		return nil, nil
	}
	charCode := strings.TrimSpace(cbr.ISOCharCode)
	if charCode == "" {
		return nil, nil
	}

	cur.ISOCode = int16(cbr.ISONumCode)
	cur.ISOCharCode = charCode
	cur.LatName = strings.TrimSpace(cbr.EngName)

	lib, ok := iso4217.LookupByAlpha3(charCode)
	if ok {
		cur.MinorUnits = int32(lib.MinorUnits)
	}

	name := strings.TrimSpace(cbr.Name)
	cur.Name = &name

	extSys := models.CBRSystem
	cur.ExtSystemId = &extSys

	ext.IntId = int64(cbr.ISONumCode)
	ext.Code = strings.TrimSpace(cbr.Id)
	ext.Type = models.ExCodeTypeCurrency
	ext.ExternalSystemId = extSys

	return &cur, &ext
}

type Valuta struct {
	Name string    `xml:"name,attr"`
	Item []valItem `xml:"Item"`
}

func (i valItem) String() string {
	return fmt.Sprintf("ID: %s, Name: %s, EngName: %s,Nominal: %d, ParentCode: %s,ISO_Num_Code: %d, ISO_Char_Code: %s", i.Id, i.Name, i.EngName, i.Nominal, i.ParentCode, i.ISONumCode, i.ISOCharCode)
}
func (val Valuta) String() string {
	return fmt.Sprintf("Name: %s, Item: %s", val.Name, val.Item)
}

func ParseCurrenciesXML(bdy []byte, logger *zap.Logger) ([]models.Currency, []models.ExternalCode, error) {
	decoder := xml.NewDecoder(bytes.NewReader(bdy))
	decoder.CharsetReader = charset.NewReaderLabel
	var val Valuta

	if err := decoder.Decode(&val); err != nil {
		return nil, nil, fmt.Errorf("Не удалось получить справочник валют: %w", err)
	}
	currencies := make([]models.Currency, 0, len(val.Item))
	extCodes := make([]models.ExternalCode, 0, len(val.Item))

	for _, item := range val.Item {
		ccy, ext := cbrCurrencyItemToCurrency(item)
		if ccy == nil || ext == nil {
			if logger != nil {
				logger.Debug("пропуск записи справочника валют ЦБР: невалидный ISO_Num_Code или пустой ISO_Char_Code",
					zap.String("id", item.Id),
					zap.Int("ISO_Num_Code", item.ISONumCode),
					zap.String("ISO_Char_Code", item.ISOCharCode),
				)
			}
			continue
		}
		currencies = append(currencies, *ccy)
		extCodes = append(extCodes, *ext)
	}
	return currencies, extCodes, nil
}
