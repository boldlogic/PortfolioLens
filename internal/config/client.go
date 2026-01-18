package config

import "fmt"

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

func (cl *ClientConfig) applyDefaults() {
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

func (cl *ClientConfig) validate() []error {
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
