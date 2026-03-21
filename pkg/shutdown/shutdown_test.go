package shutdown

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func TestIsExceeded(t *testing.T) {
	t.Run("Ошибка nil: IsExceeded возвращает false", func(t *testing.T) {
		if IsExceeded(nil) {
			t.Fatal("ожидали false")
		}
	})

	t.Run("Сама context.Canceled: IsExceeded true", func(t *testing.T) {
		if !IsExceeded(context.Canceled) {
			t.Fatal("ожидали true")
		}
	})

	t.Run("Ошибка обёрнута вокруг context.Canceled: IsExceeded true", func(t *testing.T) {
		err := fmt.Errorf("db: %w", context.Canceled)
		if !IsExceeded(err) {
			t.Fatal("ожидали true")
		}
	})

	t.Run("Сама context.DeadlineExceeded: IsExceeded true", func(t *testing.T) {
		if !IsExceeded(context.DeadlineExceeded) {
			t.Fatal("ожидали true")
		}
	})

	t.Run("Ошибка обёрнута вокруг context.DeadlineExceeded: IsExceeded true", func(t *testing.T) {
		err := fmt.Errorf("timeout: %w", context.DeadlineExceeded)
		if !IsExceeded(err) {
			t.Fatal("ожидали true")
		}
	})

	t.Run("Произвольная ошибка без отмены контекста: IsExceeded false", func(t *testing.T) {
		if IsExceeded(errors.New("boom")) {
			t.Fatal("ожидали false")
		}
	})
}
