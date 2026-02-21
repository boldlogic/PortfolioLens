package cbr

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strconv"
	"time"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/xmlconv"
	"golang.org/x/net/html/charset"
)

type Valute struct {
	ID         string          `xml:"ID,attr"`
	ISONumCode string          `xml:"NumCode"`
	CharCode   string          `xml:"CharCode"`
	Nominal    int             `xml:"Nominal"`
	Name       string          `xml:"Name"`
	Value      xmlconv.RuFloat `xml:"Value"`
	VunitRate  xmlconv.RuFloat `xml:"VunitRate"`
}

type ValCurs struct {
	Date   string   `xml:"Date,attr"`
	Name   string   `xml:"name,attr"`
	Valute []Valute `xml:"Valute"`
}

func ParseFxRatesXML(bdy []byte) ([]models.FxRate, error) {
	decoder := xml.NewDecoder(bytes.NewReader(bdy))
	decoder.CharsetReader = charset.NewReaderLabel
	var valCurs ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		return []models.FxRate{}, fmt.Errorf("Не удалось получить курсы валют: %w", err)
	}

	rateDate, err := time.Parse("02.01.2006", valCurs.Date)
	if err != nil {
		return []models.FxRate{}, fmt.Errorf("не удалось определить дату курсов валют")
	}

	rates := make([]models.FxRate, 0, len(valCurs.Valute))

	for _, item := range valCurs.Valute {

		isoCode := 0
		if item.ISONumCode != "" {
			parsed, err := strconv.Atoi(item.ISONumCode)
			if err == nil {
				isoCode = parsed
			}
		}
		if isoCode <= 0 {
			continue //наверное надо будет потом логировать кривые валюты, но пока не критично
		}
		var basePerQuoteUnit float64
		if item.VunitRate != 0 {
			basePerQuoteUnit = float64(1 / item.VunitRate)
		}
		rates = append(rates, models.FxRate{
			Date:             rateDate,
			QuoteISOCode:     643,
			BaseISOCode:      isoCode,
			Nominal:          item.Nominal,
			QuoteForNominal:  float64(item.Value),
			QuotePerUnit:     float64(item.VunitRate),
			BasePerQuoteUnit: basePerQuoteUnit,
		})
	}
	return rates, nil
}
