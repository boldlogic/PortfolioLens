package commonconfig

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, name string, content []byte) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, content, 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func Test_readFile(t *testing.T) {
	t.Run("чтение_существующего_файла", func(t *testing.T) {
		p := writeTemp(t, "x.txt", []byte("hello"))
		got, err := readFile(p)
		if err != nil {
			t.Fatalf("readFile: %v", err)
		}
		if string(got) != "hello" {
			t.Fatalf("содержимое %q, ожидали hello", got)
		}
	})

	t.Run("чтение_несуществующего_ErrFileStat", func(t *testing.T) {
		_, err := readFile(filepath.Join(t.TempDir(), "нет_такого"))
		if err == nil {
			t.Fatal("ожидали ошибку")
		}
		if !errors.Is(err, ErrFileStat) {
			t.Fatalf("ошибка %v, ожидали обёртку ErrFileStat", err)
		}
	})
}

type sampleCfg struct {
	Name string `yaml:"name" json:"name"`
}

func TestDecodeConfig(t *testing.T) {
	t.Run("расширение_yaml", func(t *testing.T) {
		p := writeTemp(t, "c.yaml", []byte("name: alpha\n"))
		got, err := DecodeConfig[sampleCfg](p)
		if err != nil {
			t.Fatalf("DecodeConfig: %v", err)
		}
		if got.Name != "alpha" {
			t.Fatalf("Name=%q, ожидали alpha", got.Name)
		}
	})

	t.Run("расширение_yml", func(t *testing.T) {
		p := writeTemp(t, "c.yml", []byte("name: beta\n"))
		got, err := DecodeConfig[sampleCfg](p)
		if err != nil {
			t.Fatalf("DecodeConfig: %v", err)
		}
		if got.Name != "beta" {
			t.Fatalf("Name=%q, ожидали beta", got.Name)
		}
	})

	t.Run("расширение_json", func(t *testing.T) {
		p := writeTemp(t, "c.json", []byte(`{"name":"gamma"}`))
		got, err := DecodeConfig[sampleCfg](p)
		if err != nil {
			t.Fatalf("DecodeConfig: %v", err)
		}
		if got.Name != "gamma" {
			t.Fatalf("Name=%q, ожидали gamma", got.Name)
		}
	})

	t.Run("файл_отсутствует_ErrFileStat", func(t *testing.T) {
		_, err := DecodeConfig[sampleCfg](filepath.Join(t.TempDir(), "missing.yaml"))
		if err == nil {
			t.Fatal("ожидали ошибку")
		}
		if !errors.Is(err, ErrFileStat) {
			t.Fatalf("ошибка %v, ожидали ErrFileStat", err)
		}
	})

	t.Run("расширение_txt_ErrWrongFileExt", func(t *testing.T) {
		p := writeTemp(t, "c.txt", []byte("x"))
		_, err := DecodeConfig[sampleCfg](p)
		if err == nil {
			t.Fatal("ожидали ошибку")
		}
		if !errors.Is(err, ErrWrongFileExt) {
			t.Fatalf("ошибка %v, ожидали ErrWrongFileExt", err)
		}
	})
}
