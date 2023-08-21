package inference

import (
	"context"
)

type PromptRequest struct {
	Prompt string `json:"prompt"`
}

type PromptResponse struct {
	Response string `json:"response"`
}

type ErrorResponse struct { // TODO: consider movingig into delphierrors package
	Error     string `json:"error"`
	ErrorType string `json:"error_type"`
}

func (c *Client) Prompt(ctx context.Context, prompt PromptRequest) (PromptResponse, error) {
	const path = "/prompt"
	req, err := c.NewRequest(ctx, "POST", path, prompt)
	if err != nil {
		return PromptResponse{}, err
	}
	var response PromptResponse
	// TODO: implement request headers
	if err := c.DoAndDecode(req, &response); err != nil {
		return PromptResponse{}, err
	}

	return response, nil
}
