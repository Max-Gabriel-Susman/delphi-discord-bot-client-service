package inference

import "context"

type MetricsRequestResponse struct {
	MetricsText string `json:"text"`
}

// curl inference-service-addr/info \
//     -X GET \
//     -H 'Content-Type: application/json'

func (c *Client) Metrics(ctx context.Context) (*MetricsRequestResponse, error) {
	return nil, nil
}
