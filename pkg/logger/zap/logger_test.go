package logger

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_writeSyncers(t *testing.T) {
	t.Run("output_file пуст: один WriteSyncer (консоль)", func(t *testing.T) {
		ss := writeSyncers(Config{OutputFile: ""})
		if len(ss) != 1 {
			t.Fatalf("ожидали 1 syncer, получили %d", len(ss))
		}
	})

	t.Run("output_file открывается: один WriteSyncer только в файл (без дубля в stdout)", func(t *testing.T) {
		ss := writeSyncers(Config{OutputFile: os.DevNull})
		if len(ss) != 1 {
			t.Fatalf("ожидали 1 syncer, получили %d", len(ss))
		}
	})

	t.Run("output_file недоступен: fallback один WriteSyncer в stdout", func(t *testing.T) {
		badDir := filepath.Join(t.TempDir(), "нет_каталога", "x.log")
		ss := writeSyncers(Config{OutputFile: badDir})
		if len(ss) != 1 {
			t.Fatalf("ожидали 1 syncer, получили %d", len(ss))
		}
	})
}

func TestNew_noPanic(t *testing.T) {
	t.Run("output_file пуст: New возвращает logger, Sync без паники", func(t *testing.T) {
		log := New(Config{Level: "info", Format: "console", OutputFile: ""})
		if log == nil {
			t.Fatal("ожидали logger")
		}
		_ = log.Sync()
	})

	t.Run("output_file открывается (DevNull): New возвращает logger, Info и Sync без паники", func(t *testing.T) {
		log := New(Config{Level: "info", Format: "json", OutputFile: os.DevNull})
		if log == nil {
			t.Fatal("ожидали logger")
		}
		log.Info("test")
		_ = log.Sync()
	})
}
