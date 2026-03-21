package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempYAML(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "cfg.yaml")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoad_minimal_yaml_без_ошибки(t *testing.T) {
	p := writeTempYAML(t, `
db:
  db_name: fxdb
  user: u
  password: p
  server: sql.local
server:
  port: 3000
`)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg == nil {
		t.Fatal("ожидали конфиг")
	}
	if cfg.Db.Name != "fxdb" {
		t.Fatalf("Db.Name=%q", cfg.Db.Name)
	}
	if cfg.Db.Host != "localhost" {
		t.Fatalf("после ApplyDefaults Host=%q, ожидали localhost", cfg.Db.Host)
	}
	if cfg.Server.Port != 3000 {
		t.Fatalf("Server.Port=%d", cfg.Server.Port)
	}
	if cfg.Client.Timeout != 60 {
		t.Fatalf("Client.Timeout=%d, ожидали дефолт httpclient 60", cfg.Client.Timeout)
	}
}

func TestLoad_validate_обёртка_некорректный_конфиг(t *testing.T) {
	p := writeTempYAML(t, `
db:
  user: u
  password: p
  server: sql.local
server:
  port: 3000
`)
	_, err := Load(p)
	if err == nil {
		t.Fatal("ожидали ошибку валидации")
	}
	if !strings.Contains(err.Error(), "некорректный конфиг") {
		t.Fatalf("ожидали префикс обёртки, получили: %v", err)
	}
}
