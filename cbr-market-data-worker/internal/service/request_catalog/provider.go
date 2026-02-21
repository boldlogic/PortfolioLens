package request_catalog

import (
	"fmt"

	"github.com/boldlogic/PortfolioLens/cbr-market-data-worker/internal/config"
)

type Provider struct {
	plans map[RequestType]RequestPlan
}

func NewProvider(cfg config.ClientConfig) *Provider {
	registry := make(map[RequestType]RequestPlan, len(cfg.Endpoints))

	for _, ep := range cfg.Endpoints {
		var params []Param
		for _, ps := range ep.QueryParams {
			params = append(params, Param{
				Name: ps.Name,
				Type: ps.Type,
			})
		}
		registry[RequestType(ep.Code)] = RequestPlan{
			Url:            fmt.Sprintf("https://%s/%s", cfg.Host, ep.Path),
			Params:         params,
			Method:         ep.Method,
			Headers:        ep.Headers,
			RequestTimeout: ep.RequestTimeout,
			RetryPolicy:    ep.RetryPolicy,
			RetryCount:     ep.RetryCount,
		}
	}
	return &Provider{
		plans: registry,
	}
}

func (r *Provider) GetPlan(reqType string) (RequestPlan, error) {
	endpoint := r.plans[RequestType(reqType)]
	if endpoint.Url == "" {
		return RequestPlan{}, fmt.Errorf("неизвестный тип запроса")
	}
	return endpoint, nil
}
