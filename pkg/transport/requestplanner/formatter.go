package requestplanner

import (
	"fmt"
	"time"
)

func formatParamValue(raw, format string) (string, error) {
	switch format {
	case "dd/MM/yyyy":
		t, err := time.Parse("2006-01-02", raw)
		if err != nil {
			return "", fmt.Errorf("неверная дата '%s', ожидается YYYY-MM-DD: %w", raw, err)
		}
		return t.Format("02/01/2006"), nil
	default:
		return raw, nil
	}
}
