package cbr

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/xmlconv"
	"golang.org/x/net/html/charset"
)

type ValCursDynamic struct {
	ID         string                `xml:"ID,attr"`
	DateRange1 string                `xml:"DateRange1,attr"`
	DateRange2 string                `xml:"DateRange2,attr"`
	Name       string                `xml:"name,attr"`
	Record     []ValCursDynRecordXML `xml:"Record"`
}

type ValCursDynRecordXML struct {
	Date      string          `xml:"Date,attr"`
	ID        string          `xml:"Id,attr"`
	Nominal   int             `xml:"Nominal"`
	Value     xmlconv.RuFloat `xml:"Value"`
	VunitRate xmlconv.RuFloat `xml:"VunitRate"`
}

func ParseFxRateDynamicXML(bdy []byte, base int) ([]models.FxRate, error) {
	decoder := xml.NewDecoder(bytes.NewReader(bdy))
	decoder.CharsetReader = charset.NewReaderLabel
	var valCurs ValCursDynamic
	if err := decoder.Decode(&valCurs); err != nil {
		return []models.FxRate{}, fmt.Errorf("Не удалось получить курсы валют: %w", err)
	}
	rates := make([]models.FxRate, 0, len(valCurs.Record))
	for _, item := range valCurs.Record {
		date, err := time.Parse("02.01.2006", item.Date)
		if err != nil {
			return []models.FxRate{}, fmt.Errorf("не удалось определить дату курсов валют")
		}
		var basePerQuoteUnit float64
		if item.VunitRate != 0 {
			basePerQuoteUnit = float64(1 / item.VunitRate)
		}
		rates = append(rates, models.FxRate{
			Date:             date,
			QuoteISOCode:     643,
			BaseISOCode:      base,
			Nominal:          item.Nominal,
			QuoteForNominal:  float64(item.Value),
			QuotePerUnit:     float64(item.VunitRate),
			BasePerQuoteUnit: basePerQuoteUnit,
		})
	}

	return rates, nil
}
