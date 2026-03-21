package service

import (
	"errors"
	"testing"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/utils"
)

func TestCheckLimitDate_после_сегодня_ErrBusinessValidation(t *testing.T) {
	today := utils.Today()
	min := today.AddDate(-1, 0, 0)
	future := today.AddDate(0, 0, 1)
	err := checkLimitDate(future, min)
	if !errors.Is(err, models.ErrBusinessValidation) {
		t.Fatalf("ожидали ErrBusinessValidation: %v", err)
	}
}

func TestMinRollForwardDate_nil_совпадает_с_Today(t *testing.T) {
	got := minRollForwardDate(nil)
	want := utils.Today()
	if utils.DateToYYYYMMDD(got) != utils.DateToYYYYMMDD(want) {
		t.Fatalf("ожидали сегодняшнюю дату: got=%v want=%v", got, want)
	}
}
