package requestplanner

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func PrepareRequest(ctx context.Context, plan RequestPlan) (*http.Request, error) {
	rawURL := strings.TrimSpace(plan.Url)
	if rawURL == "" {
		return nil, fmt.Errorf("URL плана пуст")
	}
	method := strings.TrimSpace(strings.ToUpper(plan.Method))
	if method == "" {
		method = http.MethodGet
	}

	reqURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка разбора URL '%s': %w", rawURL, err)
	}

	query, headers, err := fillParams(plan.Params)
	if err != nil {
		return nil, err
	}

	if len(query) > 0 {
		q := reqURL.Query()
		for k, v := range query {
			q.Set(k, v)
		}
		reqURL.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req, nil
}
