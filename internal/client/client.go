package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/boldlogic/cbr-market-data-worker/internal/config"
	"github.com/boldlogic/cbr-market-data-worker/internal/service/request_catalog"
)

type Client struct {
	Client *http.Client
}

func NewClient(cfg config.ClientConfig) *Client {
	tr := &http.Transport{
		MaxIdleConnsPerHost: 1,
		IdleConnTimeout:     30 * time.Second,
	}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(cfg.Timeout) * time.Second,
	}

	return &Client{
		Client: httpClient,
	}
}

type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

func (c *Client) SendRequest(ctx context.Context, req *http.Request) (Response, error) {
	resp, err := c.Client.Do(req)
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

func (c *Client) PrepareRequest(ctx context.Context, endpoint request_catalog.RequestPlan) (*http.Request, error) {

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
