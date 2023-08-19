package inference

import (
	"context"
	"fmt"
)

// curl inference-service-addr/generate \
//     -X POST \
//     -d '{"inputs":"What is Deep Learning?","parameters":{"max_new_tokens":20}}' \
//     -H 'Content-Type: application/json'

func (c *Client) Generate(ctx context.Context, prompt GenerateInferenceRequest) (GeneratedInferenceResponse, error) {
	// if span, ok := tracer.SpanFromContext(ctx); ok {
	// 	span.SetTag(ext.ManualDrop, true)
	// } TODO: add tracing
	fmt.Println("Token Generation requested") // delete
	const path = "/generate"
	req, err := c.NewRequest(ctx, "POST", path, prompt)
	if err != nil {
		return GeneratedInferenceResponse{}, err
	}

	var response GeneratedInferenceResponse
	// TODO: implement request headers
	if err := c.DoAndDecode(req, &response); err != nil {
		return GeneratedInferenceResponse{}, err
	}

	return response, nil
}
