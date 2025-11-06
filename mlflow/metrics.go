package mlflow

import "context"

// MetricsService handles communication with the metrics related methods of the MLflow API.
type MetricsService service

// MetricHistoryOptions represents the parameters for retrieving metric history
type MetricHistoryOptions struct {
	// RunID is the unique identifier for the run
	RunID string `json:"run_id,omitempty"`
	// MetricKey is the name of the metric to retrieve
	MetricKey string `json:"metric_key,omitempty"`
	// MaxResults is the maximum number of metrics to return
	MaxResults int32 `json:"max_results,omitempty"`
	// PageToken is used for pagination to fetch the next page of results
	PageToken string `json:"page_token,omitempty"`
}

// MetricHistory represents the history of a metric over time
type MetricHistory struct {
	// Metrics is the list of metric values recorded over time
	Metrics []*Metric `json:"metrics,omitempty"`
	// NextPageToken is used to retrieve the next page of results (empty if no more results)
	NextPageToken string `json:"next_page_token,omitempty"`
}

// GetHistory retrieves the history of values for a specific metric
//
// Returns all logged values of the metric in chronological order.
// Use pagination parameters for metrics with many data points.
func (s *MetricsService) GetHistory(ctx context.Context, opts *MetricHistoryOptions) (*MetricHistory, error) {
	var res MetricHistory

	_, err := s.client.Do(ctx, "GET", "metrics/get-history", nil, opts, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
