package logger

import "testing"

func TestConfig_Validate(t *testing.T) {
	t.Run("Level=info: Validate без ошибок", func(t *testing.T) {
		errs := (Config{Level: "info"}).Validate()
		if len(errs) != 0 {
			t.Fatalf("ожидали 0 ошибок, получили %v", errs)
		}
	})

	t.Run("Level не парсится zapcore: ровно одна ошибка Validate", func(t *testing.T) {
		errs := (Config{Level: "не-уровень"}).Validate()
		if len(errs) != 1 {
			t.Fatalf("ожидали 1 ошибку, получили %v", errs)
		}
	})
}
