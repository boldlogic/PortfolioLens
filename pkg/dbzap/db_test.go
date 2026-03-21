package dbzap

import (
	"testing"
)

func TestDBConfig_ApplyDefaults(t *testing.T) {
	t.Run("Пустой DBConfig после ApplyDefaults: Host=localhost, Driver=sqlserver", func(t *testing.T) {
		db := DBConfig{}
		db.ApplyDefaults()
		if db.Host != "localhost" {
			t.Fatalf("Host=%q, ожидали localhost", db.Host)
		}
		if db.Driver != "sqlserver" {
			t.Fatalf("Driver=%q, ожидали sqlserver", db.Driver)
		}
	})

	t.Run("Host и Driver уже заданы: ApplyDefaults их не меняет", func(t *testing.T) {
		db := DBConfig{Host: "db.example", Driver: "sqlserver"}
		db.ApplyDefaults()
		if db.Host != "db.example" {
			t.Fatalf("Host=%q, ожидали db.example", db.Host)
		}
	})
}

func TestDBConfig_ApplySecretsFromEnv(t *testing.T) {
	t.Run("В env задан DB_PASSWORD: перекрывает пароль из конфига", func(t *testing.T) {
		t.Setenv("DB_PASSWORD", "from-env")
		t.Setenv("DB_USER", "")
		db := DBConfig{Password: "from-yaml"}
		db.ApplySecretsFromEnv()
		if db.Password != "from-env" {
			t.Fatalf("Password=%q, ожидали from-env", db.Password)
		}
	})

	t.Run("В env задан DB_USER: перекрывает user из конфига", func(t *testing.T) {
		t.Setenv("DB_PASSWORD", "")
		t.Setenv("DB_USER", "env-user")
		db := DBConfig{User: "yaml-user"}
		db.ApplySecretsFromEnv()
		if db.User != "env-user" {
			t.Fatalf("User=%q, ожидали env-user", db.User)
		}
	})

	t.Run("пустые переменные не меняют конфиг", func(t *testing.T) {
		t.Setenv("DB_PASSWORD", "")
		t.Setenv("DB_USER", "")
		db := DBConfig{User: "u", Password: "p"}
		db.ApplySecretsFromEnv()
		if db.User != "u" || db.Password != "p" {
			t.Fatalf("User=%q Password=%q, ожидали u и p", db.User, db.Password)
		}
	})
}

func TestDBConfig_Validate(t *testing.T) {
	t.Run("Все обязательные поля и driver sqlserver: Validate без ошибок", func(t *testing.T) {
		db := DBConfig{
			Server:   "host",
			Name:     "db",
			User:     "u",
			Password: "p",
			Driver:   "sqlserver",
		}
		errs := db.Validate()
		if len(errs) != 0 {
			t.Fatalf("ожидали 0 ошибок, получили %v", errs)
		}
	})

	t.Run("Нулевой DBConfig: ровно пять ошибок валидации", func(t *testing.T) {
		var db DBConfig
		errs := db.Validate()
		if len(errs) != 5 {
			t.Fatalf("ожидали 5 ошибок, получили %d: %v", len(errs), errs)
		}
	})

	t.Run("Driver не sqlserver при остальных полях ок: одна ошибка валидации", func(t *testing.T) {
		db := DBConfig{
			Server:   "host",
			Name:     "db",
			User:     "u",
			Password: "p",
			Driver:   "postgres",
		}
		errs := db.Validate()
		if len(errs) != 1 {
			t.Fatalf("ожидали 1 ошибку, получили %v", errs)
		}
	})
}

func TestDBConfig_ApplyDefaults_then_Validate(t *testing.T) {
	t.Run("Только обязательные поля без Driver: после ApplyDefaults Validate без ошибок", func(t *testing.T) {
		db := DBConfig{
			Server:   "srv",
			Name:     "n",
			User:     "u",
			Password: "p",
		}
		db.ApplyDefaults()
		errs := db.Validate()
		if len(errs) != 0 {
			t.Fatalf("ожидали 0 ошибок, получили %v", errs)
		}
	})
}
