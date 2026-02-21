package cbr

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"golang.org/x/net/html/charset"
)

type ValItem struct {
	Id          string `xml:"ID,attr"`
	Name        string `xml:"Name"`
	EngName     string `xml:"EngName"`
	Nominal     int    `xml:"Nominal"`
	ParentCode  string `xml:"ParentCode"`
	ISONumCode  int    `xml:"ISO_Num_Code"`
	ISOCharCode string `xml:"ISO_Char_Code"`
}

type Valuta struct {
	Name string    `xml:"name,attr"`
	Item []ValItem `xml:"Item"`
}

func (i ValItem) String() string {
	return fmt.Sprintf("ID: %s, Name: %s, EngName: %s,Nominal: %d, ParentCode: %s,ISO_Num_Code: %d, ISO_Char_Code: %s", i.Id, i.Name, i.EngName, i.Nominal, i.ParentCode, i.ISONumCode, i.ISOCharCode)
}
func (val Valuta) String() string {
	return fmt.Sprintf("Name: %s, Item: %s", val.Name, val.Item)
}

func ParseCurrenciesXML(bdy []byte) ([]models.Currency, error) {
	decoder := xml.NewDecoder(bytes.NewReader(bdy))
	decoder.CharsetReader = charset.NewReaderLabel
	var val Valuta

	if err := decoder.Decode(&val); err != nil {
		return []models.Currency{}, fmt.Errorf("Не удалось получить справочник валют: %w", err)
	}
	currencies := make([]models.Currency, 0, len(val.Item))
	for _, item := range val.Item {
		currencies = append(currencies, models.Currency{
			CbCode:      item.Id,
			ISOCharCode: strings.TrimSpace(item.ISOCharCode),
			Name:        strings.TrimSpace(item.Name),
			LatName:     (item.EngName),
			Nominal:     item.Nominal,
			ParentCode:  strings.TrimSpace(item.ParentCode),
			ISOCode:     item.ISONumCode,
		})
	}
	return currencies, nil
}
