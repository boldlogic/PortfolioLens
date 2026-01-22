package request_catalog

import (
	"fmt"

	"github.com/boldlogic/cbr-market-data-worker/internal/config"
)

type RequestType string

type RequestPlan struct {
	Url            string
	Method         string
	Headers        map[string]string
	RequestTimeout int
	RetryPolicy    string
	RetryCount     int
}

type Provider struct {
	Plans map[RequestType]RequestPlan
}

func NewProvider(cfg config.ClientConfig) *Provider {
	registry := make(map[RequestType]RequestPlan, len(cfg.Endpoints))

	for _, ep := range cfg.Endpoints {
		registry[RequestType(ep.Code)] = RequestPlan{
			Url:            fmt.Sprintf("https://%s/%s", cfg.Host, ep.Path),
			Method:         ep.Method,
			Headers:        ep.Headers,
			RequestTimeout: ep.RequestTimeout,
			RetryPolicy:    ep.RetryPolicy,
			RetryCount:     ep.RetryCount,
		}
	}
	return &Provider{
		Plans: registry,
	}
}

func (r *Provider) GetPlan(reqType string) (RequestPlan, error) {
	endpoint := r.Plans[RequestType(reqType)]
	if endpoint.Url == "" {
		return RequestPlan{}, fmt.Errorf("неизвестный тип запроса")
	}
	return endpoint, nil
}
