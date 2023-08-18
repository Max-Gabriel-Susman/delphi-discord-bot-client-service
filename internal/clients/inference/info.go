package inference

import "context"

type InfoRequestResponse struct {
	DockerLabel           string  `json:"docker_label"`
	MaxBatchTotalTokens   int     `json:"max_batch_total_tokens"`
	MaxBestOf             int     `json:"max_best_of"`
	MaxConcurrentRequests int     `json:"max_concurrent_requests"`
	MaxInputLength        int     `json:"max_input_length"`
	MaxStopSequences      int     `json:"max_stop_sequences"`
	MaxTotalTokens        int     `json:"max_total_tokens"`
	MaxWaitingTokens      int     `json:"max_waiting_tokens"`
	ModelDeviceType       string  `json:"model_device_type"`
	ModelDType            string  `json:"model_dtype"`
	ModelID               string  `json:"model_id"`
	ModelPipelineTag      string  `json:"model_pipeline_tag"`
	ModelSHA              string  `json:"model_sha"`
	SHA                   string  `json:"sha"`
	ValidationWorkers     int     `json:"validation_workers"`
	Version               string  `json:"version"`
	WaitingServedRatio    float32 `json:"waiting_served_ratio"`
}

// curl inference-service-addr/info \
//     -X GET \
//     -H 'Content-Type: application/json'

func (c *Client) Info(ctx context.Context) (*InfoRequestResponse, error) {
	return nil, nil
}
