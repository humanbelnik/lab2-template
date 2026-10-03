package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	requestTimeout = 5 * time.Second
	maxRetries     = 3
	retryDelay     = 300 * time.Millisecond
)

type HTTPClient struct {
	httpClient *http.Client
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		httpClient: &http.Client{
			Timeout: requestTimeout,
		},
	}
}

func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	var requestBody []byte
	if req.Body != nil {
		var err error
		requestBody, err = io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {
			return nil, err
		}
	}

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if requestBody != nil {
			req.Body = io.NopCloser(bytes.NewReader(requestBody))
			req.ContentLength = int64(len(requestBody))
		}

		resp, err := c.httpClient.Do(req)
		if err == nil && resp.StatusCode < http.StatusInternalServerError {
			return resp, nil
		}

		if err != nil {
			lastErr = err
		} else {
			lastErr = fmt.Errorf("service responded with status %d", resp.StatusCode)
			resp.Body.Close()
		}

		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	return nil, fmt.Errorf("request to %s failed after %d attempts: %w", req.URL.String(), maxRetries, lastErr)
}
