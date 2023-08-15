package inference

import "context"

// curl 127.0.0.1:8080/generate_stream \
//     -X POST \
//     -d '{"inputs":"What is Deep Learning?","parameters":{"max_new_tokens":20}}' \
//     -H 'Content-Type: application/json'

func (c *TextGenerationInferenceClient) GenerateStream(ctx context.Context, req *Prompt) (*Inference, error) {
	return nil, nil
}
