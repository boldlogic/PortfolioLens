package service

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/boldlogic/cbr-market-data-worker/internal/config"
	"github.com/sirupsen/logrus"
)

type Client struct {
	client          *http.Client
	log             logrus.FieldLogger
	RequestRegistry map[RequestType]Endpoint
}

type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       interface{}
}

func NewClient(cfg config.ClientConfig, log logrus.FieldLogger) *Client {
	registry := make(map[RequestType]Endpoint, len(cfg.Endpoints))

	for _, ep := range cfg.Endpoints {
		registry[RequestType(ep.Code)] = Endpoint{
			url:            fmt.Sprintf("https://%s/%s", cfg.Host, ep.Path),
			method:         ep.Method,
			headers:        ep.Headers,
			requestTimeout: ep.RequestTimeout,
			retryPolicy:    ep.RetryPolicy,
			retryCount:     ep.RetryCount,
		}
	}
	return &Client{
		client:          &http.Client{},
		log:             log,
		RequestRegistry: registry,
	}
}

// func (c *Client) Get(ctx context.Context, url string, headers http.Header) (Response, error) {

// }

func (c *Client) ExecRequest(ctx context.Context, reqType string) (Response, error) {
	//h.Service.RequestRegistry[service.RequestType(cb.Type)])
	req := c.RequestRegistry[RequestType(reqType)]

	if req.url == "" {
		return Response{}, fmt.Errorf("Неизвестный тип запроса")
	}

	request, err := http.NewRequestWithContext(ctx, req.method, req.url, nil)

	if err != nil {
		return Response{}, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	headers := make(http.Header)
	for k, v := range req.headers {
		headers.Add(k, v)
		//request.Header[k]=v
	}
	resp, err := c.client.Do(request)
	if err != nil {
		return Response{}, fmt.Errorf("can't do request: %w", err)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("can't read response body: %w", err)
	}

	if err = resp.Body.Close(); err != nil {
		return Response{}, fmt.Errorf("can't close response body: %w", err)
	}

	c.log.Info("запрос на выполнение", req.url, req.headers, string(respBody))
	return Response{
		StatusCode: resp.StatusCode,
	}, nil
}
