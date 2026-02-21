package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/boldlogic/PortfolioLens/pkg/config"
	logger "github.com/boldlogic/PortfolioLens/pkg/logger/zap"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Log  logger.Config   `yaml:"log" json:"log"`
	Db   config.DBConfig `yaml:"db" json:"db"`
	Http HttpConfig      `yaml:"http_server" json:"http_server"`
}

func LoadConfig(configPath string) (*Config, error) {

	fileBody, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл конфигурации: %w", err)
	}

	var cfg Config
	if err = yaml.Unmarshal(fileBody, &cfg); err != nil {
		return nil, fmt.Errorf("не удалось разобрать конфигурацию: %w", err)
	}
	cfg.applyDefaults()
	errs := cfg.validate()
	if err := errors.Join(errs...); err != nil {
		return nil, fmt.Errorf("некорректный конфиг: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() []error {
	var errs []error

	dbErrs := c.Db.Validate()
	if len(dbErrs) > 0 {
		errs = append(errs, dbErrs...)
	}

	HttpErrs := c.Http.validate()
	if len(dbErrs) > 0 {
		errs = append(errs, HttpErrs...)
	}
	return errs
}

func (c *Config) applyDefaults() {

	c.Db.ApplyDefaults()
}
