package requestplanner

import (
	"testing"
)

func Test_formatParamValue(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		format  string
		want    string
		wantErr bool
	}{
		{"date dd/MM/yyyy", "2024-01-15", "dd/MM/yyyy", "15/01/2024", false},
		{"no format passthrough", "some-value", "", "some-value", false},
		{"empty format", "x", "", "x", false},
		{"invalid date", "not-a-date", "dd/MM/yyyy", "", true},
		{"empty raw with date format", "", "dd/MM/yyyy", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatParamValue(tt.raw, tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("formatParamValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("formatParamValue() = %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_formatParamValue_invalid_date_returns_error(t *testing.T) {
	_, err := formatParamValue("bad", "dd/MM/yyyy")
	if err == nil {
		t.Fatal("expected error for invalid date")
	}
}
