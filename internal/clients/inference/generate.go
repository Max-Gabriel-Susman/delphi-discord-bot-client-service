package inference

import "context"

// curl inference-service-addr/generate \
//     -X POST \
//     -d '{"inputs":"What is Deep Learning?","parameters":{"max_new_tokens":20}}' \
//     -H 'Content-Type: application/json'

func (c *Client) Generate(ctx context.Context, req GenerateInferenceRequest) (*GeneratedInferenceResponse, error) {
	return nil, nil
}
