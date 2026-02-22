package utils

import (
	"time"
)

const (
	DateFormat = "2006-01-02"
)

func ParseDate(date string) (*time.Time, error) {
	var dt *time.Time
	if date != "" {

		parsed, err := time.Parse(DateFormat, date) //2006-01-11
		if err != nil {
			return nil, err
		}
		dt = &parsed
	}
	return dt, nil
}

func ParseDateDefault(date string) (time.Time, error) {
	var dt time.Time
	if date != "" {

		parsed, err := time.Parse(DateFormat, date) //2006-01-11
		if err != nil {
			return dt, err
		}
		dt = parsed
	}
	return dt, nil
}
