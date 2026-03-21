package client

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/boldlogic/PortfolioLens/pkg/transport/httpclient/clientmetrics"
	"go.uber.org/zap"
)

type CommonClient interface {
	SendRequest(ctx context.Context, req *http.Request) (int, []byte, error)
	SendWithRetry(ctx context.Context, req *http.Request, retryCount int) (int, []byte, int, error)
}

type Client struct {
	commonClient CommonClient
	metrics      *clientmetrics.ClientMetrics
	target       string
	logger       *zap.Logger
}

func NewClient(commonClient CommonClient, m *clientmetrics.ClientMetrics, target string, logger *zap.Logger) *Client {
	return &Client{
		commonClient: commonClient,
		metrics:      m,
		target:       target,
		logger:       logger,
	}
}

func (c Client) SendRequest(ctx context.Context, req *http.Request) (int, []byte, error) {
	start := time.Now()
	code, body, err := c.commonClient.SendRequest(ctx, req)

	status := strconv.Itoa(code)
	if err != nil {
		status = "network_error"
	}
	c.metrics.RecordRequest(req.Method, c.target, req.URL.Path, status, time.Since(start))

	if err == nil {
		c.logResponse(req.URL.String(), code, 1)
	}
	return code, body, err
}

func (c Client) SendWithRetry(ctx context.Context, req *http.Request, retryCount int) (int, []byte, int, error) {
	start := time.Now()
	code, body, attempts, err := c.commonClient.SendWithRetry(ctx, req, retryCount)

	status := strconv.Itoa(code)
	if err != nil {
		status = "network_error"
	}
	c.metrics.RecordRequest(req.Method, c.target, req.URL.Path, status, time.Since(start))

	c.logResponse(req.URL.String(), code, attempts)
	return code, body, attempts, err
}

func (c Client) logResponse(url string, status int, attempts int) {
	fields := []zap.Field{zap.String("url", url), zap.Int("status", status), zap.Int("attempts", attempts)}
	if status != http.StatusOK {
		c.logger.Warn("HTTP ответ не 200", fields...)
		return
	}
	c.logger.Debug("HTTP запрос выполнен", fields...)
}
