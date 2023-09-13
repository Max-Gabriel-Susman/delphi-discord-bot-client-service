package application

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) HealthCheck(ctx context.Context) (*ErrorResponse, error) {
	// if span, ok := tracer.SpanFromContext(ctx); ok {
	// 	span.SetTag(ext.ManualDrop, true)
	// } TODO: add tracing
	fmt.Println("Healthcheck requested") // delete
	const path = "/health"

	req, err := c.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("could not construct helathcheck request. Address: %s. Err: %w", c.FullURL(path), err)
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send healthcheck request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected healthcheck status code response: %d", resp.StatusCode)
	}
	fmt.Println("Healthcheck completed") // delete
	return nil, nil
}