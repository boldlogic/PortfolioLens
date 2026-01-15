package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/boldlogic/cbr-market-data-worker/pkg/logger"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Log    logger.Config `yaml:"log" json:"log"`
	Server ServerConfig  `yaml:"server" json:"server"`
	Client ClientConfig  `yaml:"client" json:"client"`
}

type ServerConfig struct {
	ListenHost   string `yaml:"listen_host" json:"listen_host"`
	ExternalHost string `yaml:"external_host" json:"external_host"`
	Port         int    `yaml:"port" json:"port"`
	Timeout      int    `yaml:"timeout" json:"timeout"`
}

type ClientConfig struct {
	Host      string     `yaml:"host"`
	Endpoints []Endpoint `yaml:"endpoints" json:"endpoints"`
}

type Endpoint struct {
	Code           string            `yaml:"code" json:"code"`
	Path           string            `yaml:"path" json:"path"`
	Method         string            `yaml:"method,omitempty" json:"method,omitempty"`
	Headers        map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
	RequestTimeout int               `yaml:"request_timeout,omitempty" json:"request_timeout,omitempty"`
	RetryPolicy    string            `yaml:"retry_policy,omitempty" json:"retry_policy,omitempty"`
	RetryCount     int               `yaml:"retry_count,omitempty" json:"retry_count,omitempty"`
}

const defaultConfigPath = "config.yaml"

var err error

func ParseConfig() (*Config, error) {
	configPath := flag.String("config", defaultConfigPath, "")
	flag.Parse()

	fileBody, err := os.ReadFile(*configPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл конфигурации: %w", err)
	}

	var cfg Config
	if err = yaml.Unmarshal(fileBody, &cfg); err != nil {
		return nil, fmt.Errorf("не удалось разобрать конфигурацию: %w", err)
	}

	return &cfg, nil
}

func (cl *ClientConfig) ApplyDefaults() {
	if cl.Host == "" {
		cl.Host = "www.cbr.ru"
	}
	for i := range cl.Endpoints {
		if cl.Endpoints[i].RequestTimeout <= 0 {
			cl.Endpoints[i].RequestTimeout = 20
		}
		if cl.Endpoints[i].RetryPolicy == "" {
			cl.Endpoints[i].RetryPolicy = "fixed"
		}
		if cl.Endpoints[i].RetryCount <= 0 {
			cl.Endpoints[i].RetryCount = 0
		}
	}
}

func (cl *ClientConfig) Validate() []error {
	var errs []error
	if len(cl.Endpoints) == 0 {
		errs = append(errs, fmt.Errorf("отсутствует массив 'endpoints' в 'client'"))
		return errs
	}
	for i := range cl.Endpoints {
		if cl.Endpoints[i].Code == "" {
			errs = append(errs, fmt.Errorf("в массиве 'endpoints' не заполнен 'code'"))
		}
		if cl.Endpoints[i].Path == "" {
			errs = append(errs, fmt.Errorf("в массиве 'endpoints' не заполнен 'path'"))
		}
		if cl.Endpoints[i].Method == "" {
			errs = append(errs, fmt.Errorf("в массиве 'endpoints' не заполнен 'method'"))
		}
	}

	return errs
}

func (c *Config) Validate() []error {
	var errs []error

	clErrs := c.Client.Validate()
	if len(clErrs) > 0 {
		errs = append(errs, clErrs...)
	}
	return errs
}
