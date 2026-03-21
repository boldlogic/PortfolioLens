package converters

import (
	"errors"
	"testing"
)

type sampleJSON struct {
	Name string `json:"name"`
}

func TestDecodeJSON(t *testing.T) {
	t.Run("Валидный JSON-объект: без ошибки, поле Name заполнено", func(t *testing.T) {
		got, err := DecodeJSON[sampleJSON]([]byte(`{"name":"one"}`))
		if err != nil {
			t.Fatalf("DecodeJSON: %v", err)
		}
		if got.Name != "one" {
			t.Fatalf("Name=%q, ожидали one", got.Name)
		}
	})

	t.Run("Обрыв JSON: ErrWrongJSON", func(t *testing.T) {
		_, err := DecodeJSON[sampleJSON]([]byte(`{"name":`))
		if err == nil {
			t.Fatal("ожидали ошибку")
		}
		if !errors.Is(err, ErrWrongJSON) {
			t.Fatalf("ошибка %v, ожидали ErrWrongJSON", err)
		}
	})

	t.Run("Пустая строка как тело: ErrWrongJSON", func(t *testing.T) {
		_, err := DecodeJSON[sampleJSON]([]byte(""))
		if err == nil {
			t.Fatal("ожидали ошибку")
		}
		if !errors.Is(err, ErrWrongJSON) {
			t.Fatalf("ошибка %v, ожидали ErrWrongJSON", err)
		}
	})
}
