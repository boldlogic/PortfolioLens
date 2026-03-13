package cbr

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/xmlconv"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"golang.org/x/net/html/charset"
)

type valute struct {
	ID         string          `xml:"ID,attr"`
	ISONumCode string          `xml:"NumCode"`
	CharCode   string          `xml:"CharCode"`
	Nominal    int             `xml:"Nominal"`
	Name       string          `xml:"Name"`
	Value      xmlconv.RuFloat `xml:"Value"`
	VunitRate  xmlconv.RuFloat `xml:"VunitRate"`
}

type valCurs struct {
	Date   string   `xml:"Date,attr"`
	Name   string   `xml:"name,attr"`
	Valute []valute `xml:"Valute"`
}

func (p *Parser) ParseFxRatesXML(bdy []byte) ([]models.FxRate, error) {
	p.logger.Debug("разбор XML курсов на дату", zap.Int("body_len", len(bdy)))
	decoder := xml.NewDecoder(bytes.NewReader(bdy))
	decoder.CharsetReader = charset.NewReaderLabel
	var v valCurs
	if err := decoder.Decode(&v); err != nil {
		p.logger.Error("ошибка разбора XML курсов на дату", zap.Error(err))
		return []models.FxRate{}, fmt.Errorf("Не удалось получить курсы валют: %w", err)
	}
	p.logger.Debug("в ответе курсов на дату записей", zap.String("Date", v.Date), zap.Int("Valute", len(v.Valute)))

	rateDate, err := time.Parse(RateDateFormat, v.Date)
	if err != nil {
		p.logger.Error("ошибка при разборе даты курсов на дату",
			zap.String("Date", v.Date),
			zap.String("name", v.Name),
			zap.Error(err))
		return []models.FxRate{}, fmt.Errorf("не удалось определить дату курсов валют: %w", err)
	}

	rates := p.valCursToFXRates(v.Valute, rateDate)

	if len(rates) == 0 {
		p.logger.Warn("курсы на дату пусты", zap.String("Date", v.Date))
	}
	return rates, nil
}

func (p *Parser) valCursToFXRates(vals []valute, rateDate time.Time) []models.FxRate {
	out := make([]models.FxRate, 0, len(vals))
	skipped := 0

	for _, item := range vals {
		ccy, ok := p.valuteToFXRate(item, rateDate)
		if !ok {
			skipped++
			continue
		}
		out = append(out, ccy)
	}
	p.logger.Info("курсы на дату",
		zap.Time("Date", rateDate),
		zap.Int("получено", len(vals)),
		zap.Int("распознано", len(out)),
		zap.Int("пропущено", skipped))
	if skipped > 0 {
		p.logger.Warn("часть записей курсов на дату пропущена", zap.Int("пропущено", skipped), zap.Int("получено", len(vals)))
	}
	return out

}

func (p *Parser) valuteToFXRate(val valute, rateDate time.Time) (models.FxRate, bool) {
	var out models.FxRate
	ok := false
	isoCode := 0

	if strings.TrimSpace(val.ISONumCode) == "" {
		p.logger.Warn("пропуск записи курса на дату: пустой NumCode",
			zap.String("ID", val.ID),
			zap.String("NumCode", val.ISONumCode),
			zap.String("CharCode", val.CharCode),
			zap.String("Name", val.Name),
			zap.Int("Nominal", val.Nominal),
		)
		return models.FxRate{}, false
	}

	parsed, err := strconv.Atoi(val.ISONumCode)
	if err != nil {
		p.logger.Warn("пропуск записи курса на дату: невалидный NumCode",
			zap.String("ID", val.ID),
			zap.String("NumCode", val.ISONumCode),
			zap.String("CharCode", val.CharCode),
			zap.String("Name", val.Name),
			zap.Int("Nominal", val.Nominal),
			zap.Error(err),
		)
		return models.FxRate{}, false
	}

	isoCode = parsed
	if isoCode <= 0 {
		p.logger.Warn("пропуск записи курса на дату: NumCode должен быть положительным числом",
			zap.String("ID", val.ID),
			zap.String("NumCode", val.ISONumCode),
			zap.String("CharCode", val.CharCode),
			zap.String("Name", val.Name),
			zap.Int("Nominal", val.Nominal),
		)
		return models.FxRate{}, false
	}

	out.Date = rateDate
	out.BaseISOCode = isoCode
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
