package cbr

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/xmlconv"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"golang.org/x/net/html/charset"
)

type valCursDynamic struct {
	ID         string                `xml:"ID,attr"`
	DateRange1 string                `xml:"DateRange1,attr"`
	DateRange2 string                `xml:"DateRange2,attr"`
	Name       string                `xml:"name,attr"`
	Record     []valCursDynRecordXML `xml:"Record"`
}

type valCursDynRecordXML struct {
	Date      string          `xml:"Date,attr"`
	ID        string          `xml:"Id,attr"`
	Nominal   int             `xml:"Nominal"`
	Value     xmlconv.RuFloat `xml:"Value"`
	VunitRate xmlconv.RuFloat `xml:"VunitRate"`
}

func (p *Parser) ParseFxRateDynamicXML(bdy []byte, base int) ([]models.FxRate, error) {
	p.logger.Debug("разбор XML истории курса валюты", zap.Int("body_len", len(bdy)), zap.Int("base_iso", base))
	decoder := xml.NewDecoder(bytes.NewReader(bdy))
	decoder.CharsetReader = charset.NewReaderLabel
	var v valCursDynamic
	if err := decoder.Decode(&v); err != nil {
		p.logger.Error("ошибка разбора XML истории курса валюты", zap.Error(err))
		return []models.FxRate{}, fmt.Errorf("Не удалось получить курсы валют: %w", err)
	}
	p.logger.Debug("в ответе истории курса валюты записей",
		zap.Int("Record", len(v.Record)),
		zap.String("DateRange1", v.DateRange1),
		zap.String("DateRange2", v.DateRange2))

	rates := p.valCursDynamicToFXRates(v.Record, base)

	if len(rates) == 0 {
		p.logger.Warn("история курса валюты пуста", zap.Int("base_iso", base))
	}
	return rates, nil
}

func (p *Parser) valCursDynamicToFXRates(vals []valCursDynRecordXML, base int) []models.FxRate {
	out := make([]models.FxRate, 0, len(vals))
	skipped := 0

	for _, item := range vals {
		ccy, ok := p.valCursDynRecordToFXRate(item, base)
		if !ok {
			skipped++
			continue
		}
		out = append(out, ccy)
	}

	p.logger.Info("история курса валюты",
		zap.Int("base_iso", base),
		zap.Int("получено", len(vals)),
		zap.Int("распознано", len(out)),
		zap.Int("пропущено", skipped))

	if skipped > 0 {
		p.logger.Warn("часть записей курсов на дату пропущена", zap.Int("пропущено", skipped), zap.Int("получено", len(vals)))
	}
	return out
}

func (p *Parser) valCursDynRecordToFXRate(val valCursDynRecordXML, base int) (models.FxRate, bool) {

	var out models.FxRate
	ok := false
	date, err := time.Parse(RateDateFormat, val.Date)
	if err != nil {
		p.logger.Warn("ошибка при разборе даты в записи истории курса валюты",
			zap.String("Date", val.Date),
			zap.String("Id", val.ID),
			zap.Int("Nominal", val.Nominal),
			zap.Float64("Value", float64(val.Value)),
			zap.Int("base_iso", base),
			zap.Error(err))
		return models.FxRate{}, false
	}

	out.Date = date
	out.BaseISOCode = base
	out.QuoteISOCode = 643

	quotePerUnit := decimal.NewFromFloat(float64(val.VunitRate))
	if quotePerUnit.IsZero() && val.Nominal != 0 {
		quotePerUnit = decimal.NewFromFloat(float64(val.Value)).Div(decimal.NewFromInt(int64(val.Nominal)))
	}
	out.RateQuotePerBase = quotePerUnit

	var basePerQuoteUnit decimal.Decimal
	if !quotePerUnit.IsZero() {
		basePerQuoteUnit = decimal.NewFromInt(1).Div(quotePerUnit)
	} else if float64(val.Value) != 0 && val.Nominal != 0 {
		basePerQuoteUnit = decimal.NewFromInt(int64(val.Nominal)).Div(decimal.NewFromFloat(float64(val.Value)))
	}
	out.RateBasePerQuote = basePerQuoteUnit
	extSys := models.CBRSystem
	out.ExtSystemId = &extSys
	ok = true
	return out, ok
}
