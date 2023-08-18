package inference

import "context"

type MetricsRequestResponse struct {
	MetricsText string `json:"text"`
}

func (c *Client) Metrics(ctx context.Context) (*MetricsRequestResponse, error) {
	return nil, nil
}
