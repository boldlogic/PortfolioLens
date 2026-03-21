package converters

import (
	"errors"
	"testing"
)

type sampleYAML struct {
	Name string `yaml:"name"`
}

func TestDecodeYAML(t *testing.T) {
	t.Run("Валидный YAML с полем name, без ошибки, Name заполнено", func(t *testing.T) {
		got, err := DecodeYAML[sampleYAML]([]byte("name: one\n"))
		if err != nil {
			t.Fatalf("DecodeYAML: %v", err)
		}
		if got.Name != "one" {
			t.Fatalf("Name=%q, ожидали one", got.Name)
		}
	})

	t.Run("Некорректный YAML (обрыв списка): ErrWrongYAML", func(t *testing.T) {
		_, err := DecodeYAML[sampleYAML]([]byte("name: [\n"))
		if err == nil {
			t.Fatal("ожидали ошибку")
		}
		if !errors.Is(err, ErrWrongYAML) {
			t.Fatalf("ошибка %v, ожидали ErrWrongYAML", err)
		}
	})

	t.Run("Пустой документ YAML: без ошибки, структура нулевая", func(t *testing.T) {
		got, err := DecodeYAML[sampleYAML]([]byte(""))
		if err != nil {
			t.Fatalf("DecodeYAML: %v", err)
		}
		if got.Name != "" {
			t.Fatalf("Name=%q, ожидали пусто", got.Name)
		}
	})
}
