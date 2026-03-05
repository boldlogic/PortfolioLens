package config

import "fmt"

type ClientConfig struct {
	Host      string     `yaml:"host"`
	Timeout   int        `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Endpoints []Endpoint `yaml:"endpoints" json:"endpoints"`
}

type QueryParam struct {
	Name string `yaml:"name"`
	Type string `yaml:"type" json:"type"`
}

type Endpoint struct {
	Code           string            `yaml:"code" json:"code"`
	Path           string            `yaml:"path" json:"path"`
	Method         string            `yaml:"method,omitempty" json:"method,omitempty"`
	QueryParams    []QueryParam      `yaml:"query_params,omitempty" json:"query_params,omitempty"`
	Headers        map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
	RequestTimeout int               `yaml:"request_timeout,omitempty" json:"request_timeout,omitempty"`
	RetryPolicy    string            `yaml:"retry_policy,omitempty" json:"retry_policy,omitempty"`
	RetryCount     int               `yaml:"retry_count,omitempty" json:"retry_count,omitempty"`
}

func (cl *ClientConfig) applyDefaults() {
	if cl.Host == "" {
		cl.Host = "www.cbr.ru"
	}
	if cl.Timeout == 0 {
		cl.Timeout = 60
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
		epErrs := cl.Endpoints[i].validate()
		if len(epErrs) > 0 {
			errs = append(errs, epErrs...)
		}
	}

	return errs
}

func (ep *Endpoint) validate() []error {
	var errs []error

	if ep.Code == "" {
		errs = append(errs, fmt.Errorf("в массиве 'endpoints' не заполнен 'code'"))

	}
	if ep.Path == "" {
		errs = append(errs, fmt.Errorf("в массиве 'endpoints' не заполнен 'path'"))

	}
	if ep.Method == "" {
		errs = append(errs, fmt.Errorf("в массиве 'endpoints' не заполнен 'method'"))

	}
	for i := range ep.QueryParams {
		qpErrs := ep.QueryParams[i].validate()
		if len(qpErrs) > 0 {
			errs = append(errs, qpErrs...)
		}
	}
	return errs
}

func (p *QueryParam) validate() []error {
	var errs []error

	if p.Name == "" {
		errs = append(errs, fmt.Errorf("в массиве 'query_params' не заполнен 'name'"))
	}
	if p.Type == "" {
		errs = append(errs, fmt.Errorf("в массиве 'query_params' не заполнен 'type'"))
	}
	return errs
}
