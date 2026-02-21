package config

import (
	"errors"
	"fmt"
	"os"

	config "github.com/boldlogic/PortfolioLens/pkg/config"
	logger "github.com/boldlogic/PortfolioLens/pkg/logger/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Log    logger.Config   `yaml:"log" json:"log"`
	Server ServerConfig    `yaml:"server" json:"server"`
	Client ClientConfig    `yaml:"client" json:"client"`
	Db     config.DBConfig `yaml:"db" json:"db"`
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

	clErrs := c.Client.validate()
	if len(clErrs) > 0 {
		errs = append(errs, clErrs...)
	}

	dbErrs := c.Db.Validate()
	if len(dbErrs) > 0 {
		errs = append(errs, dbErrs...)
	}
	srvErrs := c.Server.validate()
	if len(srvErrs) > 0 {
		errs = append(errs, srvErrs...)
	}
	return errs
}

func (c *Config) applyDefaults() {
	c.Client.applyDefaults()
	c.Db.ApplyDefaults()
	c.Server.applyDefaults()
}
