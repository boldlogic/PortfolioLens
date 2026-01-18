package config

import "fmt"

type DBConfig struct {
	Driver   string `yaml:"driver,omitempty" json:"driver,omitempty"`
	Server   string `yaml:"server,omitempty" json:"server,omitempty"`
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Name     string `yaml:"db_name" json:"db_name"`
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password,omitempty" json:"password,omitempty"`
	SSLMode  string `yaml:"ssl_mode,omitempty" json:"ssl_mode,omitempty"`
}

func (db *DBConfig) applyDefaults() {
	if db.Host == "" {
		db.Host = "localhost"
	}
	if db.Driver == "" {
		db.Driver = "sqlserver"
	}

}

func (db *DBConfig) validate() []error {
	var errs []error
	if db.Name == "" {
		errs = append(errs, fmt.Errorf("в блоке 'db' не заполнено 'db_name'"))
	}
	if db.User == "" {
		errs = append(errs, fmt.Errorf("в блоке 'db' не заполнено 'user'"))
	}
	if db.Password == "" {
		errs = append(errs, fmt.Errorf("в блоке 'db' не заполнено 'password'"))
	}
	if db.Server == "" {
		errs = append(errs, fmt.Errorf("в блоке 'db' не заполнено 'server'"))
	}
	if db.Driver != "sqlserver" {
		errs = append(errs, fmt.Errorf("в блоке 'db' неподдерживаемый 'driver'"))
	}

	return errs
}

func (db *DBConfig) GetDSN() (dsn string) {
	//if db.Driver == "sqlserver" {
	return fmt.Sprintf("sqlserver://%s:%s@%s?database=%s&trustServerCertificate=true",
		db.User, db.Password, db.Server, db.Name)

}
