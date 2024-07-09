package mlflow

import "context"

type MetricsService service

type MetricHistoryOptions struct {
	RunID      string `json:"run_id,omitempty"`
	MetricKey  string `json:"metric_key,omitempty"`
	MaxResults int32  `json:"max_results,omitempty"`
	PageToken  string `json:"page_token,omitempty"`
}

type MetricHistory struct {
	Metrics   []*Metric `json:"metrics,omitempty"`
	NextToken string    `json:"next_page_token,omitempty"`
}

func (s *MetricsService) GetHistory(ctx context.Context, opts *MetricHistoryOptions) (*MetricHistory, error) {
	var res MetricHistory

	_, err := s.client.Do(ctx, "GET", "metrics/get-history", nil, opts, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
