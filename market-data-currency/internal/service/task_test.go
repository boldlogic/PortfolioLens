package service

import (
	"context"
	"strings"
	"testing"
)

func Test_handleResponse_unknown_action_returns_error(t *testing.T) {
	svc := newTestService(&fakeSchedulerRepo{})
	ctx := context.Background()

	err := svc.handleResponse(ctx, "UNKNOWN_ACTION", []byte{}, 1, nil)
	if err == nil {
		t.Fatal("handleResponse(UNKNOWN_ACTION) ожидалась ошибка")
	}
	if !strings.Contains(err.Error(), "неизвестный код действия") {
		t.Errorf("ожидалось сообщение 'неизвестный код действия', получено: %v", err)
	}
	if !strings.Contains(err.Error(), "UNKNOWN_ACTION") {
		t.Errorf("ожидалось упоминание кода в ошибке, получено: %v", err)
	}
}

