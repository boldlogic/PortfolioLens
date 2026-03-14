package httpserver

import "fmt"

type ServerConfig struct {
	ListenHost   string `yaml:"listen_host" json:"listen_host"`
	ExternalHost string `yaml:"external_host" json:"external_host"`
	Port         int    `yaml:"port" json:"port"`
	Timeout      int    `yaml:"timeout" json:"timeout"`
}

func (srv *ServerConfig) ApplyDefaults() {
	if srv.ListenHost == "" {
		srv.ListenHost = "127.0.0.1"
	}
	if srv.ExternalHost == "" {
		srv.ExternalHost = "localhost"
	}
	if srv.Port == 0 {
		srv.Port = 80
	}
	if srv.Timeout == 0 {
		srv.Timeout = 60
	}

}

func (srv *ServerConfig) Validate() []error {
	var errs []error
	if srv.Port < 1 || srv.Port > 65535 {
		errs = append(errs, fmt.Errorf("в блоке 'server' некорректный 'port': должен быть в диапазоне 1-65535, получено %d", srv.Port))
	}
	return errs
}
