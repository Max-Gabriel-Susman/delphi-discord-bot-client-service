package application

import (
	"context"
	"fmt"
)

type ScanRequest struct {
	Request string `json:"request"`
}

type ScanResponse struct {
	Results string `json:"results"`
}

func (c *Client) Scan(ctx context.Context, scan ScanRequest) (ScanResponse, error) {
	const path = "/scan"
	req, err := c.NewRequest(ctx, "POST", path, scan)
	if err != nil {
		return ScanResponse{}, err
	}
	var response ScanResponse
	// TODO: implement request headers
	if err := c.DoAndDecode(req, &response); err != nil {
		return ScanResponse{}, err
	}
	fmt.Println("Scan completed") // delete l8r
	return response, nil
}
