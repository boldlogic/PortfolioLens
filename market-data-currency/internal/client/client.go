package client

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

type CommonClient interface {
	SendRequest(ctx context.Context, req *http.Request) (int, []byte, error)
	SendWithRetry(ctx context.Context, req *http.Request, retryCount int) (int, []byte, int, error)
}

type Client struct {
	commonClient CommonClient
	logger       *zap.Logger
}

func NewClient(commonClient CommonClient, logger *zap.Logger) *Client {
	return &Client{commonClient: commonClient, logger: logger}
}

func (c Client) SendRequest(ctx context.Context, req *http.Request) (int, []byte, error) {
	code, body, err := c.commonClient.SendRequest(ctx, req)
	if err == nil {
		c.logResponse(req.URL.String(), code, 1)
	}
	return code, body, err
}

func (c Client) SendWithRetry(ctx context.Context, req *http.Request, retryCount int) (int, []byte, int, error) {
	code, body, attempts, err := c.commonClient.SendWithRetry(ctx, req, retryCount)
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
