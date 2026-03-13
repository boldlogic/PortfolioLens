package service

import (
	"context"
	"errors"
	"testing"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/models/requestplan"
	"github.com/boldlogic/PortfolioLens/pkg/models/scheduler"
	"go.uber.org/zap"
)

type fakeSchedulerRepo struct {
	externalCodeFunc func(ctx context.Context, isoCharCode string, extCodeTypeId, actionId uint8) (string, error)
}

func (f *fakeSchedulerRepo) FetchOneNewTask(ctx context.Context) (scheduler.Task, error) {
	return scheduler.Task{}, models.ErrNotFound
}

func (f *fakeSchedulerRepo) SelectAction(ctx context.Context, id uint8) (scheduler.Action, error) {
	return scheduler.Action{}, nil
}

func (f *fakeSchedulerRepo) UpdateTaskStatus(ctx context.Context, id int64, newStatus scheduler.TaskStatusID, errMsg string) error {
	return nil
}

func (f *fakeSchedulerRepo) SelectTaskParams(ctx context.Context, taskId int64) ([]scheduler.TaskParam, error) {
	return nil, nil
}

func (f *fakeSchedulerRepo) SelectExternalCodeByCurrency(ctx context.Context, isoCharCode string, extCodeTypeId, actionId uint8) (string, error) {
	if f.externalCodeFunc != nil {
		return f.externalCodeFunc(ctx, isoCharCode, extCodeTypeId, actionId)
	}
	return "", nil
}

func (f *fakeSchedulerRepo) SelectRequestPlan(ctx context.Context, actionId uint8) (requestplan.RequestPlan, error) {
	return requestplan.RequestPlan{}, nil
}

func newTestService(repo *fakeSchedulerRepo) *Service {
	return NewService(nil, nil, nil, repo, zap.NewNop())
}

func Test_resolveParamValue_static(t *testing.T) {
	svc := newTestService(&fakeSchedulerRepo{})
	ctx := context.Background()

	ep := requestplan.EndpointParam{
		ParamId:      nil,
		DefaultValue: "static-value",
		ExternalName: "X-Custom-Header",
		ParamLocation: "header",
	}
	value, skip, err := svc.resolveParamValue(ctx, ep, nil, 1)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if skip {
		t.Error("ожидался skip=false для статического параметра")
	}
	if value != "static-value" {
		t.Errorf("ожидалось значение %q, получено %q", "static-value", value)
	}
}

func Test_resolveParamValue_direct_required_missing(t *testing.T) {
	svc := newTestService(&fakeSchedulerRepo{})
	ctx := context.Background()

	paramId := 1
	ep := requestplan.EndpointParam{
		ParamId:       &paramId,
		ParamCode:     "date_from",
		ExtCodeTypeId: nil,
		IsRequired:    true,
	}
	_, _, err := svc.resolveParamValue(ctx, ep, map[string]string{}, 1)
	if err == nil {
		t.Fatal("ожидалась ошибка при отсутствии обязательного параметра")
	}
	if want := "отсутствует обязательный параметр 'date_from'"; err.Error() != want {
		t.Errorf("сообщение об ошибке: получено %q", err.Error())
	}
}

func Test_resolveParamValue_direct_required_ok(t *testing.T) {
	svc := newTestService(&fakeSchedulerRepo{})
	ctx := context.Background()

	paramId := 1
	ep := requestplan.EndpointParam{
		ParamId:       &paramId,
		ParamCode:     "date_from",
		ExtCodeTypeId: nil,
		IsRequired:    true,
	}
	value, skip, err := svc.resolveParamValue(ctx, ep, map[string]string{"date_from": "2024-01-15"}, 1)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if skip {
		t.Error("ожидался skip=false")
	}
	if value != "2024-01-15" {
		t.Errorf("ожидалось %q, получено %q", "2024-01-15", value)
	}
}

func Test_resolveParamValue_direct_optional_missing(t *testing.T) {
	svc := newTestService(&fakeSchedulerRepo{})
	ctx := context.Background()

	paramId := 2
	ep := requestplan.EndpointParam{
		ParamId:       &paramId,
		ParamCode:     "date_to",
		ExtCodeTypeId: nil,
		IsRequired:    false,
	}
	value, skip, err := svc.resolveParamValue(ctx, ep, map[string]string{}, 1)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if !skip {
		t.Error("ожидался skip=true для отсутствующего необязательного параметра")
	}
	if value != "" {
		t.Errorf("ожидалось пустое значение, получено %q", value)
	}
}

func Test_resolveParamValue_externalCode_required_missing(t *testing.T) {
	svc := newTestService(&fakeSchedulerRepo{})
	ctx := context.Background()

	paramId := 3
	extId := uint8(1)
	ep := requestplan.EndpointParam{
		ParamId:       &paramId,
		ParamCode:     "char_code",
		ExtCodeTypeId: &extId,
		IsRequired:    true,
	}
	_, _, err := svc.resolveParamValue(ctx, ep, map[string]string{}, 1)
	if err == nil {
		t.Fatal("ожидалась ошибка при отсутствии обязательного параметра с внешним кодом")
	}
}

func Test_resolveParamValue_externalCode_optional_missing(t *testing.T) {
	svc := newTestService(&fakeSchedulerRepo{})
	ctx := context.Background()

	paramId := 3
	extId := uint8(1)
	ep := requestplan.EndpointParam{
		ParamId:       &paramId,
		ParamCode:     "char_code",
		ExtCodeTypeId: &extId,
		IsRequired:    false,
	}
	value, skip, err := svc.resolveParamValue(ctx, ep, map[string]string{}, 1)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if !skip {
		t.Error("ожидался skip=true")
	}
	if value != "" {
		t.Errorf("ожидалось пустое значение, получено %q", value)
	}
}

func Test_resolveParamValue_externalCode_lookup_ok(t *testing.T) {
	extCodeTypeId := uint8(1)
	repo := &fakeSchedulerRepo{
		externalCodeFunc: func(ctx context.Context, isoCharCode string, extCodeTypeId, actionId uint8) (string, error) {
			if isoCharCode == "USD" && extCodeTypeId == 1 {
				return "840", nil
			}
			return "", models.ErrNotFound
		},
	}
	svc := newTestService(repo)
	ctx := context.Background()

	paramId := 3
	ep := requestplan.EndpointParam{
		ParamId:       &paramId,
		ParamCode:     "char_code",
		ExtCodeTypeId: &extCodeTypeId,
		IsRequired:    true,
	}
	value, skip, err := svc.resolveParamValue(ctx, ep, map[string]string{"char_code": "USD"}, 1)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if skip {
		t.Error("ожидался skip=false")
	}
	if value != "840" {
		t.Errorf("ожидалось %q, получено %q", "840", value)
	}
}

func Test_resolveParamValue_externalCode_lookup_fails(t *testing.T) {
	repoErr := errors.New("db error")
	repo := &fakeSchedulerRepo{
		externalCodeFunc: func(ctx context.Context, isoCharCode string, extCodeTypeId, actionId uint8) (string, error) {
			return "", repoErr
		},
	}
	svc := newTestService(repo)
	ctx := context.Background()

	paramId := 3
	extId := uint8(1)
	ep := requestplan.EndpointParam{
		ParamId:       &paramId,
		ParamCode:     "char_code",
		ExtCodeTypeId: &extId,
		IsRequired:    true,
	}
	_, _, err := svc.resolveParamValue(ctx, ep, map[string]string{"char_code": "USD"}, 1)
	if err == nil {
		t.Fatal("ожидалась ошибка при поиске внешнего кода")
	}
	if !errors.Is(err, repoErr) {
		t.Errorf("ожидалась обёртка repoErr в ошибке, получено %v", err)
	}
}

func Test_fillPlanForTask_empty_params(t *testing.T) {
	svc := newTestService(&fakeSchedulerRepo{})
	ctx := context.Background()

	raw := requestplan.RequestPlan{
		Url:    "https://example.com/api",
		Method: "GET",
		Params: []requestplan.EndpointParam{},
	}
	resolved, err := svc.fillPlanForTask(ctx, raw, map[string]string{}, 1)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if resolved.Url != raw.Url || resolved.Method != raw.Method {
		t.Errorf("url и method не скопированы")
	}
	if len(resolved.Params) != 0 {
		t.Errorf("ожидалось 0 параметров, получено %d", len(resolved.Params))
	}
}

func Test_fillPlanForTask_mixed_params(t *testing.T) {
	extId := uint8(1)
	repo := &fakeSchedulerRepo{
		externalCodeFunc: func(ctx context.Context, isoCharCode string, extCodeTypeId, actionId uint8) (string, error) {
			if isoCharCode == "EUR" {
				return "978", nil
			}
			return "", models.ErrNotFound
		},
	}
	svc := newTestService(repo)
	ctx := context.Background()

	paramId := 1
	raw := requestplan.RequestPlan{
		Url:    "https://cbr.ru/script",
		Method: "GET",
		Params: []requestplan.EndpointParam{
			{ParamId: nil, DefaultValue: "text/xml", ExternalName: "Accept", ParamLocation: "header"},
			{ParamId: &paramId, ParamCode: "char_code", ExtCodeTypeId: &extId, IsRequired: true, ExternalName: "date_req1", ParamLocation: "query"},
			{ParamId: &paramId, ParamCode: "date_from", ExtCodeTypeId: nil, IsRequired: false, ExternalName: "date_req1", ParamLocation: "query"},
		},
	}
	taskParams := map[string]string{"char_code": "EUR", "date_from": "2024-01-01"}

	resolved, err := svc.fillPlanForTask(ctx, raw, taskParams, 1)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if len(resolved.Params) != 3 {
		t.Fatalf("ожидалось 3 параметра, получено %d", len(resolved.Params))
	}
	if resolved.Params[0].Value != "text/xml" {
		t.Errorf("статический параметр: получено %q", resolved.Params[0].Value)
	}
	if resolved.Params[1].Value != "978" {
		t.Errorf("параметр внешнего кода: получено %q", resolved.Params[1].Value)
	}
	if resolved.Params[2].Value != "2024-01-01" {
		t.Errorf("прямой параметр: получено %q", resolved.Params[2].Value)
	}
}

func Test_fillPlanForTask_skips_optional_missing(t *testing.T) {
	svc := newTestService(&fakeSchedulerRepo{})
	ctx := context.Background()

	paramId := 1
	raw := requestplan.RequestPlan{
		Url:    "https://example.com",
		Method: "GET",
		Params: []requestplan.EndpointParam{
			{ParamId: &paramId, ParamCode: "optional", IsRequired: false, ExternalName: "opt", ParamLocation: "query"},
			{ParamId: nil, DefaultValue: "fixed", ExternalName: "required", ParamLocation: "header"},
		},
	}

	resolved, err := svc.fillPlanForTask(ctx, raw, map[string]string{}, 1)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if len(resolved.Params) != 1 {
		t.Fatalf("ожидался 1 параметр (необязательный пропущен), получено %d", len(resolved.Params))
	}
	if resolved.Params[0].Value != "fixed" {
		t.Errorf("ожидалось fixed, получено %q", resolved.Params[0].Value)
	}
}

func Test_fillPlanForTask_error_propagates(t *testing.T) {
	svc := newTestService(&fakeSchedulerRepo{})
	ctx := context.Background()

	paramId := 1
	raw := requestplan.RequestPlan{
		Url:    "https://example.com",
		Method: "GET",
		Params: []requestplan.EndpointParam{
			{ParamId: &paramId, ParamCode: "req", IsRequired: true, ExternalName: "req", ParamLocation: "query"},
		},
	}

	_, err := svc.fillPlanForTask(ctx, raw, map[string]string{}, 1)
	if err == nil {
		t.Fatal("ожидалась ошибка при отсутствии обязательного параметра")
	}
}
