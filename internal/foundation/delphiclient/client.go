package delphiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/Max-Gabriel-Susman/delphi-discord-bot-client-service/internal/foundation/delphierror"
)

var (
	ErrReqTimeout  = errors.New("outgoing request timed out")
	ErrReqCanceled = errors.New("outgoing request's context was canceled")
	ErrNotFound    = errors.New("not found")
)

type Client struct {
	BaseURL string
	*http.Client
	ResourceName string
}

func NewClient(name, address string) *Client {
	c := &Client{
		BaseURL:      address,
		Client:       &http.Client{Timeout: 60 * time.Second},
		ResourceName: name,
	}

	// c.Client.Transport = delphitrace.TracedTransport(name)
	return c
}

func (c *Client) FullURL(path string) string {
	return c.BaseURL + path
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.Client.Do(req)
	if err == nil {
		return resp, nil
	}
	if errors.Is(err, context.DeadlineExceeded) {
		// HTTP 499 in Nginx means that the client closed the connection before the server answered the request
		return nil, delphierror.WithStatusCode(ErrReqTimeout, http.StatusGatewayTimeout)
	}
	if errors.Is(err, context.Canceled) {
		// HTTP 499 in Nginx means that the client closed the connection before the server answered the request
		return nil, delphierror.WithStatusCode(ErrReqCanceled, 499)
	}
	return nil, err
}

func (c *Client) NewRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	var buf io.Reader
	if body != nil {
		switch b := body.(type) {
		case io.Reader:
			buf = b
		default:
			var tmp bytes.Buffer
			if err := json.NewEncoder(&tmp).Encode(body); err != nil {
				return nil, err
			}
			buf = &tmp
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, c.FullURL(path), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")
	return req, nil
}
