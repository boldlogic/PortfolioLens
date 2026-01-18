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
	Body       []byte
}

func NewClient(cfg config.ClientConfig, log logrus.FieldLogger) *Client {
	registry := make(map[RequestType]Endpoint, len(cfg.Endpoints))

	for _, ep := range cfg.Endpoints {
		registry[RequestType(ep.Code)] = Endpoint{
			Url:            fmt.Sprintf("https://%s/%s", cfg.Host, ep.Path),
			Method:         ep.Method,
			Headers:        ep.Headers,
			RequestTimeout: ep.RequestTimeout,
			RetryPolicy:    ep.RetryPolicy,
			RetryCount:     ep.RetryCount,
		}
	}
	return &Client{
		client:          &http.Client{},
		log:             log,
		RequestRegistry: registry,
	}
}

func (c *Client) Execute(ctx context.Context, reqType string) error {

	request, err := c.PrepareRequest(ctx, reqType)
	if err != nil {
		return err
	}
	if reqType == "CBR_CURRENCIES" {
		err = c.GetCbrCurrencies(ctx, request)
		if err != nil {
			return err
		}
	}

	c.log.Infof("Сформирован запрос %s по URL %s, заголовки: ", request.Method, request.URL, request.Header)
	return nil
}

func (c *Client) sendRequest(ctx context.Context, req *http.Request) (Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("не удалось выполнить запрос: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("не удалось прочитать тело запроса: %w", err)
	}

	return Response{
		StatusCode: resp.StatusCode,
		//Headers: resp.Header,
		Body: body,
	}, nil
}

func (c *Client) PrepareRequest(ctx context.Context, reqType string) (*http.Request, error) {
	endpoint := c.RequestRegistry[RequestType(reqType)]

	if endpoint.Url == "" {
		return nil, fmt.Errorf("неизвестный тип запроса")
	}

	request, err := http.NewRequestWithContext(ctx, endpoint.Method, endpoint.Url, nil)

	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	headers := make(http.Header)
	for k, v := range endpoint.Headers {
		headers.Add(k, v)
	}
	request.Header = headers

	return request, nil
}
