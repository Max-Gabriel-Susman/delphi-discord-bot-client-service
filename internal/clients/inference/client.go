package inference

import (
	"github.com/Max-Gabriel-Susman/delphi-go-kit/delphiclient"
)

// This client consumes the API specification documented @ https://huggingface.github.io/text-generation-inference/#/Text%20Generation%20Inference/generate

type Client struct {
	*delphiclient.Client
}

func NewClient(name, address string) *Client {
	return &Client{
		Client: delphiclient.NewClient(name, address),
	}
}
