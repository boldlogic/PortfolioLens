package service

import (
	"net/http"

	"github.com/boldlogic/cbr-market-data-worker/internal/config"
	"github.com/sirupsen/logrus"
)

type Client struct {
	client *http.Client
	log    logrus.FieldLogger
}

type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       interface{}
}

func NewClient(cfg config.ClientConfig, log logrus.FieldLogger) *Client {
	return &Client{
		client: &http.Client{},
		log:    log,
	}
}

// func (c *Client) Get(ctx context.Context, url string, headers http.Header) (Response, error) {

// }
