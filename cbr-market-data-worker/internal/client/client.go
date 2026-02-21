package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/boldlogic/PortfolioLens/cbr-market-data-worker/internal/config"
	"github.com/boldlogic/PortfolioLens/cbr-market-data-worker/internal/service/request_catalog"
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
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 YaBrowser/25.12.0.0 Safari/537.36")
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

	for k, v := range endpoint.Headers {
		request.Header.Set(k, v)
	}

	return request, nil
}

func (c *Client) PrepareRequestWithParams(ctx context.Context, endpoint request_catalog.RequestPlan, reqParams map[string]string) (*http.Request, error) {

	reqURL, err := url.Parse(endpoint.Url)
	if err != nil {
		return nil, err
	}

	query := reqURL.Query()
	for key, value := range reqParams {
		query.Set(key, value)
	}
	reqURL.RawQuery = query.Encode()

	request, err := http.NewRequestWithContext(ctx, endpoint.Method, reqURL.String(), nil)

	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	for k, v := range endpoint.Headers {
		request.Header.Set(k, v)
	}

	return request, nil
}
