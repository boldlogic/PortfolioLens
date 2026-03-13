package service

import (
	"context"
	"errors"
	"testing"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"go.uber.org/zap"
)

type fakeCurrencyRepo struct {
	selectCurrenciesFunc func(ctx context.Context) ([]models.Currency, error)
	selectCurrencyFunc   func(ctx context.Context, charCode string) (models.Currency, error)
}

func (f *fakeCurrencyRepo) SelectCurrencies(ctx context.Context) ([]models.Currency, error) {
	if f.selectCurrenciesFunc != nil {
		return f.selectCurrenciesFunc(ctx)
	}
	return nil, nil
}

func (f *fakeCurrencyRepo) SelectCurrency(ctx context.Context, charCode string) (models.Currency, error) {
	if f.selectCurrencyFunc != nil {
		return f.selectCurrencyFunc(ctx, charCode)
	}
	return models.Currency{}, nil
}

func (f *fakeCurrencyRepo) MergeCurrencies(ctx context.Context, currencies []models.Currency) error {
	return nil
}

func (f *fakeCurrencyRepo) MergeExternalCodes(ctx context.Context, codes []models.ExternalCode) error {
	return nil
}

func (f *fakeCurrencyRepo) SelectCountCurrencies(ctx context.Context) (int, error) {
	return 0, nil
}

func (f *fakeCurrencyRepo) SetEmptyCurrencyNamesFromQuik(ctx context.Context) error {
	return nil
}

func (f *fakeCurrencyRepo) MergeFxCBRRates(ctx context.Context, rates []models.FxRate) error {
	return nil
}

func (f *fakeCurrencyRepo) MergeFxCBRRatesQuik(ctx context.Context) error {
	return nil
}

func newTestServiceWithCurrencyRepo(currencyRepo *fakeCurrencyRepo) *Service {
	return NewService(nil, nil, currencyRepo, &fakeSchedulerRepo{}, zap.NewNop())
}

func Test_GetCurrencies_delegates_to_repo(t *testing.T) {
	called := false
	want := []models.Currency{{ISOCharCode: "USD", ISOCode: 840}}
	fakeRepo := &fakeCurrencyRepo{
		selectCurrenciesFunc: func(ctx context.Context) ([]models.Currency, error) {
			called = true
			return want, nil
		},
	}
	svc := newTestServiceWithCurrencyRepo(fakeRepo)
	ctx := context.Background()

	got, err := svc.GetCurrencies(ctx)
	if err != nil {
		t.Fatalf("GetCurrencies: неожиданная ошибка: %v", err)
	}
	if !called {
		t.Error("SelectCurrencies не был вызван")
	}
	if len(got) != 1 || got[0].ISOCharCode != "USD" {
		t.Errorf("GetCurrencies: получено %v, ожидалось %v", got, want)
	}
}

func Test_GetCurrency_invalid_code_returns_ErrBusinessValidation(t *testing.T) {
	selectCalled := false
	fakeRepo := &fakeCurrencyRepo{
		selectCurrencyFunc: func(ctx context.Context, charCode string) (models.Currency, error) {
			selectCalled = true
			return models.Currency{}, nil
		},
	}
	svc := newTestServiceWithCurrencyRepo(fakeRepo)
	ctx := context.Background()

	_, _, err := svc.GetCurrency(ctx, "ZZZ")
	if err != models.ErrBusinessValidation {
		t.Errorf("GetCurrency(ZZZ): получена ошибка %v, ожидалась ErrBusinessValidation", err)
	}
	if selectCalled {
		t.Error("SelectCurrency не должен вызываться при невалидном коде")
	}
}

func Test_GetCurrency_empty_code_returns_ErrBusinessValidation(t *testing.T) {
	svc := newTestServiceWithCurrencyRepo(&fakeCurrencyRepo{})
	ctx := context.Background()

	_, _, err := svc.GetCurrency(ctx, "")
	if err != models.ErrBusinessValidation {
		t.Errorf("GetCurrency(''): получена ошибка %v, ожидалась ErrBusinessValidation", err)
	}
}

func Test_GetCurrency_valid_code_calls_select(t *testing.T) {
	calledWith := ""
	fakeRepo := &fakeCurrencyRepo{
		selectCurrencyFunc: func(ctx context.Context, charCode string) (models.Currency, error) {
			calledWith = charCode
			return models.Currency{ISOCharCode: charCode, ISOCode: 840}, nil
		},
	}
	svc := newTestServiceWithCurrencyRepo(fakeRepo)
	ctx := context.Background()

	ccy, detail, err := svc.GetCurrency(ctx, "USD")
	if err != nil {
		t.Fatalf("GetCurrency(USD): неожиданная ошибка: %v", err)
	}
	if detail != "" {
		t.Errorf("GetCurrency(USD): получено detail=%q, ожидалась пустая строка", detail)
	}
	if ccy.ISOCharCode != "USD" {
		t.Errorf("GetCurrency(USD): получена валюта %q, ожидалось USD", ccy.ISOCharCode)
	}
	if calledWith != "USD" {
		t.Errorf("SelectCurrency вызван с %q, ожидалось USD", calledWith)
	}
}

func Test_GetCurrency_normalizes_to_uppercase(t *testing.T) {
	calledWith := ""
	fakeRepo := &fakeCurrencyRepo{
		selectCurrencyFunc: func(ctx context.Context, charCode string) (models.Currency, error) {
			calledWith = charCode
			return models.Currency{ISOCharCode: charCode}, nil
		},
	}
	svc := newTestServiceWithCurrencyRepo(fakeRepo)
	ctx := context.Background()

	_, _, err := svc.GetCurrency(ctx, "  usd  ")
	if err != nil {
		t.Fatalf("GetCurrency('  usd  '): неожиданная ошибка: %v", err)
	}
	if calledWith != "USD" {
		t.Errorf("SelectCurrency вызван с %q, ожидалось USD (в верхнем регистре)", calledWith)
	}
}

func Test_GetCurrency_repo_error_propagates(t *testing.T) {
	repoErr := errors.New("ошибка БД")
	fakeRepo := &fakeCurrencyRepo{
		selectCurrencyFunc: func(ctx context.Context, charCode string) (models.Currency, error) {
			return models.Currency{}, repoErr
		},
	}
	svc := newTestServiceWithCurrencyRepo(fakeRepo)
	ctx := context.Background()

	_, detail, err := svc.GetCurrency(ctx, "USD")
	if err != repoErr {
		t.Errorf("GetCurrency(USD): получена ошибка %v, ожидалась ошибка репозитория", err)
	}
	if detail != "" {
		t.Errorf("GetCurrency при ошибке репозитория: получено detail=%q, ожидалась пустая строка", detail)
	}
}
