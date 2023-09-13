package application

import "github.com/Max-Gabriel-Susman/delphi-go-kit/delphiclient"

type Client struct {
	*delphiclient.Client
}

func NewClient(name, address string) *Client {
	return &Client{
		Client: delphiclient.NewClient(name, address),
	}
}

type ErrorResponse struct { // TODO: consider movingig into delphierrors package
	Error     string `json:"error"`
	ErrorType string `json:"error_type"`
}
