package utils

import (
	"errors"
	"time"

	"github.com/boldlogic/PortfolioLens/pkg/models"
)

// const (
// 	DateFormat = "2006-01-02"
// )

var (
	ErrWrongIsoDateFormat error = errors.New("некорректный формат date. Ожидается YYYY-MM-DD")
)

func ParseDate(date string) (*time.Time, error) {
	var dt *time.Time
	if date != "" {

		parsed, err := time.Parse(models.ISODateFormat, date) //2006-01-11
		if err != nil {
			return nil, err
		}
		dt = &parsed
	}
	return dt, nil
}

// ISO 8601 (YYYY-MM-DD). Дефолт = time.Now()
func ParseIsoDateWithDefault(date string) (time.Time, error) {
	var dt time.Time

	if date != "" {
		parsed, err := time.Parse(models.ISODateFormat, date) //2006-01-11
		if err != nil {
			return dt, ErrWrongIsoDateFormat
		}

		dt = parsed
		return dt, nil
	}

	return time.Now(), nil
}

// Возвращает сегодняшнюю дату (полночь) в локальной временной зоне.
func Today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

// Обрезает t до даты (полночь) в указанной локали.
func TruncateToDateIn(t time.Time, loc *time.Location) time.Time {
	tIn := t.In(loc)
	return time.Date(tIn.Year(), tIn.Month(), tIn.Day(), 0, 0, 0, 0, loc)
}

// Возвращает дату как int64 YYYYMMDD
func DateToYYYYMMDD(t time.Time) int64 {
	y, m, d := t.Date()
	return int64(y)*10000 + int64(m)*100 + int64(d)
}

func EarliestDate(dates ...*time.Time) *time.Time {
	var result *time.Time
	for _, d := range dates {
		if d == nil {
			continue
		}
		if result == nil || d.Before(*result) {
			result = d
		}
	}
	return result
}
