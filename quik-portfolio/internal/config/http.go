package config

import (
	"fmt"
)

type HttpConfig struct {
	Port    int `yaml:"port" json:"port"`
	Timeout int `yaml:"timeout" json:"timeout"`
}

func (srv *HttpConfig) validate() []error {
	var errs []error
	if srv.Port <= 0 {
		errs = append(errs, fmt.Errorf("в блоке 'http_server' некорректный 'port'"))
	}
	if srv.Timeout <= 0 {
		errs = append(errs, fmt.Errorf("в блоке 'http_server' некорректный 'timeout'"))
	}
	return errs
}
